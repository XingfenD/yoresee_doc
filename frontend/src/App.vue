<template>
  <div class="app-container">
    <!-- 主内容区 -->
    <div class="main-container">
      <el-config-provider :locale="elementLocale">
        <router-view />
      </el-config-provider>
    </div>
  </div>
</template>

<script setup>
import { computed, onBeforeUnmount, onMounted, watch } from 'vue';
import { useUserStore } from './store/user';
import { ElConfigProvider } from 'element-plus';
import zhCn from 'element-plus/es/locale/lang/zh-cn';
import en from 'element-plus/es/locale/lang/en';
import { useI18n } from 'vue-i18n';

const userStore = useUserStore();
const { locale } = useI18n();

let darkModeObserver = null;

const applyDarkModeClass = () => {
  const enabled = Boolean(userStore.darkMode);
  document.documentElement.classList.toggle('dark-mode', enabled);
  document.body.classList.toggle('dark-mode', enabled);
};

watch(
  () => userStore.darkMode,
  () => {
    applyDarkModeClass();
  },
  { immediate: true }
);

onMounted(() => {
  applyDarkModeClass();
  darkModeObserver = new MutationObserver(() => {
    const enabled = Boolean(userStore.darkMode);
    const rootMatches = document.documentElement.classList.contains('dark-mode') === enabled;
    const bodyMatches = document.body.classList.contains('dark-mode') === enabled;
    if (rootMatches && bodyMatches) {
      return;
    }
    applyDarkModeClass();
  });
  darkModeObserver.observe(document.documentElement, { attributes: true, attributeFilter: ['class'] });
  darkModeObserver.observe(document.body, { attributes: true, attributeFilter: ['class'] });
});

onBeforeUnmount(() => {
  if (!darkModeObserver) {
    return;
  }
  darkModeObserver.disconnect();
  darkModeObserver = null;
});

const elementLocale = computed(() => (locale.value === 'zh' ? zhCn : en));
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
  display: flex;
  flex-direction: column;
}

/* 主内容区 */
.main-container {
  flex: 1;
  overflow: auto;
  padding: 0;
}

/* 全局返回按钮样式 */
.back-button {
  display: inline-flex;
  align-items: center;
  gap: var(--spacing-xs);
  font-weight: 500;
  color: var(--text-medium);
}

.back-button:hover {
  color: var(--primary-color);
}

/* 通用 Tabs 样式 */
.common-tabs .el-tabs__header {
  margin: 0 0 var(--spacing-md);
  border-bottom: 1px solid var(--border-color);
}

.common-tabs .el-tabs__nav-wrap {
  padding: 0 var(--spacing-sm);
}

.common-tabs .el-tabs__item {
  color: var(--text-medium);
  font-weight: 500;
}

.common-tabs .el-tabs__item.is-active {
  color: var(--primary-color);
}

.common-tabs .el-tabs__active-bar {
  background-color: var(--primary-color);
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

  /* message box (used by create knowledge base prompt) */
  .el-message-box {
    background-color: var(--bg-white) !important;
    border: 1px solid var(--border-color) !important;
    color: var(--text-dark) !important;
  }

  .el-message-box__header {
    background-color: var(--bg-white) !important;
    border-bottom: none !important;
  }

  .el-message-box__content {
    background-color: var(--bg-white) !important;
    color: var(--text-dark) !important;
  }

  .el-message-box__btns {
    background-color: var(--bg-white) !important;
    border-top: none !important;
  }

  .el-message-box__title {
    color: var(--text-dark) !important;
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

  .el-dropdown-item,
  .el-dropdown-menu__item {
    color: var(--text-dark) !important;
  }

  .el-dropdown-item:hover,
  .el-dropdown-item:focus,
  .el-dropdown-menu__item:hover,
  .el-dropdown-menu__item:focus,
  .el-dropdown-menu__item.is-hovering {
    background-color: var(--select-option-hover) !important;
    color: var(--text-dark) !important;
  }

  /* 日期选择器 */
  .el-picker__popper {
    background-color: var(--bg-white) !important;
    border-color: var(--border-color) !important;
  }

  .el-date-picker__header,
  .el-picker-panel__body,
  .el-picker-panel__content {
    background-color: var(--bg-white) !important;
    color: var(--text-dark) !important;
  }

  .el-date-table th,
  .el-date-table td,
  .el-date-table td .el-date-table-cell {
    color: var(--text-dark) !important;
  }

  .el-date-table td.in-range div,
  .el-date-table td.available:hover div {
    background-color: var(--bg-medium) !important;
  }

  /* 文本域 */
  .el-textarea__inner {
    background-color: var(--bg-white) !important;
    color: var(--text-dark) !important;
    border-color: var(--border-color) !important;
  }

  /* 数字输入 */
  .el-input-number__decrease,
  .el-input-number__increase {
    background-color: var(--bg-medium) !important;
    color: var(--text-dark) !important;
    border-color: var(--border-color) !important;
  }

  .el-input-number__decrease:hover,
  .el-input-number__increase:hover {
    background-color: var(--bg-light) !important;
  }

  /* 开关 */
  .el-switch {
    --el-switch-on-color: #3a7afe !important;
    --el-switch-off-color: #3a3a3a !important;
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

  /* 对话框 */
  .el-dialog {
    background-color: var(--bg-white) !important;
    border-color: var(--border-color) !important;
  }

  .el-dialog__header {
    background-color: var(--bg-white) !important;
    border-bottom: 1px solid var(--border-color) !important;
  }

  .el-dialog__title {
    color: var(--text-dark) !important;
  }

  .el-dialog__body {
    background-color: var(--bg-white) !important;
    color: var(--text-dark) !important;
  }

  .el-dialog__footer {
    background-color: var(--bg-white) !important;
    border-top: 1px solid var(--border-color) !important;
  }

  /* 输入框字数统计 */
  .el-input__count-inner {
    background-color: var(--input-bg) !important;
    color: var(--text-light) !important;
  }

  /* 加载遮罩 */
  .el-loading-mask {
    background-color: rgba(16, 18, 22, 0.72) !important;
  }

  .el-loading-spinner .path {
    stroke: var(--primary-color) !important;
  }

  .el-loading-spinner .el-loading-text {
    color: var(--text-medium) !important;
  }
}
</style>
