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
    :title="t('system.organization.title')"
    layout="list"
    @change-language="handleLanguageChange"
    @toggle-theme="toggleTheme"
    @logout="handleLogout"
    @menu-select="handleMenuSelect"
  >
    <template #actions>
      <el-button class="page-action-btn" type="primary" size="small" @click="openCreateDialog">
        {{ t('system.organization.create') }}
      </el-button>
    </template>

    <ManageLayout>
      <ManageSection :title="t('system.organization.placeholderTitle')" body-padding="md">
        <CommonList
          mode="tree"
          :rows="pagedOrgNodes"
          :columns="orgColumns"
          :is-dark="isDarkMode"
          row-key="external_id"
          :empty-text="t('message.empty')"
          tree-column-key="name"
          tree-key-field="external_id"
          :tree-loading="orgLoading"
          show-pagination
          :current-page="currentPage"
          :page-size="pageSize"
          :page-sizes="pageSizes"
          :total="totalOrgs"
          :pagination-layout="paginationLayout"
          @page-change="handlePageChange"
          @size-change="handleSizeChange"
        >
          <template #tree-cell="{ row }">
            <el-button class="link-btn" type="primary" text @click="goToDetail(row)">
              {{ row?.name || '-' }}
            </el-button>
          </template>
          <template #cell-actions="{ row }">
            <el-button size="small" text type="primary" @click="openEditDialog(row)">
              {{ t('common.edit') }}
            </el-button>
            <el-button size="small" text type="primary" @click="openCreateChildDialog(row)">
              {{ t('system.organization.createChild') }}
            </el-button>
          </template>
        </CommonList>
      </ManageSection>
    </ManageLayout>
  </PageLayout>

  <el-dialog
    v-model="showCreateDialog"
    :title="dialogTitle"
    width="420px"
    :close-on-click-modal="false"
  >
    <el-form label-position="top">
      <el-form-item :label="t('common.name')" required>
        <el-input v-model="createForm.name" maxlength="50" />
      </el-form-item>
      <el-form-item :label="t('common.description')">
        <el-input v-model="createForm.description" type="textarea" :rows="3" maxlength="200" />
      </el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="showCreateDialog = false">{{ t('common.cancel') }}</el-button>
      <el-button type="primary" :loading="creating" @click="submitCreate">
        {{ t('common.confirm') }}
      </el-button>
    </template>
  </el-dialog>

  <el-dialog
    v-model="showEditDialog"
    :title="t('system.organization.edit')"
    width="420px"
    :close-on-click-modal="false"
  >
    <el-form label-position="top">
      <el-form-item :label="t('common.name')" required>
        <el-input v-model="editForm.name" maxlength="50" />
      </el-form-item>
      <el-form-item :label="t('common.description')">
        <el-input v-model="editForm.description" type="textarea" :rows="3" maxlength="200" />
      </el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="showEditDialog = false">{{ t('common.cancel') }}</el-button>
      <el-button type="primary" :loading="editing" @click="submitEdit">
        {{ t('common.confirm') }}
      </el-button>
    </template>
  </el-dialog>
</template>

<script setup>
import { computed, ref, onMounted } from 'vue';
import { useI18n } from 'vue-i18n';
import { useRouter } from 'vue-router';
import { useUserStore } from '@/store/user';
import PageLayout from '@/components/PageLayout.vue';
import ManageLayout from '@/components/ManageLayout.vue';
import ManageSection from '@/components/ManageSection.vue';
import CommonList from '@/components/CommonList.vue';
import { useManageShell } from '@/composables/useManageShell';
import { useServerTable } from '@/composables/useServerTable';
import { usePageBoot } from '@/composables/usePageBoot';
import { listOrgNodes, createOrgNode, updateOrgNode } from '@/services/api';
import { ElMessage } from 'element-plus';

