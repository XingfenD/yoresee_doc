<template>
  <div class="app-container">
    <router-view />
  </div>
</template>

<script setup>
// App.vue 作为路由容器，负责渲染路由组件
</script>

<style>

* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

body {
  font-family: Inter, PingFang SC, Roboto, -apple-system, BlinkMacSystemFont, 'Segoe UI', Arial, sans-serif;
  font-size: 14px;
  line-height: 1.5;
  color: #4E5969;
  background-color: #F7F8FA;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}

::-webkit-scrollbar {
  width: 6px;
  height: 6px;
}

::-webkit-scrollbar-track {
  background: #F1F1F1;
  border-radius: 3px;
}

::-webkit-scrollbar-thumb {
  background: #C9CDD4;
  border-radius: 3px;
}

::-webkit-scrollbar-thumb:hover {
  background: #86909C;
}

/* 应用容器 */
.app-container {
  width: 100%;
  height: 100vh;
  overflow: hidden;
}
</style>