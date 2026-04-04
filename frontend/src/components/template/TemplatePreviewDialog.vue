<template>
  <el-dialog
    v-model="visible"
    :title="title"
    width="840px"
    append-to-body
    @closed="emit('closed')"
  >
    <div class="template-preview-dialog-body">
      <div v-if="!content" class="template-preview-empty">
        <el-empty :description="t('templates.contentEmpty')" />
      </div>
      <DocumentPreviewViewer
        v-else
        class="template-preview-render"
        :content="content"
        :document-type="documentType"
        :is-template="true"
        :is-dark-mode="isDarkMode"
      />
    </div>
    <template #footer>
      <el-button @click="visible = false">{{ t('button.cancel') }}</el-button>
    </template>
  </el-dialog>
</template>

<script setup>
import { computed } from 'vue';
import { useI18n } from 'vue-i18n';
import DocumentPreviewViewer from '@/components/document/render/DocumentPreviewViewer.vue';

const props = defineProps({
  modelValue: {
    type: Boolean,
    default: false
  },
  title: {
    type: String,
    default: ''
  },
  content: {
    type: String,
    default: ''
  },
  documentType: {
    type: [String, Number],
    default: '1'
  },
  isDarkMode: {
    type: Boolean,
    default: false
  }
});

const emit = defineEmits(['update:modelValue', 'closed']);
const { t } = useI18n();

const visible = computed({
  get: () => props.modelValue,
  set: (value) => emit('update:modelValue', value)
});
</script>

<style scoped>
.template-preview-dialog-body {
  min-height: 360px;
  max-height: 64vh;
  overflow: auto;
  border: 1px solid var(--border-color);
  border-radius: 8px;
  padding: 12px;
  background: var(--bg-white);
}

.template-preview-empty {
  min-height: 280px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.template-preview-render {
  min-height: 360px;
}

.dark-mode .template-preview-dialog-body {
  background: #141a22;
  border-color: #2a313a;
}
</style>
