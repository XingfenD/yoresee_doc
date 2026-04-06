<template>
  <div class="rich-text-editor">
    <header class="rich-text-toolbar">
      <button type="button" class="toolbar-btn" :class="{ 'is-active': isActive('bold') }" @click="toggleBold">
        B
      </button>
      <button type="button" class="toolbar-btn" :class="{ 'is-active': isActive('italic') }" @click="toggleItalic">
        I
      </button>
      <button type="button" class="toolbar-btn" :class="{ 'is-active': isActive('strike') }" @click="toggleStrike">
        S
      </button>
      <span class="toolbar-divider"></span>
      <button
        type="button"
        class="toolbar-btn"
        :class="{ 'is-active': isHeadingActive(1) }"
        @click="toggleHeading(1)"
      >
        H1
      </button>
      <button
        type="button"
        class="toolbar-btn"
        :class="{ 'is-active': isHeadingActive(2) }"
        @click="toggleHeading(2)"
      >
        H2
      </button>
      <button
        type="button"
        class="toolbar-btn"
        :class="{ 'is-active': isActive('bulletList') }"
        @click="toggleBulletList"
      >
        UL
      </button>
      <button
        type="button"
        class="toolbar-btn"
        :class="{ 'is-active': isActive('orderedList') }"
        @click="toggleOrderedList"
      >
        OL
      </button>
      <button
        type="button"
        class="toolbar-btn"
        :class="{ 'is-active': isActive('blockquote') }"
        @click="toggleBlockquote"
      >
        Quote
      </button>
      <button
        type="button"
        class="toolbar-btn"
        :class="{ 'is-active': isActive('codeBlock') }"
        @click="toggleCodeBlock"
      >
        Code
      </button>
      <button
        v-if="commentEnabled"
        type="button"
        class="toolbar-btn"
        @mousedown.prevent
        @click="addInlineComment"
      >
        Comment
      </button>
      <button
        v-for="item in componentToolbarItems"
        :key="item.key"
        type="button"
        class="toolbar-btn"
        @click="runComponentCommand(item)"
      >
        {{ item.label }}
      </button>
      <span class="toolbar-spacer"></span>
      <button type="button" class="toolbar-btn" :disabled="!canUndo" @click="undo">
        Undo
      </button>
      <button type="button" class="toolbar-btn" :disabled="!canRedo" @click="redo">
        Redo
      </button>
    </header>

    <div v-if="editor" ref="bodyScrollRef" class="rich-text-body">
      <button
        v-show="selectionCommentVisible"
        type="button"
        class="selection-comment-trigger"
        :style="selectionCommentStyle"
        :title="t('document.comments')"
        @mousedown.prevent
        @mouseenter="selectionCommentHovering = true"
        @mouseleave="selectionCommentHovering = false"
        @click="handleSelectionCommentClick"
      >
        <el-icon>
          <ChatDotRound />
        </el-icon>
      </button>
      <RichTextParagraphHandle
        :visible="paragraphHandleVisible"
        :style="paragraphHandleStyle"
        :mode="paragraphHandleMode"
        :block-type="paragraphType"
        :actions="paragraphActions"
        plus-title="新增段落"
        more-title="更多操作"
        :type-label="paragraphTypeLabel"
        @action="runParagraphAction"
        @mouseenter="handleHandleEnter"
        @mouseleave="handleHandleLeave"
      />
      <EditorContent :editor="editor" class="rich-text-content-host" />
    </div>
    <div v-else class="rich-text-placeholder">{{ placeholderText }}</div>

    <footer class="rich-text-hint">
      {{ t('document.richTextShortcutHint') }}
    </footer>
  </div>
</template>

<script setup>
import { computed, nextTick, onBeforeUnmount, onMounted, ref, watch } from 'vue';
import { Editor, EditorContent } from '@tiptap/vue-3';
import StarterKit from '@tiptap/starter-kit';
import Placeholder from '@tiptap/extension-placeholder';
import { ChatDotRound } from '@element-plus/icons-vue';
import TurndownService from 'turndown';
import { useI18n } from 'vue-i18n';
import { marked } from 'marked';
import RichTextParagraphHandle from '@/components/document/rich-text/RichTextParagraphHandle.vue';
import { resolveRichTextComponentSystem } from '@/components/document/rich-text/components/registry';
import { COMMENT_ANCHOR_ATTR, CommentAnchorExtension } from '@/components/document/rich-text/extensions/commentAnchorExtension';
import { useRichTextParagraphHandle } from '@/composables/document/editor/rich-text-editor/useRichTextParagraphHandle';

