<template>
  <div class="toolbar">
    <el-dropdown trigger="click">
      <el-button size="small" icon="Setting" circle>
      </el-button>
      <template #dropdown>
        <el-dropdown-menu>
          <!-- 主题切换 -->
          <el-dropdown-item>
            <div class="dropdown-item-content">
              <span>{{ t('theme.light') }} / {{ t('theme.dark') }}</span>
              <ThemeSwitcher />
            </div>
          </el-dropdown-item>
          <!-- 语言切换 -->
          <el-dropdown-item divided>
            <div class="language-options">
              <el-dropdown trigger="click" @command="handleLanguageChange">
                <span class="language-selector">
                  {{ currentLanguage === 'en' ? t('language.english') : t('language.chinese') }}
                  <el-icon class="el-icon--right"><ArrowDown /></el-icon>
                </span>
                <template #dropdown>
                  <el-dropdown-menu>
                    <el-dropdown-item command="en" :icon="'Flag'">
                      {{ t('language.english') }}
                    </el-dropdown-item>
                    <el-dropdown-item command="zh" :icon="'ChatLineRound'">
                      {{ t('language.chinese') }}
                    </el-dropdown-item>
                  </el-dropdown-menu>
                </template>
              </el-dropdown>
            </div>
          </el-dropdown-item>
        </el-dropdown-menu>
      </template>
    </el-dropdown>
  </div>
</template>

<script setup>
import { computed, ref } from 'vue';
import { useI18n } from 'vue-i18n';
import { Setting, ArrowDown, Flag, ChatLineRound } from '@element-plus/icons-vue';
import ThemeSwitcher from './ThemeSwitcher.vue';

const { locale, t } = useI18n();

// 计算当前语言
const currentLanguage = computed({
  get: () => locale.value,
  set: (value) => {
    locale.value = value;
    localStorage.setItem('language', value);
  }
});

// 处理语言切换
const handleLanguageChange = (command) => {
  currentLanguage.value = command;
};

// 初始化语言设置
const initLanguage = () => {
  const savedLanguage = localStorage.getItem('language');
  if (savedLanguage) {
    currentLanguage.value = savedLanguage;
  }
};

// 初始化
initLanguage();
</script>

<style scoped>
.toolbar {
  position: relative;
}

.dropdown-item-content {
  display: flex;
  align-items: center;
  justify-content: space-between;
  width: 100%;
  padding: 0 8px;
}

.language-options {
  width: 100%;
}

.language-selector {
  display: flex;
  align-items: center;
  width: 100%;
  padding: 0 8px;
  cursor: pointer;
}

.language-selector:hover {
  color: var(--el-color-primary);
}

/* 响应式设计 */
@media (max-width: 768px) {
  .toolbar {
    padding: 0;
  }
}
</style>