<template>
  <div
    v-show="visible"
    class="code-language-floating"
    :style="styleObject"
    @mouseenter="$emit('enter')"
    @mouseleave="$emit('leave')"
  >
    <el-autocomplete
      v-model="draftProxy"
      class="code-language-input"
      popper-class="rich-text-code-language-popper"
      :fetch-suggestions="querySuggestions"
      :trigger-on-focus="true"
      :debounce="0"
      value-key="label"
      placeholder="语言"
      @focus="$emit('focus')"
      @blur="$emit('blur')"
      @change="$emit('change')"
      @select="$emit('select', $event)"
    >
      <template #default="{ item }">
        <div class="code-language-option">
          <span class="code-language-option-label">{{ item.label }}</span>
          <span class="code-language-option-value">{{ item.value }}</span>
        </div>
      </template>
    </el-autocomplete>
  </div>
</template>

<script setup>
import { computed } from 'vue';

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
</script>

<style scoped>
.code-language-floating {
  position: absolute;
  z-index: 34;
  pointer-events: auto;
}

.code-language-input {
  width: 180px;
}

.code-language-input :deep(.el-input__wrapper) {
  min-height: 28px;
  border-radius: 6px;
  box-shadow: var(--shadow-sm);
}

.code-language-option {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 10px;
}

.code-language-option-label {
  color: var(--text-primary);
}

.code-language-option-value {
  color: var(--text-light);
  font-size: 12px;
}

:global(.rich-text-code-language-popper.el-popper) {
  border: 1px solid var(--border-color);
  border-radius: 8px;
  box-shadow: var(--shadow-md);
}

:global(.rich-text-code-language-popper .el-autocomplete-suggestion__wrap) {
  max-height: 260px;
}

:global(.rich-text-code-language-popper .el-autocomplete-suggestion__list li) {
  border-radius: 6px;
  margin: 2px 4px;
  padding: 6px 10px;
}

:global(.rich-text-code-language-popper .el-autocomplete-suggestion__list li.highlighted),
:global(.rich-text-code-language-popper .el-autocomplete-suggestion__list li:hover) {
  background: var(--select-option-hover);
}

:global(.dark-mode .rich-text-code-language-popper.el-popper) {
  background: var(--bg-white);
  border-color: var(--border-color);
}

:global(.dark-mode .rich-text-code-language-popper .el-popper__arrow::before) {
  background: var(--bg-white);
  border-color: var(--border-color);
}

:global(.dark-mode .rich-text-code-language-popper .el-autocomplete-suggestion__wrap) {
  background: var(--bg-white);
}

:global(.dark-mode .rich-text-code-language-popper .el-autocomplete-suggestion__list li) {
  color: var(--text-primary);
  background: transparent;
}
</style>
