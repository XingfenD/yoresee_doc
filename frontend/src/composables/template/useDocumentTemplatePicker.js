import { computed, nextTick, ref, watch } from 'vue';
import { useI18n } from 'vue-i18n';
import { useTemplateCatalog } from '@/composables/template/useTemplateCatalog';
import { matchDocumentType, normalizeDocumentType } from '@/utils/documentType';

export function useDocumentTemplatePicker(props) {
  const { t } = useI18n();

  const activeScope = ref('recent');
  const keyword = ref('');
  const showKnowledgeBaseTemplates = computed(() => Boolean(props.knowledgeBaseId));
  const selectedDocumentType = computed(() => normalizeDocumentType(props.documentType, ''));

  const normalizeScopeToTab = (scope) => {
    if (scope === 'system') {
      return 'public';
    }
    if (scope === 'private') {
      return 'my';
    }
    if (scope === 'knowledge_base') {
      return 'knowledge_base';
    }
    return 'recent';
  };

  const {
    recentTemplates,
    myTemplates,
    publicTemplates,
    kbTemplates,
    loadingRecent,
    loadingMy,
    loadingPublic,
    loadingKb,
    ensureLoaded: ensureTemplateLoaded,
    invalidateScope: invalidateTemplateScope
  } = useTemplateCatalog({
    includeKnowledgeBase: true,
    knowledgeBaseId: computed(() => props.knowledgeBaseId || ''),
    documentType: selectedDocumentType,
    onError: (error, scope) => {
      console.error(`[DocumentTemplatePicker] load ${scope} templates failed`, error);
    }
  });

  const scopeOptions = computed(() => {
    const base = [
      { name: 'recent', label: t('templates.recent') },
      { name: 'my', label: t('templates.my') },
      { name: 'public', label: t('templates.public') }
    ];
    if (showKnowledgeBaseTemplates.value) {
      base.push({ name: 'knowledge_base', label: t('templates.knowledgeBaseTab') });
    }
    return base;
  });

  const blankTemplateOption = computed(() => ({
    id: '__blank__',
    name: t('templates.blankDocument'),
    description: t('templates.blankDocumentDesc'),
    tags: [],
    is_blank: true
  }));

  const currentTemplates = computed(() => {
    if (activeScope.value === 'my') {
      return myTemplates.value;
    }
    if (activeScope.value === 'public') {
      return publicTemplates.value;
    }
    if (activeScope.value === 'knowledge_base') {
      return kbTemplates.value;
    }
    return recentTemplates.value;
  });

  const currentLoading = computed(() => {
    if (activeScope.value === 'my') {
      return loadingMy.value;
    }
    if (activeScope.value === 'public') {
      return loadingPublic.value;
    }
    if (activeScope.value === 'knowledge_base') {
      return loadingKb.value;
    }
    return loadingRecent.value;
  });

  const currentEmptyText = computed(() => {
    if (activeScope.value === 'my') {
      return t('templates.noMy');
    }
    if (activeScope.value === 'public') {
      return t('templates.noPublic');
    }
    if (activeScope.value === 'knowledge_base') {
      return t('templates.noKb');
    }
    return t('templates.noRecent');
  });

  const filteredTemplates = computed(() => {
    const text = keyword.value.trim().toLowerCase();
    const blank = blankTemplateOption.value;
    const typedTemplates = currentTemplates.value.filter((item) =>
      matchDocumentType(item?.type, selectedDocumentType.value)
    );
    if (!text) {
      return [blank, ...typedTemplates];
    }
    const list = typedTemplates.filter((item) => {
      const name = String(item?.name || '').toLowerCase();
      const description = String(item?.description || '').toLowerCase();
      const tags = Array.isArray(item?.tags) ? item.tags.join(' ').toLowerCase() : '';
      return name.includes(text) || description.includes(text) || tags.includes(text);
    });
    const blankMatched = [blank.name, blank.description].join(' ').toLowerCase().includes(text);
    if (blankMatched) {
      return [blank, ...list];
    }
    return list;
  });

  const selectedTemplateIdForPane = computed(() => props.selectedTemplateId || '__blank__');

  const initialScope = computed(() => {
    const mapped = normalizeScopeToTab(props.preferredScope);
    if (mapped === 'knowledge_base' && !showKnowledgeBaseTemplates.value) {
      return 'recent';
    }
    return mapped;
  });

  const fetchTemplates = async (scope) => {
    if (!props.visible) {
      return;
    }
    if (scope === 'knowledge_base' && !showKnowledgeBaseTemplates.value) {
      return;
    }
    await ensureTemplateLoaded(scope);
  };

  watch(
    () => props.visible,
    async (nextVisible) => {
      if (!nextVisible) {
        return;
      }
      await nextTick();
      activeScope.value = initialScope.value;
      keyword.value = '';
      if (activeScope.value === 'recent') {
        fetchTemplates('recent');
        return;
      }
      fetchTemplates('recent');
      fetchTemplates(activeScope.value);
    },
    { immediate: true }
  );

  watch(
    () => [props.visible, props.preferredScope, props.selectedTemplateId],
    ([visible]) => {
      if (!visible) {
        return;
      }
      if (!props.selectedTemplateId) {
        return;
      }
      const preferred = initialScope.value;
      if (preferred === 'recent') {
        return;
      }
      if (activeScope.value === 'recent') {
        activeScope.value = preferred;
        fetchTemplates(preferred);
      }
    }
  );

  watch(
    () => props.knowledgeBaseId,
    (next, prev) => {
      if (next === prev) {
        return;
      }
      invalidateTemplateScope('knowledge_base');
      if (activeScope.value === 'knowledge_base' && props.visible) {
        fetchTemplates('knowledge_base');
      }
    }
  );

  watch(showKnowledgeBaseTemplates, (show) => {
    if (!show && activeScope.value === 'knowledge_base') {
      activeScope.value = 'recent';
    }
  });

  watch(activeScope, (scope) => {
    fetchTemplates(scope);
  });

  watch(
    () => props.documentType,
    () => {
      if (!props.visible) {
        return;
      }
      fetchTemplates(activeScope.value);
    }
  );

  return {
    activeScope,
    keyword,
    scopeOptions,
    currentLoading,
    filteredTemplates,
    selectedTemplateIdForPane,
    currentEmptyText
  };
}
