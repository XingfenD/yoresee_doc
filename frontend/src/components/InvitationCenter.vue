<template>
  <el-tabs v-model="activeTab" class="common-tabs">
    <el-tab-pane :label="tabListLabel" name="list">
      <CommonList
        :rows="inviteList"
        :columns="inviteColumns"
        :is-dark="isDarkMode"
        row-key="code"
        :empty-text="t('message.empty')"
        :show-pagination="isSystemMode"
        :total="isSystemMode ? inviteTotal : inviteList.length"
        v-model:current-page="invitePage"
        v-model:page-size="invitePageSize"
        :page-sizes="[10, 20, 50]"
        @page-change="handleInvitePageChange"
        :show-search="isSystemMode"
        v-model:search-query="inviteKeyword"
        :search-placeholder="t('common.search')"
        @search="handleInviteSearch"
        :show-title-bar="isSystemMode"
        :title="tabListLabel"
      >
        <template #cell-status="{ value }">
          <AppTag :type="inviteStatusType(value)" size="small">
            {{ inviteStatusLabel(value) }}
          </AppTag>
        </template>
        <template #cell-usage="{ row }">
          {{ row.used }}/{{ row.max === null ? '-' : row.max }}
        </template>
        <template #cell-code="{ row }">
          <el-tooltip :content="row.note || t('user.invite.notePlaceholder')" placement="top">
            <span class="invite-code" @click="copyInviteCode(row.code)">{{ row.code }}</span>
          </el-tooltip>
        </template>
        <template #cell-actions="{ row }">
          <el-button size="small" text type="primary" @click="handlePauseInvite(row)">
            {{ row.disabled ? t('user.invite.resume') : t('user.invite.pause') }}
          </el-button>
          <el-button size="small" text type="danger" @click="handleDeleteInvite(row)">
            {{ t('user.invite.delete') }}
          </el-button>
        </template>
      </CommonList>
    </el-tab-pane>

    <el-tab-pane :label="tabRecordsLabel" name="records">
      <CommonList
        :rows="inviteRecords"
        :columns="recordColumns"
        :is-dark="isDarkMode"
        row-key="row_key"
        :empty-text="recordsEmptyText"
        :show-pagination="isSystemMode"
        :total="isSystemMode ? recordTotal : inviteRecords.length"
        v-model:current-page="recordPage"
        v-model:page-size="recordPageSize"
        :page-sizes="[10, 20, 50]"
        @page-change="handleRecordPageChange"
        :show-search="isSystemMode"
        v-model:search-query="recordKeyword"
        :search-placeholder="t('common.search')"
        @search="handleRecordSearch"
        :show-title-bar="isSystemMode"
        :title="tabRecordsLabel"
      >
        <template #cell-status="{ value }">
          <AppTag :type="value === 'success' ? 'success' : 'warning'" size="small">
            {{ value === 'success' ? recordsSuccessLabel : recordsFailedLabel }}
          </AppTag>
        </template>
      </CommonList>
    </el-tab-pane>
  </el-tabs>

  <InviteCreateDialog v-model="showCreateDialog" @submit="handleCreateInvite" />
</template>

<script setup>
import { computed, onMounted, ref, watch } from 'vue';
import { useI18n } from 'vue-i18n';
import { ElMessageBox } from 'element-plus';
import CommonList from '@/components/CommonList.vue';
import InviteCreateDialog from '@/components/InviteCreateDialog.vue';
import AppTag from '@/components/AppTag.vue';
import { useServerTable } from '@/composables/useServerTable';
import { isActionCancelled, useApiAction } from '@/composables/useApiAction';
import { listInvitations, listInvitationRecords, createInvitation, updateInvitation, deleteInvitation } from '@/services/api';

const props = defineProps({
  mode: {
    type: String,
    default: 'user',
    validator: (value) => ['user', 'system'].includes(value)
  },
  isDarkMode: {
    type: Boolean,
    default: false
  }
});

const { t } = useI18n();
const { runApi } = useApiAction({ t });

const isSystemMode = computed(() => props.mode === 'system');
const activeTab = ref('list');
const inviteLoaded = ref(false);
const recordsLoaded = ref(false);
const showCreateDialog = ref(false);

const tabListLabel = computed(() => t(isSystemMode.value ? 'system.invite.tabs.list' : 'user.invite.tabs.list'));
const tabRecordsLabel = computed(() => t(isSystemMode.value ? 'system.invite.tabs.records' : 'user.invite.tabs.records'));
const recordsEmptyText = computed(() => t(isSystemMode.value ? 'system.invite.records.empty' : 'user.invite.records.empty'));
const recordsSuccessLabel = computed(() => t(isSystemMode.value ? 'system.invite.records.success' : 'user.invite.records.success'));
const recordsFailedLabel = computed(() => t(isSystemMode.value ? 'system.invite.records.failed' : 'user.invite.records.failed'));

