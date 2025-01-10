<script setup>
import {reactive} from 'vue'
import {ScrapeSITC} from '../../wailsjs/go/main/App'

const data = reactive({
  username: "",
  password: "",
  downloadDir: "sitc_pieces",
  resultText: "Please enter your credentials below 👇",
})

function scrape() {
  data.resultText = "Scraping..."
  ScrapeSITC(data.username, data.password, data.downloadDir).then(result => {
    data.resultText = "Scraped everything :-)"
  })
}

</script>

<template>
  <main>
    <div id="result" class="result">{{ data.resultText }}</div>
    <div id="input" class="input-box">
      <input id="username" v-model="data.username" autocomplete="off" class="input" type="text" placeholder="Username"/>
      <input id="password" v-model="data.password" autocomplete="off" class="input" type="password" placeholder="Password"/>
      <button class="btn" @click="scrape">Scrape!</button>
    </div>
  </main>
</template>

<style scoped>
.result {
  height: 20px;
  line-height: 20px;
  margin: 1.5rem auto;
}

.input-box .btn {
  width: 60px;
  height: 30px;
  line-height: 30px;
  border-radius: 3px;
  border: none;
  margin: 0 0 0 20px;
  padding: 0 8px;
  cursor: pointer;
}

.input-box .btn:hover {
  background-image: linear-gradient(to top, #cfd9df 0%, #e2ebf0 100%);
  color: #333333;
}

.input-box .input {
  border: none;
  border-radius: 3px;
  outline: none;
  height: 30px;
  line-height: 30px;
  padding: 0 10px;
  background-color: rgba(240, 240, 240, 1);
  -webkit-font-smoothing: antialiased;
}

.input-box .input:hover {
  border: none;
  background-color: rgba(255, 255, 255, 1);
}

.input-box .input:focus {
  border: none;
  background-color: rgba(255, 255, 255, 1);
}
</style>
