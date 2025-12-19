<template>
  <div class="big-item-page">
    <div class="card">
      <div class="card-title">
        <el-icon :size="20" color="#1890ff"><ShoppingCart /></el-icon>
        大件求购
      </div>

      <!-- 筛选栏 -->
      <div class="filter-bar" style="display: block !important;">
        <!-- 第一行 - 平台和类别选择 -->
        <div style="display: flex; gap: 12px; flex-wrap: wrap; align-items: flex-end; margin-bottom: 16px;">
          <div class="filter-item">
            <label class="filter-label">平台</label>
            <el-select
              v-model="platform"
              style="width: 150px;"
              @change="handlePlatformChange"
            >
              <template #label="{ value }">
                <div style="display: flex; align-items: center; gap: 8px;">
                  <img :src="getPlatformIcon(value)" style="width: 18px; height: 18px; object-fit: contain;" />
                  <span>{{ getPlatformSelectLabel(value) }}</span>
                </div>
              </template>
              <el-option value="uu" label="悠悠">
                <div style="display: flex; align-items: center; gap: 8px;">
                  <img :src="uuIcon" style="width: 18px; height: 18px; object-fit: contain;" />
                  <span>悠悠</span>
                </div>
              </el-option>
              <el-option value="buff" label="BUFF">
                <div style="display: flex; align-items: center; gap: 8px;">
                  <img :src="buffIcon" style="width: 18px; height: 18px; object-fit: contain;" />
                  <span>BUFF</span>
                </div>
              </el-option>
              <el-option value="c5" label="C5GAME">
                <div style="display: flex; align-items: center; gap: 8px;">
                  <img :src="c5Icon" style="width: 18px; height: 18px; object-fit: contain;" />
                  <span>C5GAME</span>
                </div>
              </el-option>
            </el-select>
          </div>

          <div class="filter-item">
            <label class="filter-label">类别</label>
            <el-select
              v-model="category"
              style="width: 150px;"
              @change="handleCategoryChange"
            >
              <el-option label="手套和刀具" value="all" />
              <el-option label="手套" value="手套" />
              <el-option label="刀具" value="匕首" />
            </el-select>
          </div>
        </div>

        <!-- 第二行 - 搜索和排序 -->
        <div style="display: flex; gap: 12px; align-items: flex-end; margin-bottom: 16px;">
          <div style="flex: 1; min-width: 200px; max-width: 500px;">
            <label style="display: block; margin-bottom: 8px; font-size: 14px; color: #595959; font-weight: 500;">搜索</label>
            <el-input
              v-model="searchKeyword"
              placeholder="搜索饰品名称..."
              clearable
              @input="handleSearch"
            >
              <template #prefix>
                <el-icon><Search /></el-icon>
              </template>
            </el-input>
          </div>

          <div style="min-width: 160px;">
            <label style="display: block; margin-bottom: 8px; font-size: 14px; color: #595959; font-weight: 500;">排序</label>
            <el-select
              v-model="sortOption"
              @change="handleSortChange"
              style="width: 160px;"
            >
              <el-option label="默认" value="default" />
              <el-option label="价格差 ↑" value="price_diff_asc" />
              <el-option label="价格差 ↓" value="price_diff_desc" />
              <el-option label="利润率 ↑" value="profit_rate_asc" />
              <el-option label="利润率 ↓" value="profit_rate_desc" />
              <el-option label="在售价 ↑" value="sell_price_asc" />
              <el-option label="在售价 ↓" value="sell_price_desc" />
              <el-option label="求购价 ↑" value="bidding_price_asc" />
              <el-option label="求购价 ↓" value="bidding_price_desc" />
            </el-select>
          </div>

          <div>
            <button
              class="btn btn-primary"
              :disabled="loading"
              @click="refreshData"
              style="height: 40px; padding: 0 24px; font-size: 14px;"
            >
              {{ loading ? '刷新中...' : '确定并搜索' }}
            </button>
          </div>

          <div style="margin-left: auto;">
            <button
              class="btn btn-refresh"
              :disabled="loading"
              @click="reloadData"
              style="height: 40px; padding: 0 24px; font-size: 14px;"
            >
              <el-icon style="margin-right: 6px;"><Refresh /></el-icon>
              {{ loading ? '刷新中...' : '刷新数据' }}
            </button>
          </div>
        </div>

        <!-- 当前选项描述 -->
        <div style="padding: 12px; background: #f5f7fa; border-radius: 8px; font-size: 13px; color: #595959;">
          <span style="font-weight: 600;">当前选项描述：</span>
          <span>{{ getFilterDescription() }}</span>
        </div>
      </div>

      <!-- 数据表格 -->
      <div class="table-wrapper">
        <el-table
          :data="items"
          v-loading="loading"
          style="width: 100%"
        >
          <el-table-column type="index" label="#" width="60" />
          
          <el-table-column label="饰品名称" min-width="250" class-name="skin-column">
            <template #default="{ row }">
              <div class="skin-info">
                <img :src="row.image_url" @error="handleImageError" alt="饰品" class="skin-img" />
                <div>
                  <div class="skin-name clickable" @click="goToDetail(row.market_hash_name)">{{ row.name }}</div>
                  <div style="font-size: 12px; color: #8c8c8c;">{{ row.type_name }}</div>
                  <div class="update-time-inline">{{ formatUpdateTime(row.update_time) }}</div>
                </div>
              </div>
            </template>
          </el-table-column>
          
          <el-table-column label="在售价" width="140" class-name="price-column">
            <template #default="{ row }">
              <div class="price-value">¥{{ formatPrice(row.sell_price) }}</div>
              <div class="count-info">在售 {{ row.sell_count }} 件</div>
            </template>
          </el-table-column>
          
          <el-table-column label="求购价" width="140" class-name="price-column">
            <template #default="{ row }">
              <div class="price-value bidding">¥{{ formatPrice(row.bidding_price) }}</div>
              <div class="count-info">求购 {{ row.bidding_count }} 个</div>
            </template>
          </el-table-column>
          
          <el-table-column label="价差" width="120" class-name="diff-column">
            <template #default="{ row }">
              <div class="diff-value" :class="{ positive: row.price_diff > 0 }">
                {{ row.price_diff > 0 ? '+' : '' }}¥{{ formatPrice(row.price_diff) }}
              </div>
            </template>
          </el-table-column>
          
          <el-table-column label="利润率" width="120" class-name="rate-column">
            <template #default="{ row }">
              <div class="rate-value" :class="{ positive: row.profit_rate > 0 }">
                {{ formatPercent(row.profit_rate) }}
              </div>
            </template>
          </el-table-column>

          <el-table-column label="各平台数据" width="120" align="center">
            <template #default="{ row }">
              <el-popover
                v-if="row.platform_list && row.platform_list.length > 0"
                placement="right"
                :width="500"
                trigger="hover"
                :show-after="300"
                :hide-after="0"
                :persistent="false"
              >
                <template #reference>
                  <div class="platform-data-cell">
                    <span class="platform-data-link">平台数据</span>
                  </div>
                </template>
                <div style="padding: 8px; max-height: 500px; overflow-y: auto;">
                  <!-- 当前平台数据（抬头） -->
                  <div v-if="getCurrentPlatformData(row)" style="padding: 12px; background: #f5f7fa; border-radius: 8px; margin-bottom: 12px;">
                    <div style="display: flex; align-items: center; gap: 20px;">
                      <div style="display: flex; align-items: center; gap: 10px;">
                        <img 
                          :src="getPlatformIcon(platform)" 
                          style="width: 20px; height: 20px; object-fit: contain; cursor: pointer;" 
                          :alt="platform"
                          @click="openPlatformLink(row, platform)"
                        />
                        <span style="font-weight: 600; font-size: 14px; color: #262626;">{{ getPlatformSelectLabel(platform) }}</span>
                      </div>
                      <div style="display: flex; gap: 40px; flex: 1; justify-content: center;">
                        <span style="font-size: 13px; white-space: nowrap;">在售价(¥): <strong style="color: #52c41a;">{{ formatPrice(getCurrentPlatformData(row)?.sellPrice) }}</strong>({{ getCurrentPlatformData(row)?.sellCount }})</span>
                        <span style="font-size: 13px; white-space: nowrap;">求购价(¥): <strong style="color: #ff4d4f;">{{ formatPrice(getCurrentPlatformData(row)?.biddingPrice) }}</strong>({{ getCurrentPlatformData(row)?.biddingCount }})</span>
                      </div>
                    </div>
                  </div>
                  
                  <!-- 其他平台数据表格 -->
                  <table style="width: 100%; border-collapse: collapse; font-size: 13px;">
                    <thead>
                      <tr style="background: #fafafa; border-bottom: 1px solid #e8e8e8;">
                        <th style="padding: 10px 12px; text-align: left; font-weight: 600; color: #595959;">平台</th>
                        <th style="padding: 10px 12px; text-align: left; font-weight: 600; color: #595959;">在售价(¥)</th>
                        <th style="padding: 10px 12px; text-align: left; font-weight: 600; color: #595959;">求购价(¥)</th>
                      </tr>
                    </thead>
                    <tbody>
                      <tr 
                        v-for="p in getOtherPlatforms(row)" 
                        :key="p.platform"
                        style="border-bottom: 1px solid #f0f0f0;"
                      >
                        <td style="padding: 10px 12px;">
                          <div style="display: flex; align-items: center; gap: 8px;">
                            <img 
                              :src="getPlatformIconByName(p.platformName)" 
                              style="width: 18px; height: 18px; object-fit: contain; cursor: pointer;" 
                              @click="openPlatformLinkByData(p)"
                              :alt="p.platformName"
                            />
                            <span>{{ p.platformName }}</span>
                          </div>
                        </td>
                        <td style="padding: 10px 12px; color: #52c41a; font-weight: 500;">
                          {{ formatPrice(p.sellPrice) }}<span style="color: #8c8c8c;">({{ p.sellCount }})</span>
                        </td>
                        <td style="padding: 10px 12px; color: #ff4d4f; font-weight: 500;">
                          {{ formatPrice(p.biddingPrice) }}<span style="color: #8c8c8c;">({{ p.biddingCount }})</span>
                        </td>
                      </tr>
                    </tbody>
                  </table>
                </div>
              </el-popover>
              <span v-else style="color: #8c8c8c;">-</span>
            </template>
          </el-table-column>
        </el-table>
      </div>

      <!-- 分页 -->
      <div class="pagination-wrapper">
        <div class="pagination-info">
          共 {{ total }} 条记录，每页 
          <el-select v-model="pageSize" style="width: 80px; margin: 0 8px;" @change="handleSizeChange">
            <el-option :value="25" label="25" />
            <el-option :value="50" label="50" />
            <el-option :value="100" label="100" />
          </el-select>
          条
        </div>
        
        <div class="pagination-pages">
          <button class="page-btn" :disabled="pageNum === 1" @click="handleCurrentChange(pageNum - 1)">‹</button>
          <button
            v-for="page in visiblePages"
            :key="page"
            class="page-btn"
            :class="{ active: page === pageNum, ellipsis: page === '...' }"
            :disabled="page === '...'"
            @click="page !== '...' && handleCurrentChange(page as number)"
          >
            {{ page }}
          </button>
          <button class="page-btn" :disabled="pageNum >= totalPages" @click="handleCurrentChange(pageNum + 1)">›</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import { Search, Refresh, ShoppingCart } from '@element-plus/icons-vue'
