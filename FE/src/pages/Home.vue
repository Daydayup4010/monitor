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
            <div class="update-time-display" v-if="lastUpdateTime">
              <el-icon><Clock /></el-icon>
              <span>更新时间：{{ formatUpdateTime(lastUpdateTime) }}</span>
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
          :data="displayData"
          style="width: 100%"
          stripe
          size="large"
          class="modern-table"
          @sort-change="handleSortChange"
          :height="tableHeight"
          :max-height="tableHeight"
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
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useSkinStore } from '@/stores/skin'
import { formatTime, formatPrice, formatPercent, getProfitColor, debounce } from '@/utils'
import dayjs from 'dayjs'
import type { SkinItem } from '@/types'

const skinStore = useSkinStore()

// 更新时间
const lastUpdateTime = ref('')

// 搜索和筛选
const searchKeyword = ref('')
const selectedCategory = ref('')
const profitRateRange = ref([0, 100])

// 排序
const sortConfig = ref<{ prop: string; order: string }>({ prop: '', order: '' })

// 表格高度计算
const tableHeight = ref(500)

// 获取所有类别
const categories = computed(() => {
  const cats = new Set(skinStore.skinItems.map(item => item.category))
  return Array.from(cats).filter(Boolean)
})

// 筛选后的皮肤数据
const filteredSkins = computed(() => {
  let filtered = skinStore.skinItems

  // 按关键词搜索
  if (searchKeyword.value) {
    const keyword = searchKeyword.value.toLowerCase()
    filtered = filtered.filter(item =>
      item.name.toLowerCase().includes(keyword)
    )
  }

  // 按类别筛选
  if (selectedCategory.value) {
    filtered = filtered.filter(item => item.category === selectedCategory.value)
  }

  // 按利润率筛选
  const [minRate, maxRate] = profitRateRange.value
  filtered = filtered.filter(item => {
    const rate = item.profit_rate * 100
    return rate >= minRate && rate <= maxRate
  })

  return filtered
})

