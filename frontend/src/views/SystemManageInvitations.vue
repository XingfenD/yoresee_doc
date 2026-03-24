<template>
  <PageLayout
    :system-name="systemName"
    :current-language="currentLanguage"
    :is-dark-mode="isDarkMode"
    :user-avatar="userAvatar"
    :username="userInfo?.username || '用户'"
    :active-menu="activeMenu"
    :side-menu-items="manageMenuItems"
    sidebar-scene="manage"
    :title="t('system.invite.title')"
    layout="list"
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

    <el-tabs v-model="activeTab" class="common-tabs">
      <el-tab-pane :label="t('system.invite.tabs.list')" name="list">
        <CommonList
          :rows="inviteList"
          :columns="inviteColumns"
          :is-dark="isDarkMode"
          row-key="code"
          :empty-text="t('message.empty')"
        >
          <template #cell-status="{ value }">
            <el-tag :type="inviteStatusType(value)" size="small">
              {{ inviteStatusLabel(value) }}
            </el-tag>
          </template>
      <template #cell-usage="{ row }">
        {{ row.used }}/{{ row.max === null ? '-' : row.max }}
      </template>
      <template #cell-code="{ row }">
        <el-tooltip :content="row.note || t('user.invite.notePlaceholder')" placement="top">
          <span class="invite-code" @click="copyInviteCode(row.code)">{{ row.code }}</span>
        </el-tooltip>
      </template>
      <template #cell-actions="{ row }">
        <el-button size="small" text type="primary" @click="handlePauseInvite(row)">
          {{ row.disabled ? t('user.invite.resume') : t('user.invite.pause') }}
        </el-button>
        <el-button size="small" text type="danger" @click="handleDeleteInvite(row)">
              {{ t('user.invite.delete') }}
            </el-button>
          </template>
        </CommonList>
      </el-tab-pane>
      <el-tab-pane :label="t('system.invite.tabs.records')" name="records">
        <CommonList
          :rows="inviteRecords"
          :columns="recordColumns"
          :is-dark="isDarkMode"
          row-key="id"
          :empty-text="t('system.invite.records.empty')"
        >
          <template #cell-status="{ value }">
            <el-tag :type="value === 'success' ? 'success' : 'warning'" size="small">
              {{ value === 'success' ? t('system.invite.records.success') : t('system.invite.records.failed') }}
            </el-tag>
          </template>
        </CommonList>
      </el-tab-pane>
    </el-tabs>

    <InviteCreateDialog v-model="showCreateDialog" @submit="handleCreateInvite" />
  </PageLayout>
</template>

<script setup>
import { computed, ref, onMounted, watch } from 'vue';
import { useI18n } from 'vue-i18n';
import { useRouter } from 'vue-router';
import { useUserStore } from '@/store/user';
import PageLayout from '@/components/PageLayout.vue';
import CommonList from '@/components/CommonList.vue';
import InviteCreateDialog from '@/components/InviteCreateDialog.vue';
import { House, Setting, Ticket, User, UserFilled, OfficeBuilding } from '@element-plus/icons-vue';
import { ElMessage, ElMessageBox } from 'element-plus';
import { listInvitations, listInvitationRecords, createInvitation, updateInvitation, deleteInvitation } from '@/services/api';

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
  { key: 'manage-user', labelKey: 'system.menu.user', icon: User, route: '/manage/user' },
  { key: 'manage-user-group', labelKey: 'system.menu.userGroup', icon: UserFilled, route: '/manage/user_group' },
  { key: 'manage-organization', labelKey: 'system.menu.organization', icon: OfficeBuilding, route: '/manage/organization' },
  { key: 'manage-invite', labelKey: 'system.menu.invite', icon: Ticket, route: '/manage/invitations' },
  { key: 'manage-security', labelKey: 'system.menu.security', icon: Setting, route: '/manage/security' }
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

const inviteColumns = computed(() => [
  { key: 'code', label: t('user.invite.code'), minWidth: 180 },
  { key: 'status', label: t('user.invite.status'), minWidth: 120, align: 'center' },
  { key: 'usage', label: t('user.invite.usage'), minWidth: 110, align: 'center' },
  { key: 'created_at', label: t('user.invite.createdAt'), minWidth: 160 },
  { key: 'created_by', label: t('user.invite.createdBy'), minWidth: 140 },
  { key: 'expires_at', label: t('user.invite.expiresAt'), minWidth: 160 },
  { key: 'actions', label: t('user.invite.actions'), minWidth: 160, align: 'right' }
]);

const recordColumns = computed(() => [
  { key: 'code', label: t('system.invite.records.code'), minWidth: 180 },
  { key: 'used_by', label: t('system.invite.records.usedBy'), minWidth: 160 },
  { key: 'used_at', label: t('system.invite.records.usedAt'), minWidth: 180 },
  { key: 'status', label: t('system.invite.records.result'), minWidth: 120, align: 'center' }
]);

const inviteStatusType = (status) => {
  if (status === 'active') return 'success';
  if (status === 'expired') return 'info';
  return 'warning';
};

const inviteStatusLabel = (status) => {
  if (status === 'active') return t('user.invite.active');
  if (status === 'expired') return t('user.invite.expired');
  return t('user.invite.disabled');
};

const inviteList = ref([]);
const inviteRecords = ref([]);
const inviteLoaded = ref(false);
const recordsLoaded = ref(false);
const inviteLoading = ref(false);
const recordsLoading = ref(false);

