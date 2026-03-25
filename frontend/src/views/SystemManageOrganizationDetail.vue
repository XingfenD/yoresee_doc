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
            <h2 class="group-title">{{ orgInfo?.name || t('common.unknown') }}</h2>
            <p class="group-description">
              {{ orgInfo?.description || t('common.unknown') }}
            </p>
            <div class="group-stats">
              <div class="stat-item">
                <el-icon>
                  <User />
                </el-icon>
                <span>{{ t('common.members') }}: {{ orgInfo?.member_count ?? 0 }}</span>
              </div>
            </div>
          </div>
        </div>
      </section>

      <section class="manage-section">
        <div class="section-body">
          <CommonList
            :rows="memberRows"
            :columns="memberColumns"
            :is-dark="isDarkMode"
            row-key="external_id"
            :empty-text="t('message.empty')"
            :show-pagination="true"
            :total="memberTotal"
            v-model:current-page="memberPage"
            v-model:page-size="memberPageSize"
            :page-sizes="[6]"
            @page-change="handleMemberPageChange"
            :show-search="true"
            v-model:search-query="memberSearch"
            :search-placeholder="t('common.search')"
            @search="handleMemberSearch"
            :show-title-bar="true"
          >
            <template #title>{{ t('system.organization.memberList') }}</template>
            <template #toolbar-right>
              <el-button size="small" type="primary" @click="openMemberDialog">
                {{ t('system.organization.manageMembers') }}
              </el-button>
            </template>
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
        <el-button type="primary" :loading="editing" @click="submitEditOrg">
          {{ t('button.confirm') }}
        </el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="showMemberDialog" :title="t('system.organization.memberList')" width="680px">
      <div class="member-dialog">
        <div class="member-dialog__list">
          <CommonList
            :rows="memberCandidates"
            :columns="candidateColumns"
            :is-dark="isDarkMode"
            row-key="external_id"
            :empty-text="t('message.empty')"
            :show-pagination="true"
            :total="candidateTotal"
            v-model:current-page="candidatePage"
            v-model:page-size="candidatePageSize"
            :page-sizes="[3]"
            @page-change="handleCandidatePageChange"
            :show-search="true"
            v-model:search-query="candidateSearch"
            :search-placeholder="t('common.search')"
            @search="handleCandidateSearch"
            :show-title-bar="true"
            :title="t('system.organization.memberList')"
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
import { computed, ref, onMounted, onBeforeUnmount } from 'vue';
import { useI18n } from 'vue-i18n';
import { useRoute, useRouter } from 'vue-router';
import { useUserStore } from '@/store/user';
import PageLayout from '@/components/PageLayout.vue';
import TitleBar from '@/components/TitleBar.vue';
import CommonList from '@/components/CommonList.vue';
import { getOrgNode, listOrgNodeMembers, listUsers, updateOrgNode } from '@/services/api';
import { House, Setting, Ticket, User, UserFilled, OfficeBuilding } from '@element-plus/icons-vue';
import { ElMessage, ElMessageBox } from 'element-plus';

const route = useRoute();
const router = useRouter();
const userStore = useUserStore();
const { locale, t } = useI18n();

const systemName = ref('Yoresee');
const activeMenu = ref('manage-organization');
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

const orgInfo = ref(null);
const memberRows = ref([]);
const memberPage = ref(1);
const memberPageSize = ref(6);
const memberTotal = ref(0);
const memberSearch = ref('');
const memberSearchTimer = ref(null);
const candidateSearch = ref('');
const candidateSearchTimer = ref(null);
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
const candidatePageSize = ref(3);
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

const loadOrgDetail = async () => {
  const externalId = route.params.externalID;
  if (!externalId) {
    return;
  }
  try {
    const resp = await getOrgNode(externalId, { include_children: false });
    orgInfo.value = resp.org_node;
  } catch (err) {
    console.error('getOrgNode failed', err);
    orgInfo.value = null;
  }
};

