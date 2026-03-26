<template>
  <CardListSection
    :title="title"
    :items="items"
    :empty-text="emptyText"
    :fallback-description="fallbackDescription"
    :action-label="actionLabel"
    :tag-mapper="tagMapper"
    :meta-mapper="metaMapper"
    :item-key-mapper="(item) => item?.id || item?.external_id || item?.externalId"
    :item-title-mapper="(item) => item?.name || ''"
    :item-description-mapper="(item) => item?.description || ''"
    @open="$emit('open', $event)"
  />
</template>

<script setup>
import CardListSection from '@/components/CardListSection.vue';

defineProps({
  title: { type: String, default: '' },
  items: { type: Array, default: () => [] },
  emptyText: { type: String, default: '' },
  fallbackDescription: { type: String, default: '' },
  actionLabel: { type: String, default: '' },
  tagMapper: { type: Function, default: null },
  metaMapper: { type: Function, default: null }
});

defineEmits(['open']);
</script>
