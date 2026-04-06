export function useRichTextToolbarActions(options = {}) {
  const {
    editorRef,
    codeLanguageDraftRef,
    resolveCodeBlockLanguage,
    hideCodeLanguageFloating,
    updateCodeLanguageFloating
  } = options;

  const runCommand = (command) => {
    if (!editorRef.value) {
      return;
    }
    command(editorRef.value.chain().focus()).run();
  };

  const isActive = (name) => editorRef.value?.isActive(name) ?? false;
  const isHeadingActive = (level) => editorRef.value?.isActive('heading', { level }) ?? false;

  const toggleBold = () => runCommand((chain) => chain.toggleBold());
  const toggleUnderline = () => {
    if (!editorRef.value) {
      return;
    }
    const chain = editorRef.value.chain().focus();
    if (typeof chain.toggleUnderline !== 'function') {
      return;
    }
    chain.toggleUnderline().run();
  };
  const toggleItalic = () => runCommand((chain) => chain.toggleItalic());
  const toggleStrike = () => runCommand((chain) => chain.toggleStrike());
  const toggleInlineCode = () => runCommand((chain) => chain.toggleCode());
  const toggleHeading = (level) => runCommand((chain) => chain.toggleHeading({ level }));
  const toggleBulletList = () => runCommand((chain) => chain.toggleBulletList());
  const toggleOrderedList = () => runCommand((chain) => chain.toggleOrderedList());
  const toggleBlockquote = () => runCommand((chain) => chain.toggleBlockquote());
  const resolveLinkChain = (pos) => {
    const chain = editorRef.value.chain();
    if (typeof pos === 'number' && Number.isFinite(pos)) {
      return chain.focus(pos).extendMarkRange('link');
    }
    return chain.focus().extendMarkRange('link');
  };

  const setLink = (href = '', options = {}) => {
    if (!editorRef.value) {
      return;
    }
    const url = String(href || '').trim();
    if (!url) {
      resolveLinkChain(options.pos).unsetLink().run();
      return;
    }
    resolveLinkChain(options.pos).setLink({ href: url }).run();
  };
  const unsetLink = (options = {}) => {
    if (!editorRef.value) {
      return;
    }
    resolveLinkChain(options.pos).unsetLink().run();
  };
  const setTextColor = (color = '') => {
    if (!editorRef.value) {
      return;
    }
    const value = String(color || '').trim();
    if (!value) {
      editorRef.value.chain().focus().unsetColor().run();
      return;
    }
    editorRef.value.chain().focus().setColor(value).run();
  };
  const setBackgroundColor = (color = '') => {
    if (!editorRef.value) {
      return;
    }
    const value = String(color || '').trim();
    if (!value) {
      editorRef.value.chain().focus().unsetHighlight().run();
      return;
    }
    editorRef.value.chain().focus().setHighlight({ color: value }).run();
  };

  const toggleCodeBlock = () => {
    if (!editorRef.value) {
      return;
    }
    const chain = editorRef.value.chain().focus();
    if (editorRef.value.isActive('codeBlock')) {
      chain.toggleCodeBlock().run();
      hideCodeLanguageFloating?.();
      return;
    }
    const nextLanguage = resolveCodeBlockLanguage?.(codeLanguageDraftRef.value) || 'plaintext';
    codeLanguageDraftRef.value = nextLanguage;
    chain.setCodeBlock({ language: nextLanguage }).run();
    updateCodeLanguageFloating?.();
  };

  const undo = () => runCommand((chain) => chain.undo());
  const redo = () => runCommand((chain) => chain.redo());

  const runComponentCommand = (item) => {
    if (!editorRef.value || !item?.command) {
      return;
    }
    const chain = editorRef.value.chain().focus();
    const commandMethod = chain[item.command];
    if (typeof commandMethod !== 'function') {
      return;
    }
    commandMethod.call(chain, item.commandArgs || {});
    chain.run();
  };

  return {
    isActive,
    isHeadingActive,
    toggleBold,
    toggleUnderline,
    toggleItalic,
    toggleStrike,
    toggleInlineCode,
    toggleHeading,
    toggleBulletList,
    toggleOrderedList,
    toggleBlockquote,
    toggleCodeBlock,
    setLink,
    unsetLink,
    setTextColor,
    setBackgroundColor,
    undo,
    redo,
    runComponentCommand
  };
}
