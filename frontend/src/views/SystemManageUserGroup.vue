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

    <ManageLayout>
      <ManageSection>
        <CommonList
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
        </CommonList>
      </ManageSection>
    </ManageLayout>

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
import { createUserGroup, deleteUserGroup, listUserGroups, updateUserGroup } from '@/services/api';
import { ElMessage, ElMessageBox } from 'element-plus';

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

const showCreateDialog = ref(false);
const creating = ref(false);
const createForm = ref({
  name: '',
  description: ''
});
const showEditDialog = ref(false);
const editing = ref(false);
const editForm = ref({
  external_id: '',
  name: '',
  description: ''
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

const openCreateDialog = () => {
  createForm.value = {
    name: '',
    description: ''
  };
  showCreateDialog.value = true;
};

const openEditDialog = (row) => {
  if (!row?.external_id) {
    return;
  }
  editForm.value = {
    external_id: row.external_id,
    name: row.name || '',
    description: row.description || ''
  };
  showEditDialog.value = true;
};

const submitCreateGroup = async () => {
  if (creating.value) {
    return;
  }
  if (!createForm.value.name.trim()) {
    ElMessage.warning(t('message.warning'));
    return;
  }
  try {
    creating.value = true;
    await createUserGroup({
      name: createForm.value.name.trim(),
      description: createForm.value.description.trim()
    });
    showCreateDialog.value = false;
    await loadUserGroups();
    ElMessage.success(t('message.success'));
  } catch (err) {
    console.error('createUserGroup failed', err);
    ElMessage.error(t('common.requestFailed'));
  } finally {
    creating.value = false;
  }
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
    await loadUserGroups();
    ElMessage.success(t('message.success'));
  } catch (err) {
    console.error('updateUserGroup failed', err);
    ElMessage.error(t('common.requestFailed'));
  } finally {
    editing.value = false;
  }
};

const handleDeleteGroup = async (row) => {
  if (!row?.external_id) {
    return;
  }
  try {
    await ElMessageBox.confirm(t('message.confirmDelete'), t('document.delete'), {
      confirmButtonText: t('button.confirm'),
      cancelButtonText: t('button.cancel'),
      type: 'warning'
    });
    await deleteUserGroup(row.external_id);
    await loadUserGroups();
    ElMessage.success(t('message.deleteSuccess'));
  } catch (err) {
    if (err) {
      console.error('deleteUserGroup failed', err);
    }
  }
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
