<template>
  <div class="app-container">
    <router-view />
  </div>
</template>

<script setup>
import { onMounted, watch } from 'vue';
import { useUserStore } from './store/user';

const userStore = useUserStore();

// 初始化暗黑模式
const initDarkMode = () => {
  if (userStore.darkMode) {
    document.documentElement.classList.add('dark-mode');
  } else {
    document.documentElement.classList.remove('dark-mode');
  }
};

// 监听暗黑模式变化
watch(
  () => userStore.darkMode,
  (newValue) => {
    initDarkMode();
  }
);

onMounted(() => {
  initDarkMode();
});
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
  color: var(--text-medium);
  background-color: var(--bg-light);
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  transition: all 0.3s ease;
}

::-webkit-scrollbar {
  width: 6px;
  height: 6px;
}

::-webkit-scrollbar-track {
  background: var(--bg-medium);
  border-radius: 3px;
}

::-webkit-scrollbar-thumb {
  background: var(--border-color);
  border-radius: 3px;
}

::-webkit-scrollbar-thumb:hover {
  background: var(--text-light);
}

/* 应用容器 */
.app-container {
  width: 100%;
  height: 100vh;
  overflow: hidden;
}

/* 深色模式下的Element Plus组件样式 */
.dark-mode {
  /* 输入框 */
  .el-input__wrapper {
    background-color: var(--input-bg) !important;
    border-color: var(--input-border) !important;
  }
  
  .el-input__wrapper .el-input__input {
    color: var(--input-text) !important;
  }
  
  .el-input__wrapper .el-input__placeholder {
    color: var(--input-placeholder) !important;
  }
  
  /* 确保输入框文字颜色正确显示 */
  input[type="text"],
  input[type="password"],
  input[type="email"] {
    color: var(--input-text) !important;
  }
  
  /* 下拉列表 */
  .el-select__wrapper {
    background-color: var(--select-bg) !important;
    border-color: var(--select-border) !important;
  }
  
  .el-select__input {
    color: var(--select-text) !important;
  }
  
  .el-select__placeholder {
    color: var(--input-placeholder) !important;
  }
  
  .el-select-dropdown {
    background-color: var(--select-option-bg) !important;
    border-color: var(--select-border) !important;
  }
  
  .el-select-dropdown__item {
    color: var(--select-text) !important;
    background-color: var(--select-option-bg) !important;
  }
  
  .el-select-dropdown__item:hover {
    background-color: var(--select-option-hover) !important;
  }
  
  .el-select-dropdown__item.selected {
    background-color: var(--primary-light) !important;
    color: var(--primary-color) !important;
  }
  
  .el-select-dropdown__item:focus {
    background-color: var(--select-option-hover) !important;
  }
  
  .el-select-dropdown__item.hover,
  .el-select-dropdown__item:hover {
    background-color: var(--select-option-hover) !important;
  }
  
  /* 按钮 */
  .el-button--text {
    color: var(--text-medium) !important;
  }
  
  .el-button--text:hover {
    color: var(--primary-color) !important;
  }
  
  /* 卡片 */
  .el-card {
    background-color: var(--bg-white) !important;
    border-color: var(--border-color) !important;
  }
  
  /* 表单 */
  .el-form-item__label {
    color: var(--text-medium) !important;
  }
  
  /* 下拉菜单 */
  .el-dropdown-menu {
    background-color: var(--bg-white) !important;
    border-color: var(--border-color) !important;
  }
  
  .el-dropdown-item {
    color: var(--text-medium) !important;
  }
  
  .el-dropdown-item:hover {
    background-color: var(--bg-medium) !important;
  }
  
  /* 菜单 */
  .el-menu {
    background-color: var(--bg-white) !important;
    border-color: var(--border-color) !important;
  }
  
  .el-menu-item {
    color: var(--text-medium) !important;
  }
  
  .el-menu-item:hover {
    background-color: var(--bg-medium) !important;
    color: var(--primary-color) !important;
  }
  
  .el-menu-item.is-active {
    background-color: var(--primary-light) !important;
    color: var(--primary-color) !important;
  }
}
</style>