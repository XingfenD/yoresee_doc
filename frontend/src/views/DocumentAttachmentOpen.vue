<template>
  <PageLayout
    :system-name="systemName"
    :current-language="currentLanguage"
    :is-dark-mode="isDarkMode"
    :user-avatar="userAvatar"
    :username="userInfo?.username || t('common.user')"
    :active-menu="activeMenu"
    :title="''"
    content-padding="xl"
    @change-language="handleLanguageChange"
    @toggle-theme="toggleTheme"
    @logout="handleLogout"
    @menu-select="handleMenuSelect"
  >
    <TitleBar :show-back="true" :compact="true" :back-text="t('common.back')" @back="goBackToAttachments">
      <template #title>
        {{ t('document.attachments.previewTitle') }}
      </template>
    </TitleBar>

    <div class="open-wrapper" v-loading="loading">
      <el-empty v-if="errorText" :description="errorText" />
      <div v-else class="preview-panel">
        <div class="preview-meta">
          <div class="preview-name">{{ attachmentName || '-' }}</div>
          <div class="preview-type">{{ attachmentType || '-' }}</div>
        </div>

        <div class="preview-body">
          <el-image
            v-if="previewMode === 'image'"
            class="preview-image"
            :src="previewUrl"
            :preview-src-list="[previewUrl]"
            :initial-index="0"
            fit="contain"
            preview-teleported
          />
          <video v-else-if="previewMode === 'video'" class="preview-media" controls :src="previewUrl"></video>
          <audio v-else-if="previewMode === 'audio'" class="preview-audio" controls :src="previewUrl"></audio>
          <pre v-else-if="previewMode === 'text'" class="preview-text">{{ textContent }}</pre>
          <iframe v-else-if="previewMode === 'iframe'" class="preview-iframe" :src="previewUrl"></iframe>
          <el-empty v-else :description="t('document.attachments.previewUnsupported')">
            <el-button type="primary" @click="openInCurrentPage">
              {{ t('document.attachments.download') }}
            </el-button>
          </el-empty>
        </div>
      </div>
    </div>
  </PageLayout>
</template>

<script setup>
import { computed, onMounted, ref } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { useI18n } from 'vue-i18n';
import { ElMessage } from 'element-plus';
import PageLayout from '@/components/PageLayout.vue';
import TitleBar from '@/components/TitleBar.vue';
import { useWorkspaceShell } from '@/composables/useWorkspaceShell';
import { useUserStore } from '@/store/user';
import { listDocumentAttachments } from '@/services/api';
import { resolveAttachmentUrl } from '@/utils/attachmentUrl';

const route = useRoute();
const router = useRouter();
const userStore = useUserStore();
const { locale, t } = useI18n();

const kbId = computed(() => route.params.kbId || 'personal');
const docId = computed(() => route.params.docId || '');
const attachmentId = computed(() => route.params.attachmentId || '');

const {
  systemName,
  activeMenu,
  isDarkMode,
  userInfo,
  userAvatar,
  currentLanguage,
  initLanguage,
  handleLanguageChange,
  toggleTheme,
  handleLogout,
  handleMenuSelect,
  fetchSystemInfo
} = useWorkspaceShell({
  locale,
  router,
  userStore,
  defaultActiveMenu: kbId.value === 'personal' ? 'documents' : 'knowledge-base'
});

const loading = ref(false);
const errorText = ref('');
const previewUrl = ref('');
const attachmentType = ref('');
const attachmentName = ref('');
const textContent = ref('');
const previewMode = ref('');

const goBackToAttachments = () => {
  if (kbId.value === 'personal') {
    router.push(`/mydocument/${docId.value}/attachments`);
    return;
  }
  router.push(`/knowledge-base/${kbId.value}/document/${docId.value}/attachments`);
};

const detectPreviewMode = (mimeType, fileName) => {
  const type = `${mimeType || ''}`.toLowerCase();
  const name = `${fileName || ''}`.toLowerCase();

  if (type.startsWith('image/')) return 'image';
  if (type.startsWith('video/')) return 'video';
  if (type.startsWith('audio/')) return 'audio';
  if (type === 'application/pdf' || name.endsWith('.pdf')) return 'iframe';
  if (
    type.startsWith('text/') ||
    type.includes('json') ||
    type.includes('xml') ||
    name.endsWith('.md') ||
    name.endsWith('.go') ||
    name.endsWith('.js') ||
    name.endsWith('.ts') ||
    name.endsWith('.java') ||
    name.endsWith('.py')
  ) {
    return 'text';
  }
  return 'iframe';
};

const openInCurrentPage = () => {
  if (!previewUrl.value) return;
  window.location.assign(previewUrl.value);
};

const loadAttachment = async () => {
  if (!docId.value || !attachmentId.value) {
    errorText.value = t('document.attachments.urlMissing');
    return;
  }

  loading.value = true;
  try {
    const response = await listDocumentAttachments(docId.value);
    const current = (response.attachments || []).find((item) => item?.external_id === attachmentId.value);
    const targetUrl = resolveAttachmentUrl(current?.url);

    if (!targetUrl) {
      errorText.value = t('document.attachments.urlMissing');
      return;
    }
    previewUrl.value = targetUrl;
    attachmentType.value = current?.mime_type || '';
    attachmentName.value = current?.name || '';
    previewMode.value = detectPreviewMode(current?.mime_type, current?.name);

    if (previewMode.value === 'text') {
      try {
        const resp = await fetch(targetUrl, { method: 'GET', credentials: 'include' });
        if (!resp.ok) {
          throw new Error('fetch attachment text failed');
        }
        textContent.value = await resp.text();
      } catch (error) {
        previewMode.value = 'iframe';
        ElMessage.warning(t('document.attachments.contentLoadFailed'));
      }
    }
  } catch (error) {
    errorText.value = t('document.attachments.loadFailed');
    ElMessage.error(t('document.attachments.loadFailed'));
  } finally {
    loading.value = false;
  }
};

onMounted(async () => {
  await initLanguage();
  await fetchSystemInfo();
  await loadAttachment();
});
</script>

<style scoped>
.open-wrapper {
  min-height: 520px;
}

.preview-panel {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-md);
  height: calc(100vh - 220px);
}

.preview-meta {
  display: flex;
  justify-content: space-between;
  gap: var(--spacing-md);
  padding: var(--spacing-md);
  border: 1px solid var(--border-color);
  border-radius: var(--border-radius-md);
  background: var(--bg-white);
}

.preview-name {
  font-weight: 600;
  color: var(--text-dark);
}

.preview-type {
  color: var(--text-secondary);
  font-size: 13px;
}

.preview-body {
  flex: 1;
  min-height: 0;
  border: 1px solid var(--border-color);
  border-radius: var(--border-radius-md);
  overflow: hidden;
  background: var(--bg-white);
}

.preview-image,
.preview-media,
.preview-iframe {
  width: 100%;
  height: 100%;
  border: 0;
}

.preview-audio {
  width: 100%;
  padding: var(--spacing-lg);
}

.preview-text {
  margin: 0;
  height: 100%;
  overflow: auto;
  padding: var(--spacing-lg);
  font-size: 13px;
  line-height: 1.6;
  background: #0f172a;
  color: #e2e8f0;
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, "Liberation Mono", "Courier New", monospace;
}

:global(.dark-mode) .preview-meta,
:global(.dark-mode) .preview-body {
  background: var(--bg-secondary);
}
</style>
