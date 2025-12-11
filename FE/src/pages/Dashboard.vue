<template>
  <div class="dashboard-page">
    <!-- 标语区域 -->
    <section class="hero-section">
      <h1 class="hero-title">专注于CS2饰品搬砖</h1>
      <p class="hero-subtitle">提供专业的饰品数据变化分析、各平台搬砖比价数据</p>
    </section>

    <!-- 搜索框 -->
    <div class="search-section">
      <div class="search-box">
        <el-input
          v-model="searchKeyword"
          placeholder="搜索饰品名称..."
          size="large"
          clearable
          @input="handleSearchInput"
          @clear="clearSearch"
        >
          <template #prefix>
            <el-icon><Search /></el-icon>
          </template>
        </el-input>
        <!-- 搜索结果下拉 -->
        <div class="search-results" v-if="showSearchResults && searchResults.length > 0">
          <div 
            class="search-result-item" 
            v-for="item in searchResults" 
            :key="item.marketHashName"
            @click="goToDetail(item)"
          >
            <img :src="item.iconUrl" :alt="item.name" class="result-image" @error="handleImageError" />
            <span class="result-name">{{ item.name }}</span>
          </div>
        </div>
        <div class="search-results" v-else-if="showSearchResults && searchKeyword && !searchLoading">
          <div class="no-results">未找到相关饰品</div>
        </div>
      </div>
    </div>

    <!-- 数据展示区域 -->
    <div class="data-section">
      <!-- 饰品涨幅榜 -->
      <div class="data-card">
        <div class="card-header">
          <h2 class="card-title">
            <el-icon><TrendCharts /></el-icon>
            饰品涨幅榜
          </h2>
          <span class="card-badge">TOP 10</span>
        </div>
        <div class="card-body" v-loading="rankingLoading">
          <div class="ranking-list" v-if="rankingList.length">
            <div 
              class="ranking-item clickable" 
              v-for="(item, index) in rankingList" 
              :key="item.marketHashName"
              @click="goToDetailByHash(item.marketHashName)"
            >
              <div class="rank-num" :class="getRankClass(index)">{{ index + 1 }}</div>
              <div class="item-image">
                <img :src="item.iconUrl" :alt="item.name" @error="handleImageError" />
              </div>
              <div class="item-info">
                <div class="item-name">{{ item.name }}</div>
                <div class="item-price">¥{{ formatPrice(item.todayPrice) }}</div>
              </div>
              <div class="item-rate" :class="item.increaseRate1D >= 0 ? 'rate-up' : 'rate-down'">
                {{ item.increaseRate1D >= 0 ? '+' : '' }}{{ item.increaseRate1D.toFixed(2) }}%
              </div>
            </div>
          </div>
          <el-empty v-else description="暂无数据" />
        </div>
        <div class="card-footer">
          <el-button type="primary" link @click="goToRanking">
            查看完整榜单 <el-icon><ArrowRight /></el-icon>
          </el-button>
        </div>
      </div>

      <!-- 搬砖利润榜 -->
      <div class="data-card">
        <div class="card-header">
          <h2 class="card-title">
            <el-icon><DataAnalysis /></el-icon>
            搬砖利润榜
          </h2>
          <span class="card-badge">TOP 10</span>
        </div>
        <div class="card-body" v-loading="brickLoading">
          <div class="brick-list" v-if="brickList.length">
            <div 
              class="brick-item clickable" 
              v-for="(item, index) in brickList" 
              :key="item.market_hash_name"
              @click="goToDetailByHash(item.market_hash_name)"
            >
              <div class="rank-num" :class="getRankClass(index)">{{ index + 1 }}</div>
              <div class="item-image">
                <img :src="item.image_url" :alt="item.name" @error="handleImageError" />
              </div>
              <div class="item-info">
                <div class="item-name">{{ item.name }}</div>
                <div class="item-prices">
                  <span class="price-source">买入: ¥{{ formatPrice(item.source_price) }}</span>
                  <span class="price-arrow">→</span>
                  <span class="price-target">卖出: ¥{{ formatPrice(item.target_price) }}</span>
                </div>
              </div>
              <div class="item-profit">
                <div class="profit-rate">+{{ (item.profit_rate * 100).toFixed(2) }}%</div>
                <div class="profit-diff">赚 ¥{{ formatPrice(item.price_diff) }}</div>
              </div>
            </div>
          </div>
          <el-empty v-else description="暂无数据" />
        </div>
        <div class="card-footer">
          <el-button type="primary" link @click="goToHome">
            查看更多搬砖数据 <el-icon><ArrowRight /></el-icon>
          </el-button>
        </div>
      </div>
    </div>

    <!-- 特点介绍 -->
    <section class="features-section">
      <h2 class="section-title">为什么选择我们</h2>
      <div class="features-grid">
        <div class="feature-item">
          <div class="feature-icon">📊</div>
          <h3>实时数据</h3>
          <p>热门饰品分钟级实时更新</p>
        </div>
        <div class="feature-item">
          <div class="feature-icon">📈</div>
          <h3>涨跌分析</h3>
          <p>多维度数据分析，发现投资机会</p>
        </div>
        <div class="feature-item">
          <div class="feature-icon">💰</div>
          <h3>搬砖比价</h3>
          <p>跨平台价差对比，轻松赚取利润</p>
        </div>
        <div class="feature-item">
          <div class="feature-icon">🔔</div>
          <h3>专业工具</h3>
          <p>丰富的筛选条件，精准定位目标</p>
        </div>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { TrendCharts, DataAnalysis, ArrowRight, Search } from '@element-plus/icons-vue'
