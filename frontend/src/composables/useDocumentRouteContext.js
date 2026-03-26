import { computed, ref } from 'vue';

export function useDocumentRouteContext({ props, route }) {
  const kbId = ref(props.kbId || route.params.kbId);
  const docId = ref(props.docId || route.params.docId);

  const resolveActiveMenu = (currentKbId) => {
    if (currentKbId === 'personal') return 'documents';
    if (currentKbId) return 'knowledge-base';
    return 'home';
  };

  const activeMenu = ref(resolveActiveMenu(kbId.value));
  const collabEnabled = computed(() => !!docId.value && docId.value !== 'example');
  const collabRoom = computed(() => (docId.value ? `${docId.value}` : ''));
  const collabUrl = computed(() => '/ws/doc');
  const collabToken = computed(() => localStorage.getItem('token') || '');
  const collabReady = ref(false);
  const lastSyncedDocId = ref('');

  return {
    kbId,
    docId,
    activeMenu,
    resolveActiveMenu,
    collabEnabled,
    collabRoom,
    collabUrl,
    collabToken,
    collabReady,
    lastSyncedDocId
  };
}