const props = defineProps({
  modelValue: {
    type: [String, Object],
    default: ''
  },
  valueFormat: {
    type: String,
    default: 'markdown',
    validator: (value) => ['markdown', 'json'].includes(value)
  },
  placeholder: {
    type: String,
    default: ''
  },
  externalExtensions: {
    type: Array,
    default: () => []
  },
  enabledComponents: {
    type: Array,
    default: () => ['mindmap', 'drawio']
  },
  commentEnabled: {
    type: Boolean,
    default: false
  },
  reservedBridgeOptions: {
    type: Object,
    default: () => ({})
  }
});

const emit = defineEmits([
  'update:modelValue',
  'commit',
  'ready',
  'comment-add',
  'comment-remove',
  'comment-changed'
]);
const { t } = useI18n();

const editor = ref(null);
const bodyScrollRef = ref(null);
const applyingModelValue = ref(false);
const lastEmittedValue = ref('');
const selectionCommentVisible = ref(false);
const selectionCommentStyle = ref({ top: '0px', left: '0px' });
const selectionCommentHovering = ref(false);
const selectionUpdateRaf = ref(0);

const isJsonValueMode = computed(() => props.valueFormat === 'json');
const placeholderText = computed(() => props.placeholder || t('document.editorPlaceholder'));
const canUndo = computed(() => editor.value?.can().chain().focus().undo().run() ?? false);
const canRedo = computed(() => editor.value?.can().chain().focus().redo().run() ?? false);
const componentToolbarItems = ref([]);

const safeClone = (value) => {
  try {
    return JSON.parse(JSON.stringify(value));
  } catch (_) {
    return { type: 'doc', content: [] };
  }
};

const normalizeJsonDoc = (value) => {
  if (value && typeof value === 'object' && value.type === 'doc') {
    return safeClone(value);
  }
  return { type: 'doc', content: [] };
};

const serializeJsonDoc = (value) => {
  try {
    return JSON.stringify(normalizeJsonDoc(value));
  } catch (_) {
    return '{"type":"doc","content":[]}';
  }
};

const turndown = new TurndownService({
  headingStyle: 'atx',
  bulletListMarker: '-',
  codeBlockStyle: 'fenced',
  emDelimiter: '*'
});

const decodeMindmapSource = (value) => {
  if (!value) {
    return '';
  }
  try {
    return decodeURIComponent(String(value));
  } catch (_) {
    return String(value);
  }
};

const encodeDrawioSourceToBase64 = (value) => {
  const source = String(value || '');
  if (!source) {
    return '';
  }
  try {
    if (typeof TextEncoder !== 'undefined') {
      const bytes = new TextEncoder().encode(source);
      let binary = '';
      const chunkSize = 0x8000;
      for (let index = 0; index < bytes.length; index += chunkSize) {
        const chunk = bytes.subarray(index, index + chunkSize);
        binary += String.fromCharCode(...chunk);
      }
      return btoa(binary);
    }
    return btoa(unescape(encodeURIComponent(source)));
  } catch (_) {
    return '';
  }
};

const decodeDrawioSourceFromBase64 = (value) => {
  const source = String(value || '').trim();
  if (!source) {
    return '';
  }
  try {
    const binary = atob(source);
    if (typeof TextDecoder !== 'undefined') {
      const bytes = Uint8Array.from(binary, (char) => char.charCodeAt(0));
      return new TextDecoder().decode(bytes);
    }
    return decodeURIComponent(escape(binary));
  } catch (_) {
    return '';
  }
};

const normalizeDrawioSource = (value) => {
  const source = String(value || '').trim();
  if (!source) {
    return '';
  }
  if (source.startsWith('base64:')) {
    const decoded = decodeDrawioSourceFromBase64(source.slice('base64:'.length));
    return decoded || '';
  }
  return source;
};

