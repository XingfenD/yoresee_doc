<template>
  <span :class="['app-tag', `app-tag--${resolvedType}`, `app-tag--${size}`]">
    <slot />
  </span>
</template>

<script setup>
import { computed } from 'vue';

const props = defineProps({
  type: {
    type: String,
    default: 'info'
  },
  size: {
    type: String,
    default: 'small'
  }
});

const resolvedType = computed(() => {
  const allowed = ['success', 'warning', 'danger', 'info', 'primary'];
  return allowed.includes(props.type) ? props.type : 'info';
});
</script>

<style scoped>
.app-tag {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  border-radius: 999px;
  border: 1px solid transparent;
  font-weight: 600;
  line-height: 1;
  white-space: nowrap;
  user-select: none;
}

.app-tag--small {
  height: 22px;
  padding: 0 8px;
  font-size: 12px;
}

.app-tag--success {
  background: rgba(16, 185, 129, 0.12);
  color: #10b981;
  border-color: rgba(16, 185, 129, 0.24);
}

.app-tag--warning {
  background: rgba(245, 158, 11, 0.12);
  color: #d97706;
  border-color: rgba(245, 158, 11, 0.24);
}

.app-tag--danger {
  background: rgba(239, 68, 68, 0.12);
  color: #ef4444;
  border-color: rgba(239, 68, 68, 0.24);
}

.app-tag--info {
  background: rgba(100, 116, 139, 0.12);
  color: #475569;
  border-color: rgba(100, 116, 139, 0.22);
}

.app-tag--primary {
  background: rgba(59, 130, 246, 0.12);
  color: #2563eb;
  border-color: rgba(59, 130, 246, 0.24);
}

:global(.dark-mode) .app-tag--success {
  background: rgba(16, 185, 129, 0.2);
  color: #34d399;
  border-color: rgba(16, 185, 129, 0.35);
}

:global(.dark-mode) .app-tag--warning {
  background: rgba(245, 158, 11, 0.2);
  color: #fbbf24;
  border-color: rgba(245, 158, 11, 0.35);
}

:global(.dark-mode) .app-tag--danger {
  background: rgba(239, 68, 68, 0.2);
  color: #f87171;
  border-color: rgba(239, 68, 68, 0.35);
}

:global(.dark-mode) .app-tag--info {
  background: rgba(148, 163, 184, 0.16);
  color: #cbd5e1;
  border-color: rgba(148, 163, 184, 0.28);
}

:global(.dark-mode) .app-tag--primary {
  background: rgba(59, 130, 246, 0.2);
  color: #93c5fd;
  border-color: rgba(59, 130, 246, 0.35);
}
</style>
