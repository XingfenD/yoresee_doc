import { ref } from 'vue';
import * as Y from 'yjs';
import { WebsocketProvider } from 'y-websocket';

export function useYjsCollaboration({
  props,
  emit,
  editorRef,
  vditorRef,
  isVditorReady,
  suppressInput,
  getEditableElement,
  getValueSafely,
  setValueSafely
}) {
  const ydocRef = ref(null);
  const providerRef = ref(null);
  const ytextRef = ref(null);
  const ycommentMetaRef = ref(null);
  const commentMetaObserverRef = ref(null);
  const activeRoomRef = ref('');
  const collabSyncedRef = ref(false);
  const pendingSeedRef = ref('');

  const getCaretOffset = (container) => {
    const selection = window.getSelection();
    if (!selection || selection.rangeCount === 0) {
      return null;
    }
    const range = selection.getRangeAt(0);
    if (!container.contains(range.startContainer)) {
      return null;
    }
    const preRange = range.cloneRange();
    preRange.selectNodeContents(container);
    preRange.setEnd(range.startContainer, range.startOffset);
    return preRange.toString().length;
  };

  const restoreCaretOffset = (container, offset) => {
    if (offset === null || offset === undefined) {
      return;
    }
    const selection = window.getSelection();
    if (!selection) {
      return;
    }
    const walker = document.createTreeWalker(container, NodeFilter.SHOW_TEXT);
    let currentOffset = 0;
    let node = walker.nextNode();
    while (node) {
      const nodeLength = node.nodeValue ? node.nodeValue.length : 0;
      if (currentOffset + nodeLength >= offset) {
        const range = document.createRange();
        range.setStart(node, Math.max(0, offset - currentOffset));
        range.collapse(true);
        selection.removeAllRanges();
        selection.addRange(range);
        return;
      }
      currentOffset += nodeLength;
      node = walker.nextNode();
    }
    if (container.lastChild) {
      const range = document.createRange();
      range.selectNodeContents(container);
      range.collapse(false);
      selection.removeAllRanges();
      selection.addRange(range);
    }
  };

  const resolveCollabUrl = () => {
    if (!props.collabUrl) {
      return '';
    }
    if (props.collabUrl.startsWith('ws://') || props.collabUrl.startsWith('wss://')) {
      return props.collabUrl;
    }
    const protocol = window.location.protocol === 'https:' ? 'wss' : 'ws';
    const host = window.location.host;
    const path = props.collabUrl.startsWith('/') ? props.collabUrl : `/${props.collabUrl}`;
    return `${protocol}://${host}${path}`;
  };

  const getAwarenessPeerCount = () => {
    const provider = providerRef.value;
    if (!provider || !provider.awareness) {
      return 0;
    }
    try {
      return provider.awareness.getStates().size;
    } catch (error) {
      return 0;
    }
  };

  const syncContentToYjs = (content = '') => {
    if (!props.collabEnabled || !ytextRef.value) {
      return;
    }
    ytextRef.value.delete(0, ytextRef.value.length);
    ytextRef.value.insert(0, content);
  };

  const applyRemoteValue = (remoteValue) => {
    const vditor = vditorRef.value;
    if (!vditor || !isVditorReady.value || typeof vditor.setValue !== 'function') {
      return;
    }
    if (suppressInput.value) {
      return;
    }

    const currentValue = getValueSafely('');
    if (currentValue === remoteValue) {
      emit('update:modelValue', remoteValue);
      return;
    }

    const element = getEditableElement() || editorRef.value;
    const hasFocus = element && document.activeElement && element.contains(document.activeElement);
    const caretOffset = hasFocus ? getCaretOffset(element) : null;

    suppressInput.value = true;
    setValueSafely(remoteValue);
    emit('update:modelValue', remoteValue);
    suppressInput.value = false;

    if (hasFocus && element) {
      requestAnimationFrame(() => {
        restoreCaretOffset(element, caretOffset);
      });
    }
  };

  const setupCollaboration = () => {
    if (!props.collabEnabled || !props.collabRoom) {
      return;
    }
    const url = resolveCollabUrl();
    if (!url) {
      return;
    }

    ydocRef.value = new Y.Doc();
    ytextRef.value = ydocRef.value.getText('content');
    ycommentMetaRef.value = ydocRef.value.getMap('comment_meta');
    commentMetaObserverRef.value = (event) => {
      if (event?.transaction?.local) {
        return;
      }
      emit('comment-changed');
    };
    ycommentMetaRef.value.observe(commentMetaObserverRef.value);
    activeRoomRef.value = props.collabRoom;

    providerRef.value = new WebsocketProvider(url, props.collabRoom, ydocRef.value, {
      params: props.collabToken ? { token: props.collabToken } : {}
    });

    providerRef.value.on('sync', (isSynced) => {
      collabSyncedRef.value = isSynced;
      emit('collab-sync', isSynced);
      if (!isSynced || !ytextRef.value) {
        return;
      }
      const remote = ytextRef.value.toString();
      if (ytextRef.value.length === 0) {
        const peerCount = getAwarenessPeerCount();
        if (peerCount > 1) {
          return;
        }
        const seed = pendingSeedRef.value || props.modelValue || '';
        if (seed) {
          ytextRef.value.insert(0, seed);
        }
      } else if (remote && remote !== props.modelValue) {
        applyRemoteValue(remote);
      }
      pendingSeedRef.value = '';
    });

    ytextRef.value.observe((event) => {
      if (event?.transaction?.local) {
        return;
      }
      const remoteValue = ytextRef.value?.toString() || '';
      applyRemoteValue(remoteValue);
    });
  };

  const teardownCollaboration = () => {
    if (ycommentMetaRef.value && commentMetaObserverRef.value) {
      ycommentMetaRef.value.unobserve(commentMetaObserverRef.value);
    }
    ycommentMetaRef.value = null;
    commentMetaObserverRef.value = null;

    if (providerRef.value) {
      providerRef.value.destroy();
      providerRef.value = null;
    }

    if (ydocRef.value) {
      ydocRef.value.destroy();
      ydocRef.value = null;
    }

    ytextRef.value = null;
    activeRoomRef.value = '';
    collabSyncedRef.value = false;
    pendingSeedRef.value = '';
    emit('collab-sync', false);
  };

  const handleModelValueChange = (newValue) => {
    if (props.collabEnabled && ytextRef.value) {
      if (ytextRef.value.toString() === newValue) {
        return;
      }
      if (!collabSyncedRef.value) {
        pendingSeedRef.value = newValue || '';
        return;
      }
      if (collabSyncedRef.value && ytextRef.value.length === 0) {
        const peerCount = getAwarenessPeerCount();
        if (peerCount > 1) {
          return;
        }
        const seed = newValue || pendingSeedRef.value || '';
        pendingSeedRef.value = '';
        if (seed) {
          ytextRef.value.insert(0, seed);
        }
        return;
      }
      if (collabSyncedRef.value && ytextRef.value.length > 0 && activeRoomRef.value === props.collabRoom) {
        return;
      }
    }

    const vditor = vditorRef.value;
    if (!vditor || !isVditorReady.value || typeof vditor.getValue !== 'function') {
      return;
    }

    const currentValue = getValueSafely('');
    if (currentValue !== newValue) {
      setValueSafely(newValue);
    }
  };

  const handleRoomChange = () => {
    if (!props.collabEnabled) {
      return;
    }
    teardownCollaboration();
    setupCollaboration();
  };

  return {
    ydocRef,
    syncContentToYjs,
    setupCollaboration,
    teardownCollaboration,
    handleModelValueChange,
    handleRoomChange
  };
}
