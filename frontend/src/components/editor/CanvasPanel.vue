<template>
  <div class="canvas-panel" @click.self="editorStore.selectQuestion(null)">
    <!-- 表单标题和描述 -->
    <div class="form-header" @click.self="editorStore.selectQuestion(null)">
      <input
        v-model="editorStore.formDefinition.title"
        class="form-title-input"
        placeholder="请输入表单标题"
        @click="editorStore.selectQuestion(null)"
      />
      <textarea
        v-model="editorStore.formDefinition.description"
        class="form-description-textarea"
        placeholder="请输入表单描述"
        rows="1"
        @click="editorStore.selectQuestion(null)"
      ></textarea>
    </div>

    <!-- 问题列表 -->
    <draggable
      v-model="editorStore.formDefinition.questions"
      item-key="id"
      class="question-list"
      handle=".drag-handle"
    >
      <template #item="{ element: question, index }">
        <div
          class="question-wrapper"
          :class="{ 'is-selected': editorStore.selectedQuestionId === question.id }"
          @click="editorStore.selectQuestion(question.id)"
        >
          <!-- 拖拽手柄 -->
          <div class="drag-handle">
            <el-icon><Rank /></el-icon>
          </div>
          <!-- 删除按钮 -->
          <div class="delete-handle">
            <el-button
              type="danger"
              :icon="Delete"
              circle
              plain
              size="small"
              @click.stop="handleDeleteQuestion(question.id)"
            />
          </div>
          <!-- 使用 QuestionRenderer 组件来真实渲染问题 -->
          <QuestionRenderer :question="question" :index="index" />
        </div>
      </template>
    </draggable>

    <el-empty v-if="editorStore.formDefinition.questions.length === 0" description="从左侧拖拽或点击添加题目" />
  </div>
</template>

<script setup lang="ts">
import { useEditorStore } from '@/stores/editor'
import draggable from 'vuedraggable'
import { Rank, Delete } from '@element-plus/icons-vue'
import QuestionRenderer from './QuestionRenderer.vue'

const editorStore = useEditorStore()

const handleDeleteQuestion = (id: string) => {
  editorStore.deleteQuestion(id)
}
</script>

<style lang="scss" scoped>
.canvas-panel {
  background-color: #fff;
  border-radius: 5px;
  padding: 24px;
  min-height: 100%;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
}
.form-header {
  border-bottom: 2px solid #e0e0e0;
  padding-bottom: 20px;
  margin-bottom: 20px;

  .form-title-input,
  .form-description-textarea {
    width: 100%;
    border: none;
    outline: none;
    background-color: transparent;
    resize: none;
    box-sizing: border-box;
    cursor: pointer;
    &:focus {
      background-color: #f0f8ff;
      border-radius: 4px;
      cursor: text;
    }
  }
  .form-title-input {
    font-size: 24px;
    font-weight: bold;
    padding: 10px;
  }
  .form-description-textarea {
    font-size: 14px;
    padding: 10px;
    color: #606266;
    font-family: inherit;
  }
}
.question-list { min-height: 300px; }
.question-wrapper {
  position: relative;
  padding: 20px;
  padding-left: 50px;
  border: 1px dashed transparent;
  border-radius: 4px;
  margin-bottom: 15px;
  cursor: pointer;
  transition: all 0.2s;

  .delete-handle {
    position: absolute;
    right: -12px;
    top: -12px;
    z-index: 10;
    visibility: hidden;
    opacity: 0;
    transition: all 0.2s;
  }
  &:hover {
    border-color: #409eff;
    .delete-handle {
      visibility: visible;
      opacity: 1;
    }
  }
  &.is-selected {
    border-color: #409eff;
    border-style: solid;
    background-color: #f4f8fe;
    .delete-handle {
      visibility: visible;
      opacity: 1;
    }
  }
  .drag-handle {
    position: absolute;
    left: 10px;
    top: 50%;
    transform: translateY(-50%);
    cursor: grab;
    color: #909399;
    font-size: 20px;
    &:hover { color: #409eff; }
  }
}
</style>
