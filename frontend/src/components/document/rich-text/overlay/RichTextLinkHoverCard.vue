<template>
  <div
    v-show="visible"
    class="link-hover-card"
    :style="styleObject"
    @mouseenter="$emit('enter')"
    @mouseleave="$emit('leave')"
  >
    <div class="link-hover-title">超链接</div>
    <div class="link-hover-url">{{ href }}</div>
    <div class="link-hover-actions">
      <button type="button" class="link-hover-btn" @mousedown.prevent @click="$emit('edit')">
        编辑
      </button>
      <button
        type="button"
        class="link-hover-btn link-hover-btn--danger"
        @mousedown.prevent
        @click="$emit('remove')"
      >
        删除
      </button>
    </div>
  </div>
</template>

<script setup>
defineProps({
  visible: {
    type: Boolean,
    default: false
  },
  styleObject: {
    type: Object,
    default: () => ({})
  },
  href: {
    type: String,
    default: ''
  }
});

defineEmits(['enter', 'leave', 'edit', 'remove']);
</script>

<style scoped>
.link-hover-card {
  position: absolute;
  z-index: 36;
  width: 280px;
  border: 1px solid var(--border-color);
  border-radius: 10px;
  background: var(--bg-white);
  box-shadow: var(--shadow-md);
  padding: 10px;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.link-hover-title {
  font-size: 12px;
  font-weight: 600;
  color: var(--text-light);
}

.link-hover-url {
  font-size: 13px;
  color: var(--text-primary);
  line-height: 1.4;
  word-break: break-all;
}

.link-hover-actions {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: 8px;
}

.link-hover-btn {
  height: 28px;
  border-radius: 7px;
  border: 1px solid var(--border-color);
  background: var(--bg-white);
  color: var(--text-medium);
  padding: 0 10px;
  font-size: 12px;
  cursor: pointer;
  transition: all 0.18s ease;
}

.link-hover-btn:hover {
  border-color: color-mix(in srgb, var(--primary-color) 48%, var(--border-color) 52%);
  color: var(--primary-color);
}

.link-hover-btn--danger:hover {
  border-color: color-mix(in srgb, #ef4444 48%, var(--border-color) 52%);
  color: #ef4444;
}
</style>