import { bigItemApi, type BigItemBiddingItem, type PlatformItem } from '@/api'
import { formatPrice, formatPercent, debounce } from '@/utils'
import buffIcon from '@/assets/icons/buff.png'
import uuIcon from '@/assets/icons/uu.png'
import c5Icon from '@/assets/icons/c5.png'
import steamIcon from '@/assets/icons/steam.png'

const router = useRouter()

// 数据
const items = ref<BigItemBiddingItem[]>([])
const total = ref(0)
const loading = ref(false)

// 分页
const pageNum = ref(1)
const pageSize = ref(50)

// 筛选
const searchKeyword = ref('')
const platform = ref(localStorage.getItem('big_item_platform') || 'uu')
const category = ref(localStorage.getItem('big_item_category') || 'all')
const sortOption = ref(localStorage.getItem('big_item_sort') || 'default')

// 监听变化，保存到 localStorage
watch(platform, (val) => localStorage.setItem('big_item_platform', val))
watch(category, (val) => localStorage.setItem('big_item_category', val))
watch(sortOption, (val) => localStorage.setItem('big_item_sort', val))

// 排序配置映射
const sortMap: Record<string, { field: string; desc: boolean }> = {
  'default': { field: 'profit_rate', desc: true },
  'profit_rate_desc': { field: 'profit_rate', desc: true },
  'profit_rate_asc': { field: 'profit_rate', desc: false },
  'price_diff_desc': { field: 'price_diff', desc: true },
  'price_diff_asc': { field: 'price_diff', desc: false },
  'sell_price_desc': { field: 'sell_price', desc: true },
  'sell_price_asc': { field: 'sell_price', desc: false },
  'bidding_price_desc': { field: 'bidding_price', desc: true },
  'bidding_price_asc': { field: 'bidding_price', desc: false },
}

