<template>
  <div class="vertical-section">
    <div class="section-header">
      <h3 class="section-title">{{ title }}</h3>
    </div>
    <div class="section-content">
      <el-empty v-if="showEmpty" :description="emptyText" />
      <el-card
        v-for="kb in items"
        :key="kb.externalId"
        class="knowledge-base-item"
      >
        <template #header>
          <div class="card-header">
            <span class="kb-name">{{ kb.name }}</span>
            <el-tag v-if="getTag(kb)" :type="getTag(kb).type" size="small">
              {{ getTag(kb).label }}
            </el-tag>
          </div>
        </template>

        <p class="kb-description">
          {{ kb.description || fallbackDescription }}
        </p>

        <div v-if="getMetaRows(kb).length > 0" class="kb-details">
          <div
            v-for="(row, index) in getMetaRows(kb)"
            :key="`${row.label}-${index}`"
            class="detail-item"
          >
            <span class="detail-label">{{ row.label }}:</span>
            <span class="detail-value">{{ row.value }}</span>
          </div>
        </div>

        <div class="kb-actions">
          <el-button size="small" type="primary" @click="handleOpen(kb)">
            {{ actionLabel }}
          </el-button>
        </div>
      </el-card>

      <div class="load-more" v-if="showLoadMore">
        <el-button @click="handleLoadMore" :loading="loading" plain>
          {{ loading ? loadingLabel : loadMoreLabel }}
        </el-button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue';
const props = defineProps({
  title: { type: String, default: '' },
  items: { type: Array, default: () => [] },
  emptyText: { type: String, default: '' },
  tagType: { type: String, default: '' },
  tagLabel: { type: String, default: '' },
  tagMapper: { type: Function, default: null },
  fallbackDescription: { type: String, default: '' },
  metaMapper: { type: Function, default: null },
  showLoadMore: { type: Boolean, default: false },
  loading: { type: Boolean, default: false },
  loadMoreLabel: { type: String, default: '' },
  loadingLabel: { type: String, default: '' },
  actionLabel: { type: String, default: '' }
});

const emit = defineEmits(['open', 'load-more']);

const showEmpty = computed(() => props.items.length === 0 && props.emptyText);

const getTag = (kb) => {
  if (props.tagMapper) {
    return props.tagMapper(kb);
  }
  if (props.tagType && props.tagLabel) {
    return { type: props.tagType, label: props.tagLabel };
  }
  return null;
};

const getMetaRows = (kb) => {
  if (props.metaMapper) {
    return props.metaMapper(kb) || [];
  }
  return [];
};

const handleOpen = (kb) => {
  emit('open', kb);
};

const handleLoadMore = () => {
  emit('load-more');
};
</script>

<style scoped>
.vertical-section {
  display: flex;
  flex-direction: column;
  background-color: var(--bg-white);
  border-radius: var(--border-radius-md);
  box-shadow: var(--shadow-sm);
  overflow: hidden;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: var(--spacing-md);
  border-bottom: 1px solid var(--border-color);
  background-color: var(--bg-white);
}

.section-title {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
  color: var(--text-primary);
}

.section-content {
  padding: var(--spacing-md);
}

.knowledge-base-item {
  margin-bottom: var(--spacing-md);
  transition: box-shadow 0.3s ease;
}

.knowledge-base-item:hover {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.kb-name {
  font-size: 16px;
  font-weight: 600;
  color: var(--text-dark);
}

.kb-description {
  margin: var(--spacing-sm) 0;
  color: var(--text-medium);
}

.kb-details {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-xs);
  margin-top: var(--spacing-sm);
}

.detail-item {
  display: flex;
  gap: var(--spacing-xs);
  color: var(--text-light);
  font-size: 12px;
}

.detail-label {
  color: var(--text-medium);
}

.kb-actions {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
  margin-top: 15px;
}

.load-more {
  text-align: center;
  margin-top: var(--spacing-md);
}

@media (max-width: 768px) {
  .section-header {
    flex-direction: column;
    align-items: stretch;
    gap: var(--spacing-md);
  }
}

.dark-mode .detail-label {
  color: var(--text-light);
}

.dark-mode .detail-value {
  color: var(--text-medium);
}

.dark-mode .knowledge-base-item {
  background-color: var(--bg-white);
  border: 1px solid var(--border-color);
}

.dark-mode .kb-name {
  color: var(--text-dark);
}

.dark-mode .kb-description {
  color: var(--text-medium);
}

.dark-mode .knowledge-base-item {
  box-shadow: var(--shadow-sm);
}
</style>
