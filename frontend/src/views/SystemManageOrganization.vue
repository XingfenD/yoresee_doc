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

    <div class="manage-layout">
      <section class="manage-section">
        <div class="section-header">
          <h3 class="section-title">{{ t('system.organization.placeholderTitle') }}</h3>
        </div>
        <div class="section-body">
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
        </div>
      </section>
    </div>
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
import CommonList from '@/components/CommonList.vue';
import { House, Setting, Ticket, User, UserFilled, OfficeBuilding } from '@element-plus/icons-vue';
import { listOrgNodes, createOrgNode, updateOrgNode } from '@/services/api';
import { ElMessage } from 'element-plus';

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

const orgTreeData = ref([]);
const orgLoading = ref(false);
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

const currentPage = ref(1);
const pageSize = ref(6);
const pageSizes = [6, 10, 20];
const paginationLayout = 'total, prev, pager, next';
const totalOrgs = ref(0);

const pagedOrgNodes = computed(() => orgTreeData.value);

const fetchOrgNodes = async () => {
  orgLoading.value = true;
  try {
    const resp = await listOrgNodes({
      parent_external_id: '',
      page: currentPage.value,
      page_size: pageSize.value,
      include_children: true
    });
    totalOrgs.value = Number(resp.total) || 0;
    orgTreeData.value = resp.org_nodes || [];
  } catch (error) {
    orgTreeData.value = [];
    totalOrgs.value = 0;
  } finally {
    orgLoading.value = false;
  }
};

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

const handlePageChange = (page) => {
  currentPage.value = page;
  fetchOrgNodes();
};

const handleSizeChange = (size) => {
  pageSize.value = size;
  currentPage.value = 1;
  fetchOrgNodes();
};

onMounted(() => {
  initLanguage();
  fetchOrgNodes();
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

.section-header {
  padding: var(--spacing-md);
  border-bottom: 1px solid var(--border-color);
  background: var(--bg-white);
}

.section-title {
  margin: 0;
  font-size: 16px;
  font-weight: 600;
  color: var(--text-dark);
}

.section-body {
  padding: var(--spacing-md);
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
