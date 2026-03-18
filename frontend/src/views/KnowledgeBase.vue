<template>
  <PageLayout
    :system-name="systemName"
    :current-language="currentLanguage"
    :is-dark-mode="isDarkMode"
    :user-avatar="userAvatar"
    :username="userInfo?.username || t('common.unknown')"
    :active-menu="activeMenu"
    :title="t('knowledgeBase.title')"
    @change-language="handleLanguageChange"
    @toggle-theme="toggleTheme"
    @logout="handleLogout"
    @menu-select="handleMenuSelect"
  >
    <template #actions>
      <el-button class="page-action-btn" type="primary" size="small" @click="createKnowledgeBase">
        {{ t("knowledgeBase.createNew") }}
      </el-button>
    </template>

    <div class="knowledge-base-vertical-layout">
      <el-tabs v-model="activeTab" class="knowledge-base-tabs">
        <el-tab-pane :label="t('knowledgeBase.my')" name="my">
          <KnowledgeBaseListSection
            :title="t('knowledgeBase.my')"
            :items="myKnowledgeBases"
            :empty-text="t('knowledgeBase.noRecent')"
            :tag-mapper="myTagMapper"
            :fallback-description="t('knowledgeBase.noDescription')"
            :meta-mapper="myMetaMapper"
            :show-load-more="myHasMore"
            :loading="myLoading"
            :load-more-label="t('common.loadMore')"
            :loading-label="t('common.loading')"
            :action-label="t('common.open')"
            @open="viewKnowledgeBase"
            @load-more="loadMoreMyKnowledgeBases"
          />
        </el-tab-pane>

        <el-tab-pane :label="t('knowledgeBase.recent')" name="recent">
          <KnowledgeBaseListSection
            :title="t('knowledgeBase.recent')"
            :items="recentKnowledgeBases"
            :empty-text="t('knowledgeBase.noRecent')"
            :tag-mapper="recentTagMapper"
            :fallback-description="t('knowledgeBase.noDescription')"
            :meta-mapper="null"
            :show-load-more="false"
            :action-label="t('common.open')"
            @open="accessKnowledgeBase"
          />
        </el-tab-pane>

        <el-tab-pane :label="t('knowledgeBase.publicList')" name="public">
          <KnowledgeBaseListSection
            :title="t('knowledgeBase.publicList')"
            :items="publicKnowledgeBases"
            :empty-text="t('knowledgeBase.noRecent')"
            :tag-type="'success'"
            :tag-label="t('knowledgeBase.public')"
            :fallback-description="t('knowledgeBase.noDescription')"
            :meta-mapper="publicMetaMapper"
            :show-load-more="publicHasMore"
            :loading="publicLoading"
            :load-more-label="t('common.loadMore')"
            :loading-label="t('common.loading')"
            :action-label="t('common.open')"
            @open="viewKnowledgeBase"
            @load-more="loadMorePublicKnowledgeBases"
          />
        </el-tab-pane>
      </el-tabs>
    </div>
  </PageLayout>
</template>

<script setup>
import { ref, onMounted, computed, watch } from "vue";
import { useRouter } from "vue-router";
import { useUserStore } from "@/store/user";
import { ElMessage } from "element-plus";
import { useI18n } from "vue-i18n";
import PageLayout from "@/components/PageLayout.vue";
import KnowledgeBaseListSection from "@/components/KnowledgeBaseListSection.vue";
import * as api from "@/services/api";
import { House, Collection, Plus } from "@element-plus/icons-vue";

// 国际化
const { locale, t } = useI18n();
const router = useRouter();
const userStore = useUserStore();

// 系统信息
const systemName = ref("Yoresee");

// 导航相关
const activeMenu = ref("knowledge-base");
const activeTab = ref("my");
const isDarkMode = ref(false);
const currentLanguage = computed({
  get: () => locale.value,
  set: (value) => {
    locale.value = value;
    localStorage.setItem("language", value);
  },
});

// 用户信息
const userInfo = computed(() => userStore.userInfo);
const userAvatar = computed(() => userInfo.value?.avatar || 'https://cube.elemecdn.com/0/88/03b0d39583f48206768a7534e55bcpng.png');

// 最近访问的知识库
const recentKnowledgeBases = ref([]);
const recentPage = ref(1);
const recentPageSize = ref(10);
const recentTotal = ref(0);
const recentLoading = ref(false);
const recentLoaded = ref(false);

// 我的知识库
const myKnowledgeBases = ref([]);
const myPage = ref(1);
const myPageSize = ref(10);
const myTotal = ref(0);
const myLoading = ref(false);
const myHasMore = computed(() => myKnowledgeBases.value.length < myTotal.value);
const myTagMapper = (kb) =>
  kb.isPublic ? null : { type: "info", label: t("knowledgeBase.private") };
const myMetaMapper = (kb) => [
  { label: t("knowledgeBase.documentsCount"), value: kb.documents_count || 0 },
  { label: t("knowledgeBase.updatedAt"), value: formatDate(kb.updated_at) },
];

