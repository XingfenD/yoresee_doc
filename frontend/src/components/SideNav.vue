<template>
  <aside class="side-nav" :class="{ 'collapsed': isCollapsed }">
    <el-menu :default-active="currentActiveMenu" class="side-menu" @select="handleMenuSelect">
      <el-menu-item index="home">
        <el-icon>
          <House />
        </el-icon>
        <span class="menu-text">{{ t('navigation.home') }}</span>
      </el-menu-item>
      <el-menu-item index="documents">
        <el-icon>
          <Document />
        </el-icon>
        <span class="menu-text">{{ t('navigation.myDocuments') }}</span>
      </el-menu-item>
      <el-menu-item index="knowledge-base">
        <el-icon>
          <Collection />
        </el-icon>
        <span class="menu-text">{{ t('navigation.knowledgeBase') }}</span>
      </el-menu-item>
      <el-menu-item index="templates">
        <el-icon>
          <Tickets />
        </el-icon>
        <span class="menu-text">{{ t('navigation.templates') }}</span>
      </el-menu-item>
    </el-menu>
    <button class="collapse-btn" @click="toggleCollapse"
      :title="isCollapsed ? t('common.expand') : t('common.collapse')">
      <el-icon>
        <DArrowRight v-if="isCollapsed" />
        <DArrowLeft v-else />
      </el-icon>
    </button>
  </aside>
</template>

<script setup>
import { ref, computed } from 'vue';
import { useRouter } from 'vue-router';
import { useI18n } from 'vue-i18n';
import { House, Collection, Document, Tickets, DArrowRight, DArrowLeft } from '@element-plus/icons-vue';

const router = useRouter();
const { t } = useI18n();

const isCollapsed = ref(localStorage.getItem('sideNavCollapsed') === 'true');

const toggleCollapse = () => {
  isCollapsed.value = !isCollapsed.value;
  localStorage.setItem('sideNavCollapsed', isCollapsed.value);
};

// 接收当前激活的菜单作为 props
const props = defineProps({
  activeMenu: {
    type: String,
    default: 'home'
  }
});

// 定义 emit 事件
const emit = defineEmits(['menuSelect']);

// 计算属性：当前激活的菜单项
const currentActiveMenu = computed(() => props.activeMenu);

// 处理菜单选择
const handleMenuSelect = (key) => {
  emit('menuSelect', key);
  if (key === 'home') {
    router.push('/');
  } else if (key === 'documents') {
    router.push('/mydocuments');
  } else if (key === 'knowledge-base') {
    router.push('/knowledge-base');
  } else if (key === 'templates') {
    router.push('/templates');
  }
};
</script>

<style scoped>
.side-nav {
  width: 240px;
  background-color: var(--bg-white);
  border-right: 1px solid var(--border-color);
  overflow: hidden;
  transition: all 0.3s ease;
  position: relative;
}

.side-nav.collapsed {
  width: 60px;
}

.side-menu {
  border-right: none;
}

.side-menu .el-menu-item {
  height: 48px;
  line-height: 48px;
  margin: 0;
  border-radius: 0;
  color: var(--text-medium);
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
  padding-left: var(--spacing-md);
}

.side-menu .el-menu-item:hover {
  background-color: var(--primary-light);
  color: var(--primary-color);
}

.side-menu .el-menu-item.is-active {
  background-color: var(--primary-light);
  color: var(--primary-color);
  border-right: 3px solid var(--primary-color);
}

.menu-text {
  display: inline-block;
  max-width: 160px;
  white-space: nowrap;
  overflow: hidden;
  opacity: 1;
  transform: translateX(0);
  transition: max-width 0.25s ease, opacity 0.2s ease, transform 0.25s ease;
}

.side-nav.collapsed .menu-text {
  max-width: 0;
  opacity: 0;
  transform: translateX(-8px);
}

.collapse-btn {
  position: absolute;
  right: -12px;
  top: 50%;
  transform: translateY(-50%);
  width: 24px;
  height: 24px;
  border-radius: 50%;
  background-color: var(--bg-white);
  border: 1px solid var(--border-color);
  color: var(--text-medium);
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.3s ease;
  z-index: 10;
  box-shadow: var(--shadow-sm);
}

.collapse-btn:hover {
  background-color: var(--primary-color);
  border-color: var(--primary-color);
  color: white;
  box-shadow: var(--shadow-md);
}

.collapse-btn .el-icon {
  font-size: 14px;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .side-nav {
    width: 60px;
  }

  .side-nav.collapsed {
    width: 60px;
  }

  .menu-text {
    display: none;
  }

  .collapse-btn {
    display: none;
  }
}
</style>