const resolveMindmapSourceFromNode = (node) => {
  if (!node || typeof node.getAttribute !== 'function') {
    return '';
  }
  const rawSource =
    node.getAttribute('data-source') ||
    node.getAttribute('source') ||
    '';
  if (rawSource) {
    return decodeMindmapSource(rawSource);
  }

  if (typeof node.querySelector === 'function') {
    const textarea = node.querySelector('textarea');
    if (textarea?.value) {
      return String(textarea.value);
    }
  }

  return '';
};

const resolveDrawioSourceFromNode = (node) => {
  if (!node || typeof node.getAttribute !== 'function') {
    return '';
  }
  const source =
    node.getAttribute('data-diagram') ||
    node.getAttribute('diagram') ||
    '';
  if (!source) {
    return '';
  }
  try {
    return decodeURIComponent(String(source));
  } catch (_) {
    return String(source);
  }
};

const extractDrawioBlocksFromHtml = (sourceHtml) => {
  const source = String(sourceHtml || '');
  if (!source) {
    return [];
  }

  const blocks = [];
  const pattern = /<([a-z0-9-]+)\b([^>]*)>/gi;
  let matched = pattern.exec(source);
  while (matched) {
    const tagName = String(matched[1] || '').toLowerCase();
    const attrsSource = String(matched[2] || '');
    const dataTypeMatch = attrsSource.match(/\bdata-type=["']([^"']+)["']/i);
    const dataType = String(dataTypeMatch?.[1] || '').toLowerCase();
    const hasDrawioType = tagName === 'yoresee-drawio' || dataType.includes('drawio');
    if (!hasDrawioType) {
      matched = pattern.exec(source);
      continue;
    }
    const dataDiagramMatch = attrsSource.match(/\b(?:data-diagram|diagram)=["']([^"']+)["']/i);
    if (!dataDiagramMatch?.[1]) {
      matched = pattern.exec(source);
      continue;
    }

    let decodedSource = '';
    try {
      decodedSource = decodeURIComponent(String(dataDiagramMatch[1]));
    } catch (_) {
      decodedSource = String(dataDiagramMatch[1]);
    }
    decodedSource = decodedSource.trim();
    if (!decodedSource) {
      matched = pattern.exec(source);
      continue;
    }

    const encoded = encodeDrawioSourceToBase64(decodedSource);
    if (encoded) {
      blocks.push(`\`\`\`drawio\nbase64:${encoded}\n\`\`\``);
    }

    matched = pattern.exec(source);
  }
  return blocks;
};

turndown.addRule('yoreseeMindmap', {
  filter: (node) => {
    if (!node) {
      return false;
    }
    const nodeName = String(node.nodeName || '').toLowerCase();
    if (nodeName === 'yoresee-mindmap') {
      return true;
    }

    if (typeof node.getAttribute !== 'function') {
      return false;
    }
    const dataType = String(node.getAttribute('data-type') || '').toLowerCase();
    const hasDataSource = Boolean(node.getAttribute('data-source') || node.getAttribute('source'));
    return hasDataSource && dataType.includes('mindmap');
  },
  replacement: (_, node) => {
    const source = resolveMindmapSourceFromNode(node).trim();
    if (!source) {
      return '';
    }
    return `\n\n\`\`\`mindmap\n${source}\n\`\`\`\n\n`;
  }
});

turndown.addRule('yoreseeDrawio', {
  filter: (node) => {
    if (!node) {
      return false;
    }
    const nodeName = String(node.nodeName || '').toLowerCase();
    if (nodeName === 'yoresee-drawio') {
      return true;
    }

    if (typeof node.getAttribute !== 'function') {
      return false;
    }
    const dataType = String(node.getAttribute('data-type') || '').toLowerCase();
    const hasData = Boolean(node.getAttribute('data-diagram') || node.getAttribute('diagram'));
    return hasData && dataType.includes('drawio');
  },
  replacement: (_, node) => {
    const source = resolveDrawioSourceFromNode(node).trim();
    if (!source) {
      return '';
    }
    const encoded = encodeDrawioSourceToBase64(source);
    if (!encoded) {
      return '';
    }
    return `\n\n\`\`\`drawio\nbase64:${encoded}\n\`\`\`\n\n`;
  }
});

