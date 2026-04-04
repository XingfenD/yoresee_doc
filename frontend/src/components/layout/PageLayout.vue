<template>
  <div class="page-container">
    <TopNav
      :system-name="systemName"
      :current-language="currentLanguage"
      :is-dark-mode="isDarkMode"
      :user-avatar="userAvatar"
      :username="username"
      @change-language="$emit('change-language', $event)"
      @toggle-theme="$emit('toggle-theme')"
      @logout="$emit('logout')"
    />

    <div class="page-main">
      <SideNav
        :active-menu="activeMenu"
        :menu-items="sideMenuItems"
        :scene="sidebarScene"
        @menu-select="$emit('menu-select', $event)"
      />

      <div class="page-content" :class="contentPaddingClass">
        <TitleBar v-if="layout === 'list' && showHeader">
          <template #title>
            {{ title }}
          </template>
          <template v-if="$slots.actions" #actions>
            <slot name="actions" />
          </template>
        </TitleBar>
        <div class="page-header" v-else-if="showHeader">
          <h2 class="page-title">{{ title }}</h2>
          <div class="page-actions">
            <slot name="actions" />
          </div>
        </div>
        <div class="page-body">
          <slot />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, useSlots } from 'vue';
import SideNav from '@/components/layout/SideNav.vue';
import TopNav from '@/components/layout/TopNav.vue';
import TitleBar from '@/components/layout/TitleBar.vue';

const props = defineProps({
  systemName: {
    type: String,
    default: ''
  },
  currentLanguage: {
    type: String,
    default: 'zh-CN'
  },
  isDarkMode: {
    type: Boolean,
    default: false
  },
  userAvatar: {
    type: String,
    default: ''
  },
  username: {
    type: String,
    default: ''
  },
  activeMenu: {
    type: String,
    default: 'home'
  },
  sideMenuItems: {
    type: Array,
    default: () => []
  },
  sidebarScene: {
    type: String,
    default: 'home'
  },
  title: {
    type: String,
    default: ''
  },
  layout: {
    type: String,
    default: 'default'
  },
  contentPadding: {
    type: String,
    default: 'lg'
  }
});

defineEmits(['change-language', 'toggle-theme', 'logout', 'menu-select']);

const slots = useSlots();
const showHeader = computed(() => Boolean(props.title || slots.actions));
const contentPaddingClass = computed(() => `page-content--${props.contentPadding || 'lg'}`);
</script>

<style scoped>
.page-container {
  height: 100vh;
  display: flex;
  flex-direction: column;
  background-color: var(--bg-light);
  transition: all 0.3s ease;
}

.page-main {
  flex: 1;
  display: flex;
  overflow: hidden;
}

.page-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  min-height: 0;
  overflow-y: auto;
}

.page-content--lg {
  padding: var(--spacing-lg);
}

.page-content--xl {
  padding: var(--spacing-xl);
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: var(--spacing-lg);
  padding-bottom: var(--spacing-md);
  border-bottom: 1px solid var(--border-color);
}

.page-title {
  font-size: 20px;
  font-weight: 600;
  color: var(--text-dark);
  margin: 0;
}

.page-actions {
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
}

.page-actions :deep(.page-action-btn) {
  padding: 8px 14px;
  height: 32px;
  border-radius: var(--border-radius-md);
  font-weight: 500;
}

.page-actions :deep(.page-action-btn.el-button--primary) {
  background-color: var(--primary-color);
  border-color: var(--primary-color);
  color: #fff;
}

.page-actions :deep(.page-action-btn.el-button--primary:hover),
.page-actions :deep(.page-action-btn.el-button--primary:focus) {
  background-color: var(--primary-color);
  border-color: var(--primary-color);
  color: #fff;
  opacity: 0.9;
}

.page-body {
  display: flex;
  flex-direction: column;
  flex: 1;
  min-height: 0;
  gap: var(--spacing-lg);
}

@media (max-width: 1024px) {
  .page-content {
    padding: var(--spacing-md);
  }
}
</style>
