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
import type { Question } from '@/api/form'
import { v4 as uuidv4 } from 'uuid'
import {
  CircleCheck,
  Document,
} from '@element-plus/icons-vue'

const editorStore = useEditorStore()

const componentGroups = [
  {
    title: '选择题',
    components: [
      { type: 'single_choice', label: '单选题', icon: CircleCheck },
    ]
  },
  {
    title: '填空题',
    components: [
      { type: 'text_input', label: '文本题', icon: Document },
    ]
  }
]

const handleAddComponent = (type: 'single_choice' | 'text_input') => {
  let newQuestion: Question;

  if (type === 'single_choice') {
    newQuestion = {
      id: uuidv4(),
      type: 'single_choice',
      title: '单选题',
      options: [
        { id: uuidv4(), text: '选项1' },
        { id: uuidv4(), text: '选项2' },
      ]
    }
  } else if (type === 'text_input') {
    newQuestion = {
      id: uuidv4(),
      type: 'text_input',
      title: '文本题',
    }
  } else {
    return
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
