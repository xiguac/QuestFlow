<template>
  <div class="question-renderer">
    <!-- 题干 -->
    <h4 class="question-title">
      <span class="question-index">{{ index + 1 }}.</span>
      {{ question.title }}
    </h4>

    <!-- 根据题型渲染不同的预览 -->
    <div v-if="question.type === 'single_choice'" class="options-container">
      <el-radio-group disabled>
        <el-radio
          v-for="option in question.options"
          :key="option.id"
          :label="option.id"
          class="option-item"
        >
          {{ option.text }}
        </el-radio>
      </el-radio-group>
    </div>

    <div v-if="question.type === 'text_input'" class="text-input-container">
      <el-input
        type="textarea"
        :rows="3"
        placeholder="用户在此输入回答"
        disabled
      />
    </div>

  </div>
</template>

<script setup lang="ts">
import type { Question } from '@/api/form'
import type { PropType } from 'vue'

defineProps({
  question: {
    type: Object as PropType<Question>,
    required: true
  },
  index: {
    type: Number,
    required: true
  }
})
</script>

<style lang="scss" scoped>
.question-renderer {
  text-align: left;
}
.question-title {
  font-size: 16px;
  margin: 0 0 15px;
  color: #303133;
  word-break: break-all;
}
.question-index {
  margin-right: 8px;
}
.options-container {
  .el-radio-group {
    display: flex;
    flex-direction: column;
    align-items: flex-start;
  }
  .option-item {
    margin-bottom: 10px;
    // 穿透 el-radio 的样式，使其在 disabled 状态下依然清晰
    :deep(.el-radio__input.is-disabled .el-radio__inner) {
      background-color: #f5f7fa;
      border-color: #e4e7ed;
    }
    :deep(.el-radio__label) {
      color: #606266;
    }
  }
}
.text-input-container {
  :deep(.el-textarea__inner[disabled]) {
    background-color: #f5f7fa !important;
    color: #c0c4cc !important;
    cursor: not-allowed;
  }
}
</style>
