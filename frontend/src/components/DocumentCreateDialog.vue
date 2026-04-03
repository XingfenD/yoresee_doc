<template>
  <el-dialog
    v-model="visible"
    :title="t('knowledgeBase.createDocument')"
    :width="dialogWidth"
    class="document-create-dialog"
    :close-on-click-modal="false"
    @keydown.esc.prevent="handleCancel"
  >
    <el-form
      :model="formState"
      label-position="top"
      class="create-form"
      @submit.prevent
      @keydown.enter.prevent="handleCreate"
    >
      <el-form-item v-if="showTitleInput" :label="t('knowledgeBase.documentTitle')" required>
        <el-input
          v-model="formState.title"
          :placeholder="t('knowledgeBase.enterDocumentTitle')"
          maxlength="100"
          show-word-limit
        />
      </el-form-item>

      <el-form-item v-if="showPublicSwitch" :label="t('document.settings.publicLabel')">
        <el-switch v-model="formState.isPublic" />
      </el-form-item>

      <el-form-item v-if="showLocationSelector" :label="t('document.createLocationLabel')" required>
        <el-radio-group v-model="formState.containerType" class="location-segment">
          <el-radio-button
            v-for="item in locationOptions"
            :key="item.value"
            :label="item.value"
          >
            {{ item.label }}
          </el-radio-button>
        </el-radio-group>
        <el-select
          v-if="formState.containerType === 'knowledge_base'"
          v-model="formState.targetKnowledgeBaseId"
          class="location-kb-select"
          :placeholder="t('document.createLocationKnowledgeBase')"
          :loading="loadingKnowledgeBases"
          filterable
        >
          <el-option
            v-for="item in knowledgeBaseOptions"
            :key="item.value"
            :label="item.label"
            :value="item.value"
          />
        </el-select>
      </el-form-item>

      <el-form-item v-if="showTemplatePicker" :label="t('knowledgeBase.template')" class="template-picker-field">
        <DocumentTemplatePicker
          :visible="visible"
          :preferred-scope="formState.templateMeta?.scope || ''"
          :selected-template-id="formState.template"
          :knowledge-base-id="knowledgeBaseId"
          @select="selectTemplate"
        />
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
import DocumentTemplatePicker from '@/components/DocumentTemplatePicker.vue';
import { getKnowledgeBases } from '@/services/api';

const props = defineProps({
  modelValue: {
    type: Boolean,
    default: false
  },
  loading: {
    type: Boolean,
    default: false
  },
  showTemplatePicker: {
    type: Boolean,
    default: true
  },
  showTitleInput: {
    type: Boolean,
    default: false
  },
  showPublicSwitch: {
    type: Boolean,
    default: false
  },
  showLocationSelector: {
    type: Boolean,
    default: false
  },
  initialTitle: {
    type: String,
    default: ''
  },
  initialIsPublic: {
    type: Boolean,
    default: false
  },
  initialContainerType: {
    type: String,
    default: 'own'
  },
  initialTargetKnowledgeBaseId: {
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
    type: [String, Number, BigInt],
    default: ''
  },
  initialTemplateMeta: {
    type: Object,
    default: null
  }
});

const emit = defineEmits(['update:modelValue', 'submit', 'cancel']);
const { t } = useI18n();

const dialogWidth = computed(() => {
  if (props.showTemplatePicker) {
    return '980px';
  }
  return '620px';
});

const visible = computed({
  get: () => props.modelValue,
  set: (value) => emit('update:modelValue', value)
});

const loadingKnowledgeBases = ref(false);
const knowledgeBaseOptions = ref([]);

const locationOptions = computed(() => [
  { label: t('document.createLocationOwn'), value: 'own' },
  { label: t('document.createLocationKnowledgeBase'), value: 'knowledge_base' }
]);

const formState = reactive({
  title: '',
  isPublic: false,
  containerType: 'own',
  targetKnowledgeBaseId: '',
  template: '',
  parentExternalId: '',
  templateMeta: null
});

const loadKnowledgeBaseOptions = async () => {
  if (!props.showLocationSelector) {
    return;
  }
  loadingKnowledgeBases.value = true;
  try {
    const resp = await getKnowledgeBases({
      only_mine: true,
      page: 1,
      page_size: 200,
      order_by: 'updated_at',
      order_desc: true
    });
    knowledgeBaseOptions.value = (resp.knowledge_bases || []).map((item) => ({
      value: item.external_id,
      label: item.name || item.external_id
    }));
  } catch (error) {
    console.error('[DocumentCreateDialog] load knowledge base options failed', error);
    knowledgeBaseOptions.value = [];
  } finally {
    loadingKnowledgeBases.value = false;
  }
};

const resetForm = () => {
  formState.title = props.initialTitle || t('document.untitledDefaultTitle');
  formState.isPublic = Boolean(props.initialIsPublic);
  formState.containerType = props.initialContainerType === 'knowledge_base' ? 'knowledge_base' : 'own';
  formState.targetKnowledgeBaseId =
    props.initialTargetKnowledgeBaseId || props.knowledgeBaseId || '';
  const hasInitialTemplateId =
    props.initialTemplateId !== '' &&
    props.initialTemplateId !== null &&
    props.initialTemplateId !== undefined;
  formState.template = hasInitialTemplateId ? String(props.initialTemplateId) : '';
  formState.templateMeta = props.initialTemplateMeta || null;
  formState.parentExternalId = props.parentExternalId || '';
};

const handleCancel = () => {
  emit('cancel');
  visible.value = false;
};

const handleCreate = () => {
  if (props.showTitleInput && !formState.title.trim()) {
    ElMessage.error(t('knowledgeBase.titleRequired'));
    return;
  }
  if (
    props.showLocationSelector &&
    formState.containerType === 'knowledge_base' &&
    !formState.targetKnowledgeBaseId
  ) {
    ElMessage.error(t('knowledgeBase.selectKnowledgeBase'));
    return;
  }

  const title = formState.title.trim() || t('document.untitledDefaultTitle');

  emit('submit', {
    title,
    template: formState.template,
    template_meta: formState.templateMeta,
    parent_external_id: formState.parentExternalId || undefined,
    is_public: props.showPublicSwitch ? Boolean(formState.isPublic) : false,
    container_type: props.showLocationSelector ? formState.containerType : undefined,
    knowledge_base_external_id:
      props.showLocationSelector && formState.containerType === 'knowledge_base'
        ? formState.targetKnowledgeBaseId || undefined
        : undefined
  });
};

const selectTemplate = (tpl) => {
  if (tpl?.is_blank) {
    formState.template = '';
    formState.templateMeta = null;
    return;
  }
  const templateId = String(tpl.id);
  if (formState.template === templateId) {
    formState.template = '';
    formState.templateMeta = null;
    return;
  }
  formState.template = templateId;
  formState.templateMeta = tpl;
};

watch(
  () => props.modelValue,
  async (nextVisible) => {
    if (nextVisible) {
      resetForm();
      await loadKnowledgeBaseOptions();
    }
  }
);
</script>

<style scoped>
.document-create-dialog :deep(.el-dialog__body) {
  padding-top: 10px;
}

.create-form {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.template-picker-field {
  margin-bottom: 0;
}

.location-segment {
  display: flex;
  width: 100%;
}

.location-segment :deep(.el-radio-button) {
  flex: 1;
}

.location-segment :deep(.el-radio-button__inner) {
  width: 100%;
  text-align: center;
}

.location-kb-select {
  width: 100%;
  margin-top: 10px;
}
</style>
