package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"path"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/alitto/pond/v2"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

var baseUri = "https://portal.summer-hannover.de"
var mtx sync.Mutex

// App struct
type App struct {
	ctx                   context.Context
	client                *http.Client
	pool                  pond.Pool
	totalInstruments      float64
	downloadedInstruments float64
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

func (a *App) SelectDownloadDir(currentSelection string) string {
	dir, _ := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{})
	if dir != "" {
		return dir
	}
	return currentSelection
}

func (a *App) ScrapeSITC(username, password, downloadDir string) {
	a.pool = pond.NewPool(4)
	jar, err := cookiejar.New(nil)
	if err != nil {
		panic(err)
	}

	a.client = &http.Client{
		Jar: jar,
	}
	postUser := []string{username}
	postPass := []string{password}
	_, err = a.client.PostForm(baseUri+"/index.php/site/login", url.Values{"LoginForm[username]": postUser, "LoginForm[password]": postPass})
	fmt.Println(err)

	respGet, _ := a.client.Get(baseUri + "/index.php/listmusic/showi/i_id/-")
	doc, err := goquery.NewDocumentFromReader(respGet.Body)
	if err != nil {
		panic(err)
	}
	instumentInfos := make([]SelectInstrumentInfo, 0)
	doc.Find("select option").Each(func(i int, s *goquery.Selection) {
		value, _ := s.Attr("value")
		name := s.Text()
		instumentInfos = append(instumentInfos, SelectInstrumentInfo{name, value})
	})
	a.totalInstruments = float64(len(instumentInfos))
	a.downloadedInstruments = 0
	//fmt.Println(instumentInfos)

	for _, instrumentInfo := range instumentInfos {
		fmt.Println(instrumentInfo)
		a.pool.Submit(func() {
			a.scrapeInstrumentPage(instrumentInfo, downloadDir)
		})
	}
	a.pool.StopAndWait()
}

type SelectInstrumentInfo struct {
	name  string
	value string
}

type PieceLink struct {
	name       string
	instrument string
	href       string
}

type ProgressInfo struct {
	InstrumentName  string
	PercentProgress int
}

func (a *App) scrapeInstrumentPage(instrumentInfo SelectInstrumentInfo, downloadDir string) error {
	fmt.Println("Getting pieces for instrument", instrumentInfo.name)
	respGet, _ := a.client.Get(baseUri + "/index.php/listmusic/showi/i_id/" + instrumentInfo.value)
	doc, err := goquery.NewDocumentFromReader(respGet.Body)
	if err != nil {
		panic(err)
	}

	doc.Find("table tbody tr").Each(func(count int, s *goquery.Selection) {
		pieceLink := PieceLink{}
		for i, s := range s.Children().EachIter() {
			if i == 0 {
				pieceLink.name = s.Text()
			}
			if i == 1 {
				pieceLink.instrument = s.Text()
			}
			if i == 4 {
				pieceLink.href = s.Find("a").First().AttrOr("href", "")
			}
		}
		//pool.Submit(func() {
		a.downloadPdf(downloadDir, pieceLink)
		//})
	})
	mtx.Lock()
	a.downloadedInstruments++
	runtime.EventsEmit(a.ctx, "downloadProgress", ProgressInfo{InstrumentName: instrumentInfo.name, PercentProgress: int((a.downloadedInstruments / a.totalInstruments * 100))})
	mtx.Unlock()
	return nil
}

func (a *App) downloadPdf(downloadDir string, link PieceLink) {
	if !strings.HasSuffix(link.href, ".pdf") {
		//fmt.Println("Not downloading", link.name)
		return
	}

	dirNameWithInstrument := path.Join(downloadDir, link.instrument)
	os.MkdirAll(dirNameWithInstrument, 0755)

	filename := path.Join(dirNameWithInstrument, fmt.Sprintf("%s %s (SITC).pdf", link.name, link.instrument))
	f, err := os.Create(filename)
	if err != nil {
		fmt.Printf("error creating file '%s': %s\n", filename, err.Error())
	}
	pdfResp, err := a.client.Get(baseUri + link.href)
	if err != nil {
		fmt.Printf("error downloading file %s: %s", link.href, err.Error())
	}
	io.Copy(f, pdfResp.Body)
}
