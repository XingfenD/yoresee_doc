<template>
  <PageLayout
    :system-name="systemName"
    :current-language="currentLanguage"
    :is-dark-mode="isDarkMode"
    :user-avatar="userAvatar"
    :username="userInfo?.username || '用户'"
    :active-menu="activeMenu"
    :side-menu-items="userMenuItems"
    :title="t('user.invite.title')"
    @change-language="handleLanguageChange"
    @toggle-theme="toggleTheme"
    @logout="handleLogout"
    @menu-select="handleMenuSelect"
  >
    <template #actions>
      <el-button class="page-action-btn" type="primary" size="small" @click="openCreateDialog">
        {{ t('user.invite.create') }}
      </el-button>
    </template>

    <div class="invite-section" :class="{ 'is-dark': isDarkMode }">
      <div class="invite-list">
        <div class="invite-row invite-row--head">
          <div class="cell cell-code">{{ t('user.invite.code') }}</div>
          <div class="cell cell-status">{{ t('user.invite.status') }}</div>
          <div class="cell cell-usage">{{ t('user.invite.usage') }}</div>
          <div class="cell cell-created-at">{{ t('user.invite.createdAt') }}</div>
          <div class="cell cell-created-by">{{ t('user.invite.createdBy') }}</div>
          <div class="cell cell-expires-at">{{ t('user.invite.expiresAt') }}</div>
          <div class="cell cell-actions">{{ t('user.invite.actions') }}</div>
        </div>
        <div v-for="item in inviteList" :key="item.code" class="invite-row">
          <div class="cell cell-code">{{ item.code }}</div>
          <div class="cell cell-status">
            <el-tag :type="statusType(item.status)" size="small">
              {{ statusLabel(item.status) }}
            </el-tag>
          </div>
          <div class="cell cell-usage">{{ item.used }}/{{ item.max }}</div>
          <div class="cell cell-created-at">{{ item.created_at }}</div>
          <div class="cell cell-created-by">{{ item.created_by }}</div>
          <div class="cell cell-expires-at">{{ item.expires_at }}</div>
          <div class="cell cell-actions">
            <el-button size="small" text type="primary">{{ t('user.invite.pause') }}</el-button>
            <el-button size="small" text type="danger">{{ t('user.invite.delete') }}</el-button>
          </div>
        </div>
      </div>
    </div>

    <InviteCreateDialog v-model="showCreateDialog" @submit="handleCreateInvite" />
  </PageLayout>
</template>

<script setup>
import { computed, ref, onMounted, onBeforeUnmount } from 'vue';
import { useI18n } from 'vue-i18n';
import { useRouter } from 'vue-router';
import { useUserStore } from '@/store/user';
import PageLayout from '@/components/PageLayout.vue';
import InviteCreateDialog from '@/components/InviteCreateDialog.vue';
import { House, User, Ticket, Setting } from '@element-plus/icons-vue';

const router = useRouter();
const userStore = useUserStore();
const { locale, t } = useI18n();

const systemName = ref('Yoresee');
const activeMenu = ref('user-invite');
const isDarkMode = ref(document.documentElement.classList.contains('dark-mode'));
let classObserver = null;

const userInfo = computed(() => userStore.userInfo);
const userAvatar = computed(() => userInfo.value?.avatar || 'https://cube.elemecdn.com/0/88/03b0d39583f48206768a7534e55bcpng.png');

const userMenuItems = [
  { key: 'home', labelKey: 'navigation.home', icon: House, route: '/' },
  { key: 'user-center', labelKey: 'user.menu.center', icon: User, route: '/user_info/example' },
  { key: 'user-invite', labelKey: 'user.menu.invite', icon: Ticket, route: '/user_info/invatations' },
  { key: 'user-security', labelKey: 'user.menu.security', icon: Setting, route: '/user_info/example' }
];

const currentLanguage = computed({
  get: () => locale.value,
  set: (value) => {
    locale.value = value;
    localStorage.setItem('language', value);
  }
});

