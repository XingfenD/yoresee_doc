import { computed, onBeforeUnmount, onMounted, ref, unref } from 'vue';

export function useEditorFullscreen(options = {}) {
  const {
    editorLayoutRef,
    docId,
    onChange = null
  } = options;

  const isEditorFullscreen = ref(false);
  const canToggleEditorFullscreen = computed(() => Boolean(unref(docId)));
  let fullscreenHandler = null;

  const syncEditorFullscreenState = () => {
    const target = unref(editorLayoutRef);
    isEditorFullscreen.value = Boolean(target) && document.fullscreenElement === target;
    if (typeof onChange === 'function') {
      onChange(isEditorFullscreen.value);
    }
  };

  const toggleEditorFullscreen = async () => {
    const target = unref(editorLayoutRef);
    if (!target) {
      return;
    }
    try {
      if (document.fullscreenElement === target) {
        await document.exitFullscreen();
      } else {
        await target.requestFullscreen();
      }
    } catch (error) {
      // ignore unsupported browser fullscreen errors
    } finally {
      syncEditorFullscreenState();
    }
  };

  onMounted(() => {
    fullscreenHandler = () => syncEditorFullscreenState();
    document.addEventListener('fullscreenchange', fullscreenHandler);
    syncEditorFullscreenState();
  });

  onBeforeUnmount(() => {
    if (!fullscreenHandler) {
      return;
    }
    document.removeEventListener('fullscreenchange', fullscreenHandler);
    fullscreenHandler = null;
  });

  return {
    isEditorFullscreen,
    canToggleEditorFullscreen,
    toggleEditorFullscreen,
    syncEditorFullscreenState
  };
}
