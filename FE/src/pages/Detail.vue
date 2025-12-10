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
          <div class="goods-meta">
            <span class="meta-item">
              <span class="meta-label">品质｜</span>
              <span class="meta-value" :style="{ color: getRarityColor(goodsDetail.rarityName) }">{{ goodsDetail.rarityName || '-' }}</span>
            </span>
            <span class="meta-divider"></span>
            <span class="meta-item">
              <span class="meta-label">类别｜</span>
              <span class="meta-value">{{ goodsDetail.qualityName || '-' }}</span>
            </span>
          </div>
        </div>
      </div>

      <!-- 主要内容区域 -->
      <div class="content-wrapper">
        <!-- 左侧: 价格信息和市场对比 -->
        <div class="left-panel">
          <!-- 价格涨幅卡片 -->
          <div class="price-change-card" v-if="goodsDetail?.priceChange?.length">
            <div class="price-change-nav">
              <button 
                class="nav-btn nav-btn-left" 
                :disabled="priceChangeIndex === 0"
                @click="priceChangeIndex = Math.max(0, priceChangeIndex - 1)"
              >
                <el-icon><ArrowLeft /></el-icon>
              </button>
              <div class="price-change-wrapper">
                <div class="price-change-list" :style="{ transform: `translateX(-${priceChangeIndex * 197}px)` }">
                  <div 
                    class="price-change-item" 
                    :class="{ 'is-up': item.isUp }"
                    v-for="item in goodsDetail.priceChange" 
                    :key="item.label"
                  >
                    <div class="change-label">
                      {{ item.label }} 
                      <span class="change-icon" :class="item.isUp ? 'icon-up' : 'icon-down'">
                        {{ item.isUp ? '▲' : '▼' }}
                      </span>
                    </div>
                    <div class="change-value" :class="item.isUp ? 'text-red' : 'text-green'">
                      <span class="price-diff">¥{{ formatPriceDiff(item.priceDiff) }}</span>
                      <span class="change-rate"> ({{ item.changeRate >= 0 ? '+' : '' }}{{ item.changeRate.toFixed(2) }}%)</span>
                    </div>
                  </div>
                </div>
              </div>
              <button 
                class="nav-btn nav-btn-right" 
                :disabled="priceChangeIndex >= (goodsDetail?.priceChange?.length || 0) - 2"
                @click="priceChangeIndex = Math.min((goodsDetail?.priceChange?.length || 0) - 2, priceChangeIndex + 1)"
              >
                <el-icon><ArrowRight /></el-icon>
              </button>
            </div>
          </div>

          <!-- 市场对比卡片 -->
          <div class="market-compare-card">
            <div class="card-title">市场对比</div>
            <div class="market-list">
              <div 
                class="market-item" 
                v-for="platform in goodsDetail?.platformList" 
                :key="platform.platform"
              >
                <div class="market-item-header">
                  <a :href="platform.link" target="_blank" class="platform-link">
                    <img :src="getPlatformIcon(platform.platform)" class="platform-icon" />
                    <span class="platform-name">{{ platform.platformName }}</span>
                  </a>
                  <div class="mini-chart-container" :ref="el => setMiniChartRef(platform.platform, el)"></div>
                </div>
                <div class="market-item-price">
                  <a :href="platform.link" target="_blank" class="price-link">
                    <span class="price-value">¥ {{ platform.sellPrice.toFixed(2) }}</span>
                    <span class="lowest-badge" v-if="isLowestPrice(platform)">底</span>
                  </a>
                </div>
                <div class="market-item-footer">
                  <span class="sell-count">在售：{{ platform.sellCount }}</span>
                  <span class="update-time">{{ formatUpdateTime(platform.updateTime) }}</span>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- 右侧: 饰品走势 -->
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
import { dataApi, type GoodsDetailResponse, type GoodsPlatformInfo } from '@/api'
import { ShoppingCart, Loading, ArrowLeft, ArrowRight } from '@element-plus/icons-vue'
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
const priceChangeIndex = ref(0)  // 价格涨幅卡片当前索引
let chartInstance: echarts.ECharts | null = null

// 迷你图表引用
const miniChartRefs = ref<Record<string, HTMLElement | null>>({})
const miniChartInstances: Record<string, echarts.ECharts> = {}

const setMiniChartRef = (platform: string, el: any) => {
  miniChartRefs.value[platform] = el as HTMLElement
}

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

