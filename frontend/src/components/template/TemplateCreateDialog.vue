<template>
  <el-dialog
    v-model="visible"
    :title="title"
    width="520px"
    :close-on-click-modal="false"
    @keydown.esc.prevent="handleCancel"
  >
    <el-form label-position="top" :model="formState" @submit.prevent>
      <el-form-item :label="t('templates.nameLabel')" required>
        <el-input v-model="formState.name" maxlength="100" />
      </el-form-item>
      <el-form-item :label="t('templates.descLabel')">
        <el-input v-model="formState.description" type="textarea" :rows="3" maxlength="255" />
      </el-form-item>
      <el-form-item :label="t('templates.scopeLabel')">
        <el-select v-model="formState.scope" style="width: 100%">
          <el-option value="own" :label="t('templates.scopeOwn')" />
          <el-option v-if="showKbScope" value="knowledge_base" :label="t('templates.scopeKb')" />
          <el-option value="public" :label="t('templates.scopePublic')" />
        </el-select>
      </el-form-item>
      <el-form-item :label="t('templates.tagsLabel')">
        <el-input v-model="formState.tags" :placeholder="t('templates.tagsPlaceholder')" />
      </el-form-item>
      <el-form-item v-if="showContent" :label="t('templates.contentLabel')" required>
        <el-input
          v-model="formState.content"
          type="textarea"
          :rows="6"
          :placeholder="t('templates.contentPlaceholder')"
        />
      </el-form-item>
    </el-form>

    <template #footer>
      <el-button @click="handleCancel">{{ t('button.cancel') }}</el-button>
      <el-button type="primary" :loading="loading" @click="handleSubmit">
        {{ t('button.confirm') }}
      </el-button>
    </template>
  </el-dialog>
</template>

<script setup>
import { computed, reactive, watch } from 'vue';
import { ElMessage } from 'element-plus';
import { useI18n } from 'vue-i18n';

const props = defineProps({
  modelValue: { type: Boolean, default: false },
  loading: { type: Boolean, default: false },
  title: { type: String, default: '' },
  showContent: { type: Boolean, default: true },
  showKbScope: { type: Boolean, default: true },
  initialName: { type: String, default: '' },
  initialDescription: { type: String, default: '' },
  initialScope: { type: String, default: 'own' },
  initialTags: { type: [String, Array], default: '' },
  initialContent: { type: String, default: '' }
});

const emit = defineEmits(['update:modelValue', 'submit', 'cancel']);
const { t } = useI18n();

const visible = computed({
  get: () => props.modelValue,
  set: (value) => emit('update:modelValue', value)
});

const formState = reactive({
  name: '',
  description: '',
  scope: 'own',
  tags: '',
  content: ''
});

const normalizeTags = (tags) => {
  if (Array.isArray(tags)) {
    return tags.join(', ');
  }
  return tags || '';
};

const resetForm = () => {
  formState.name = props.initialName || '';
  formState.description = props.initialDescription || '';
  formState.scope = props.initialScope || 'own';
  formState.tags = normalizeTags(props.initialTags);
  formState.content = props.initialContent || '';
};

const handleCancel = () => {
  emit('cancel');
  visible.value = false;
};

const handleSubmit = () => {
  if (!formState.name.trim()) {
    ElMessage.error(t('templates.nameRequired'));
    return;
  }
  if (props.showContent && !formState.content.trim()) {
    ElMessage.error(t('templates.emptyContent'));
    return;
  }

  const tags = formState.tags
    ? formState.tags.split(',').map((tag) => tag.trim()).filter(Boolean)
    : [];

  emit('submit', {
    name: formState.name.trim(),
    description: formState.description || '',
    scope: formState.scope,
    tags,
    content: formState.content
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