const totalPages = computed(() => Math.ceil(total.value / pageSize.value))

const visiblePages = computed(() => {
  const current = pageNum.value
  const totalP = totalPages.value
  const pages: (number | string)[] = []
  
  if (totalP <= 7) {
    for (let i = 1; i <= totalP; i++) pages.push(i)
  } else {
    pages.push(1)
    if (current <= 4) {
      for (let i = 2; i <= 5; i++) pages.push(i)
      pages.push('...')
      pages.push(totalP)
    } else if (current >= totalP - 3) {
      pages.push('...')
      for (let i = totalP - 4; i <= totalP; i++) pages.push(i)
    } else {
      pages.push('...')
      for (let i = current - 1; i <= current + 1; i++) pages.push(i)
      pages.push('...')
      pages.push(totalP)
    }
  }
  return pages
})

// 获取数据
const fetchData = async () => {
  loading.value = true
  try {
    const sortConfig = sortMap[sortOption.value] || sortMap['default']
    const response = await bigItemApi.getBigItemBidding({
      page_num: pageNum.value,
      page_size: pageSize.value,
      sort: sortConfig.field,
      desc: sortConfig.desc,
      search: searchKeyword.value,
      platform: platform.value,
      category: category.value,
    })
    items.value = response.data || []
    total.value = (response as any).total || 0
  } catch (error) {
    console.error('获取大件求购数据失败:', error)
  } finally {
    loading.value = false
  }
}

