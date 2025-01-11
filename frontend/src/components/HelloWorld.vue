<script setup>
import { reactive } from 'vue'
import { ScrapeSITC, SelectDownloadDir } from '../../wailsjs/go/main/App'
import {EventsOn, EventsOff} from "../../wailsjs/runtime/runtime";

const data = reactive({
  username: "",
  password: "",
  downloadDir: "sitc_pieces",
  resultText: "",
  percentDone: 0,
  scraping: false,
})

function scrape() {
  data.resultText = "Scraping...";
  data.scraping = true;
  let cancelEventId = EventsOn("downloadProgress", status => {
    data.resultText = `Downloaded pieces for instrument ${status.InstrumentName}. Percent done: ${status.PercentProgress}`
    data.percentDone = status.PercentProgress;
  })
  console.log(cancelEventId)
  ScrapeSITC(data.username, data.password, data.downloadDir).then(result => {
    EventsOff("downloadProgress")
    data.resultText = "Downloaded everything :-)"
  })
}

function selectDownloadDir() {
  SelectDownloadDir(data.downloadDir).then(result => {
    data.downloadDir = result;
  })
}

</script>

<template>
  <main>
    <div class="row">
      <div class="col">
        <p>Welcome to the Summer in the City Portal scraper. This program will download all available pieces from the
          portal to your computer.</p>
        <p>Please enter your portal login information below and click the "Download!" button to start the download.</p>
      </div>
    </div>
    <div class="row">
      <div class="col">
        <label for="downloadDir">Download directory</label>
        <p class="grouped">
          <input id="downloadDir" v-model="data.downloadDir" type="text" disabled placeholder="Download directory">
          <button class="button" @click="selectDownloadDir">Select</button>
        </p>
        <p>
          <label for="username">Username</label>
          <input id="username" v-model="data.username" autocomplete="off" class="input" type="text"
            placeholder="Username" />
        </p>
        <p>
          <label for="password">Password</label>
          <input id="password" v-model="data.password" autocomplete="off" class="input" type="password"
            placeholder="Password" />
        </p>
        <button class="button" @click="scrape">Download!</button>
      </div>
    </div>
    <p></p>
    <p></p>
    <template v-if="data.scraping">
      <div class="row">
        <div class="col">
          <meter max="100" min="0" :value="data.percentDone"></meter>
        </div>
      </div>
      <div class="row">
        <div class="col">
          {{ data.resultText }}
        </div>
      </div>
    </template>
  </main>
</template>

<style scoped>
main {
  padding: 100px;
}
meter {
  width: 100%;
}
</style>
