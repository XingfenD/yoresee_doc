<template>
  <el-button :size="'small'" :icon="isDarkMode ? 'Moon' : 'Sunny'" @click="toggleTheme" :class="['theme-toggle-btn']" :style="{ color: isDarkMode ? 'var(--text-light)' : 'var(--text-medium)' }">
  </el-button>
</template>

<script setup>
import { ref, watch } from 'vue';
import { Moon, Sunny } from '@element-plus/icons-vue';

// 响应式数据
const isDarkMode = ref(false);

// 初始化主题
const initTheme = () => {
  const savedDarkMode = localStorage.getItem('darkMode');
  if (savedDarkMode === 'true') {
    isDarkMode.value = true;
    document.documentElement.classList.add('dark-mode');
  }
};

// 处理主题切换
const toggleTheme = () => {
  isDarkMode.value = !isDarkMode.value;
  if (isDarkMode.value) {
    document.documentElement.classList.add('dark-mode');
    localStorage.setItem('darkMode', 'true');
  } else {
    document.documentElement.classList.remove('dark-mode');
    localStorage.setItem('darkMode', 'false');
  }
};

// 初始化
initTheme();
</script>

<style scoped>
.theme-toggle-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  transition: color 0.3s ease;
  
  &:hover {
    color: var(--primary-color);
  }
}

/* 确保图标大小一致 */
.theme-toggle-btn .el-icon {
  font-size: 18px;
}
</style>