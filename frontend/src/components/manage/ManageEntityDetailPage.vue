<template>
  <PageLayout
    :system-name="systemName"
    :current-language="currentLanguage"
    :is-dark-mode="isDarkMode"
    :user-avatar="userAvatar"
    :username="userInfo?.username || t('common.user')"
    :active-menu="activeMenu"
    :side-menu-items="manageMenuItems"
    sidebar-scene="manage"
    :title="''"
    content-padding="xl"
    @change-language="handleLanguageChange"
    @toggle-theme="toggleTheme"
    @logout="handleLogout"
    @menu-select="handleMenuSelect"
  >
    <ManageLayout>
      <TitleBar :show-back="true" :compact="true" :back-text="t('common.back')" @back="router.back()">
        <template #actions>
          <el-button type="primary" @click="openEditDialog">
            {{ t('document.edit') }}
          </el-button>
        </template>
      </TitleBar>

      <ManageSection plain>
        <InfoStatsCard
          :title="entityInfo?.name || t('common.unknown')"
          :description="entityInfo?.description || t('common.unknown')"
          :stats="entityStats"
        />
      </ManageSection>

      <ManageSection>
        <CommonList
          :rows="memberRows"
          :columns="memberColumns"
          :is-dark="isDarkMode"
          row-key="external_id"
          :empty-text="t('message.empty')"
          :show-pagination="true"
          :total="memberTotal"
          v-model:current-page="memberPage"
          v-model:page-size="memberPageSize"
          :page-sizes="[6]"
          @page-change="handleMemberPageChange"
          :show-search="true"
          v-model:search-query="memberSearch"
          :search-placeholder="t('common.search')"
          @search="handleMemberSearch"
          :show-title-bar="true"
        >
          <template #title>{{ t(entityLabels.memberListKey) }}</template>
          <template #toolbar-right>
            <el-button size="small" type="primary" @click="openMemberDialog">
              {{ t(entityLabels.manageMembersKey) }}
            </el-button>
          </template>
          <template #cell-actions="{ row }">
            <el-button size="small" text type="danger" @click="removeMember(row)">
              {{ t('document.delete') }}
            </el-button>
          </template>
        </CommonList>
      </ManageSection>
    </ManageLayout>

    <el-dialog v-model="showEditDialog" :title="t('document.edit')" width="480px">
      <el-form label-position="top" :model="editForm">
        <el-form-item :label="t('common.name')">
          <el-input v-model="editForm.name" maxlength="50" show-word-limit />
        </el-form-item>
        <el-form-item :label="t('common.description')">
          <el-input
            v-model="editForm.description"
            type="textarea"
            :rows="3"
            maxlength="200"
            show-word-limit
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showEditDialog = false">{{ t('button.cancel') }}</el-button>
        <el-button type="primary" :loading="editing" @click="submitEdit">
          {{ t('button.confirm') }}
        </el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="showMemberDialog" :title="t(entityLabels.memberListKey)" width="680px">
      <div class="member-dialog">
        <div class="member-dialog__list">
          <CommonList
            :rows="memberCandidates"
            :columns="candidateColumns"
            :is-dark="isDarkMode"
            row-key="external_id"
            :empty-text="t('message.empty')"
            :show-pagination="true"
            :total="candidateTotal"
            v-model:current-page="candidatePage"
            v-model:page-size="candidatePageSize"
            :page-sizes="[3]"
            @page-change="handleCandidatePageChange"
            :show-search="true"
            v-model:search-query="candidateSearch"
            :search-placeholder="t('common.search')"
            @search="handleCandidateSearch"
            :show-title-bar="true"
            :title="t(entityLabels.memberListKey)"
          >
            <template #cell-actions="{ row }">
              <el-checkbox
                v-model="selectedMemberIds"
                :label="row.external_id"
                class="checkbox-only"
              />
            </template>
          </CommonList>
        </div>
      </div>
      <template #footer>
        <el-button @click="showMemberDialog = false">{{ t('button.cancel') }}</el-button>
        <el-button type="primary" :loading="savingMembers" @click="submitMemberUpdate">
          {{ t('button.confirm') }}
        </el-button>
      </template>
    </el-dialog>
  </PageLayout>
</template>

<script setup>
import PageLayout from '@/components/layout/PageLayout.vue';
import TitleBar from '@/components/layout/TitleBar.vue';
import ManageLayout from '@/components/manage/ManageLayout.vue';
import ManageSection from '@/components/manage/ManageSection.vue';
import CommonList from '@/components/list/CommonList.vue';
import InfoStatsCard from '@/components/shared/InfoStatsCard.vue';
import { useEntityDetailPage } from '@/composables/manage/useEntityDetailPage';

const props = defineProps({
  entityType: {
    type: String,
    required: true,
    validator: (value) => ['organization', 'user-group'].includes(value)
  }
});

const {
  router,
  systemName,
  currentLanguage,
  isDarkMode,
  userAvatar,
  userInfo,
  activeMenu,
  manageMenuItems,
  handleLanguageChange,
  toggleTheme,
  handleLogout,
  handleMenuSelect,
  entityLabels,
  entityInfo,
  entityStats,
  memberRows,
  memberColumns,
  memberTotal,
  memberPage,
  memberPageSize,
  handleMemberPageChange,
  memberSearch,
  handleMemberSearch,
  openMemberDialog,
  memberCandidates,
  candidateColumns,
  candidateTotal,
  candidatePage,
  candidatePageSize,
  handleCandidatePageChange,
  candidateSearch,
  handleCandidateSearch,
  selectedMemberIds,
  savingMembers,
  submitMemberUpdate,
  showMemberDialog,
  showEditDialog,
  editing,
  editForm,
  openEditDialog,
  submitEdit,
  removeMember
} = useEntityDetailPage(props);
</script>

<style scoped>
.member-dialog {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.member-dialog__list {
  max-height: 360px;
  overflow: auto;
}

.checkbox-only :deep(.el-checkbox__label) {
  display: none;
}
</style>