// 排序后的显示数据
const displayData = computed(() => {
  let data = [...filteredSkins.value]
  
  if (sortConfig.value.prop && sortConfig.value.order) {
    const { prop, order } = sortConfig.value
    data.sort((a: any, b: any) => {
      let aVal = a[prop]
      let bVal = b[prop]
      
      // 处理数字类型
      if (typeof aVal === 'number' && typeof bVal === 'number') {
        return order === 'ascending' ? aVal - bVal : bVal - aVal
      }
      
      // 处理字符串类型
      if (typeof aVal === 'string' && typeof bVal === 'string') {
        return order === 'ascending' 
          ? aVal.localeCompare(bVal) 
          : bVal.localeCompare(aVal)
      }
      
      return 0
    })
  }
  
  return data
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

// 处理排序
const handleSortChange = ({ column, prop, order }: any) => {
  sortConfig.value = {
    prop: prop || '',
    order: order || ''
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
  padding: 32px 48px;
  min-height: 100vh;
  position: relative;
  animation: fadeInUp 0.6s ease-out;
}

.home-page::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: 
    radial-gradient(circle at 10% 20%, rgba(24, 144, 255, 0.05) 0%, transparent 20%),
    radial-gradient(circle at 90% 80%, rgba(114, 46, 209, 0.05) 0%, transparent 20%);
  pointer-events: none;
  border-radius: 24px;
}

/* 筛选区域 */
.filter-section {
  margin-bottom: 40px;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(20px);
  border-radius: 24px;
  padding: 40px 48px;
  box-shadow: 
    0 8px 32px rgba(24, 144, 255, 0.08),
    0 4px 16px rgba(0, 0, 0, 0.04);
  border: 1px solid rgba(24, 144, 255, 0.1);
  position: relative;
  overflow: hidden;
  transition: all 0.3s ease;
  animation: slideInLeft 0.8s ease-out 0.1s both;
}

.filter-section::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: linear-gradient(135deg, rgba(24, 144, 255, 0.02) 0%, rgba(114, 46, 209, 0.02) 100%);
  pointer-events: none;
}

.filter-section:hover {
  transform: translateY(-3px);
  box-shadow: 
    0 12px 40px rgba(24, 144, 255, 0.12),
    0 6px 20px rgba(0, 0, 0, 0.06);
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
  border-radius: 16px;
  box-shadow: 0 4px 20px rgba(24, 144, 255, 0.1);
  transition: all 0.3s ease;
}

.filter-input:hover,
.filter-select:hover {
  box-shadow: 0 6px 25px rgba(24, 144, 255, 0.15);
  transform: translateY(-1px);
}

:deep(.filter-input .el-input__wrapper) {
  height: 48px;
  padding: 0 16px;
  border: 2px solid rgba(24, 144, 255, 0.1);
  background: rgba(255, 255, 255, 0.95);
  border-radius: 16px;
  transition: all 0.3s ease;
  box-shadow: none;
}

:deep(.filter-input .el-input__wrapper:hover) {
  border-color: rgba(24, 144, 255, 0.3);
  background: rgba(255, 255, 255, 1);
}

:deep(.filter-input .el-input__wrapper.is-focus) {
  border-color: #1890ff;
  box-shadow: 0 0 0 3px rgba(24, 144, 255, 0.1);
}

:deep(.filter-select .el-select__wrapper) {
  height: 48px;
  border: 2px solid rgba(24, 144, 255, 0.1);
  background: rgba(255, 255, 255, 0.95);
  border-radius: 16px;
  transition: all 0.3s ease;
  box-shadow: none;
}

:deep(.filter-select .el-select__wrapper:hover) {
  border-color: rgba(24, 144, 255, 0.3);
  background: rgba(255, 255, 255, 1);
}

:deep(.filter-select .el-select__wrapper.is-focused) {
  border-color: #1890ff;
  box-shadow: 0 0 0 3px rgba(24, 144, 255, 0.1);
}

.slider-container {
  background: rgba(255, 255, 255, 0.95);
  border: 2px solid rgba(24, 144, 255, 0.1);
  border-radius: 16px;
  padding: 20px;
  transition: all 0.3s ease;
  height: 48px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 20px;
  max-width: 500px;
}

.slider-container:hover {
  border-color: rgba(24, 144, 255, 0.3);
  background: rgba(255, 255, 255, 1);
  transform: translateY(-1px);
  box-shadow: 0 6px 25px rgba(24, 144, 255, 0.15);
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

.update-time-display {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  font-size: 13px;
  color: #666;
  margin-bottom: 12px;
  padding: 8px 16px;
  background: rgba(24, 144, 255, 0.05);
  border-radius: 8px;
  border: 1px solid rgba(24, 144, 255, 0.1);
  transition: all 0.3s ease;
}

.update-time-display:hover {
  background: rgba(24, 144, 255, 0.1);
  border-color: rgba(24, 144, 255, 0.2);
}

.update-time-display .el-icon {
  color: #1890ff;
}

.refresh-btn {
  height: 48px;
  border-radius: 12px;
  font-weight: 600;
  background: linear-gradient(135deg, #1890ff, #40a9ff);
  border: none;
  box-shadow: 0 4px 15px rgba(24, 144, 255, 0.3);
  transition: all 0.3s ease;
  width: 100%;
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
  height: calc(100vh - 280px);
  border-radius: 24px;
  overflow: hidden;
  background: rgba(255, 255, 255, 0.98);
  backdrop-filter: blur(20px);
  box-shadow: 
    0 8px 32px rgba(24, 144, 255, 0.08),
    0 4px 16px rgba(0, 0, 0, 0.04);
  border: 1px solid rgba(24, 144, 255, 0.1);
  position: relative;
  transition: all 0.3s ease;
  animation: slideInRight 0.8s ease-out 0.2s both;
  display: flex;
  flex-direction: column;
}

.table-card::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 2px;
  background: linear-gradient(90deg, 
    transparent 0%, 
    #1890ff 25%, 
    #40a9ff 50%, 
    #1890ff 75%, 
    transparent 100%
  );
  background-size: 200% 100%;
  animation: tableShimmer 4s ease-in-out infinite;
}

@keyframes tableShimmer {
  0%, 100% { background-position: -200% 0; }
  50% { background-position: 200% 0; }
}

.table-card:hover {
  transform: translateY(-2px);
  box-shadow: 
    0 12px 40px rgba(24, 144, 255, 0.12),
    0 6px 20px rgba(0, 0, 0, 0.06);
}

.table-container {
  flex: 1;
  border-radius: 20px;
  position: relative;
  overflow: hidden;
  min-height: 0;
}

.modern-table {
  font-size: 16px;
  border-radius: 20px;
  position: relative;
}

.modern-table::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 4px;
  background: linear-gradient(90deg, 
    #1890ff 0%, 
    #40a9ff 25%, 
    #52c41a 50%, 
    #faad14 75%, 
    #722ed1 100%
  );
  z-index: 10;
}

.loading-container {
  min-height: 400px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, 
    rgba(24, 144, 255, 0.03) 0%, 
    rgba(114, 46, 209, 0.02) 50%,
    rgba(24, 144, 255, 0.03) 100%
  );
  border-radius: 20px;
  position: relative;
  overflow: hidden;
}

