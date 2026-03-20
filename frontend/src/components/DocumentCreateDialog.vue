<template>
  <el-dialog v-model="visible" :title="t('knowledgeBase.createDocument')" width="560px" :close-on-click-modal="false"
    @keydown.esc.prevent="handleCancel">
    <el-form :model="formState" label-position="top" @submit.prevent @keydown.enter.prevent="handleCreate">
      <el-form-item :label="t('knowledgeBase.documentTitle')" required>
        <el-input v-model="formState.title" :placeholder="t('knowledgeBase.enterDocumentTitle')" maxlength="100"
          show-word-limit />
      </el-form-item>

      <el-form-item :label="t('knowledgeBase.template')">
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
                  :class="{ 'is-selected': formState.template === tpl.id }"
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
                  :class="{ 'is-selected': formState.template === tpl.id }"
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
                  :class="{ 'is-selected': formState.template === tpl.id }"
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
                  :class="{ 'is-selected': formState.template === tpl.id }"
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
import { listTemplates } from '@/services/api';

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
const recentTemplates = ref([]);
const myTemplates = ref([]);
const publicTemplates = ref([]);
const kbTemplates = ref([]);
const loadingRecent = ref(false);
const loadingMy = ref(false);
const loadingPublic = ref(false);
const loadingKb = ref(false);
const loadedRecent = ref(false);
const loadedMy = ref(false);
const loadedPublic = ref(false);
const loadedKb = ref(false);

const showKnowledgeBaseTemplates = computed(() => !!props.knowledgeBaseId);

const resetForm = () => {
  formState.title = props.initialTitle || '';
  formState.template = '';
  formState.templateMeta = null;
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
  if (formState.template === tpl.id) {
    formState.template = '';
    formState.templateMeta = null;
    return;
  }
  formState.template = tpl.id;
  formState.templateMeta = tpl;
};

const fetchTemplates = async (tab) => {
  if (tab === 'recent') {
    if (loadedRecent.value) return;
    if (loadingRecent.value) return;
    loadingRecent.value = true;
    try {
      const data = await listTemplates({
        order_by: 'updated_at',
        order_desc: true,
        page: 1,
        page_size: 50
      });
      recentTemplates.value = data.templates || [];
      loadedRecent.value = true;
    } finally {
      loadingRecent.value = false;
    }
    return;
  }
  if (tab === 'my') {
    if (loadedMy.value) return;
    if (loadingMy.value) return;
    loadingMy.value = true;
    try {
      const data = await listTemplates({
        only_mine: true,
        target_container: 'own',
        order_by: 'updated_at',
        order_desc: true,
        page: 1,
        page_size: 50
      });
      myTemplates.value = data.templates || [];
      loadedMy.value = true;
    } finally {
      loadingMy.value = false;
    }
    return;
  }
  if (tab === 'public') {
    if (loadedPublic.value) return;
    if (loadingPublic.value) return;
    loadingPublic.value = true;
    try {
      const data = await listTemplates({
        target_container: 'public',
        order_by: 'updated_at',
        order_desc: true,
        page: 1,
        page_size: 50
      });
      publicTemplates.value = data.templates || [];
      loadedPublic.value = true;
    } finally {
      loadingPublic.value = false;
    }
    return;
  }
  if (tab === 'knowledge_base' && showKnowledgeBaseTemplates.value) {
    if (loadedKb.value) return;
    if (loadingKb.value) return;
    loadingKb.value = true;
    try {
      const data = await listTemplates({
        target_container: 'knowledge_base',
        knowledge_base_id: props.knowledgeBaseId,
        order_by: 'updated_at',
        order_desc: true,
        page: 1,
        page_size: 50
      });
      kbTemplates.value = data.templates || [];
      loadedKb.value = true;
    } finally {
      loadingKb.value = false;
    }
  }
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
      loadedKb.value = false;
      kbTemplates.value = [];
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
