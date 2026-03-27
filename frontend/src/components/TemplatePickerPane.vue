<template>
  <div v-loading="loading">
    <div v-if="items.length === 0" class="template-empty">
      <el-empty :description="emptyText" />
    </div>
    <div v-else class="template-list">
      <div
        v-for="item in items"
        :key="item.id"
        class="template-card"
        :class="{ 'is-selected': selectedTemplateId === String(item.id) }"
        @click="emit('select', item)"
      >
        <div class="template-card-title">
          {{ item.name }}
        </div>
        <div class="template-card-desc">
          {{ item.description || fallbackDescription }}
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
defineProps({
  loading: {
    type: Boolean,
    default: false
  },
  items: {
    type: Array,
    default: () => []
  },
  selectedTemplateId: {
    type: String,
    default: ''
  },
  emptyText: {
    type: String,
    default: ''
  },
  fallbackDescription: {
    type: String,
    default: ''
  }
});

const emit = defineEmits(['select']);
</script>

<style scoped>
.template-list {
  display: grid;
  grid-template-columns: 1fr;
  gap: 0;
}

.template-card {
  cursor: pointer;
  border-bottom: 1px solid #eef0f4;
  padding: 12px 8px 12px 12px;
  transition: background-color 0.2s ease, border-color 0.2s ease;
  position: relative;
}

.template-card:hover {
  background-color: #f7f8fa;
}

.template-card.is-selected {
  background-color: #f0f5ff;
}

.template-card-title {
  font-weight: 600;
  color: #111827;
  display: flex;
  align-items: center;
  gap: 8px;
}

.template-card-desc {
  margin-top: 4px;
  color: #6b7280;
  font-size: 12px;
}

.template-card.is-selected::before {
  content: '';
  position: absolute;
  left: 0;
  top: 8px;
  bottom: 8px;
  width: 3px;
  border-radius: 2px;
  background: #3370ff;
}

.template-empty {
  padding: var(--spacing-md) 0;
}

@media (max-width: 640px) {
  .template-list {
    grid-template-columns: 1fr;
  }
}

.dark-mode .template-card {
  background-color: transparent;
  border-color: #1f2937;
}

.dark-mode .template-card:hover {
  background-color: #111827;
}

.dark-mode .template-card.is-selected {
  background-color: #0b1f3a;
}

.dark-mode .template-card-title {
  color: #e5e7eb;
}

.dark-mode .template-card-desc {
  color: #9ca3af;
}
</style>
