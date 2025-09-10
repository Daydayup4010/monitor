import { defineStore } from 'pinia'
import { ref } from 'vue'
import { ElMessage } from 'element-plus'
import { tokenApi } from '@/api'
import type { UUToken, BuffToken, TokenStatus } from '@/types'

export const useTokenStore = defineStore('token', () => {
  const tokenStatus = ref<TokenStatus>({
    uu: 'yes',  // 默认无效
    buff: 'yes'
  })

  const loading = ref(false)

  // 更新UU Token
  const updateUUToken = async (data: UUToken) => {
    loading.value = true
    try {
      await tokenApi.updateUUToken(data)
      ElMessage.success('UU Token更新成功')
      await verifyTokens()
    } catch (error) {
      console.error('Update UU token failed:', error)
    } finally {
      loading.value = false
    }
  }

  // 更新Buff Token
  const updateBuffToken = async (data: BuffToken) => {
    loading.value = true
    try {
      await tokenApi.updateBuffToken(data)
      ElMessage.success('Buff Token更新成功')
      await verifyTokens()
    } catch (error) {
      console.error('Update Buff token failed:', error)
    } finally {
      loading.value = false
    }
  }

  // 验证Token状态（获取状态用）
  const verifyTokens = async () => {
    try {
      const response = await tokenApi.verifyTokens()
      tokenStatus.value = response.data || { uu: 'yes', buff: 'yes' }
    } catch (error) {
      console.error('Verify tokens failed:', error)
    }
  }

  // 手动验证所有Token（验证按钮用）
  const manualVerifyTokens = async () => {
    try {
      const response = await tokenApi.manualVerifyTokens()
      tokenStatus.value = response.data || { uu: 'yes', buff: 'yes' }
    } catch (error) {
      console.error('Manual verify tokens failed:', error)
    }
  }

  // 检查Token是否有效 (no=有效，yes=无效)
  const isTokenValid = (platform: 'uu' | 'buff'): boolean => {
    return tokenStatus.value[platform] === 'no'
  }

  return {
    tokenStatus,
    loading,
    updateUUToken,
    updateBuffToken,
    verifyTokens,
    manualVerifyTokens,
    isTokenValid,
  }
})
