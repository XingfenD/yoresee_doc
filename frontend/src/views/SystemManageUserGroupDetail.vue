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
    :title="t('system.userGroup.title')"
    @change-language="handleLanguageChange"
    @toggle-theme="toggleTheme"
    @logout="handleLogout"
    @menu-select="handleMenuSelect"
  >
    <template #actions>
      <el-button class="page-action-btn" size="small" @click="openEditDialog">
        {{ t('document.edit') }}
      </el-button>
      <el-button class="page-action-btn" size="small" @click="router.back()">
        {{ t('common.back') }}
      </el-button>
    </template>

    <div class="manage-layout">
      <section class="manage-section">
        <div class="section-header">
          <h3 class="section-title">{{ groupInfo?.name || t('common.unknown') }}</h3>
        </div>
        <div class="section-body">
          <div class="group-meta">
            <div class="group-meta__item">
              <span class="group-meta__label">{{ t('common.description') }}</span>
              <span class="group-meta__value">{{ groupInfo?.description || '-' }}</span>
            </div>
            <div class="group-meta__item">
              <span class="group-meta__label">{{ t('common.members') }}</span>
              <span class="group-meta__value">{{ groupInfo?.member_count ?? 0 }}</span>
            </div>
          </div>
        </div>
      </section>

      <section class="manage-section">
        <div class="section-header section-header--split">
          <h3 class="section-title">{{ t('common.members') }}</h3>
          <el-button size="small" type="primary" @click="openMemberDialog">
            {{ t('button.create') }}
          </el-button>
        </div>
        <div class="section-body">
          <CommonList
            :rows="memberRows"
            :columns="memberColumns"
            :is-dark="isDarkMode"
            row-key="external_id"
            :empty-text="t('message.empty')"
          >
            <template #cell-actions="{ row }">
              <el-button size="small" text type="danger" @click="removeMember(row)">
                {{ t('document.delete') }}
              </el-button>
            </template>
          </CommonList>
        </div>
      </section>
    </div>

    <el-dialog v-model="showEditDialog" :title="t('document.edit')" width="480px">
      <el-form label-position="top" :model="editForm">
        <el-form-item :label="t('common.name')">
          <el-input v-model="editForm.name" maxlength="50" show-word-limit />
        </el-form-item>
        <el-form-item :label="t('common.description')">
          <el-input
            v-model="editForm.description"
            type="textarea"
            :rows="3"
            maxlength="200"
            show-word-limit
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showEditDialog = false">{{ t('button.cancel') }}</el-button>
        <el-button type="primary" :loading="editing" @click="submitEditGroup">
          {{ t('button.confirm') }}
        </el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="showMemberDialog" :title="t('common.members')" width="680px">
      <div class="member-dialog">
        <el-input
          v-model="memberSearch"
          :placeholder="t('common.search')"
          clearable
          @input="handleMemberSearch"
        />
        <div class="member-dialog__list">
          <CommonList
            :rows="memberCandidates"
            :columns="candidateColumns"
            :is-dark="isDarkMode"
            row-key="external_id"
            :empty-text="t('message.empty')"
          >
            <template #cell-actions="{ row }">
              <el-checkbox
                v-model="selectedMemberIds"
                :label="row.external_id"
                class="checkbox-only"
              />
            </template>
          </CommonList>
        </div>
      </div>
      <template #footer>
        <el-button @click="showMemberDialog = false">{{ t('button.cancel') }}</el-button>
        <el-button type="primary" :loading="savingMembers" @click="submitMemberUpdate">
          {{ t('button.confirm') }}
        </el-button>
      </template>
    </el-dialog>
  </PageLayout>
</template>

<script setup>
import { computed, ref, onMounted } from 'vue';
import { useI18n } from 'vue-i18n';
import { useRoute, useRouter } from 'vue-router';
import { useUserStore } from '@/store/user';
import PageLayout from '@/components/PageLayout.vue';
import CommonList from '@/components/CommonList.vue';
import { getUserGroup, listUsers, updateUserGroup } from '@/services/api';
import { House, Setting, Ticket, User, UserFilled, OfficeBuilding } from '@element-plus/icons-vue';
import { ElMessage, ElMessageBox } from 'element-plus';

const route = useRoute();
const router = useRouter();
const userStore = useUserStore();
const { locale, t } = useI18n();

const systemName = ref('Yoresee');
const activeMenu = ref('manage-user-group');
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

