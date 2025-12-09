<template>
  <div class="detail-page">
    <div class="detail-container">
      <!-- 商品基础信息 -->
      <div class="goods-header" v-if="goodsDetail">
        <div class="goods-image">
          <img :src="goodsDetail.iconUrl" :alt="goodsDetail.name" />
        </div>
        <div class="goods-info">
          <h1 class="goods-name">{{ goodsDetail.name }}</h1>
        </div>
      </div>

      <!-- 主要内容区域 -->
      <div class="content-wrapper">
        <!-- 左侧: 饰品走势 -->
        <div class="chart-card">
          <el-tabs v-model="activeTab" class="my-tabs">
            <el-tab-pane label="饰品走势" name="trend">
              <!-- 图表控制区 -->
              <div class="charts-menu">
                <div class="filters">
                  <!-- 平台选择 -->
                  <el-select v-model="selectedPlatform" class="platform-select" @change="handlePlatformChange">
                    <template #label="{ value }">
                      <div style="display: flex; align-items: center; gap: 8px;">
                        <img :src="getPlatformIcon(value)" style="width: 18px; height: 18px; object-fit: contain;" />
                        <span>{{ getPlatformLabel(value) }}</span>
                      </div>
                    </template>
                    <el-option
                      v-for="p in platforms"
                      :key="p.value"
                      :label="p.label"
                      :value="p.value"
                    >
                      <div style="display: flex; align-items: center; gap: 8px;">
                        <img :src="p.icon" style="width: 18px; height: 18px; object-fit: contain;" />
                        <span>{{ p.label }}</span>
                      </div>
                    </el-option>
                  </el-select>
                  
                  <!-- 时间范围选择 -->
                  <el-select v-model="selectedDays" class="time-select" @change="handleDaysChange">
                    <el-option
                      v-for="t in timeRanges"
                      :key="t.value"
                      :label="t.label"
                      :value="t.value"
                    />
                  </el-select>
                </div>
              </div>

              <!-- 图表名称 -->
              <div class="chart-title" v-if="goodsDetail">{{ goodsDetail.name }}</div>

              <!-- 图表区域 -->
              <div class="charts-area">
                <div ref="chartRef" class="chart-container"></div>
              </div>
            </el-tab-pane>

            <el-tab-pane label="当前在售" name="onsale">
              <div class="onsale-table" v-if="goodsDetail">
                <div class="table-header">
                  <div class="col-source">来源</div>
                  <div class="col-price">售价</div>
                  <div class="col-bidding">求购价</div>
                  <div class="col-sell">当前在售</div>
                  <div class="col-bid-count">当前求购</div>
                  <div class="col-action"></div>
                </div>
                <div class="table-body">
                  <div
                    class="table-row"
                    v-for="platform in goodsDetail.platformList"
                    :key="platform.platform"
                  >
                    <div class="col-source">
                      <img :src="getPlatformIcon(platform.platform)" class="platform-icon-img" />
                      <span class="platform-name">{{ platform.platformName }}</span>
                    </div>
                    <div class="col-price text-orange">¥ {{ platform.sellPrice.toFixed(2) }}</div>
                    <div class="col-bidding text-orange">¥ {{ platform.biddingPrice.toFixed(2) }}</div>
                    <div class="col-sell">{{ platform.sellCount }}件</div>
                    <div class="col-bid-count">{{ platform.biddingCount }}件</div>
                    <div class="col-action">
                      <a :href="platform.link" target="_blank" class="buy-btn" v-if="platform.link">
                        <el-icon><ShoppingCart /></el-icon> 前去购买
                      </a>
                    </div>
                  </div>
                </div>
              </div>
            </el-tab-pane>
          </el-tabs>
        </div>
      </div>
    </div>

    <!-- 加载状态 -->
    <div class="loading-mask" v-if="loading">
      <el-icon class="is-loading"><Loading /></el-icon>
      <span>加载中...</span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, watch, nextTick } from 'vue'
import { useRoute } from 'vue-router'
import { dataApi, type GoodsDetailResponse } from '@/api'
import { ShoppingCart, Loading } from '@element-plus/icons-vue'
import * as echarts from 'echarts'
import buffIcon from '@/assets/icons/buff.png'
import uuIcon from '@/assets/icons/uu.png'
import c5Icon from '@/assets/icons/c5.png'
import steamIcon from '@/assets/icons/steam.png'

const route = useRoute()

// 数据
const goodsDetail = ref<GoodsDetailResponse | null>(null)
const loading = ref(false)
const activeTab = ref('trend')
const selectedPlatform = ref('YOUPIN')
const selectedDays = ref(30)
const chartRef = ref<HTMLElement | null>(null)
let chartInstance: echarts.ECharts | null = null

// 平台选项
const platforms = [
  { value: 'YOUPIN', label: '悠悠', icon: uuIcon },
  { value: 'BUFF', label: 'BUFF', icon: buffIcon },
  { value: 'C5', label: 'C5GAME', icon: c5Icon },
  { value: 'STEAM', label: 'Steam', icon: steamIcon },
]

