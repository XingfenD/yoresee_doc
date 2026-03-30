<template>
  <header class="top-nav">
    <div class="nav-left">
      <router-link class="system-link" to="/">
        <h1 class="system-title">{{ systemName }}</h1>
      </router-link>
    </div>
    <div class="nav-right">
      <el-dropdown trigger="click" @command="emit('change-language', $event)" class="nav-item">
        <span class="nav-link">
          <el-icon :size="18">
            <Flag v-if="currentLanguage === 'en'" />
            <ChatLineRound v-else />
          </el-icon>
        </span>
        <template #dropdown>
          <el-dropdown-menu>
            <el-dropdown-item command="en" :icon="'Flag'">
              {{ t('language.english') }}
            </el-dropdown-item>
            <el-dropdown-item command="zh" :icon="'ChatLineRound'">
              {{ t('language.chinese') }}
            </el-dropdown-item>
          </el-dropdown-menu>
        </template>
      </el-dropdown>

      <div class="nav-item theme-switch">
        <span class="nav-link" @click="emit('toggle-theme')">
          <el-icon :size="18">
            <Moon v-if="isDarkMode" />
            <Sunny v-else />
          </el-icon>
        </span>
      </div>

      <div class="nav-item">
        <span class="nav-link notification-link" @click="goToNotifications">
          <span v-if="hasUnread" class="notification-dot" />
          <el-icon :size="18">
            <Bell />
          </el-icon>
        </span>
      </div>

      <el-dropdown trigger="click" class="nav-item">
        <span class="user-info">
          <el-avatar v-if="userAvatar" size="small" :src="userAvatar"></el-avatar>
          <span class="username">{{ username }}</span>
          <el-icon class="el-icon--right">
            <ArrowDown />
          </el-icon>
        </span>
        <template #dropdown>
          <el-dropdown-menu>
            <el-dropdown-item v-if="showUserCenter" @click="goToUserCenter">{{ t('user.center') }}</el-dropdown-item>
            <el-dropdown-item v-if="showSystemManage" @click="handleSystemManage">
              {{ t('system.management') }}
            </el-dropdown-item>
            <el-dropdown-item divided @click="emit('logout')">{{ t('button.logout') }}</el-dropdown-item>
          </el-dropdown-menu>
        </template>
      </el-dropdown>
    </div>
  </header>
</template>

<script setup>
import { onMounted, onBeforeUnmount, ref } from 'vue';
import { useRouter } from 'vue-router';
import { useI18n } from 'vue-i18n';
import { ArrowDown, Flag, ChatLineRound, Moon, Sunny, Bell } from '@element-plus/icons-vue';
import { queryTopNavDisplay } from '@/services/auth';
import { listNotifications } from '@/services/api';
import { useApiAction } from '@/composables/useApiAction';

const props = defineProps({
  systemName: {
    type: String,
    default: ''
  },
  currentLanguage: {
    type: String,
    default: 'zh'
  },
  isDarkMode: {
    type: Boolean,
    default: false
  },
  userAvatar: {
    type: String,
    default: ''
  },
  username: {
    type: String,
    default: ''
  }
});

const emit = defineEmits(['change-language', 'toggle-theme', 'logout']);
const { t } = useI18n();
const router = useRouter();
const { runSilent } = useApiAction({ t });
const showSystemManage = ref(false);
const showUserCenter = ref(false);
const hasUnread = ref(false);

const goToUserCenter = () => {
  router.push('/user_info/example');
};

const handleSystemManage = () => {
  router.push('/manage');
};

const goToNotifications = () => {
  router.push('/user_info/notifications');
};

const loadTopNavDisplay = async () => {
  await runSilent(
    () => queryTopNavDisplay(),
    {
      context: 'loadTopNavDisplay',
      onSuccess: (resp) => {
        const menus = resp.display_menus || [];
        showUserCenter.value = menus.includes('user-center');
        showSystemManage.value = menus.includes('system-manage');
      },
      onError: () => {
        showUserCenter.value = false;
        showSystemManage.value = false;
      }
    }
  );
};

const loadUnreadNotifications = async () => {
  await runSilent(
    () => listNotifications({ page: 1, page_size: 1, status: 'unread' }),
    {
      context: 'loadUnreadNotifications',
      onSuccess: (resp) => {
        hasUnread.value = Number(resp.total) > 0;
      },
      onError: () => {
        hasUnread.value = false;
      }
    }
  );
};

const handleUnreadEvent = (event) => {
  if (!event?.detail) {
    return;
  }
  hasUnread.value = Boolean(event.detail.hasUnread);
};

onMounted(() => {
  loadTopNavDisplay();
  loadUnreadNotifications();
  window.addEventListener('notifications:unread', handleUnreadEvent);
});

onBeforeUnmount(() => {
  window.removeEventListener('notifications:unread', handleUnreadEvent);
});
</script>

<style scoped>
.top-nav {
  height: 60px;
  background-color: var(--bg-white);
  border-bottom: 1px solid var(--border-color);
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 var(--spacing-xl);
  box-shadow: var(--shadow-sm);
  transition: all 0.3s ease;
}

.system-title {
  font-size: 18px;
  font-weight: 600;
  color: var(--primary-color);
  margin: 0;
}

.system-link {
  text-decoration: none;
  display: inline-flex;
  align-items: center;
}

.system-link:hover .system-title {
  color: var(--primary-color);
  opacity: 0.9;
}

.nav-right {
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
}

.nav-item {
  display: flex;
  align-items: center;
  margin-left: var(--spacing-sm);
}

.nav-link {
  display: flex;
  align-items: center;
  padding: var(--spacing-xs) var(--spacing-sm);
  border-radius: var(--border-radius-md);
  color: var(--text-medium);
  transition: all 0.3s ease;
  cursor: pointer;
}

.nav-link:hover {
  background-color: var(--bg-medium);
  color: var(--primary-color);
}

.notification-link {
  position: relative;
}

.notification-dot {
  position: absolute;
  top: 6px;
  right: 6px;
  width: 6px;
  height: 6px;
  background: #ef4444;
  border-radius: 50%;
  box-shadow: 0 0 0 2px var(--bg-white);
}

.dark-mode .notification-dot {
  box-shadow: 0 0 0 2px var(--bg-white);
}

.theme-switch {
  padding: var(--spacing-xs) var(--spacing-sm);
}

.user-info {
  display: flex;
  align-items: center;
  cursor: pointer;
  padding: var(--spacing-sm) var(--spacing-md);
  border-radius: var(--border-radius-md);
  transition: background-color 0.3s;
}

.user-info:hover {
  background-color: var(--bg-medium);
}

.username {
  margin-left: var(--spacing-sm);
  margin-right: 4px;
  color: var(--text-medium);
  font-size: 14px;
}

@media (max-width: 768px) {
  .top-nav {
    padding: 0 var(--spacing-md);
  }

  .system-title {
    font-size: 16px;
  }
}
</style>
