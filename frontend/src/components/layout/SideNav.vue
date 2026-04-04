<template>
  <aside class="side-nav" :class="{ 'collapsed': isCollapsed }">
    <el-menu :default-active="currentActiveMenu" class="side-menu" @select="handleMenuSelect">
      <el-menu-item v-for="item in filteredMenuItems" :key="item.key" :index="item.key">
        <el-icon v-if="item.icon">
          <component :is="item.icon" />
        </el-icon>
        <span class="menu-text">{{ getMenuLabel(item) }}</span>
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
import { ref, computed, watch } from 'vue';
import { useRouter } from 'vue-router';
import { useI18n } from 'vue-i18n';
import { House, Collection, Document, Tickets, Search, DArrowRight, DArrowLeft } from '@element-plus/icons-vue';
import { querySideBarDisplay } from '@/services/auth';
import { useApiAction } from '@/composables/actions/useApiAction';
import {
  getSideBarDisplayTabsCache,
  setSideBarDisplayTabsCache
} from '@/composables/layout/useSideBarDisplayCache';

const router = useRouter();
const { t } = useI18n();
const { runSilent, runWithLoading } = useApiAction({ t });

const isCollapsed = ref(localStorage.getItem('sideNavCollapsed') === 'true');

const toggleCollapse = () => {
  isCollapsed.value = !isCollapsed.value;
  localStorage.setItem('sideNavCollapsed', isCollapsed.value);
};

// 接收当前激活的菜单作为 props
const defaultMenuItems = [
  { key: 'home', labelKey: 'navigation.home', icon: House, route: '/' },
  { key: 'search', labelKey: 'navigation.search', icon: Search, route: '/search' },
  { key: 'documents', labelKey: 'navigation.myDocuments', icon: Document, route: '/mydocuments' },
  { key: 'knowledge-base', labelKey: 'navigation.knowledgeBase', icon: Collection, route: '/knowledge-base' },
  { key: 'templates', labelKey: 'navigation.templates', icon: Tickets, route: '/templates' }
];

const props = defineProps({
  activeMenu: {
    type: String,
    default: 'home'
  },
  menuItems: {
    type: Array,
    default: () => []
  },
  scene: {
    type: String,
    default: ''
  }
});

// 定义 emit 事件
const emit = defineEmits(['menuSelect']);

// 计算属性：当前激活的菜单项
const currentActiveMenu = computed(() => props.activeMenu);
const resolvedMenuItems = computed(() => (
  Array.isArray(props.menuItems) && props.menuItems.length ? props.menuItems : defaultMenuItems
));
const filterLoading = ref(false);
const sceneKey = computed(() => (props.scene || '').trim());
const displayTabs = ref(getSideBarDisplayTabsCache(sceneKey.value) || []);

const filteredMenuItems = computed(() => {
  if (!sceneKey.value) {
    return resolvedMenuItems.value;
  }
  const tabsSet = new Set(displayTabs.value);
  return resolvedMenuItems.value.filter((item) => tabsSet.has(item.key));
});

const loadDisplayTabs = async () => {
  const scene = sceneKey.value;
  if (!scene) {
    displayTabs.value = [];
    return;
  }

  const cachedTabs = getSideBarDisplayTabsCache(scene);
  if (Array.isArray(cachedTabs)) {
    displayTabs.value = cachedTabs;
    return;
  }

  displayTabs.value = [];
  await runWithLoading(
    filterLoading,
    () =>
      runSilent(
        () => querySideBarDisplay(scene),
        {
          context: 'loadSideBarDisplay',
          onSuccess: (resp) => {
            const tabs = resp.display_tabs || [];
            setSideBarDisplayTabsCache(scene, tabs);
            displayTabs.value = tabs;
          },
          onError: () => {
            // fail closed so privileged menu is not exposed by client fallback.
            displayTabs.value = [];
          }
        }
      )
  );
};

const getMenuLabel = (item) => {
  if (item.label) {
    return item.label;
  }
  if (item.labelKey) {
    return t(item.labelKey);
  }
  return item.key;
};

watch(
  () => sceneKey.value,
  () => {
    loadDisplayTabs();
  },
  { immediate: true }
);

// 处理菜单选择
const handleMenuSelect = (key) => {
  emit('menuSelect', key);
  const selected = filteredMenuItems.value.find((item) => item.key === key);
  if (selected?.route) {
    router.push(selected.route);
    return;
  }
  if (key === 'home') {
    router.push('/');
  } else if (key === 'search') {
    router.push('/search');
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
  transition: width 0.3s ease, background-color 0.2s ease, border-color 0.2s ease;
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
  transition: background-color 0.2s ease, color 0.2s ease;
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
