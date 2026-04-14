<template>
  <div ref="wrapperRef" class="mention-textarea-wrapper">
    <el-input
      ref="inputRef"
      :model-value="modelValue"
      type="textarea"
      :autosize="autosize"
      :placeholder="placeholder"
      @update:model-value="onModelValueUpdate"
      @keydown="onKeydown"
    />
    <!-- Mention dropdown anchored to wrapper -->
    <div
      v-if="dropdownVisible && users.length > 0"
      class="mention-dropdown"
      :style="dropdownStyle"
    >
      <div
        v-for="(user, index) in users"
        :key="user.external_id"
        class="mention-dropdown-item"
        :class="{ 'is-active': index === activeIndex }"
        @mousedown.prevent="selectUser(user)"
      >
        <img v-if="user.avatar" :src="user.avatar" class="mention-item-avatar" />
        <span v-else class="mention-item-avatar mention-item-avatar--fallback">
          {{ (user.nickname || user.username || '?').charAt(0).toUpperCase() }}
        </span>
        <span class="mention-item-name">{{ user.nickname || user.username }}</span>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, watch, nextTick } from 'vue';
import { listUsers } from '@/services/api/membership.js';

const props = defineProps({
  modelValue: { type: String, default: '' },
  placeholder: { type: String, default: '' },
  autosize: { type: [Boolean, Object], default: () => ({ minRows: 2, maxRows: 4 }) },
});

const emit = defineEmits(['update:modelValue', 'mention-users-change']);

const inputRef = ref(null);
const wrapperRef = ref(null);

const dropdownVisible = ref(false);
const dropdownStyle = ref({});
const users = ref([]);
const activeIndex = ref(0);
const keyword = ref('');
const mentionedUsers = ref([]);

let mentionStart = -1;
let fetchTimer = null;

function getTextareaEl() {
  return inputRef.value?.textarea ?? inputRef.value?.$el?.querySelector('textarea') ?? null;
}

function onModelValueUpdate(val) {
  emit('update:modelValue', val);
  detectMention(val);
}

function detectMention(value) {
  const el = getTextareaEl();
  if (!el) return;

  const cursor = el.selectionStart ?? value.length;

  // Scan backwards from cursor for @ on same line
  let atPos = -1;
  for (let i = cursor - 1; i >= 0; i--) {
    const ch = value[i];
    if (ch === '@') { atPos = i; break; }
    if (ch === '\n' || ch === ' ') break;
  }

  if (atPos === -1) {
    closeDrop();
    return;
  }

  mentionStart = atPos;
  keyword.value = value.slice(atPos + 1, cursor);
  activeIndex.value = 0;

  // Position dropdown below the @ character in the textarea
  positionDropdown(el, atPos);
  dropdownVisible.value = true;

  // Debounce the API fetch
  clearTimeout(fetchTimer);
  fetchTimer = setTimeout(() => fetchUsers(keyword.value), 150);
}

async function fetchUsers(kw) {
  try {
    const result = await listUsers({ keyword: kw, page: 1, page_size: 8 });
    users.value = result.users || [];
    if (users.value.length === 0) closeDrop();
  } catch {
    users.value = [];
    closeDrop();
  }
}

