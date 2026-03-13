<template>
  <div class="template-section">
    <div class="section-header">
      <h3 class="section-title">{{ title }}</h3>
    </div>
    <div class="section-content">
      <el-empty v-if="showEmpty" :description="emptyText" />
      <el-card v-for="tpl in items" :key="tpl.id" class="template-item">
        <template #header>
          <div class="card-header">
            <span class="template-name">{{ tpl.name }}</span>
            <el-tag v-if="tagMapper" :type="tagMapper(tpl).type" size="small">
              {{ tagMapper(tpl).label }}
            </el-tag>
          </div>
        </template>

        <p class="template-description">
          {{ tpl.description || fallbackDescription }}
        </p>

        <div class="template-meta" v-if="metaMapper">
          <div
            v-for="(row, index) in metaMapper(tpl)"
            :key="`${row.label}-${index}`"
            class="meta-item"
          >
            <span class="meta-label">{{ row.label }}:</span>
            <span class="meta-value">{{ row.value }}</span>
          </div>
        </div>

        <div class="template-actions">
          <el-button size="small" type="primary" @click="$emit('open', tpl)">
            {{ actionLabel }}
          </el-button>
        </div>
      </el-card>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue';

const props = defineProps({
  title: { type: String, default: '' },
  items: { type: Array, default: () => [] },
  emptyText: { type: String, default: '' },
  fallbackDescription: { type: String, default: '' },
  actionLabel: { type: String, default: '' },
  tagMapper: { type: Function, default: null },
  metaMapper: { type: Function, default: null }
});

defineEmits(['open']);

const showEmpty = computed(() => props.items.length === 0 && props.emptyText);
</script>

<style scoped>
.template-section {
  display: flex;
  flex-direction: column;
  background-color: var(--bg-white);
  border-radius: var(--border-radius-md);
  box-shadow: var(--shadow-sm);
  overflow: hidden;
}

.section-header {
  display: flex;
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

.template-item {
  margin-bottom: var(--spacing-md);
  transition: box-shadow 0.3s ease;
}

.template-item:hover {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: var(--spacing-sm);
}

.template-name {
  font-size: 16px;
  font-weight: 600;
  color: var(--text-dark);
}

.template-description {
  margin: var(--spacing-sm) 0;
  color: var(--text-medium);
}

.template-meta {
  border-top: 1px solid var(--border-color);
  padding-top: var(--spacing-sm);
  margin-top: var(--spacing-sm);
  display: flex;
  flex-direction: column;
  gap: var(--spacing-xs);
}

.meta-item {
  display: flex;
  justify-content: space-between;
  font-size: 12px;
  color: var(--text-light);
}

.meta-label {
  color: var(--text-medium);
}

.template-actions {
  display: flex;
  justify-content: flex-end;
  margin-top: var(--spacing-sm);
}

.dark-mode .template-name {
  color: var(--text-dark);
}

.dark-mode .template-description {
  color: var(--text-medium);
}

.dark-mode .template-item {
  background-color: var(--bg-white);
  border: 1px solid var(--border-color);
}
</style>
