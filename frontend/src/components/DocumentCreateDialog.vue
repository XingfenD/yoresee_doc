<template>
  <el-dialog
    v-model="visible"
    :title="t('knowledgeBase.createDocument')"
    width="980px"
    class="document-create-dialog"
    :close-on-click-modal="false"
    @keydown.esc.prevent="handleCancel"
  >
    <DocumentTemplatePicker
      :visible="visible"
      :selected-template-id="formState.template"
      :knowledge-base-id="knowledgeBaseId"
      @select="selectTemplate"
    />

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
import { computed, reactive, watch } from 'vue';
import { useI18n } from 'vue-i18n';
import DocumentTemplatePicker from '@/components/DocumentTemplatePicker.vue';

const props = defineProps({
  modelValue: {
    type: Boolean,
    default: false
  },
  loading: {
    type: Boolean,
    default: false
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
  }
});

const emit = defineEmits(['update:modelValue', 'submit', 'cancel']);
const { t } = useI18n();

const visible = computed({
  get: () => props.modelValue,
  set: (value) => emit('update:modelValue', value)
});

const formState = reactive({
  template: '',
  parentExternalId: '',
  templateMeta: null
});

const resetForm = () => {
  formState.template = props.initialTemplateId ? String(props.initialTemplateId) : '';
  formState.templateMeta = props.initialTemplateMeta || null;
  formState.parentExternalId = props.parentExternalId || '';
};

const handleCancel = () => {
  emit('cancel');
  visible.value = false;
};

const handleCreate = () => {
  emit('submit', {
    title: t('document.untitledDefaultTitle'),
    template: formState.template,
    template_meta: formState.templateMeta,
    parent_external_id: formState.parentExternalId || undefined,
    is_public: false
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
  (nextVisible) => {
    if (nextVisible) {
      resetForm();
    }
  }
);
</script>

<style scoped>
.document-create-dialog :deep(.el-dialog__body) {
  padding-top: 10px;
}
</style>
