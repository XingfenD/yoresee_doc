import { createApp } from 'vue';
import { createPinia } from 'pinia';
import ElementPlus from 'element-plus';
import 'element-plus/dist/index.css';
import * as ElementPlusIconsVue from '@element-plus/icons-vue';
import App from './App.vue';
import router from './router';
import './styles/variables.css';
import './styles/mention.css';
import i18n from './i18n';

const sanitizeAuthStorage = () => {
  try {
    const token = localStorage.getItem('token');
    const rawUserInfo = localStorage.getItem('userInfo');

    if (!token || token === 'null' || token === 'undefined') {
      localStorage.removeItem('token');
      localStorage.removeItem('userInfo');
      return;
    }

    if (!rawUserInfo) {
      localStorage.removeItem('token');
      return;
    }

    const parsed = JSON.parse(rawUserInfo);
    if (!parsed || typeof parsed !== 'object' || !parsed.username) {
      localStorage.removeItem('token');
      localStorage.removeItem('userInfo');
    }
  } catch (error) {
    localStorage.removeItem('token');
    localStorage.removeItem('userInfo');
  }
};

sanitizeAuthStorage();

const app = createApp(App);

// 注册Element Plus图标
for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
  app.component(key, component);
}

// 使用插件
app.use(createPinia());
app.use(router);
app.use(ElementPlus);
app.use(i18n);

app.mount('#app');
