<template>
  <Teleport to="body">
    <div
      v-if="visible && users.length > 0"
      class="mention-user-list"
      :style="{ top: position.y + 'px', left: position.x + 'px' }"
    >
      <div
        v-for="(user, index) in users"
        :key="user.external_id"
        class="mention-user-item"
        :class="{ 'is-active': index === activeIndex }"
        @mousedown.prevent="emit('select', user)"
      >
        <img v-if="user.avatar" :src="user.avatar" class="mention-user-avatar" />
        <span v-else class="mention-user-avatar mention-user-avatar--fallback">
          {{ (user.nickname || user.username || '?').charAt(0).toUpperCase() }}
        </span>
        <span class="mention-user-name">{{ user.nickname || user.username }}</span>
      </div>
    </div>
  </Teleport>
</template>

<script setup>
import { ref, watch } from 'vue';
import { listUsers } from '@/services/api/membership.js';

const props = defineProps({
  keyword: { type: String, default: '' },
  position: { type: Object, default: () => ({ x: 0, y: 0 }) },
  visible: { type: Boolean, default: false },
  activeIndex: { type: Number, default: 0 },
});

const emit = defineEmits(['select']);

const users = ref([]);

watch(
  () => [props.visible, props.keyword],
  async ([visible, keyword]) => {
    if (!visible) {
      users.value = [];
      return;
    }
    try {
      const result = await listUsers({ keyword, page: 1, page_size: 8 });
      users.value = result.users || [];
    } catch {
      users.value = [];
    }
  },
  { immediate: true },
);
</script>

<style scoped>
.mention-user-list {
  position: fixed;
  z-index: 9999;
  background: #fff;
  border: 1px solid var(--border-color, #c9cdd4);
  border-radius: 6px;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.12);
  min-width: 180px;
  max-width: 260px;
  overflow: hidden;
}

.mention-user-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 6px 12px;
  cursor: pointer;
  font-size: 14px;
  color: var(--text-primary, #1d2129);
}

.mention-user-item:hover,
.mention-user-item.is-active {
  background: var(--primary-light, #e8f0ff);
}

.mention-user-avatar {
  width: 24px;
  height: 24px;
  border-radius: 50%;
  object-fit: cover;
  flex-shrink: 0;
}

.mention-user-avatar--fallback {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  background: var(--primary-color, #165dff);
  color: #fff;
  font-size: 12px;
  font-weight: 600;
}

.mention-user-name {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
</style>
