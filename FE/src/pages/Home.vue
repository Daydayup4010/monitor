<template>
  <div class="home-page">
    <!-- 筛选区域 -->
    <div class="filter-section">
      <div class="filter-content">
        <div class="filter-row">
          <div class="filter-item search-item">
            <div class="filter-label">搜索饰品</div>
            <el-input
              v-model="searchKeyword"
              placeholder="搜索饰品名称..."
              clearable
              size="large"
              @input="handleSearch"
              class="filter-input"
            >
              <template #prefix>
                <el-icon><Search /></el-icon>
              </template>
            </el-input>
          </div>
          
          <div class="filter-item select-item">
            <div class="filter-label">饰品类别</div>
            <el-select
              v-model="selectedCategory"
              placeholder="选择类别"
              clearable
              size="large"
              class="filter-select"
            >
              <el-option
                v-for="category in categories"
                :key="category"
                :label="category"
                :value="category"
              />
            </el-select>
          </div>
          
          <div class="filter-item slider-item">
            <div class="filter-label">利润率范围</div>
            <div class="slider-container">
              <el-slider
                v-model="profitRateRange"
                range
                :min="0"
                :max="100"
                :step="1"
                @change="handleProfitRateChange"
                class="filter-slider"
              />
              <div class="range-display">
                {{ profitRateRange[0] }}% - {{ profitRateRange[1] }}%
              </div>
            </div>
          </div>
          
          <div class="filter-item refresh-item">
            <div class="update-time-wrapper" v-if="lastUpdateTime">
              <div class="update-time-label">
                <el-icon><Clock /></el-icon>
                <span>更新时间</span>
              </div>
              <div class="update-time-value">{{ formatUpdateTime(lastUpdateTime) }}</div>
            </div>
            <el-button
              type="primary"
              size="large"
              :loading="skinStore.loading"
              @click="refreshData"
              class="refresh-btn"
            >
              <el-icon><Refresh /></el-icon>
              刷新
            </el-button>
          </div>
        </div>
      </div>
    </div>

    <!-- 数据表格 -->
    <el-card class="table-card">
      <div v-if="skinStore.loading" class="loading-container" v-loading="true" element-loading-text="加载饰品数据中..." element-loading-background="transparent">
      </div>
      
      <div v-else class="table-container">
        <el-table
          ref="tableRef"
          :data="displayData"
          style="width: 100%; min-width: 600px"
          stripe
          size="large"
          class="modern-table"
          @sort-change="handleSortChange"
          :height="tableHeight"
          :max-height="tableHeight"
          :default-sort="{ prop: sortConfig.prop || '', order: sortConfig.order || '' }"
          :key="tableKey"
        >
        <el-table-column type="index" label="#" width="70" class-name="index-column" />
        
        <el-table-column prop="name" label="饰品名称" min-width="300" fixed="left" class-name="name-column">
          <template #default="{ row }">
            <div class="skin-info">
              <div class="skin-image">
                <img :src="row.image_url" @error="handleImageError" class="skin-avatar" />
              </div>
              <div class="skin-details">
                <div class="skin-name">{{ row.name }}</div>
                <div class="skin-category">{{ row.category }}</div>
              </div>
            </div>
          </template>
        </el-table-column>
        
        <el-table-column prop="u_price" label="UU价格" width="140" sortable class-name="price-column">
          <template #default="{ row }">
            <div class="price-cell">
              <span class="price-value">¥{{ formatPrice(row.u_price) }}</span>
            </div>
          </template>
        </el-table-column>
        
        <el-table-column prop="buff_price" label="Buff价格" width="140" sortable class-name="price-column">
          <template #default="{ row }">
            <div class="price-cell">
              <span class="price-value">¥{{ formatPrice(row.buff_price) }}</span>
            </div>
          </template>
        </el-table-column>
        
        <el-table-column prop="price_diff" label="价格差" width="140" sortable class-name="diff-column">
          <template #default="{ row }">
            <div class="diff-cell">
              <span
                class="diff-value"
                :style="{ 
                  color: getProfitColor(row.profit_rate)
                }"
              >
                ¥{{ formatPrice(row.price_diff) }}
              </span>
            </div>
          </template>
        </el-table-column>
        
        <el-table-column prop="profit_rate" label="利润率" width="140" sortable class-name="profit-column">
          <template #default="{ row }">
            <div class="profit-cell">
              <el-tag 
                :type="getProfitTagType(row.profit_rate)"
                size="large"
                class="profit-tag"
              >
                {{ formatPercent(row.profit_rate) }}
              </el-tag>
            </div>
          </template>
        </el-table-column>
        </el-table>
      </div>

      <!-- 自定义分页 -->
      <div class="custom-pagination">
        <div class="pagination-info">
          共 {{ skinStore.total }} 条数据
        </div>
        
        <div class="pagination-size">
          <span class="size-label">每页显示</span>
          <div class="size-options">
            <div
              v-for="size in [20, 50, 100, 200]"
              :key="size"
              class="size-option"
              :class="{ active: skinStore.pagination.page_size === size }"
              @click="handleSizeChange(size)"
            >
              {{ size }}
            </div>
          </div>
          <span class="size-label">条</span>
        </div>
        
        <div class="pagination-controls">
          <button
            class="page-btn"
            :disabled="skinStore.pagination.page_num === 1"
            @click="handleCurrentChange(skinStore.pagination.page_num - 1)"
          >
            <el-icon><ArrowLeft /></el-icon>
          </button>
          
          <div class="page-numbers">
            <span
              v-for="page in visiblePages"
              :key="page"
              class="page-number"
              :class="{ 
                active: page === skinStore.pagination.page_num,
                ellipsis: page === '...'
              }"
              @click="page !== '...' && handleCurrentChange(page as number)"
            >
              {{ page }}
            </span>
          </div>
          
          <button
            class="page-btn"
            :disabled="skinStore.pagination.page_num >= totalPages"
            @click="handleCurrentChange(skinStore.pagination.page_num + 1)"
          >
            <el-icon><ArrowRight /></el-icon>
          </button>
        </div>
        
        <div class="page-info">
          第 {{ skinStore.pagination.page_num }} 页，共 {{ totalPages }} 页
        </div>
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, nextTick, watch } from 'vue'
import { useSkinStore } from '@/stores/skin'
import { formatTime, formatPrice, formatPercent, getProfitColor, debounce } from '@/utils'
import dayjs from 'dayjs'
import type { SkinItem } from '@/types'

