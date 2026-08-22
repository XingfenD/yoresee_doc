import { computed, onMounted, ref, watch } from 'vue';
import { useI18n } from 'vue-i18n';
import { ElMessageBox } from 'element-plus';
import { useServerTable } from '@/composables/list/useServerTable';
import { isActionCancelled, useApiAction } from '@/composables/actions/useApiAction';
import {
  listInvitations,
  listInvitationRecords,
  createInvitation,
  updateInvitation,
  deleteInvitation
} from '@/services/api';
import {
  inviteStatusType,
  inviteStatusLabel as rawInviteStatusLabel,
  formatDateYYYYMMDD,
  toRfc3339EndOfDay,
  mapInviteRow,
  mapRecordRow
} from './invitationData';

export function useInvitationCenter(props) {
  const { t } = useI18n();
  const { runApi } = useApiAction({ t });

  const isSystemMode = computed(() => props.mode === 'system');
  const activeTab = ref('list');
  const inviteLoaded = ref(false);
  const recordsLoaded = ref(false);
  const showCreateDialog = ref(false);

  const tabListLabel = computed(() =>
    t(isSystemMode.value ? 'system.invite.tabs.list' : 'user.invite.tabs.list')
  );
  const tabRecordsLabel = computed(() =>
    t(isSystemMode.value ? 'system.invite.tabs.records' : 'user.invite.tabs.records')
  );
  const recordsEmptyText = computed(() =>
    t(isSystemMode.value ? 'system.invite.records.empty' : 'user.invite.records.empty')
  );
  const recordsSuccessLabel = computed(() =>
    t(isSystemMode.value ? 'system.invite.records.success' : 'user.invite.records.success')
  );
  const recordsFailedLabel = computed(() =>
    t(isSystemMode.value ? 'system.invite.records.failed' : 'user.invite.records.failed')
  );

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
      {
        key: 'actions',
        label: t('user.invite.actions'),
        minWidth: 160,
        align: 'center',
        headerAlign: 'center'
      }
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

  const inviteStatusLabel = (status) => rawInviteStatusLabel(status, t);

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
          ? { page, page_size: page_size, keyword }
          : { only_mine: true, page: 1, page_size: 50 }
      ),
    mapRows: (resp) =>
      (resp.invitations || []).map((invite) => mapInviteRow(invite, { isSystemMode: isSystemMode.value, t })),
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
          ? { page, page_size: page_size, keyword }
          : { only_mine: true, page: 1, page_size: 100 }
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

  return {
    isSystemMode,
    activeTab,
    inviteList,
    inviteColumns,
    tabListLabel,
    inviteTotal,
    invitePage,
    invitePageSize,
    handleInvitePageChange,
    inviteKeyword,
    handleInviteSearch,
    tabRecordsLabel,
    inviteRecords,
    recordColumns,
    recordsEmptyText,
    recordTotal,
    recordPage,
    recordPageSize,
    handleRecordPageChange,
    recordKeyword,
    handleRecordSearch,
    inviteStatusType,
    inviteStatusLabel,
    recordsSuccessLabel,
    recordsFailedLabel,
    copyInviteCode,
    handlePauseInvite,
    handleDeleteInvite,
    showCreateDialog,
    handleCreateInvite,
    openCreateDialog
  };
}
