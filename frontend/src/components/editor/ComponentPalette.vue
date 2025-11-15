<template>
  <div class="component-palette">
    <h3>组件库</h3>
    <el-card shadow="never">
      <div v-for="group in componentGroups" :key="group.title" class="group">
        <div class="group-title">{{ group.title }}</div>
        <div class="component-list">
          <div
            v-for="component in group.components"
            :key="component.type"
            class="component-item"
            :class="{ 'full-width': group.components.length === 1 }"
            @click="handleAddComponent(component.type)"
          >
            <el-icon><component :is="component.icon" /></el-icon>
            <span>{{ component.label }}</span>
          </div>
        </div>
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { useEditorStore } from '@/stores/editor'
import { type Question, type QuestionType } from '@/api/form'
import { v4 as uuidv4 } from 'uuid'
import {
  CircleCheck,
  Document,
  Finished,
  SwitchButton,
} from '@element-plus/icons-vue'

const editorStore = useEditorStore()

// 【核心改动】调整组件分组
const componentGroups = [
  {
    title: '选择题',
    components: [
      { type: 'single_choice', label: '单选题', icon: CircleCheck },
      { type: 'multi_choice', label: '多选题', icon: Finished },
    ]
  },
  {
    title: '判断题',
    components: [
      { type: 'judgment', label: '判断题', icon: SwitchButton },
    ]
  },
  {
    title: '填空题',
    components: [
      { type: 'text_input', label: '文本题', icon: Document },
    ]
  }
]

// 处理添加组件的点击事件
const handleAddComponent = (type: QuestionType) => {
  let newQuestion: Question;

  switch (type) {
    case 'single_choice':
      newQuestion = {
        id: uuidv4(),
        type: 'single_choice',
        title: '单选题',
        options: [
          { id: uuidv4(), text: '选项1' },
          { id: uuidv4(), text: '选项2' },
        ]
      };
      break;

    case 'multi_choice':
      newQuestion = {
        id: uuidv4(),
        type: 'multi_choice',
        title: '多选题',
        options: [
          { id: uuidv4(), text: '选项A' },
          { id: uuidv4(), text: '选项B' },
          { id: uuidv4(), text: '选项C' },
        ]
      };
      break;

    case 'judgment':
      newQuestion = {
        id: uuidv4(),
        type: 'judgment',
        title: '判断题',
        options: [
          { id: 'true', text: '正确' },
          { id: 'false', text: '错误' },
        ]
      };
      break;

    case 'text_input':
      newQuestion = {
        id: uuidv4(),
        type: 'text_input',
        title: '文本题',
      };
      break;

    default:
      return;
  }

  editorStore.addQuestion(newQuestion)
}
</script>

<style lang="scss" scoped>
.component-palette {
  h3 {
    margin-top: 0;
    margin-bottom: 15px;
    font-size: 16px;
  }
}
.group { margin-bottom: 15px; }
.group-title {
  color: #606266;
  font-size: 14px;
  margin-bottom: 10px;
}
.component-list {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 10px;
}
.component-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 10px;
  border: 1px solid #e4e7ed;
  border-radius: 4px;
  cursor: pointer;
  transition: all 0.2s;

  &.full-width {
    grid-column: span 2; // 让只有一个组件的分组占据整行
  }

  .el-icon {
    font-size: 24px;
    margin-bottom: 8px;
  }
  span { font-size: 12px; }
  &:hover {
    border-color: #409eff;
    color: #409eff;
    background-color: #ecf5ff;
  }
}
</style>
