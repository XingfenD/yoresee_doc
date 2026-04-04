import { computed, ref, unref, watch } from 'vue';
import { listRecentTemplates, listTemplates } from '@/services/api';

const normalizeScope = (scope) => {
  if (scope === 'kb') {
    return 'knowledge_base';
  }
  return scope;
};

const resolveRefValue = (valueOrRefOrGetter) => {
  if (typeof valueOrRefOrGetter === 'function') {
    return valueOrRefOrGetter();
  }
  return unref(valueOrRefOrGetter);
};

export function useTemplateCatalog(options = {}) {
  const {
    includeKnowledgeBase = false,
    knowledgeBaseId = '',
    documentType = '',
    pageSize = 50,
    onError
  } = options;

  const recentTemplates = ref([]);
  const myTemplates = ref([]);
  const publicTemplates = ref([]);
  const kbTemplates = ref([]);

  const loadingRecent = ref(false);
  const loadingMy = ref(false);
  const loadingPublic = ref(false);
  const loadingKb = ref(false);

  const loadedRecent = ref(false);
  const loadedMy = ref(false);
  const loadedPublic = ref(false);
  const loadedKb = ref(false);

  const templatesMap = {
    recent: recentTemplates,
    my: myTemplates,
    public: publicTemplates,
    knowledge_base: kbTemplates
  };

  const loadingMap = {
    recent: loadingRecent,
    my: loadingMy,
    public: loadingPublic,
    knowledge_base: loadingKb
  };

  const loadedMap = {
    recent: loadedRecent,
    my: loadedMy,
    public: loadedPublic,
    knowledge_base: loadedKb
  };

  const kbId = computed(() => `${resolveRefValue(knowledgeBaseId) || ''}`);
  const normalizedDocType = computed(() => `${resolveRefValue(documentType) || ''}`.trim());

  const requestTemplates = async (scope) => {
    if (scope === 'recent') {
      const resp = await listRecentTemplates({
        page: 1,
        page_size: pageSize
      });
      return resp.templates || [];
    }

    if (scope === 'my') {
      const resp = await listTemplates({
        only_mine: true,
        target_container: 'own',
        type: normalizedDocType.value || undefined,
        order_by: 'updated_at',
        order_desc: true,
        page: 1,
        page_size: pageSize
      });
      return resp.templates || [];
    }

    if (scope === 'public') {
      const resp = await listTemplates({
        target_container: 'public',
        type: normalizedDocType.value || undefined,
        order_by: 'updated_at',
        order_desc: true,
        page: 1,
        page_size: pageSize
      });
      return resp.templates || [];
    }

    if (scope === 'knowledge_base') {
      if (!includeKnowledgeBase || !kbId.value) {
        return [];
      }
      const resp = await listTemplates({
        target_container: 'knowledge_base',
        knowledge_base_id: kbId.value,
        type: normalizedDocType.value || undefined,
        order_by: 'updated_at',
        order_desc: true,
        page: 1,
        page_size: pageSize
      });
      return resp.templates || [];
    }

    return [];
  };

  const isLoaded = (scope) => {
    const normalizedScope = normalizeScope(scope);
    return Boolean(loadedMap[normalizedScope]?.value);
  };

  const invalidateScope = (scope, options = {}) => {
    const normalizedScope = normalizeScope(scope);
    const { clear = true } = options;
    if (!loadedMap[normalizedScope]) {
      return;
    }
    loadedMap[normalizedScope].value = false;
    if (clear && templatesMap[normalizedScope]) {
      templatesMap[normalizedScope].value = [];
    }
  };

  const ensureLoaded = async (scope, options = {}) => {
    const normalizedScope = normalizeScope(scope);
    const { force = false } = options;
    const templatesRef = templatesMap[normalizedScope];
    const loadingRef = loadingMap[normalizedScope];
    const loadedRef = loadedMap[normalizedScope];

    if (!templatesRef || !loadingRef || !loadedRef) {
      return [];
    }

    if (loadingRef.value) {
      return templatesRef.value;
    }

    if (!force && loadedRef.value) {
      return templatesRef.value;
    }

    loadingRef.value = true;
    try {
      const nextTemplates = await requestTemplates(normalizedScope);
      templatesRef.value = nextTemplates;
      loadedRef.value = true;
      return nextTemplates;
    } catch (error) {
      templatesRef.value = [];
      loadedRef.value = false;
      if (typeof onError === 'function') {
        onError(error, normalizedScope);
      } else {
        console.error(`[useTemplateCatalog] load ${normalizedScope} failed`, error);
      }
      return [];
    } finally {
      loadingRef.value = false;
    }
  };

  watch(
    kbId,
    (next, prev) => {
      if (next === prev) {
        return;
      }
      invalidateScope('knowledge_base');
    }
  );

  watch(
    normalizedDocType,
    (next, prev) => {
      if (next === prev) {
        return;
      }
      invalidateScope('my');
      invalidateScope('public');
      invalidateScope('knowledge_base');
    }
  );

  return {
    recentTemplates,
    myTemplates,
    publicTemplates,
    kbTemplates,
    loadingRecent,
    loadingMy,
    loadingPublic,
    loadingKb,
    loadedRecent,
    loadedMy,
    loadedPublic,
    loadedKb,
    ensureLoaded,
    invalidateScope,
    isLoaded
  };
}
