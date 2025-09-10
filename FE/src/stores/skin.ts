import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { dataApi } from '@/api'
import type { SkinItem, PaginationParams, Statistics } from '@/types'

export const useSkinStore = defineStore('skin', () => {
  const skinItems = ref<SkinItem[]>([])
  const total = ref(0)
  const loading = ref(false)

  // 分页参数 - 默认每页50条
  const pagination = ref<PaginationParams>({
    page_size: 50,
    page_num: 1,
  })

  // 获取饰品数据
  const getSkinItems = async (params?: Partial<PaginationParams>) => {
    loading.value = true
    try {
      const queryParams = { ...pagination.value, ...params }
      const response = await dataApi.getSkinItems(queryParams)
      skinItems.value = response.data || []
      total.value = response.total || 0
      pagination.value = queryParams
    } catch (error) {
      console.error('Get skin items failed:', error)
    } finally {
      loading.value = false
    }
  }

  // 计算统计数据
  const statistics = computed<Statistics>(() => {
    const items = skinItems.value
    if (!items.length) {
      return {
        totalItems: 0,
        avgProfitRate: 0,
        maxPriceDiff: 0,
        lastUpdateTime: '',
      }
    }

    const avgProfitRate = items.reduce((sum, item) => sum + item.profit_rate, 0) / items.length
    const maxPriceDiff = Math.max(...items.map(item => item.price_diff))
    const lastUpdateTime = items[0]?.updated_at || ''

    return {
      totalItems: total.value,
      avgProfitRate: Number((avgProfitRate * 100).toFixed(2)),
      maxPriceDiff,
      lastUpdateTime,
    }
  })

  return {
    skinItems,
    total,
    loading,
    pagination,
    statistics,
    getSkinItems,
  }
})
