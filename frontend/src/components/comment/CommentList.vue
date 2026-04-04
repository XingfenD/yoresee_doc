<template>
  <div class="comment-list-panel">
    <div v-if="showHeader" class="comment-list-header">
      <div class="comment-list-title">{{ title }}</div>
      <div class="comment-list-actions">
        <slot name="header-actions" />
      </div>
    </div>
    <div class="comment-list-body">
      <div v-if="!items || items.length === 0" class="comment-list-empty">
        {{ emptyText }}
      </div>
      <div v-else class="comment-list-items">
        <template v-for="item in items" :key="resolveKey(item)">
          <slot name="item" :item="item" />
        </template>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, useSlots } from 'vue';

const props = defineProps({
  title: {
    type: String,
    default: ''
  },
  emptyText: {
    type: String,
    default: ''
  },
  items: {
    type: Array,
    default: () => []
  },
  keyField: {
    type: String,
    default: 'external_id'
  },
  showTitle: {
    type: Boolean,
    default: true
  }
});

const slots = useSlots();
const showHeader = computed(() => props.showTitle || !!slots['header-actions']);

const resolveKey = (item) => {
  if (!item) return Math.random().toString(36);
  if (props.keyField && item[props.keyField] !== undefined) {
    return item[props.keyField];
  }
  return item.external_id || item.id || item.anchor_id || Math.random().toString(36);
};
</script>

<style scoped>
.comment-list-panel {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.comment-list-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.comment-list-title {
  font-size: 14px;
  font-weight: 600;
  color: var(--text-dark);
}

.comment-list-actions {
  display: flex;
  align-items: center;
  gap: 8px;
}

.comment-list-body {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.comment-list-empty {
  font-size: 13px;
  color: var(--text-light);
}

.comment-list-items {
  display: flex;
  flex-direction: column;
  gap: 12px;
}
</style>