const skinStore = useSkinStore()

// 表格ref
const tableRef = ref()

// 更新时间
const lastUpdateTime = ref('')

// 搜索和筛选
const searchKeyword = ref('')
const selectedCategory = ref('')
const profitRateRange = ref([0, 100])

// 排序配置 - 现在通过后端处理
const sortConfig = ref<{ prop: string; order: string }>({ prop: '', order: '' })

// 表格高度计算
const tableHeight = ref(500)

// 表格key，用于强制重新渲染
const tableKey = ref(0)

// 获取所有类别
const categories = computed(() => {
  const cats = new Set(skinStore.skinItems.map(item => item.category))
  return Array.from(cats).filter(Boolean)
})

// 注意：搜索和筛选功能现在需要在后端实现
// 这里保留前端状态用于UI显示，但实际筛选逻辑应该通过API参数传递给后端

// 显示数据 - 排序现在由后端处理
const displayData = computed(() => {
  // 直接使用store中的数据，排序已经在后端完成
  return skinStore.skinItems
})

// 分页相关计算属性
const totalPages = computed(() => {
  return Math.ceil(skinStore.total / skinStore.pagination.page_size)
})

const visiblePages = computed(() => {
  const current = skinStore.pagination.page_num
  const total = totalPages.value
  const pages: (number | string)[] = []
  
  if (total <= 7) {
    // 如果总页数少于7页，显示所有页码
    for (let i = 1; i <= total; i++) {
      pages.push(i)
    }
  } else {
    // 总是显示第一页
    pages.push(1)
    
    if (current <= 4) {
      // 当前页在前面时
      for (let i = 2; i <= 5; i++) {
        pages.push(i)
      }
      pages.push('...')
      pages.push(total)
    } else if (current >= total - 3) {
      // 当前页在后面时
      pages.push('...')
      for (let i = total - 4; i <= total; i++) {
        pages.push(i)
      }
    } else {
      // 当前页在中间时
      pages.push('...')
      for (let i = current - 1; i <= current + 1; i++) {
        pages.push(i)
      }
      pages.push('...')
      pages.push(total)
    }
  }
  
  return pages
})

