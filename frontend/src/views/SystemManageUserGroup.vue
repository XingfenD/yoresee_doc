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

    <div class="manage-layout">
      <section class="manage-section">
        <div class="section-body">
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
        </div>
      </section>
    </div>

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
import { computed, ref, onMounted, onBeforeUnmount } from 'vue';
import { useI18n } from 'vue-i18n';
import { useRouter } from 'vue-router';
import { useUserStore } from '@/store/user';
import PageLayout from '@/components/PageLayout.vue';
import CommonList from '@/components/CommonList.vue';
import { createUserGroup, deleteUserGroup, listUserGroups, updateUserGroup } from '@/services/api';
import { House, Setting, Ticket, User, UserFilled, OfficeBuilding } from '@element-plus/icons-vue';
import { ElMessage, ElMessageBox } from 'element-plus';

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

const groupRows = ref([]);
const groupPage = ref(1);
const groupPageSize = ref(10);
const groupTotal = ref(0);
const groupKeyword = ref('');
const groupSearchTimer = ref(null);
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

const loadUserGroups = async () => {
  try {
    const resp = await listUserGroups({
      page: groupPage.value,
      page_size: groupPageSize.value,
      keyword: groupKeyword.value.trim() || undefined
    });
    groupRows.value = resp.user_groups || [];
    groupTotal.value = Number(resp.total) || 0;
  } catch (err) {
    console.error('listUserGroups failed', err);
    groupRows.value = [];
    groupTotal.value = 0;
  }
};

const handleGroupPageChange = async (page) => {
  groupPage.value = page;
  await loadUserGroups();
};

const handleGroupSearch = () => {
  if (groupSearchTimer.value) {
    clearTimeout(groupSearchTimer.value);
  }
  groupSearchTimer.value = setTimeout(async () => {
    groupPage.value = 1;
    await loadUserGroups();
  }, 300);
};

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
  initLanguage();
  loadUserGroups();
});

onBeforeUnmount(() => {
  if (groupSearchTimer.value) {
    clearTimeout(groupSearchTimer.value);
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

.link-btn {
  padding: 0;
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
