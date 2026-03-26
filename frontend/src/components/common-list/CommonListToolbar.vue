<template>
  <div v-if="showTitleBar || showSearch" class="list-titlebar">
    <div class="list-title">
      <slot name="title">
        {{ title }}
      </slot>
    </div>
    <div class="list-title-actions">
      <slot name="toolbar-actions" />
      <el-input
        v-if="showSearch"
        v-model="innerSearchQuery"
        :placeholder="searchPlaceholder"
        clearable
        class="list-search"
        @clear="emitSearch"
        @input="emitSearch"
      />
      <slot name="toolbar-right" />
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue';

const props = defineProps({
  showTitleBar: { type: Boolean, default: false },
  showSearch: { type: Boolean, default: false },
  title: { type: String, default: '' },
  searchPlaceholder: { type: String, default: '' },
  searchQuery: { type: String, default: '' }
});

const emit = defineEmits(['update:searchQuery', 'search']);

const innerSearchQuery = computed({
  get: () => props.searchQuery,
  set: (value) => emit('update:searchQuery', value)
});

const emitSearch = () => {
  emit('search', innerSearchQuery.value);
};
</script>

<style scoped>
.list-titlebar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  padding: 10px 12px;
  border-bottom: 1px solid var(--border-color);
  background: var(--bg-white);
}

.list-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 17px;
  font-weight: 600;
  color: var(--text-dark);
  line-height: 1.2;
}

.list-title-actions {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-left: auto;
}

.list-search {
  max-width: 320px;
  flex: 1;
  margin-top: -2px;
}

:global(.common-list.is-dark) .list-titlebar {
  background: #161b22;
  border-bottom-color: #2a313a;
}
</style>