// 事件处理
const handlePlatformChange = () => {
  pageNum.value = 1
  fetchData()
}

const handleCategoryChange = () => {
  pageNum.value = 1
  fetchData()
}

const handleSortChange = () => {
  pageNum.value = 1
  fetchData()
}

const handleSearch = debounce(() => {
  pageNum.value = 1
  fetchData()
}, 300)

const refreshData = () => {
  pageNum.value = 1
  fetchData()
}

const reloadData = () => {
  fetchData()
}

const handleSizeChange = () => {
  pageNum.value = 1
  fetchData()
}

const handleCurrentChange = (page: number) => {
  pageNum.value = page
  fetchData()
}

const goToDetail = (marketHashName: string) => {
  router.push({
    path: '/app/detail',
    query: { market_hash_name: marketHashName }
  })
}

// 工具函数
const getPlatformIcon = (p: string) => {
  const icons: Record<string, string> = { buff: buffIcon, uu: uuIcon, c5: c5Icon }
  return icons[p] || uuIcon
}

const getPlatformSelectLabel = (p: string) => {
  const labels: Record<string, string> = { buff: 'BUFF', uu: '悠悠', c5: 'C5GAME' }
  return labels[p] || p
}

const getFilterDescription = () => {
  const platformNames: Record<string, string> = { buff: 'BUFF', uu: '悠悠', c5: 'C5GAME' }
  const categoryNames: Record<string, string> = { all: '手套和刀具', '手套': '手套', '匕首': '刀具' }
  const platformName = platformNames[platform.value] || platform.value
  const categoryName = categoryNames[category.value] || category.value
  return `查看${platformName}平台${categoryName}的求购价差（价差 = 在售价 - 求购价）`
}

const formatUpdateTime = (timestamp: number) => {
  if (!timestamp) return '-'
  const date = new Date(timestamp * 1000)
  const now = new Date()
  const diff = Math.floor((now.getTime() - date.getTime()) / 1000)
  
  if (diff < 60) return '刚刚'
  if (diff < 3600) return `${Math.floor(diff / 60)} 分钟前`
  if (diff < 86400) return `${Math.floor(diff / 3600)} 小时前`
  
  return date.toLocaleDateString('zh-CN', { month: '2-digit', day: '2-digit', hour: '2-digit', minute: '2-digit' })
}

const handleImageError = (e: Event) => {
  const img = e.target as HTMLImageElement
  if (!img.dataset.fallback) {
    img.dataset.fallback = 'true'
    img.src = '/favicon.png'
  }
}

// 平台数据相关函数
const getCurrentPlatformData = (row: BigItemBiddingItem) => {
  const platformNameMap: Record<string, string> = {
    'buff': 'BUFF',
    'uu': '悠悠',
    'c5': 'C5GAME'
  }
  return row.platform_list?.find((p: PlatformItem) => p.platformName === platformNameMap[platform.value])
}

const getOtherPlatforms = (row: BigItemBiddingItem) => {
  const platformNameMap: Record<string, string> = {
    'buff': 'BUFF',
    'uu': '悠悠',
    'c5': 'C5GAME'
  }
  const currentName = platformNameMap[platform.value]
  return row.platform_list?.filter((p: PlatformItem) => p.platformName !== currentName) || []
}

const openPlatformLink = (row: BigItemBiddingItem, p: string) => {
  const platformNameMap: Record<string, string> = {
    'buff': 'BUFF',
    'uu': '悠悠',
    'c5': 'C5GAME'
  }
  const platformData = row.platform_list?.find((item: PlatformItem) => item.platformName === platformNameMap[p])
  if (platformData?.link) {
    window.open(platformData.link, '_blank')
  }
}

