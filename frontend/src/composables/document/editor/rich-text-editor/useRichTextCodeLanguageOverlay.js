import { ref } from 'vue';
import { common, createLowlight } from 'lowlight';

const BUILTIN_CODE_LANGUAGES = Object.freeze([
  { value: 'plaintext', label: 'Plain Text' },
  { value: 'javascript', label: 'JavaScript' },
  { value: 'typescript', label: 'TypeScript' },
  { value: 'json', label: 'JSON' },
  { value: 'html', label: 'HTML' },
  { value: 'css', label: 'CSS' },
  { value: 'markdown', label: 'Markdown' },
  { value: 'bash', label: 'Bash' },
  { value: 'sql', label: 'SQL' },
  { value: 'yaml', label: 'YAML' },
  { value: 'go', label: 'Go' },
  { value: 'python', label: 'Python' },
  { value: 'java', label: 'Java' },
  { value: 'rust', label: 'Rust' }
]);

const CODE_LANGUAGE_ALIAS_MAP = Object.freeze({
  text: 'plaintext',
  plain: 'plaintext',
  txt: 'plaintext',
  js: 'javascript',
  ts: 'typescript',
  yml: 'yaml',
  sh: 'bash',
  shell: 'bash',
  py: 'python',
  md: 'markdown'
});

const builtinCodeLanguageSet = new Set(BUILTIN_CODE_LANGUAGES.map((item) => item.value));

