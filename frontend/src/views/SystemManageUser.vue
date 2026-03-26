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
    :title="t('system.user.title')"
    layout="list"
    @change-language="handleLanguageChange"
    @toggle-theme="toggleTheme"
    @logout="handleLogout"
    @menu-select="handleMenuSelect"
  >
    <ManageLayout>
      <ManageSection>
        <CommonList
          :rows="userRows"
          :columns="userColumns"
          :is-dark="isDarkMode"
          row-key="external_id"
          :empty-text="t('message.empty')"
          :show-pagination="true"
          :total="userTotal"
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :page-sizes="[10]"
          @page-change="handlePageChange"
          :show-search="true"
          v-model:search-query="keyword"
          :search-placeholder="t('common.search')"
          @search="handleSearch"
          :show-title-bar="true"
          :title="t('system.user.placeholderTitle')"
        >
          <template #cell-status="{ row }">
            <AppTag :type="row.status === 1 ? 'success' : 'danger'" size="small">
              {{ row.status === 1 ? t('user.active') : t('user.disabled') }}
            </AppTag>
          </template>
          <template #cell-actions="{ row }">
            <el-button size="small" text type="primary" @click="openEditDialog(row)">
              {{ t('document.edit') }}
            </el-button>
            <el-button
              size="small"
              text
              :type="row.status === 1 ? 'danger' : 'success'"
              @click="toggleUserStatus(row)"
            >
              {{ row.status === 1 ? t('user.ban') : t('user.unban') }}
            </el-button>
          </template>
        </CommonList>
      </ManageSection>
    </ManageLayout>

    <el-dialog v-model="showEditDialog" :title="t('document.edit')" width="480px">
      <el-form label-position="top" :model="editForm">
        <el-form-item :label="t('user.name')">
          <el-input v-model="editForm.username" maxlength="50" show-word-limit />
        </el-form-item>
        <el-form-item :label="t('user.email')">
          <el-input v-model="editForm.email" maxlength="100" show-word-limit />
        </el-form-item>
        <el-form-item :label="t('common.name')">
          <el-input v-model="editForm.nickname" maxlength="50" show-word-limit />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showEditDialog = false">{{ t('button.cancel') }}</el-button>
        <el-button type="primary" :loading="editing" @click="submitEditUser">
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
import AppTag from '@/components/AppTag.vue';
import { useManageShell } from '@/composables/useManageShell';
import { useServerTable } from '@/composables/useServerTable';
import { usePageBoot } from '@/composables/usePageBoot';
import { listUsers, updateUser } from '@/services/api';
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
  defaultActiveMenu: 'manage-user'
});
const { boot } = usePageBoot({ initLanguage, fetchSystemInfo });

const {
  rows: userRows,
  page: currentPage,
  pageSize,
  total: userTotal,
  keyword,
  load: loadUsers,
  handlePageChange,
  handleSearch
} = useServerTable({
  initialPageSize: 10,
  fetcher: ({ page, page_size, keyword: currentKeyword }) =>
    listUsers({
      page,
      page_size,
      keyword: currentKeyword
    }),
  mapRows: (resp) => resp.users || [],
  onError: (err) => {
    console.error('listUsers failed', err);
  }
});

const showEditDialog = ref(false);
const editing = ref(false);
const editForm = ref({
  external_id: '',
  username: '',
  email: '',
  nickname: ''
});

const userColumns = computed(() => [
  { key: 'username', label: t('user.name'), minWidth: 140 },
  { key: 'email', label: t('user.email'), minWidth: 220, flex: 1.4 },
  { key: 'nickname', label: t('common.name'), minWidth: 140 },
  { key: 'status', label: t('user.status'), minWidth: 120 },
  { key: 'created_at', label: t('common.createdAt'), minWidth: 140 },
  { key: 'actions', label: t('common.actions'), minWidth: 140, align: 'center' }
]);

const openEditDialog = (row) => {
  if (!row?.external_id) {
    return;
  }
  editForm.value = {
    external_id: row.external_id,
    username: row.username || '',
    email: row.email || '',
    nickname: row.nickname || ''
  };
  showEditDialog.value = true;
};

const submitEditUser = async () => {
  if (editing.value) {
    return;
  }
  if (!editForm.value.username.trim()) {
    ElMessage.warning(t('message.warning'));
    return;
  }
  try {
    editing.value = true;
    await updateUser({
      external_id: editForm.value.external_id,
      username: editForm.value.username.trim(),
      email: editForm.value.email.trim(),
      nickname: editForm.value.nickname.trim()
    });
    showEditDialog.value = false;
    await loadUsers();
    ElMessage.success(t('message.success'));
  } catch (err) {
    console.error('updateUser failed', err);
    ElMessage.error(t('common.requestFailed'));
  } finally {
    editing.value = false;
  }
};

const toggleUserStatus = async (row) => {
  if (!row?.external_id) {
    return;
  }
  try {
    const isActive = row.status === 1;
    await ElMessageBox.confirm(
      isActive ? t('message.confirmBan') : t('message.confirmUnban'),
      isActive ? t('user.ban') : t('user.unban'),
      {
      confirmButtonText: t('button.confirm'),
      cancelButtonText: t('button.cancel'),
      type: 'warning'
      }
    );
    const nextStatus = isActive ? 0 : 1;
    await updateUser({
      external_id: row.external_id,
      status: nextStatus
    });
    await loadUsers();
    ElMessage.success(t('message.success'));
  } catch (err) {
    if (err) {
      console.error('toggleUserStatus failed', err);
    }
  }
};

onMounted(() => {
  boot(loadUsers);
});
</script>