// 公开知识库
const publicKnowledgeBases = ref([]);
const publicPage = ref(1);
const publicPageSize = ref(10);
const publicTotal = ref(0);
const publicLoading = ref(false);
const publicHasMore = computed(
  () => publicKnowledgeBases.value.length < publicTotal.value
);
const publicMetaMapper = (kb) => [
  { label: t("knowledgeBase.documentsCount"), value: kb.documents_count || 0 },
  { label: t("knowledgeBase.owner"), value: kb.creator_name || t("common.unknown") },
];

const recentTagMapper = (kb) =>
  kb.isPublic ? { type: "success", label: t("knowledgeBase.public") } : null;

const myLoaded = ref(false);
const publicLoaded = ref(false);

// 获取最近访问的知识库
const fetchRecentKnowledgeBases = async (page = 1, pageSize = 10) => {
  if (recentLoading.value) return;

  recentLoading.value = true;

  try {
    const params = {
      page: page,
      page_size: pageSize,
    };
    const data = await api.getRecentKnowledgeBases(params);

    if (page === 1) {
      recentKnowledgeBases.value = data.knowledge_bases || [];
    } else {
      recentKnowledgeBases.value.push(...(data.knowledge_bases || []));
    }

    recentTotal.value = data.total || 0;
    recentLoaded.value = true;
  } catch (error) {
    console.error("Failed to fetch recent knowledge bases:", error);
    ElMessage.error(t("knowledgeBase.fetchError"));
  } finally {
    recentLoading.value = false;
  }
};

// 获取我的知识库
const fetchMyKnowledgeBases = async (page = 1, pageSize = 10) => {
  if (myLoading.value) return;

  myLoading.value = true;

  try {
    const params = {
      page: page,
      page_size: pageSize,
      only_mine: true,
    };

    const data = await api.getKnowledgeBases(params);

    if (page === 1) {
      myKnowledgeBases.value = data.knowledge_bases || [];
    } else {
      myKnowledgeBases.value.push(...(data.knowledge_bases || []));
    }

    myTotal.value = data.total || 0;
    myLoaded.value = true;
  } catch (error) {
    console.error("Failed to fetch my knowledge bases:", error);
    ElMessage.error(t("knowledgeBase.fetchError"));
  } finally {
    myLoading.value = false;
  }
};

// 获取公开知识库
const fetchPublicKnowledgeBases = async (page = 1, pageSize = 10) => {
  if (publicLoading.value) return;

  publicLoading.value = true;

  try {
    const params = {
      page: page,
      page_size: pageSize,
      is_public: true, // 获取公开知识库
    };

    const data = await api.getKnowledgeBases(params);

    if (page === 1) {
      publicKnowledgeBases.value = data.knowledge_bases || [];
    } else {
      publicKnowledgeBases.value.push(...(data.knowledge_bases || []));
    }

    publicTotal.value = data.total || 0;
    publicLoaded.value = true;
  } catch (error) {
    console.error("Failed to fetch public knowledge bases:", error);
    ElMessage.error(t("knowledgeBase.fetchError"));
  } finally {
    publicLoading.value = false;
  }
};

// 加载更多我的知识库
const loadMoreMyKnowledgeBases = async () => {
  if (!myHasMore.value || myLoading.value) return;

  myPage.value++;
  await fetchMyKnowledgeBases(myPage.value, myPageSize.value);
};

// 加载更多公开知识库
const loadMorePublicKnowledgeBases = async () => {
  if (!publicHasMore.value || publicLoading.value) return;

  publicPage.value++;
  await fetchPublicKnowledgeBases(publicPage.value, publicPageSize.value);
};

// 创建知识库
const createKnowledgeBase = () => {
  // TODO: 实现创建知识库功能
  ElMessage.info(t("knowledgeBase.createComingSoon"));
};

// 格式化日期
const formatDate = (dateString) => {
  if (!dateString) return t("common.unknown");

  const date = new Date(dateString);
  return date.toLocaleDateString();
};

// 查看知识库详情
const viewKnowledgeBase = (kb) => {
  router.push(`/knowledge-base/${kb.external_id}`);
};

// 访问知识库
const accessKnowledgeBase = (kb) => {
  router.push(`/knowledge-base/${kb.external_id}`);
};

// 处理菜单选择
const handleMenuSelect = (key) => {
  activeMenu.value = key;
};

// 处理语言切换
const handleLanguageChange = (command) => {
  currentLanguage.value = command;
};

// 处理主题切换
const toggleTheme = () => {
  isDarkMode.value = !isDarkMode.value;
  if (isDarkMode.value) {
    document.documentElement.classList.add("dark-mode");
    localStorage.setItem("darkMode", "true");
  } else {
    document.documentElement.classList.remove("dark-mode");
    localStorage.setItem("darkMode", "false");
  }
};

