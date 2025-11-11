<template>
  <div class="settings-panel">
    <!-- 当没有选中任何问题时，显示表单全局设置 -->
    <div v-if="selectedQuestionIndex === -1">
      <h3>表单设置</h3>
      <el-card shadow="never">
        <el-form label-position="top">
          <el-form-item label="表单标题">
            <el-input v-model="editorStore.formDefinition.title" />
          </el-form-item>
          <el-form-item label="表单描述">
            <el-input v-model="editorStore.formDefinition.description" type="textarea" :rows="3" />
          </el-form-item>
        </el-form>
      </el-card>
    </div>

    <!-- 当选中一个问题时，显示该问题的特定设置 -->
    <div v-else>
      <h3>题目设置</h3>
      <el-card shadow="never">
        <el-form label-position="top">
          <el-form-item label="题型">
            <el-input :model-value="getQuestionTypeText(selectedQuestion.type)" disabled />
          </el-form-item>
          <el-form-item label="题干">
            <el-input v-model="editorStore.formDefinition.questions[selectedQuestionIndex].title" type="textarea" :rows="2" />
          </el-form-item>

          <!-- 特定题型的设置，例如单选题的选项 -->
          <div v-if="selectedQuestion.type === 'single_choice'">
            <el-divider>选项设置</el-divider>
            <div
              v-for="(option, index) in selectedQuestion.options"
              :key="option.id"
              class="option-item"
            >
              <el-input v-model="editorStore.formDefinition.questions[selectedQuestionIndex].options[index].text" placeholder="请输入选项内容">
                <template #prepend>
                  <el-radio :label="option.id" :model-value="null" disabled />
                </template>
                <template #append>
                  <el-button
                    :icon="Delete"
                    circle
                    plain
                    type="danger"
                    size="small"
                    @click="removeOption(index)"
                  />
                </template>
              </el-input>
            </div>
            <el-button type="primary" link @click="addOption">
              <el-icon><Plus></Plus></el-icon>
              添加选项
            </el-button>
          </div>
        </el-form>
      </el-card>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useEditorStore } from '@/stores/editor'
import { v4 as uuidv4 } from 'uuid'
import { Delete, Plus } from '@element-plus/icons-vue'

const editorStore = useEditorStore()

const selectedQuestionIndex = computed(() => {
  if (!editorStore.selectedQuestionId) {
    return -1
  }
  return editorStore.formDefinition.questions.findIndex(
    (q) => q.id === editorStore.selectedQuestionId
  )
})

const selectedQuestion = computed(() => {
  if (selectedQuestionIndex.value > -1) {
    return editorStore.formDefinition.questions[selectedQuestionIndex.value]
  }
  return null
})

const getQuestionTypeText = (type: string) => {
  switch (type) {
    case 'single_choice': return '单选题'
    case 'text_input': return '文本题'
    default: return '未知题型'
  }
}

const addOption = () => {
  if (selectedQuestion.value && selectedQuestion.value.type === 'single_choice') {
    if (!selectedQuestion.value.options) {
      editorStore.formDefinition.questions[selectedQuestionIndex.value].options = []
    }

    editorStore.formDefinition.questions[selectedQuestionIndex.value].options.push({
      id: uuidv4(),
      text: `新选项 ${selectedQuestion.value.options.length + 1}`
    })
  }
}

const removeOption = (index: number) => {
  if (selectedQuestion.value && selectedQuestion.value.options) {
    editorStore.formDefinition.questions[selectedQuestionIndex.value].options.splice(index, 1)
  }
}
</script>

<style lang="scss" scoped>
.settings-panel {
  h3 {
    margin-top: 0;
    margin-bottom: 15px;
    font-size: 16px;
  }
}
.option-item {
  margin-bottom: 10px;
  :deep(.el-input-group__prepend) {
    .el-radio {
      margin-right: 0;
      .el-radio__label {
        display: none;
      }
    }
  }
}
</style>