const showCreateDialog = ref(false);

const openCreateDialog = () => {
  showCreateDialog.value = true;
};

const formatDateYYYYMMDD = (date) => {
  const pad = (num) => `${num}`.padStart(2, '0');
  return `${date.getFullYear()}-${pad(date.getMonth() + 1)}-${pad(date.getDate())}`;
};

const toRfc3339EndOfDay = (dateStr) => {
  if (!dateStr) return '';
  if (dateStr.includes('T')) {
    return dateStr;
  }
  const [year, month, day] = dateStr.split('-').map((value) => Number(value));
  if (!year || !month || !day) return '';
  const local = new Date(year, month - 1, day, 23, 59, 59);
  return local.toISOString();
};

const resolveInviteStatus = (invite) => {
  if (invite.disabled) return 'disabled';
  if (invite.expires_at) {
    const expireTs = Date.parse(invite.expires_at);
    if (!Number.isNaN(expireTs) && expireTs < Date.now()) return 'expired';
  }
  if (typeof invite.max_used_cnt === 'number' && typeof invite.used_cnt === 'number') {
    if (invite.used_cnt >= invite.max_used_cnt) return 'expired';
  }
  return 'active';
};

const toNumber = (value, fallback = null) => {
  if (value === null || value === undefined) return fallback;
  const num = Number(value);
  return Number.isFinite(num) ? num : fallback;
};

const mapInviteRow = (invite) => ({
  code: invite.code,
  created_at: invite.created_at || '-',
  created_by: invite.created_by_name || invite.created_by_external_id || t('common.unknown'),
  expires_at: invite.expires_at || '-',
  status: resolveInviteStatus(invite),
  max: toNumber(invite.max_used_cnt, null),
  used: toNumber(invite.used_cnt, 0),
  disabled: invite.disabled,
  note: invite.note || ''
});

const mapRecordRow = (record) => ({
  id: record.id,
  code: record.code,
  used_by: record.used_by || '-',
  used_at: record.used_at || '-',
  status: record.status
});

const fetchInvitations = async () => {
  if (inviteLoading.value) return;
  inviteLoading.value = true;
  try {
    const resp = await listInvitations({ page: 1, page_size: 50 });
    inviteList.value = (resp.invitations || []).map(mapInviteRow);
    inviteLoaded.value = true;
  } catch (err) {
    console.error('listInvitations failed', err);
    inviteList.value = [];
  } finally {
    inviteLoading.value = false;
  }
};

const fetchInvitationRecords = async () => {
  if (recordsLoading.value) return;
  recordsLoading.value = true;
  try {
    const resp = await listInvitationRecords({ page: 1, page_size: 50 });
    inviteRecords.value = (resp.records || []).map(mapRecordRow);
    recordsLoaded.value = true;
  } catch (err) {
    console.error('listInvitationRecords failed', err);
    inviteRecords.value = [];
  } finally {
    recordsLoading.value = false;
  }
};

const handleCreateInvite = async (payload) => {
  try {
    let expiresAt = payload.expires_at;
    if (payload.expire_type === 'days' && payload.expire_days) {
      const target = new Date();
      target.setDate(target.getDate() + Number(payload.expire_days));
      expiresAt = formatDateYYYYMMDD(target);
    }
    await createInvitation({
      expires_at: expiresAt ? toRfc3339EndOfDay(expiresAt) : undefined,
      max_used_cnt: payload.limit_enabled ? payload.max_usage : undefined,
      note: payload.note
    });
    ElMessage.success(t('message.success'));
    await fetchInvitations();
  } catch (err) {
    console.error('createInvitation failed', err);
    ElMessage.error(t('common.requestFailed'));
  }
};

const handlePauseInvite = async (row) => {
  if (!row?.code) return;
  try {
    await updateInvitation({
      code: row.code,
      disabled: !row.disabled
    });
    ElMessage.success(t('message.success'));
    await fetchInvitations();
  } catch (err) {
    console.error('updateInvitation failed', err);
    ElMessage.error(t('common.requestFailed'));
  }
};

const handleDeleteInvite = async (row) => {
  if (!row?.code) return;
  try {
    await ElMessageBox.confirm(t('message.confirmDelete'), t('user.invite.delete'), {
      confirmButtonText: t('button.confirm'),
      cancelButtonText: t('button.cancel'),
      type: 'warning'
    });
    await deleteInvitation(row.code);
    ElMessage.success(t('message.deleteSuccess'));
    await fetchInvitations();
  } catch (err) {
    if (err) {
      console.error('deleteInvitation failed', err);
    }
  }
};

const copyInviteCode = async (code) => {
  if (!code) return;
  try {
    await navigator.clipboard.writeText(code);
    ElMessage.success(t('common.copySuccess'));
  } catch (err) {
    console.error('copy invite code failed', err);
    ElMessage.error(t('common.copyFailed'));
  }
};

const initLanguage = () => {
  const savedLanguage = localStorage.getItem('language');
  if (savedLanguage) {
    currentLanguage.value = savedLanguage;
  }
};

watch(activeTab, async (tab) => {
  if (tab === 'list' && !inviteLoaded.value) {
    await fetchInvitations();
  }
  if (tab === 'records' && !recordsLoaded.value) {
    await fetchInvitationRecords();
  }
});

onMounted(async () => {
  initLanguage();
  await fetchInvitations();
});
</script>

<style scoped>
.invite-code {
  cursor: pointer;
  color: var(--primary-color);
}
</style>
