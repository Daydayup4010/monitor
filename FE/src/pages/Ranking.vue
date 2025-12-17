<template>
  <div class="ranking-page">
    <div class="card">
      <div class="card-title">
        <el-icon :size="20" color="#1890ff"><TrendCharts /></el-icon>
        饰品榜单
      </div>

      <!-- Tab 切换 -->
      <div class="tab-wrapper">
        <div class="tab-content">
          <div 
            class="tab-item" 
            :class="{ active: activeTab === 'increase' }"
            @click="switchTab('increase')"
          >在售价涨幅榜</div>
          <div 
            class="tab-item" 
            :class="{ active: activeTab === 'decrease' }"
            @click="switchTab('decrease')"
          >在售价跌幅榜</div>
        </div>
      </div>

      <!-- 数据表格 -->
      <div class="table-wrapper">
        <el-table
          :data="rankingData"
          v-loading="loading"
          style="width: 100%"
          :row-class-name="getRowClassName"
          :show-header="false"
        >
          <el-table-column type="index" label="#" width="60">
            <template #default="{ $index }">
              <span :class="getRankClass($index)">{{ $index + 1 }}</span>
            </template>
          </el-table-column>
          
          <el-table-column label="饰品" class-name="skin-column">
            <template #default="{ row }">
              <div class="skin-flex">
                <div class="item-img-box">
                  <img :src="row.iconUrl" @error="handleImageError" alt="饰品" />
                </div>
                <div class="skin-info-box">
                  <div class="skin-name clickable" @click="goToDetail(row.marketHashName)">{{ row.name }}</div>
                  <div class="skin-price">¥{{ formatPrice(row.todayPrice) }}</div>
                  <div class="skin-meta">
                    <span class="meta-item">
                      价格涨幅: 
                      <span :class="getChangeClass(row.increaseRate1D)">{{ formatChangeRate(row.increaseRate1D) }}</span>
                      <el-tooltip placement="right" :show-after="200">
                        <template #content>
                          <div class="tooltip-content">
                            <div class="tooltip-header">当前价格: ¥{{ formatPrice(row.todayPrice) }}</div>
                            <div class="tooltip-row">与1天前对比 <span class="tooltip-price">¥{{ formatPrice(row.yesterdayPrice) }}</span> <span :class="getChangeClass(row.increaseRate1D)">{{ formatChangeRate(row.increaseRate1D) }}</span></div>
                            <div class="tooltip-row" v-if="row.price3DaysAgo !== null">与3天前对比 <span class="tooltip-price">¥{{ formatPrice(row.price3DaysAgo) }}</span> <span :class="getChangeClass(row.increaseRate3D)">{{ formatChangeRate(row.increaseRate3D) }}</span></div>
                            <div class="tooltip-row" v-if="row.price7DaysAgo !== null">与7天前对比 <span class="tooltip-price">¥{{ formatPrice(row.price7DaysAgo) }}</span> <span :class="getChangeClass(row.increaseRate7D)">{{ formatChangeRate(row.increaseRate7D) }}</span></div>
                            <div class="tooltip-row" v-if="row.price30DaysAgo !== null">与30天前对比 <span class="tooltip-price">¥{{ formatPrice(row.price30DaysAgo) }}</span> <span :class="getChangeClass(row.increaseRate30D)">{{ formatChangeRate(row.increaseRate30D) }}</span></div>
                          </div>
                        </template>
                        <el-icon class="meta-icon"><InfoFilled /></el-icon>
                      </el-tooltip>
                    </span>
                    <span class="meta-item">
                      在售数涨幅: 
                      <span :class="getChangeClass(row.sellCountRate1D)">{{ formatChangeRate(row.sellCountRate1D) }}</span>
                      <el-tooltip placement="right" :show-after="200">
                        <template #content>
                          <div class="tooltip-content">
                            <div class="tooltip-header">当前在售数: {{ row.todaySellCount }}</div>
                            <div class="tooltip-row">与1天前对比 <span class="tooltip-price">{{ row.yesterdaySellCount }}</span> <span :class="getChangeClass(row.sellCountRate1D)">{{ formatChangeRate(row.sellCountRate1D) }}</span></div>
                            <div class="tooltip-row" v-if="row.sellCount3DaysAgo !== null">与3天前对比 <span class="tooltip-price">{{ row.sellCount3DaysAgo }}</span> <span :class="getChangeClass(row.sellCountRate3D)">{{ formatChangeRate(row.sellCountRate3D) }}</span></div>
                            <div class="tooltip-row" v-if="row.sellCount7DaysAgo !== null">与7天前对比 <span class="tooltip-price">{{ row.sellCount7DaysAgo }}</span> <span :class="getChangeClass(row.sellCountRate7D)">{{ formatChangeRate(row.sellCountRate7D) }}</span></div>
                            <div class="tooltip-row" v-if="row.sellCount30DaysAgo !== null">与30天前对比 <span class="tooltip-price">{{ row.sellCount30DaysAgo }}</span> <span :class="getChangeClass(row.sellCountRate30D)">{{ formatChangeRate(row.sellCountRate30D) }}</span></div>
                          </div>
                        </template>
                        <el-icon class="meta-icon"><InfoFilled /></el-icon>
                      </el-tooltip>
                    </span>
                  </div>
                </div>
              </div>
            </template>
          </el-table-column>
          
          <el-table-column label="" width="280">
            <template #default="{ row }">
              <div class="rate-cell">
                <div class="rate-top" :class="getChangeClass(row.increaseRate1D)">
                  <span class="rate-percent">{{ formatChangeRate(row.increaseRate1D) }}</span>
                  <span class="rate-diff">({{ row.priceChange >= 0 ? '+' : '' }}¥{{ formatPrice(Math.abs(row.priceChange)) }})</span>
                </div>
                <div class="rate-bottom">
                  ¥{{ formatPrice(row.yesterdayPrice) }} <span class="arrow">&gt;&gt;</span> ¥{{ formatPrice(row.todayPrice) }}
                </div>
              </div>
            </template>
          </el-table-column>
        </el-table>
      </div>

      <!-- 数据说明 -->
      <div class="data-info">
        <el-icon><InfoFilled /></el-icon>
        <span>数据来源：悠悠有品平台 | 每日凌晨更新</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { TrendCharts, InfoFilled } from '@element-plus/icons-vue'