import { dataApi, type PriceIncreaseItem, type SearchResult } from '@/api'
import { useSkinStore } from '@/stores/skin'

const router = useRouter()
const skinStore = useSkinStore()

// 搜索相关
const searchKeyword = ref('')
const searchResults = ref<SearchResult[]>([])
const showSearchResults = ref(false)
const searchLoading = ref(false)
let searchTimer: number | null = null

// 榜单数据
const rankingList = ref<PriceIncreaseItem[]>([])
const brickList = ref<any[]>([])
const rankingLoading = ref(false)
const brickLoading = ref(false)

// 搜索输入处理（防抖）
const handleSearchInput = () => {
  if (searchTimer) {
    clearTimeout(searchTimer)
  }
  
  if (!searchKeyword.value.trim()) {
    searchResults.value = []
    showSearchResults.value = false
    return
  }
  
  searchTimer = window.setTimeout(async () => {
    await doSearch()
  }, 300)
}

// 执行搜索
const doSearch = async () => {
  if (!searchKeyword.value.trim()) return
  
  searchLoading.value = true
  showSearchResults.value = true
  
  try {
    const res = await dataApi.searchGoods({ 
      keyword: searchKeyword.value.trim(),
      limit: 50 
    })
    if (res.code === 1 && res.data) {
      searchResults.value = res.data
    }
  } catch (error) {
    console.error('搜索失败:', error)
  } finally {
    searchLoading.value = false
  }
}

// 清除搜索
const clearSearch = () => {
  searchKeyword.value = ''
  searchResults.value = []
  showSearchResults.value = false
}

// 点击搜索结果跳转到详情
const goToDetail = (item: SearchResult) => {
  showSearchResults.value = false
  const url = router.resolve({
    path: '/app/detail',
    query: { market_hash_name: item.marketHashName }
  }).href
  window.open(url, '_blank')
}

// 根据 hash name 跳转详情
const goToDetailByHash = (marketHashName: string) => {
  const url = router.resolve({
    path: '/app/detail',
    query: { market_hash_name: marketHashName }
  }).href
  window.open(url, '_blank')
}

// 获取涨幅榜数据
const fetchRankingData = async () => {
  rankingLoading.value = true
  try {
    const res = await dataApi.getPriceIncrease({ is_desc: true, limit: 10 })
    if (res.code === 1 && res.data) {
      rankingList.value = res.data
    }
  } catch (error) {
    console.error('获取涨幅榜失败:', error)
  } finally {
    rankingLoading.value = false
  }
}

