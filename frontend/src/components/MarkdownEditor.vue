<template>
  <div class="markdown-editor">
    <div ref="editorRef" class="vditor-container"></div>
  </div>
</template>

<script setup>
import { ref, onMounted, onBeforeUnmount, watch } from 'vue';
import Vditor from 'vditor';
import 'vditor/dist/index.css';
import * as Y from 'yjs';
import { WebsocketProvider } from 'y-websocket';

const props = defineProps({
  modelValue: {
    type: String,
    default: ''
  },
  placeholder: {
    type: String,
    default: ''
  },
  height: {
    type: [String, Number],
    default: '100%'
  },
  collabEnabled: {
    type: Boolean,
    default: false
  },
  collabRoom: {
    type: String,
    default: ''
  },
  collabUrl: {
    type: String,
    default: '/collab'
  },
  collabToken: {
    type: String,
    default: ''
  },
  commentEnabled: {
    type: Boolean,
    default: false
  }
});

const emit = defineEmits([
  'update:modelValue',
  'collab-sync',
  'ready',
  'comment-add',
  'comment-remove',
  'comment-scroll',
  'comment-adjust',
  'comment-changed'
]);

const editorRef = ref(null);
let vditor = null;
let themeObserver = null;
let ydoc = null;
let provider = null;
let ytext = null;
let ycommentMeta = null;
let commentMetaObserver = null;
let activeRoom = '';
let isVditorReady = false;
let suppressInput = false;
let collabSynced = false;
let pendingSeed = '';

const getVditorInstance = () => vditor;
const removeCommentIds = (ids = []) => {
  if (!vditor || typeof vditor.removeCommentIds !== 'function') {
    return;
  }
  if (!Array.isArray(ids) || ids.length === 0) {
    return;
  }
  vditor.removeCommentIds(ids);
  queueMicrotask(() => {
    syncEditorToYjs();
  });
};
const broadcastCommentChange = () => {
  if (!props.collabEnabled || !ydoc) {
    return;
  }
  if (!ycommentMeta) {
    ycommentMeta = ydoc.getMap('comment_meta');
  }
  ycommentMeta.set('tick', `${Date.now()}_${Math.random().toString(36).slice(2)}`);
};
defineExpose({ getVditor: getVditorInstance, removeCommentIds, broadcastCommentChange });

const getEditableElement = () => {
  if (!vditor) {
    return null;
  }
  return vditor?.vditor?.wysiwyg?.element || vditor?.vditor?.ir?.element || editorRef.value;
};

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

const syncEditorToYjs = (fallbackValue = '') => {
  let nextValue = fallbackValue;
  if (vditor && typeof vditor.getValue === 'function') {
    try {
      nextValue = vditor.getValue();
    } catch (error) {
      nextValue = fallbackValue;
    }
  }
  if (props.collabEnabled && ytext) {
    ytext.delete(0, ytext.length);
    ytext.insert(0, nextValue);
  }
  emit('update:modelValue', nextValue);
};

const getAwarenessPeerCount = () => {
  if (!provider || !provider.awareness) {
    return 0;
  }
  try {
    return provider.awareness.getStates().size;
  } catch (error) {
    return 0;
  }
};

