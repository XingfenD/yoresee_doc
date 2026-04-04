<template>
  <el-dialog
    v-model="visible"
    :title="t('user.invite.createTitle')"
    width="520px"
    :close-on-click-modal="false"
    @keydown.esc.prevent="handleCancel"
  >
    <el-form :model="formState" label-position="top" @submit.prevent>
      <el-form-item :label="t('user.invite.expiresAt')" required>
        <div class="invite-expire-row">
          <el-radio-group v-model="formState.expireType" class="invite-expire-options">
            <el-radio value="days">{{ t('user.invite.expireByDays') }}</el-radio>
            <el-radio value="date">{{ t('user.invite.expireByDate') }}</el-radio>
          </el-radio-group>
          <div class="invite-expire-input">
            <el-input-number
              v-if="formState.expireType === 'days'"
              v-model="formState.expireDays"
              :min="1"
              :max="365"
              style="width: 100%"
            />
            <el-date-picker
              v-else
              v-model="formState.expiresAt"
              type="date"
              value-format="YYYY-MM-DD"
              :placeholder="t('user.invite.expiresAt')"
              style="width: 100%"
            />
          </div>
        </div>
      </el-form-item>

      <el-form-item :label="t('user.invite.limitUsage')">
        <div class="invite-field-row invite-field-row--inline">
          <el-switch v-model="formState.limitEnabled" />
          <el-input-number
            v-model="formState.maxUsage"
            :min="1"
            :max="9999"
            :disabled="!formState.limitEnabled"
            style="width: 160px"
          />
        </div>
      </el-form-item>

      <el-form-item :label="t('user.invite.note')">
        <el-input
          v-model="formState.note"
          type="textarea"
          :rows="3"
          :placeholder="t('user.invite.notePlaceholder')"
        />
      </el-form-item>
    </el-form>

    <template #footer>
      <span class="dialog-footer">
        <el-button @click="handleCancel">{{ t('button.cancel') }}</el-button>
        <el-button type="primary" @click="handleSubmit">
          {{ t('button.confirm') }}
        </el-button>
      </span>
    </template>
  </el-dialog>
</template>

<script setup>
import { computed, reactive, watch } from 'vue';
import { ElMessage } from 'element-plus';
import { useI18n } from 'vue-i18n';

const props = defineProps({
  modelValue: {
    type: Boolean,
    default: false
  }
});

const emit = defineEmits(['update:modelValue', 'submit', 'cancel']);
const { t } = useI18n();

const visible = computed({
  get: () => props.modelValue,
  set: (value) => emit('update:modelValue', value)
});

const formState = reactive({
  expireType: 'days',
  expireDays: 7,
  expiresAt: '',
  limitEnabled: false,
  maxUsage: 10,
  note: ''
});

const resetForm = () => {
  formState.expireType = 'days';
  formState.expireDays = 7;
  formState.expiresAt = '';
  formState.limitEnabled = false;
  formState.maxUsage = 10;
  formState.note = '';
};

const handleCancel = () => {
  emit('cancel');
  visible.value = false;
};

const handleSubmit = () => {
  if (formState.expireType === 'date' && !formState.expiresAt) {
    ElMessage.error(t('user.invite.expiresAtRequired'));
    return;
  }
  if (formState.expireType === 'days' && (!formState.expireDays || formState.expireDays < 1)) {
    ElMessage.error(t('user.invite.expireDaysRequired'));
    return;
  }
  if (formState.limitEnabled && (!formState.maxUsage || formState.maxUsage < 1)) {
    ElMessage.error(t('user.invite.maxUsageRequired'));
    return;
  }
  emit('submit', {
    expire_type: formState.expireType,
    expire_days: formState.expireType === 'days' ? formState.expireDays : undefined,
    expires_at: formState.expireType === 'date' ? formState.expiresAt : undefined,
    limit_enabled: formState.limitEnabled,
    max_usage: formState.limitEnabled ? formState.maxUsage : undefined,
    note: formState.note?.trim() || undefined
  });
  visible.value = false;
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
.invite-field-row {
  margin-top: var(--spacing-sm);
}

.invite-expire-row {
  display: flex;
  align-items: center;
  gap: var(--spacing-md);
  margin-top: var(--spacing-xs);
}

.invite-expire-options {
  display: flex;
  align-items: center;
  gap: var(--spacing-md);
  flex-shrink: 0;
}

.invite-expire-input {
  flex: 1;
  min-width: 220px;
}

.invite-field-row--inline {
  display: flex;
  align-items: center;
  gap: var(--spacing-md);
}

.dark-mode :deep(.el-input-number__decrease),
.dark-mode :deep(.el-input-number__increase) {
  background-color: var(--bg-medium);
  color: var(--text-dark);
  border-color: var(--border-color);
}

.dark-mode :deep(.el-switch) {
  --el-switch-on-color: #3a7afe;
  --el-switch-off-color: #3a3a3a;
}

.dark-mode :deep(.el-date-editor),
.dark-mode :deep(.el-date-editor .el-input__wrapper) {
  background-color: var(--bg-medium);
  border-color: var(--border-color);
}

.dark-mode :deep(.el-date-editor .el-input__inner) {
  color: var(--text-dark);
}
</style>
