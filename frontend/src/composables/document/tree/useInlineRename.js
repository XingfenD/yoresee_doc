import { nextTick, ref } from 'vue';

const resolveBool = (valueOrGetterOrRef, fallback = false) => {
  if (typeof valueOrGetterOrRef === 'function') {
    return valueOrGetterOrRef();
  }
  if (valueOrGetterOrRef && typeof valueOrGetterOrRef === 'object' && 'value' in valueOrGetterOrRef) {
    return valueOrGetterOrRef.value;
  }
  if (valueOrGetterOrRef === undefined) {
    return fallback;
  }
  return Boolean(valueOrGetterOrRef);
};

export function useInlineRename(options = {}) {
  const { revertRename = true, onConfirm = null } = options;

  const renamingInputRef = ref(null);

  const setRenamingInputRef = (el, data) => {
    if (el && data?.isRenaming) {
      renamingInputRef.value = el;
    }
  };

  const focusRenameInput = () => {
    const inputEl = renamingInputRef.value?.input || renamingInputRef.value;
    if (!inputEl || typeof inputEl.focus !== 'function') {
      return;
    }
    inputEl.focus();
    if (typeof inputEl.select === 'function') {
      inputEl.select();
    }
  };

  const startInlineRename = (node) => {
    if (!node || node.isRenaming) {
      return;
    }

    node.isRenaming = true;
    node.originalLabel = node.label;
    node.renameValue = node.label;

    nextTick(() => {
      focusRenameInput();
    });
  };

  const clearRenameState = (node) => {
    node.renameValue = '';
    node.isRenaming = false;
  };

  const cancelInlineRename = (node) => {
    if (!node?.isRenaming) {
      return;
    }
    node.label = node.originalLabel || node.label;
    clearRenameState(node);
  };

  const confirmInlineRename = (node) => {
    if (!node?.isRenaming) {
      return;
    }

    const nextName = node.renameValue?.trim();
    if (!nextName) {
      cancelInlineRename(node);
      return;
    }

    clearRenameState(node);

    if (typeof onConfirm === 'function') {
      onConfirm({
        node,
        title: nextName
      });
    }

    if (resolveBool(revertRename, true)) {
      node.label = node.originalLabel || node.label;
    } else {
      node.label = nextName;
    }
  };

  return {
    renamingInputRef,
    setRenamingInputRef,
    startInlineRename,
    cancelInlineRename,
    confirmInlineRename
  };
}
