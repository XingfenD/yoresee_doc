<template>
  <div class="diff-viewer">
    <div class="diff-head">
      <div class="diff-head-cell">{{ leftTitle }}</div>
      <div class="diff-head-cell">{{ rightTitle }}</div>
    </div>
    <div ref="diffRoot" class="diff-body"></div>
  </div>
</template>

<script setup>
import { nextTick, onBeforeUnmount, onMounted, ref, watch } from 'vue';
import { createTwoFilesPatch } from 'diff';
import { html as renderDiffHtml } from 'diff2html';
import 'diff2html/bundles/css/diff2html.min.css';

const props = defineProps({
  leftText: {
    type: String,
    default: ''
  },
  rightText: {
    type: String,
    default: ''
  },
  leftTitle: {
    type: String,
    default: ''
  },
  rightTitle: {
    type: String,
    default: ''
  }
});

const diffRoot = ref(null);

const renderDiff = async () => {
  await nextTick();
  if (!diffRoot.value) return;
  const patch = createTwoFilesPatch(
    'document',
    'document',
    `${props.leftText || ''}`,
    `${props.rightText || ''}`,
    '',
    '',
    { context: 3 }
  );
  const html = renderDiffHtml(patch, {
    outputFormat: 'side-by-side',
    drawFileList: false,
    matching: 'lines',
    diffStyle: 'word',
    renderNothingWhenEmpty: false
  });
  if (diffRoot.value) {
    diffRoot.value.innerHTML = html;
  }
};

onMounted(async () => {
  await renderDiff();
});

watch(
  () => [props.leftText, props.rightText, props.leftTitle, props.rightTitle],
  async () => {
    await renderDiff();
  }
);

onBeforeUnmount(() => {
  if (diffRoot.value) {
    diffRoot.value.innerHTML = '';
  }
});
</script>

<style scoped>
.diff-viewer {
  border: 1px solid var(--border-color);
  border-radius: var(--border-radius-sm);
  overflow: hidden;
  background: var(--bg-white);
}

.diff-head {
  display: grid;
  grid-template-columns: 1fr 1fr;
  border-bottom: 1px solid var(--border-color);
  background: #f7f9fc;
}

.diff-head-cell {
  padding: 8px 12px;
  font-size: 13px;
  color: var(--text-medium);
  border-right: 1px solid var(--border-color);
}

.diff-head-cell:last-child {
  border-right: 0;
}

.diff-body {
  padding: 12px;
  overflow: visible;
}

:deep(.d2h-wrapper) {
  margin: 0;
  border: 1px solid var(--border-color);
  border-radius: var(--border-radius-sm);
  overflow: hidden;
  background: var(--bg-white);
  font-size: 14px;
}

:deep(.d2h-file-wrapper) {
  margin-bottom: 0;
}

:deep(.d2h-file-side-diff) {
  max-height: 640px;
  overflow-x: auto;
  overflow-y: auto;
}

:global(.dark-mode) .diff-head {
  background: #1b2330;
}

:global(.dark-mode) :deep(.d2h-wrapper),
:global(.dark-mode) :deep(.d2h-file-wrapper) {
  background: #111827;
  color: #e5e7eb;
}

:global(.dark-mode) :deep(.d2h-file-header),
:global(.dark-mode) :deep(.d2h-code-linenumber),
:global(.dark-mode) :deep(.d2h-code-side-linenumber),
:global(.dark-mode) :deep(.d2h-code-side-emptyplaceholder),
:global(.dark-mode) :deep(.d2h-file-side-diff) {
  background: #1f2937;
  color: #d1d5db;
  border-color: #374151;
}

@media (max-width: 1200px) {
  .diff-body {
    padding: 8px;
  }

  :deep(.d2h-file-side-diff) {
    max-height: 520px;
  }
}
</style>
