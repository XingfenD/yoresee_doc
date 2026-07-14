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

  /* 输入框文字颜色（EP 变量覆盖不到原生 input） */
  input[type="text"],
  input[type="password"],
  input[type="email"] {
    color: var(--el-input-text-color);
  }

  /* message box 背景统一 */
  .el-message-box,
  .el-message-box__header,
  .el-message-box__content,
  .el-message-box__btns {
    background-color: var(--el-fill-color-blank);
  }

  .el-message-box__header { border-bottom: none; }
  .el-message-box__btns { border-top: none; }

  /* select dropdown 项 */
  .el-select-dropdown__item.selected {
    background-color: var(--el-color-primary-light-9);
    color: var(--el-color-primary);
  }

  /* date picker panel */
  .el-picker-panel__body,
  .el-picker-panel__content {
    background-color: var(--el-fill-color-blank);
    color: var(--el-text-color-regular);
  }

  .el-date-table td.in-range div,
  .el-date-table td.available:hover div {
    background-color: var(--bg-medium);
  }

  /* 数字输入按钮 */
  .el-input-number__decrease,
  .el-input-number__increase {
    background-color: var(--bg-medium);
    color: var(--text-dark);
    border-color: var(--border-color);
  }

  .el-input-number__decrease:hover,
  .el-input-number__increase:hover {
    background-color: var(--bg-light);
  }

  /* 对话框边框 */
  .el-dialog__header {
    border-bottom: 1px solid var(--el-border-color);
  }

  .el-dialog__footer {
    border-top: 1px solid var(--el-border-color);
  }

  /* 加载动画 */
  .el-loading-spinner .path {
    stroke: var(--primary-color);
  }
}
</style>
