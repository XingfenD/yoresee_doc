<template>
  <div class="template-picker-shell">
    <aside class="template-scope-nav">
      <button
        v-for="scope in scopeOptions"
        :key="scope.name"
        type="button"
        class="scope-item"
        :class="{ 'is-active': activeScope === scope.name }"
        @click="activeScope = scope.name"
      >
        {{ scope.label }}
      </button>
    </aside>

    <section class="template-main-panel">
      <div class="template-toolbar">
        <el-input
          v-model="keyword"
          clearable
          class="template-search-input"
          :placeholder="t('templates.searchPlaceholder')"
        >
          <template #prefix>
            <el-icon><Search /></el-icon>
          </template>
        </el-input>
      </div>

      <div class="template-content">
        <TemplatePickerPane
          :loading="currentLoading"
          :items="filteredTemplates"
          :selected-template-id="selectedTemplateIdForPane"
          :empty-text="currentEmptyText"
          :fallback-description="t('templates.noDescription')"
          layout="grid"
          @select="(tpl) => emit('select', tpl)"
        />
      </div>
    </section>
  </div>
</template>

<script setup>
import { Search } from '@element-plus/icons-vue';
import { useI18n } from 'vue-i18n';
import TemplatePickerPane from '@/components/template/TemplatePickerPane.vue';
import { useDocumentTemplatePicker } from '@/composables/template/useDocumentTemplatePicker';

const props = defineProps({
  visible: {
    type: Boolean,
    default: false
  },
  preferredScope: {
    type: String,
    default: ''
  },
  selectedTemplateId: {
    type: String,
    default: ''
  },
  documentType: {
    type: [String, Number],
    default: ''
  },
  knowledgeBaseId: {
    type: String,
    default: ''
  }
});

defineEmits(['select']);
const { t } = useI18n();

const {
  activeScope,
  keyword,
  scopeOptions,
  currentLoading,
  filteredTemplates,
  selectedTemplateIdForPane,
  currentEmptyText
} = useDocumentTemplatePicker(props);
</script>

<style scoped>
.template-picker-shell {
  display: grid;
  grid-template-columns: 180px minmax(0, 1fr);
  border: 1px solid var(--border-color);
  border-radius: 12px;
  overflow: hidden;
  height: 460px;
  width: 100%;
  background: var(--bg-white);
}

.template-scope-nav {
  border-right: 1px solid var(--border-color);
  padding: 12px 10px;
  display: flex;
  flex-direction: column;
  gap: 6px;
  background: color-mix(in srgb, var(--bg-white) 92%, #f3f5f9 8%);
}

.scope-item {
  border: 0;
  background: transparent;
  color: var(--text-secondary);
  border-radius: 8px;
  height: 36px;
  text-align: left;
  padding: 0 12px;
  font-size: 14px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.scope-item:hover {
  background: color-mix(in srgb, #3370ff 10%, transparent);
  color: var(--text-primary);
}

.scope-item.is-active {
  background: color-mix(in srgb, #3370ff 16%, transparent);
  color: #2f65e2;
  font-weight: 600;
}

.template-main-panel {
  padding: 12px;
  height: 100%;
  min-height: 0;
  display: flex;
  flex-direction: column;
}

.template-toolbar {
  margin-bottom: 12px;
}

.template-content {
  flex: 1;
  min-height: 0;
  overflow-y: auto;
  padding-right: 2px;
}

.template-search-input {
  max-width: 340px;
}

.dark-mode .template-picker-shell {
  background: #0f1218;
}

.dark-mode .template-scope-nav {
  background: #0b0f14;
}

.dark-mode .scope-item {
  color: #9ca3af;
}

.dark-mode .scope-item:hover {
  background: rgba(76, 141, 255, 0.14);
  color: #d1d5db;
}

.dark-mode .scope-item.is-active {
  background: rgba(76, 141, 255, 0.2);
  color: #88b0ff;
}

@media (max-width: 900px) {
  .template-picker-shell {
    grid-template-columns: 1fr;
  }

  .template-scope-nav {
    border-right: none;
    border-bottom: 1px solid var(--border-color);
    flex-direction: row;
    flex-wrap: wrap;
  }

  .template-search-input {
    max-width: none;
  }
}
</style>