const applyVditorTheme = () => {
  if (!vditor || typeof vditor.setTheme !== 'function') {
    return;
  }
  try {
    const isDarkMode = document.documentElement.classList.contains('dark-mode');
    vditor.setTheme(
      isDarkMode ? 'dark' : 'classic',
      isDarkMode ? 'dark' : 'light',
      isDarkMode ? 'dark' : 'github'
    );
  } catch (error) {
    // Ignore theme apply errors during init/destroy
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

const setupCollaboration = () => {
  if (!props.collabEnabled || !props.collabRoom) {
    return;
  }
  const url = resolveCollabUrl();
  if (!url) {
    return;
  }

  ydoc = new Y.Doc();
  ytext = ydoc.getText('content');
  ycommentMeta = ydoc.getMap('comment_meta');
  commentMetaObserver = (event) => {
    if (event?.transaction?.local) {
      return;
    }
    emit('comment-changed');
  };
  ycommentMeta.observe(commentMetaObserver);
  activeRoom = props.collabRoom;
  provider = new WebsocketProvider(url, props.collabRoom, ydoc, {
    params: props.collabToken ? { token: props.collabToken } : {}
  });

  provider.on('sync', (isSynced) => {
    collabSynced = isSynced;
    emit('collab-sync', isSynced);
    if (!isSynced || !ytext) {
      return;
    }
    const remote = ytext.toString();
    if (ytext.length === 0) {
      const peerCount = getAwarenessPeerCount();
      if (peerCount > 1) {
        return;
      }
      const seed = pendingSeed || props.modelValue || '';
      if (seed) {
        ytext.insert(0, seed);
      }
    } else if (remote && remote !== props.modelValue) {
      emit('update:modelValue', remote);
      if (vditor && isVditorReady) {
        suppressInput = true;
        vditor.setValue(remote);
        suppressInput = false;
      }
    }
    pendingSeed = '';
  });

  ytext.observe((event) => {
    if (!vditor || !isVditorReady || typeof vditor.setValue !== 'function') {
      return;
    }
    if (suppressInput) {
      return;
    }
    if (event?.transaction?.local) {
      return;
    }
    const remoteValue = ytext.toString();
    let currentValue = '';
    try {
      currentValue = vditor.getValue();
    } catch (error) {
      currentValue = '';
    }
    if (currentValue === remoteValue) {
      emit('update:modelValue', remoteValue);
      return;
    }
    const editorElement = getEditableElement();
    const hasFocus = editorElement && document.activeElement && editorElement.contains(document.activeElement);
    const caretOffset = hasFocus ? getCaretOffset(editorElement) : null;
    suppressInput = true;
    vditor.setValue(remoteValue);
    emit('update:modelValue', remoteValue);
    suppressInput = false;
    if (hasFocus && editorElement) {
      requestAnimationFrame(() => {
        restoreCaretOffset(editorElement, caretOffset);
      });
    }
  });
};

const teardownCollaboration = () => {
  if (ycommentMeta && commentMetaObserver) {
    ycommentMeta.unobserve(commentMetaObserver);
  }
  ycommentMeta = null;
  commentMetaObserver = null;
  if (provider) {
    provider.destroy();
    provider = null;
  }
  if (ydoc) {
    ydoc.destroy();
    ydoc = null;
    ytext = null;
  }
  activeRoom = '';
  collabSynced = false;
  pendingSeed = '';
  emit('collab-sync', false);
};

onMounted(() => {
  vditor = new Vditor(editorRef.value, {
    height: props.height,
    value: props.modelValue,
    placeholder: props.placeholder,
    mode: 'wysiwyg',
    theme: 'classic',
    icon: 'ant',
    toolbar: [
      'headings',
      'bold',
      'italic',
      'strike',
      '|',
      'list',
      'ordered-list',
      'check',
      '|',
      'link',
      'quote',
      'code',
      'table',
      '|',
      'edit-mode',
      'undo',
      'redo',
      'outline',
      'fullscreen'
    ],
    toolbarConfig: {
      hide: false
    },
    customWysiwygToolbar: () => [],
    counter: {
      enable: true
    },
    cache: {
      enable: false
    },
    comment: props.commentEnabled
      ? {
          enable: true,
          add: (id, text, commentsData) => {
            emit('comment-add', { id, text, commentsData });
            queueMicrotask(() => {
              syncEditorToYjs();
            });
          },
          remove: (ids) => {
            emit('comment-remove', ids);
            queueMicrotask(() => {
              syncEditorToYjs();
            });
          },
          scroll: (top) => {
            emit('comment-scroll', top);
          },
          adjustTop: (commentsData) => {
            emit('comment-adjust', commentsData);
          }
        }
      : {
          enable: false
        },
    upload: {
      handler: (files) => {
        return Promise.reject('上传功能暂未实现');
      }
    },
    after: () => {
      suppressInput = true;
      vditor.setValue(props.modelValue);
      suppressInput = false;
      isVditorReady = true;
      applyVditorTheme();
      emit('ready', vditor);
    },
    input: (value) => {
      if (suppressInput) {
        return;
      }
      queueMicrotask(() => {
        syncEditorToYjs(value);
      });
    }
  });

  setupCollaboration();
  themeObserver = new MutationObserver(() => {
    applyVditorTheme();
  });
  themeObserver.observe(document.documentElement, {
    attributes: true,
    attributeFilter: ['class']
  });
});

onBeforeUnmount(() => {
  if (themeObserver) {
    themeObserver.disconnect();
    themeObserver = null;
  }
  if (vditor) {
    vditor.destroy();
    vditor = null;
  }
  teardownCollaboration();
});

watch(() => props.modelValue, (newValue) => {
  if (props.collabEnabled && ytext) {
    if (ytext.toString() === newValue) {
      return;
    }
    if (!collabSynced) {
      pendingSeed = newValue || '';
      return;
    }
    if (collabSynced && ytext.length === 0) {
      const peerCount = getAwarenessPeerCount();
      if (peerCount > 1) {
        return;
      }
      const seed = newValue || pendingSeed || '';
      pendingSeed = '';
      if (seed) {
        ytext.insert(0, seed);
      }
      return;
    }
    if (collabSynced && ytext.length > 0 && activeRoom === props.collabRoom) {
      return;
    }
  }
  if (!vditor || !isVditorReady || typeof vditor.getValue !== 'function') {
    return;
  }
  let currentValue = '';
  try {
    currentValue = vditor.getValue();
  } catch (error) {
    return;
  }
  if (currentValue !== newValue) {
    suppressInput = true;
    vditor.setValue(newValue);
    suppressInput = false;
  }
});

watch(
  () => props.collabRoom,
  () => {
    if (!props.collabEnabled) {
      return;
    }
    teardownCollaboration();
    setupCollaboration();
  }
);
</script>

<style scoped>
.markdown-editor {
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
}

.vditor-container {
  flex: 1;
  min-height: 500px;
}

.markdown-editor :deep(.vditor) {
  border: none;
  background-color: var(--bg-white);
}

.dark-mode .markdown-editor :deep(.vditor) {
  background-color: var(--bg-white);
}

.markdown-editor :deep(.vditor-toolbar) {
  background-color: var(--bg-white);
  border-bottom: 1px solid var(--border-color);
}

.dark-mode .markdown-editor :deep(.vditor-toolbar) {
  background-color: var(--bg-medium);
  border-color: var(--border-color);
}

.markdown-editor :deep(.vditor-toolbar__item) {
  color: var(--text-medium);
}

.dark-mode .markdown-editor :deep(.vditor-toolbar__item) {
  color: var(--text-medium);
}

.markdown-editor :deep(.vditor-toolbar__item:hover) {
  color: var(--primary-color);
}

.markdown-editor :deep(.vditor-toolbar__item--current) {
  color: var(--primary-color);
}

.markdown-editor :deep(.vditor-content) {
  background-color: var(--bg-white);
  color: var(--text-dark);
}

.dark-mode .markdown-editor :deep(.vditor-content) {
  background-color: var(--bg-medium);
  color: var(--text-dark);
}

.markdown-editor :deep(.vditor-ir) {
  background-color: var(--bg-white);
  color: var(--text-dark);
}

.dark-mode .markdown-editor :deep(.vditor-ir) {
  background-color: var(--bg-medium);
  color: var(--text-dark);
}

.markdown-editor :deep(.vditor-ir__node) {
  color: var(--text-dark);
}

.dark-mode .markdown-editor :deep(.vditor-ir__node) {
  color: var(--text-dark);
}

.markdown-editor :deep(.vditor-ir__link) {
  color: var(--primary-color);
}

.markdown-editor :deep(.vditor-ir__link:hover) {
  color: var(--primary-color);
}

.markdown-editor :deep(.vditor-ir__marker) {
  color: var(--text-light);
}

.markdown-editor :deep(.vditor-ir__marker--heading) {
  color: var(--text-medium);
}

.dark-mode .markdown-editor :deep(.vditor-ir__marker) {
  color: var(--text-light);
}
</style>
