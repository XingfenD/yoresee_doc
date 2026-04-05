<template>
  <div class="not-found-page">
    <div class="not-found-card">
      <p class="error-code">404</p>
      <h1 class="error-title">页面不存在</h1>
      <p class="error-desc">
        抱歉，未找到你访问的页面
        <span v-if="missingPath" class="missing-path">（{{ missingPath }}）</span>
      </p>
      <div class="error-actions">
        <el-button @click="goBack">返回上一页</el-button>
        <el-button type="primary" @click="goHome">回到首页</el-button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue';
import { useRoute, useRouter } from 'vue-router';

const router = useRouter();
const route = useRoute();

const missingPath = computed(() => {
  const from = route.query.from;
  if (typeof from !== 'string') {
    return '';
  }
  return from;
});

const goHome = () => {
  router.push('/');
};

const goBack = () => {
  if (window.history.length > 1) {
    router.back();
    return;
  }
  goHome();
};
</script>

<style scoped>
.not-found-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 24px;
  background: var(--bg-light);
}

.not-found-card {
  width: min(520px, 100%);
  padding: 48px 40px;
  text-align: center;
  border: 1px solid var(--border-color);
  border-radius: 16px;
  background: var(--bg-white);
  box-shadow: var(--shadow-sm);
}

.error-code {
  margin: 0;
  font-size: 72px;
  line-height: 1;
  font-weight: 700;
  color: var(--primary-color);
}

.error-title {
  margin: 16px 0 8px;
  font-size: 28px;
  font-weight: 700;
  color: var(--text-dark);
}

.error-desc {
  margin: 0;
  color: var(--text-medium);
}

.missing-path {
  color: var(--text-light);
}

.error-actions {
  margin-top: 24px;
  display: flex;
  justify-content: center;
  gap: 12px;
}

:global(body.dark-mode) .not-found-card {
  box-shadow: none;
}

@media (max-width: 640px) {
  .not-found-card {
    padding: 36px 24px;
  }

  .error-code {
    font-size: 56px;
  }

  .error-title {
    font-size: 22px;
  }
}
</style>