// 处理搜索
const handleSearch = debounce(() => {
  // 搜索逻辑在computed中处理
}, 300)

// 处理利润率范围变化
const handleProfitRateChange = () => {
  // 筛选逻辑在computed中处理
}

// 监听数据变化，强制表格重新渲染以应用排序状态
watch(() => skinStore.skinItems, () => {
  tableKey.value++
})

// 处理排序 - 调用后端API
const handleSortChange = ({ column, prop, order }: any) => {
  if (!prop) return
  
  // 更新排序配置
  sortConfig.value = {
    prop: prop || '',
    order: order || ''
  }
  
  // 将前端排序字段映射到后端字段
  const fieldMap: Record<string, string> = {
    'buff_price': 'buff_price',
    'u_price': 'u_price', 
    'price_diff': 'price_diff',
    'profit_rate': 'profit_rate'
  }
  
  const sortField = fieldMap[prop] || ''
  let isDesc = false
  
  if (order === 'descending') {
    isDesc = true
  } else if (order === 'ascending') {
    isDesc = false
  }
  
  // 调用后端API进行排序
  if (sortField && order) {
    skinStore.getSkinItems({ 
      sort: sortField, 
      desc: isDesc,
      page_num: 1 // 排序时重置到第一页
    })
  } else {
    // 如果没有排序字段或order为null，清除排序
    skinStore.getSkinItems({ 
      sort: '', 
      desc: false,
      page_num: 1
    })
  }
}

// 处理页面大小变化
const handleSizeChange = (size: number) => {
  skinStore.getSkinItems({ page_size: size, page_num: 1 })
}

// 处理当前页变化
const handleCurrentChange = (page: number) => {
  if (page >= 1 && page <= totalPages.value) {
    skinStore.getSkinItems({ page_num: page })
  }
}

// 刷新数据
const refreshData = async () => {
  await skinStore.getSkinItems()
  // 更新时间：取第一个饰品的update_at值
  if (skinStore.skinItems.length > 0) {
    lastUpdateTime.value = skinStore.skinItems[0].updated_at
  }
}

// 格式化更新时间
const formatUpdateTime = (timeStr: string) => {
  if (!timeStr) return ''
  return dayjs(timeStr).format('YYYY-MM-DD HH:mm:ss')
}

// 处理图片加载失败
const handleImageError = (e: Event) => {
  const img = e.target as HTMLImageElement
  img.src = 'data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iNjAiIGhlaWdodD0iNjAiIHZpZXdCb3g9IjAgMCA2MCA2MCIgZmlsbD0ibm9uZSIgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj4KPHJlY3Qgd2lkdGg9IjYwIiBoZWlnaHQ9IjYwIiBmaWxsPSIjRjhGQUZGIiByeD0iNCIvPgo8cGF0aCBkPSJNMzAgMTVMMzggMjhIMjJMMzAgMTVaIiBmaWxsPSIjMTg5MEZGIi8+CjxjaXJjbGUgY3g9IjIyIiBjeT0iMjIiIHI9IjQiIGZpbGw9IiM0MEE5RkYiLz4KPHJlY3QgeD0iMTUiIHk9IjM1IiB3aWR0aD0iMzAiIGhlaWdodD0iNCIgcng9IjIiIGZpbGw9IiNGQUFEMTQiLz4KPHJlY3QgeD0iMjAiIHk9IjQyIiB3aWR0aD0iMjAiIGhlaWdodD0iMyIgcng9IjEuNSIgZmlsbD0iIzcyMkVEMSIvPgo8L3N2Zz4K' // 更大的默认饰品图标
}

