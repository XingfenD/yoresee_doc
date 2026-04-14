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
        <AppAvatar
          :src="user.avatar"
          :name="user.nickname || user.username"
          :size="24"
        />
        <span class="mention-user-name">{{ user.nickname || user.username }}</span>
      </div>
    </div>
  </Teleport>
</template>

<script setup>
import { ref, watch } from 'vue';
import { listUsers } from '@/services/api/membership.js';
import AppAvatar from '@/components/base/AppAvatar.vue';

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
  background: var(--bg-white, #fff);
  border: 1px solid var(--border-color, #c9cdd4);
  border-radius: 6px;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.18);
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

.mention-user-name {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
</style>