.loading-container::before {
  content: '';
  position: absolute;
  top: 0;
  left: -100%;
  width: 100%;
  height: 100%;
  background: linear-gradient(90deg, 
    transparent, 
    rgba(24, 144, 255, 0.1), 
    transparent
  );
  animation: loadingShimmer 2s ease-in-out infinite;
}

@keyframes loadingShimmer {
  0% { left: -100%; }
  100% { left: 100%; }
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
  border: 3px solid rgba(24, 144, 255, 0.2);
  border-radius: 10px;
  overflow: hidden;
  box-shadow: 
    0 6px 16px rgba(24, 144, 255, 0.15),
    0 3px 8px rgba(0, 0, 0, 0.1);
  transition: all 0.3s ease;
  position: relative;
  background: rgba(248, 250, 252, 0.8);
}

.skin-image:hover {
  border-color: #1890ff;
  transform: scale(1.08) rotate(1deg);
  box-shadow: 
    0 8px 24px rgba(24, 144, 255, 0.3),
    0 4px 12px rgba(0, 0, 0, 0.15);
}

.skin-avatar {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: all 0.3s ease;
  border-radius: 6px;
}

.skin-avatar:hover {
  transform: scale(1.1);
}

.skin-image::before {
  content: '';
  position: absolute;
  top: 0;
  left: -100%;
  width: 100%;
  height: 100%;
  background: linear-gradient(90deg, 
    transparent, 
    rgba(255, 255, 255, 0.4), 
    transparent
  );
  transition: left 0.5s ease;
  z-index: 1;
}

.skin-info:hover .skin-image::before {
  left: 100%;
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
  background: linear-gradient(135deg, rgba(24, 144, 255, 0.12), rgba(64, 169, 255, 0.08));
  padding: 8px 16px;
  border-radius: 12px;
  display: inline-block;
  min-width: 90px;
  text-align: center;
  border: 1px solid rgba(24, 144, 255, 0.2);
  transition: all 0.3s ease;
  position: relative;
  overflow: hidden;
}

.price-value::before {
  content: '';
  position: absolute;
  top: 0;
  left: -100%;
  width: 100%;
  height: 100%;
  background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.3), transparent);
  transition: left 0.5s ease;
}

.price-cell:hover .price-value::before {
  left: 100%;
}

.price-cell:hover .price-value {
  transform: scale(1.05);
  box-shadow: 0 4px 12px rgba(24, 144, 255, 0.2);
}

.diff-value {
  font-size: 16px;
  font-weight: 700;
  padding: 8px 16px;
  border-radius: 12px;
  display: inline-block;
  min-width: 90px;
  text-align: center;
  background: linear-gradient(135deg, rgba(0, 0, 0, 0.08), rgba(0, 0, 0, 0.04));
  border: 1px solid rgba(0, 0, 0, 0.1);
  transition: all 0.3s ease;
  position: relative;
  overflow: hidden;
}

