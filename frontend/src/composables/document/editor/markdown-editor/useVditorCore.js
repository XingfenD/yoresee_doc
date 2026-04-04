import Vditor from 'vditor';

const DEFAULT_TOOLBAR = [
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
  'outline'
];

export function useVditorCore({
  props,
  emit,
  editorRef,
  vditorRef,
  isVditorReady,
  suppressInput,
  getCommentOptions,
  onEditorInput
}) {
  let themeObserver = null;

  const getEditableElement = () => {
    const vditor = vditorRef.value;
    if (!vditor) {
      return null;
    }
    return vditor?.vditor?.wysiwyg?.element || vditor?.vditor?.ir?.element || editorRef.value;
  };

  const getValueSafely = (fallback = '') => {
    const vditor = vditorRef.value;
    if (!vditor || typeof vditor.getValue !== 'function') {
      return fallback;
    }
    try {
      return vditor.getValue();
    } catch (error) {
      return fallback;
    }
  };

  const setValueSafely = (value = '') => {
    const vditor = vditorRef.value;
    if (!vditor || typeof vditor.setValue !== 'function') {
      return;
    }
    suppressInput.value = true;
    vditor.setValue(value);
    suppressInput.value = false;
  };

  const applyVditorTheme = () => {
    const vditor = vditorRef.value;
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

  const initVditor = () => {
    const initIsDarkMode = document.documentElement.classList.contains('dark-mode');
    vditorRef.value = new Vditor(editorRef.value, {
      height: props.height,
      value: props.modelValue,
      placeholder: props.placeholder,
      mode: 'wysiwyg',
      theme: initIsDarkMode ? 'dark' : 'classic',
      icon: 'ant',
      toolbar: DEFAULT_TOOLBAR,
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
      comment: getCommentOptions(),
      upload: {
        handler: () => Promise.reject('上传功能暂未实现')
      },
      after: () => {
        setValueSafely(props.modelValue);
        isVditorReady.value = true;
        applyVditorTheme();
        emit('ready', vditorRef.value);
      },
      input: (value) => {
        if (suppressInput.value) {
          return;
        }
        onEditorInput?.(value);
      }
    });

    themeObserver = new MutationObserver(() => {
      applyVditorTheme();
    });
    themeObserver.observe(document.documentElement, {
      attributes: true,
      attributeFilter: ['class']
    });
  };

  const destroyVditor = () => {
    if (themeObserver) {
      themeObserver.disconnect();
      themeObserver = null;
    }
    if (vditorRef.value) {
      vditorRef.value.destroy();
      vditorRef.value = null;
    }
    isVditorReady.value = false;
  };

  return {
    getEditableElement,
    getValueSafely,
    setValueSafely,
    initVditor,
    destroyVditor
  };
}
