<template>
  <span
    ref="triggerRef"
    class="document-type-menu-trigger"
    :class="{ 'is-disabled': disabled }"
    @click.capture="handleTriggerClick"
  >
    <slot />
  </span>

  <AppMenu
    ref="menuRef"
    :visible="visible"
    :x="position.x"
    :y="position.y"
    :ignore-elements="[triggerRef]"
    @close="closeMenu"
  >
    <AppMenuItem
      v-for="option in resolvedOptions"
      :key="option.value"
      :disabled="Boolean(option.disabled)"
      @click="handleSelect(option)"
    >
      {{ option.label }}
    </AppMenuItem>
  </AppMenu>
</template>

<script setup>
import { computed, nextTick, ref } from 'vue';
import { useI18n } from 'vue-i18n';
import AppMenu from '@/components/base/AppMenu.vue';
import AppMenuItem from '@/components/base/AppMenuItem.vue';
import { buildDocumentTypeOptions, normalizeDocumentType } from '@/utils/documentType';

const props = defineProps({
  trigger: {
    type: String,
    default: 'click'
  },
  placement: {
    type: String,
    default: 'bottom-start'
  },
  disabled: {
    type: Boolean,
    default: false
  },
  options: {
    type: Array,
    default: () => []
  }
});

const emit = defineEmits(['select', 'open', 'close', 'visible-change']);
const { t } = useI18n();

const visible = ref(false);
const triggerRef = ref(null);
const menuRef = ref(null);
const position = ref({ x: 0, y: 0 });

const defaultOptions = computed(() => buildDocumentTypeOptions(t));
const resolvedOptions = computed(() => (props.options.length ? props.options : defaultOptions.value));
const updatePosition = async () => {
  const triggerEl = triggerRef.value;
  if (!triggerEl) {
    return;
  }

  const rect = triggerEl.getBoundingClientRect();
  const gap = 6;
  let x = rect.left;
  let y = rect.bottom + gap;

  if (props.placement === 'right-start') {
    x = rect.right + gap;
    y = rect.top;
  }

  position.value = { x, y };
  await nextTick();

  const menuEl = menuRef.value?.getMenuElement?.();
  if (!menuEl) {
    return;
  }

  const menuWidth = menuEl.offsetWidth || 150;
  const menuHeight = menuEl.offsetHeight || 120;
  const maxX = Math.max(8, window.innerWidth - menuWidth - 8);
  const maxY = Math.max(8, window.innerHeight - menuHeight - 8);

  position.value = {
    x: Math.min(Math.max(8, x), maxX),
    y: Math.min(Math.max(8, y), maxY)
  };
};

const setVisible = (nextVisible) => {
  if (visible.value === nextVisible) {
    return;
  }
  visible.value = nextVisible;
  emit('visible-change', nextVisible);
  emit(nextVisible ? 'open' : 'close');
};

const closeMenu = () => {
  setVisible(false);
};

const handleTriggerClick = async (event) => {
  if (props.disabled) {
    return;
  }
  event.stopPropagation();
  const nextVisible = !visible.value;
  setVisible(nextVisible);
  if (nextVisible) {
    await updatePosition();
  }
};

const handleSelect = (option) => {
  if (!option || option.disabled) {
    return;
  }
  emit('select', normalizeDocumentType(option.value));
  closeMenu();
};
</script>

<style scoped>
.document-type-menu-trigger {
  display: inline-flex;
  align-items: center;
}

.document-type-menu-trigger.is-disabled {
  cursor: not-allowed;
}
</style>