turndown.addRule('yoreseeCommentAnchor', {
  filter: (node) => node.nodeName === 'SPAN' && node.getAttribute(COMMENT_ANCHOR_ATTR),
  replacement: (content, node) => {
    const anchorId = String(node.getAttribute(COMMENT_ANCHOR_ATTR) || '').trim();
    if (!anchorId) {
      return content || '';
    }
    return `<span ${COMMENT_ANCHOR_ATTR}="${anchorId}">${content || ''}</span>`;
  }
});

marked.setOptions({
  gfm: true,
  breaks: true
});

const markdownToHtml = (value) => {
  const markdownSource = String(value || '');
  if (!markdownSource.trim()) {
    return '<p></p>';
  }

  const sourceWithComponentBlocks = markdownSource
    .replace(
      /```drawio\s*\n([\s\S]*?)```/gi,
      (_, drawioSource) => {
        const normalized = normalizeDrawioSource(drawioSource);
        if (!normalized) {
          return '';
        }
        const encoded = encodeURIComponent(normalized);
        return `\n<yoresee-drawio data-diagram="${encoded}" data-type="drawio"></yoresee-drawio>\n`;
      }
    )
    .replace(
      /```mindmap\s*\n([\s\S]*?)```/gi,
      (_, mindmapSource) => {
        const encoded = encodeURIComponent(String(mindmapSource || '').trim());
        return `\n<yoresee-mindmap data-source="${encoded}"></yoresee-mindmap>\n`;
      }
    );

  const parsed = marked.parse(sourceWithComponentBlocks, { async: false });
  return typeof parsed === 'string' ? parsed : markdownSource;
};

const htmlToMarkdown = (value) => {
  const source = String(value || '');
  if (!source.trim()) {
    return '';
  }
  const drawioFallbackBlocks = extractDrawioBlocksFromHtml(source);
  let markdown = turndown.turndown(source);
  if (drawioFallbackBlocks.length > 0 && !/```drawio[\s\S]*?```/i.test(markdown)) {
    markdown = `${markdown.trimEnd()}\n\n${drawioFallbackBlocks.join('\n\n')}`.trim();
  }
  return markdown.replace(/\n{3,}/g, '\n\n').trimEnd();
};

const serializeModelValue = (value) => {
  if (isJsonValueMode.value) {
    return serializeJsonDoc(value);
  }
  return String(value || '');
};

const modelValueFromEditor = (instance) => {
  if (isJsonValueMode.value) {
    return safeClone(instance.getJSON());
  }
  return htmlToMarkdown(instance.getHTML());
};

const resolveInitialEditorContent = () => {
  if (isJsonValueMode.value) {
    return normalizeJsonDoc(props.modelValue);
  }
  return markdownToHtml(String(props.modelValue || ''));
};

const applyModelValueToEditor = (modelValue) => {
  if (!editor.value) {
    return;
  }
  const serialized = serializeModelValue(modelValue);
  if (serialized === lastEmittedValue.value) {
    return;
  }
  applyingModelValue.value = true;
  if (isJsonValueMode.value) {
    editor.value.commands.setContent(normalizeJsonDoc(modelValue), false);
  } else {
    editor.value.commands.setContent(markdownToHtml(String(modelValue || '')), false, {
      preserveWhitespace: true
    });
  }
  lastEmittedValue.value = serialized;
  applyingModelValue.value = false;
};

const runCommand = (command) => {
  if (!editor.value) {
    return;
  }
  command(editor.value.chain().focus()).run();
};

const isActive = (name) => editor.value?.isActive(name) ?? false;
const isHeadingActive = (level) => editor.value?.isActive('heading', { level }) ?? false;

const toggleBold = () => runCommand((chain) => chain.toggleBold());
const toggleItalic = () => runCommand((chain) => chain.toggleItalic());
const toggleStrike = () => runCommand((chain) => chain.toggleStrike());
const toggleHeading = (level) => runCommand((chain) => chain.toggleHeading({ level }));
const toggleBulletList = () => runCommand((chain) => chain.toggleBulletList());
const toggleOrderedList = () => runCommand((chain) => chain.toggleOrderedList());
const toggleBlockquote = () => runCommand((chain) => chain.toggleBlockquote());
const toggleCodeBlock = () => runCommand((chain) => chain.toggleCodeBlock());
const undo = () => runCommand((chain) => chain.undo());
const redo = () => runCommand((chain) => chain.redo());

