<template>
  <div class="recent-documents-section">
    <div class="section-header">
      <h3 class="section-title">{{ title }}</h3>
      <el-button v-if="showViewAll" link type="primary" @click="handleViewAll">
        {{ t('common.viewAll') }}
      </el-button>
    </div>
    <div class="section-content">
      <el-empty :description="emptyText" v-if="items.length === 0" />
      <el-card v-for="doc in items" :key="doc.id" class="document-item">
        <div class="item-card-header">
          <h4 class="item-title">{{ doc.title }}</h4>
          <el-tag v-if="doc.status === 'draft'" type="warning" size="small">
            {{ t('document.draft') }}
          </el-tag>
          <el-tag v-else type="success" size="small">
            {{ t('document.published') }}
          </el-tag>
        </div>
        
        <div class="item-meta">
          <span class="meta-item">
            <el-icon><User /></el-icon>
            {{ doc.author }}
          </span>
          <span class="meta-item">
            <el-icon><Timer /></el-icon>
            {{ formatDate(doc.updatedAt) }}
          </span>
          <span class="meta-item">
            <el-icon><View /></el-icon>
            {{ doc.views || 0 }} {{ t('document.views') }}
          </span>
        </div>

        <div class="item-actions">
          <el-button v-if="singleAction" size="small" type="primary" @click="handleView(doc)">
            {{ primaryActionLabel || t('common.open') }}
          </el-button>
          <template v-else>
            <el-button size="small" @click="handleView(doc)">
              {{ t('document.view') }}
            </el-button>
            <el-button size="small" type="primary" @click="handleEdit(doc)">
              {{ t('document.edit') }}
            </el-button>
          </template>
        </div>
      </el-card>
    </div>
  </div>
</template>

<script setup>
import { useI18n } from 'vue-i18n';
import { User, Timer, View } from '@element-plus/icons-vue';

const { t } = useI18n();

// 定义组件属性
const props = defineProps({
  title: {
    type: String,
    default: ''
  },
  items: {
    type: Array,
    default: () => []
  },
  emptyText: {
    type: String,
    default: ''
  },
  showViewAll: {
    type: Boolean,
    default: false
  },
  singleAction: {
    type: Boolean,
    default: false
  },
  primaryActionLabel: {
    type: String,
    default: ''
  }
});

// 定义事件
const emit = defineEmits(['view-all', 'view-item', 'edit-item']);

// 格式化日期
const formatDate = (dateString) => {
  if (!dateString) return t('common.unknown');

  const date = new Date(dateString);
  return date.toLocaleDateString();
};

// 处理查看全部
const handleViewAll = () => {
  emit('view-all');
};

// 处理查看项目
const handleView = (doc) => {
  emit('view-item', doc);
};

// 处理编辑项目
const handleEdit = (doc) => {
  emit('edit-item', doc);
};
</script>

<style scoped>
.recent-documents-section {
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

.document-item {
  margin-bottom: var(--spacing-md);
  transition: box-shadow 0.3s ease;
}

.document-item:hover {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.item-card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
  gap: 10px;
}

.item-title {
  font-weight: bold;
  font-size: 16px;
  color: var(--el-text-color-primary);
  flex: 1;
  margin: 0;
}

.item-meta {
  border-top: 1px solid var(--el-border-color-light);
  padding-top: 10px;
  margin-top: 10px;
  display: flex;
  gap: var(--spacing-md);
  flex-wrap: wrap;
}

.meta-item {
  display: flex;
  align-items: center;
  font-size: 12px;
  color: var(--text-light);
  gap: 4px;
}

.item-actions {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
  margin-top: 15px;
}

/* 深色模式支持 */
.dark-mode .item-title {
  color: var(--text-dark); /* 使用更亮的文字颜色，确保在深色背景下清晰可见 */
}

.dark-mode .meta-item {
  color: var(--text-light); /* 确保元数据文字在深色背景下也清晰可见 */
}
</style>
