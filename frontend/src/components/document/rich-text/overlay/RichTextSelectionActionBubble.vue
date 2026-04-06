<template>
  <div
    v-show="visible"
    class="selection-action-bubble"
    :style="styleObject"
    @mouseenter="$emit('enter')"
    @mouseleave="$emit('leave')"
  >
    <button
      type="button"
      class="selection-action-btn"
      :class="{ 'is-active': selectionBoldActive }"
      title="Bold"
      @mousedown.prevent
      @click="$emit('toggle-bold')"
    >
      B
    </button>
    <button
      type="button"
      class="selection-action-btn"
      :class="{ 'is-active': selectionUnderlineActive }"
      title="Underline"
      @mousedown.prevent
      @click="$emit('toggle-underline')"
    >
      U
    </button>
    <button
      type="button"
      class="selection-action-btn"
      :class="{ 'is-active': selectionItalicActive }"
      title="Italic"
      @mousedown.prevent
      @click="$emit('toggle-italic')"
    >
      I
    </button>
    <button
      type="button"
      class="selection-action-btn"
      :class="{ 'is-active': selectionStrikeActive }"
      title="Strike"
      @mousedown.prevent
      @click="$emit('toggle-strike')"
    >
      S
    </button>
    <button
      type="button"
      class="selection-action-btn"
      :class="{ 'is-active': selectionInlineCodeActive }"
      title="Inline Code"
      @mousedown.prevent
      @click="$emit('toggle-inline-code')"
    >
      &lt;/&gt;
    </button>
    <button
      type="button"
      class="selection-action-btn selection-action-btn--icon"
      :class="{ 'is-active': selectionLinkActive }"
      title="Link"
      @mousedown.prevent
      @click="$emit('link-click')"
    >
      <el-icon>
        <Link />
      </el-icon>
    </button>

    <div class="selection-color-group">
      <button
        ref="colorTriggerRef"
        type="button"
        class="selection-color-trigger"
        :class="{ 'is-active': selectionHighlightActive }"
        title="字体与背景颜色"
        @mousedown.prevent
        @click="toggleColorPanel"
      >
        <span class="selection-color-trigger-a" :style="colorTriggerAStyle">A</span>
        <el-icon class="selection-color-trigger-arrow">
          <component :is="colorPanelVisible ? ArrowUp : ArrowDown" />
        </el-icon>
      </button>

      <div
        v-if="colorPanelVisible"
        ref="colorPanelRef"
        class="selection-color-panel"
        @mousedown.stop
      >
        <div class="selection-color-panel-title">字体颜色</div>
        <div class="selection-color-text-row">
          <button
            v-for="item in textColorItems"
            :key="item"
            type="button"
            class="selection-color-text-item"
            :class="{ 'is-active': isTextColorActive(item) }"
            :style="{ color: item }"
            @mousedown.prevent
            @click="pickTextColor(item)"
          >
            A
          </button>
        </div>

        <div class="selection-color-panel-title">背景颜色</div>
        <div class="selection-color-bg-grid">
          <button
            v-for="item in backgroundColorItems"
            :key="item || 'none'"
            type="button"
            class="selection-color-bg-item"
            :class="{
              'is-active': isBackgroundColorActive(item),
              'is-none': !item
            }"
            :style="item ? { backgroundColor: item } : {}"
            @mousedown.prevent
            @click="pickBackgroundColor(item)"
          ></button>
        </div>

        <button
          type="button"
          class="selection-color-reset-btn"
          @mousedown.prevent
          @click="resetColorToDefault"
        >
          恢复默认
        </button>
      </div>
    </div>

    <button
      v-if="commentEnabled"
      type="button"
      class="selection-action-btn selection-action-btn--icon"
      :title="commentTitle"
      @mousedown.prevent
      @click="$emit('comment-click')"
    >
      <el-icon>
        <ChatDotRound />
      </el-icon>
    </button>
  </div>
</template>

<script setup>
import { computed } from 'vue';
import { ArrowDown, ArrowUp, ChatDotRound, Link } from '@element-plus/icons-vue';
import { useRichTextSelectionColorPanel } from './useRichTextSelectionColorPanel';

const props = defineProps({
  visible: {
    type: Boolean,
    default: false
  },
  styleObject: {
    type: Object,
    default: () => ({})
  },
  commentEnabled: {
    type: Boolean,
    default: false
  },
  commentTitle: {
    type: String,
    default: ''
  },
  selectionBoldActive: {
    type: Boolean,
    default: false
  },
  selectionUnderlineActive: {
    type: Boolean,
    default: false
  },
  selectionItalicActive: {
    type: Boolean,
    default: false
  },
  selectionStrikeActive: {
    type: Boolean,
    default: false
  },
  selectionInlineCodeActive: {
    type: Boolean,
    default: false
  },
  selectionLinkActive: {
    type: Boolean,
    default: false
  },
  selectionHighlightActive: {
    type: Boolean,
    default: false
  },
  selectionTextColor: {
    type: String,
    default: '#1f2937'
  },
  selectionBackgroundColor: {
    type: String,
    default: ''
  }
});

const emit = defineEmits([
  'enter',
  'leave',
  'toggle-bold',
  'toggle-underline',
  'toggle-italic',
  'toggle-strike',
  'toggle-inline-code',
  'link-click',
  'text-color-change',
  'background-color-change',
  'comment-click'
]);

const visibleRef = computed(() => props.visible);
const selectionTextColorRef = computed(() => props.selectionTextColor);
const selectionBackgroundColorRef = computed(() => props.selectionBackgroundColor);

const {
  textColorItems,
  backgroundColorItems,
  colorPanelVisible,
  colorPanelRef,
  colorTriggerRef,
  colorTriggerAStyle,
  pickTextColor,
  pickBackgroundColor,
  isTextColorActive,
  isBackgroundColorActive,
  resetColorToDefault,
  toggleColorPanel
} = useRichTextSelectionColorPanel({
  visibleRef,
  selectionTextColorRef,
  selectionBackgroundColorRef,
  onTextColorChange: (value) => emit('text-color-change', value),
  onBackgroundColorChange: (value) => emit('background-color-change', value)
});
</script>
<style scoped src="./RichTextSelectionActionBubble.css"></style>
