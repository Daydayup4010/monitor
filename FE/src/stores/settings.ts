import { defineStore } from 'pinia'
import { ref } from 'vue'
import { settingsApi } from '@/api'
import { showMessage } from '@/utils/message'
import type { Settings } from '@/types'

export const useSettingsStore = defineStore('settings', () => {
  const settings = ref<Settings>({
    min_sell_num: 0,
    min_diff: 0,
    max_sell_price: 10000,
    min_sell_price: 0,
  })

  const loading = ref(false)

  // 获取设置
  const getSettings = async () => {
    loading.value = true
    try {
      const response = await settingsApi.getSettings()
      settings.value = response.data || settings.value
    } catch (error) {
      console.error('Get settings failed:', error)
    } finally {
      loading.value = false
    }
  }

  // 更新设置
  const updateSettings = async (data: Settings) => {
    loading.value = true
    try {
      await settingsApi.updateSettings(data)
      settings.value = data
      showMessage.success('设置更新成功')
    } catch (error) {
      console.error('Update settings failed:', error)
    } finally {
      loading.value = false
    }
  }

  return {
    settings,
    loading,
    getSettings,
    updateSettings,
  }
})
