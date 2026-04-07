<template>
  <div
    v-show="visible"
    class="code-language-floating"
    :style="styleObject"
    @mouseenter="$emit('enter')"
    @mouseleave="$emit('leave')"
  >
    <div class="code-language-panel">
      <el-input
        v-model="draftProxy"
        size="small"
        class="code-language-input"
        placeholder="语言"
        @focus="handleFocus"
        @blur="handleBlur"
        @input="handleInput"
        @keydown.stop
        @keydown.down.prevent="handleArrowDown"
        @keydown.up.prevent="handleArrowUp"
        @keydown.enter.prevent="handleEnter"
        @keydown.esc.prevent="closeMenu"
      />
      <div
        v-show="menuVisible && options.length > 0"
        class="code-language-menu"
        @mousedown.prevent
        @mouseenter="handleMenuEnter"
        @mouseleave="handleMenuLeave"
      >
        <AppMenuItem
          v-for="(item, index) in options"
          :key="item.value"
          class="code-language-option-row"
          :class="{ 'is-highlighted': highlightedIndex === index }"
          @mouseenter="highlightedIndex = index"
          @click="selectOption(item)"
        >
          <span class="code-language-option-label">{{ item.label }}</span>
          <span class="code-language-option-value">{{ item.value }}</span>
        </AppMenuItem>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, onBeforeUnmount, ref, watch } from 'vue';
import AppMenuItem from '@/components/base/AppMenuItem.vue';

const props = defineProps({
  visible: {
    type: Boolean,
    default: false
  },
  styleObject: {
    type: Object,
    default: () => ({})
  },
  modelValue: {
    type: String,
    default: 'plaintext'
  },
  querySuggestions: {
    type: Function,
    required: true
  }
});

const emit = defineEmits([
  'update:modelValue',
  'enter',
  'leave',
  'focus',
  'blur',
  'change',
  'select'
]);

const draftProxy = computed({
  get: () => props.modelValue,
  set: (value) => emit('update:modelValue', value)
});

const options = ref([]);
const menuVisible = ref(false);
const highlightedIndex = ref(-1);
const inputFocused = ref(false);
const menuHovering = ref(false);
let blurTimer = 0;

const clearBlurTimer = () => {
  if (blurTimer) {
    window.clearTimeout(blurTimer);
    blurTimer = 0;
  }
};

const updateOptions = (keyword = '') => {
  props.querySuggestions?.(keyword, (items) => {
    options.value = Array.isArray(items) ? items : [];
    highlightedIndex.value = options.value.length > 0 ? 0 : -1;
  });
};

const handleFocus = () => {
  clearBlurTimer();
  inputFocused.value = true;
  emit('focus');
  updateOptions(draftProxy.value || '');
  menuVisible.value = true;
};

const handleInput = (value) => {
  updateOptions(value);
  menuVisible.value = true;
};

const closeMenu = () => {
  menuVisible.value = false;
};

const finalizeBlur = () => {
  if (inputFocused.value || menuHovering.value) {
    return;
  }
  menuVisible.value = false;
  emit('change');
  emit('blur');
};

const selectOption = (item) => {
  if (!item) {
    return;
  }
  emit('update:modelValue', item.value || '');
  emit('select', item);
  menuVisible.value = false;
};

const handleBlur = () => {
  inputFocused.value = false;
  clearBlurTimer();
  blurTimer = window.setTimeout(() => {
    blurTimer = 0;
    finalizeBlur();
  }, 120);
};

const handleMenuEnter = () => {
  menuHovering.value = true;
  clearBlurTimer();
  emit('enter');
};

const handleMenuLeave = () => {
  menuHovering.value = false;
  finalizeBlur();
  emit('leave');
};

const moveHighlight = (delta) => {
  const length = options.value.length;
  if (length <= 0) {
    highlightedIndex.value = -1;
    return;
  }
  const current = highlightedIndex.value < 0 ? 0 : highlightedIndex.value;
  highlightedIndex.value = (current + delta + length) % length;
};

const handleArrowDown = () => {
  if (!menuVisible.value) {
    menuVisible.value = true;
    updateOptions(draftProxy.value || '');
  }
  moveHighlight(1);
};

const handleArrowUp = () => {
  if (!menuVisible.value) {
    menuVisible.value = true;
    updateOptions(draftProxy.value || '');
  }
  moveHighlight(-1);
};

const handleEnter = () => {
  if (menuVisible.value && highlightedIndex.value >= 0) {
    const next = options.value[highlightedIndex.value];
    if (next) {
      selectOption(next);
      return;
    }
  }
  emit('change');
};

watch(
  () => props.visible,
  (nextVisible) => {
    if (nextVisible) {
      updateOptions(draftProxy.value || '');
      menuVisible.value = true;
      return;
    }
    clearBlurTimer();
    menuVisible.value = false;
    inputFocused.value = false;
    menuHovering.value = false;
  }
);

onBeforeUnmount(() => {
  clearBlurTimer();
});
</script>

<style scoped>
.code-language-floating {
  position: absolute;
  z-index: 34;
  pointer-events: auto;
}

.code-language-panel {
  position: relative;
}

.code-language-input {
  width: 132px;
}

:deep(.code-language-input .el-input__wrapper) {
  min-height: 26px;
  border-radius: 6px;
  box-shadow: var(--shadow-sm);
}

:deep(.code-language-input .el-input__inner) {
  font-size: 12px;
}

.code-language-menu {
  position: absolute;
  right: 0;
  top: calc(100% + 6px);
  z-index: 1;
  min-width: 176px;
  max-height: 220px;
  overflow: auto;
  padding: 4px;
  border: 1px solid var(--border-color);
  border-radius: 8px;
  background: var(--bg-white);
  box-shadow: var(--shadow-md);
}

.code-language-menu::-webkit-scrollbar {
  width: 6px;
}

.code-language-menu::-webkit-scrollbar-thumb {
  background: color-mix(in srgb, var(--border-color) 85%, transparent);
  border-radius: 4px;
}

.code-language-option-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 8px;
  border-radius: 6px;
}

.code-language-option-label {
  color: var(--text-primary);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.code-language-option-value {
  color: var(--text-light);
  font-size: 11px;
  margin-left: 8px;
}

.code-language-menu :deep(.app-menu-item.is-highlighted) {
  background: var(--select-option-hover);
  color: var(--primary-color);
}
</style>