import { dataApi, type PriceIncreaseItem } from '@/api'
import { formatPrice } from '@/utils'
import { showMessage } from '@/utils/message'

const router = useRouter()
const activeTab = ref<'increase' | 'decrease'>('increase')
const loading = ref(false)
const rankingData = ref<PriceIncreaseItem[]>([])

// 获取排行数据
const fetchRankingData = async () => {
  loading.value = true
  try {
    const isDesc = activeTab.value === 'increase' // 涨幅榜降序，跌幅榜升序
    const response = await dataApi.getPriceIncrease({ is_desc: isDesc, limit: 300 })
    if (response.code === 1) {
      rankingData.value = response.data || []
    } else {
      showMessage.error('获取数据失败')
    }
  } catch (error) {
    console.error('获取排行数据失败:', error)
  } finally {
    loading.value = false
  }
}

// 切换 Tab
const switchTab = (tab: 'increase' | 'decrease') => {
  if (activeTab.value !== tab) {
    activeTab.value = tab
    fetchRankingData()
  }
}

// 计算价格差
const getPriceDiff = (todayPrice: number, oldPrice: number | null) => {
  if (oldPrice === null) return 0
  return todayPrice - oldPrice
}

// 格式化涨跌幅
const formatChangeRate = (rate: number | null) => {
  if (rate === null) return '-'
  const sign = rate >= 0 ? '+' : ''
  return `${sign}${rate.toFixed(2)}%`
}

// 获取涨跌幅样式类
const getChangeClass = (rate: number | null) => {
  if (rate === null) return ''
  if (rate > 0) return 'change-up'
  if (rate < 0) return 'change-down'
  return 'change-neutral'
}

// 获取排名样式类
const getRankClass = (index: number) => {
  if (index === 0) return 'rank rank-1'
  if (index === 1) return 'rank rank-2'
  if (index === 2) return 'rank rank-3'
  return 'rank'
}

// 获取行样式类
const getRowClassName = ({ rowIndex }: { rowIndex: number }) => {
  if (rowIndex < 3) return 'top-row'
  return ''
}

// 图片加载失败处理
const handleImageError = (e: Event) => {
  const img = e.target as HTMLImageElement
  img.src = 'data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iNjAiIGhlaWdodD0iNjAiIHZpZXdCb3g9IjAgMCA2MCA2MCIgZmlsbD0ibm9uZSIgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj48cmVjdCB3aWR0aD0iNjAiIGhlaWdodD0iNjAiIGZpbGw9IiNGNUY1RjUiLz48dGV4dCB4PSIzMCIgeT0iMzUiIGZvbnQtc2l6ZT0iMTQiIGZpbGw9IiNCRkJGQkYiIHRleHQtYW5jaG9yPSJtaWRkbGUiPuaXoOWbvjwvdGV4dD48L3N2Zz4='
}

// 跳转到详情页
const goToDetail = (marketHashName: string) => {
  const routeData = router.resolve({
    path: '/app/detail',
    query: { market_hash_name: marketHashName }
  })
  window.open(routeData.href, '_blank')
}

onMounted(() => {
  fetchRankingData()
})
</script>

<style scoped>
/* 所有样式在unified.css中 */
.ranking-page {
  padding: 0;
  width: 1800px;
  max-width: 100%;
  margin: 0 auto;
  overflow-x: hidden;
}

.ranking-page .card {
  overflow: hidden;
}

.tab-wrapper {
  display: flex;
  justify-content: space-between;
  border-bottom: 1px solid #eee;
  padding-bottom: 10px;
  margin-bottom: 20px;
}

.tab-content {
  display: flex;
  height: 40px;
  gap: 10px;
  border-radius: 10px;
  background: #e6f4ff;
  padding: 3px;
  align-items: center;
}