const createCommentAnchorId = () =>
  `rt_${Date.now()}_${Math.random().toString(36).slice(2, 8)}`;

const {
  paragraphHandleVisible,
  paragraphHandleStyle,
  paragraphHandleMode,
  paragraphType,
  paragraphActions,
  runParagraphAction,
  handleHandleEnter,
  handleHandleLeave
} = useRichTextParagraphHandle({
  editorRef: editor,
  scrollContainerRef: bodyScrollRef,
  labels: {
    addAbove: '上方插入空段落',
    addBelow: '下方插入空段落',
    delete: t('document.delete')
  },
  onMutated: () => {
    emit('comment-changed');
  }
});

const paragraphTypeLabelMap = {
  paragraph: '段落',
  heading: '标题',
  list: '列表',
  quote: '引用',
  code: '代码块',
  table: '表格',
  divider: '分割线',
  mindmap: 'Mindmap',
  drawio: 'Draw.io'
};
const paragraphTypeLabel = computed(() => paragraphTypeLabelMap[paragraphType.value] || '段落');

const addInlineComment = () => {
  if (!props.commentEnabled || !editor.value) {
    return;
  }
  const { from, to, empty } = editor.value.state.selection;
  if (empty || from === to) {
    return;
  }
  const selectedText = editor.value.state.doc.textBetween(from, to, ' ');
  const anchorId = createCommentAnchorId();
  const chain = editor.value.chain().focus();
  if (typeof chain.setCommentAnchor === 'function') {
    chain.setCommentAnchor(anchorId).run();
  } else {
    chain.setMark('commentAnchor', { anchorId }).run();
  }
  emit('comment-add', {
    id: anchorId,
    text: selectedText || ''
  });
  emit('comment-changed');
  selectionCommentVisible.value = false;
};

const isSelectionInsideEditor = () => {
  const root = editor.value?.view?.dom;
  const selection = window.getSelection?.();
  if (!root || !selection || selection.rangeCount === 0 || selection.isCollapsed) {
    return false;
  }
  const range = selection.getRangeAt(0);
  const common = range.commonAncestorContainer;
  const commonElement = common?.nodeType === Node.ELEMENT_NODE ? common : common?.parentElement;
  return Boolean(commonElement && root.contains(commonElement));
};

const updateSelectionCommentTrigger = () => {
  if (!props.commentEnabled || !editor.value || !bodyScrollRef.value) {
    selectionCommentVisible.value = false;
    return;
  }
  if (!isSelectionInsideEditor()) {
    selectionCommentVisible.value = false;
    return;
  }
  const instance = editor.value;
  const selection = instance.state.selection;
  const { from, to, empty } = selection;
  if (empty || from === to) {
    selectionCommentVisible.value = false;
    return;
  }
  const selectedText = instance.state.doc.textBetween(from, to, ' ').trim();
  if (!selectedText) {
    selectionCommentVisible.value = false;
    return;
  }

  let fromCoords;
  let toCoords;
  try {
    fromCoords = instance.view.coordsAtPos(from);
    toCoords = instance.view.coordsAtPos(to);
  } catch (_) {
    selectionCommentVisible.value = false;
    return;
  }
  const container = bodyScrollRef.value;
  const containerRect = container.getBoundingClientRect();
  const triggerWidth = 28;
  const triggerHeight = 28;
  const rawLeft = toCoords.right - containerRect.left + container.scrollLeft + 8;
  const rawTop = Math.min(fromCoords.top, toCoords.top) - containerRect.top + container.scrollTop - triggerHeight - 6;
  const minLeft = container.scrollLeft + 8;
  const maxLeft = container.scrollLeft + container.clientWidth - triggerWidth - 8;
  const minTop = container.scrollTop + 8;
  const maxTop = container.scrollTop + container.clientHeight - triggerHeight - 8;

  selectionCommentStyle.value = {
    left: `${Math.min(maxLeft, Math.max(minLeft, rawLeft))}px`,
    top: `${Math.min(maxTop, Math.max(minTop, rawTop))}px`
  };
  selectionCommentVisible.value = true;
};

const handleSelectionCommentClick = () => {
  addInlineComment();
};

const handleSelectionOverlayScrollOrResize = () => {
  if (!selectionCommentVisible.value) {
    return;
  }
  updateSelectionCommentTrigger();
};

