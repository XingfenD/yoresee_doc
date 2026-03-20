<template>
  <PageLayout
    :system-name="systemName"
    :current-language="currentLanguage"
    :is-dark-mode="isDarkMode"
    :user-avatar="userAvatar"
    :username="userInfo?.username || '用户'"
    :active-menu="activeMenu"
    :side-menu-items="manageMenuItems"
    :title="t('system.invite.title')"
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

    <el-tabs v-model="activeTab" class="manage-tabs">
      <el-tab-pane :label="t('system.invite.tabs.list')" name="list">
        <InviteList :items="inviteList" :is-dark="isDarkMode" />
      </el-tab-pane>
      <el-tab-pane :label="t('system.invite.tabs.records')" name="records">
        <div class="records-section" :class="{ 'is-dark': isDarkMode }">
          <div class="records-row records-row--head">
            <div class="cell">{{ t('system.invite.records.code') }}</div>
            <div class="cell">{{ t('system.invite.records.usedBy') }}</div>
            <div class="cell">{{ t('system.invite.records.usedAt') }}</div>
            <div class="cell">{{ t('system.invite.records.result') }}</div>
          </div>
          <div v-for="record in inviteRecords" :key="record.id" class="records-row">
            <div class="cell">{{ record.code }}</div>
            <div class="cell">{{ record.used_by }}</div>
            <div class="cell">{{ record.used_at }}</div>
            <div class="cell">
              <el-tag :type="record.status === 'success' ? 'success' : 'warning'" size="small">
                {{ record.status === 'success' ? t('system.invite.records.success') : t('system.invite.records.failed') }}
              </el-tag>
            </div>
          </div>
          <el-empty v-if="inviteRecords.length === 0" :description="t('system.invite.records.empty')" />
        </div>
      </el-tab-pane>
    </el-tabs>

    <InviteCreateDialog v-model="showCreateDialog" @submit="handleCreateInvite" />
  </PageLayout>
</template>

<script setup>
import { computed, ref, onMounted } from 'vue';
import { useI18n } from 'vue-i18n';
import { useRouter } from 'vue-router';
import { useUserStore } from '@/store/user';
import PageLayout from '@/components/PageLayout.vue';
import InviteList from '@/components/InviteList.vue';
import InviteCreateDialog from '@/components/InviteCreateDialog.vue';
import { House, Setting, Ticket } from '@element-plus/icons-vue';

const router = useRouter();
const userStore = useUserStore();
const { locale, t } = useI18n();

const systemName = ref('Yoresee');
const activeMenu = ref('manage-invite');
const activeTab = ref('list');
const isDarkMode = computed(() => userStore.darkMode);

const userInfo = computed(() => userStore.userInfo);
const userAvatar = computed(() => userInfo.value?.avatar || 'https://cube.elemecdn.com/0/88/03b0d39583f48206768a7534e55bcpng.png');

const manageMenuItems = [
  { key: 'home', labelKey: 'navigation.home', icon: House, route: '/' },
  { key: 'manage-security', labelKey: 'system.menu.security', icon: Setting, route: '/manage/security' },
  { key: 'manage-invite', labelKey: 'system.menu.invite', icon: Ticket, route: '/manage/invitations' }
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

const toggleTheme = () => {
  userStore.toggleDarkMode();
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
    code: 'SYS-9KD2-8LMQ',
    created_at: '2026-03-19 09:12',
    created_by: 'admin',
    expires_at: '2026-05-01 23:59',
    status: 'active',
    max: 50,
    used: 12
  },
  {
    code: 'SYS-3PA9-2XKD',
    created_at: '2026-03-01 11:20',
    created_by: 'admin',
    expires_at: '2026-03-25 23:59',
    status: 'expired',
    max: 5,
    used: 3
  }
]);

const inviteRecords = ref([
  {
    id: 1,
    code: 'SYS-9KD2-8LMQ',
    used_by: 'user_lee',
    used_at: '2026-03-19 10:01',
    status: 'success'
  },
  {
    id: 2,
    code: 'SYS-9KD2-8LMQ',
    used_by: 'user_wang',
    used_at: '2026-03-19 10:05',
    status: 'failed'
  }
]);

const showCreateDialog = ref(false);

const openCreateDialog = () => {
  showCreateDialog.value = true;
};

const handleCreateInvite = (payload) => {
  // TODO: hook to backend
  console.log('create system invite payload', payload);
};

const initLanguage = () => {
  const savedLanguage = localStorage.getItem('language');
  if (savedLanguage) {
    currentLanguage.value = savedLanguage;
  }
};

onMounted(() => {
  initLanguage();
});
</script>

<style scoped>
.manage-tabs :deep(.el-tabs__header) {
  margin-bottom: var(--spacing-lg);
}

.records-section {
  background: var(--bg-white);
  border: 1px solid var(--border-color);
  border-radius: var(--border-radius-md);
  padding: var(--spacing-md);
  display: flex;
  flex-direction: column;
  gap: var(--spacing-sm);
}

.records-section.is-dark {
  background: #1e1e1e;
  border-color: #2c2c2c;
}

.records-row {
  display: grid;
  grid-template-columns: 1.2fr 1fr 1.2fr 0.8fr;
  gap: var(--spacing-md);
  padding: 10px 12px;
  border-bottom: 1px solid var(--border-color);
}

.records-row--head {
  font-weight: 600;
  background: var(--bg-light);
}

.records-section.is-dark .records-row {
  border-bottom-color: #2c2c2c;
  color: #f5f5f5;
}

.records-section.is-dark .records-row--head {
  background: #2a2a2a;
  color: #f5f5f5;
}
</style>
