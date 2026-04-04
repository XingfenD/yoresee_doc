<template>
  <div class="slide-editor">
    <section class="slide-panel slide-panel--source">
      <header class="slide-panel-header">
        <span>{{ t('document.slide.sourceTitle') }}</span>
        <span class="slide-panel-subtitle">{{ t('document.slide.sourceTip') }}</span>
      </header>
      <el-input
        :model-value="draft"
        type="textarea"
        class="slide-source-input"
        resize="none"
        :autosize="false"
        :placeholder="placeholder || t('document.slide.placeholder')"
        @update:model-value="handleInput"
        @blur="emit('commit')"
      />
    </section>

    <section class="slide-panel slide-panel--preview">
      <header class="slide-panel-header">
        <span>{{ t('document.slide.previewTitle') }}</span>
        <span class="slide-panel-subtitle">{{ t('document.slide.count', { count: renderedSlides.length }) }}</span>
      </header>
      <div class="slide-preview-body">
        <swiper
          ref="swiperRef"
          class="slide-swiper"
          :modules="swiperModules"
          :slides-per-view="1"
          :space-between="12"
          :navigation="true"
          :pagination="{ clickable: true }"
          :keyboard="{ enabled: true }"
          :observer="true"
          :observe-parents="true"
          @slideChange="handleSlideChange"
        >
          <swiper-slide
            v-for="slide in renderedSlides"
            :key="slide.id"
            class="slide-swiper-item"
          >
            <article class="slide-card" v-html="slide.html"></article>
          </swiper-slide>
        </swiper>
      </div>
      <div v-if="activeIndex >= 0" class="slide-index">
        {{ `${activeIndex + 1} / ${renderedSlides.length}` }}
      </div>
      <div class="slide-toolbar-hint">
        {{ t('document.slide.previewHint') }}
      </div>
    </section>
  </div>
</template>

<script setup>
import { computed, nextTick, ref, watch } from 'vue';
import { useI18n } from 'vue-i18n';
import { marked } from 'marked';
import { Swiper, SwiperSlide } from 'swiper/vue';
import { Keyboard, Navigation, Pagination } from 'swiper/modules';
import 'swiper/css';
import 'swiper/css/navigation';
import 'swiper/css/pagination';

const props = defineProps({
  modelValue: {
    type: String,
    default: ''
  },
  placeholder: {
    type: String,
    default: ''
  }
});

const emit = defineEmits(['update:modelValue', 'commit']);
const { t } = useI18n();

const swiperRef = ref(null);
const draft = ref(props.modelValue || '');
const activeIndex = ref(0);
const swiperModules = [Navigation, Pagination, Keyboard];

marked.setOptions({
  gfm: true,
  breaks: true
});

const splitSlides = (content) => {
  const normalized = String(content || '').replace(/\r\n/g, '\n');
  const chunks = normalized
    .split(/\n-{3,}\n/g)
    .map((item) => item.trim())
    .filter(Boolean);

  if (chunks.length > 0) {
    return chunks;
  }
  return [t('document.slide.emptySlide')];
};

const renderedSlides = computed(() =>
  splitSlides(draft.value).map((slide, index) => ({
    id: `${index}-${slide.length}`,
    html: marked.parse(slide)
  }))
);

const updateSwiperLayout = async () => {
  await nextTick();
  const instance = swiperRef.value?.swiper;
  if (!instance) {
    return;
  }
  instance.update?.();
  if (instance.activeIndex >= renderedSlides.value.length) {
    const targetIndex = Math.max(0, renderedSlides.value.length - 1);
    instance.slideTo(targetIndex, 0);
  }
  activeIndex.value = instance.activeIndex || 0;
};

const handleInput = (value) => {
  draft.value = value;
  emit('update:modelValue', value);
  updateSwiperLayout();
};

const handleSlideChange = (swiper) => {
  if (!swiper) {
    return;
  }
  activeIndex.value = Number(swiper.activeIndex) || 0;
};