function positionDropdown(el, atCharPos) {
  // Use a mirror div to find approximate caret Y position
  const style = window.getComputedStyle(el);
  const mirror = document.createElement('div');
  const props_to_copy = [
    'fontSize', 'fontFamily', 'fontWeight', 'letterSpacing',
    'lineHeight', 'paddingTop', 'paddingBottom', 'paddingLeft', 'paddingRight',
    'borderTopWidth', 'borderLeftWidth', 'boxSizing', 'wordWrap',
  ];
  for (const p of props_to_copy) mirror.style[p] = style[p];
  mirror.style.position = 'absolute';
  mirror.style.visibility = 'hidden';
  mirror.style.top = '-9999px';
  mirror.style.left = '-9999px';
  mirror.style.width = el.offsetWidth + 'px';
  mirror.style.whiteSpace = 'pre-wrap';
  mirror.style.overflowWrap = 'break-word';

  const before = el.value.slice(0, atCharPos);
  mirror.textContent = before;
  const caret = document.createElement('span');
  caret.textContent = '@';
  mirror.appendChild(caret);
  document.body.appendChild(mirror);

  const elRect = el.getBoundingClientRect();
  const caretTop = caret.offsetTop - el.scrollTop;
  const caretLeft = caret.offsetLeft;
  const lineH = parseInt(style.lineHeight) || parseInt(style.fontSize) * 1.4 || 20;
  document.body.removeChild(mirror);

  // Position relative to wrapper
  const wRect = wrapperRef.value?.getBoundingClientRect() ?? elRect;
  const relTop = elRect.top - wRect.top + caretTop + lineH + 2;
  const relLeft = elRect.left - wRect.left + caretLeft;

  dropdownStyle.value = {
    top: relTop + 'px',
    left: Math.max(0, relLeft) + 'px',
  };
}

function onKeydown(event) {
  if (!dropdownVisible.value || users.value.length === 0) return;

  if (event.key === 'ArrowDown') {
    event.preventDefault();
    activeIndex.value = (activeIndex.value + 1) % users.value.length;
  } else if (event.key === 'ArrowUp') {
    event.preventDefault();
    activeIndex.value = (activeIndex.value - 1 + users.value.length) % users.value.length;
  } else if (event.key === 'Enter') {
    const user = users.value[activeIndex.value];
    if (user) {
      event.preventDefault();
      selectUser(user);
    }
  } else if (event.key === 'Escape') {
    closeDrop();
  }
}

function selectUser(user) {
  const el = getTextareaEl();
  if (!el || mentionStart === -1) { closeDrop(); return; }

  const cursor = el.selectionStart ?? el.value.length;
  const displayName = user.nickname || user.username;
  const before = el.value.slice(0, mentionStart);
  const after = el.value.slice(cursor);
  const inserted = `@${displayName} `;
  const newValue = before + inserted + after;

  emit('update:modelValue', newValue);
  nextTick(() => {
    const newCursor = mentionStart + inserted.length;
    el.value = newValue;
    el.setSelectionRange(newCursor, newCursor);
    el.focus();
  });

  if (!mentionedUsers.value.find(u => u.external_id === user.external_id)) {
    mentionedUsers.value.push({ external_id: user.external_id, nickname: displayName });
    emit('mention-users-change', mentionedUsers.value.map(u => u.external_id));
  }

  mentionStart = -1;
  closeDrop();
}

function closeDrop() {
  dropdownVisible.value = false;
  users.value = [];
  keyword.value = '';
}

function getMentionedUsers() {
  return mentionedUsers.value;
}

function reset() {
  mentionedUsers.value = [];
  closeDrop();
}

defineExpose({ getMentionedUsers, reset });
</script>

<style scoped>
.mention-textarea-wrapper {
  position: relative;
  width: 100%;
}

.mention-dropdown {
  position: absolute;
  z-index: 9999;
  background: var(--bg-white, #fff);
  border: 1px solid var(--border-color, #c9cdd4);
  border-radius: 6px;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.12);
  min-width: 180px;
  max-width: 260px;
  overflow: hidden;
}

.mention-dropdown-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 6px 12px;
  cursor: pointer;
  font-size: 14px;
  color: var(--text-primary, #1d2129);
}

.mention-dropdown-item:hover,
.mention-dropdown-item.is-active {
  background: var(--primary-light, #e8f0ff);
}

.mention-item-avatar {
  width: 24px;
  height: 24px;
  border-radius: 50%;
  object-fit: cover;
  flex-shrink: 0;
}

.mention-item-avatar--fallback {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  background: var(--primary-color, #165dff);
  color: #fff;
  font-size: 12px;
  font-weight: 600;
}

.mention-item-name {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
</style>