const handleLanguageChange = (command) => {
  currentLanguage.value = command;
};

const initLanguage = () => {
  const savedLanguage = localStorage.getItem('language');
  if (savedLanguage) {
    currentLanguage.value = savedLanguage;
  }
};

const toggleTheme = () => {
  const next = !document.documentElement.classList.contains('dark-mode');
  document.documentElement.classList.toggle('dark-mode', next);
  localStorage.setItem('darkMode', next ? 'true' : 'false');
  isDarkMode.value = next;
};

const handleLogout = () => {
  userStore.logout();
  router.push('/login');
};

const handleMenuSelect = (key) => {
  activeMenu.value = key;
};

const inviteList = ref([
  {
    code: 'YORE-8K2P-9Q1M',
    created_at: '2026-03-18 10:12',
    created_by: 'admin',
    expires_at: '2026-04-18 23:59',
    status: 'active',
    max: 20,
    used: 6
  },
  {
    code: 'YORE-7D4X-1ZA3',
    created_at: '2026-03-10 09:30',
    created_by: 'admin',
    expires_at: '2026-03-31 23:59',
    status: 'active',
    max: 5,
    used: 5
  },
  {
    code: 'YORE-2M9K-4L8N',
    created_at: '2026-02-20 15:40',
    created_by: 'admin',
    expires_at: '2026-03-05 23:59',
    status: 'expired',
    max: 10,
    used: 7
  }
]);

const statusType = (status) => {
  if (status === 'active') return 'success';
  if (status === 'expired') return 'info';
  return 'warning';
};

const statusLabel = (status) => {
  if (status === 'active') return t('user.invite.active');
  if (status === 'expired') return t('user.invite.expired');
  return t('user.invite.disabled');
};

const showCreateDialog = ref(false);

const openCreateDialog = () => {
  showCreateDialog.value = true;
};

const handleCreateInvite = (payload) => {
  // TODO: hook to backend
  console.log('create invite payload', payload);
};

const sectionStyle = computed(() => ({}));

onMounted(() => {
  initLanguage();
  classObserver = new MutationObserver(() => {
    isDarkMode.value = document.documentElement.classList.contains('dark-mode');
  });
  classObserver.observe(document.documentElement, { attributes: true, attributeFilter: ['class'] });
});

onBeforeUnmount(() => {
  classObserver?.disconnect();
});
</script>

<style scoped>
.invite-section {
  background: var(--bg-white);
  border: 1px solid var(--border-color);
  border-radius: var(--border-radius-md);
  box-shadow: var(--shadow-sm);
  overflow: hidden;
  padding: var(--spacing-md);
}

.invite-section.is-dark {
  background: #1e1e1e;
  border-color: #2c2c2c;
}

.invite-list {
  display: flex;
  flex-direction: column;
  gap: 0;
}

.invite-row {
  display: grid;
  grid-template-columns: 1.3fr 0.9fr 0.9fr 1fr 0.8fr 1fr 0.9fr;
  align-items: center;
  padding: 10px 12px;
  border-bottom: 1px solid var(--border-color);
  background: var(--bg-white);
  color: var(--text-dark);
}

.invite-row--head {
  font-weight: 600;
  background: var(--bg-light);
  color: var(--text-dark);
}

.invite-section.is-dark .invite-row {
  background: #1e1e1e;
  color: #f5f5f5;
  border-bottom-color: #2c2c2c;
}

.invite-section.is-dark .invite-row--head {
  background: #2a2a2a;
  color: #f5f5f5;
}

.cell {
  min-width: 0;
  font-size: 13px;
}

.cell-actions {
  display: flex;
  gap: 8px;
}

@media (max-width: 1200px) {
  .invite-row {
    grid-template-columns: 1.2fr 0.8fr 0.8fr 1fr 0.8fr 1fr;
    row-gap: 6px;
  }

  .cell-actions {
    grid-column: 1 / -1;
  }
}
</style>