.tab-item {
  cursor: pointer;
  padding: 7px 30px;
  border-radius: 10px;
  font-size: 14px;
  color: #333;
  transition: all 0.2s ease;
}

.tab-item:hover {
  color: #1890ff;
}

.tab-item.active {
  background: white;
  color: #1890ff;
}

/* 饰品 flex 布局 - 仿照竞品网站 */
.skin-flex {
  display: flex;
  align-items: center;
}

.item-img-box {
  width: 90px;
  height: 60px;
  margin-right: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 8px;
  padding: 4px;
  background: #f5f7fa;
  flex-shrink: 0;
}

.item-img-box img {
  max-width: 100%;
  max-height: 100%;
  object-fit: contain;
}

.skin-info-box {
  flex: 1;
}

.skin-name {
  font-size: 16px;
  color: #000;
  margin-bottom: 5px;
  line-height: 1.4;
}

.skin-name.clickable {
  cursor: pointer;
  transition: color 0.2s;
}

.skin-name.clickable:hover {
  color: #1890ff;
}

.skin-price {
  font-size: 16px;
  color: #FD802C;
  margin-bottom: 5px;
}

.skin-meta {
  color: #999;
  font-size: 14px;
}

.meta-item {
  margin-right: 20px;
}

.meta-icon {
  margin-left: 8px;
  color: #1890ff;
  cursor: pointer;
  vertical-align: middle;
  font-size: 16px;
  background: #e6f4ff;
  border-radius: 50%;
  padding: 2px;
}

.meta-icon:hover {
  color: #fff;
  background: #1890ff;
}

/* tooltip 样式 */
.tooltip-content {
  font-size: 14px;
  line-height: 2;
}

.tooltip-header {
  font-size: 15px;
  margin-bottom: 8px;
  padding-bottom: 8px;
  border-bottom: 1px solid rgba(255,255,255,0.2);
}

.tooltip-row {
  white-space: nowrap;
}

.tooltip-price {
  margin: 0 12px;
}

/* 右侧涨跌幅样式 */
.rate-cell {
  text-align: left;
  padding-left: 20px;
}

.rate-top {
  margin-bottom: 10px;
  white-space: nowrap;
}

.rate-percent {
  font-size: 20px;
}

.rate-diff {
  font-size: 15px;
  margin-left: 2px;
}

.rate-bottom {
  color: #999;
  font-size: 15px;
  white-space: nowrap;
}

.rate-bottom .arrow {
  margin: 0 5px;
}

/* 涨跌幅颜色 */
.change-up {
  color: #f5222d;
}

.change-down {
  color: #52c41a;
}

.change-neutral {
  color: #999;
}

/* 排名样式 */
.rank {
  display: inline-block;
  width: 24px;
  height: 24px;
  line-height: 24px;
  text-align: center;
  font-size: 16px;
  color: #333;
}

.rank-1,
.rank-2,
.rank-3 {
  width: 24px;
  height: 24px;
  background-size: contain;
  background-repeat: no-repeat;
  background-position: center;
  color: transparent;
}

.rank-1 {
  background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24'%3E%3Cpath fill='%23FFD700' d='M12 2L15.09 8.26L22 9.27L17 14.14L18.18 21.02L12 17.77L5.82 21.02L7 14.14L2 9.27L8.91 8.26L12 2Z'/%3E%3C/svg%3E");
}

.rank-2 {
  background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24'%3E%3Cpath fill='%23C0C0C0' d='M12 2L15.09 8.26L22 9.27L17 14.14L18.18 21.02L12 17.77L5.82 21.02L7 14.14L2 9.27L8.91 8.26L12 2Z'/%3E%3C/svg%3E");
}

.rank-3 {
  background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24'%3E%3Cpath fill='%23CD7F32' d='M12 2L15.09 8.26L22 9.27L17 14.14L18.18 21.02L12 17.77L5.82 21.02L7 14.14L2 9.27L8.91 8.26L12 2Z'/%3E%3C/svg%3E");
}

/* 数据说明 */
.data-info {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 16px;
  background: #f5f7fa;
  border-radius: 8px;
  margin-top: 16px;
  font-size: 14px;
  color: #666;
}

.data-info .el-icon {
  color: #1890ff;
}

/* 表格样式 - 覆盖全局样式 */
.table-wrapper {
  margin-top: 16px;
  overflow-x: hidden !important;
  border: none;
}

:deep(.el-table) {
  font-size: 14px;
  width: 100% !important;
}

:deep(.el-table__header-wrapper),
:deep(.el-table__body-wrapper) {
  overflow-x: hidden !important;
}

/* 右侧涨跌列不换行 */
.rate-cell {
  white-space: nowrap;
}

:deep(.el-table th) {
  background: #fafafa !important;
  color: #333;
  font-size: 14px;
  font-weight: normal;
  border-bottom: 1px solid #eee;
}

:deep(.el-table td) {
  border-bottom: 1px solid #f0f0f0;
}

:deep(.el-table__row) {
  height: 80px;
}

:deep(.el-table__row:hover > td) {
  background-color: #fafafa !important;
}
</style>
