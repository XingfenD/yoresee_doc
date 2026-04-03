<template>
  <el-dialog v-model="visible" :title="title" width="520px">
    <el-form label-position="top" class="template-settings-form">
      <el-form-item :label="t('templates.nameLabel')" required>
        <el-input
          v-model="formState.name"
          maxlength="100"
          show-word-limit
          :placeholder="t('templates.nameLabel')"
        />
      </el-form-item>
      <el-form-item :label="t('templates.descLabel')">
        <el-input
          v-model="formState.description"
          type="textarea"
          :rows="3"
          maxlength="300"
          show-word-limit
          :placeholder="t('templates.descLabel')"
        />
      </el-form-item>
      <el-form-item :label="t('templates.public')">
        <el-switch v-model="formState.is_public" />
      </el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="visible = false">{{ t('button.cancel') }}</el-button>
      <el-button type="primary" :loading="loading" @click="emit('submit')">
        {{ t('button.confirm') }}
      </el-button>
    </template>
  </el-dialog>
</template>

<script setup>
import { computed } from 'vue';
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
  title: {
    type: String,
    default: ''
  },
  formState: {
    type: Object,
    required: true
  }
});

const emit = defineEmits(['update:modelValue', 'submit']);
const { t } = useI18n();

const visible = computed({
  get: () => props.modelValue,
  set: (value) => emit('update:modelValue', value)
});
</script>

<style scoped>
.template-settings-form :deep(.el-switch) {
  vertical-align: middle;
}
</style>
