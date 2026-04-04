<template>
  <div class="template-pane" v-loading="loading">
    <div v-if="items.length === 0" class="template-empty">
      <el-empty :description="emptyText" />
    </div>
    <div v-else class="template-list" :class="`template-list--${layout}`">
      <TemplatePickerCard
        v-for="item in items"
        :key="item.id"
        :item="item"
        :layout="layout"
        :selected-template-id="selectedTemplateId"
        :fallback-description="fallbackDescription"
        @select="emit('select', $event)"
        @preview="openPreviewDialog"
      />
    </div>
  </div>

  <TemplatePreviewDialog
    v-model="showPreviewDialog"
    :title="previewDialogTitle"
    :content="previewContent"
    :document-type="previewingTemplate?.type || '1'"
    :is-dark-mode="isDarkMode"
    @closed="closePreviewDialog"
  />
</template>

<script setup>
import { computed, ref } from 'vue';
import { useI18n } from 'vue-i18n';
import { useUserStore } from '@/store/user';
import TemplatePickerCard from '@/components/template/TemplatePickerCard.vue';
import TemplatePreviewDialog from '@/components/template/TemplatePreviewDialog.vue';

defineProps({
  loading: {
    type: Boolean,
    default: false
  },
  items: {
    type: Array,
    default: () => []
  },
  selectedTemplateId: {
    type: String,
    default: ''
  },
  emptyText: {
    type: String,
    default: ''
  },
  fallbackDescription: {
    type: String,
    default: ''
  },
  layout: {
    type: String,
    default: 'list'
  }
});

const emit = defineEmits(['select']);
const { t } = useI18n();
const userStore = useUserStore();
const isDarkMode = computed(() => Boolean(userStore.darkMode));
const showPreviewDialog = ref(false);
const previewingTemplate = ref(null);

const previewContent = computed(() => String(previewingTemplate.value?.content || ''));
const previewDialogTitle = computed(() => {
  const name = String(previewingTemplate.value?.name || '').trim();
  if (!name) {
    return t('templates.previewTitle');
  }
  return `${t('templates.previewTitle')} · ${name}`;
});

const openPreviewDialog = (tpl) => {
  previewingTemplate.value = tpl || null;
  showPreviewDialog.value = true;
};

const closePreviewDialog = () => {
  showPreviewDialog.value = false;
  previewingTemplate.value = null;
};
</script>

<style scoped>
.template-pane {
  height: 100%;
}

.template-list {
  display: grid;
  grid-template-columns: 1fr;
  gap: 0;
}

.template-list--grid {
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 12px;
}

.template-empty {
  padding: var(--spacing-md) 0;
  min-height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
}

.template-empty :deep(.el-empty) {
  margin: 0;
}

@media (max-width: 640px) {
  .template-list {
    grid-template-columns: 1fr;
  }

  .template-list--grid {
    grid-template-columns: 1fr;
  }
}
</style>
