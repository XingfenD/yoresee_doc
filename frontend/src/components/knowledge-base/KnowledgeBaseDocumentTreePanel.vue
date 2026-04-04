<template>
  <section class="kb-tree-panel">
    <div class="section-header">
      <h3 class="section-title">{{ title }}</h3>

      <div class="tree-controls">
        <el-input
          v-model="searchValue"
          :placeholder="searchPlaceholder"
          prefix-icon="Search"
          clearable
          class="search-input"
        />

        <el-select v-model="sortValue" :placeholder="sortPlaceholder" class="sort-select">
          <el-option
            v-for="option in sortOptions"
            :key="option.value"
            :label="option.label"
            :value="option.value"
          />
        </el-select>
      </div>
    </div>

    <div class="tree-content" v-loading="loading">
      <DocumentTree
        v-if="nodes.length > 0"
        :nodes="nodes"
        :loading="loading"
        :show-toolbar="false"
        :show-create="false"
        :show-delete="false"
        :context-menu-enabled="false"
        @node-click="(data) => emit('node-click', data)"
      >
        <template #node-extra="{ data }">
          <AppTag v-if="data.tags && data.tags.length > 0" size="small" type="info" class="node-tag">
            {{ data.tags[0] }}
          </AppTag>
        </template>

        <template #node-actions="{ data }">
          <el-button size="small" type="primary" text @click.stop="emit('open', data)">
            {{ openLabel }}
          </el-button>

          <AppDropdown trigger="click" @command="(command) => emit('node-action', command, data)">
            <el-button size="small" text @click.stop>
              <el-icon>
                <MoreFilled />
              </el-icon>
            </el-button>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="rename">
                  {{ renameLabel }}
                </el-dropdown-item>
                <el-dropdown-item command="share" divided>
                  {{ shareLabel }}
                </el-dropdown-item>
                <el-dropdown-item command="delete" divided>
                  {{ deleteLabel }}
                </el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </AppDropdown>
        </template>
      </DocumentTree>

      <div v-else-if="!loading" class="empty-tree-state">
        <el-empty :description="emptyText" :image-size="64" />
      </div>
    </div>

    <div class="pagination-container" v-if="total > pageSize">
      <el-pagination
        :current-page="currentPage"
        :page-size="pageSize"
        :page-sizes="[20, 50, 100]"
        :total="total"
        layout="total, sizes, prev, pager, next, jumper"
        @size-change="(value) => emit('size-change', value)"
        @current-change="(value) => emit('current-change', value)"
      />
    </div>
  </section>
</template>

<script setup>
import { computed } from 'vue';
import { MoreFilled } from '@element-plus/icons-vue';
import DocumentTree from '@/components/document/DocumentTree.vue';
import AppTag from '@/components/base/AppTag.vue';
import AppDropdown from '@/components/base/AppDropdown.vue';

const props = defineProps({
  title: { type: String, default: '' },
  searchKeyword: { type: String, default: '' },
  searchPlaceholder: { type: String, default: '' },
  sortBy: { type: String, default: '' },
  sortPlaceholder: { type: String, default: '' },
  sortOptions: { type: Array, default: () => [] },
  nodes: { type: Array, default: () => [] },
  loading: { type: Boolean, default: false },
  emptyText: { type: String, default: '' },
  total: { type: Number, default: 0 },
  currentPage: { type: Number, default: 1 },
  pageSize: { type: Number, default: 50 },
  openLabel: { type: String, default: '' },
  renameLabel: { type: String, default: '' },
  shareLabel: { type: String, default: '' },
  deleteLabel: { type: String, default: '' }
});

const emit = defineEmits([
  'update:searchKeyword',
  'update:sortBy',
  'node-click',
  'open',
  'node-action',
  'size-change',
  'current-change'
]);

const searchValue = computed({
  get: () => props.searchKeyword,
  set: (value) => emit('update:searchKeyword', value)
});

const sortValue = computed({
  get: () => props.sortBy,
  set: (value) => emit('update:sortBy', value)
});
</script>

<style scoped>
.kb-tree-panel {
  display: flex;
  flex-direction: column;
  min-height: 0;
  background-color: var(--bg-white);
  border-radius: var(--border-radius-md);
  box-shadow: var(--shadow-sm);
  overflow: hidden;
  flex: 1;
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
  color: var(--text-dark);
}

.tree-controls {
  display: flex;
  gap: var(--spacing-md);
  align-items: center;
}

.search-input {
  width: 200px;
}

.sort-select {
  width: 150px;
}

.tree-content {
  flex: 1;
  min-height: 0;
  padding: var(--spacing-md);
  overflow-y: auto;
}

.empty-tree-state {
  display: flex;
  height: 100%;
  justify-content: center;
  align-items: center;
  background-color: var(--bg-white);
  border-radius: var(--border-radius-md);
  box-shadow: var(--shadow-sm);
  border: 1px solid var(--border-color);
}

.pagination-container {
  padding: var(--spacing-md);
  border-top: 1px solid var(--border-color);
  background-color: var(--bg-white);
  display: flex;
  justify-content: center;
  align-items: center;
}

@media (max-width: 768px) {
  .tree-controls {
    flex-direction: column;
    align-items: stretch;
  }

  .search-input,
  .sort-select {
    width: 100%;
  }
}
</style>