// 获取利润率标签类型
const getProfitTagType = (rate: number) => {
  if (rate > 0.2) return 'success'
  if (rate > 0.1) return 'warning'
  if (rate > 0.05) return 'info'
  return 'danger'
}

// 计算表格高度
const calculateTableHeight = () => {
  // 计算表格区域的可用高度，为分页留出空间
  const windowHeight = window.innerHeight
  const headerHeight = 70      // 顶部导航
  const filterHeight = 140     // 筛选区域
  const paginationHeight = 100 // 分页区域
  const padding = 100          // 内边距
  
  // 表格高度 = 卡片高度 - 分页高度 - 卡片内边距
  const cardHeight = windowHeight - headerHeight - filterHeight - padding
  tableHeight.value = cardHeight - paginationHeight - 40
  
  // 最小高度限制
  if (tableHeight.value < 300) {
    tableHeight.value = 300
  }
}

// 初始化数据
onMounted(async () => {
  await skinStore.getSkinItems()
  // 初始化更新时间
  if (skinStore.skinItems.length > 0) {
    lastUpdateTime.value = skinStore.skinItems[0].updated_at
  }
  calculateTableHeight()
  
  // 监听窗口大小变化
  window.addEventListener('resize', calculateTableHeight)
})

// 组件卸载时清理事件监听
onUnmounted(() => {
  window.removeEventListener('resize', calculateTableHeight)
})
</script>

<style scoped>
.home-page {
  padding: 32px 48px 40px 48px;
  min-height: 100vh;
  position: relative;
}

/* 筛选区域 */
.filter-section {
  margin-bottom: 40px;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border-radius: 16px;
  padding: 32px 40px;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.08);
  border: 1px solid rgba(24, 144, 255, 0.1);
}

.filter-content {
  width: 100%;
}

.filter-row {
  display: grid;
  grid-template-columns: 300px 200px 1fr 140px;
  gap: 48px;
  align-items: end;
  max-width: none;
  justify-content: start;
}

