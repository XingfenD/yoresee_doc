export function useDocumentHeaderRouting(options = {}) {
  const {
    router,
    kbId,
    docId,
    onCommand
  } = options;

  const pushDocumentSubRoute = (segment) => {
    if (!docId.value) {
      return;
    }
    if (kbId.value === 'personal') {
      router.push(`/mydocument/${docId.value}/${segment}`);
      return;
    }
    router.push(`/knowledge-base/${kbId.value}/document/${docId.value}/${segment}`);
  };

  const handleHeaderCommand = (command) => {
    if (typeof onCommand === 'function' && onCommand(command)) {
      return;
    }
    if (command === 'document_settings') {
      pushDocumentSubRoute('setting');
      return;
    }
    if (command === 'manage_attachments') {
      pushDocumentSubRoute('attachments');
      return;
    }
    if (command === 'show_history') {
      pushDocumentSubRoute('history');
    }
  };

  return {
    handleHeaderCommand
  };
}