.diff-value::before {
  content: '';
  position: absolute;
  top: 0;
  left: -100%;
  width: 100%;
  height: 100%;
  background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.4), transparent);
  transition: left 0.5s ease;
}

.diff-cell:hover .diff-value::before {
  left: 100%;
}

.diff-cell:hover .diff-value {
  transform: scale(1.05);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.profit-tag {
  font-size: 15px;
  font-weight: 700;
  padding: 10px 20px;
  border-radius: 14px;
  min-width: 80px;
  text-align: center;
  transition: all 0.3s ease;
  position: relative;
  overflow: hidden;
  border: 1px solid currentColor;
  opacity: 0.9;
}

.profit-tag::before {
  content: '';
  position: absolute;
  top: 0;
  left: -100%;
  width: 100%;
  height: 100%;
  background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.3), transparent);
  transition: left 0.6s ease;
}

.profit-cell:hover .profit-tag::before {
  left: 100%;
}

.profit-cell:hover .profit-tag {
  transform: scale(1.08);
  opacity: 1;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.15);
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
  background: linear-gradient(135deg, 
    rgba(248, 250, 255, 0.95) 0%, 
    rgba(230, 247, 255, 0.95) 50%,
    rgba(240, 249, 255, 0.95) 100%
  );
  font-weight: 700;
  font-size: 15px;
  color: #1890ff;
  border-bottom: 3px solid #1890ff;
  height: 65px;
  position: relative;
  transition: all 0.3s ease;
}

:deep(.el-table th::after) {
  content: '';
  position: absolute;
  bottom: -3px;
  left: 0;
  width: 100%;
  height: 3px;
  background: linear-gradient(90deg, #1890ff, #40a9ff, #1890ff);
  background-size: 200% 100%;
  animation: headerShine 3s ease-in-out infinite;
}

@keyframes headerShine {
  0%, 100% { background-position: -200% 0; }
  50% { background-position: 200% 0; }
}

:deep(.el-table td) {
  border-bottom: 1px solid rgba(24, 144, 255, 0.08);
  height: 85px;
  transition: all 0.2s ease;
  position: relative;
}

:deep(.el-table--striped .el-table__body tr.el-table__row--striped td) {
  background: linear-gradient(135deg, 
    rgba(24, 144, 255, 0.02) 0%, 
    rgba(114, 46, 209, 0.01) 100%
  );
}

:deep(.el-table__body tr:hover td) {
  background: linear-gradient(135deg, 
    rgba(24, 144, 255, 0.08) 0%, 
    rgba(64, 169, 255, 0.06) 100%
  ) !important;
  transform: scale(1.002);
  box-shadow: 0 2px 8px rgba(24, 144, 255, 0.1);
}

:deep(.el-table__body tr) {
  transition: all 0.3s ease;
  cursor: pointer;
}

:deep(.el-table__body tr:hover) {
  transform: translateY(-1px);
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
  }
  
  .refresh-btn {
    max-width: 200px;
  }
  
  .filter-section {
    padding: 28px 32px;
  }
}

@media (max-width: 768px) {
  .home-page {
    padding: 20px 24px;
  }
  
  .filter-row {
    grid-template-columns: 1fr;
    gap: 20px;
  }
  
  .filter-item {
    width: 100%;
  }
  
  .custom-pagination {
    flex-direction: column;
    gap: 16px;
    padding: 20px;
  }
  
  .pagination-controls {
    order: -1;
  }
  
  .skin-info {
    gap: 12px;
  }
  
  .skin-name {
    font-size: 14px;
  }
  
  .price-value, .diff-value {
    font-size: 14px;
    min-width: 70px;
    padding: 6px 12px;
  }
  
  .modern-table {
    font-size: 14px;
  }
  
  :deep(.el-table th) {
    height: 50px;
    font-size: 13px;
  }
  
  :deep(.el-table td) {
    height: 60px;
  }
}
</style>
