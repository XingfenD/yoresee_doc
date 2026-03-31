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
    <TitleBar :show-back="true" :compact="true" :back-text="t('common.back')" @back="goBackToDocument">
      <template #title>
        {{ t('document.attachments.title') }}
      </template>
      <template #actions>
        <el-upload
          :auto-upload="false"
          :show-file-list="false"
          :before-upload="() => false"
          :on-change="handleChooseFile"
        >
          <el-button type="primary" :loading="uploading" :disabled="!docId">
            {{ t('document.attachments.upload') }}
          </el-button>
        </el-upload>
      </template>
    </TitleBar>

    <div class="detail-content" v-loading="loading">
      <div class="detail-columns">
        <CommonList
          class="detail-list-panel"
          :rows="attachments"
          :columns="columns"
          :row-key="'external_id'"
          :is-dark="isDarkMode"
          :empty-text="t('document.attachments.empty')"
          :show-pagination="false"
          :show-title-bar="true"
          :title="t('document.attachments.title')"
        >
          <template #cell-size="{ value }">
            {{ formatSize(value) }}
          </template>
          <template #cell-actions="{ row }">
            <el-button link type="primary" @click="openAttachment(row)">
              {{ t('document.attachments.open') }}
            </el-button>
            <el-button link type="danger" @click="removeAttachment(row)">
              {{ t('document.delete') }}
            </el-button>
          </template>
        </CommonList>
      </div>
    </div>
  </PageLayout>
</template>

<script setup>
import { computed, onMounted, ref } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { useI18n } from 'vue-i18n';
import { ElMessage, ElMessageBox } from 'element-plus';
import PageLayout from '@/components/PageLayout.vue';
import TitleBar from '@/components/TitleBar.vue';
import CommonList from '@/components/CommonList.vue';
import { useWorkspaceShell } from '@/composables/useWorkspaceShell';
import { useUserStore } from '@/store/user';
import {
  uploadDocumentAttachment,
  listDocumentAttachments,
  deleteDocumentAttachment
} from '@/services/api';

const route = useRoute();
const router = useRouter();
const userStore = useUserStore();
const { locale, t } = useI18n();

const kbId = computed(() => route.params.kbId || 'personal');
const docId = computed(() => route.params.docId || '');

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

const attachments = ref([]);
const uploading = ref(false);
const loading = ref(false);

const columns = computed(() => ([
  { key: 'name', label: t('document.attachments.name'), minWidth: '220px' },
  { key: 'mime_type', label: t('document.attachments.type'), minWidth: '140px' },
  { key: 'size', label: t('document.attachments.size'), width: '120px' },
  { key: 'created_at', label: t('document.attachments.createdAt'), minWidth: '180px' },
  { key: 'actions', label: t('document.attachments.actions'), width: '140px', align: 'center' }
]));

const fetchAttachments = async () => {
  if (!docId.value) {
    attachments.value = [];
    return;
  }
  loading.value = true;
  try {
    const response = await listDocumentAttachments(docId.value);
    attachments.value = response.attachments || [];
  } catch (error) {
    attachments.value = [];
    ElMessage.error(t('document.attachments.loadFailed'));
  } finally {
    loading.value = false;
  }
};

const handleChooseFile = async (uploadFile) => {
  const rawFile = uploadFile?.raw;
  if (!rawFile || !docId.value) return;

  uploading.value = true;
  try {
    const fileBytes = new Uint8Array(await rawFile.arrayBuffer());
    await uploadDocumentAttachment({
      document_external_id: docId.value,
      file_name: rawFile.name,
      content_type: rawFile.type || undefined,
      file_bytes: fileBytes
    });
    ElMessage.success(t('document.attachments.uploadSuccess'));
    await fetchAttachments();
  } catch (error) {
    ElMessage.error(t('document.attachments.uploadFailed'));
  } finally {
    uploading.value = false;
  }
};

const removeAttachment = async (attachment) => {
  try {
    await ElMessageBox.confirm(
      t('document.attachments.deleteConfirm'),
      t('document.delete'),
      {
        confirmButtonText: t('button.confirm'),
        cancelButtonText: t('button.cancel'),
        type: 'warning'
      }
    );
    await deleteDocumentAttachment(docId.value, attachment.external_id);
    ElMessage.success(t('document.attachments.deleteSuccess'));
    await fetchAttachments();
  } catch (error) {
    if (error === 'cancel' || error === 'close') return;
    ElMessage.error(t('document.attachments.deleteFailed'));
  }
};

const openAttachment = (attachment) => {
  const attachmentId = attachment?.external_id;
  if (!attachmentId) {
    ElMessage.warning(t('document.attachments.urlMissing'));
    return;
  }

  const path = kbId.value === 'personal'
    ? `/mydocument/${docId.value}/attachment/${attachmentId}`
    : `/knowledge-base/${kbId.value}/document/${docId.value}/attachment/${attachmentId}`;
  router.push(path);
};

const formatSize = (bytes) => {
  const value = Number(bytes || 0);
  if (!Number.isFinite(value) || value <= 0) return '0 B';
  if (value < 1024) return `${value} B`;
  if (value < 1024 * 1024) return `${(value / 1024).toFixed(1)} KB`;
  return `${(value / (1024 * 1024)).toFixed(1)} MB`;
};

const goBackToDocument = () => {
  if (kbId.value === 'personal') {
    router.push(`/mydocument/${docId.value}`);
    return;
  }
  router.push(`/knowledge-base/${kbId.value}/document/${docId.value}`);
};

onMounted(async () => {
  await initLanguage();
  await fetchSystemInfo();
  await fetchAttachments();
});
</script>

<style scoped>
.detail-content {
  display: flex;
  flex-direction: column;
  flex: 1;
  min-height: 0;
  gap: var(--spacing-lg);
}

.detail-columns {
  display: flex;
  flex: 1;
  min-height: 0;
  gap: var(--spacing-lg);
  align-items: stretch;
}

.detail-list-panel {
  flex: 1;
  min-width: 0;
  min-height: 0;
}

@media (max-width: 1200px) {
  .detail-columns {
    flex-direction: column;
  }
}
</style>