// 获取搬砖利润榜数据（使用 goods/data 接口）
const fetchBrickData = async () => {
  brickLoading.value = true
  try {
    // 使用 store 中保存的平台和排序设置（与搬砖页面保持一致）
    const sortConfig = skinStore.getSortConfig()
    const res = await dataApi.getSkinItems({ 
      page_size: 10, 
      page_num: 1,
      sort: sortConfig.field,
      desc: sortConfig.desc,
      source: skinStore.sourcePlatform,
      target: skinStore.targetPlatform
    })
    if (res.code === 1 && res.data) {
      brickList.value = res.data
    }
  } catch (error) {
    console.error('获取搬砖榜失败:', error)
  } finally {
    brickLoading.value = false
  }
}

// 格式化价格
const formatPrice = (price: number) => {
  return price?.toFixed(2) || '0.00'
}

// 获取排名样式
const getRankClass = (index: number) => {
  if (index === 0) return 'rank-gold'
  if (index === 1) return 'rank-silver'
  if (index === 2) return 'rank-bronze'
  return ''
}

// 图片加载失败处理
const handleImageError = (e: Event) => {
  const target = e.target as HTMLImageElement
  // 防止无限循环：如果已经是默认图片则不再处理
  if (target.dataset.fallback) return
  target.dataset.fallback = 'true'
  // 使用本地默认图片
  target.src = '/favicon.png'
}

// 跳转到涨跌榜单
const goToRanking = () => {
  router.push('/app/ranking')
}

// 跳转到搬砖页面
const goToHome = () => {
  router.push('/app/home')
}

// 点击其他地方关闭搜索结果
const handleClickOutside = (e: MouseEvent) => {
  const searchBox = document.querySelector('.search-box')
  if (searchBox && !searchBox.contains(e.target as Node)) {
    showSearchResults.value = false
  }
}

onMounted(() => {
  fetchRankingData()
  fetchBrickData()
  document.addEventListener('click', handleClickOutside)
})

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)
  if (searchTimer) {
    clearTimeout(searchTimer)
  }
})
</script>

<style scoped>
.dashboard-page {
  padding: 20px 40px;
  max-width: 1400px;
  margin: 0 auto;
}

/* 标语区域 */
.hero-section {
  text-align: center;
  padding: 40px 0 30px;
}

.hero-title {
  font-size: 32px;
  font-weight: 600;
  color: #262626;
  margin-bottom: 12px;
}

.hero-subtitle {
  font-size: 16px;
  color: #8c8c8c;
  margin: 0;
}

/* 搜索区域 */
.search-section {
  margin-bottom: 30px;
}

.search-box {
  max-width: 600px;
  margin: 0 auto;
  position: relative;
}

.search-box :deep(.el-input__wrapper) {
  border-radius: 24px;
  padding: 4px 20px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
}

.search-box :deep(.el-input__wrapper:hover),
.search-box :deep(.el-input__wrapper.is-focus) {
  box-shadow: 0 4px 16px rgba(24, 144, 255, 0.15);
}

.search-results {
  position: absolute;
  top: 100%;
  left: 0;
  right: 0;
  margin-top: 8px;
  background: #fff;
  border-radius: 12px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.12);
  z-index: 100;
  max-height: 400px;
  overflow-y: auto;
}

.search-result-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  cursor: pointer;
  transition: background 0.2s;
}

.search-result-item:hover {
  background: #f5f7fa;
}

.search-result-item:first-child {
  border-radius: 12px 12px 0 0;
}

.search-result-item:last-child {
  border-radius: 0 0 12px 12px;
}

.result-image {
  width: 40px;
  height: 30px;
  object-fit: contain;
  background: #f5f7fa;
  border-radius: 4px;
}

.result-name {
  flex: 1;
  font-size: 14px;
  color: #262626;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.no-results {
  padding: 20px;
  text-align: center;
  color: #8c8c8c;
  font-size: 14px;
}

