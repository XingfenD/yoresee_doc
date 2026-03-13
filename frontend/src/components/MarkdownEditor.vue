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
  }
});

const emit = defineEmits(['update:modelValue']);

const editorRef = ref(null);
let vditor = null;
let themeObserver = null;
let ydoc = null;
let provider = null;
let ytext = null;
let isApplyingRemote = false;
let isVditorReady = false;
let collabSynced = false;
let suppressInput = false;
const debugCollab = false;
let pendingSeed = '';

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
  provider = new WebsocketProvider(url, props.collabRoom, ydoc, {
    params: props.collabToken ? { token: props.collabToken } : {}
  });

  provider.on('sync', (isSynced) => {
    collabSynced = isSynced;
    if (debugCollab) {
      console.log('[collab] sync', { isSynced, room: props.collabRoom });
    }
    if (!isSynced || !ytext) {
      return;
    }
    const remote = ytext.toString();
    if (ytext.length === 0) {
      const seed = pendingSeed || props.modelValue || '';
      if (seed) {
        ytext.insert(0, seed);
        if (debugCollab) {
          console.log('[collab] seed from pending', { length: seed.length });
        }
      }
    } else if (pendingSeed && remote !== pendingSeed) {
      if (debugCollab) {
        console.log('[collab] remote wins over pending seed', { remoteLength: remote.length });
      }
      emit('update:modelValue', remote);
      if (vditor && isVditorReady) {
        isApplyingRemote = true;
        vditor.setValue(remote);
        isApplyingRemote = false;
      }
    } else if (remote && remote !== props.modelValue) {
      if (debugCollab) {
        console.log('[collab] align local to remote', { remoteLength: remote.length });
      }
      emit('update:modelValue', remote);
      if (vditor && isVditorReady) {
        isApplyingRemote = true;
        vditor.setValue(remote);
        isApplyingRemote = false;
      }
    }
    pendingSeed = '';
  });

  ytext.observe(() => {
    if (!vditor || !isVditorReady || typeof vditor.setValue !== 'function') {
      return;
    }
    if (suppressInput) {
      return;
    }
    if (debugCollab) {
      console.log('[collab] ytext update', { length: ytext.length, room: props.collabRoom });
    }
    isApplyingRemote = true;
    vditor.setValue(ytext.toString());
    emit('update:modelValue', ytext.toString());
    isApplyingRemote = false;
  });
};

const teardownCollaboration = () => {
  if (provider) {
    provider.destroy();
    provider = null;
  }
  if (ydoc) {
    ydoc.destroy();
    ydoc = null;
    ytext = null;
  }
  collabSynced = false;
  pendingSeed = '';
};