const requestSelectionCommentTriggerUpdate = () => {
  if (selectionUpdateRaf.value) {
    cancelAnimationFrame(selectionUpdateRaf.value);
  }
  selectionUpdateRaf.value = requestAnimationFrame(() => {
    updateSelectionCommentTrigger();
    selectionUpdateRaf.value = 0;
  });
};

const runComponentCommand = (item) => {
  if (!editor.value || !item?.command) {
    return;
  }
  const chain = editor.value.chain().focus();
  const commandMethod = chain[item.command];
  if (typeof commandMethod !== 'function') {
    return;
  }
  commandMethod.call(chain, item.commandArgs || {});
  chain.run();
};

const getCommentAnchorElements = (ids = null) => {
  const root = editor.value?.view?.dom;
  if (!root) {
    return [];
  }
  const all = Array.from(root.querySelectorAll(`span[${COMMENT_ANCHOR_ATTR}]`));
  if (!Array.isArray(ids) || ids.length === 0) {
    return all;
  }
  const idSet = new Set(ids.map((id) => String(id || '').trim()).filter(Boolean));
  return all.filter((element) => idSet.has(String(element.getAttribute(COMMENT_ANCHOR_ATTR) || '').trim()));
};

const getCommentIds = () => {
  const elements = getCommentAnchorElements();
  const container = bodyScrollRef.value;
  const containerRect = container?.getBoundingClientRect?.();
  const resultMap = new Map();

  elements.forEach((element) => {
    const anchorId = String(element.getAttribute(COMMENT_ANCHOR_ATTR) || '').trim();
    if (!anchorId) {
      return;
    }
    let top = 0;
    if (container && containerRect) {
      const rect = element.getBoundingClientRect();
      top = rect.top - containerRect.top + container.scrollTop;
    }
    if (!resultMap.has(anchorId) || top < resultMap.get(anchorId)) {
      resultMap.set(anchorId, top);
    }
  });

  return Array.from(resultMap.entries()).map(([id, top]) => ({ id, top }));
};

const hlCommentIds = (ids = []) => {
  getCommentAnchorElements(ids).forEach((element) => {
    element.classList.add('comment-anchor-highlight');
  });
};

const unHlCommentIds = (ids = []) => {
  getCommentAnchorElements(ids).forEach((element) => {
    element.classList.remove('comment-anchor-highlight');
  });
};

const removeCommentIds = (ids = []) => {
  if (!editor.value || !Array.isArray(ids) || ids.length === 0) {
    return;
  }
  const idSet = new Set(ids.map((id) => String(id || '').trim()).filter(Boolean));
  if (idSet.size === 0) {
    return;
  }
  let transaction = editor.value.state.tr;
  editor.value.state.doc.descendants((node, pos) => {
    if (!node.isText || !Array.isArray(node.marks) || node.marks.length === 0) {
      return;
    }
    node.marks.forEach((mark) => {
      if (mark.type.name !== 'commentAnchor') {
        return;
      }
      const anchorId = String(mark.attrs?.anchorId || '').trim();
      if (!idSet.has(anchorId)) {
        return;
      }
      transaction = transaction.removeMark(pos, pos + node.nodeSize, mark);
    });
  });

  if (!transaction.docChanged) {
    return;
  }
  editor.value.view.dispatch(transaction);
  emit('comment-remove', Array.from(idSet));
  emit('comment-changed');
};

const scrollToCommentId = async (id) => {
  const targetId = String(id || '').trim();
  if (!targetId) {
    return;
  }
  await nextTick();
  const container = bodyScrollRef.value;
  if (!container) {
    return;
  }
  const target = getCommentAnchorElements([targetId])[0];
  if (!target) {
    return;
  }
  const containerRect = container.getBoundingClientRect();
  const targetRect = target.getBoundingClientRect();
  const top = Math.max(targetRect.top - containerRect.top + container.scrollTop - 24, 0);
  container.scrollTo({ top, behavior: 'smooth' });
  hlCommentIds([targetId]);
};

const buildReservedBridge = () => ({
  ...props.reservedBridgeOptions,
  comments: null,
  collaboration: null
});

