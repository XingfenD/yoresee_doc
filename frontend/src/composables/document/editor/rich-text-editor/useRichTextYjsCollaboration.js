import { computed, onBeforeUnmount, onMounted, ref, watch } from 'vue';
import * as Y from 'yjs';
import { WebsocketProvider } from 'y-websocket';
import Collaboration from '@tiptap/extension-collaboration';
import {
  createYjsContentBridge,
  YJS_CONTENT_CARRIER,
  YJS_CONTENT_FIELD
} from '@/composables/document/editor/collab/useYjsContentBridge';

const resolveCollabUrl = (rawUrl) => {
  const input = String(rawUrl || '').trim();
  if (!input) {
    return '';
  }
  if (input.startsWith('ws://') || input.startsWith('wss://')) {
    return input;
  }
  const protocol = window.location.protocol === 'https:' ? 'wss' : 'ws';
  const host = window.location.host;
  const path = input.startsWith('/') ? input : `/${input}`;
  return `${protocol}://${host}${path}`;
};

const getAwarenessPeerCount = (provider) => {
  if (!provider?.awareness) {
    return 0;
  }
  try {
    return provider.awareness.getStates().size;
  } catch (_) {
    return 0;
  }
};

const hasSeedContent = (content) => {
  if (typeof content === 'string') {
    return content.trim().length > 0;
  }
  if (!content || typeof content !== 'object') {
    return false;
  }
  if (Array.isArray(content.content)) {
    return content.content.length > 0;
  }
  return true;
};

export function useRichTextYjsCollaboration(options = {}) {
  const {
    collabEnabledRef,
    collabRoomRef,
    collabUrlRef,
    collabTokenRef,
    modelValueRef,
    editorRef,
    resolveInitialEditorContent,
    onSync
  } = options;

  const ydocRef = ref(null);
  const contentBridgeRef = ref(null);
  const providerRef = ref(null);
  const collabSyncedRef = ref(false);
  const pendingSeedRef = ref(false);

  const emitSync = (isSynced) => {
    collabSyncedRef.value = Boolean(isSynced);
    onSync?.(Boolean(isSynced));
  };

  const seedEditorFromModelValue = () => {
    const instance = editorRef.value;
    if (!instance) {
      return false;
    }
    const seedContent = resolveInitialEditorContent(modelValueRef.value);
    if (!hasSeedContent(seedContent)) {
      return false;
    }
    instance.commands.setContent(seedContent, false, {
      preserveWhitespace: true
    });
    return true;
  };

  const maybeSeedAfterSync = () => {
    const bridge = contentBridgeRef.value;
    if (!collabSyncedRef.value) {
      return;
    }
    if (!bridge || !bridge.isEmpty()) {
      pendingSeedRef.value = false;
      return;
    }
    if (getAwarenessPeerCount(providerRef.value) > 1) {
      pendingSeedRef.value = false;
      return;
    }
    if (!editorRef.value) {
      pendingSeedRef.value = true;
      return;
    }
    pendingSeedRef.value = false;
    seedEditorFromModelValue();
  };

  const teardownCollaboration = () => {
    if (providerRef.value) {
      providerRef.value.destroy();
      providerRef.value = null;
    }
    if (ydocRef.value) {
      ydocRef.value.destroy();
      ydocRef.value = null;
    }
    contentBridgeRef.value = null;
    pendingSeedRef.value = false;
    emitSync(false);
  };

  const setupCollaboration = () => {
    teardownCollaboration();

    if (!collabEnabledRef.value || !collabRoomRef.value) {
      return;
    }

    const url = resolveCollabUrl(collabUrlRef.value);
    if (!url) {
      return;
    }

    ydocRef.value = new Y.Doc();
    contentBridgeRef.value = createYjsContentBridge({
      doc: ydocRef.value,
      carrier: YJS_CONTENT_CARRIER.XML_FRAGMENT
    });

    providerRef.value = new WebsocketProvider(url, collabRoomRef.value, ydocRef.value, {
      params: collabTokenRef.value ? { token: collabTokenRef.value } : {}
    });

    providerRef.value.on('sync', (isSynced) => {
      emitSync(isSynced);
      if (!isSynced) {
        return;
      }
      maybeSeedAfterSync();
    });
  };

  const collaborationExtensionsRef = computed(() => {
    if (!collabEnabledRef.value || !ydocRef.value) {
      return [];
    }
    return [
      Collaboration.configure({
        document: ydocRef.value,
        field: YJS_CONTENT_FIELD
      })
    ];
  });

  const collaborationBridgeRef = computed(() => {
    if (!ydocRef.value || !providerRef.value) {
      return null;
    }
    return {
      ydoc: ydocRef.value,
      provider: providerRef.value
    };
  });

  watch(editorRef, (instance) => {
    if (!instance || !pendingSeedRef.value) {
      return;
    }
    maybeSeedAfterSync();
  });

  onMounted(() => {
    setupCollaboration();
  });

  onBeforeUnmount(() => {
    teardownCollaboration();
  });

  return {
    collabSyncedRef,
    collaborationExtensionsRef,
    collaborationBridgeRef,
    setupCollaboration,
    teardownCollaboration
  };
}
