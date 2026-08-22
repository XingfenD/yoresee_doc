import { computed, reactive, ref, watch } from 'vue';
import { ElMessage } from 'element-plus';
import { useI18n } from 'vue-i18n';
import { getKnowledgeBases } from '@/services/api';
import { DEFAULT_DOCUMENT_TYPE, normalizeDocumentType } from '@/utils/documentType';

export function useDocumentCreateDialog(props, emit) {
  const { t } = useI18n();

  const dialogWidth = computed(() => {
    if (props.showTemplatePicker) {
      return '980px';
    }
    return '620px';
  });

  const visible = computed({
    get: () => props.modelValue,
    set: (value) => emit('update:modelValue', value)
  });

  const loadingKnowledgeBases = ref(false);
  const knowledgeBaseOptions = ref([]);

  const locationOptions = computed(() => [
    { label: t('document.createLocationOwn'), value: 'own' },
    { label: t('document.createLocationKnowledgeBase'), value: 'knowledge_base' }
  ]);

  const formState = reactive({
    title: '',
    isPublic: false,
    containerType: 'own',
    targetKnowledgeBaseId: '',
    template: '',
    parentExternalId: '',
    templateMeta: null,
    documentType: DEFAULT_DOCUMENT_TYPE
  });

  const loadKnowledgeBaseOptions = async () => {
    if (!props.showLocationSelector) {
      return;
    }
    loadingKnowledgeBases.value = true;
    try {
      const resp = await getKnowledgeBases({
        only_mine: true,
        page: 1,
        page_size: 200,
        order_by: 'updated_at',
        order_desc: true
      });
      knowledgeBaseOptions.value = (resp.knowledge_bases || []).map((item) => ({
        value: item.external_id,
        label: item.name || item.external_id
      }));
    } catch (error) {
      console.error('[DocumentCreateDialog] load knowledge base options failed', error);
      knowledgeBaseOptions.value = [];
    } finally {
      loadingKnowledgeBases.value = false;
    }
  };

  const resetForm = () => {
    formState.title = props.initialTitle || t('document.untitledDefaultTitle');
    formState.isPublic = Boolean(props.initialIsPublic);
    formState.containerType = props.initialContainerType === 'knowledge_base' ? 'knowledge_base' : 'own';
    formState.targetKnowledgeBaseId = props.initialTargetKnowledgeBaseId || props.knowledgeBaseId || '';
    const hasInitialTemplateId =
      props.initialTemplateId !== '' && props.initialTemplateId !== null && props.initialTemplateId !== undefined;
    formState.template = hasInitialTemplateId ? String(props.initialTemplateId) : '';
    formState.templateMeta = props.initialTemplateMeta || null;
    const inferredType = props.initialTemplateMeta?.type;
    formState.documentType = normalizeDocumentType(
      inferredType || props.initialDocumentType || DEFAULT_DOCUMENT_TYPE
    );
    formState.parentExternalId = props.parentExternalId || '';
  };

  const handleCancel = () => {
    emit('cancel');
    visible.value = false;
  };

  const handleCreate = () => {
    if (props.showTitleInput && !formState.title.trim()) {
      ElMessage.error(t('knowledgeBase.titleRequired'));
      return;
    }
    if (
      props.showLocationSelector &&
      formState.containerType === 'knowledge_base' &&
      !formState.targetKnowledgeBaseId
    ) {
      ElMessage.error(t('knowledgeBase.selectKnowledgeBase'));
      return;
    }

    const title = formState.title.trim() || t('document.untitledDefaultTitle');

    emit('submit', {
      title,
      type: formState.documentType,
      template: formState.template,
      template_meta: formState.templateMeta,
      parent_external_id: formState.parentExternalId || undefined,
      is_public: props.showPublicSwitch ? Boolean(formState.isPublic) : false,
      container_type: props.showLocationSelector ? formState.containerType : undefined,
      knowledge_base_external_id:
        props.showLocationSelector && formState.containerType === 'knowledge_base'
          ? formState.targetKnowledgeBaseId || undefined
          : undefined
    });
  };

  const selectTemplate = (tpl) => {
    if (tpl?.is_blank) {
      formState.template = '';
      formState.templateMeta = null;
      return;
    }
    const templateId = String(tpl.id);
    if (formState.template === templateId) {
      formState.template = '';
      formState.templateMeta = null;
      return;
    }
    formState.template = templateId;
    formState.templateMeta = tpl;
    if (tpl?.type !== undefined && tpl?.type !== null && tpl?.type !== '') {
      formState.documentType = normalizeDocumentType(tpl.type, formState.documentType);
    }
  };

  watch(
    () => props.modelValue,
    async (nextVisible) => {
      if (nextVisible) {
        resetForm();
        await loadKnowledgeBaseOptions();
      }
    }
  );

  return {
    dialogWidth,
    visible,
    loadingKnowledgeBases,
    knowledgeBaseOptions,
    locationOptions,
    formState,
    handleCancel,
    handleCreate,
    selectTemplate
  };
}