// 获取平台图标
const getPlatformIcon = (platform: string) => {
  const iconMap: Record<string, string> = {
    'YOUPIN': uuIcon,
    'BUFF': buffIcon,
    'C5': c5Icon,
    'STEAM': steamIcon
  }
  return iconMap[platform] || ''
}

// 获取平台显示名称
const getPlatformLabel = (platform: string) => {
  const labelMap: Record<string, string> = {
    'YOUPIN': '悠悠',
    'BUFF': 'BUFF',
    'C5': 'C5GAME',
    'STEAM': 'Steam'
  }
  return labelMap[platform] || platform
}

// 时间范围选项
const timeRanges = [
  { value: 30, label: '近1个月' },
  { value: 90, label: '近3个月' },
  { value: 180, label: '近6个月' },
  { value: 365, label: '近1年' },
]

// 获取商品详情
const fetchGoodsDetail = async () => {
  const marketHashName = route.query.market_hash_name as string
  if (!marketHashName) return

  loading.value = true
  try {
    const res = await dataApi.getGoodsDetail({
      market_hash_name: marketHashName,
      days: selectedDays.value,
    })
    if (res.code === 1 && res.data) {
      goodsDetail.value = res.data
      await nextTick()
      initChart()
    }
  } catch (error) {
    console.error('获取商品详情失败:', error)
  } finally {
    loading.value = false
  }
}

// 初始化图表
const initChart = () => {
  if (!chartRef.value || !goodsDetail.value) return

  if (chartInstance) {
    chartInstance.dispose()
  }

  chartInstance = echarts.init(chartRef.value)

  // 根据选中的平台获取历史数据
  const historyData = goodsDetail.value.priceHistory[selectedPlatform.value] || []
  const dates = historyData.map(item => item.date)
  const prices = historyData.map(item => item.sellPrice)
  const counts = historyData.map(item => item.sellCount)

  const option: echarts.EChartsOption = {
    tooltip: {
      trigger: 'axis',
      backgroundColor: '#fff',
      borderColor: '#fff',
      borderWidth: 1,
      padding: 10,
      textStyle: {
        color: '#666',
        fontSize: 14,
      },
      formatter: (params: any) => {
        const date = params[0].axisValue
        let html = `<div style="font-weight:500;margin-bottom:8px;">${date}</div>`
        params.forEach((item: any) => {
          const color = item.color
          const name = item.seriesName
          const value = item.seriesName === '价格' ? `¥${item.value}` : `${item.value}件`
          html += `<div style="display:flex;justify-content:space-between;align-items:center;margin:4px 0;">
            <span style="display:inline-block;width:10px;height:10px;border-radius:50%;background:${color};margin-right:8px;"></span>
            <span style="flex:1;">${name}</span>
            <span style="margin-left:20px;font-weight:500;">${value}</span>
          </div>`
        })
        return html
      },
    },
    legend: {
      data: ['价格', '在售数量'],
      top: 5,
      right: 100,
      itemWidth: 18,
      itemHeight: 7,
      textStyle: {
        fontSize: 14,
        color: '#666',
      },
    },
    grid: {
      left: 60,
      right: 60,
      top: 60,
      bottom: 60,
    },
    xAxis: {
      type: 'category',
      data: dates,
      axisLine: {
        lineStyle: { color: '#6E7079' },
      },
      axisTick: {
        show: true,
      },
      axisLabel: {
        color: '#979CAE',
        fontSize: 12,
        formatter: (value: string) => {
          // 简化日期显示
          return value.slice(5)
        },
      },
    },
    yAxis: [
      {
        type: 'value',
        name: '金额（¥）',
        nameTextStyle: {
          color: '#979CAE',
          fontSize: 14,
          padding: [0, 0, 10, 0],
        },
        position: 'left',
        axisLine: { show: false },
        axisTick: { show: false },
        axisLabel: {
          color: '#979CAE',
          fontSize: 12,
        },
        splitLine: {
          lineStyle: { color: '#EFF1F5' },
        },
      },
      {
        type: 'value',
        name: '在售数量（件）',
        nameTextStyle: {
          color: '#979CAE',
          fontSize: 14,
          padding: [0, 0, 10, 0],
        },
        position: 'right',
        axisLine: { show: false },
        axisTick: { show: false },
        axisLabel: {
          color: '#979CAE',
          fontSize: 12,
        },
        splitLine: { show: false },
      },
    ],
    series: [
      {
        name: '价格',
        type: 'line',
        yAxisIndex: 0,
        data: prices,
        smooth: true,
        symbol: 'circle',
        symbolSize: 6,
        showSymbol: false,
        lineStyle: {
          color: 'rgb(248, 118, 0)',
          width: 2,
        },
        itemStyle: {
          color: 'rgb(248, 118, 0)',
        },
        areaStyle: {
          color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
            { offset: 0, color: 'rgba(248, 118, 0, 0.3)' },
            { offset: 1, color: 'rgba(248, 118, 0, 0.05)' },
          ]),
        },
      },
      {
        name: '在售数量',
        type: 'line',
        yAxisIndex: 1,
        data: counts,
        smooth: true,
        symbol: 'circle',
        symbolSize: 6,
        showSymbol: false,
        lineStyle: {
          color: 'rgb(0, 152, 255)',
          width: 2,
        },
        itemStyle: {
          color: 'rgb(0, 152, 255)',
        },
      },
    ],
  }

  chartInstance.setOption(option)
}

