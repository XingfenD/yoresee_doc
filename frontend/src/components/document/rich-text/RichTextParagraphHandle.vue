<template>
  <div
    v-show="visible"
    ref="rootRef"
    class="rich-text-paragraph-handle"
    :style="style"
    @mouseenter="onHandleEnter"
    @mouseleave="onHandleLeave"
  >
    <button
      v-if="mode === 'empty'"
      ref="triggerRef"
      type="button"
      class="handle-btn handle-btn--plus"
      :title="plusTitle"
      @mousedown.prevent
      @click.stop="toggleMenu"
    >
      <el-icon><Plus /></el-icon>
    </button>
    <template v-else>
      <button
        type="button"
        class="handle-btn handle-btn--type"
        :title="typeLabel"
        @mousedown.prevent
      >
        <el-icon><component :is="currentTypeIcon" /></el-icon>
      </button>
      <button
        ref="triggerRef"
        type="button"
        class="handle-btn handle-btn--more"
        :title="moreTitle"
        @mousedown.prevent
        @click.stop="toggleMenu"
      >
        <el-icon><MoreFilled /></el-icon>
      </button>
    </template>
  </div>

  <AppMenu
    ref="menuRef"
    :visible="menuVisible"
    :x="menuPosition.x"
    :y="menuPosition.y"
    :ignore-elements="[triggerRef, rootRef]"
    @close="closeMenu"
  >
    <AppMenuItem
      v-for="item in actions"
      :key="item.key"
      :danger="Boolean(item.danger)"
      :disabled="Boolean(item.disabled)"
      @click="selectAction(item.key)"
    >
      <template #icon>
        <el-icon v-if="resolveActionIcon(item)">
          <component :is="resolveActionIcon(item)" />
        </el-icon>
      </template>
      {{ item.label }}
    </AppMenuItem>
  </AppMenu>
</template>

<script setup>
import { computed, nextTick, ref, watch } from 'vue';
import {
  ArrowDownBold,
  ArrowUpBold,
  ChatLineRound,
  Connection,
  Delete,
  Document,
  Grid,
  List,
  Menu,
  Minus,
  MoreFilled,
  Plus,
  Tickets
} from '@element-plus/icons-vue';
import AppMenu from '@/components/base/AppMenu.vue';
import AppMenuItem from '@/components/base/AppMenuItem.vue';

const props = defineProps({
  visible: {
    type: Boolean,
    default: false
  },
  style: {
    type: Object,
    default: () => ({})
  },
  mode: {
    type: String,
    default: 'normal'
  },
  blockType: {
    type: String,
    default: 'paragraph'
  },
  actions: {
    type: Array,
    default: () => []
  },
  plusTitle: {
    type: String,
    default: ''
  },
  moreTitle: {
    type: String,
    default: ''
  },
  typeLabel: {
    type: String,
    default: ''
  }
});

const emit = defineEmits(['action', 'mouseenter', 'mouseleave']);

const triggerRef = ref(null);
const rootRef = ref(null);
const menuRef = ref(null);
const menuVisible = ref(false);
const menuPosition = ref({ x: 0, y: 0 });

const blockTypeIconMap = {
  paragraph: Document,
  heading: Menu,
  list: List,
  quote: ChatLineRound,
  code: Tickets,
  table: Grid,
  divider: Minus,
  mindmap: Connection,
  drawio: Connection
};

const actionIconMap = {
  'add-above': ArrowUpBold,
  'add-below': ArrowDownBold,
  delete: Delete
};

const currentTypeIcon = computed(() => blockTypeIconMap[props.blockType] || Document);

const resolveActionIcon = (item) => {
  if (!item) {
    return null;
  }
  if (item.icon && typeof item.icon === 'object') {
    return item.icon;
  }
  return actionIconMap[item.iconKey || item.key] || null;
};

const updateMenuPosition = async () => {
  const triggerEl = triggerRef.value;
  if (!triggerEl) {
    return;
  }
  const rect = triggerEl.getBoundingClientRect();
  const baseX = rect.right + 6;
  const baseY = rect.top;
  menuPosition.value = { x: baseX, y: baseY };
  await nextTick();

  const menuEl = menuRef.value?.getMenuElement?.();
  if (!menuEl) {
    return;
  }
  const maxX = Math.max(8, window.innerWidth - menuEl.offsetWidth - 8);
  const maxY = Math.max(8, window.innerHeight - menuEl.offsetHeight - 8);
  menuPosition.value = {
    x: Math.max(8, Math.min(baseX, maxX)),
    y: Math.max(8, Math.min(baseY, maxY))
  };
};

const openMenu = async () => {
  menuVisible.value = true;
  emit('mouseenter');
  await updateMenuPosition();
};

const closeMenu = () => {
  menuVisible.value = false;
  emit('mouseleave');
};

const toggleMenu = () => {
  if (menuVisible.value) {
    closeMenu();
    return;
  }
  openMenu();
};

const selectAction = (actionKey) => {
  emit('action', actionKey);
  closeMenu();
};

const onHandleEnter = () => {
  emit('mouseenter');
};

const onHandleLeave = () => {
  if (menuVisible.value) {
    return;
  }
  emit('mouseleave');
};

watch(
  () => props.visible,
  (next) => {
    if (!next && menuVisible.value) {
      closeMenu();
    }
  }
);
</script>

<style scoped>
.rich-text-paragraph-handle {
  position: absolute;
  z-index: 26;
  display: inline-flex;
  align-items: flex-start;
  gap: 4px;
}

.handle-btn {
  width: 22px;
  height: 22px;
  border-radius: 6px;
  border: 1px solid var(--border-color);
  background: var(--bg-white);
  color: var(--text-light);
  box-shadow: var(--shadow-sm);
  display: inline-flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.18s ease;
}

.handle-btn:hover {
  border-color: color-mix(in srgb, var(--primary-color) 55%, var(--border-color) 45%);
  color: var(--primary-color);
}

.handle-btn--plus:hover {
  border-color: color-mix(in srgb, var(--primary-color) 65%, var(--border-color) 35%);
  background: color-mix(in srgb, var(--primary-color) 9%, var(--bg-white) 91%);
}
</style>
