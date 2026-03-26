<template>
  <div class="card-list-section">
    <div class="section-header">
      <h3 class="section-title">{{ title }}</h3>
    </div>
    <div class="section-content">
      <el-empty v-if="showEmpty" :description="emptyText" />
      <el-card
        v-for="item in items"
        :key="resolveItemKey(item)"
        class="section-item"
      >
        <template #header>
          <div class="card-header">
            <span class="item-name">{{ resolveItemTitle(item) }}</span>
            <AppTag v-if="resolveTag(item)" :type="resolveTag(item).type" size="small">
              {{ resolveTag(item).label }}
            </AppTag>
          </div>
        </template>

        <p class="item-description">
          {{ resolveItemDescription(item) || fallbackDescription }}
        </p>

        <div v-if="resolveMetaRows(item).length > 0" class="item-meta">
          <div
            v-for="(row, index) in resolveMetaRows(item)"
            :key="`${row.label}-${index}`"
            class="meta-item"
          >
            <span class="meta-label">{{ row.label }}:</span>
            <span class="meta-value">{{ row.value }}</span>
          </div>
        </div>

        <div class="item-actions">
          <el-button size="small" type="primary" @click="$emit('open', item)">
            {{ actionLabel }}
          </el-button>
        </div>
      </el-card>

      <div v-if="showLoadMore" class="load-more">
        <el-button :loading="loading" plain @click="$emit('load-more')">
          {{ loading ? loadingLabel : loadMoreLabel }}
        </el-button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue';
import AppTag from '@/components/AppTag.vue';

const props = defineProps({
  title: { type: String, default: '' },
  items: { type: Array, default: () => [] },
  emptyText: { type: String, default: '' },
  fallbackDescription: { type: String, default: '' },
  actionLabel: { type: String, default: '' },
  tagType: { type: String, default: '' },
  tagLabel: { type: String, default: '' },
  tagMapper: { type: Function, default: null },
  metaMapper: { type: Function, default: null },
  showLoadMore: { type: Boolean, default: false },
  loading: { type: Boolean, default: false },
  loadMoreLabel: { type: String, default: '' },
  loadingLabel: { type: String, default: '' },
  itemKeyMapper: { type: Function, default: null },
  itemTitleMapper: { type: Function, default: null },
  itemDescriptionMapper: { type: Function, default: null }
});

defineEmits(['open', 'load-more']);

const showEmpty = computed(() => props.items.length === 0 && props.emptyText);

const resolveItemKey = (item) => {
  if (props.itemKeyMapper) return props.itemKeyMapper(item);
  return item?.externalId || item?.id || item?.external_id || item?.name;
};

const resolveItemTitle = (item) => {
  if (props.itemTitleMapper) return props.itemTitleMapper(item);
  return item?.name || '';
};

const resolveItemDescription = (item) => {
  if (props.itemDescriptionMapper) return props.itemDescriptionMapper(item);
  return item?.description || '';
};

const resolveTag = (item) => {
  if (props.tagMapper) return props.tagMapper(item);
  if (props.tagType && props.tagLabel) {
    return { type: props.tagType, label: props.tagLabel };
  }
  return null;
};

const resolveMetaRows = (item) => {
  if (!props.metaMapper) return [];
  return props.metaMapper(item) || [];
};
</script>

<style scoped>
.card-list-section {
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

.section-item {
  margin-bottom: var(--spacing-md);
  transition: box-shadow 0.3s ease;
}

.section-item:hover {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: var(--spacing-sm);
}

.item-name {
  font-size: 16px;
  font-weight: 600;
  color: var(--text-dark);
}

.item-description {
  margin: var(--spacing-sm) 0;
  color: var(--text-medium);
}

.item-meta {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-xs);
  margin-top: var(--spacing-sm);
}

.meta-item {
  display: flex;
  gap: var(--spacing-xs);
  font-size: 12px;
  color: var(--text-light);
}

.meta-label {
  color: var(--text-medium);
}

.item-actions {
  display: flex;
  justify-content: flex-end;
  margin-top: var(--spacing-sm);
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

.dark-mode .section-item {
  background-color: var(--bg-white);
  border: 1px solid var(--border-color);
}

.dark-mode .item-name {
  color: var(--text-dark);
}

.dark-mode .item-description {
  color: var(--text-medium);
}

.dark-mode .meta-label {
  color: var(--text-light);
}

.dark-mode .meta-value {
  color: var(--text-medium);
}
</style>
