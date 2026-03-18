<template>
  <div class="recent-knowledge-base-section">
    <div class="section-header">
      <h3 class="section-title">{{ title }}</h3>
      <el-button v-if="showViewAll" link type="primary" @click="handleViewAll">
        {{ t("common.viewAll") }}
      </el-button>
    </div>
    <div class="section-content">
      <el-empty :description="emptyText" v-if="items.length === 0" />
      <el-card v-for="kb in items" :key="kb.externalId" class="knowledge-base-item">
        <template #header>
          <div class="card-header">
            <span class="kb-name">{{ kb.name }}</span>
            <el-tag v-if="kb.is_public" type="success" size="small">
              {{ t("knowledgeBase.public") }}
            </el-tag>
            <el-tag v-else type="info" size="small">
              {{ t("knowledgeBase.private") }}
            </el-tag>
          </div>
        </template>

        <p class="kb-description">
          {{ kb.description || t("knowledgeBase.noDescription") }}
        </p>

        <div class="kb-details">
          <div class="detail-item">
            <span class="detail-label">{{ t("knowledgeBase.documentsCount") }}:</span>
            <span class="detail-value">{{ kb.documents_count || 0 }}</span>
          </div>
          <div class="detail-item">
            <span class="detail-label">{{ t("knowledgeBase.updatedAt") }}:</span>
            <span class="detail-value">{{ formatDate(kb.updated_at) }}</span>
          </div>
        </div>

        <div class="kb-actions">
          <el-button size="small" @click="handleView(kb)">
            {{ t("common.view") }}
          </el-button>
          <el-button size="small" type="primary" @click="handleAccess(kb)">
            {{ t("knowledgeBase.access") }}
          </el-button>
        </div>
      </el-card>
    </div>
  </div>
</template>

<script setup>
import { useI18n } from "vue-i18n";

const { t } = useI18n();

// 定义组件属性
const props = defineProps({
  title: {
    type: String,
    default: "",
  },
  items: {
    type: Array,
    default: () => [],
  },
  emptyText: {
    type: String,
    default: "",
  },
  showViewAll: {
    type: Boolean,
    default: false,
  },
});

// 定义事件
const emit = defineEmits(["view-all", "view-item", "access-item"]);

// 格式化日期
const formatDate = (dateString) => {
  if (!dateString) return t("common.unknown");

  const date = new Date(dateString);
  return date.toLocaleDateString();
};

// 处理查看全部
const handleViewAll = () => {
  emit("view-all");
};

// 处理查看项目
const handleView = (kb) => {
  emit("view-item", kb);
};

// 处理访问项目
const handleAccess = (kb) => {
  emit("access-item", kb);
};
</script>

<style scoped>
.recent-knowledge-base-section {
  display: flex;
  flex-direction: column;
  background-color: var(--bg-white);
  border-radius: var(--border-radius-md);
  box-shadow: var(--shadow-sm);
  overflow: hidden;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: var(--spacing-md);
  border-bottom: 1px solid var(--border-color);
  background-color: var(--bg-white);
}

.section-title {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
  color: var(--text-primary);
}

.section-content {
  flex: 1;
  overflow-y: auto;
  padding: var(--spacing-md);
}

.knowledge-base-item {
  margin-bottom: var(--spacing-md);
  transition: box-shadow 0.3s ease;
}

.knowledge-base-item:hover {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
  gap: 10px;
}

.kb-name {
  font-weight: bold;
  font-size: 16px;
  color: var(--el-text-color-primary);
  flex: 1;
}

.kb-description {
  color: var(--el-text-color-regular);
  margin: 10px 0;
  line-height: 1.5;
  word-break: break-word;
}

.kb-details {
  border-top: 1px solid var(--el-border-color-light);
  padding-top: 15px;
  margin-top: 10px;
}

.detail-item {
  display: flex;
  justify-content: space-between;
  margin-bottom: 8px;
}

.detail-label {
  color: var(--el-text-color-secondary);
  font-weight: 500;
}

.detail-value {
  color: var(--el-text-color-primary);
  font-weight: 400;
}

.kb-actions {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
  margin-top: 15px;
}

/* 深色模式支持 */
.dark-mode .kb-name {
  color: var(--text-dark); /* 使用更亮的文字颜色，确保在深色背景下清晰可见 */
}

.dark-mode .kb-description {
  color: var(--text-medium); /* 确保描述文字在深色背景下也清晰可见 */
}

.dark-mode .detail-label {
  color: var(--text-light); /* 确保详情标签在深色背景下也清晰可见 */
}

.dark-mode .detail-value {
  color: var(--text-medium); /* 确保详情值在深色背景下也清晰可见 */
}
</style>
