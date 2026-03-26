<template>
  <div class="info-stats-card">
    <h2 class="info-title">{{ title }}</h2>
    <p class="info-description">{{ description }}</p>

    <div class="info-stats">
      <div
        v-for="(stat, index) in stats"
        :key="stat.key || `${stat.label}-${index}`"
        class="stat-item"
      >
        <el-icon v-if="stat.icon">
          <component :is="stat.icon" />
        </el-icon>
        <span>{{ formatStatText(stat) }}</span>
      </div>
    </div>
  </div>
</template>

<script setup>
const props = defineProps({
  title: { type: String, default: '' },
  description: { type: String, default: '' },
  stats: {
    type: Array,
    default: () => []
  }
});

const formatStatText = (stat) => {
  if (stat?.text) return stat.text;
  if (stat?.label && Object.prototype.hasOwnProperty.call(stat, 'value')) {
    return `${stat.label}: ${stat.value}`;
  }
  return stat?.label || '';
};
</script>

<style scoped>
.info-stats-card {
  background-color: var(--bg-white);
  border-radius: var(--border-radius-md);
  box-shadow: var(--shadow-sm);
  padding: var(--spacing-lg);
  border: 1px solid var(--border-color);
}

.info-title {
  margin: 0 0 var(--spacing-md) 0;
  font-size: 24px;
  font-weight: 600;
  color: var(--text-dark);
}

.info-description {
  margin: 0 0 var(--spacing-lg) 0;
  font-size: 16px;
  color: var(--text-medium);
  line-height: 1.6;
}

.info-stats {
  display: flex;
  gap: var(--spacing-lg);
  flex-wrap: wrap;
}

.stat-item {
  display: flex;
  align-items: center;
  gap: var(--spacing-xs);
  color: var(--text-medium);
  font-size: 14px;
}

.stat-item .el-icon {
  color: var(--primary-color);
}

@media (max-width: 768px) {
  .info-stats {
    flex-direction: column;
  }
}

:global(.dark-mode) .info-stats-card {
  background: var(--bg-medium);
  border-color: var(--border-color);
}

:global(.dark-mode) .info-title {
  color: var(--text-dark);
}

:global(.dark-mode) .info-description {
  color: var(--text-medium);
}

:global(.dark-mode) .stat-item {
  color: var(--text-medium);
}
</style>
