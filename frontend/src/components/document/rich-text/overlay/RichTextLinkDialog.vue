<template>
  <Teleport to="body">
    <div
      v-if="visible"
      class="link-dialog-mask"
      @mousedown.self="$emit('cancel')"
    >
      <div class="link-dialog-card" @mousedown.stop>
        <div class="link-dialog-title">{{ title || '编辑超链接' }}</div>
        <input
          v-model="valueProxy"
          class="link-dialog-input"
          type="text"
          placeholder="https://yoresee.cc"
          @keydown.enter.prevent="$emit('confirm')"
        >
        <div class="link-dialog-actions">
          <button type="button" class="link-dialog-btn" @click="$emit('cancel')">取消</button>
          <button type="button" class="link-dialog-btn link-dialog-btn--primary" @click="$emit('confirm')">
            确定
          </button>
        </div>
      </div>
    </div>
  </Teleport>
</template>

<script setup>
import { computed } from 'vue';

const props = defineProps({
  visible: {
    type: Boolean,
    default: false
  },
  title: {
    type: String,
    default: ''
  },
  modelValue: {
    type: String,
    default: ''
  }
});

const emit = defineEmits(['update:modelValue', 'cancel', 'confirm']);

const valueProxy = computed({
  get: () => props.modelValue,
  set: (value) => emit('update:modelValue', value)
});
</script>

<style scoped>
.link-dialog-mask {
  position: fixed;
  inset: 0;
  z-index: 90;
  background: rgba(0, 0, 0, 0.2);
  display: flex;
  align-items: center;
  justify-content: center;
}

.link-dialog-card {
  width: min(460px, calc(100vw - 40px));
  border-radius: 12px;
  border: 1px solid var(--border-color);
  background: var(--bg-white);
  box-shadow: var(--shadow-lg);
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.link-dialog-title {
  font-size: 16px;
  font-weight: 600;
  color: var(--text-primary);
}

.link-dialog-input {
  height: 36px;
  border-radius: 8px;
  border: 1px solid var(--border-color);
  background: var(--bg-white);
  color: var(--text-primary);
  padding: 0 10px;
  outline: none;
}

.link-dialog-input:focus {
  border-color: color-mix(in srgb, var(--primary-color) 52%, var(--border-color) 48%);
}

.link-dialog-actions {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
}

.link-dialog-btn {
  height: 32px;
  border-radius: 8px;
  border: 1px solid var(--border-color);
  background: var(--bg-white);
  color: var(--text-medium);
  font-size: 13px;
  padding: 0 12px;
  cursor: pointer;
}

.link-dialog-btn--primary {
  border-color: color-mix(in srgb, var(--primary-color) 52%, var(--border-color) 48%);
  background: var(--primary-color);
  color: #fff;
}
</style>