const inviteColumns = computed(() => {
  const baseColumns = [
    { key: 'code', label: t('user.invite.code'), minWidth: 180 },
    { key: 'status', label: t('user.invite.status'), minWidth: 120, align: 'center' },
    { key: 'usage', label: t('user.invite.usage'), minWidth: 110, align: 'center' },
    { key: 'created_at', label: t('user.invite.createdAt'), minWidth: 160 }
  ];
  if (isSystemMode.value) {
    baseColumns.push({ key: 'created_by', label: t('user.invite.createdBy'), minWidth: 140 });
  }
  baseColumns.push(
    { key: 'expires_at', label: t('user.invite.expiresAt'), minWidth: 160 },
    { key: 'actions', label: t('user.invite.actions'), minWidth: 160, align: 'center', headerAlign: 'center' }
  );
  return baseColumns;
});

const recordColumns = computed(() => {
  if (isSystemMode.value) {
    return [
      { key: 'code', label: t('system.invite.records.code'), minWidth: 180 },
      { key: 'used_by', label: t('system.invite.records.usedBy'), minWidth: 160 },
      { key: 'used_at', label: t('system.invite.records.usedAt'), minWidth: 180 },
      { key: 'status', label: t('system.invite.records.result'), minWidth: 120, align: 'center' }
    ];
  }

  return [
    { key: 'code', label: t('user.invite.records.code'), minWidth: 180 },
    { key: 'used_by', label: t('user.invite.records.usedBy'), minWidth: 160 },
    { key: 'used_at', label: t('user.invite.records.usedAt'), minWidth: 180 },
    { key: 'status', label: t('user.invite.records.result'), minWidth: 120, align: 'center' }
  ];
});

const inviteStatusType = (status) => {
  if (status === 'active') return 'success';
  if (status === 'expired') return 'info';
  return 'warning';
};

const inviteStatusLabel = (status) => {
  if (status === 'active') return t('user.invite.active');
  if (status === 'expired') return t('user.invite.expired');
  return t('user.invite.disabled');
};

const formatDateYYYYMMDD = (date) => {
  const pad = (num) => `${num}`.padStart(2, '0');
  return `${date.getFullYear()}-${pad(date.getMonth() + 1)}-${pad(date.getDate())}`;
};

const toRfc3339EndOfDay = (dateStr) => {
  if (!dateStr) return '';
  if (dateStr.includes('T')) {
    return dateStr;
  }
  const [year, month, day] = dateStr.split('-').map((value) => Number(value));
  if (!year || !month || !day) return '';
  const local = new Date(year, month - 1, day, 23, 59, 59);
  return local.toISOString();
};

const resolveInviteStatus = (invite) => {
  if (invite.disabled) return 'disabled';
  if (invite.expires_at) {
    const expireTs = Date.parse(invite.expires_at);
    if (!Number.isNaN(expireTs) && expireTs < Date.now()) return 'expired';
  }
  if (typeof invite.max_used_cnt === 'number' && typeof invite.used_cnt === 'number' && invite.used_cnt >= invite.max_used_cnt) {
    return 'expired';
  }
  return 'active';
};

const toNumber = (value, fallback = null) => {
  if (value === null || value === undefined) return fallback;
  const num = Number(value);
  return Number.isFinite(num) ? num : fallback;
};

const mapInviteRow = (invite) => {
  const row = {
    code: invite.code,
    created_at: invite.created_at || '-',
    expires_at: invite.expires_at || '-',
    status: resolveInviteStatus(invite),
    max: toNumber(invite.max_used_cnt, null),
    used: toNumber(invite.used_cnt, 0),
    disabled: invite.disabled,
    note: invite.note || ''
  };

  if (isSystemMode.value) {
    row.created_by = invite.created_by_name || invite.created_by_external_id || t('common.unknown');
  }

  return row;
};

const mapRecordRow = (record) => ({
  row_key: record.row_key || `${record.code || ''}_${record.used_at || ''}_${record.status || ''}`,
  code: record.code,
  used_by: record.used_by || '-',
  used_at: record.used_at || '-',
  status: record.status || 'failed'
});

