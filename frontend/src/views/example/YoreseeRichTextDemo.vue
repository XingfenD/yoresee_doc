<template>
  <div class="demo-page">
    <header class="demo-header">
      <h1>YoreseeRichText Demo</h1>
      <p>纯前端样例：左侧编辑，右侧实时看 raw/diff。点击工具栏 `Mindmap` 或 `Draw.io` 可插入组件。</p>
    </header>

    <div class="demo-grid">
      <section class="editor-panel">
        <YoreseeRichTextEditor
          v-model="content"
          value-format="json"
          :enabled-components="['mindmap', 'drawio']"
          :comment-enabled="true"
        />
      </section>

      <section class="raw-panel">
        <div class="panel-title">Raw Data (JSON)</div>
        <pre class="raw-content">{{ rawJson }}</pre>
      </section>

      <section class="diff-text-panel">
        <div class="panel-title">Diff Text (Noise Stripped)</div>
        <pre class="raw-content">{{ diffText }}</pre>
      </section>

    </div>
  </div>
</template>

<script setup>
import { computed, ref } from 'vue';
import YoreseeRichTextEditor from '@/components/document/YoreseeRichTextEditor.vue';
import { DEFAULT_MINDMAP_SOURCE } from '@/components/document/rich-text/components/mindmap/mindmapExtension';
import { DEFAULT_DRAWIO_XML } from '@/components/document/rich-text/components/drawio/drawioExtension';
import { serializeRichTextJsonToDiffText } from '@/components/document/rich-text/diff/serializeRichTextJsonToDiffText';
import { resolveRichTextPreviewDiffAdapterRegistry } from '@/components/document/rich-text/components/registry';

const content = ref({
  type: 'doc',
  content: [
    {
      type: 'heading',
      attrs: { level: 1 },
      content: [{ type: 'text', text: 'YoreseeRichText' }]
    },
    {
      type: 'paragraph',
      content: [{ type: 'text', text: '这是一个纯前端示例页面。' }]
    },
    {
      type: 'bulletList',
      content: [
        { type: 'listItem', content: [{ type: 'paragraph', content: [{ type: 'text', text: '支持标题、列表、引用、代码块' }] }] },
        { type: 'listItem', content: [{ type: 'paragraph', content: [{ type: 'text', text: '实时看到 raw data（JSON）' }] }] },
        { type: 'listItem', content: [{ type: 'paragraph', content: [{ type: 'text', text: '组件节点由 type + attrs 表达' }] }] }
      ]
    },
    {
      type: 'blockquote',
      content: [{ type: 'paragraph', content: [{ type: 'text', text: '这里是引用' }] }]
    },
    {
      type: 'codeBlock',
      attrs: { language: 'js' },
      content: [{ type: 'text', text: "console.log('hello yoresee rich text');" }]
    },
    {
      type: 'mindmapBlock',
      attrs: { source: DEFAULT_MINDMAP_SOURCE }
    },
    {
      type: 'drawioBlock',
      attrs: { diagram: DEFAULT_DRAWIO_XML }
    }
  ]
});

const rawJson = computed(() => JSON.stringify(content.value, null, 2));
const previewDiffAdapters = resolveRichTextPreviewDiffAdapterRegistry(['mindmap', 'drawio']);

const diffText = computed(() => serializeRichTextJsonToDiffText(content.value, previewDiffAdapters));
</script>

<style scoped>
.demo-page {
  height: 100vh;
  min-height: 0;
  display: flex;
  flex-direction: column;
  background: var(--bg-light);
  color: var(--text-primary);
  padding: 16px;
  box-sizing: border-box;
  gap: 12px;
}

.demo-header h1 {
  margin: 0;
  font-size: 22px;
  line-height: 1.2;
}

.demo-header p {
  margin: 6px 0 0;
  color: var(--text-medium);
  font-size: 13px;
}

.demo-grid {
  flex: 1;
  min-height: 0;
  display: grid;
  grid-template-columns: minmax(440px, 1.2fr) minmax(320px, 1fr);
  grid-template-rows: minmax(200px, 1fr) minmax(200px, 1fr);
  grid-template-areas:
    'editor raw'
    'editor diff';
  gap: 12px;
}

.editor-panel,
.raw-panel,
.diff-text-panel {
  min-height: 0;
  min-width: 0;
  border: 1px solid var(--border-color);
  border-radius: 10px;
  background: var(--bg-white);
  overflow: hidden;
}

.editor-panel {
  grid-area: editor;
}

.raw-panel {
  grid-area: raw;
}

.diff-text-panel {
  grid-area: diff;
}

.panel-title {
  border-bottom: 1px solid var(--border-color);
  padding: 10px 12px;
  font-size: 13px;
  font-weight: 600;
  color: var(--text-medium);
}

.raw-content {
  margin: 0;
  padding: 12px;
  height: calc(100% - 42px);
  overflow: auto;
  font-size: 12px;
  line-height: 1.6;
  color: var(--text-primary);
  background: transparent;
  white-space: pre-wrap;
  word-break: break-word;
  box-sizing: border-box;
}

@media (max-width: 1180px) {
  .demo-page {
    height: auto;
    min-height: 100vh;
  }

  .demo-grid {
    grid-template-columns: 1fr;
    grid-template-rows: 420px 260px 260px;
    grid-template-areas:
      'editor'
      'raw'
      'diff';
  }

  .editor-panel {
    grid-area: editor;
  }
}
</style>
