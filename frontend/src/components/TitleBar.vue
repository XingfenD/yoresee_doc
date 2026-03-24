<template>
  <div class="title-bar">
    <div class="title-bar__left">
      <el-button v-if="showBack" text class="title-bar__back" @click="emit('back')">
        <el-icon><ArrowLeft /></el-icon>
        {{ backText }}
      </el-button>
      <div v-if="$slots.title" class="title-bar__title">
        <slot name="title" />
      </div>
    </div>
    <div class="title-bar__actions">
      <slot name="actions" />
    </div>
  </div>
</template>

<script setup>
import { ArrowLeft } from '@element-plus/icons-vue';

defineProps({
  showBack: {
    type: Boolean,
    default: false
  },
  backText: {
    type: String,
    default: ''
  }
});

const emit = defineEmits(['back']);
</script>

<style scoped>
.title-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: var(--spacing-lg);
  padding-bottom: var(--spacing-md);
  border-bottom: 1px solid var(--border-color);
}

.title-bar__left {
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
  min-width: 0;
}

.title-bar__title {
  font-size: 20px;
  font-weight: 600;
  color: var(--text-dark);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.title-bar__actions {
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
}

.title-bar__back {
  height: 36px;
  padding: 8px 14px;
  font-size: 15px;
  font-weight: 500;
  border-radius: 10px;
}

.title-bar__back:hover,
.title-bar__back:focus {
  background: #cfe0ff !important;
  border-color: #cfe0ff !important;
  color: inherit !important;
}

.dark-mode .title-bar__back:hover,
.dark-mode .title-bar__back:focus {
  background: rgba(59, 130, 246, 0.24) !important;
  border-color: rgba(59, 130, 246, 0.24) !important;
  color: #e5e7eb !important;
}

.title-bar__back :deep(.el-icon) {
  font-size: 16px;
}

.title-bar__actions :deep(.el-button),
.title-bar__actions :deep(.el-button--small) {
  height: 32px;
  padding: 8px 14px;
  font-size: 14px;
  font-weight: 500;
  border-radius: 10px;
}

.title-bar__actions :deep(.el-button--primary) {
  background-color: var(--primary-color);
  border-color: var(--primary-color);
  color: #fff;
}

.title-bar__actions :deep(.el-button--primary:hover),
.title-bar__actions :deep(.el-button--primary:focus) {
  background-color: var(--primary-color);
  border-color: var(--primary-color);
  color: #fff;
  opacity: 0.9;
}

@media (max-width: 768px) {
  .title-bar {
    flex-direction: column;
    align-items: stretch;
    gap: var(--spacing-md);
  }
}
</style>
