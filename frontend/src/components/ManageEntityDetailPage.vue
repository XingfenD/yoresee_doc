<template>
  <PageLayout
    :system-name="systemName"
    :current-language="currentLanguage"
    :is-dark-mode="isDarkMode"
    :user-avatar="userAvatar"
    :username="userInfo?.username || t('common.user')"
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
    <ManageLayout>
      <TitleBar :show-back="true" :compact="true" :back-text="t('common.back')" @back="router.back()">
        <template #actions>
          <el-button type="primary" @click="openEditDialog">
            {{ t('document.edit') }}
          </el-button>
        </template>
      </TitleBar>

      <ManageSection plain>
        <InfoStatsCard
          :title="entityInfo?.name || t('common.unknown')"
          :description="entityInfo?.description || t('common.unknown')"
          :stats="entityStats"
        />
      </ManageSection>

      <ManageSection>
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
      </ManageSection>
    </ManageLayout>

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
import { computed, ref, onMounted } from 'vue';
import { useI18n } from 'vue-i18n';
import { useRoute, useRouter } from 'vue-router';
import { useUserStore } from '@/store/user';
import { useManageShell } from '@/composables/useManageShell';
import { useServerTable } from '@/composables/useServerTable';
import { usePageBoot } from '@/composables/usePageBoot';
import { isActionCancelled, useApiAction } from '@/composables/useApiAction';
import PageLayout from '@/components/PageLayout.vue';
import TitleBar from '@/components/TitleBar.vue';
import ManageLayout from '@/components/ManageLayout.vue';
import ManageSection from '@/components/ManageSection.vue';
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
const { runApi } = useApiAction({ t });

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
  fetchSystemInfo,
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
const { boot } = usePageBoot({ initLanguage, fetchSystemInfo });

const entityInfo = ref(null);
const entityStats = computed(() => [
  { key: 'members', icon: User, label: t('common.members'), value: entityInfo.value?.member_count ?? 0 }
]);
const {
  rows: memberRows,
  page: memberPage,
  pageSize: memberPageSize,
  total: memberTotal,
  keyword: memberSearch,
  load: loadEntityMembers,
  handlePageChange: handleMemberPageChange,
  handleSearch: handleMemberSearch
} = useServerTable({
  initialPageSize: 6,
  fetcher: async ({ page, page_size, keyword }) => {
    const externalId = getExternalId();
    if (!externalId) {
      return { users: [], total: 0 };
    }
    return entityAdapter.value.loadMembers({
      external_id: externalId,
      keyword,
      page,
      page_size
    });
  },
  mapRows: (resp) => resp.users || [],
  onError: (err) => {
    console.error('load entity members failed', err);
  }
});

const {
  rows: memberCandidates,
  page: candidatePage,
  pageSize: candidatePageSize,
  total: candidateTotal,
  keyword: candidateSearch,
  load: loadMemberCandidates,
  handlePageChange: handleCandidatePageChange,
  handleSearch: handleCandidateSearch
} = useServerTable({
  initialPageSize: 3,
  fetcher: ({ page, page_size, keyword }) =>
    listUsers({
      page,
      page_size,
      keyword
    }),
  mapRows: (resp) => resp.users || [],
  onError: (err) => {
    console.error('listUsers failed', err);
  }
});

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
  if (!externalId) {
    entityInfo.value = null;
    return;
  }

  entityInfo.value = await runApi(
    async () => entityAdapter.value.loadDetail(externalId),
    {
      context: 'load entity detail',
      showErrorMessage: false,
      fallback: null
    }
  );
};

const openMemberDialog = async () => {
  showMemberDialog.value = true;
  candidateSearch.value = '';
  candidatePage.value = 1;
  selectedMemberIds.value = memberRows.value.map((member) => member.external_id).filter(Boolean);
  await loadMemberCandidates();
};

const submitMemberUpdate = async () => {
  if (savingMembers.value || !entityInfo.value?.external_id) {
    return;
  }

  savingMembers.value = true;
  await runApi(
    async () => {
      const memberIds = Array.from(new Set(selectedMemberIds.value)).filter(Boolean);
      await entityAdapter.value.update({
        external_id: entityInfo.value.external_id,
        sync_members: true,
        member_user_external_ids: memberIds
      });
      showMemberDialog.value = false;
      await loadEntityMembers();
    },
    {
      context: 'update entity members',
      successMessage: t('message.success'),
      onFinally: () => {
        savingMembers.value = false;
      }
    }
  );
};

const removeMember = async (row) => {
  if (!row?.external_id || !entityInfo.value?.external_id) {
    return;
  }

  await runApi(
    async () => {
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
    },
    {
      context: 'remove member',
      successMessage: t('message.deleteSuccess'),
      ignoreError: isActionCancelled
    }
  );
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
  editing.value = true;
  await runApi(
    async () => {
      await entityAdapter.value.update({
        external_id: editForm.value.external_id,
        name: editForm.value.name.trim(),
        description: editForm.value.description.trim()
      });
      showEditDialog.value = false;
      await loadEntityDetail();
    },
    {
      context: 'update entity',
      successMessage: t('message.success'),
      onFinally: () => {
        editing.value = false;
      }
    }
  );
};

onMounted(() => {
  boot(loadEntityDetail, loadEntityMembers);
});
</script>

<style scoped>
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
</style>
