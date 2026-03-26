<template>
  <section class="manage-section" :class="{ 'manage-section--plain': plain }">
    <div v-if="showHeader" class="section-header">
      <slot name="header">
        <h3 class="section-title">{{ title }}</h3>
      </slot>
    </div>
    <div class="section-body" :class="{ 'section-body--padded': bodyPadding === 'md' }">
      <slot />
    </div>
  </section>
</template>

<script setup>
import { computed, useSlots } from 'vue';

const props = defineProps({
  title: {
    type: String,
    default: ''
  },
  plain: {
    type: Boolean,
    default: false
  },
  bodyPadding: {
    type: String,
    default: 'none'
  }
});

const slots = useSlots();
const showHeader = computed(() => Boolean(props.title || slots.header));
</script>

<style scoped>
.manage-section {
  background: var(--bg-white);
  border: 1px solid var(--border-color);
  border-radius: var(--border-radius-md);
  box-shadow: var(--shadow-sm);
  overflow: hidden;
}

.manage-section--plain {
  background: transparent;
  border-color: transparent;
  box-shadow: none;
}

.section-header {
  padding: var(--spacing-md);
  border-bottom: 1px solid var(--border-color);
  background: var(--bg-white);
}

.section-title {
  margin: 0;
  font-size: 16px;
  font-weight: 600;
  color: var(--text-dark);
}

.section-body {
  padding: 0;
}

.section-body--padded {
  padding: var(--spacing-md);
}

.dark-mode .manage-section {
  background: #161b22;
  border-color: #2b2f36;
}

.dark-mode .manage-section--plain {
  background: transparent;
  border-color: transparent;
}

.dark-mode .section-header {
  background: #161b22;
  border-bottom-color: #2b2f36;
}
</style>
