<template>
  <CardListSection
    class="vertical-section"
    :title="title"
    :items="items"
    :empty-text="emptyText"
    :tag-type="tagType"
    :tag-label="tagLabel"
    :tag-mapper="tagMapper"
    :fallback-description="fallbackDescription"
    :meta-mapper="metaMapper"
    :show-load-more="showLoadMore"
    :loading="loading"
    :load-more-label="loadMoreLabel"
    :loading-label="loadingLabel"
    :action-label="actionLabel"
    :item-key-mapper="(item) => item?.externalId || item?.external_id || item?.id"
    :item-title-mapper="(item) => item?.name || ''"
    :item-description-mapper="(item) => item?.description || ''"
    @open="$emit('open', $event)"
    @load-more="$emit('load-more')"
  />
</template>

<script setup>
import CardListSection from '@/components/list/CardListSection.vue';

defineProps({
  title: { type: String, default: '' },
  items: { type: Array, default: () => [] },
  emptyText: { type: String, default: '' },
  tagType: { type: String, default: '' },
  tagLabel: { type: String, default: '' },
  tagMapper: { type: Function, default: null },
  fallbackDescription: { type: String, default: '' },
  metaMapper: { type: Function, default: null },
  showLoadMore: { type: Boolean, default: false },
  loading: { type: Boolean, default: false },
  loadMoreLabel: { type: String, default: '' },
  loadingLabel: { type: String, default: '' },
  actionLabel: { type: String, default: '' }
});

defineEmits(['open', 'load-more']);
</script>