/* 数据展示区域 */
.data-section {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 24px;
}

.data-card {
  background: #fff;
  border-radius: 12px;
  border: 1px solid #e8e8e8;
  overflow: hidden;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px 24px;
  background: #fafafa;
  border-bottom: 1px solid #f0f0f0;
}

.card-title {
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 18px;
  font-weight: 600;
  color: #262626;
  margin: 0;
}

.card-title .el-icon {
  color: #1890ff;
}

.card-badge {
  padding: 4px 12px;
  background: #e6f4ff;
  color: #1890ff;
  border-radius: 20px;
  font-size: 12px;
  font-weight: 600;
}

.card-body {
  padding: 16px 24px;
  min-height: 400px;
}

.card-footer {
  padding: 16px 24px;
  text-align: center;
  border-top: 1px solid #f0f0f0;
  background: #fafafa;
}

/* 排行榜列表 */
.ranking-list,
.brick-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.ranking-item,
.brick-item {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 12px 16px;
  background: #fafafa;
  border-radius: 8px;
  transition: all 0.2s;
}

.ranking-item.clickable,
.brick-item.clickable {
  cursor: pointer;
}

.ranking-item:hover,
.brick-item:hover {
  background: #f0f5ff;
}

.rank-num {
  width: 28px;
  height: 28px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 600;
  background: #f0f0f0;
  color: #8c8c8c;
}

.rank-gold {
  background: linear-gradient(135deg, #ffd700, #ffaa00);
  color: #fff;
}

.rank-silver {
  background: linear-gradient(135deg, #c0c0c0, #a0a0a0);
  color: #fff;
}

.rank-bronze {
  background: linear-gradient(135deg, #cd7f32, #b87333);
  color: #fff;
}

.item-image {
  width: 60px;
  height: 45px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #f5f7fa;
  border-radius: 6px;
}

.item-image img {
  max-width: 100%;
  max-height: 100%;
  object-fit: contain;
}

.item-info {
  flex: 1;
  min-width: 0;
}

.item-name {
  font-size: 14px;
  font-weight: 500;
  color: #262626;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.item-price {
  font-size: 13px;
  color: #8c8c8c;
  margin-top: 4px;
}

.item-prices {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 12px;
  color: #8c8c8c;
  margin-top: 4px;
}

.price-arrow {
  color: #bfbfbf;
}

.item-rate {
  font-size: 15px;
  font-weight: 600;
  padding: 4px 10px;
  border-radius: 6px;
}

.rate-up {
  color: #52c41a;
  background: #f6ffed;
}

.rate-down {
  color: #ff4d4f;
  background: #fff2f0;
}

.item-profit {
  text-align: right;
}

.profit-rate {
  font-size: 15px;
  font-weight: 600;
  color: #52c41a;
}

.profit-diff {
  font-size: 12px;
  color: #8c8c8c;
  margin-top: 2px;
}

/* 特点介绍 */
.features-section {
  padding: 60px 0;
}

.section-title {
  text-align: center;
  font-size: 24px;
  font-weight: 600;
  color: #262626;
  margin-bottom: 40px;
}

.features-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 24px;
}

.feature-item {
  text-align: center;
  padding: 30px 20px;
  background: #fff;
  border-radius: 12px;
  border: 1px solid #e8e8e8;
  transition: all 0.3s;
}

.feature-item:hover {
  transform: translateY(-4px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.08);
  border-color: #1890ff;
}

.feature-icon {
  font-size: 40px;
  margin-bottom: 16px;
}

.feature-item h3 {
  font-size: 18px;
  font-weight: 600;
  color: #262626;
  margin-bottom: 8px;
}

.feature-item p {
  font-size: 14px;
  color: #8c8c8c;
  margin: 0;
}

/* 响应式 */
@media (max-width: 1200px) {
  .data-section {
    grid-template-columns: 1fr;
  }
  
  .features-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 768px) {
  .dashboard-page {
    padding: 20px;
  }
  
  .features-grid {
    grid-template-columns: 1fr;
  }
}
</style>