.filter-item {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.refresh-item {
  align-items: center;
  max-width: 200px;
}

.filter-label {
  font-size: 15px;
  font-weight: 700;
  color: #1890ff;
  margin-bottom: 4px;
  text-align: left;
}

.filter-input,
.filter-select {
  width: 100%;
  height: 48px;
  border-radius: 12px;
}

:deep(.filter-input .el-input__wrapper) {
  height: 48px;
  padding: 0 16px;
  border: 1px solid rgba(24, 144, 255, 0.2);
  background: rgba(255, 255, 255, 0.95);
  border-radius: 12px;
}

:deep(.filter-input .el-input__wrapper.is-focus) {
  border-color: #1890ff;
}

:deep(.filter-select .el-select__wrapper) {
  height: 48px;
  border: 1px solid rgba(24, 144, 255, 0.2);
  background: rgba(255, 255, 255, 0.95);
  border-radius: 12px;
}

:deep(.filter-select .el-select__wrapper.is-focused) {
  border-color: #1890ff;
}

.slider-container {
  background: rgba(255, 255, 255, 0.95);
  border: 1px solid rgba(24, 144, 255, 0.2);
  border-radius: 12px;
  padding: 20px;
  height: 48px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 20px;
  max-width: 500px;
}

.filter-slider {
  flex: 1;
  margin: 0;
}

.range-display {
  font-size: 13px;
  color: #1890ff;
  font-weight: 600;
  background: rgba(24, 144, 255, 0.1);
  padding: 4px 8px;
  border-radius: 8px;
  white-space: nowrap;
  min-width: 80px;
  text-align: center;
}

.update-time-wrapper {
  background: rgba(255, 255, 255, 0.95);
  border: 2px solid rgba(24, 144, 255, 0.1);
  border-radius: 12px;
  padding: 8px 16px;
  margin-bottom: 12px;
  text-align: center;
  transition: all 0.3s ease;
  backdrop-filter: blur(10px);
  box-shadow: 0 4px 20px rgba(24, 144, 255, 0.08);
  width: 100%;
  max-width: 200px;
  min-width: 180px;
}

.update-time-wrapper:hover {
  border-color: rgba(24, 144, 255, 0.25);
  background: rgba(255, 255, 255, 1);
  transform: translateY(-1px);
  box-shadow: 0 6px 25px rgba(24, 144, 255, 0.15);
}

.update-time-label {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  font-size: 12px;
  color: #8c8c8c;
  font-weight: 500;
  margin-bottom: 4px;
  white-space: nowrap;
}

.update-time-label .el-icon {
  color: #1890ff;
  font-size: 14px;
}

.update-time-value {
  font-size: 12px;
  color: #262626;
  font-weight: 600;
  font-family: 'Courier New', Monaco, monospace;
  white-space: nowrap;
  letter-spacing: 0.1px;
  line-height: 1.2;
}

.refresh-btn {
  height: 44px;
  border-radius: 12px;
  font-weight: 600;
  background: linear-gradient(135deg, #1890ff, #40a9ff);
  border: none;
  box-shadow: 0 4px 15px rgba(24, 144, 255, 0.3);
  transition: all 0.3s ease;
  width: 100%;
  max-width: 120px;
  font-size: 14px;
}

.refresh-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(24, 144, 255, 0.4);
}

.refresh-btn:active {
  transform: translateY(0);
}

/* 表格美化 */
.table-card {
  height: calc(100vh - 400px);
  border-radius: 16px;
  overflow: hidden;
  background: rgba(255, 255, 255, 0.98);
  backdrop-filter: blur(10px);
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.08);
  border: 1px solid rgba(24, 144, 255, 0.1);
  position: relative;
  display: flex;
  flex-direction: column;
}

.table-container {
  flex: 1;
  border-radius: 20px;
  position: relative;
  overflow: auto;
  min-height: 0;
}

.modern-table {
  font-size: 16px;
  border-radius: 12px;
  position: relative;
}

.loading-container {
  min-height: 400px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(248, 250, 252, 0.5);
  border-radius: 12px;
}

/* 表格单元格样式 */
.skin-info {
  display: flex;
  align-items: center;
  gap: 20px;
  padding: 12px 0;
}

.skin-image {
  width: 60px;
  height: 60px;
  border: 2px solid rgba(24, 144, 255, 0.2);
  border-radius: 8px;
  overflow: hidden;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  background: rgba(248, 250, 252, 0.8);
}

.skin-avatar {
  width: 100%;
  height: 100%;
  object-fit: cover;
  border-radius: 6px;
}

