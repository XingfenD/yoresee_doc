<template>
  <div class="slide-preview-viewer">
    <el-empty v-if="renderedSlides.length === 0" />
    <template v-else>
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
          class="slide-item"
        >
          <article class="slide-card" v-html="slide.html"></article>
        </swiper-slide>
      </swiper>
      <div class="slide-index">{{ `${activeIndex + 1} / ${renderedSlides.length}` }}</div>
    </template>
  </div>
</template>

<script setup>
import { computed, nextTick, ref, watch } from 'vue';
import { marked } from 'marked';
import { Swiper, SwiperSlide } from 'swiper/vue';
import { Keyboard, Navigation, Pagination } from 'swiper/modules';
import 'swiper/css';
import 'swiper/css/navigation';
import 'swiper/css/pagination';

const props = defineProps({
  slides: {
    type: Array,
    default: () => []
  }
});

const swiperRef = ref(null);
const activeIndex = ref(0);
const swiperModules = [Navigation, Pagination, Keyboard];

marked.setOptions({
  gfm: true,
  breaks: true
});

const renderedSlides = computed(() =>
  (Array.isArray(props.slides) ? props.slides : [])
    .map((slide) => String(slide || '').trim())
    .filter(Boolean)
    .map((slide, index) => ({
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
    const target = Math.max(0, renderedSlides.value.length - 1);
    instance.slideTo(target, 0);
  }
  activeIndex.value = instance.activeIndex || 0;
};

const handleSlideChange = (swiper) => {
  if (!swiper) {
    return;
  }
  activeIndex.value = Number(swiper.activeIndex) || 0;
};

watch(
  renderedSlides,
  () => {
    updateSwiperLayout();
  },
  { deep: true }
);
</script>

<style scoped>
.slide-preview-viewer {
  width: 100%;
  height: 100%;
  min-height: 0;
  border: 1px solid var(--border-color);
  border-radius: var(--border-radius-sm);
  background: var(--bg-white);
  padding: 12px;
  box-sizing: border-box;
  display: flex;
  flex-direction: column;
}

.slide-swiper {
  width: 100%;
  flex: 1;
  min-height: 0;
}

.slide-item {
  height: 100%;
}

.slide-card {
  height: 100%;
  min-height: 220px;
  overflow: auto;
  border: 1px solid var(--border-color);
  border-radius: var(--border-radius-md);
  padding: 20px;
  background: var(--bg-white);
  color: var(--text-dark);
}

.slide-card :deep(h1),
.slide-card :deep(h2),
.slide-card :deep(h3) {
  margin-top: 0;
}

.slide-card :deep(p),
.slide-card :deep(li) {
  color: var(--text-medium);
  line-height: 1.6;
}

.slide-index {
  margin-top: 8px;
  text-align: right;
  color: var(--text-light);
  font-size: 12px;
}

.slide-swiper :deep(.swiper-button-next),
.slide-swiper :deep(.swiper-button-prev) {
  color: var(--text-light);
}

.slide-swiper :deep(.swiper-pagination-bullet) {
  background: var(--text-light);
  opacity: 0.35;
}

.slide-swiper :deep(.swiper-pagination-bullet-active) {
  opacity: 0.95;
}
</style>

