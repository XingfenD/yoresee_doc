<template>
  <div class="auth-container">
    <header class="auth-nav">
      <div class="nav-right">
        <AppDropdown trigger="click" class="nav-item" @command="$emit('change-language', $event)">
          <span class="nav-link">
            <el-icon :size="18"><Flag v-if="currentLanguage === 'en'" /><ChatLineRound v-else /></el-icon>
          </span>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item command="en" :icon="'Flag'">
                {{ t('language.english') }}
              </el-dropdown-item>
              <el-dropdown-item command="zh" :icon="'ChatLineRound'">
                {{ t('language.chinese') }}
              </el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </AppDropdown>

        <div class="nav-item theme-switch">
          <span class="nav-link" @click="$emit('toggle-theme')">
            <el-icon :size="18"><Moon v-if="isDarkMode" /><Sunny v-else /></el-icon>
          </span>
        </div>
      </div>
    </header>

    <div class="auth-form-wrapper">
      <div class="auth-header">
        <h2>{{ systemName }}</h2>
        <p>{{ subtitle }}</p>
      </div>
      <slot />
    </div>
  </div>
</template>

<script setup>
import { useI18n } from 'vue-i18n';
import { Flag, ChatLineRound, Moon, Sunny } from '@element-plus/icons-vue';
import AppDropdown from '@/components/base/AppDropdown.vue';

defineProps({
  currentLanguage: { type: String, default: 'en' },
  isDarkMode: { type: Boolean, default: false },
  systemName: { type: String, default: 'Yoresee' },
  subtitle: { type: String, default: '' }
});

defineEmits(['change-language', 'toggle-theme']);

const { t } = useI18n();
</script>

<style scoped>
.auth-container {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  background-color: var(--bg-light);
  padding: var(--spacing-md);
  transition: all 0.3s ease;
  position: relative;
}

.auth-nav {
  position: absolute;
  top: 0;
  right: 0;
  padding: var(--spacing-md);
}

.nav-right {
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
}

.nav-item {
  display: flex;
  align-items: center;
  margin-left: var(--spacing-sm);
}

.nav-link {
  display: flex;
  align-items: center;
  padding: var(--spacing-xs) var(--spacing-sm);
  border-radius: var(--border-radius-md);
  color: var(--text-medium);
  transition: all 0.3s ease;
  cursor: pointer;
}

.nav-link:hover {
  background-color: var(--bg-medium);
  color: var(--primary-color);
}

.theme-switch {
  padding: var(--spacing-xs) var(--spacing-sm);
}

.auth-form-wrapper {
  margin: auto;
  width: 100%;
  max-width: 400px;
  background-color: var(--bg-white);
  border-radius: var(--border-radius-lg);
  box-shadow: var(--shadow-md);
  padding: var(--spacing-xl);
  transition: all 0.3s ease;
}

.auth-form-wrapper:hover {
  box-shadow: var(--shadow-lg);
}

.auth-header {
  text-align: center;
  margin-bottom: var(--spacing-lg);
}

.auth-header h2 {
  font-size: 24px;
  font-weight: 600;
  color: var(--text-dark);
  margin-bottom: var(--spacing-sm);
}

.auth-header p {
  font-size: 14px;
  color: var(--text-light);
  margin: 0;
}

@media (max-width: 768px) {
  .auth-form-wrapper {
    padding: var(--spacing-lg) var(--spacing-md);
  }

  .auth-header h2 {
    font-size: 20px;
  }
}
</style>
