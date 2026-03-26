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
    <div class="manage-layout">
      <section class="manage-section">
        <div class="section-body">
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
              <span :class="['status-pill', row.status === 1 ? 'is-active' : 'is-disabled']">
                {{ row.status === 1 ? t('user.active') : t('user.disabled') }}
              </span>
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
        </div>
      </section>
    </div>

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
import { computed, ref, onMounted, onBeforeUnmount } from 'vue';
import { useI18n } from 'vue-i18n';
import { useRouter } from 'vue-router';
import { useUserStore } from '@/store/user';
import PageLayout from '@/components/PageLayout.vue';
import CommonList from '@/components/CommonList.vue';
import { useManageShell } from '@/composables/useManageShell';
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

const userRows = ref([]);
const currentPage = ref(1);
const pageSize = ref(10);
const userTotal = ref(0);
const keyword = ref('');
const searchTimer = ref(null);
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

const loadUsers = async () => {
  try {
    const resp = await listUsers({
      page: currentPage.value,
      page_size: pageSize.value,
      keyword: keyword.value.trim() || undefined
    });
    userRows.value = resp.users || [];
    const totalNumber = Number(resp.total);
    userTotal.value = Number.isFinite(totalNumber) ? totalNumber : 0;
  } catch (err) {
    console.error('listUsers failed', err);
    userRows.value = [];
    userTotal.value = 0;
  }
};

const handlePageChange = async (page) => {
  currentPage.value = page;
  await loadUsers();
};

const handleSearch = async () => {
  if (searchTimer.value) {
    clearTimeout(searchTimer.value);
  }
  searchTimer.value = setTimeout(async () => {
    currentPage.value = 1;
    await loadUsers();
  }, 300);
};

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
  initLanguage();
  loadUsers();
});

onBeforeUnmount(() => {
  if (searchTimer.value) {
    clearTimeout(searchTimer.value);
  }
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

.section-body {
  padding: 0;
}

.status-pill {
  display: inline-flex;
  align-items: center;
  padding: 2px 8px;
  border-radius: 999px;
  font-size: 12px;
  font-weight: 600;
}

.status-pill.is-active {
  background: rgba(16, 185, 129, 0.12);
  color: #10b981;
}

.status-pill.is-disabled {
  background: rgba(239, 68, 68, 0.12);
  color: #ef4444;
}

.dark-mode .manage-section {
  background: #161b22;
  border-color: #2b2f36;
}

.dark-mode .section-header {
  background: #161b22;
  border-bottom-color: #2b2f36;
}
</style>
