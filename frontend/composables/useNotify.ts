import { ElMessage, ElMessageBox } from 'element-plus'

const MSG_ZINDEX = 10000
const BOX_ZINDEX = 10001

function extractError(err: any): string {
  return err?.data?.message || err?.message || ''
}

export const useNotify = () => {
  const success = (msg: string) => ElMessage({ message: msg, type: 'success', zIndex: MSG_ZINDEX })

  const error = (err: any, fallback = '操作失败') => {
    const msg = extractError(err) || fallback
    ElMessage({ message: msg, type: 'error', zIndex: MSG_ZINDEX })
  }

  const confirm = (msg: string, title = '提示') =>
    ElMessageBox.confirm(msg, title, {
      type: 'warning',
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      zIndex: BOX_ZINDEX,
    })

  return { success, error, confirm }
}