const router = useRouter();
const userStore = useUserStore();
const { locale, t } = useI18n();

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
  defaultActiveMenu: 'manage-organization'
});
const { boot } = usePageBoot({ initLanguage, fetchSystemInfo });

const {
  rows: orgTreeData,
  loading: orgLoading,
  page: currentPage,
  pageSize,
  total: totalOrgs,
  load: fetchOrgNodes,
  handlePageChange,
  handleSizeChange
} = useServerTable({
  initialPageSize: 6,
  fetcher: ({ page, page_size }) =>
    listOrgNodes({
      parent_external_id: '',
      page,
      page_size,
      include_children: true
    }),
  mapRows: (resp) => resp.org_nodes || [],
  onError: (error) => {
    console.error('fetchOrgNodes failed', error);
  }
});

const showCreateDialog = ref(false);
const creating = ref(false);
const createParent = ref(null);
const showEditDialog = ref(false);
const editing = ref(false);
const editForm = ref({
  external_id: '',
  name: '',
  description: ''
});
const createForm = ref({
  name: '',
  description: ''
});

const orgColumns = computed(() => [
  { key: 'name', label: t('common.name'), minWidth: 220 },
  { key: 'member_count', label: t('common.members'), minWidth: 120, align: 'center' },
  { key: 'description', label: t('common.description'), minWidth: 220 },
  { key: 'actions', label: t('common.actions'), minWidth: 140, align: 'center' }
]);

const pageSizes = [6, 10, 20];
const paginationLayout = 'total, prev, pager, next';

const pagedOrgNodes = computed(() => orgTreeData.value);

const dialogTitle = computed(() => {
  if (createParent.value?.name) {
    return t('system.organization.createChildWith', { name: createParent.value.name });
  }
  return t('system.organization.create');
});

const openCreateDialog = () => {
  createForm.value = {
    name: '',
    description: ''
  };
  createParent.value = null;
  showCreateDialog.value = true;
};

const openCreateChildDialog = (row) => {
  createForm.value = {
    name: '',
    description: ''
  };
  createParent.value = row || null;
  showCreateDialog.value = true;
};

const openEditDialog = (row) => {
  editForm.value = {
    external_id: row?.external_id || '',
    name: row?.name || '',
    description: row?.description || ''
  };
  showEditDialog.value = true;
};

const goToDetail = (row) => {
  if (!row?.external_id) {
    return;
  }
  router.push(`/manage/organization/${row.external_id}`);
};

const submitCreate = async () => {
  if (!createForm.value.name.trim()) {
    ElMessage.warning(t('message.nameRequiredGeneric'));
    return;
  }
  creating.value = true;
  try {
    await createOrgNode({
      creator_user_external_id: userInfo.value?.external_id || '',
      name: createForm.value.name.trim(),
      description: createForm.value.description.trim(),
      parent_external_id: createParent.value?.external_id || '',
      member_user_external_ids: []
    });
    ElMessage.success(t('message.createSuccessGeneric'));
    showCreateDialog.value = false;
    fetchOrgNodes();
  } catch (error) {
    ElMessage.error(t('message.createFailedGeneric'));
  } finally {
    creating.value = false;
  }
};

const submitEdit = async () => {
  if (!editForm.value.name.trim()) {
    ElMessage.warning(t('message.nameRequiredGeneric'));
    return;
  }
  editing.value = true;
  try {
    await updateOrgNode({
      external_id: editForm.value.external_id,
      name: editForm.value.name.trim(),
      description: editForm.value.description.trim(),
      sync_members: false,
      member_user_external_ids: []
    });
    ElMessage.success(t('message.saveSuccessGeneric'));
    showEditDialog.value = false;
    fetchOrgNodes();
  } catch (error) {
    ElMessage.error(t('message.saveFailedGeneric'));
  } finally {
    editing.value = false;
  }
};

onMounted(() => {
  boot(fetchOrgNodes);
});
</script>

<style scoped>
.link-btn {
  padding: 0;
}
</style>
