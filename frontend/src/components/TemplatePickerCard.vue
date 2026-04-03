<template>
  <div
    class="template-card"
    :class="{
      'is-selected': isSelected,
      'is-blank': Boolean(item?.is_blank),
      'is-grid': layout === 'grid'
    }"
    @click="emit('select', item)"
  >
    <div v-if="layout === 'grid'" class="template-card-preview">
      {{ previewText }}
    </div>
    <div class="template-card-title">
      {{ item?.name }}
    </div>
    <div class="template-card-desc">
      {{ item?.description || fallbackDescription }}
    </div>

    <div v-if="showPreviewButton" class="template-card-actions">
      <el-button class="template-preview-btn" size="small" @click.stop="emit('preview', item)">
        {{ t('common.preview') }}
      </el-button>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue';
import { useI18n } from 'vue-i18n';

const props = defineProps({
  item: {
    type: Object,
    default: () => ({})
  },
  selectedTemplateId: {
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

const emit = defineEmits(['select', 'preview']);
const { t } = useI18n();

const isSelected = computed(() => props.selectedTemplateId === String(props.item?.id || ''));

const previewText = computed(() => {
  if (props.item?.is_blank) {
    return '+';
  }
  const text = String(props.item?.name || '').trim();
  return text ? text.slice(0, 1).toUpperCase() : 'T';
});

const showPreviewButton = computed(() => !props.item?.is_blank);
</script>

<style scoped>
.template-card {
  cursor: pointer;
  border-bottom: 1px solid #eef0f4;
  padding: 12px 8px 12px 12px;
  transition: background-color 0.2s ease, border-color 0.2s ease;
  position: relative;
}

.template-card:hover {
  background-color: #f7f8fa;
}

.template-card.is-selected {
  background-color: #f0f5ff;
}

.template-card-preview {
  height: 72px;
  border-radius: 10px;
  margin-bottom: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 22px;
  font-weight: 700;
  color: #2f65e2;
  background: linear-gradient(135deg, #e9f0ff 0%, #f7faff 100%);
}

.template-card.is-blank .template-card-preview {
  border: 1px dashed #91b0ff;
  background: linear-gradient(135deg, #edf3ff 0%, #f8fbff 100%);
}

.template-card-title {
  font-weight: 600;
  color: #111827;
  display: flex;
  align-items: center;
  gap: 8px;
}

.template-card-desc {
  margin-top: 4px;
  color: #6b7280;
  font-size: 12px;
}

.template-card::before {
  content: '';
  position: absolute;
  left: 0;
  top: 8px;
  bottom: 8px;
  width: 3px;
  border-radius: 2px;
  background: #3370ff;
  opacity: 0;
}

.template-card.is-selected::before {
  opacity: 1;
}

.template-card.is-grid {
  border: 1px solid #e7ebf2;
  border-radius: 12px;
  min-height: 158px;
  background: #fff;
  padding: 12px;
}

.template-card.is-grid:hover {
  background-color: #f8faff;
  border-color: #d7e2ff;
}

.template-card.is-grid.is-selected {
  border-color: #6b93ff;
  background: #f2f7ff;
}

.template-card.is-grid::before {
  left: auto;
  top: auto;
  bottom: 12px;
  right: 12px;
  width: 8px;
  height: 8px;
  border-radius: 999px;
}

.template-card-actions {
  position: absolute;
  right: 8px;
  top: 8px;
  opacity: 0;
  transform: translateY(-2px);
  transition: opacity 0.16s ease, transform 0.16s ease;
}

.template-card:hover .template-card-actions {
  opacity: 1;
  transform: translateY(0);
}

.template-preview-btn {
  border-radius: 8px;
  border-color: color-mix(in srgb, #3370ff 35%, white);
  background: color-mix(in srgb, #3370ff 8%, white);
  color: #2f65e2;
}

.template-preview-btn:hover,
.template-preview-btn:focus {
  border-color: #3370ff;
  color: #2456d0;
}

.dark-mode .template-card {
  background-color: transparent;
  border-color: #1f2937;
}

.dark-mode .template-card:hover {
  background-color: #111827;
}

.dark-mode .template-card.is-selected {
  background-color: #0b1f3a;
}

.dark-mode .template-card-title {
  color: #e5e7eb;
}

.dark-mode .template-card-desc {
  color: #9ca3af;
}

.dark-mode .template-card-preview {
  background: linear-gradient(135deg, #152235 0%, #1f314c 100%);
  color: #8db0ff;
}

.dark-mode .template-card.is-blank .template-card-preview {
  border-color: #4c8dff;
  background: linear-gradient(135deg, #102038 0%, #173054 100%);
}

.dark-mode .template-card.is-grid {
  border-color: #243040;
  background: #0f1218;
}

.dark-mode .template-card.is-grid:hover {
  border-color: #35517e;
  background: #121a26;
}

.dark-mode .template-card.is-grid.is-selected {
  border-color: #4c8dff;
  background: #142237;
}

.dark-mode .template-preview-btn {
  border-color: rgba(76, 141, 255, 0.6);
  background: rgba(76, 141, 255, 0.16);
  color: #9ec0ff;
}

.dark-mode .template-preview-btn:hover,
.dark-mode .template-preview-btn:focus {
  border-color: rgba(118, 169, 255, 0.92);
  color: #dbe9ff;
}
</style>
