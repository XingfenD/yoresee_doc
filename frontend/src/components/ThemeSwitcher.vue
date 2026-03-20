<template>
  <el-button
    :size="'small'"
    :icon="isDarkMode ? 'Moon' : 'Sunny'"
    @click="toggleTheme"
    :class="['theme-toggle-btn']"
    :style="{ color: isDarkMode ? 'var(--text-light)' : 'var(--text-medium)' }"
  >
  </el-button>
</template>

<script setup>
import { computed } from 'vue';
import { Moon, Sunny } from '@element-plus/icons-vue';
import { useUserStore } from '@/store/user';

const userStore = useUserStore();
const isDarkMode = computed(() => userStore.darkMode);

const toggleTheme = () => {
  userStore.toggleDarkMode();
};
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
