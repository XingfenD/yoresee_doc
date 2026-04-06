<template>
  <div class="handle-menu-branch" :class="{ 'is-submenu': level > 0 }">
    <div
      v-for="(item, index) in actions"
      :key="buildPath(item, index)"
      class="handle-menu-node"
      @mouseenter="handleNodeEnter(item, index)"
    >
      <button
        type="button"
        class="handle-menu-row"
        :class="{
          'is-danger': Boolean(item?.danger),
          'is-disabled': Boolean(item?.disabled),
          'is-group': isGroupAction(item)
        }"
        :disabled="Boolean(item?.disabled)"
        @click.stop.prevent="handleRowClick(item)"
      >
        <span class="handle-menu-row-icon">
          <el-icon v-if="resolveIcon(item)">
            <component :is="resolveIcon(item)" />
          </el-icon>
        </span>
        <span class="handle-menu-row-label">{{ item?.label || '' }}</span>
        <span v-if="isGroupAction(item)" class="handle-menu-row-arrow">
          <el-icon><ArrowRight /></el-icon>
        </span>
      </button>

      <div
        v-if="isGroupAction(item) && isGroupOpen(buildPath(item, index))"
        class="handle-menu-subpanel"
      >
        <RichTextHandleMenuBranch
          :actions="item.children"
          :level="level + 1"
          :path-prefix="buildPath(item, index)"
          :open-paths="openPaths"
          :resolve-icon="resolveIcon"
          @select="$emit('select', $event)"
          @group-enter="$emit('group-enter', $event)"
        />
      </div>
    </div>
  </div>
</template>

<script setup>
import { ArrowRight } from '@element-plus/icons-vue';

const props = defineProps({
  actions: {
    type: Array,
    default: () => []
  },
  level: {
    type: Number,
    default: 0
  },
  pathPrefix: {
    type: String,
    default: ''
  },
  openPaths: {
    type: Array,
    default: () => []
  },
  resolveIcon: {
    type: Function,
    default: () => null
  }
});

const emit = defineEmits(['select', 'group-enter']);

const isGroupAction = (item) => Array.isArray(item?.children) && item.children.length > 0;

const buildPath = (item, index) => {
  const token = item?.key || `idx-${index}`;
  return props.pathPrefix ? `${props.pathPrefix}.${token}` : token;
};

const isGroupOpen = (path) => props.openPaths[props.level] === path;

const handleNodeEnter = (item, index) => {
  const path = isGroupAction(item) ? buildPath(item, index) : null;
  emit('group-enter', {
    level: props.level,
    path
  });
};

const handleRowClick = (item) => {
  if (!item || item.disabled || isGroupAction(item)) {
    return;
  }
  emit('select', item.key);
};
</script>

<style scoped>
.handle-menu-branch {
  min-width: 170px;
}

.handle-menu-branch.is-submenu {
  border: 1px solid var(--border-color);
  border-radius: var(--border-radius-sm);
  background: var(--bg-white);
  box-shadow: var(--shadow-md);
  padding: var(--spacing-xs) 0;
}

.handle-menu-node {
  position: relative;
}

.handle-menu-row {
  width: 100%;
  padding: var(--spacing-xs) var(--spacing-md);
  background: transparent;
  border: none;
  text-align: left;
  display: flex;
  align-items: center;
  gap: var(--spacing-xs);
  color: var(--text-medium);
  cursor: pointer;
}

.handle-menu-row-icon {
  width: 16px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  color: var(--text-light);
}

.handle-menu-row-label {
  flex: 1;
}

.handle-menu-row-arrow {
  display: inline-flex;
  align-items: center;
  color: var(--text-light);
}

.handle-menu-row:hover {
  background-color: var(--bg-light);
  color: var(--primary-color);
}

.handle-menu-row:hover .handle-menu-row-icon,
.handle-menu-row:hover .handle-menu-row-arrow {
  color: var(--primary-color);
}

.handle-menu-row.is-danger:hover {
  color: #f56c6c;
}

.handle-menu-row.is-danger:hover .handle-menu-row-icon,
.handle-menu-row.is-danger:hover .handle-menu-row-arrow {
  color: #f56c6c;
}

.handle-menu-row.is-disabled {
  color: var(--text-light);
  cursor: not-allowed;
}

.handle-menu-row.is-disabled:hover {
  color: var(--text-light);
  background: transparent;
}

.handle-menu-row.is-disabled:hover .handle-menu-row-icon,
.handle-menu-row.is-disabled:hover .handle-menu-row-arrow {
  color: var(--text-light);
}

.handle-menu-row.is-group {
  font-weight: 600;
}

.handle-menu-subpanel {
  position: absolute;
  top: -6px;
  left: calc(100% - 3px);
  z-index: 1;
}
</style>
