<template>
  <aside class="kb-templates-panel">
    <div class="kb-templates-inner" v-loading="loading">
      <TemplateListSection
        :title="title"
        :items="items"
        :empty-text="emptyText"
        :fallback-description="fallbackDescription"
        :tag-mapper="tagMapper"
        :meta-mapper="metaMapper"
        :action-label="actionLabel"
        @open="(item) => emit('open', item)"
      />
    </div>
  </aside>
</template>

<script setup>
import TemplateListSection from '@/components/template/TemplateListSection.vue';

defineProps({
  loading: { type: Boolean, default: false },
  title: { type: String, default: '' },
  items: { type: Array, default: () => [] },
  emptyText: { type: String, default: '' },
  fallbackDescription: { type: String, default: '' },
  tagMapper: { type: Function, default: null },
  metaMapper: { type: Function, default: null },
  actionLabel: { type: String, default: '' }
});

const emit = defineEmits(['open']);
</script>

<style scoped>
.kb-templates-panel {
  display: flex;
  min-height: 0;
  width: 320px;
  flex-shrink: 0;
}

.kb-templates-inner {
  flex: 1;
  min-height: 0;
}

.kb-templates-inner :deep(.card-list-section) {
  height: 100%;
  display: flex;
  flex-direction: column;
}

.kb-templates-inner :deep(.section-content) {
  flex: 1;
  min-height: 0;
  overflow-y: auto;
}

@media (max-width: 1200px) {
  .kb-templates-panel {
    width: 100%;
  }
}
</style>
