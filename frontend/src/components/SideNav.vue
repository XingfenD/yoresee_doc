<template>
  <aside class="side-nav">
    <el-menu :default-active="currentActiveMenu" class="side-menu" @select="handleMenuSelect">
      <el-menu-item index="home">
        <el-icon>
          <House />
        </el-icon>
        <span>{{ t('navigation.home') }}</span>
      </el-menu-item>
      <el-menu-item index="knowledge-base">
        <el-icon>
          <Collection />
        </el-icon>
        <span>{{ t('navigation.knowledgeBase') }}</span>
      </el-menu-item>
    </el-menu>
  </aside>
</template>

<script setup>
import { computed } from 'vue';
import { useRouter } from 'vue-router';
import { useI18n } from 'vue-i18n';
import { House, Collection } from '@element-plus/icons-vue';

const router = useRouter();
const { t } = useI18n();

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
  } else if (key === 'knowledge-base') {
    router.push('/knowledge-base');
  }
};
</script>

<style scoped>
.side-nav {
  width: 240px;
  background-color: var(--bg-white);
  border-right: 1px solid var(--border-color);
  overflow-y: auto;
  transition: all 0.3s ease;
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

/* 响应式设计 */
@media (max-width: 768px) {
  .side-nav {
    width: 60px;
  }

  .side-menu .el-menu-item span {
    display: none;
  }
}
</style>