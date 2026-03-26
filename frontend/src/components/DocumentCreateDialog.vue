<template>
  <el-dialog v-model="visible" :title="t('knowledgeBase.createDocument')" width="560px" :close-on-click-modal="false"
    @keydown.esc.prevent="handleCancel">
    <el-form :model="formState" label-position="top" @submit.prevent @keydown.enter.prevent="handleCreate">
      <el-form-item :label="t('knowledgeBase.documentTitle')" required>
        <el-input v-model="formState.title" :placeholder="t('knowledgeBase.enterDocumentTitle')" maxlength="100"
          show-word-limit />
      </el-form-item>

      <el-form-item v-if="showTemplateSelector" :label="t('knowledgeBase.template')">
        <el-tabs v-model="activeTab" class="template-tabs">
          <el-tab-pane :label="t('templates.recent')" name="recent">
            <div v-loading="loadingRecent">
              <div v-if="recentTemplates.length === 0" class="template-empty">
                <el-empty :description="t('templates.noRecent')" />
              </div>
              <div v-else class="template-list">
                <div
                  v-for="tpl in recentTemplates"
                  :key="tpl.id"
                  class="template-card"
                  :class="{ 'is-selected': formState.template === String(tpl.id) }"
                  @click="selectTemplate(tpl)"
                >
                  <div class="template-card-title">
                    {{ tpl.name }}
                  </div>
                  <div class="template-card-desc">{{ tpl.description || t('templates.noDescription') }}</div>
                </div>
              </div>
            </div>
          </el-tab-pane>
          <el-tab-pane :label="t('templates.my')" name="my">
            <div v-loading="loadingMy">
              <div v-if="myTemplates.length === 0" class="template-empty">
                <el-empty :description="t('templates.noMy')" />
              </div>
              <div v-else class="template-list">
                <div
                  v-for="tpl in myTemplates"
                  :key="tpl.id"
                  class="template-card"
                  :class="{ 'is-selected': formState.template === String(tpl.id) }"
                  @click="selectTemplate(tpl)"
                >
                  <div class="template-card-title">
                    {{ tpl.name }}
                  </div>
                  <div class="template-card-desc">{{ tpl.description || t('templates.noDescription') }}</div>
                </div>
              </div>
            </div>
          </el-tab-pane>
          <el-tab-pane :label="t('templates.public')" name="public">
            <div v-loading="loadingPublic">
              <div v-if="publicTemplates.length === 0" class="template-empty">
                <el-empty :description="t('templates.noPublic')" />
              </div>
              <div v-else class="template-list">
                <div
                  v-for="tpl in publicTemplates"
                  :key="tpl.id"
                  class="template-card"
                  :class="{ 'is-selected': formState.template === String(tpl.id) }"
                  @click="selectTemplate(tpl)"
                >
                  <div class="template-card-title">
                    {{ tpl.name }}
                  </div>
                  <div class="template-card-desc">{{ tpl.description || t('templates.noDescription') }}</div>
                </div>
              </div>
            </div>
          </el-tab-pane>
          <el-tab-pane v-if="showKnowledgeBaseTemplates" :label="t('templates.knowledgeBaseTab')" name="knowledge_base">
            <div v-loading="loadingKb">
              <div v-if="kbTemplates.length === 0" class="template-empty">
                <el-empty :description="t('templates.noKb')" />
              </div>
              <div v-else class="template-list">
                <div
                  v-for="tpl in kbTemplates"
                  :key="tpl.id"
                  class="template-card"
                  :class="{ 'is-selected': formState.template === String(tpl.id) }"
                  @click="selectTemplate(tpl)"
                >
                  <div class="template-card-title">
                    {{ tpl.name }}
                  </div>
                  <div class="template-card-desc">{{ tpl.description || t('templates.noDescription') }}</div>
                </div>
              </div>
            </div>
          </el-tab-pane>
        </el-tabs>
      </el-form-item>
    </el-form>

    <template #footer>
      <span class="dialog-footer">
        <el-button @click="handleCancel">{{ t('button.cancel') }}</el-button>
        <el-button type="primary" :loading="loading" @click="handleCreate">
          {{ t('button.create') }}
        </el-button>
      </span>
    </template>
  </el-dialog>
</template>

<script setup>
import { computed, reactive, ref, watch } from 'vue';
import { ElMessage } from 'element-plus';
import { useI18n } from 'vue-i18n';
import { useTemplateCatalog } from '@/composables/useTemplateCatalog';