watch(
  () => props.modelValue,
  (value) => {
    if (value === draft.value) {
      return;
    }
    draft.value = value || '';
    updateSwiperLayout();
  }
);

watch(
  renderedSlides,
  () => {
    updateSwiperLayout();
  },
  { deep: true }
);

const reRender = () => {
  updateSwiperLayout();
};

defineExpose({
  reRender
});
</script>

<style scoped>
.slide-editor {
  display: flex;
  gap: var(--spacing-sm);
  width: 100%;
  height: 100%;
  min-width: 0;
  min-height: 0;
  padding: var(--spacing-sm);
  box-sizing: border-box;
}

.slide-panel {
  display: flex;
  flex-direction: column;
  min-width: 0;
  min-height: 0;
  border: 1px solid var(--border-color);
  border-radius: var(--border-radius-md);
  overflow: hidden;
  background: var(--bg-white);
}

.slide-panel--source {
  width: min(40%, 460px);
  flex-shrink: 0;
}

.slide-panel--preview {
  flex: 1;
}

.slide-panel-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: var(--spacing-sm);
  padding: 10px 12px;
  border-bottom: 1px solid var(--border-color);
  font-size: 13px;
  font-weight: 600;
  color: var(--text-dark);
}

.slide-panel-subtitle {
  font-weight: 400;
  color: var(--text-light);
  font-size: 12px;
}

.slide-source-input {
  flex: 1;
  min-height: 0;
}

.slide-source-input :deep(.el-textarea__inner) {
  height: 100%;
  min-height: 100%;
  resize: none;
  border: none;
  box-shadow: none;
  border-radius: 0;
  padding: 12px;
  font-size: 14px;
  line-height: 1.6;
  background: transparent;
}

.slide-preview-body {
  flex: 1;
  min-height: 0;
  width: 100%;
  padding: 12px;
  background: color-mix(in srgb, var(--bg-light) 75%, transparent);
}

.slide-swiper {
  width: 100%;
  height: 100%;
}

.slide-swiper-item {
  height: 100%;
}

.slide-card {
  height: 100%;
  min-height: 220px;
  overflow: auto;
  border: 1px solid var(--border-color);
  border-radius: var(--border-radius-md);
  background: var(--bg-white);
  padding: 24px;
  color: var(--text-dark);
}

.slide-card :deep(h1),
.slide-card :deep(h2),
.slide-card :deep(h3) {
  margin: 0 0 12px;
  line-height: 1.25;
}

.slide-card :deep(p),
.slide-card :deep(li) {
  margin: 0 0 8px;
  line-height: 1.6;
  color: var(--text-medium);
}

.slide-index {
  padding: 6px 12px 0;
  font-size: 12px;
  color: var(--text-light);
  text-align: right;
}

.slide-toolbar-hint {
  padding: 4px 12px 10px;
  font-size: 12px;
  color: var(--text-light);
}

:global(.dark-mode) .slide-panel {
  background: var(--bg-white);
  border-color: var(--border-color);
}

:global(.dark-mode) .slide-panel-header {
  border-color: var(--border-color);
  color: var(--text-dark);
}

:global(.dark-mode) .slide-source-input :deep(.el-textarea__inner) {
  color: var(--text-dark);
}

:global(.dark-mode) .slide-card {
  background: color-mix(in srgb, var(--bg-white) 95%, #05070a 5%);
  border-color: var(--border-color);
}

:global(.dark-mode) .slide-card :deep(p),
:global(.dark-mode) .slide-card :deep(li) {
  color: var(--text-medium);
}

.slide-panel--preview {
  position: relative;
}

.slide-swiper :deep(.swiper-button-next),
.slide-swiper :deep(.swiper-button-prev) {
  color: var(--text-light);
}

.slide-swiper :deep(.swiper-pagination-bullet) {
  background: var(--text-light);
  opacity: 0.45;
}

.slide-swiper :deep(.swiper-pagination-bullet-active) {
  background: var(--primary-color);
  opacity: 1;
}
</style>
