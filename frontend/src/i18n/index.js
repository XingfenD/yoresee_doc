import { createI18n } from 'vue-i18n';

// 导入语言包
import enUS from './locales/en-US';
import zhCN from './locales/zh-CN';

// 从localStorage获取语言设置
const lang = localStorage.getItem('language') || 'en';

const i18n = createI18n({
  legacy: false, // 使用组合式 API
  locale: lang,
  fallbackLocale: 'en',
  messages: {
    en: enUS,
    zh: zhCN
  }
});

export default i18n;