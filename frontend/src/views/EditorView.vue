<template>
  <div class="editor-container">
    <div class="editor-header">
      <span>{{ isEditing ? '编辑表单' : '新建表单' }}</span>
      <div class="header-actions">
        <el-button @click="handlePreview">预览</el-button>
        <el-button type="success" @click="handlePublish" :loading="isPublishing">
          发布
        </el-button>
        <el-button type="primary" @click="handleSave" :loading="isSaving">
          保存
        </el-button>
      </div>
    </div>
    <div class="editor-main">
      <div class="left-panel">
        <ComponentPalette />
      </div>
      <div class="center-panel" v-loading="isLoadingForm" element-loading-text="正在加载表单...">
        <!-- 这里必须是 CanvasPanel -->
        <CanvasPanel v-if="!isLoadingForm" />
      </div>
      <div class="right-panel">
        <SettingsPanel />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useEditorStore } from '@/stores/editor'
import { createFormAPI, updateFormAPI, getFormDetailsAPI, updateFormStatusAPI } from '@/api/form'
import ComponentPalette from '@/components/editor/ComponentPalette.vue'
import CanvasPanel from '@/components/editor/CanvasPanel.vue'
import SettingsPanel from '@/components/editor/SettingsPanel.vue'
import { ElMessage, ElMessageBox } from 'element-plus'

const editorStore = useEditorStore()
const route = useRoute()
const router = useRouter()

const formId = ref<number | null>(null)
const isSaving = ref(false)
const isPublishing = ref(false)
const isLoadingForm = ref(false)

const isEditing = computed(() => formId.value !== null)

// 基础校验
const validateForm = () => {
  const { title, questions } = editorStore.formDefinition
  if (!title.trim() || title === '请输入表单标题') {
    ElMessage.warning('表单标题不能为空')
    return false
  }
  if (questions.length === 0) {
    ElMessage.warning('表单至少需要一个题目')
    return false
  }
  return true
}

// 保存逻辑（创建或更新）
const handleSave = async (redirectOnSuccess = true) => {
  if (!validateForm()) return
  isSaving.value = true
  try {
    const formData = {
      title: editorStore.formDefinition.title,
      description: editorStore.formDefinition.description,
      definition: editorStore.formDefinition.questions
    }

    if (isEditing.value && formId.value) {
      await updateFormAPI(formId.value, formData)
    } else {
      const res = await createFormAPI(formData)
      formId.value = res.ID // 创建成功后，获取ID，进入编辑模式
      editorStore.formId = res.ID
    }

    ElMessage.success('保存成功！')
    if (redirectOnSuccess) {
      router.push('/dashboard')
    }
  } catch (error) {
    console.error('保存失败:', error)
  } finally {
    isSaving.value = false
  }
}

// 发布逻辑
const handlePublish = async () => {
  if (!validateForm()) return

  isPublishing.value = true;
  try {
    // 如果是新表单，先保存
    if (!isEditing.value || !formId.value) {
      await handleSave(false); // 调用保存但不跳转
      if (!formId.value) { // 如果保存后仍然没有ID，说明保存失败
        isPublishing.value = false;
        return;
      }
    } else {
      // 如果是旧表单，也先保存一下确保数据最新
      await handleSave(false);
    }

    // 调用发布接口
    await updateFormStatusAPI(formId.value, 2);
    ElMessage.success('发布成功！');
    router.push('/dashboard');
  } catch (error) {
    console.error('发布失败:', error);
  } finally {
    isPublishing.value = false;
  }
};


// 预览逻辑
const handlePreview = () => {
  sessionStorage.setItem('questflow_form_preview', JSON.stringify(editorStore.formDefinition));
  const previewUrl = router.resolve({ name: 'preview' }).href;
  window.open(previewUrl, '_blank');
}

// 加载表单数据
onMounted(async () => {
  const id = Number(route.params.formId)
  if (!isNaN(id) && id > 0) {
    formId.value = id
    isLoadingForm.value = true
    try {
      const res = await getFormDetailsAPI(id)
      editorStore.setForm(res.ID, {
        title: res.Title,
        description: res.Description,
        questions: res.Definition.questions || []
      })
    } catch (error) {
      ElMessage.error('加载表单失败！')
      router.push('/dashboard')
    } finally {
      isLoadingForm.value = false
    }
  } else {
    editorStore.initNewForm()
  }
})
</script>

<style lang="scss" scoped>
.editor-container {
  display: flex;
  flex-direction: column;
  height: 100vh;
  background-color: #f0f2f5;
}

.editor-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  height: 60px;
  padding: 0 20px;
  background-color: #fff;
  border-bottom: 1px solid #e0e0e0;
  flex-shrink: 0;
}

.editor-main {
  display: flex;
  flex-grow: 1;
  overflow: hidden;
}

.left-panel,
.right-panel {
  width: 280px;
  background-color: #fff;
  padding: 15px;
  flex-shrink: 0;
  overflow-y: auto;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
}

.left-panel {
  border-right: 1px solid #e0e0e0;
}

.right-panel {
  border-left: 1px solid #e0e0e0;
}

.center-panel {
  flex-grow: 1;
  padding: 20px;
  overflow-y: auto;
}
</style>
