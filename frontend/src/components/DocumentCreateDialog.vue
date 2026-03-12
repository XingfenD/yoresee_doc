<template>
  <el-dialog v-model="visible" :title="t('knowledgeBase.createDocument')" width="500px" :close-on-click-modal="false"
    @keydown.esc.prevent="handleCancel">
    <el-form :model="formState" label-position="top" @submit.prevent @keydown.enter.prevent="handleCreate">
      <el-form-item :label="t('knowledgeBase.documentTitle')" required>
        <el-input v-model="formState.title" :placeholder="t('knowledgeBase.enterDocumentTitle')" maxlength="100"
          show-word-limit />
      </el-form-item>

      <el-form-item :label="t('knowledgeBase.documentType')">
        <el-select v-model="formState.type" :placeholder="t('knowledgeBase.selectDocumentType')" style="width: 100%">
          <el-option label="Markdown" value="markdown" />
          <el-option label="Text" value="text" />
        </el-select>
      </el-form-item>

      <el-form-item :label="t('knowledgeBase.template')">
        <el-select v-model="formState.template" :placeholder="t('knowledgeBase.selectTemplate')" style="width: 100%"
          disabled>
          <el-option label="空白文档" value="" />
        </el-select>
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
import { computed, reactive, watch } from 'vue';
import { useI18n } from 'vue-i18n';

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
  type: 'markdown',
  template: '',
  parentExternalId: ''
});

const resetForm = () => {
  formState.title = props.initialTitle || '';
  formState.type = 'markdown';
  formState.template = '';
  formState.parentExternalId = props.parentExternalId || '';
};

const handleCancel = () => {
  emit('cancel');
  visible.value = false;
};

const handleCreate = () => {
  emit('submit', {
    title: formState.title.trim(),
    type: formState.type,
    template: formState.template,
    parent_external_id: formState.parentExternalId || undefined
  });
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