onMounted(() => {
  vditor = new Vditor(editorRef.value, {
    height: props.height,
    value: props.modelValue,
    placeholder: props.placeholder,
    mode: 'wysiwyg',
    theme: 'classic',
    icon: 'ant',
    customWysiwygToolbar: () => [],
    counter: {
      enable: true
    },
    cache: {
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
      if (debugCollab) {
        console.log('[collab] vditor ready', { room: props.collabRoom });
      }
      applyVditorTheme();
    },
    input: (value) => {
      if (suppressInput) {
        return;
      }
      if (props.collabEnabled && ytext && !isApplyingRemote) {
        if (debugCollab) {
          console.log('[collab] input -> ytext', { length: value.length, room: props.collabRoom });
        }
        ytext.delete(0, ytext.length);
        ytext.insert(0, value);
      }
      emit('update:modelValue', value);
    }
  });

  setupCollaboration();
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
    // Defer seeding until sync; never overwrite active collaborative edits.
    if (!collabSynced || ytext.length === 0) {
      pendingSeed = newValue || '';
      if (debugCollab) {
        console.log('[collab] modelValue pending seed', { length: pendingSeed.length, room: props.collabRoom });
      }
    }
    return;
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
    vditor.setValue(newValue);
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

onMounted(() => {
  themeObserver = new MutationObserver(() => {
    applyVditorTheme();
  });
  themeObserver.observe(document.documentElement, {
    attributes: true,
    attributeFilter: ['class']
  });
});
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

.dark-mode .markdown-editor :deep(.vditor-ir__marker) {
  color: var(--text-light);
}

.markdown-editor :deep(.vditor-ir__heading) {
  color: var(--text-dark);
  font-weight: 600;
}

.dark-mode .markdown-editor :deep(.vditor-ir__heading) {
  color: var(--text-dark);
}

.dark-mode .markdown-editor :deep(.vditor-content h1),
.dark-mode .markdown-editor :deep(.vditor-content h2),
.dark-mode .markdown-editor :deep(.vditor-content h3),
.dark-mode .markdown-editor :deep(.vditor-content h4),
.dark-mode .markdown-editor :deep(.vditor-content h5),
.dark-mode .markdown-editor :deep(.vditor-content h6),
.dark-mode .markdown-editor :deep(.vditor-wysiwyg h1),
.dark-mode .markdown-editor :deep(.vditor-wysiwyg h2),
.dark-mode .markdown-editor :deep(.vditor-wysiwyg h3),
.dark-mode .markdown-editor :deep(.vditor-wysiwyg h4),
.dark-mode .markdown-editor :deep(.vditor-wysiwyg h5),
.dark-mode .markdown-editor :deep(.vditor-wysiwyg h6) {
  color: var(--text-dark);
}

.markdown-editor :deep(.vditor-ir__codeblock) {
  background-color: var(--bg-light);
  color: var(--text-dark);
}

.dark-mode .markdown-editor :deep(.vditor-ir__codeblock) {
  background-color: var(--bg-medium);
}

.markdown-editor :deep(.vditor-ir__quote) {
  border-left-color: var(--primary-color);
  color: var(--text-medium);
}

.dark-mode .markdown-editor :deep(.vditor-ir__quote) {
  color: var(--text-medium);
}

.markdown-editor :deep(.vditor-ir__strong) {
  font-weight: 600;
  color: var(--text-dark);
}

.dark-mode .markdown-editor :deep(.vditor-ir__strong) {
  color: var(--text-dark);
}

.markdown-editor :deep(.vditor-ir__em) {
  font-style: italic;
  color: var(--text-dark);
}

.dark-mode .markdown-editor :deep(.vditor-ir__em) {
  color: var(--text-dark);
}

.markdown-editor :deep(.vditor-wysiwyg) {
  background-color: var(--bg-white);
  color: var(--text-dark);
}

.dark-mode .markdown-editor :deep(.vditor-wysiwyg) {
  background-color: var(--bg-medium);
  color: var(--text-dark);
}

.dark-mode .markdown-editor :deep(.vditor-content p),
.dark-mode .markdown-editor :deep(.vditor-content li),
.dark-mode .markdown-editor :deep(.vditor-content blockquote),
.dark-mode .markdown-editor :deep(.vditor-content td),
.dark-mode .markdown-editor :deep(.vditor-content th),
.dark-mode .markdown-editor :deep(.vditor-ir__node),
.dark-mode .markdown-editor :deep(.vditor-wysiwyg),
.dark-mode .markdown-editor :deep(.vditor-wysiwyg p),
.dark-mode .markdown-editor :deep(.vditor-wysiwyg li),
.dark-mode .markdown-editor :deep(.vditor-wysiwyg blockquote),
.dark-mode .markdown-editor :deep(.vditor-wysiwyg td),
.dark-mode .markdown-editor :deep(.vditor-wysiwyg th) {
  color: var(--text-dark);
}

.markdown-editor :deep(.vditor-wysiwyg pre) {
  background-color: var(--bg-light);
}

.dark-mode .markdown-editor :deep(.vditor-wysiwyg pre) {
  background-color: var(--bg-medium);
}

.markdown-editor :deep(.vditor-wysiwyg code) {
  background-color: var(--bg-light);
  color: var(--text-dark);
}

.dark-mode .markdown-editor :deep(.vditor-wysiwyg code) {
  background-color: var(--bg-medium);
}

.markdown-editor :deep(.vditor-wysiwyg blockquote) {
  border-left-color: var(--primary-color);
  color: var(--text-medium);
}

.dark-mode .markdown-editor :deep(.vditor-wysiwyg blockquote) {
  color: var(--text-medium);
}

.markdown-editor :deep(.vditor-wysiwyg a) {
  color: var(--primary-color);
}

.markdown-editor :deep(.vditor-wysiwyg table) {
  border-color: var(--border-color);
}

.markdown-editor :deep(.vditor-wysiwyg th) {
  background-color: var(--bg-light);
  border-color: var(--border-color);
}

.dark-mode .markdown-editor :deep(.vditor-wysiwyg th) {
  background-color: var(--bg-medium);
}

.markdown-editor :deep(.vditor-wysiwyg td) {
  border-color: var(--border-color);
}

.dark-mode .markdown-editor :deep(.vditor-wysiwyg td) {
  background-color: var(--bg-light);
}

.markdown-editor :deep(.vditor-preview) {
  background-color: var(--bg-light);
}

.dark-mode .markdown-editor :deep(.vditor-preview) {
  background-color: var(--bg-medium);
}

.markdown-editor :deep(.vditor-preview h1),
.markdown-editor :deep(.vditor-preview h2),
.markdown-editor :deep(.vditor-preview h3) {
  color: var(--text-dark);
}

.dark-mode .markdown-editor :deep(.vditor-preview h1),
.dark-mode .markdown-editor :deep(.vditor-preview h2),
.dark-mode .markdown-editor :deep(.vditor-preview h3) {
  color: var(--text-dark);
}

.markdown-editor :deep(.vditor-preview code) {
  background-color: var(--bg-white);
  color: var(--text-dark);
}

.dark-mode .markdown-editor :deep(.vditor-preview code) {
  background-color: var(--bg-medium);
}

.markdown-editor :deep(.vditor-preview pre) {
  background-color: var(--bg-white);
}

.dark-mode .markdown-editor :deep(.vditor-preview pre) {
  background-color: var(--bg-medium);
}

.markdown-editor :deep(.vditor-preview blockquote) {
  border-left-color: var(--primary-color);
  color: var(--text-medium);
}

.dark-mode .markdown-editor :deep(.vditor-preview blockquote) {
  color: var(--text-medium);
}

.markdown-editor :deep(.vditor-preview a) {
  color: var(--primary-color);
}

.markdown-editor :deep(.vditor-preview table) {
  border-color: var(--border-color);
}

.markdown-editor :deep(.vditor-preview th) {
  background-color: var(--bg-light);
  border-color: var(--border-color);
}

.dark-mode .markdown-editor :deep(.vditor-preview th) {
  background-color: var(--bg-medium);
}

.markdown-editor :deep(.vditor-preview td) {
  border-color: var(--border-color);
}

.dark-mode .markdown-editor :deep(.vditor-preview td) {
  background-color: var(--bg-light);
}

.markdown-editor :deep(.vditor-statusbar) {
  background-color: var(--bg-white);
  border-top: 1px solid var(--border-color);
  color: var(--text-light);
}

.dark-mode .markdown-editor :deep(.vditor-statusbar) {
  background-color: var(--bg-medium);
  border-color: var(--border-color);
  color: var(--text-light);
}
</style>