// 初始化主题
const initTheme = () => {
  const savedDarkMode = localStorage.getItem("darkMode");
  if (savedDarkMode === "true") {
    isDarkMode.value = true;
    document.documentElement.classList.add("dark-mode");
  }
};

// 初始化语言
const initLanguage = () => {
  const savedLanguage = localStorage.getItem("language");
  if (savedLanguage) {
    currentLanguage.value = savedLanguage;
  }
};

// 获取系统信息
const fetchSystemInfo = async () => {
  try {
    const info = await userStore.fetchSystemInfo();
    systemName.value = info.system_name;
  } catch (err) {
    console.error("获取系统信息失败:", err);
  }
};

// 登出处理
const handleLogout = () => {
  userStore.logout();
  router.push("/login");
};

onMounted(async () => {
  // 获取系统信息
  await fetchSystemInfo();

  // 初始化主题和语言
  initTheme();
  initLanguage();

  // 默认加载我的知识库
  await fetchMyKnowledgeBases(myPage.value, myPageSize.value);
});

watch(activeTab, async (tab) => {
  if (tab === "my" && !myLoaded.value) {
    await fetchMyKnowledgeBases(myPage.value, myPageSize.value);
  }
  if (tab === "recent" && !recentLoaded.value) {
    await fetchRecentKnowledgeBases(recentPage.value, recentPageSize.value);
  }
  if (tab === "public" && !publicLoaded.value) {
    await fetchPublicKnowledgeBases(publicPage.value, publicPageSize.value);
  }
});
</script>

<style scoped>
/* 垂直布局 */
.knowledge-base-vertical-layout {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-lg);
  height: auto;
}

.knowledge-base-tabs :deep(.el-tabs__header) {
  margin: 0 0 var(--spacing-md);
  border-bottom: 1px solid var(--border-color);
}

.knowledge-base-tabs :deep(.el-tabs__nav-wrap) {
  padding: 0 var(--spacing-sm);
}

.knowledge-base-tabs :deep(.el-tabs__item) {
  color: var(--text-medium);
  font-weight: 500;
}

.knowledge-base-tabs :deep(.el-tabs__item.is-active) {
  color: var(--primary-color);
}

.knowledge-base-tabs :deep(.el-tabs__active-bar) {
  background-color: var(--primary-color);
}

/* 深色模式支持 */
.dark-mode .search-input :deep(.el-input__wrapper) {
  background-color: var(--input-bg);
  border-color: var(--input-border);
  color: var(--input-text);
}

.dark-mode .search-input :deep(.el-input__inner) {
  background-color: var(--input-bg);
  border-color: var(--input-border);
  color: var(--input-text);
}

.dark-mode .filter-select :deep(.el-input__wrapper) {
  background-color: var(--select-bg);
  border-color: var(--select-border);
  color: var(--select-text);
}

.dark-mode .filter-select :deep(.el-input__inner) {
  background-color: var(--select-bg);
  border-color: var(--select-border);
  color: var(--select-text);
}

.dark-mode .filter-select :deep(.el-select__dropdown) {
  background-color: var(--select-bg);
  border-color: var(--select-border);
}

.dark-mode .filter-select :deep(.el-popper) {
  background-color: var(--select-bg);
  border-color: var(--select-border);
  color: var(--select-text);
}

.dark-mode .filter-select :deep(.el-select-dropdown__item) {
  background-color: var(--select-option-bg);
  color: var(--select-text);
}

.dark-mode .filter-select :deep(.el-select-dropdown__item:hover) {
  background-color: var(--select-option-hover);
}

.dark-mode .el-pagination :deep(.el-pager li) {
  background-color: var(--bg-white);
  color: var(--text-primary);
  border-color: var(--border-color);
}

.dark-mode .el-pagination :deep(button) {
  background-color: var(--bg-white);
  color: var(--text-primary);
  border-color: var(--border-color);
}

.dark-mode .el-pagination.is-background :deep(.btn-next),
.dark-mode .el-pagination.is-background :deep(.btn-prev),
.dark-mode .el-pagination.is-background :deep(.el-pager li) {
  background-color: var(--bg-medium);
  color: var(--text-primary);
}

.dark-mode .el-pagination.is-background :deep(.el-pager li:not(.is-disabled):hover) {
  color: var(--primary-color);
}

.dark-mode .el-pagination.is-background :deep(.el-pager li.is-active) {
  background-color: var(--primary-color);
  color: white;
}

.dark-mode .el-button {
  border-color: var(--border-color);
  color: var(--text-primary);
}

.dark-mode .el-button--primary {
  background-color: var(--primary-color);
  border-color: var(--primary-color);
  color: white;
}

.dark-mode .el-button--primary:hover {
  background-color: var(--primary-light);
  border-color: var(--primary-light);
  color: white;
}

.dark-mode .column {
  background-color: var(--bg-medium);
  border: 1px solid var(--border-color);
}

.dark-mode .column-header {
  background-color: var(--bg-medium);
  border-bottom: 1px solid var(--border-color);
}

</style>