onMounted(() => {
  const initialValue = isJsonValueMode.value ? normalizeJsonDoc(props.modelValue) : String(props.modelValue || '');
  lastEmittedValue.value = serializeModelValue(initialValue);
  const componentSystem = resolveRichTextComponentSystem(props.enabledComponents);
  componentToolbarItems.value = componentSystem.toolbarItems;

    editor.value = new Editor({
      extensions: [
      StarterKit.configure({
        codeBlock: true,
        blockquote: true,
        heading: { levels: [1, 2, 3] }
      }),
      Placeholder.configure({
        placeholder: placeholderText.value
      }),
      ...(props.commentEnabled ? [CommentAnchorExtension] : []),
      ...componentSystem.extensions,
      ...props.externalExtensions
    ],
    content: resolveInitialEditorContent(),
      editorProps: {
        attributes: {
          class: 'yoresee-rich-text-content'
        },
        handleKeyDown: (view, event) => {
          if (event.key !== 'Tab' || event.metaKey || event.ctrlKey || event.altKey) {
            return false;
          }
          const instance = editor.value;
          if (instance?.isActive('listItem')) {
            event.preventDefault();
            if (event.shiftKey) {
              instance.chain().focus().liftListItem('listItem').run();
            } else {
              instance.chain().focus().sinkListItem('listItem').run();
            }
            return true;
          }
          if (event.shiftKey) {
            return false;
          }
          event.preventDefault();
          view.dispatch(view.state.tr.insertText('    '));
          return true;
        }
      },
    onUpdate: ({ editor: instance }) => {
      if (applyingModelValue.value) {
        return;
      }
      const nextValue = modelValueFromEditor(instance);
      const serialized = serializeModelValue(nextValue);
      if (serialized === lastEmittedValue.value) {
        return;
      }
      lastEmittedValue.value = serialized;
      emit('update:modelValue', nextValue);
      updateSelectionCommentTrigger();
    },
    onSelectionUpdate: () => {
      updateSelectionCommentTrigger();
    },
    onBlur: () => {
      if (!selectionCommentHovering.value) {
        selectionCommentVisible.value = false;
      }
      emit('commit');
    }
  });

  emit('ready', {
    editor: editor.value,
    reservedBridge: buildReservedBridge()
  });
});

watch(
  () => props.modelValue,
  (nextValue) => {
    applyModelValueToEditor(nextValue);
  }
);

defineExpose({
  reRender: () => {
    const instance = editor.value;
    if (!instance?.view) {
      return;
    }
    instance.view.updateState(instance.state);
  },
  getEditor: () => editor.value,
  getCommentIds,
  hlCommentIds,
  unHlCommentIds,
  removeCommentIds,
  scrollToCommentId,
  broadcastCommentChange: () => {
    emit('comment-changed');
  }
});

onBeforeUnmount(() => {
  if (selectionUpdateRaf.value) {
    cancelAnimationFrame(selectionUpdateRaf.value);
    selectionUpdateRaf.value = 0;
  }
  bodyScrollRef.value?.removeEventListener('mouseup', requestSelectionCommentTriggerUpdate);
  bodyScrollRef.value?.removeEventListener('keyup', requestSelectionCommentTriggerUpdate);
  bodyScrollRef.value?.removeEventListener('scroll', handleSelectionOverlayScrollOrResize);
  document.removeEventListener('selectionchange', requestSelectionCommentTriggerUpdate);
  window.removeEventListener('resize', handleSelectionOverlayScrollOrResize);
  editor.value?.destroy();
  editor.value = null;
});

onMounted(() => {
  nextTick(() => {
    bodyScrollRef.value?.addEventListener('scroll', handleSelectionOverlayScrollOrResize, { passive: true });
    bodyScrollRef.value?.addEventListener('mouseup', requestSelectionCommentTriggerUpdate);
    bodyScrollRef.value?.addEventListener('keyup', requestSelectionCommentTriggerUpdate);
  });
  document.addEventListener('selectionchange', requestSelectionCommentTriggerUpdate);
  window.addEventListener('resize', handleSelectionOverlayScrollOrResize);
});
</script>

<style scoped>
.rich-text-editor {
  display: flex;
  flex-direction: column;
  width: 100%;
  height: 100%;
  min-width: 0;
  min-height: 0;
  background: var(--bg-white);
}

.rich-text-toolbar {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 10px;
  border-bottom: 1px solid var(--border-color);
  flex-wrap: wrap;
}

