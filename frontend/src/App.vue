<template>
  <RouterView />
</template>

<script setup lang="ts">
import { RouterView } from 'vue-router'
</script>

<style>
/* 
  你可以在这里放置一些真正意义上的全局样式，
  它会影响到应用中的所有组件。
  例如，重置一些盒模型行为。
*/
html, body, #app {
  height: 100%;
  margin: 0;
  padding: 0;
}
</style>