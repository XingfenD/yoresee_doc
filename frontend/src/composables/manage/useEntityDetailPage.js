import { computed, onMounted, ref } from 'vue';
import { useI18n } from 'vue-i18n';
import { useRoute, useRouter } from 'vue-router';
import { useUserStore } from '@/store/user';
import { useManageShell } from '@/composables/shell/useManageShell';
import { useServerTable } from '@/composables/list/useServerTable';
import { usePageBoot } from '@/composables/shell/usePageBoot';
import { isActionCancelled, useApiAction } from '@/composables/actions/useApiAction';
import { usePageTitle } from '@/composables/usePageTitle';
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

export function useEntityDetailPage(props) {
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
  const pageTitleLabel = computed(() =>
    props.entityType === 'user-group' ? t('pageTitle.userGroup') : t('pageTitle.organization')
  );
  usePageTitle(pageTitleLabel, computed(() => entityInfo.value?.name || ''));
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

  return {
    router,
    systemName,
    currentLanguage,
    isDarkMode,
    userAvatar,
    userInfo,
    activeMenu,
    manageMenuItems,
    handleLanguageChange,
    toggleTheme,
    handleLogout,
    handleMenuSelect,
    entityLabels,
    entityInfo,
    entityStats,
    memberRows,
    memberColumns,
    memberTotal,
    memberPage,
    memberPageSize,
    handleMemberPageChange,
    memberSearch,
    handleMemberSearch,
    openMemberDialog,
    memberCandidates,
    candidateColumns,
    candidateTotal,
    candidatePage,
    candidatePageSize,
    handleCandidatePageChange,
    candidateSearch,
    handleCandidateSearch,
    selectedMemberIds,
    savingMembers,
    submitMemberUpdate,
    showMemberDialog,
    showEditDialog,
    editing,
    editForm,
    openEditDialog,
    submitEdit,
    removeMember
  };
}
