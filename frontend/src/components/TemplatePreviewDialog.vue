<template>
  <el-dialog
    v-model="visible"
    :title="title"
    width="840px"
    @closed="emit('closed')"
  >
    <div class="template-preview-dialog-body">
      <div v-if="!content" class="template-preview-empty">
        <el-empty :description="t('templates.contentEmpty')" />
      </div>
      <div v-else ref="previewRef" class="template-preview-render"></div>
    </div>
    <template #footer>
      <el-button @click="visible = false">{{ t('button.cancel') }}</el-button>
    </template>
  </el-dialog>
</template>

<script setup>
import { computed, nextTick, ref, watch } from 'vue';
import { useI18n } from 'vue-i18n';
import Vditor from 'vditor';
import 'vditor/dist/index.css';

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
  isDarkMode: {
    type: Boolean,
    default: false
  }
});

const emit = defineEmits(['update:modelValue', 'closed']);
const { t } = useI18n();
const previewRef = ref(null);

const visible = computed({
  get: () => props.modelValue,
  set: (value) => emit('update:modelValue', value)
});

const renderPreviewContent = async () => {
  if (!props.modelValue) {
    return;
  }
  await nextTick();
  if (!previewRef.value) {
    return;
  }
  if (!props.content) {
    previewRef.value.innerHTML = '';
    return;
  }
  await Vditor.preview(previewRef.value, props.content, {
    mode: props.isDarkMode ? 'dark' : 'light',
    theme: {
      current: props.isDarkMode ? 'dark' : 'light'
    },
    hljs: {
      style: props.isDarkMode ? 'monokai' : 'github'
    }
  });
};

watch(
  () => [props.modelValue, props.content, props.isDarkMode],
  () => {
    renderPreviewContent();
  }
);
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
  color: var(--text-primary);
}

.dark-mode .template-preview-dialog-body {
  background: #141a22;
  border-color: #2a313a;
}

.dark-mode .template-preview-render.vditor-reset {
  color: #e5e7eb;
}
</style>
