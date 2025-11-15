/**
 * 从 Blob 数据下载文件
 * @param blob 文件内容的 Blob 对象
 * @param defaultFileName 下载时默认的文件名
 */
export function downloadBlob(blob: Blob, defaultFileName: string) {
  const url = window.URL.createObjectURL(blob);
  const link = document.createElement('a');
  link.href = url;

  // 设置下载的文件名
  link.setAttribute('download', defaultFileName);

  // 将链接添加到 DOM 中，模拟点击，然后移除
  document.body.appendChild(link);
  link.click();
  document.body.removeChild(link);

  // 释放 URL 对象
  window.URL.revokeObjectURL(url);
}