// 获取品质对应的颜色
const getRarityColor = (rarity: string) => {
  const colorMap: Record<string, string> = {
    '隐秘': '#8B0000',     // 深红色
    '保密': '#FF69B4',     // 粉色
    '受限': '#8B008B',     // 紫色
    '军规级': '#00008B',   // 深蓝色
    '工业级': '#4169E1',   // 浅蓝色
    '消费级': '#808080',   // 灰色
  }
  return colorMap[rarity] || '#333333'  // 默认黑色
}

// 时间范围选项
const timeRanges = [
  { value: 30, label: '近1个月' },
  { value: 90, label: '近3个月' },
  { value: 180, label: '近6个月' },
  { value: 365, label: '近1年' },
]

// 判断是否是最低价
const isLowestPrice = (platform: GoodsPlatformInfo) => {
  if (!goodsDetail.value?.platformList) return false
  const prices = goodsDetail.value.platformList
    .filter(p => p.sellPrice > 0)
    .map(p => p.sellPrice)
  if (prices.length === 0) return false
  const minPrice = Math.min(...prices)
  return platform.sellPrice === minPrice && platform.sellPrice > 0
}

// 格式化价格差（下跌带负号，上涨不带+号）
const formatPriceDiff = (priceDiff: number) => {
  if (priceDiff < 0) {
    return priceDiff.toFixed(2)  // 负数自带负号
  }
  return priceDiff.toFixed(2)  // 正数不带+号
}

// 格式化更新时间
const formatUpdateTime = (timestamp: number) => {
  if (!timestamp) return ''
  const now = Date.now()
  const diff = now - timestamp * 1000
  const minutes = Math.floor(diff / 60000)
  const hours = Math.floor(diff / 3600000)
  
  if (minutes < 1) return '刚刚'
  if (minutes < 60) return `${minutes} 分钟前`
  if (hours < 24) return `${hours} 小时前`
  return `${Math.floor(hours / 24)} 天前`
}

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
      // 设置页面标题为饰品名称
      document.title = res.data.name || '饰品详情'
      await nextTick()
      initChart()
      initMiniCharts()
    }
  } catch (error) {
    console.error('获取商品详情失败:', error)
  } finally {
    loading.value = false
  }
}

// 初始化迷你图表
const initMiniCharts = () => {
  if (!goodsDetail.value?.platformList) return

  goodsDetail.value.platformList.forEach(platform => {
    const el = miniChartRefs.value[platform.platform]
    if (!el) return

    // 获取该平台7天的价格数据
    const historyData = goodsDetail.value?.priceHistory[platform.platform] || []
    const last7Days = historyData.slice(-7)
    const prices = last7Days.map(item => item.sellPrice)

    if (prices.length === 0) return

    // 销毁旧实例
    if (miniChartInstances[platform.platform]) {
      miniChartInstances[platform.platform].dispose()
    }

    const chart = echarts.init(el)
    miniChartInstances[platform.platform] = chart

    // 判断涨跌趋势
    const isUp = prices.length >= 2 && prices[prices.length - 1] >= prices[0]
    const color = isUp ? '#F56C6C' : '#0DAB62'

    const option: echarts.EChartsOption = {
      tooltip: {
        trigger: 'axis',
        backgroundColor: 'rgba(255, 255, 255, 0.95)',
        borderColor: '#eee',
        borderWidth: 1,
        padding: [8, 12],
        textStyle: {
          color: '#333',
          fontSize: 12,
        },
        formatter: (params: any) => {
          const data = params[0]
          const date = data.axisValue
          const price = data.value
          return `<div style="font-size:12px;">
            <div style="color:#999;margin-bottom:4px;">${date}</div>
            <div style="color:${color};font-weight:500;">¥${price.toFixed(2)}</div>
          </div>`
        },
        axisPointer: {
          type: 'line',
          lineStyle: {
            color: '#ddd',
            width: 1,
          },
        },
      },
      grid: {
        left: 0,
        right: 0,
        top: 0,
        bottom: 0,
      },
      xAxis: {
        type: 'category',
        show: false,
        data: last7Days.map(item => item.date),
      },
      yAxis: {
        type: 'value',
        show: false,
        min: Math.min(...prices) * 0.98,
        max: Math.max(...prices) * 1.02,
      },
      series: [{
        type: 'line',
        data: prices,
        smooth: true,
        symbol: 'circle',
        symbolSize: 4,
        showSymbol: false,
        lineStyle: {
          color: color,
          width: 1.5,
        },
        itemStyle: {
          color: color,
        },
        emphasis: {
          scale: true,
          itemStyle: {
            borderWidth: 2,
            borderColor: '#fff',
          },
        },
        areaStyle: {
          color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
            { offset: 0, color: isUp ? 'rgba(245, 108, 108, 0.3)' : 'rgba(13, 171, 98, 0.3)' },
            { offset: 1, color: isUp ? 'rgba(245, 108, 108, 0.05)' : 'rgba(13, 171, 98, 0.05)' },
          ]),
        },
      }],
    }

    chart.setOption(option)
  })
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
      left: 'center',
      itemWidth: 18,
      itemHeight: 7,
      itemGap: 30,
      textStyle: {
        fontSize: 14,
        color: '#666',
      },
    },
    grid: {
      left: 80,
      right: 100,
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
  Object.values(miniChartInstances).forEach(chart => chart?.resize())
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
  Object.values(miniChartInstances).forEach(chart => chart?.dispose())
})
</script>