const initLanguage = () => {
  const savedLanguage = localStorage.getItem('language');
  if (savedLanguage) {
    currentLanguage.value = savedLanguage;
  }
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

const groupInfo = ref(null);
const memberRows = ref([]);
const showEditDialog = ref(false);
const editing = ref(false);
const editForm = ref({
  external_id: '',
  name: '',
  description: ''
});
const showMemberDialog = ref(false);
const memberCandidates = ref([]);
const memberSearch = ref('');
const selectedMemberIds = ref([]);
const savingMembers = ref(false);

const memberColumns = computed(() => [
  { key: 'username', label: t('common.name'), minWidth: 160 },
  { key: 'email', label: t('user.email') || 'Email', minWidth: 220, flex: 1.4 },
  { key: 'actions', label: t('common.actions'), minWidth: 120, align: 'center' }
]);

const candidateColumns = computed(() => [
  { key: 'username', label: t('common.name'), minWidth: 160 },
  { key: 'email', label: t('user.email') || 'Email', minWidth: 220, flex: 1.4 },
  { key: 'actions', label: t('common.actions'), minWidth: 140, align: 'center' }
]);

const loadGroupDetail = async () => {
  const externalId = route.params.externalID;
  if (!externalId) {
    return;
  }
  try {
    const resp = await getUserGroup(externalId);
    groupInfo.value = resp.user_group;
    memberRows.value = resp.user_group?.members || [];
  } catch (err) {
    console.error('getUserGroup failed', err);
    groupInfo.value = null;
    memberRows.value = [];
  }
};

const openMemberDialog = async () => {
  showMemberDialog.value = true;
  selectedMemberIds.value = memberRows.value.map((member) => member.external_id).filter(Boolean);
  await loadMemberCandidates();
};

const handleMemberSearch = async () => {
  await loadMemberCandidates();
};

const loadMemberCandidates = async () => {
  try {
    const resp = await listUsers({ keyword: memberSearch.value, page: 1, page_size: 200 });
    memberCandidates.value = resp.users || [];
  } catch (err) {
    console.error('listUsers failed', err);
    memberCandidates.value = [];
  }
};

const submitMemberUpdate = async () => {
  if (savingMembers.value) {
    return;
  }
  try {
    savingMembers.value = true;
    const memberIds = Array.from(new Set(selectedMemberIds.value)).filter(Boolean);
    await updateUserGroup({
      external_id: groupInfo.value.external_id,
      sync_members: true,
      member_user_external_ids: memberIds
    });
    showMemberDialog.value = false;
    await loadGroupDetail();
    ElMessage.success(t('message.success'));
  } catch (err) {
    console.error('updateUserGroup members failed', err);
    ElMessage.error(t('common.requestFailed'));
  } finally {
    savingMembers.value = false;
  }
};

const removeMember = async (row) => {
  if (!row?.external_id || !groupInfo.value?.external_id) {
    return;
  }
  try {
    await ElMessageBox.confirm(t('message.confirmDelete'), t('document.delete'), {
      confirmButtonText: t('button.confirm'),
      cancelButtonText: t('button.cancel'),
      type: 'warning'
    });
    const remaining = memberRows.value
      .filter((member) => member.external_id !== row.external_id)
      .map((member) => member.external_id);
    await updateUserGroup({
      external_id: groupInfo.value.external_id,
      sync_members: true,
      member_user_external_ids: remaining
    });
    await loadGroupDetail();
    ElMessage.success(t('message.deleteSuccess'));
  } catch (err) {
    if (err) {
      console.error('remove member failed', err);
    }
  }
};

const openEditDialog = () => {
  if (!groupInfo.value?.external_id) {
    return;
  }
  editForm.value = {
    external_id: groupInfo.value.external_id,
    name: groupInfo.value.name || '',
    description: groupInfo.value.description || ''
  };
  showEditDialog.value = true;
};

const submitEditGroup = async () => {
  if (editing.value) {
    return;
  }
  if (!editForm.value.name.trim()) {
    ElMessage.warning(t('message.warning'));
    return;
  }
  try {
    editing.value = true;
    await updateUserGroup({
      external_id: editForm.value.external_id,
      name: editForm.value.name.trim(),
      description: editForm.value.description.trim()
    });
    showEditDialog.value = false;
    await loadGroupDetail();
    ElMessage.success(t('message.success'));
  } catch (err) {
    console.error('updateUserGroup failed', err);
    ElMessage.error(t('common.requestFailed'));
  } finally {
    editing.value = false;
  }
};

onMounted(() => {
  initLanguage();
  loadGroupDetail();
});
</script>

<style scoped>
.manage-layout {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-lg);
}

.manage-section {
  background: var(--bg-white);
  border: 1px solid var(--border-color);
  border-radius: var(--border-radius-md);
  box-shadow: var(--shadow-sm);
  overflow: hidden;
}

.section-header {
  padding: var(--spacing-md);
  border-bottom: 1px solid var(--border-color);
  background: var(--bg-white);
}

.section-header--split {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: var(--spacing-md);
}

.section-title {
  margin: 0;
  font-size: 16px;
  font-weight: 600;
  color: var(--text-dark);
}

.section-body {
  padding: var(--spacing-md);
}

.group-meta {
  display: grid;
  gap: 12px;
}

.group-meta__item {
  display: flex;
  align-items: baseline;
  gap: 12px;
}

.group-meta__label {
  min-width: 80px;
  color: var(--text-muted);
  font-size: 12px;
}

.group-meta__value {
  color: var(--text-dark);
  font-size: 14px;
}

.member-dialog {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.member-dialog__list {
  max-height: 360px;
  overflow: auto;
}

.checkbox-only :deep(.el-checkbox__label) {
  display: none;
}

.dark-mode .manage-section {
  background: #161b22;
  border-color: #2b2f36;
}

.dark-mode .section-header {
  background: #161b22;
  border-bottom-color: #2b2f36;
}

.dark-mode .group-meta__label {
  color: #9aa4b2;
}

.dark-mode .group-meta__value {
  color: #e5e7eb;
}
</style>
