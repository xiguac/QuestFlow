<template>
  <div class="dashboard-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>我的表单</span>
          <el-button type="primary" :icon="Plus" @click="handleCreateForm">创建新表单</el-button>
        </div>
      </template>

      <el-table :data="formList" v-loading="loading" style="width: 100%">
        <el-table-column prop="Title" label="标题" min-width="180" />
        <el-table-column prop="Status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.Status)">
              {{ getStatusText(row.Status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="CreatedAt" label="创建时间" width="180">
          <template #default="{ row }">
            {{ formatTime(row.CreatedAt) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="280" fixed="right">
          <template #default="{ row }">
            <div class="action-buttons">
              <el-button type="primary" link :icon="Edit" @click="handleEdit(row)">编辑</el-button>
              <el-button
                type="warning"
                link
                :icon="Share"
                @click="handleShare(row)"
                :disabled="row.Status !== 2"
              >
                分享
              </el-button>
              <el-button type="success" link :icon="DataLine" @click="handleStats(row)">统计</el-button>

              <!-- “更多操作”下拉菜单 -->
              <el-dropdown @command="(command) => handleMoreCommand(command, row)" trigger="click">
                <el-button type="info" link :icon="MoreFilled">
                  更多
                </el-button>
                <template #dropdown>
                  <el-dropdown-menu>
                    <el-dropdown-item command="publish" :disabled="row.Status === 2" :icon="Promotion">
                      发布
                    </el-dropdown-item>
                    <el-dropdown-item command="unpublish" :disabled="row.Status !== 2" :icon="EditPen">
                      转为草稿
                    </el-dropdown-item>
                    <el-dropdown-item command="close" :disabled="row.Status === 3" :icon="CircleClose">
                      关闭
                    </el-dropdown-item>
                    <el-dropdown-item command="delete" divided :icon="Delete">
                      <span style="color: #F56C6C;">删除</span>
                    </el-dropdown-item>
                  </el-dropdown-menu>
                </template>
              </el-dropdown>
            </div>
          </template>
        </el-table-column>
      </el-table>

      <el-empty v-if="!loading && formList.length === 0" description="您还没有创建任何表单" />
    </el-card>

    <!-- 分享链接对话框 -->
    <el-dialog v-model="shareDialogVisible" title="分享表单" width="500">
      <p>任何人都可以通过以下链接填写此表单：</p>
      <el-input v-model="shareUrl" readonly>
        <template #append>
          <el-button @click="copyShareUrl">复制</el-button>
        </template>
      </el-input>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="shareDialogVisible = false">关闭</el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { getMyFormsAPI, deleteFormAPI, updateFormStatusAPI, type FormInfo } from '@/api/form'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Plus, Edit, DataLine, Share, MoreFilled,
  Promotion, EditPen, CircleClose, Delete
} from '@element-plus/icons-vue'

const router = useRouter()
const formList = ref<FormInfo[]>([])
const loading = ref(true)

// --- 分享功能状态 ---
const shareDialogVisible = ref(false)
const shareUrl = ref('')

const handleCreateForm = () => {
  router.push('/editor/new')
}

const fetchFormList = async () => {
  try {
    loading.value = true
    const res = await getMyFormsAPI()
    formList.value = res || []
  } catch (error) {
    ElMessage.error('获取表单列表失败')
    console.error(error)
  } finally {
    loading.value = false
  }
}

// --- 表格内按钮点击处理函数 ---
const handleEdit = (row: FormInfo) => {
  router.push(`/editor/${row.ID}`)
}

const handleShare = (row: FormInfo) => {
  shareUrl.value = `${window.location.origin}/form/${row.FormKey}`
  shareDialogVisible.value = true
}

const copyShareUrl = async () => {
  try {
    await navigator.clipboard.writeText(shareUrl.value)
    ElMessage.success('链接已复制到剪贴板！')
  } catch (err) {
    ElMessage.error('复制失败，请手动复制。')
    console.error('Failed to copy: ', err)
  }
}

const handleStats = (row: FormInfo) => {
  router.push(`/stats/${row.ID}`)
}

const handleDelete = (row: FormInfo) => {
  ElMessageBox.confirm(
    `确定要删除表单 "${row.Title}" 吗？此操作不可恢复。`,
    '警告',
    {
      confirmButtonText: '确定删除',
      cancelButtonText: '取消',
      type: 'warning',
    }
  ).then(async () => {
    try {
      await deleteFormAPI(row.ID)
      ElMessage.success('删除成功！')
      fetchFormList()
    } catch (error) {
      console.error('删除失败:', error)
    }
  }).catch(() => {
    ElMessage.info('已取消删除')
  })
}

// “更多操作”的统一处理函数
const handleMoreCommand = (command: string, row: FormInfo) => {
  if (command === 'delete') {
    handleDelete(row)
  } else {
    // 处理状态变更
    let status: 1 | 2 | 3 | undefined;
    let actionText = '';
    switch (command) {
      case 'publish':
        status = 2;
        actionText = '发布';
        break;
      case 'unpublish':
        status = 1;
        actionText = '转为草稿';
        break;
      case 'close':
        status = 3;
        actionText = '关闭';
        break;
    }
    if (status) {
      updateStatus(row.ID, status, actionText)
    }
  }
}

const updateStatus = async (formId: number, status: 1 | 2 | 3, actionText: string) => {
  try {
    await updateFormStatusAPI(formId, status)
    ElMessage.success(`表单已${actionText}！`)
    fetchFormList()
  } catch (error) {
    console.error(`${actionText}失败:`, error)
  }
}

// --- -------------------- ---

const getStatusText = (status: number) => {
  switch (status) {
    case 1: return '草稿'
    case 2: return '已发布'
    case 3: return '已关闭'
    default: return '未知'
  }
}

const getStatusType = (status: number) => {
  switch (status) {
    case 1: return 'warning'
    case 2: return 'success'
    case 3: return 'info'
    default: return 'danger'
  }
}

const formatTime = (timeStr: string) => {
  const date = new Date(timeStr)
  return date.toLocaleString()
}

onMounted(() => {
  fetchFormList()
})
</script>

<style scoped>
.dashboard-container {
  padding: 10px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.action-buttons {
  display: flex;
  align-items: center;
  gap: 8px; /* 给按钮之间增加一些间距 */
}
</style>