<style scoped>
.detail-page {
  padding: 20px;
  background: #f5f7fa;
  min-height: 100vh;
}

.detail-container {
  max-width: 1400px;
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

.goods-meta {
  display: flex;
  align-items: center;
  gap: 16px;
  font-size: 14px;
}

.meta-item {
  display: flex;
  align-items: center;
}

.meta-label {
  color: #999;
}

.meta-value {
  font-weight: 500;
}

.meta-divider {
  width: 1px;
  height: 14px;
  background: #e0e0e0;
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

/* 左侧面板 */
.left-panel {
  width: 486px;
  flex-shrink: 0;
  display: flex;
  flex-direction: column;
  gap: 20px;
}

/* 价格涨幅卡片 */
.price-change-card {
  background: #fff;
  border-radius: 10px;
  padding: 15px 10px;
}

.price-change-nav {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 6px;
}

.nav-btn {
  width: 24px;
  height: 24px;
  border: none;
  background: transparent;
  color: #bbb;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  transition: all 0.2s;
  flex-shrink: 0;
}

.nav-btn:hover:not(:disabled) {
  background: #f5f5f5;
  color: #333;
}

.nav-btn:disabled {
  color: #ddd;
  cursor: not-allowed;
}

.nav-btn .el-icon {
  font-size: 18px;
}

.price-change-wrapper {
  width: 382px;  /* 两个卡片: 185*2 + 间距12 */
  flex-shrink: 0;
  overflow: hidden;
}

.price-change-list {
  display: flex;
  gap: 12px;
  transition: transform 0.3s ease;
}

.price-change-item {
  flex-shrink: 0;
  width: 185px;
  height: 90px;
  background: #F5F8FD;
  border-radius: 10px;
  padding: 16px;
  text-align: center;
  display: flex;
  flex-direction: column;
  justify-content: center;
  box-sizing: border-box;
}

.price-change-item.is-up {
  background: #FFF5F5;
}

.change-label {
  font-size: 16px;
  color: #666;
  margin-bottom: 10px;
}

.change-icon {
  font-size: 14px;
  margin-left: 6px;
}

.change-icon.icon-up {
  color: #F56C6C;
}

.change-icon.icon-down {
  color: #0DAB62;
}

.change-value {
  font-size: 16px;
  font-weight: 500;
}

.price-diff {
  font-weight: 600;
}

.text-red {
  color: #F56C6C;
}

.text-green {
  color: #0DAB62;
}

/* 市场对比卡片 */
.market-compare-card {
  background: #fff;
  border-radius: 10px;
  padding: 20px;
  flex: 1;
}

.card-title {
  font-size: 16px;
  font-weight: 500;
  color: #333;
  margin-bottom: 16px;
}

.market-list {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
}

.market-item {
  width: calc(50% - 6px);
  background: #F5F8FD;
  border-radius: 10px;
  padding: 12px;
}

.market-item-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 8px;
}

.platform-link {
  display: flex;
  align-items: center;
  gap: 6px;
  text-decoration: none;
}

.platform-icon {
  width: 18px;
  height: 18px;
  object-fit: contain;
}

.platform-name {
  font-size: 14px;
  color: #999;
}

.mini-chart-container {
  width: 60px;
  height: 30px;
}

.market-item-price {
  margin-bottom: 8px;
}

.price-link {
  display: flex;
  align-items: center;
  gap: 6px;
  text-decoration: none;
}

.price-value {
  font-size: 18px;
  font-weight: 500;
  color: #ff6b00;
  white-space: nowrap;
}

.lowest-badge {
  display: inline-block;
  width: 20px;
  height: 20px;
  line-height: 20px;
  text-align: center;
  background: #0DAB62;
  color: #fff;
  font-size: 12px;
  border-radius: 4px;
}

.market-item-footer {
  display: flex;
  justify-content: space-between;
  font-size: 12px;
  color: #999;
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