const props = defineProps({
  modelValue: {
    type: Boolean,
    default: false
  },
  loading: {
    type: Boolean,
    default: false
  },
  initialTitle: {
    type: String,
    default: ''
  },
  parentExternalId: {
    type: String,
    default: ''
  },
  knowledgeBaseId: {
    type: String,
    default: ''
  },
  initialTemplateId: {
    type: [String, Number],
    default: ''
  },
  initialTemplateMeta: {
    type: Object,
    default: null
  },
  showTemplateSelector: {
    type: Boolean,
    default: true
  }
});

const emit = defineEmits(['update:modelValue', 'submit', 'cancel']);
const { t } = useI18n();

const visible = computed({
  get: () => props.modelValue,
  set: (value) => emit('update:modelValue', value)
});

const formState = reactive({
  title: '',
  template: '',
  parentExternalId: '',
  templateMeta: null
});

const activeTab = ref('recent');
const showKnowledgeBaseTemplates = computed(() => !!props.knowledgeBaseId);
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
    console.error(`[DocumentCreateDialog] load ${scope} templates failed`, error);
  }
});

const resetForm = () => {
  formState.title = props.initialTitle || '';
  formState.template = props.initialTemplateId ? String(props.initialTemplateId) : '';
  formState.templateMeta = props.initialTemplateMeta || null;
  formState.parentExternalId = props.parentExternalId || '';
};

const handleCancel = () => {
  emit('cancel');
  visible.value = false;
};

const handleCreate = () => {
  if (!formState.title.trim()) {
    ElMessage.error(t('knowledgeBase.titleRequired'));
    return;
  }
  emit('submit', {
    title: formState.title.trim(),
    template: formState.template,
    template_meta: formState.templateMeta,
    parent_external_id: formState.parentExternalId || undefined
  });
};

const selectTemplate = (tpl) => {
  const templateId = String(tpl.id);
  if (formState.template === templateId) {
    formState.template = '';
    formState.templateMeta = null;
    return;
  }
  formState.template = templateId;
  formState.templateMeta = tpl;
};

const fetchTemplates = async (tab) => {
  if (!props.showTemplateSelector) {
    return;
  }
  if (tab === 'knowledge_base' && !showKnowledgeBaseTemplates.value) {
    return;
  }
  await ensureTemplateLoaded(tab);
};

watch(
  () => props.modelValue,
  (nextVisible) => {
    if (nextVisible) {
      resetForm();
      activeTab.value = 'recent';
      fetchTemplates('recent');
    }
  }
);

watch(
  () => props.knowledgeBaseId,
  (next, prev) => {
    if (next !== prev) {
      invalidateTemplateScope('knowledge_base');
      if (activeTab.value === 'knowledge_base' && props.modelValue) {
        fetchTemplates('knowledge_base');
      }
    }
  }
);

watch(activeTab, (tab) => {
  fetchTemplates(tab);
});
</script>

<style scoped>
.template-tabs {
  --tab-active: #3370ff;
  --tab-text: #4b5563;
  --tab-border: #e5e7eb;
}

.template-tabs :deep(.el-tabs__header) {
  margin: 0 0 var(--spacing-md);
  border-bottom: 1px solid var(--tab-border);
}

.template-tabs :deep(.el-tabs__nav-wrap) {
  padding: 0;
  margin-bottom: 0;
  border-bottom: none;
}

.template-tabs :deep(.el-tabs__nav) {
  gap: 24px;
}

.template-tabs :deep(.el-tabs__item) {
  height: 36px;
  line-height: 36px;
  padding: 0 2px;
  font-weight: 500;
  color: var(--tab-text);
}

.template-tabs :deep(.el-tabs__item.is-active) {
  color: var(--tab-active);
  font-weight: 600;
}

.template-tabs :deep(.el-tabs__active-bar) {
  height: 2px;
  background: var(--tab-active);
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

.dark-mode .template-tabs {
  --tab-active: #4c8dff;
  --tab-text: #9ca3af;
  --tab-border: #1f2937;
}

.dark-mode .template-tabs :deep(.el-tabs__item) {
  color: var(--tab-text);
}

.dark-mode .template-tabs :deep(.el-tabs__item.is-active) {
  color: var(--tab-active);
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
