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
    layout="list"
    @change-language="handleLanguageChange"
    @toggle-theme="toggleTheme"
    @logout="handleLogout"
    @menu-select="handleMenuSelect"
  >
    <template #actions>
      <el-button class="page-action-btn" type="primary" size="small" @click="openCreateDialog">
        {{ t('system.userGroup.create') }}
      </el-button>
    </template>

    <ManageCrudListSection
      :rows="groupRows"
      :columns="groupColumns"
      :is-dark="isDarkMode"
      row-key="external_id"
      :empty-text="t('message.empty')"
      :show-pagination="true"
      :total="groupTotal"
      v-model:current-page="groupPage"
      v-model:page-size="groupPageSize"
      :page-sizes="[10, 20, 50]"
      @page-change="handleGroupPageChange"
      :show-search="true"
      v-model:search-query="groupKeyword"
      :search-placeholder="t('common.search')"
      @search="handleGroupSearch"
      :show-title-bar="true"
      :title="t('system.userGroup.placeholderTitle')"
    >
      <template #cell-name="{ row }">
        <el-button class="link-btn" type="primary" text @click="goToDetail(row)">
          {{ row?.name || '-' }}
        </el-button>
      </template>
      <template #cell-actions="{ row }">
        <el-button size="small" text type="primary" @click="openEditDialog(row)">
          {{ t('document.edit') }}
        </el-button>
        <el-button size="small" text type="danger" @click="handleDeleteGroup(row)">
          {{ t('document.delete') }}
        </el-button>
      </template>
    </ManageCrudListSection>

    <el-dialog v-model="showCreateDialog" :title="t('system.userGroup.title')" width="480px">
      <el-form label-position="top" :model="createForm">
        <el-form-item :label="t('common.name')">
          <el-input v-model="createForm.name" />
        </el-form-item>
        <el-form-item :label="t('common.description')">
          <el-input v-model="createForm.description" type="textarea" :rows="3" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showCreateDialog = false">{{ t('button.cancel') }}</el-button>
        <el-button type="primary" :loading="creating" @click="submitCreateGroup">
          {{ t('button.confirm') }}
        </el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="showEditDialog" :title="t('document.edit')" width="480px">
      <el-form label-position="top" :model="editForm">
        <el-form-item :label="t('common.name')">
          <el-input v-model="editForm.name" />
        </el-form-item>
        <el-form-item :label="t('common.description')">
          <el-input v-model="editForm.description" type="textarea" :rows="3" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showEditDialog = false">{{ t('button.cancel') }}</el-button>
        <el-button type="primary" :loading="editing" @click="submitEditGroup">
          {{ t('button.confirm') }}
        </el-button>
      </template>
    </el-dialog>
  </PageLayout>
</template>

<script setup>
import { computed, onMounted } from 'vue';
import { useI18n } from 'vue-i18n';
import { useRouter } from 'vue-router';
import { useUserStore } from '@/store/user';
import PageLayout from '@/components/PageLayout.vue';
import ManageCrudListSection from '@/components/ManageCrudListSection.vue';
import { useManageShell } from '@/composables/useManageShell';
import { useServerTable } from '@/composables/useServerTable';
import { usePageBoot } from '@/composables/usePageBoot';
import { isActionCancelled, useApiAction } from '@/composables/useApiAction';
import { useCrudDialog } from '@/composables/useCrudDialog';
import { createUserGroup, deleteUserGroup, listUserGroups, updateUserGroup } from '@/services/api';
import { ElMessage, ElMessageBox } from 'element-plus';

const router = useRouter();
const userStore = useUserStore();
const { locale, t } = useI18n();
const { runApi, createApiErrorHandler } = useApiAction({ t });

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
  defaultActiveMenu: 'manage-user-group'
});
const { boot } = usePageBoot({ initLanguage, fetchSystemInfo });

const {
  rows: groupRows,
  page: groupPage,
  pageSize: groupPageSize,
  total: groupTotal,
  keyword: groupKeyword,
  load: loadUserGroups,
  handlePageChange: handleGroupPageChange,
  handleSearch: handleGroupSearch
} = useServerTable({
  initialPageSize: 10,
  fetcher: ({ page, page_size, keyword }) =>
    listUserGroups({
      page,
      page_size,
      keyword
    }),
  mapRows: (resp) => resp.user_groups || [],
  onError: (err) => {
    console.error('listUserGroups failed', err);
  }
});

const {
  visible: showCreateDialog,
  submitting: creating,
  form: createForm,
  open: openCreateDialog,
  submit: submitCreateGroup
} = useCrudDialog({
  initialForm: () => ({
    name: '',
    description: ''
  }),
  validate: (form) => {
    if (!form.name.trim()) {
      ElMessage.warning(t('message.warning'));
      return false;
    }
    return true;
  },
  submitRequest: async (form) => {
    await createUserGroup({
      name: form.name.trim(),
      description: form.description.trim()
    });
  },
  onSuccess: async () => {
    await loadUserGroups();
    ElMessage.success(t('message.success'));
  },
  onError: createApiErrorHandler({ context: 'createUserGroup' })
});

const {
  visible: showEditDialog,
  submitting: editing,
  form: editForm,
  open: openEditGroupDialog,
  submit: submitEditGroup
} = useCrudDialog({
  initialForm: () => ({
    external_id: '',
    name: '',
    description: ''
  }),
  mapOpenForm: (row) => ({
    external_id: row?.external_id || '',
    name: row?.name || '',
    description: row?.description || ''
  }),
  validate: (form) => {
    if (!form.name.trim()) {
      ElMessage.warning(t('message.warning'));
      return false;
    }
    return true;
  },
  submitRequest: async (form) => {
    await updateUserGroup({
      external_id: form.external_id,
      name: form.name.trim(),
      description: form.description.trim()
    });
  },
  onSuccess: async () => {
    await loadUserGroups();
    ElMessage.success(t('message.success'));
  },
  onError: createApiErrorHandler({ context: 'updateUserGroup' })
});

const groupColumns = computed(() => [
  { key: 'name', label: t('common.name'), minWidth: 180 },
  { key: 'member_count', label: t('common.members'), minWidth: 120, align: 'center' },
  { key: 'description', label: t('common.description'), minWidth: 260, flex: 1.6 },
  { key: 'actions', label: t('common.actions'), minWidth: 120, align: 'center' }
]);

const goToDetail = (row) => {
  if (!row?.external_id) {
    return;
  }
  router.push(`/manage/user_group/${row.external_id}`);
};

const openEditDialog = (row) => {
  if (!row?.external_id) {
    return;
  }
  openEditGroupDialog(row);
};

const handleDeleteGroup = async (row) => {
  if (!row?.external_id) {
    return;
  }

  await runApi(
    async () => {
      await ElMessageBox.confirm(t('message.confirmDelete'), t('document.delete'), {
        confirmButtonText: t('button.confirm'),
        cancelButtonText: t('button.cancel'),
        type: 'warning'
      });
      await deleteUserGroup(row.external_id);
      await loadUserGroups();
    },
    {
      context: 'deleteUserGroup',
      successMessage: t('message.deleteSuccess'),
      ignoreError: isActionCancelled
    }
  );
};

onMounted(() => {
  boot(loadUserGroups);
});
</script>

<style scoped>
.link-btn {
  padding: 0;
}
</style>