const loadOrgMembers = async () => {
  const externalId = route.params.externalID;
  if (!externalId) {
    return;
  }
  try {
    const resp = await listOrgNodeMembers({
      external_id: externalId,
      keyword: memberSearch.value,
      page: memberPage.value,
      page_size: memberPageSize.value
    });
    memberRows.value = resp.users || [];
    const totalNumber = Number(resp.total);
    memberTotal.value = Number.isFinite(totalNumber) ? totalNumber : 0;
  } catch (err) {
    console.error('listOrgNodeMembers failed', err);
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
  if (memberSearchTimer.value) {
    clearTimeout(memberSearchTimer.value);
  }
  memberSearchTimer.value = setTimeout(async () => {
    memberPage.value = 1;
    await loadOrgMembers();
  }, 300);
};

const handleCandidateSearch = async () => {
  if (candidateSearchTimer.value) {
    clearTimeout(candidateSearchTimer.value);
  }
  candidateSearchTimer.value = setTimeout(async () => {
    candidatePage.value = 1;
    await loadMemberCandidates();
  }, 300);
};

const handleMemberPageChange = async () => {
  await loadOrgMembers();
};

const handleCandidatePageChange = async () => {
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
    const totalNumber = Number(resp.total);
    candidateTotal.value = Number.isFinite(totalNumber) ? totalNumber : 0;
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
    await updateOrgNode({
      external_id: orgInfo.value.external_id,
      sync_members: true,
      member_user_external_ids: memberIds
    });
    showMemberDialog.value = false;
    await loadOrgMembers();
    ElMessage.success(t('message.success'));
  } catch (err) {
    console.error('updateOrgNode members failed', err);
    ElMessage.error(t('common.requestFailed'));
  } finally {
    savingMembers.value = false;
  }
};

const removeMember = async (row) => {
  if (!row?.external_id || !orgInfo.value?.external_id) {
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
    await updateOrgNode({
      external_id: orgInfo.value.external_id,
      sync_members: true,
      member_user_external_ids: remaining
    });
    await loadOrgMembers();
    ElMessage.success(t('message.deleteSuccess'));
  } catch (err) {
    if (err) {
      console.error('remove member failed', err);
    }
  }
};

const openEditDialog = () => {
  if (!orgInfo.value?.external_id) {
    return;
  }
  editForm.value = {
    external_id: orgInfo.value.external_id,
    name: orgInfo.value.name || '',
    description: orgInfo.value.description || ''
  };
  showEditDialog.value = true;
};

const submitEditOrg = async () => {
  if (editing.value) {
    return;
  }
  if (!editForm.value.name.trim()) {
    ElMessage.warning(t('message.warning'));
    return;
  }
  try {
    editing.value = true;
    await updateOrgNode({
      external_id: editForm.value.external_id,
      name: editForm.value.name.trim(),
      description: editForm.value.description.trim()
    });
    showEditDialog.value = false;
    await loadOrgDetail();
    ElMessage.success(t('message.success'));
  } catch (err) {
    console.error('updateOrgNode failed', err);
    ElMessage.error(t('common.requestFailed'));
  } finally {
    editing.value = false;
  }
};

onMounted(() => {
  initLanguage();
  loadOrgDetail();
  loadOrgMembers();
});

onBeforeUnmount(() => {
  if (memberSearchTimer.value) {
    clearTimeout(memberSearchTimer.value);
  }
  if (candidateSearchTimer.value) {
    clearTimeout(candidateSearchTimer.value);
  }
});
</script>

<style scoped>
.manage-layout {
  display: flex;
  flex-direction: column;
  gap: 0;
}

.manage-section + .manage-section {
  margin-top: var(--spacing-lg);
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

.section-body {
  padding: 0;
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

.dark-mode .manage-section {
  background: #161b22;
  border-color: #2b2f36;
}

.dark-mode .manage-section--plain {
  background: transparent;
  border-color: transparent;
}

.dark-mode .group-info-card {
  background: var(--bg-medium);
  border-color: var(--border-color);
}

.dark-mode .group-title {
  color: #e5e7eb;
}

.dark-mode .group-description {
  color: #9aa4b2;
}

.setting-desc {
  color: var(--text-light);
  font-size: 12px;
}
</style>
