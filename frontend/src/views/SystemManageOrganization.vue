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
            row-key="id"
            :empty-text="t('message.empty')"
            tree-column-key="name"
            tree-key-field="id"
            show-pagination
            :current-page="currentPage"
            :page-size="pageSize"
            :page-sizes="pageSizes"
            :total="totalOrgs"
            :pagination-layout="paginationLayout"
            @page-change="handlePageChange"
            @size-change="handleSizeChange"
          >
            <template #cell-actions="{ row }">
              <el-button size="small" text type="primary">
                {{ t('common.edit') }}
              </el-button>
            </template>
          </CommonList>
        </div>
      </section>
    </div>
  </PageLayout>
</template>

<script setup>
import { computed, ref, onMounted } from 'vue';
import { useI18n } from 'vue-i18n';
import { useRouter } from 'vue-router';
import { useUserStore } from '@/store/user';
import PageLayout from '@/components/PageLayout.vue';
import CommonList from '@/components/CommonList.vue';
import { House, Setting, Ticket, User, UserFilled, OfficeBuilding } from '@element-plus/icons-vue';

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

const orgTreeData = ref([
  {
    id: 'org-1',
    name: '产品中心',
    manager: 'Mia Liu',
    members: 18,
    level: 'L1',
    isFolder: true,
    children: [
      {
        id: 'org-1-1',
        name: '产品规划组',
        manager: '韩东',
        members: 6,
        level: 'L2',
        isFolder: true,
        children: [
          { id: 'org-1-1-1', name: '王晨 · 产品经理', manager: '-', members: 1, level: 'L3', isFolder: false, isLeaf: true },
          { id: 'org-1-1-2', name: '李娜 · 需求分析', manager: '-', members: 1, level: 'L3', isFolder: false, isLeaf: true }
        ]
      },
      {
        id: 'org-1-2',
        name: '设计组',
        manager: '安晴',
        members: 5,
        level: 'L2',
        isFolder: true,
        children: [
          { id: 'org-1-2-1', name: '周舟 · 视觉设计', manager: '-', members: 1, level: 'L3', isFolder: false, isLeaf: true },
          { id: 'org-1-2-2', name: '安晴 · 交互设计', manager: '-', members: 1, level: 'L3', isFolder: false, isLeaf: true }
        ]
      }
    ]
  },
  {
    id: 'org-2',
    name: '研发中心',
    manager: 'Alex Chen',
    members: 32,
    level: 'L1',
    isFolder: true,
    children: [
      {
        id: 'org-2-1',
        name: '平台组',
        manager: 'Sam Zhao',
        members: 10,
        level: 'L2',
        isFolder: true,
        children: [
          {
            id: 'org-2-1-1',
            name: '后端小组',
            manager: '陈可',
            members: 4,
            level: 'L3',
            isFolder: true,
            children: [
              { id: 'org-2-1-1-1', name: '陈可 · 后端工程师', manager: '-', members: 1, level: 'L4', isFolder: false, isLeaf: true },
              { id: 'org-2-1-1-2', name: '杨柳 · 服务端', manager: '-', members: 1, level: 'L4', isFolder: false, isLeaf: true }
            ]
          },
          {
            id: 'org-2-1-2',
            name: 'DevOps 小组',
            manager: '赵敏',
            members: 3,
            level: 'L3',
            isFolder: true,
            children: [
              { id: 'org-2-1-2-1', name: '赵敏 · DevOps', manager: '-', members: 1, level: 'L4', isFolder: false, isLeaf: true },
              { id: 'org-2-1-2-2', name: '彭飞 · 运维工程师', manager: '-', members: 1, level: 'L4', isFolder: false, isLeaf: true }
            ]
          }
        ]
      },
      {
        id: 'org-2-2',
        name: '移动端组',
        manager: '王雨',
        members: 8,
        level: 'L2',
        isFolder: true,
        children: [
          {
            id: 'org-2-2-1',
            name: 'iOS 小组',
            manager: '孙凯',
            members: 3,
            level: 'L3',
            isFolder: true,
            children: [
              { id: 'org-2-2-1-1', name: '孙凯 · iOS', manager: '-', members: 1, level: 'L4', isFolder: false, isLeaf: true },
              { id: 'org-2-2-1-2', name: '陆川 · iOS', manager: '-', members: 1, level: 'L4', isFolder: false, isLeaf: true }
            ]
          },
          {
            id: 'org-2-2-2',
            name: 'Android 小组',
            manager: '王雨',
            members: 3,
            level: 'L3',
            isFolder: true,
            children: [
              { id: 'org-2-2-2-1', name: '王雨 · Android', manager: '-', members: 1, level: 'L4', isFolder: false, isLeaf: true },
              { id: 'org-2-2-2-2', name: '周童 · Android', manager: '-', members: 1, level: 'L4', isFolder: false, isLeaf: true }
            ]
          }
        ]
      }
    ]
  },
  {
    id: 'org-3',
    name: '增长中心',
    manager: 'Lina Wu',
    members: 12,
    level: 'L1',
    isFolder: true,
    children: [
      {
        id: 'org-3-1',
        name: '品牌组',
        manager: '赵然',
        members: 4,
        level: 'L2',
        isFolder: true,
        children: [
          { id: 'org-3-1-1', name: '赵然 · 品牌经理', manager: '-', members: 1, level: 'L3', isFolder: false, isLeaf: true }
        ]
      },
      {
        id: 'org-3-2',
        name: '内容组',
        manager: '简宁',
        members: 5,
        level: 'L2',
        isFolder: true,
        children: [
          {
            id: 'org-3-2-1',
            name: '内容策划小组',
            manager: '马丁',
            members: 2,
            level: 'L3',
            isFolder: true,
            children: [
              { id: 'org-3-2-1-1', name: '马丁 · 内容策划', manager: '-', members: 1, level: 'L4', isFolder: false, isLeaf: true },
              { id: 'org-3-2-1-2', name: '林安 · 内容策划', manager: '-', members: 1, level: 'L4', isFolder: false, isLeaf: true }
            ]
          },
          {
            id: 'org-3-2-2',
            name: '社媒运营小组',
            manager: '简宁',
            members: 2,
            level: 'L3',
            isFolder: true,
            children: [
              { id: 'org-3-2-2-1', name: '简宁 · 社媒运营', manager: '-', members: 1, level: 'L4', isFolder: false, isLeaf: true },
              { id: 'org-3-2-2-2', name: '高楠 · 社媒运营', manager: '-', members: 1, level: 'L4', isFolder: false, isLeaf: true }
            ]
          }
        ]
      }
    ]
  },
  {
    id: 'org-4',
    name: '商业化中心',
    manager: '董卓',
    members: 9,
    level: 'L1',
    isFolder: true,
    children: [
      {
        id: 'org-4-1',
        name: '销售组',
        manager: '董卓',
        members: 6,
        level: 'L2',
        isFolder: true,
        children: [
          {
            id: 'org-4-1-1',
            name: '大客户组',
            manager: '董卓',
            members: 3,
            level: 'L3',
            isFolder: true,
            children: [
              { id: 'org-4-1-1-1', name: '董卓 · 销售经理', manager: '-', members: 1, level: 'L4', isFolder: false, isLeaf: true },
              { id: 'org-4-1-1-2', name: '胡杨 · 大客户', manager: '-', members: 1, level: 'L4', isFolder: false, isLeaf: true }
            ]
          },
          {
            id: 'org-4-1-2',
            name: '客户成功组',
            manager: '宋韵',
            members: 2,
            level: 'L3',
            isFolder: true,
            children: [
              { id: 'org-4-1-2-1', name: '宋韵 · 客户成功', manager: '-', members: 1, level: 'L4', isFolder: false, isLeaf: true },
              { id: 'org-4-1-2-2', name: '邵雪 · 客户成功', manager: '-', members: 1, level: 'L4', isFolder: false, isLeaf: true }
            ]
          }
        ]
      }
    ]
  },
  {
    id: 'org-5',
    name: '人力资源中心',
    manager: '刘林',
    members: 7,
    level: 'L1',
    isFolder: true,
    children: [
      {
        id: 'org-5-1',
        name: 'HRBP 小组',
        manager: '刘林',
        members: 3,
        level: 'L2',
        isFolder: true,
        children: [
          { id: 'org-5-1-1', name: '刘林 · HRBP', manager: '-', members: 1, level: 'L3', isFolder: false, isLeaf: true },
          { id: 'org-5-1-2', name: '唐可 · HRBP', manager: '-', members: 1, level: 'L3', isFolder: false, isLeaf: true }
        ]
      },
      {
        id: 'org-5-2',
        name: '招聘组',
        manager: '赵媛',
        members: 2,
        level: 'L2',
        isFolder: true,
        children: [
          { id: 'org-5-2-1', name: '赵媛 · 招聘', manager: '-', members: 1, level: 'L3', isFolder: false, isLeaf: true },
          { id: 'org-5-2-2', name: '余墨 · 招聘', manager: '-', members: 1, level: 'L3', isFolder: false, isLeaf: true }
        ]
      }
    ]
  },
  {
    id: 'org-6',
    name: '财务与法务',
    manager: '胡娜',
    members: 6,
    level: 'L1',
    isFolder: true,
    children: [
      { id: 'org-6-1', name: '胡娜 · 财务经理', manager: '-', members: 1, level: 'L2', isFolder: false, isLeaf: true },
      { id: 'org-6-2', name: '关然 · 法务', manager: '-', members: 1, level: 'L2', isFolder: false, isLeaf: true }
    ]
  }
]);

const orgColumns = computed(() => [
  { key: 'name', label: t('common.name'), minWidth: 200 },
  { key: 'manager', label: t('common.manager'), minWidth: 140 },
  { key: 'members', label: t('common.members'), minWidth: 120, align: 'center' },
  { key: 'level', label: t('common.level'), minWidth: 120, align: 'center' },
  { key: 'actions', label: t('common.actions'), minWidth: 140, align: 'center' }
]);

const currentPage = ref(1);
const pageSize = ref(6);
const pageSizes = [6, 10, 20];
const paginationLayout = 'total, prev, pager, next';

const totalOrgs = computed(() => orgTreeData.value.length);

const pagedOrgNodes = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value;
  return orgTreeData.value.slice(start, start + pageSize.value);
});

const handlePageChange = (page) => {
  currentPage.value = page;
};

const handleSizeChange = (size) => {
  pageSize.value = size;
  currentPage.value = 1;
};

onMounted(() => {
  initLanguage();
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

.dark-mode .manage-section {
  background: #161b22;
  border-color: #2b2f36;
}

.dark-mode .section-header {
  background: #161b22;
  border-bottom-color: #2b2f36;
}
</style>
