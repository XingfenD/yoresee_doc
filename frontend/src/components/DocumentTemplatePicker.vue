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
import { computed, nextTick, ref, watch } from 'vue';
import { useI18n } from 'vue-i18n';
import TemplatePickerPane from '@/components/TemplatePickerPane.vue';
import { useTemplateCatalog } from '@/composables/useTemplateCatalog';

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
  knowledgeBaseId: {
    type: String,
    default: ''
  }
});

const emit = defineEmits(['select']);
const { t } = useI18n();

const activeScope = ref('recent');
const keyword = ref('');
const showKnowledgeBaseTemplates = computed(() => Boolean(props.knowledgeBaseId));

const normalizeScopeToTab = (scope) => {
  if (scope === 'system') {
    return 'public';
  }
  if (scope === 'private') {
    return 'my';
  }
  if (scope === 'knowledge_base') {
    return 'knowledge_base';
  }
  return 'recent';
};

const {
  recentTemplates,
  myTemplates,
  publicTemplates,
  kbTemplates,
  loadingRecent,
  loadingMy,
  loadingPublic,
  loadingKb,
  ensureLoaded: ensureTemplateLoaded,
  invalidateScope: invalidateTemplateScope
} = useTemplateCatalog({
  includeKnowledgeBase: true,
  knowledgeBaseId: computed(() => props.knowledgeBaseId || ''),
  onError: (error, scope) => {
    console.error(`[DocumentTemplatePicker] load ${scope} templates failed`, error);
  }
});

const scopeOptions = computed(() => {
  const base = [
    { name: 'recent', label: t('templates.recent') },
    { name: 'my', label: t('templates.my') },
    { name: 'public', label: t('templates.public') }
  ];
  if (showKnowledgeBaseTemplates.value) {
    base.push({ name: 'knowledge_base', label: t('templates.knowledgeBaseTab') });
  }
  return base;
});

const blankTemplateOption = computed(() => ({
  id: '__blank__',
  name: t('templates.blankDocument'),
  description: t('templates.blankDocumentDesc'),
  tags: [],
  is_blank: true
}));

const currentTemplates = computed(() => {
  if (activeScope.value === 'my') {
    return myTemplates.value;
  }
  if (activeScope.value === 'public') {
    return publicTemplates.value;
  }
  if (activeScope.value === 'knowledge_base') {
    return kbTemplates.value;
  }
  return recentTemplates.value;
});

const currentLoading = computed(() => {
  if (activeScope.value === 'my') {
    return loadingMy.value;
  }
  if (activeScope.value === 'public') {
    return loadingPublic.value;
  }
  if (activeScope.value === 'knowledge_base') {
    return loadingKb.value;
  }
  return loadingRecent.value;
});

const currentEmptyText = computed(() => {
  if (activeScope.value === 'my') {
    return t('templates.noMy');
  }
  if (activeScope.value === 'public') {
    return t('templates.noPublic');
  }
  if (activeScope.value === 'knowledge_base') {
    return t('templates.noKb');
  }
  return t('templates.noRecent');
});

const filteredTemplates = computed(() => {
  const text = keyword.value.trim().toLowerCase();
  const blank = blankTemplateOption.value;
  if (!text) {
    return [blank, ...currentTemplates.value];
  }
  const list = currentTemplates.value.filter((item) => {
    const name = String(item?.name || '').toLowerCase();
    const description = String(item?.description || '').toLowerCase();
    const tags = Array.isArray(item?.tags) ? item.tags.join(' ').toLowerCase() : '';
    return name.includes(text) || description.includes(text) || tags.includes(text);
  });
  const blankMatched = [blank.name, blank.description]
    .join(' ')
    .toLowerCase()
    .includes(text);
  if (blankMatched) {
    return [blank, ...list];
  }
  return list;
});

const selectedTemplateIdForPane = computed(() => props.selectedTemplateId || '__blank__');

const initialScope = computed(() => {
  const mapped = normalizeScopeToTab(props.preferredScope);
  if (mapped === 'knowledge_base' && !showKnowledgeBaseTemplates.value) {
    return 'recent';
  }
  return mapped;
});

const fetchTemplates = async (scope) => {
  if (!props.visible) {
    return;
  }
  if (scope === 'knowledge_base' && !showKnowledgeBaseTemplates.value) {
    return;
  }
  await ensureTemplateLoaded(scope);
};

watch(
  () => props.visible,
  async (nextVisible) => {
    if (!nextVisible) {
      return;
    }
    // Wait one tick so parent dialog can finish resetForm and propagate preferredScope/template id.
    await nextTick();
    activeScope.value = initialScope.value;
    keyword.value = '';
    fetchTemplates(activeScope.value);
  }
);

watch(
  () => [props.visible, props.preferredScope, props.selectedTemplateId],
  ([visible]) => {
    if (!visible) {
      return;
    }
    if (!props.selectedTemplateId) {
      return;
    }
    const preferred = initialScope.value;
    if (preferred === 'recent') {
      return;
    }
    // If scope arrives late, auto-correct only when user still stays on default tab.
    if (activeScope.value === 'recent') {
      activeScope.value = preferred;
      fetchTemplates(preferred);
    }
  }
);

watch(
  () => props.knowledgeBaseId,
  (next, prev) => {
    if (next === prev) {
      return;
    }
    invalidateTemplateScope('knowledge_base');
    if (activeScope.value === 'knowledge_base' && props.visible) {
      fetchTemplates('knowledge_base');
    }
  }
);

watch(showKnowledgeBaseTemplates, (show) => {
  if (!show && activeScope.value === 'knowledge_base') {
    activeScope.value = 'recent';
  }
});

watch(activeScope, (scope) => {
  fetchTemplates(scope);
});
</script>

<style scoped>
.template-picker-shell {
  display: grid;
  grid-template-columns: 180px minmax(0, 1fr);
  border: 1px solid var(--border-color);
  border-radius: 12px;
  overflow: hidden;
  height: 460px;
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
