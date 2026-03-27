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
            <TemplatePickerPane
              :loading="loadingRecent"
              :items="recentTemplates"
              :selected-template-id="formState.template"
              :empty-text="t('templates.noRecent')"
              :fallback-description="t('templates.noDescription')"
              @select="selectTemplate"
            />
          </el-tab-pane>
          <el-tab-pane :label="t('templates.my')" name="my">
            <TemplatePickerPane
              :loading="loadingMy"
              :items="myTemplates"
              :selected-template-id="formState.template"
              :empty-text="t('templates.noMy')"
              :fallback-description="t('templates.noDescription')"
              @select="selectTemplate"
            />
          </el-tab-pane>
          <el-tab-pane :label="t('templates.public')" name="public">
            <TemplatePickerPane
              :loading="loadingPublic"
              :items="publicTemplates"
              :selected-template-id="formState.template"
              :empty-text="t('templates.noPublic')"
              :fallback-description="t('templates.noDescription')"
              @select="selectTemplate"
            />
          </el-tab-pane>
          <el-tab-pane v-if="showKnowledgeBaseTemplates" :label="t('templates.knowledgeBaseTab')" name="knowledge_base">
            <TemplatePickerPane
              :loading="loadingKb"
              :items="kbTemplates"
              :selected-template-id="formState.template"
              :empty-text="t('templates.noKb')"
              :fallback-description="t('templates.noDescription')"
              @select="selectTemplate"
            />
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
import TemplatePickerPane from '@/components/TemplatePickerPane.vue';
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
</style>