// 平台变更 - 只需重新渲染图表，数据已包含所有平台
const handlePlatformChange = () => {
  initChart()
}

// 时间范围变更
const handleDaysChange = () => {
  fetchGoodsDetail()
}

// 窗口大小变化时调整图表
const handleResize = () => {
  chartInstance?.resize()
}

// 监听activeTab变化，切换到走势时重新渲染图表
watch(activeTab, (val) => {
  if (val === 'trend') {
    nextTick(() => {
      initChart()
    })
  }
})

onMounted(() => {
  fetchGoodsDetail()
  window.addEventListener('resize', handleResize)
})

onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
  chartInstance?.dispose()
})
</script>

<style scoped>
.detail-page {
  padding: 20px;
  background: #f5f7fa;
  min-height: 100vh;
}

.detail-container {
  max-width: 1200px;
  margin: 0 auto;
}

/* 商品头部 */
.goods-header {
  display: flex;
  align-items: center;
  gap: 24px;
  padding: 24px;
  background: #fff;
  border-radius: 10px;
  margin-bottom: 20px;
}

.goods-image {
  width: 120px;
  height: 90px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 10px;
}

.goods-image img {
  max-width: 100%;
  max-height: 100%;
  object-fit: contain;
}

.goods-info {
  flex: 1;
}

.goods-name {
  font-size: 20px;
  font-weight: 500;
  color: #333;
  margin: 0 0 8px 0;
}

.goods-hash {
  font-size: 14px;
  color: #999;
  margin: 0;
}

/* 内容区 */
.content-wrapper {
  display: flex;
  gap: 20px;
}

.chart-card {
  flex: 1;
  background: #fff;
  border-radius: 10px;
  padding: 20px;
}

/* 图表控制区 */
.charts-menu {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px 0;
  border-bottom: 1px solid #f0f0f0;
  margin-bottom: 16px;
}

.filters {
  display: flex;
  align-items: center;
  gap: 20px;
}

.platform-select {
  width: 130px;
}

.time-select {
  width: 120px;
}

.chart-title {
  font-size: 16px;
  color: #333;
  padding: 10px 0;
}

/* 图表区域 */
.charts-area {
  padding: 20px 0;
}

.chart-container {
  width: 100%;
  height: 500px;
}

/* 在售表格 */
.onsale-table {
  padding: 10px 0;
}

.table-header,
.table-row {
  display: flex;
  align-items: center;
  padding: 16px 0;
}

.table-header {
  font-size: 14px;
  color: #999;
  border-bottom: 1px solid #f0f0f0;
}

.table-row {
  border-bottom: 1px solid #f5f5f5;
}

.table-row:last-child {
  border-bottom: none;
}

.col-source {
  flex: 2;
  display: flex;
  align-items: center;
  gap: 8px;
}

.col-price,
.col-bidding {
  flex: 1.5;
  font-size: 16px;
}

.col-sell,
.col-bid-count {
  flex: 1;
  font-size: 14px;
  color: #666;
}

.col-action {
  flex: 1.5;
  text-align: right;
}

.platform-icon-img {
  width: 24px;
  height: 24px;
  object-fit: contain;
}

.platform-name {
  font-size: 16px;
  color: #262626;
}

.text-orange {
  color: #ff6b00;
}

.buy-btn {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 6px 16px;
  background: transparent;
  border: 1px solid #1890ff;
  border-radius: 6px;
  color: #1890ff;
  font-size: 14px;
  text-decoration: none;
  transition: all 0.2s;
}

.buy-btn:hover {
  background: #1890ff;
  color: #fff;
}

/* 加载状态 */
.loading-mask {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(255, 255, 255, 0.8);
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 12px;
  z-index: 1000;
}

.loading-mask .is-loading {
  font-size: 32px;
  color: #1890ff;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  from {
    transform: rotate(0deg);
  }
  to {
    transform: rotate(360deg);
  }
}

/* Tabs 样式覆盖 */
:deep(.el-tabs__header) {
  margin-bottom: 0;
}

:deep(.el-tabs__item) {
  font-size: 16px;
  color: #666;
}

:deep(.el-tabs__item.is-active) {
  color: #1890ff;
}

:deep(.el-tabs__active-bar) {
  background-color: #1890ff;
}
</style>

