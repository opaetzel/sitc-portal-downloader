<script setup>
import HelloWorld from './components/HelloWorld.vue'</script>

<template>
  <img id="logo" alt="Wails logo" src="./assets/images/sitc_logo_lila.png"/>
  <HelloWorld/>
</template>

<style>
#logo {
  display: block;
  width: 50%;
  position: fixed;
  max-height: 70vh;
  object-fit: contain;
  z-index: 5;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  opacity: 0.1;
  pointer-events:none;
}
</style>
