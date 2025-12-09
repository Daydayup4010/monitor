<template>
  <div class="landing-page">
    <!-- 头部区域 -->
    <header class="landing-header">
      <div class="header-content">
        <div class="logo">
          <span class="logo-text">CS Goods</span>
          <span class="logo-subtitle">CS2饰品搬砖网站</span>
        </div>
        <div class="header-actions">
          <el-button type="primary" @click="showLoginDialog = true">登录</el-button>
          <el-button @click="goToRegister">注册</el-button>
        </div>
      </div>
    </header>

    <!-- 登录弹窗 -->
    <LoginDialog v-model="showLoginDialog" />

    <!-- 主要内容区域 -->
    <main class="landing-main">
      <!-- 标语区域 -->
      <section class="hero-section">
        <h1 class="hero-title">专注于CS2饰品搬砖</h1>
        <p class="hero-subtitle">提供专业的饰品数据变化分析、各平台搬砖比价数据</p>
      </section>

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
          <div class="card-body" v-loading="loading">
            <div class="ranking-list" v-if="homeData?.rankingList?.length">
              <div 
                class="ranking-item" 
                v-for="(item, index) in homeData.rankingList" 
                :key="item.marketHashName"
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
                  {{ item.increaseRate1D >= 0 ? '+' : '' }}{{ (item.increaseRate1D * 100).toFixed(2) }}%
                </div>
              </div>
            </div>
            <el-empty v-else description="暂无数据" />
          </div>
          <div class="card-footer">
            <el-button type="primary" link @click="goToLoginWithTip">
              登录查看完整榜单 <el-icon><ArrowRight /></el-icon>
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
          <div class="card-body" v-loading="loading">
            <div class="brick-list" v-if="homeData?.brickMoving?.length">
              <div 
                class="brick-item" 
                v-for="(item, index) in homeData.brickMoving" 
                :key="item.market_hash_name"
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
            <el-button type="primary" link @click="goToLoginWithTip">
              登录查看更多搬砖商机 <el-icon><ArrowRight /></el-icon>
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
    </main>

    <!-- 底部区域 -->
    <Footer />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { TrendCharts, DataAnalysis, ArrowRight } from '@element-plus/icons-vue'
import { publicApi, type PublicHomeData } from '@/api'
import { showMessage } from '@/utils/message'
import LoginDialog from '@/components/LoginDialog.vue'
import Footer from '@/components/Footer.vue'

const router = useRouter()
const loading = ref(false)
const homeData = ref<PublicHomeData | null>(null)
const showLoginDialog = ref(false)

// 获取首页数据
const fetchHomeData = async () => {
  loading.value = true
  try {
    const res = await publicApi.getHomeData()
    if (res.code === 1 && res.data) {
      homeData.value = res.data
    }
  } catch (error) {
    console.error('获取首页数据失败:', error)
  } finally {
    loading.value = false
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
  target.src = 'https://steamcommunity-a.akamaihd.net/economy/image/-9a81dlWLwJ2UUGcVs_nsVtzdOEdtWwKGZZLQHTxDZ7I56KU0Zwwo4NUX4oFJZEHLbXQ9QVcJY8gulRPQV6CF7b9mMnYZh9SHY27gZKBl_JbMKyJI24H65S1xtXZwKb2YOqHxj4F68Nz2L2Y9oj2jQDm_RY4am-mctWXdFc5NQuDqAHqx-fmg5_v7oOJlyU1fQmQdw/360fx360f'
}

// 跳转到注册页
const goToRegister = () => {
  router.push('/register')
}

// 提示后打开登录弹窗
const goToLoginWithTip = () => {
  showMessage.info('请先登录后查看更多内容')
  showLoginDialog.value = true
}

onMounted(() => {
  fetchHomeData()
})
</script>

<style scoped>
.landing-page {
  min-height: 100vh;
  background: #f5f7fa;
}

/* 头部 */
.landing-header {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  z-index: 100;
  background: #fff;
  border-bottom: 1px solid #e8e8e8;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
}

.header-content {
  max-width: 1400px;
  margin: 0 auto;
  padding: 16px 40px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.logo {
  display: flex;
  flex-direction: column;
}

.logo-text {
  font-size: 24px;
  font-weight: 700;
  color: #1890ff;
}

.logo-subtitle {
  font-size: 12px;
  color: #8c8c8c;
  margin-top: 2px;
}

.header-actions {
  display: flex;
  gap: 12px;
}

/* 主要内容 */
.landing-main {
  padding-top: 80px;
  max-width: 1400px;
  margin: 0 auto;
  padding-left: 40px;
  padding-right: 40px;
}

/* 标语区域 */
.hero-section {
  text-align: center;
  padding: 60px 0 40px;
}

.hero-title {
  font-size: 36px;
  font-weight: 600;
  color: #262626;
  margin-bottom: 16px;
}

.hero-subtitle {
  font-size: 16px;
  color: #8c8c8c;
}

/* 数据展示区域 */
.data-section {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 24px;
  margin-bottom: 60px;
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
  .header-content {
    padding: 12px 20px;
  }
  
  .landing-main {
    padding-left: 20px;
    padding-right: 20px;
  }
  
  .hero-title {
    font-size: 24px;
  }
  
  .hero-subtitle {
    font-size: 14px;
  }
  
  .features-grid {
    grid-template-columns: 1fr;
  }
}
</style>
