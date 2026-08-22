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
import { useI18n } from 'vue-i18n';
import CommonList from '@/components/list/CommonList.vue';
import InviteCreateDialog from '@/components/manage/InviteCreateDialog.vue';
import AppTag from '@/components/base/AppTag.vue';
import { useInvitationCenter } from '@/composables/manage/useInvitationCenter';

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

const {
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
} = useInvitationCenter(props);

defineExpose({ openCreateDialog });
</script>

<style scoped>
.invite-code {
  cursor: pointer;
  color: var(--primary-color);
}
</style>
