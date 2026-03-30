<template>
  <CardListSection
    class="document-list-section"
    :title="title"
    :items="items"
    :empty-text="emptyText"
    :action-label="primaryAction"
    :secondary-action-label="secondaryAction"
    :show-view-all="showViewAll"
    :view-all-label="t('common.viewAll')"
    :item-key-mapper="itemKeyMapper"
    :item-title-mapper="itemTitleMapper"
    :item-description-mapper="itemDescriptionMapper"
    :meta-mapper="metaMapper"
    @view-all="emit('view-all')"
    @open="handleView"
    @secondary-action="handleEdit"
  />
</template>

<script setup>
import { computed } from 'vue';
import { useI18n } from 'vue-i18n';
import CardListSection from '@/components/CardListSection.vue';

const { t } = useI18n();

const props = defineProps({
  title: {
    type: String,
    default: ''
  },
  items: {
    type: Array,
    default: () => []
  },
  emptyText: {
    type: String,
    default: ''
  },
  showViewAll: {
    type: Boolean,
    default: false
  },
  singleAction: {
    type: Boolean,
    default: false
  },
  primaryActionLabel: {
    type: String,
    default: ''
  }
});

const emit = defineEmits(['view-all', 'view-item', 'edit-item']);

const formatDate = (dateString) => {
  if (!dateString) return t('common.unknown');
  const date = new Date(dateString);
  if (Number.isNaN(date.getTime())) return dateString;
  return date.toLocaleDateString();
};

const itemKeyMapper = (doc) => doc?.id || doc?.external_id || doc?.externalId || doc?.title;
const itemTitleMapper = (doc) => doc?.title || t('document.title');
const itemDescriptionMapper = () => '';

const metaMapper = (doc) => [
  { label: t('knowledgeBase.owner'), value: doc?.author || t('common.unknown') },
  { label: t('common.updatedAt'), value: formatDate(doc?.updatedAt || doc?.updated_at) }
];

const primaryAction = computed(() => {
  if (props.singleAction) {
    return props.primaryActionLabel || t('common.open');
  }
  return t('document.view');
});

const secondaryAction = computed(() => {
  if (props.singleAction) {
    return '';
  }
  return t('document.edit');
});

const handleView = (doc) => {
  emit('view-item', doc);
};

const handleEdit = (doc) => {
  emit('edit-item', doc);
};
</script>