const {
  rows: inviteList,
  page: invitePage,
  pageSize: invitePageSize,
  total: inviteTotal,
  keyword: inviteKeyword,
  load: loadInvitations,
  handlePageChange: changeInvitePage,
  handleSearch: triggerInviteSearch
} = useServerTable({
  initialPageSize: 10,
  fetcher: ({ page, page_size, keyword }) =>
    listInvitations(
      isSystemMode.value
        ? {
            page,
            page_size,
            keyword
          }
        : {
            only_mine: true,
            page: 1,
            page_size: 50
          }
    ),
  mapRows: (resp) => (resp.invitations || []).map(mapInviteRow),
  mapTotal: (resp, rows) => (isSystemMode.value ? resp.total : rows.length),
  onError: (err) => {
    console.error('listInvitations failed', err);
  }
});

const {
  rows: inviteRecords,
  page: recordPage,
  pageSize: recordPageSize,
  total: recordTotal,
  keyword: recordKeyword,
  load: loadInvitationRecords,
  handlePageChange: changeRecordPage,
  handleSearch: triggerRecordSearch
} = useServerTable({
  initialPageSize: 10,
  fetcher: ({ page, page_size, keyword }) =>
    listInvitationRecords(
      isSystemMode.value
        ? {
            page,
            page_size,
            keyword
          }
        : {
            only_mine: true,
            page: 1,
            page_size: 100
          }
    ),
  mapRows: (resp) => {
    const records = [...(resp.records || [])];
    if (!isSystemMode.value) {
      records.sort((a, b) => (b.used_at || '').localeCompare(a.used_at || ''));
    }
    return records.map(mapRecordRow);
  },
  mapTotal: (resp, rows) => (isSystemMode.value ? resp.total : rows.length),
  onError: (err) => {
    console.error('listInvitationRecords failed', err);
  }
});

const fetchInvitations = async () => {
  const resp = await loadInvitations();
  if (resp) {
    inviteLoaded.value = true;
  }
};

const fetchInvitationRecords = async () => {
  const resp = await loadInvitationRecords();
  if (resp) {
    recordsLoaded.value = true;
  }
};

const handleCreateInvite = async (payload) => {
  let expiresAt = payload.expires_at;
  if (payload.expire_type === 'days' && payload.expire_days) {
    const target = new Date();
    target.setDate(target.getDate() + Number(payload.expire_days));
    expiresAt = formatDateYYYYMMDD(target);
  }

  await runApi(
    async () => {
      await createInvitation({
        expires_at: expiresAt ? toRfc3339EndOfDay(expiresAt) : undefined,
        max_used_cnt: payload.limit_enabled ? payload.max_usage : undefined,
        note: payload.note
      });
      await fetchInvitations();
    },
    {
      context: 'createInvitation',
      successMessage: t('message.success')
    }
  );
};

const handlePauseInvite = async (row) => {
  if (!row?.code) return;

  await runApi(
    async () => {
      await updateInvitation({
        code: row.code,
        disabled: !row.disabled
      });
      await fetchInvitations();
    },
    {
      context: 'updateInvitation',
      successMessage: t('message.success')
    }
  );
};

const handleDeleteInvite = async (row) => {
  if (!row?.code) return;

  await runApi(
    async () => {
      await ElMessageBox.confirm(t('message.confirmDelete'), t('user.invite.delete'), {
        confirmButtonText: t('button.confirm'),
        cancelButtonText: t('button.cancel'),
        type: 'warning'
      });
      await deleteInvitation(row.code);
      await fetchInvitations();
    },
    {
      context: 'deleteInvitation',
      successMessage: t('message.deleteSuccess'),
      ignoreError: isActionCancelled
    }
  );
};

const handleInvitePageChange = async (page) => {
  if (!isSystemMode.value) return;
  await changeInvitePage(page);
};

const handleRecordPageChange = async (page) => {
  if (!isSystemMode.value) return;
  await changeRecordPage(page);
};

const handleInviteSearch = () => {
  if (!isSystemMode.value) return;
  triggerInviteSearch();
};

const handleRecordSearch = () => {
  if (!isSystemMode.value) return;
  triggerRecordSearch();
};

const copyInviteCode = async (code) => {
  if (!code) return;

  await runApi(
    async () => {
      await navigator.clipboard.writeText(code);
    },
    {
      context: 'copy invite code',
      successMessage: t('common.copySuccess'),
      errorMessage: t('common.copyFailed')
    }
  );
};

const openCreateDialog = () => {
  showCreateDialog.value = true;
};

watch(activeTab, async (tab) => {
  if (tab === 'list') {
    if (!inviteLoaded.value || !isSystemMode.value) {
      await fetchInvitations();
    }
    return;
  }

  if (tab === 'records') {
    if (!recordsLoaded.value || !isSystemMode.value) {
      await fetchInvitationRecords();
    }
  }
});

onMounted(async () => {
  await fetchInvitations();
});

defineExpose({
  openCreateDialog
});
</script>

<style scoped>
.invite-code {
  cursor: pointer;
  color: var(--primary-color);
}
</style>
