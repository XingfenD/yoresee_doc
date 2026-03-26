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
          <InfoStatsCard
            :title="entityInfo?.name || t('common.unknown')"
            :description="entityInfo?.description || t('common.unknown')"
            :stats="entityStats"
          />
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
            <template #title>{{ t(entityLabels.memberListKey) }}</template>
            <template #toolbar-right>
              <el-button size="small" type="primary" @click="openMemberDialog">
                {{ t(entityLabels.manageMembersKey) }}
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
        <el-button type="primary" :loading="editing" @click="submitEdit">
          {{ t('button.confirm') }}
        </el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="showMemberDialog" :title="t(entityLabels.memberListKey)" width="680px">
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
            :title="t(entityLabels.memberListKey)"
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
import { useManageShell } from '@/composables/useManageShell';
import PageLayout from '@/components/PageLayout.vue';
import TitleBar from '@/components/TitleBar.vue';
import CommonList from '@/components/CommonList.vue';
import InfoStatsCard from '@/components/InfoStatsCard.vue';
import {
  getOrgNode,
  listOrgNodeMembers,
  updateOrgNode,
  getUserGroup,
  listUserGroupMembers,
  updateUserGroup,
  listUsers
} from '@/services/api';
import { User } from '@element-plus/icons-vue';
import { ElMessage, ElMessageBox } from 'element-plus';

const props = defineProps({
  entityType: {
    type: String,
    required: true,
    validator: (value) => ['organization', 'user-group'].includes(value)
  }
});

const route = useRoute();
const router = useRouter();
const userStore = useUserStore();
const { locale, t } = useI18n();

const entityAdapters = {
  organization: {
    activeMenu: 'manage-organization',
    memberListKey: 'system.organization.memberList',
    manageMembersKey: 'system.organization.manageMembers',
    loadDetail: async (externalId) => {
      const resp = await getOrgNode(externalId, { include_children: false });
      return resp.org_node;
    },
    loadMembers: listOrgNodeMembers,
    update: updateOrgNode
  },
  'user-group': {
    activeMenu: 'manage-user-group',
    memberListKey: 'system.userGroup.memberList',
    manageMembersKey: 'system.userGroup.manageMembers',
    loadDetail: async (externalId) => {
      const resp = await getUserGroup(externalId);
      return resp.user_group;
    },
    loadMembers: listUserGroupMembers,
    update: updateUserGroup
  }
};

const entityAdapter = computed(() => entityAdapters[props.entityType]);
const entityLabels = computed(() => ({
  memberListKey: entityAdapter.value.memberListKey,
  manageMembersKey: entityAdapter.value.manageMembersKey
}));

const {
  systemName,
  activeMenu,
  isDarkMode,
  userInfo,
  userAvatar,
  manageMenuItems,
  currentLanguage,
  initLanguage,
  handleLanguageChange,
  toggleTheme,
  handleLogout,
  handleMenuSelect
} = useManageShell({
  locale,
  router,
  userStore,
  defaultActiveMenu: entityAdapter.value.activeMenu
});

const entityInfo = ref(null);
const entityStats = computed(() => [
  { key: 'members', icon: User, label: t('common.members'), value: entityInfo.value?.member_count ?? 0 }
]);
const memberRows = ref([]);
const memberPage = ref(1);
const memberPageSize = ref(6);
const memberTotal = ref(0);
const memberSearch = ref('');
const memberSearchTimer = ref(null);

const candidateSearch = ref('');
const candidateSearchTimer = ref(null);
const candidatePage = ref(1);
const candidatePageSize = ref(3);
const candidateTotal = ref(0);
const memberCandidates = ref([]);
const selectedMemberIds = ref([]);
const savingMembers = ref(false);

const showEditDialog = ref(false);
const editing = ref(false);
const editForm = ref({
  external_id: '',
  name: '',
  description: ''
});

const showMemberDialog = ref(false);

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

const getExternalId = () => route.params.externalID;

const loadEntityDetail = async () => {
  const externalId = getExternalId();
  if (!externalId) return;
  try {
    entityInfo.value = await entityAdapter.value.loadDetail(externalId);
  } catch (err) {
    console.error('load entity detail failed', err);
    entityInfo.value = null;
  }
};

const loadEntityMembers = async () => {
  const externalId = getExternalId();
  if (!externalId) return;
  try {
    const resp = await entityAdapter.value.loadMembers({
      external_id: externalId,
      keyword: memberSearch.value,
      page: memberPage.value,
      page_size: memberPageSize.value
    });
    memberRows.value = resp.users || [];
    const totalNumber = Number(resp.total);
    memberTotal.value = Number.isFinite(totalNumber) ? totalNumber : 0;
  } catch (err) {
    console.error('load entity members failed', err);
    memberRows.value = [];
    memberTotal.value = 0;
  }
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
    await loadEntityMembers();
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
  await loadEntityMembers();
};

const handleCandidatePageChange = async () => {
  await loadMemberCandidates();
};

const submitMemberUpdate = async () => {
  if (savingMembers.value || !entityInfo.value?.external_id) {
    return;
  }
  try {
    savingMembers.value = true;
    const memberIds = Array.from(new Set(selectedMemberIds.value)).filter(Boolean);
    await entityAdapter.value.update({
      external_id: entityInfo.value.external_id,
      sync_members: true,
      member_user_external_ids: memberIds
    });
    showMemberDialog.value = false;
    await loadEntityMembers();
    ElMessage.success(t('message.success'));
  } catch (err) {
    console.error('update entity members failed', err);
    ElMessage.error(t('common.requestFailed'));
  } finally {
    savingMembers.value = false;
  }
};

const removeMember = async (row) => {
  if (!row?.external_id || !entityInfo.value?.external_id) {
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
    await entityAdapter.value.update({
      external_id: entityInfo.value.external_id,
      sync_members: true,
      member_user_external_ids: remaining
    });
    await loadEntityMembers();
    ElMessage.success(t('message.deleteSuccess'));
  } catch (err) {
    if (err) {
      console.error('remove member failed', err);
    }
  }
};

const openEditDialog = () => {
  if (!entityInfo.value?.external_id) {
    return;
  }
  editForm.value = {
    external_id: entityInfo.value.external_id,
    name: entityInfo.value.name || '',
    description: entityInfo.value.description || ''
  };
  showEditDialog.value = true;
};

const submitEdit = async () => {
  if (editing.value) {
    return;
  }
  if (!editForm.value.name.trim()) {
    ElMessage.warning(t('message.warning'));
    return;
  }
  try {
    editing.value = true;
    await entityAdapter.value.update({
      external_id: editForm.value.external_id,
      name: editForm.value.name.trim(),
      description: editForm.value.description.trim()
    });
    showEditDialog.value = false;
    await loadEntityDetail();
    ElMessage.success(t('message.success'));
  } catch (err) {
    console.error('update entity failed', err);
    ElMessage.error(t('common.requestFailed'));
  } finally {
    editing.value = false;
  }
};

onMounted(() => {
  initLanguage();
  loadEntityDetail();
  loadEntityMembers();
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

</style>
