<template>
  <div class="success-page-container">
    <el-result
      icon="success"
      title="提交成功"
      sub-title="感谢您的参与！"
    >
      <template #extra>
        <el-button type="primary" @click="closePage">关闭页面</el-button>
      </template>
    </el-result>
  </div>
</template>

<script setup lang="ts">
const closePage = () => {
  // 尝试关闭当前标签页
  // 注意：出于安全原因，浏览器可能不会允许脚本关闭不是由它自己打开的窗口
  window.close()
}
</script>

<style scoped>
.success-page-container {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100vh;
  background-color: #f0f2f5;
}
</style>
