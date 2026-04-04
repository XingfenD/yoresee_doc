<template>
  <div class="document-preview-viewer">
    <TablePreviewViewer
      v-if="previewKind === 'table'"
      :rows="tableRows"
    />
    <SlidePreviewViewer
      v-else-if="previewKind === 'slide'"
      :slides="slideSections"
    />
    <div v-else class="markdown-preview-shell">
      <el-empty v-if="!previewMarkdown" />
      <div v-else ref="markdownRef" class="markdown-preview-body"></div>
    </div>
  </div>
</template>

<script setup>
import { computed, nextTick, ref, watch } from 'vue';
import Vditor from 'vditor';
import 'vditor/dist/index.css';
import {
  resolveDocumentPreviewContent,
  resolveDocumentPreviewKind,
  resolveSlidePreviewSections,
  resolveTablePreviewRows
} from '@/composables/document/render/useDocumentRenderBridge';
import TablePreviewViewer from '@/components/document/render/TablePreviewViewer.vue';
import SlidePreviewViewer from '@/components/document/render/SlidePreviewViewer.vue';

const props = defineProps({
  content: {
    type: String,
    default: ''
  },
  documentType: {
    type: [String, Number],
    default: '1'
  },
  isTemplate: {
    type: Boolean,
    default: false
  },
  isDarkMode: {
    type: Boolean,
    default: false
  }
});

const markdownRef = ref(null);

const previewKind = computed(() => resolveDocumentPreviewKind(props.documentType));
const previewMarkdown = computed(() => resolveDocumentPreviewContent({
  content: props.content,
  documentType: props.documentType,
  isTemplate: props.isTemplate
}));
const tableRows = computed(() => resolveTablePreviewRows({
  content: props.content,
  isTemplate: props.isTemplate
}));
const slideSections = computed(() => resolveSlidePreviewSections({
  content: props.content,
  isTemplate: props.isTemplate
}));

const renderMarkdown = async () => {
  if (previewKind.value !== 'markdown') {
    return;
  }
  await nextTick();
  if (!markdownRef.value) {
    return;
  }
  if (!previewMarkdown.value) {
    markdownRef.value.innerHTML = '';
    return;
  }
  await Vditor.preview(markdownRef.value, previewMarkdown.value, {
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
  () => [previewKind.value, previewMarkdown.value, props.isDarkMode],
  () => {
    renderMarkdown();
  },
  { immediate: true }
);
</script>

<style scoped>
.document-preview-viewer {
  width: 100%;
  height: 100%;
  min-height: 0;
}

.markdown-preview-shell {
  width: 100%;
  height: 100%;
  min-height: 0;
  border: 1px solid var(--border-color);
  border-radius: var(--border-radius-sm);
  background: var(--bg-white);
  overflow: auto;
}

.markdown-preview-body {
  padding: 12px;
  color: var(--text-primary);
  min-height: 100%;
  box-sizing: border-box;
}

.dark-mode .markdown-preview-body.vditor-reset {
  color: #e5e7eb;
}
</style>

