<template>
  <img
    v-if="src"
    :src="src"
    class="app-avatar"
    :style="sizeStyle"
    :alt="name"
    @error="onImgError"
  />
  <span
    v-else
    class="app-avatar app-avatar--fallback"
    :style="sizeStyle"
    :aria-label="name"
  >
    {{ initial }}
  </span>
</template>

<script setup>
import { ref, computed, watch } from 'vue';

const props = defineProps({
  src: {
    type: String,
    default: ''
  },
  name: {
    type: String,
    default: ''
  },
  size: {
    type: Number,
    default: 28
  }
});

// If the image errors out, fall back to the badge
const imgFailed = ref(false);
watch(() => props.src, () => { imgFailed.value = false; });

const onImgError = () => { imgFailed.value = true; };

const src = computed(() => (!imgFailed.value && props.src) ? props.src : '');

const initial = computed(() => {
  const n = (props.name || '').trim();
  return n ? n.charAt(0).toUpperCase() : '?';
});

const sizeStyle = computed(() => ({
  width: `${props.size}px`,
  height: `${props.size}px`,
  fontSize: `${Math.round(props.size * 0.43)}px`
}));
</script>

<style scoped>
.app-avatar {
  border-radius: 50%;
  object-fit: cover;
  flex-shrink: 0;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  vertical-align: middle;
}

.app-avatar--fallback {
  background: var(--primary-color, #165dff);
  color: #fff;
  font-weight: 600;
  user-select: none;
}
</style>