.toolbar-btn {
  border: 1px solid var(--border-color);
  background: var(--bg-white);
  color: var(--text-medium);
  border-radius: 6px;
  height: 28px;
  min-width: 32px;
  padding: 0 8px;
  font-size: 12px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.toolbar-btn:hover {
  border-color: color-mix(in srgb, var(--primary-color) 52%, var(--border-color) 48%);
  color: var(--primary-color);
}

.toolbar-btn.is-active {
  background: color-mix(in srgb, var(--primary-color) 14%, transparent);
  border-color: color-mix(in srgb, var(--primary-color) 52%, var(--border-color) 48%);
  color: var(--primary-color);
}

.toolbar-btn:disabled {
  opacity: 0.45;
  cursor: not-allowed;
}

.toolbar-divider {
  width: 1px;
  height: 20px;
  background: var(--border-color);
}

.toolbar-spacer {
  flex: 1;
}

.rich-text-body {
  flex: 1;
  min-height: 0;
  overflow: auto;
  position: relative;
}

.selection-comment-trigger {
  position: absolute;
  z-index: 35;
  width: 28px;
  height: 28px;
  border-radius: 8px;
  border: 1px solid var(--border-color);
  background: var(--bg-white);
  color: var(--text-medium);
  box-shadow: var(--shadow-sm);
  display: inline-flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.18s ease;
}

.selection-comment-trigger:hover {
  color: var(--primary-color);
  border-color: color-mix(in srgb, var(--primary-color) 52%, var(--border-color) 48%);
  background: color-mix(in srgb, var(--primary-color) 10%, var(--bg-white) 90%);
}

.rich-text-placeholder {
  flex: 1;
  min-height: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--text-light);
}

.rich-text-hint {
  border-top: 1px solid var(--border-color);
  padding: 6px 10px;
  color: var(--text-light);
  font-size: 12px;
}

.rich-text-body :deep(.yoresee-rich-text-content) {
  min-height: 100%;
  outline: none;
  padding: 14px 18px 14px 56px;
  color: var(--text-primary);
  line-height: 1.7;
  box-sizing: border-box;
}

.rich-text-body :deep(.yoresee-rich-text-content h1),
.rich-text-body :deep(.yoresee-rich-text-content h2),
.rich-text-body :deep(.yoresee-rich-text-content h3) {
  margin: 0 0 12px;
  line-height: 1.3;
}

.rich-text-body :deep(.yoresee-rich-text-content p) {
  margin: 0 0 10px;
}

.rich-text-body :deep(.yoresee-rich-text-content ul),
.rich-text-body :deep(.yoresee-rich-text-content ol) {
  margin: 0 0 10px;
  padding-left: 22px;
}

.rich-text-body :deep(.yoresee-rich-text-content blockquote) {
  margin: 0 0 10px;
  padding-left: 12px;
  border-left: 3px solid color-mix(in srgb, var(--primary-color) 38%, var(--border-color) 62%);
  color: var(--text-medium);
}

.rich-text-body :deep(.yoresee-rich-text-content .yoresee-comment-anchor) {
  background: color-mix(in srgb, #f59e0b 20%, transparent);
  border-radius: 3px;
  transition: background-color 0.18s ease;
}

.rich-text-body :deep(.yoresee-rich-text-content .yoresee-comment-anchor.comment-anchor-highlight) {
  background: color-mix(in srgb, #f59e0b 38%, transparent);
}

.rich-text-body :deep(.yoresee-rich-text-content pre) {
  background: color-mix(in srgb, var(--bg-light) 75%, transparent);
  border: 1px solid var(--border-color);
  border-radius: 8px;
  padding: 10px 12px;
  overflow: auto;
}

:global(.dark-mode) .rich-text-editor {
  background: var(--bg-white);
}

:global(.dark-mode) .toolbar-btn {
  background: color-mix(in srgb, var(--bg-white) 96%, #000 4%);
  color: var(--text-medium);
}

:global(.dark-mode) .rich-text-body :deep(.yoresee-rich-text-content) {
  color: var(--text-dark);
}

:global(.dark-mode) .rich-text-body :deep(.yoresee-rich-text-content pre) {
  background: color-mix(in srgb, var(--bg-light) 20%, #05070a 80%);
}
</style>