export function useRichTextCodeLanguageOverlay(options = {}) {
  const { editorRef, scrollContainerRef } = options;
  const lowlight = createLowlight(common);

  const codeLanguageDraft = ref('plaintext');
  const codeLanguageFloatingVisible = ref(false);
  const codeLanguageFloatingStyle = ref({ top: '0px', left: '0px' });
  const codeLanguageInteracting = ref(false);
  const codeBlockAnchorPos = ref(null);

  const sanitizeCodeLanguage = (value) => {
    const raw = String(value || '').trim().toLowerCase();
    if (!raw) {
      return '';
    }
    const aliased = CODE_LANGUAGE_ALIAS_MAP[raw] || raw;
    return aliased.replace(/[^a-z0-9_+-]/g, '');
  };

  const isSupportedHighlightLanguage = (language) => {
    if (!language) {
      return false;
    }
    if (builtinCodeLanguageSet.has(language)) {
      return true;
    }
    if (typeof lowlight.registered === 'function') {
      return Boolean(lowlight.registered(language));
    }
    return false;
  };

  const resolveCodeBlockLanguage = (value = codeLanguageDraft.value) => {
    const normalized = sanitizeCodeLanguage(value);
    if (!normalized) {
      return 'plaintext';
    }
    return isSupportedHighlightLanguage(normalized) ? normalized : 'plaintext';
  };

  const resolveSelectionCodeBlockMeta = () => {
    const instance = editorRef.value;
    if (!instance) {
      return null;
    }
    const { $from } = instance.state.selection;
    for (let depth = $from.depth; depth >= 0; depth -= 1) {
      const node = $from.node(depth);
      if (node?.type?.name !== 'codeBlock') {
        continue;
      }
      const pos = depth > 0 ? $from.before(depth) : 0;
      const rawDom = instance.view.nodeDOM(pos);
      let element = rawDom instanceof HTMLElement ? rawDom : null;
      if (!element && rawDom && typeof Node !== 'undefined' && rawDom.nodeType === Node.TEXT_NODE) {
        element = rawDom.parentElement;
      }
      if (!element) {
        return {
          pos,
          node
        };
      }
      if (element.tagName.toLowerCase() !== 'pre') {
        element = element.querySelector?.('pre') || element;
      }
      return {
        pos,
        node,
        element
      };
    }
    return null;
  };

  const resolveAnchorCodeBlockMeta = () => {
    const instance = editorRef.value;
    const pos = Number(codeBlockAnchorPos.value);
    if (!instance || !Number.isFinite(pos)) {
      return null;
    }
    const node = instance.state.doc.nodeAt(pos);
    if (!node || node.type?.name !== 'codeBlock') {
      return null;
    }
    const rawDom = instance.view.nodeDOM(pos);
    let element = rawDom instanceof HTMLElement ? rawDom : null;
    if (!element && rawDom && typeof Node !== 'undefined' && rawDom.nodeType === Node.TEXT_NODE) {
      element = rawDom.parentElement;
    }
    if (element && element.tagName.toLowerCase() !== 'pre') {
      element = element.querySelector?.('pre') || element;
    }
    return {
      pos,
      node,
      element: element || null
    };
  };

  const syncCodeLanguageDraft = () => {
    const selected = resolveSelectionCodeBlockMeta();
    if (!selected) {
      return;
    }
    codeBlockAnchorPos.value = selected.pos;
    const current = sanitizeCodeLanguage(selected.node?.attrs?.language);
    codeLanguageDraft.value = isSupportedHighlightLanguage(current) ? current : 'plaintext';
  };

  const resolveActiveCodeBlockDom = () => {
    const selected = resolveSelectionCodeBlockMeta();
    if (selected) {
      codeBlockAnchorPos.value = selected.pos;
      return selected.element || null;
    }
    const anchored = resolveAnchorCodeBlockMeta();
    return anchored?.element || null;
  };

  const updateCodeLanguageFloating = () => {
    const container = scrollContainerRef.value;
    const blockElement = resolveActiveCodeBlockDom();
    if (!container || !blockElement) {
      if (!codeLanguageInteracting.value) {
        codeLanguageFloatingVisible.value = false;
      }
      return;
    }
    const containerRect = container.getBoundingClientRect();
    const blockRect = blockElement.getBoundingClientRect();
    const floatingWidth = 132;
    const floatingHeight = 30;
    const rawLeft = blockRect.right - containerRect.left + container.scrollLeft - floatingWidth - 8;
    const rawTop = blockRect.top - containerRect.top + container.scrollTop + 8;
    const minLeft = container.scrollLeft + 8;
    const maxLeft = container.scrollLeft + container.clientWidth - floatingWidth - 8;
    const minTop = container.scrollTop + 8;
    const maxTop = container.scrollTop + container.clientHeight - floatingHeight - 8;
    codeLanguageFloatingStyle.value = {
      left: `${Math.min(maxLeft, Math.max(minLeft, rawLeft))}px`,
      top: `${Math.min(maxTop, Math.max(minTop, rawTop))}px`
    };
    codeLanguageFloatingVisible.value = true;
  };

  const queryCodeLanguageSuggestions = (queryString, cb) => {
    const keyword = String(queryString || '').trim().toLowerCase();
    const next = !keyword
      ? BUILTIN_CODE_LANGUAGES
      : BUILTIN_CODE_LANGUAGES.filter(
          (item) => item.value.includes(keyword) || item.label.toLowerCase().includes(keyword)
        );
    cb(next);
  };

  const applyCodeLanguage = (value) => {
    const instance = editorRef.value;
    if (!instance) {
      return;
    }
    const nextLanguage = resolveCodeBlockLanguage(value);
    if (!nextLanguage) {
      return;
    }
    const selected = resolveSelectionCodeBlockMeta() || resolveAnchorCodeBlockMeta();
    if (!selected || !Number.isFinite(selected.pos)) {
      return;
    }

    const currentNode = instance.state.doc.nodeAt(selected.pos);
    if (!currentNode || currentNode.type?.name !== 'codeBlock') {
      return;
    }
    const nextAttrs = {
      ...currentNode.attrs,
      language: nextLanguage
    };
    instance.view.dispatch(
      instance.state.tr.setNodeMarkup(selected.pos, currentNode.type, nextAttrs)
    );
    codeBlockAnchorPos.value = selected.pos;
    codeLanguageDraft.value = nextLanguage;
    updateCodeLanguageFloating();
  };

  const applyCodeLanguageFromDraft = () => {
    applyCodeLanguage(codeLanguageDraft.value);
  };

  const handleCodeLanguageSelect = (item) => {
    const value = item?.value || item?.label || codeLanguageDraft.value;
    applyCodeLanguage(value);
  };

  const handleCodeLanguageInputBlur = () => {
    setTimeout(() => {
      if (!codeLanguageInteracting.value) {
        updateCodeLanguageFloating();
      }
    }, 0);
  };

  const hideCodeLanguageFloating = () => {
    codeBlockAnchorPos.value = null;
    codeLanguageFloatingVisible.value = false;
  };

  const handleEditorBlur = () => {
    if (!codeLanguageInteracting.value) {
      hideCodeLanguageFloating();
    }
  };

  return {
    lowlight,
    codeLanguageDraft,
    codeLanguageFloatingVisible,
    codeLanguageFloatingStyle,
    codeLanguageInteracting,
    queryCodeLanguageSuggestions,
    applyCodeLanguageFromDraft,
    handleCodeLanguageSelect,
    handleCodeLanguageInputBlur,
    syncCodeLanguageDraft,
    updateCodeLanguageFloating,
    hideCodeLanguageFloating,
    handleEditorBlur,
    resolveCodeBlockLanguage
  };
}