const openPlatformLinkByData = (p: PlatformItem) => {
  if (p?.link) {
    window.open(p.link, '_blank')
  }
}

const getPlatformIconByName = (platformName: string) => {
  const iconMap: Record<string, string> = {
    'BUFF': buffIcon,
    '悠悠': uuIcon,
    'C5GAME': c5Icon,
    'STEAM': steamIcon,
    'Steam': steamIcon
  }
  return iconMap[platformName] || ''
}

onMounted(() => {
  fetchData()
})
</script>

<style scoped>
.big-item-page {
  padding: 0;
  width: 1800px;
  max-width: 100%;
  margin: 0 auto;
}

/* 使用 unified.css 中的全局 .card 和 .card-title 样式 */

.filter-bar {
  margin-bottom: 24px;
}

.filter-item {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.filter-label {
  font-size: 14px;
  color: #595959;
  font-weight: 500;
}

.table-wrapper {
  margin-bottom: 24px;
  border-radius: 8px;
  overflow: hidden;
  border: 1px solid #e8e8e8;
}

/* 使用 unified.css 中的全局 .skin-info, .skin-img, .skin-name 样式 */

.skin-name.clickable {
  cursor: pointer;
  transition: color 0.2s;
}

.skin-name.clickable:hover {
  color: #1890ff;
}

.update-time {
  font-size: 13px;
  color: #8c8c8c;
}

.update-time-inline {
  font-size: 11px;
  color: #bfbfbf;
  margin-top: 2px;
}

/* 价格相关样式 - 与Home.vue保持一致 */
.price-value {
  color: #52c41a;
  font-weight: 600;
  font-size: 14px;
}

.price-value.bidding {
  color: #ff4d4f;
}

.count-info {
  font-size: 12px;
  color: #8c8c8c;
  margin-top: 2px;
}

.diff-value {
  color: #faad14;
  font-weight: 600;
  font-size: 15px;
}

.diff-value.positive {
  color: #52c41a;
}

.rate-value {
  font-weight: 600;
  font-size: 15px;
  color: #595959;
}

.rate-value.positive {
  color: #52c41a;
}

/* 平台数据单元格 */
.platform-data-cell {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
}

.platform-data-link {
  color: #1890ff;
  cursor: pointer;
  font-size: 13px;
}

.platform-data-link:hover {
  color: #40a9ff;
}

.pagination-wrapper {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 0;
}

.pagination-info {
  font-size: 14px;
  color: #595959;
  display: flex;
  align-items: center;
}

.pagination-pages {
  display: flex;
  gap: 4px;
}

.page-btn {
  min-width: 32px;
  height: 32px;
  padding: 0 8px;
  border: 1px solid #d9d9d9;
  background: #fff;
  border-radius: 6px;
  cursor: pointer;
  font-size: 14px;
  color: #595959;
  transition: all 0.2s;
}

.page-btn:hover:not(:disabled) {
  border-color: #3b82f6;
  color: #3b82f6;
}

.page-btn.active {
  background: #3b82f6;
  border-color: #3b82f6;
  color: #fff;
}

.page-btn:disabled {
  cursor: not-allowed;
  opacity: 0.5;
}

.page-btn.ellipsis {
  border: none;
  cursor: default;
}

.btn {
  border: none;
  border-radius: 8px;
  cursor: pointer;
  font-weight: 500;
  transition: all 0.2s;
  display: flex;
  align-items: center;
  justify-content: center;
}

.btn-primary {
  background: linear-gradient(135deg, #3b82f6 0%, #2563eb 100%);
  color: #fff;
}

.btn-primary:hover:not(:disabled) {
  background: linear-gradient(135deg, #2563eb 0%, #1d4ed8 100%);
}

.btn-primary:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.btn-refresh {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #52c41a, #73d13d);
  color: white;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  font-weight: 500;
  transition: all 0.3s ease;
}

.btn-refresh:hover:not(:disabled) {
  background: linear-gradient(135deg, #389e0d, #52c41a);
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(82, 196, 26, 0.4);
}

.btn-refresh:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

:deep(.el-table) {
  --el-table-header-bg-color: #f8fafc;
  --el-table-header-text-color: #64748b;
  font-size: 14px;
}

:deep(.el-table th) {
  font-weight: 600;
  padding: 16px 12px;
}

:deep(.el-table td) {
  padding: 16px 12px;
}

:deep(.el-table__row:hover > td) {
  background: #f8fafc !important;
}
</style>

