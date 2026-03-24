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
    :title="''"
    content-padding="xl"
    @change-language="handleLanguageChange"
    @toggle-theme="toggleTheme"
    @logout="handleLogout"
    @menu-select="handleMenuSelect"
  >
    <div class="manage-layout">
      <TitleBar :show-back="true" :back-text="t('common.back')" @back="router.back()">
        <template #actions>
          <el-button type="primary" @click="openEditDialog">
            {{ t('document.edit') }}
          </el-button>
        </template>
      </TitleBar>

      <section class="manage-section manage-section--plain">
        <div class="section-body section-body--card">
          <div class="group-info-card">
            <h2 class="group-title">{{ groupInfo?.name || t('common.unknown') }}</h2>
            <p class="group-description">
              {{ groupInfo?.description || t('common.unknown') }}
            </p>
            <div class="group-stats">
              <div class="stat-item">
                <el-icon>
                  <User />
                </el-icon>
                <span>{{ t('common.members') }}: {{ groupInfo?.member_count ?? 0 }}</span>
              </div>
            </div>
          </div>
        </div>
      </section>

      <section class="manage-section">
        <div class="section-header section-header--split">
          <h3 class="section-title">{{ t('system.userGroup.memberList') }}</h3>
          <div class="member-actions">
            <el-input
              v-model="memberSearch"
              :placeholder="t('common.search')"
              clearable
              class="member-search"
              @input="handleMemberSearch"
            />
            <el-button size="small" type="primary" @click="openMemberDialog">
              {{ t('system.userGroup.manageMembers') }}
            </el-button>
          </div>
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
          <div class="pagination-container" v-if="memberTotal > memberPageSize">
            <el-pagination
              v-model:current-page="memberPage"
              v-model:page-size="memberPageSize"
              :page-sizes="[20, 50, 100]"
              :total="memberTotal"
              layout="total, sizes, prev, pager, next, jumper"
              @size-change="handleMemberPageSize"
              @current-change="handleMemberPageChange"
            />
          </div>
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
          v-model="candidateSearch"
          :placeholder="t('common.search')"
          clearable
          @input="handleCandidateSearch"
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
          <div class="pagination-container" v-if="candidateTotal > candidatePageSize">
            <el-pagination
              v-model:current-page="candidatePage"
              v-model:page-size="candidatePageSize"
              :page-sizes="[20, 50, 100]"
              :total="candidateTotal"
              layout="total, sizes, prev, pager, next, jumper"
              @size-change="handleCandidatePageSize"
              @current-change="handleCandidatePageChange"
            />
          </div>
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
import TitleBar from '@/components/TitleBar.vue';
import CommonList from '@/components/CommonList.vue';
import { getUserGroup, listUserGroupMembers, listUsers, updateUserGroup } from '@/services/api';
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
const memberPage = ref(1);
const memberPageSize = ref(20);
const memberTotal = ref(0);
const memberSearch = ref('');
const candidateSearch = ref('');
const showEditDialog = ref(false);
const editing = ref(false);
const editForm = ref({
  external_id: '',
  name: '',
  description: ''
});
const showMemberDialog = ref(false);
const memberCandidates = ref([]);
const selectedMemberIds = ref([]);
const savingMembers = ref(false);
const candidatePage = ref(1);
const candidatePageSize = ref(20);
const candidateTotal = ref(0);

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
  } catch (err) {
    console.error('getUserGroup failed', err);
    groupInfo.value = null;
  }
};

const loadGroupMembers = async () => {
  const externalId = route.params.externalID;
  if (!externalId) {
    return;
  }
  try {
    const resp = await listUserGroupMembers({
      external_id: externalId,
      keyword: memberSearch.value,
      page: memberPage.value,
      page_size: memberPageSize.value
    });
    memberRows.value = resp.users || [];
    memberTotal.value = resp.total ?? 0;
  } catch (err) {
    console.error('listUserGroupMembers failed', err);
    memberRows.value = [];
    memberTotal.value = 0;
  }
};

const openMemberDialog = async () => {
  showMemberDialog.value = true;
  candidateSearch.value = '';
  candidatePage.value = 1;
  selectedMemberIds.value = memberRows.value.map((member) => member.external_id).filter(Boolean);
  await loadMemberCandidates();
};

const handleMemberSearch = async () => {
  memberPage.value = 1;
  await loadGroupMembers();
};

const handleCandidateSearch = async () => {
  candidatePage.value = 1;
  await loadMemberCandidates();
};

const handleMemberPageChange = async () => {
  await loadGroupMembers();
};

const handleMemberPageSize = async () => {
  memberPage.value = 1;
  await loadGroupMembers();
};

const handleCandidatePageChange = async () => {
  await loadMemberCandidates();
};

const handleCandidatePageSize = async () => {
  candidatePage.value = 1;
  await loadMemberCandidates();
};

const loadMemberCandidates = async () => {
  try {
    const resp = await listUsers({
      keyword: candidateSearch.value,
      page: candidatePage.value,
      page_size: candidatePageSize.value
    });
    memberCandidates.value = resp.users || [];
    candidateTotal.value = resp.total ?? 0;
  } catch (err) {
    console.error('listUsers failed', err);
    memberCandidates.value = [];
    candidateTotal.value = 0;
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
    await loadGroupMembers();
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
    await loadGroupMembers();
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
  loadGroupMembers();
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

.manage-section--plain {
  background: transparent;
  border: none;
  box-shadow: none;
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

.section-body--card {
  padding: 0;
}

.group-info-card {
  background-color: var(--bg-white);
  border-radius: var(--border-radius-md);
  box-shadow: var(--shadow-sm);
  padding: var(--spacing-lg);
  border: 1px solid var(--border-color);
}

.group-title {
  margin: 0 0 var(--spacing-md) 0;
  font-size: 24px;
  font-weight: 600;
  color: var(--text-dark);
}

.group-description {
  margin: 0 0 var(--spacing-lg) 0;
  font-size: 16px;
  color: var(--text-medium);
  line-height: 1.6;
}

.group-stats {
  display: flex;
  gap: var(--spacing-lg);
  flex-wrap: wrap;
}

.stat-item {
  display: flex;
  align-items: center;
  gap: var(--spacing-xs);
  color: var(--text-medium);
  font-size: 14px;
}

.stat-item .el-icon {
  color: var(--primary-color);
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

.member-actions {
  display: flex;
  align-items: center;
  gap: 12px;
}

.member-search {
  width: 220px;
}

.dark-mode .manage-section {
  background: #161b22;
  border-color: #2b2f36;
}

.dark-mode .manage-section--plain {
  background: transparent;
  border-color: transparent;
}

.dark-mode .section-header {
  background: #161b22;
  border-bottom-color: #2b2f36;
}

.dark-mode .group-info-card {
  background: #161b22;
  border-color: #2b2f36;
}

.dark-mode .group-title {
  color: #e5e7eb;
}

.dark-mode .group-description {
  color: #9aa4b2;
}
</style>
