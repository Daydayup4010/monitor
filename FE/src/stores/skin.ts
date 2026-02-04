import { defineStore } from 'pinia'
import { ref, computed, watch } from 'vue'
import { dataApi } from '@/api'
import type { SkinItem, PaginationParams, Statistics } from '@/types'

export const useSkinStore = defineStore('skin', () => {
  const skinItems = ref<SkinItem[]>([])
  const total = ref(0)
  const loading = ref(false)

  // 从 localStorage 读取保存的设置
  const savedSource = localStorage.getItem('brick_source') || 'buff'
  const savedTarget = localStorage.getItem('brick_target') || 'uu'
  const savedSort = localStorage.getItem('brick_sort') || 'default'
  // 类别改为多选数组，从 localStorage 读取（存储为逗号分隔的字符串）
  const savedCategories = localStorage.getItem('brick_categories')
  const defaultCategories: string[] = []
  // 购买方案和出售方案
  const savedBuyType = localStorage.getItem('brick_buy_type') || 'sell'
  const savedSellType = localStorage.getItem('brick_sell_type') || 'sell'

  // 平台和排序设置（持久化）
  const sourcePlatform = ref(savedSource)
  const targetPlatform = ref(savedTarget)
  const sortOption = ref(savedSort)
  const categories = ref<string[]>(savedCategories ? savedCategories.split(',').filter(c => c) : defaultCategories)
  // 购买方案: sell(在售价购买) / bidding(求购价购买)
  // 出售方案: sell(在售价出售) / bidding(求购价出售)
  const buyType = ref(savedBuyType)
  const sellType = ref(savedSellType)

  // 监听变化，保存到 localStorage
  watch(sourcePlatform, (val) => localStorage.setItem('brick_source', val))
  watch(targetPlatform, (val) => localStorage.setItem('brick_target', val))
  watch(sortOption, (val) => localStorage.setItem('brick_sort', val))
  watch(categories, (val) => localStorage.setItem('brick_categories', val.join(',')), { deep: true })
  watch(buyType, (val) => localStorage.setItem('brick_buy_type', val))
  watch(sellType, (val) => localStorage.setItem('brick_sell_type', val))

  // 排序配置映射
  const sortMap: Record<string, { field: string; desc: boolean }> = {
    'default': { field: 'profit_rate', desc: true },
    'profit_rate_desc': { field: 'profit_rate', desc: true },
    'profit_rate_asc': { field: 'profit_rate', desc: false },
    'price_diff_desc': { field: 'price_diff', desc: true },
    'price_diff_asc': { field: 'price_diff', desc: false },
    'target_price_desc': { field: 'target_price', desc: true },
    'target_price_asc': { field: 'target_price', desc: false },
    'source_price_desc': { field: 'source_price', desc: true },
    'source_price_asc': { field: 'source_price', desc: false },
    'sell_count_desc': { field: 'sell_count', desc: true },
    'sell_count_asc': { field: 'sell_count', desc: false },
    'turn_over_desc': { field: 'turn_over', desc: true },
    'turn_over_asc': { field: 'turn_over', desc: false },
  }

  // 获取当前排序配置
  const getSortConfig = () => {
    return sortMap[sortOption.value] || sortMap['default']
  }

  // 分页参数 - 默认每页50条
  const pagination = ref<PaginationParams>({
    page_size: 50,
    page_num: 1,
    sort: '', // 排序字段
    desc: false, // 是否降序
  })

  // 获取饰品数据
  const getSkinItems = async (params?: Partial<PaginationParams>) => {
    loading.value = true
    try {
      const sortConfig = getSortConfig()
      // 多选类别转为逗号分隔字符串
      const categoryParam = categories.value.length > 0 ? categories.value.join(',') : ''
      const queryParams = { 
        ...pagination.value, 
        sort: sortConfig.field,
        desc: sortConfig.desc,
        source: sourcePlatform.value,
        target: targetPlatform.value,
        category: categoryParam,
        buy_type: buyType.value,
        sell_type: sellType.value,
        ...params 
      }
      const response = await dataApi.getSkinItems(queryParams)
      skinItems.value = response.data || []
      total.value = response.total || 0
      pagination.value = { ...pagination.value, ...params }
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
    sourcePlatform,
    targetPlatform,
    sortOption,
    categories,
    buyType,
    sellType,
    getSortConfig,
    getSkinItems,
  }
})
