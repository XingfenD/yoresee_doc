<template>
  <div class="template-pane" v-loading="loading">
    <div v-if="items.length === 0" class="template-empty">
      <el-empty :description="emptyText" />
    </div>
    <div v-else class="template-list" :class="`template-list--${layout}`">
      <div
        v-for="item in items"
        :key="item.id"
        class="template-card"
        :class="{ 'is-selected': selectedTemplateId === String(item.id), 'is-blank': Boolean(item.is_blank) }"
        @click="emit('select', item)"
      >
        <div v-if="layout === 'grid'" class="template-card-preview">
          {{ getPreviewText(item) }}
        </div>
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
  },
  layout: {
    type: String,
    default: 'list'
  }
});

const emit = defineEmits(['select']);

const getPreviewText = (item) => {
  if (item?.is_blank) {
    return '+';
  }
  const text = String(item?.name || '').trim();
  if (!text) {
    return 'T';
  }
  return text.slice(0, 1).toUpperCase();
};
</script>

<style scoped>
.template-pane {
  height: 100%;
}

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

.template-card-preview {
  height: 72px;
  border-radius: 10px;
  margin-bottom: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 22px;
  font-weight: 700;
  color: #2f65e2;
  background: linear-gradient(135deg, #e9f0ff 0%, #f7faff 100%);
}

.template-card.is-blank .template-card-preview {
  border: 1px dashed #91b0ff;
  background: linear-gradient(135deg, #edf3ff 0%, #f8fbff 100%);
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

.template-list--grid {
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 12px;
}

.template-list--grid .template-card {
  border: 1px solid #e7ebf2;
  border-bottom: 1px solid #e7ebf2;
  border-radius: 12px;
  padding: 12px;
  min-height: 158px;
  background: #fff;
}

.template-list--grid .template-card:hover {
  background-color: #f8faff;
  border-color: #d7e2ff;
}

.template-list--grid .template-card.is-selected {
  border-color: #6b93ff;
  background: #f2f7ff;
}

.template-list--grid .template-card.is-selected::before {
  left: auto;
  top: auto;
  bottom: 12px;
  right: 12px;
  width: 8px;
  height: 8px;
  border-radius: 999px;
  background: #3370ff;
}

.template-empty {
  padding: var(--spacing-md) 0;
  min-height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
}

.template-empty :deep(.el-empty) {
  margin: 0;
}

@media (max-width: 640px) {
  .template-list {
    grid-template-columns: 1fr;
  }

  .template-list--grid {
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

.dark-mode .template-card-preview {
  background: linear-gradient(135deg, #152235 0%, #1f314c 100%);
  color: #8db0ff;
}

.dark-mode .template-card.is-blank .template-card-preview {
  border-color: #4c8dff;
  background: linear-gradient(135deg, #102038 0%, #173054 100%);
}

.dark-mode .template-list--grid .template-card {
  border-color: #243040;
  background: #0f1218;
}

.dark-mode .template-list--grid .template-card:hover {
  border-color: #35517e;
  background: #121a26;
}

.dark-mode .template-list--grid .template-card.is-selected {
  border-color: #4c8dff;
  background: #142237;
}
</style>