.skin-details {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.skin-name {
  font-size: 16px;
  color: #262626;
  font-weight: 600;
  line-height: 1.4;
}

.skin-category {
  font-size: 13px;
  color: #8c8c8c;
  font-weight: 500;
}

.price-cell, .diff-cell, .profit-cell {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 8px 0;
}

.price-value {
  font-size: 16px;
  font-weight: 700;
  color: #1890ff;
  background: rgba(24, 144, 255, 0.1);
  padding: 8px 16px;
  border-radius: 8px;
  display: inline-block;
  min-width: 90px;
  text-align: center;
  border: 1px solid rgba(24, 144, 255, 0.2);
}

.diff-value {
  font-size: 16px;
  font-weight: 700;
  padding: 8px 16px;
  border-radius: 8px;
  display: inline-block;
  min-width: 90px;
  text-align: center;
  background: rgba(0, 0, 0, 0.05);
  border: 1px solid rgba(0, 0, 0, 0.1);
}

.profit-tag {
  font-size: 15px;
  font-weight: 700;
  padding: 8px 16px;
  border-radius: 8px;
  min-width: 80px;
  text-align: center;
}

/* 表格列样式 */
:deep(.index-column) {
  text-align: center;
  font-weight: 600;
  color: #666;
  background: rgba(0, 0, 0, 0.02);
}

:deep(.name-column) {
  background: rgba(24, 144, 255, 0.02);
}

:deep(.price-column) {
  background: rgba(24, 144, 255, 0.02);
  text-align: center;
}

:deep(.diff-column) {
  background: rgba(250, 173, 20, 0.02);
  text-align: center;
}

:deep(.profit-column) {
  background: rgba(82, 196, 26, 0.02);
  text-align: center;
}

/* 表格头部样式 */
:deep(.el-table th) {
  background: rgba(248, 250, 255, 0.95);
  font-weight: 700;
  font-size: 15px;
  color: #1890ff;
  border-bottom: 2px solid #1890ff;
  height: 60px;
}

:deep(.el-table td) {
  border-bottom: 1px solid rgba(24, 144, 255, 0.08);
  height: 80px;
}

:deep(.el-table--striped .el-table__body tr.el-table__row--striped td) {
  background: rgba(248, 250, 252, 0.5);
}

:deep(.el-table__body tr:hover td) {
  background: rgba(24, 144, 255, 0.05) !important;
}

/* 确保表格滚动 */
:deep(.el-table__body-wrapper) {
  overflow-y: auto !important;
  max-height: 100% !important;
}

:deep(.el-table__body) {
  overflow: visible;
}

/* 滚动条美化 */
:deep(.el-table__body-wrapper::-webkit-scrollbar) {
  width: 8px;
}

:deep(.el-table__body-wrapper::-webkit-scrollbar-track) {
  background: rgba(0, 0, 0, 0.05);
  border-radius: 4px;
}

:deep(.el-table__body-wrapper::-webkit-scrollbar-thumb) {
  background: linear-gradient(135deg, #1890ff, #40a9ff);
  border-radius: 4px;
  transition: all 0.3s ease;
}

:deep(.el-table__body-wrapper::-webkit-scrollbar-thumb:hover) {
  background: linear-gradient(135deg, #40a9ff, #1890ff);
  box-shadow: 0 2px 8px rgba(24, 144, 255, 0.3);
}

/* 自定义分页样式 */
.custom-pagination {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 32px;
  margin-top: 16px;
  padding: 20px 24px;
  background: linear-gradient(135deg, rgba(248, 250, 255, 0.9), rgba(230, 247, 255, 0.8));
  border-radius: 20px;
  box-shadow: 0 4px 20px rgba(24, 144, 255, 0.1);
  backdrop-filter: blur(10px);
  flex-shrink: 0;
  min-height: 80px;
}

.pagination-info, .page-info {
  font-size: 15px;
  font-weight: 600;
  color: #1890ff;
  background: rgba(255, 255, 255, 0.8);
  padding: 12px 20px;
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(24, 144, 255, 0.1);
}

.pagination-size {
  display: flex;
  align-items: center;
  gap: 12px;
  background: rgba(255, 255, 255, 0.9);
  padding: 12px 20px;
  border-radius: 16px;
  box-shadow: 0 4px 15px rgba(24, 144, 255, 0.1);
}

.size-label {
  font-size: 14px;
  font-weight: 600;
  color: #666;
}

.size-options {
  display: flex;
  gap: 8px;
}

.size-option {
  width: 40px;
  height: 36px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(248, 250, 252, 0.8);
  border: 2px solid transparent;
  border-radius: 10px;
  font-size: 14px;
  font-weight: 600;
  color: #666;
  cursor: pointer;
  transition: all 0.3s ease;
}

.size-option:hover {
  background: rgba(24, 144, 255, 0.1);
  border-color: rgba(24, 144, 255, 0.3);
  color: #1890ff;
  transform: translateY(-1px);
}

.size-option.active {
  background: linear-gradient(135deg, #1890ff, #40a9ff);
  color: white;
  border-color: #1890ff;
  box-shadow: 0 4px 12px rgba(24, 144, 255, 0.3);
}

.pagination-controls {
  display: flex;
  align-items: center;
  gap: 8px;
  background: rgba(255, 255, 255, 0.9);
  padding: 12px 16px;
  border-radius: 16px;
  box-shadow: 0 4px 15px rgba(24, 144, 255, 0.1);
}

.page-btn {
  width: 40px;
  height: 40px;
  border: none;
  border-radius: 10px;
  background: rgba(248, 250, 252, 0.8);
  color: #666;
  cursor: pointer;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 16px;
}

.page-btn:hover:not(:disabled) {
  background: rgba(24, 144, 255, 0.1);
  color: #1890ff;
  transform: translateY(-1px);
}

.page-btn:disabled {
  opacity: 0.4;
  cursor: not-allowed;
}

.page-numbers {
  display: flex;
  gap: 4px;
}

.page-number {
  width: 36px;
  height: 36px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 600;
  color: #666;
  cursor: pointer;
  transition: all 0.3s ease;
}

.page-number:hover:not(.ellipsis) {
  background: rgba(24, 144, 255, 0.1);
  color: #1890ff;
  transform: translateY(-1px);
}

.page-number.active {
  background: linear-gradient(135deg, #1890ff, #40a9ff);
  color: white;
  box-shadow: 0 3px 10px rgba(24, 144, 255, 0.3);
}

.page-number.ellipsis {
  cursor: default;
  color: #ccc;
}

.page-number.ellipsis:hover {
  background: transparent;
  transform: none;
}

/* 响应式设计 */
@media (max-width: 1400px) {
  .filter-row {
    grid-template-columns: 280px 180px 1fr 130px;
    gap: 40px;
  }
}

@media (max-width: 1200px) {
  .filter-row {
    grid-template-columns: 260px 160px 1fr 120px;
    gap: 32px;
  }
  
  .filter-section {
    padding: 32px 40px;
  }
}

@media (max-width: 900px) {
  .filter-row {
    grid-template-columns: 1fr 1fr;
    gap: 24px;
  }
  
  .slider-item {
    grid-column: 1 / -1;
  }
  
  .slider-container {
    max-width: 500px;
    margin: 0 auto;
  }
  
  .refresh-item {
    grid-column: 1 / -1;
    justify-self: center;
    max-width: 200px;
  }
  
  .refresh-btn {
    max-width: 120px;
  }
  
  .filter-section {
    padding: 28px 32px;
  }
}

@media (max-width: 768px) {
  .home-page {
    padding: 12px 12px 24px 12px;
  }
  
  .filter-section {
    padding: 20px 16px;
    margin-bottom: 20px;
  }
  
  .filter-row {
    grid-template-columns: 1fr;
    gap: 16px;
  }
  
  .filter-item {
    width: 100%;
  }
  
  .filter-input,
  .filter-select {
    height: 44px;
  }
  
  :deep(.filter-input .el-input__wrapper) {
    height: 44px;
  }
  
  :deep(.filter-select .el-select__wrapper) {
    height: 44px;
  }
  
  .slider-container {
    height: 44px;
    padding: 16px;
  }
  
  .table-card {
    height: calc(100vh - 320px);
  }
  
  .custom-pagination {
    flex-direction: column;
    gap: 12px;
    padding: 16px;
  }
  
  .pagination-controls {
    order: -1;
  }
  
  .pagination-size {
    padding: 8px 12px;
  }
  
  .size-option {
    width: 35px;
    height: 32px;
    font-size: 12px;
  }
  
  .page-btn {
    width: 35px;
    height: 35px;
  }
  
  .page-number {
    width: 32px;
    height: 32px;
    font-size: 12px;
  }
  
  .skin-info {
    gap: 10px;
    padding: 8px 0;
  }
  
  .skin-image {
    width: 45px;
    height: 45px;
  }
  
  .skin-name {
    font-size: 13px;
    line-height: 1.3;
  }
  
  .skin-category {
    font-size: 11px;
  }
  
  .price-value, .diff-value {
    font-size: 13px;
    min-width: 65px;
    padding: 5px 10px;
  }
  
  .profit-tag {
    font-size: 12px;
    padding: 6px 12px;
    min-width: 65px;
  }
  
  .modern-table {
    font-size: 13px;
  }
  
  :deep(.el-table th) {
    height: 45px;
    font-size: 12px;
  }
  
  :deep(.el-table td) {
    height: 55px;
  }
}

@media (max-width: 480px) {
  .home-page {
    padding: 8px 8px 20px 8px;
  }
  
  .filter-section {
    padding: 16px 12px;
    margin-bottom: 16px;
    border-radius: 16px;
  }
  
  .filter-label {
    font-size: 13px;
  }
  
  .filter-input,
  .filter-select {
    height: 40px;
  }
  
  :deep(.filter-input .el-input__wrapper) {
    height: 40px;
    padding: 0 12px;
  }
  
  :deep(.filter-select .el-select__wrapper) {
    height: 40px;
  }
  
  .slider-container {
    height: 40px;
    padding: 12px;
  }
  
  .range-display {
    font-size: 11px;
    min-width: 70px;
  }
  
  .update-time-wrapper {
    padding: 6px 12px;
    min-width: 160px;
  }
  
  .update-time-label {
    font-size: 10px;
  }
  
  .update-time-value {
    font-size: 10px;
  }
  
  .refresh-btn {
    height: 40px;
    font-size: 12px;
    max-width: 100px;
  }
  
  .table-card {
    height: calc(100vh - 300px);
    border-radius: 16px;
  }
  
  .custom-pagination {
    padding: 12px;
    gap: 8px;
  }
  
  .pagination-info, .page-info {
    font-size: 12px;
    padding: 8px 12px;
  }
  
  .size-label {
    font-size: 12px;
  }
  
  .size-option {
    width: 30px;
    height: 28px;
    font-size: 11px;
  }
  
  .page-btn {
    width: 32px;
    height: 32px;
    font-size: 14px;
  }
  
  .page-number {
    width: 28px;
    height: 28px;
    font-size: 11px;
  }
  
  .skin-info {
    gap: 8px;
    padding: 6px 0;
  }
  
  .skin-image {
    width: 40px;
    height: 40px;
  }
  
  .skin-name {
    font-size: 12px;
    line-height: 1.2;
  }
  
  .skin-category {
    font-size: 10px;
  }
  
  .price-value, .diff-value {
    font-size: 11px;
    min-width: 55px;
    padding: 4px 8px;
    border-radius: 8px;
  }
  
  .profit-tag {
    font-size: 10px;
    padding: 4px 8px;
    min-width: 55px;
    border-radius: 10px;
  }
  
  .modern-table {
    font-size: 11px;
  }
  
  :deep(.el-table th) {
    height: 40px;
    font-size: 11px;
    padding: 4px 8px;
  }
  
  :deep(.el-table td) {
    height: 50px;
    padding: 4px 8px;
  }
  
  /* 移动端表格水平滚动优化 */
  .table-container {
    overflow-x: auto;
    overflow-y: hidden;
  }
  
  .table-container::-webkit-scrollbar {
    height: 6px;
  }
  
  .table-container::-webkit-scrollbar-track {
    background: rgba(0, 0, 0, 0.05);
    border-radius: 3px;
  }
  
  .table-container::-webkit-scrollbar-thumb {
    background: linear-gradient(135deg, #1890ff, #40a9ff);
    border-radius: 3px;
  }
  
  /* 固定列在移动端的处理 */
  :deep(.el-table__fixed) {
    box-shadow: 2px 0 8px rgba(0, 0, 0, 0.1);
  }
  
  /* 表格列宽度优化 */
  :deep(.el-table-column--selection) {
    width: 50px !important;
  }
  
  :deep(.index-column) {
    width: 50px !important;
  }
}
</style>
