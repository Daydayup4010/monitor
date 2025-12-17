<template>
  <div class="home-page">
    <div class="card">
      <div class="card-title">
        <el-icon :size="20" color="#1890ff"><DataAnalysis /></el-icon>
        挂刀/搬砖
      </div>

      <!-- 筛选栏 -->
      <div class="filter-bar" style="display: block !important;">
        <!-- 第一行 - 平台选择 -->
        <div style="display: flex; gap: 12px; flex-wrap: wrap; align-items: flex-end; margin-bottom: 16px;">
          <div class="filter-item">
            <label class="filter-label">买入平台</label>
            <el-select
              v-model="sourcePlatform"
              style="width: 150px;"
              @change="handleSourceChange"
            >
              <template #label="{ value }">
                <div style="display: flex; align-items: center; gap: 8px;">
                  <img :src="getPlatformIcon(value)" style="width: 18px; height: 18px; object-fit: contain;" />
                  <span>{{ getPlatformSelectLabel(value) }}</span>
                </div>
              </template>
              <el-option value="buff" label="BUFF">
                <div style="display: flex; align-items: center; gap: 8px;">
                  <img :src="buffIcon" style="width: 18px; height: 18px; object-fit: contain;" />
                  <span>BUFF</span>
                </div>
              </el-option>
              <el-option value="uu" label="悠悠">
                <div style="display: flex; align-items: center; gap: 8px;">
                  <img :src="uuIcon" style="width: 18px; height: 18px; object-fit: contain;" />
                  <span>悠悠</span>
                </div>
              </el-option>
              <el-option value="c5" label="C5GAME">
                <div style="display: flex; align-items: center; gap: 8px;">
                  <img :src="c5Icon" style="width: 18px; height: 18px; object-fit: contain;" />
                  <span>C5GAME</span>
                </div>
              </el-option>
              <el-option value="steam" label="Steam">
                <div style="display: flex; align-items: center; gap: 8px;">
                  <img :src="steamIcon" style="width: 18px; height: 18px; object-fit: contain;" />
                  <span>Steam</span>
                </div>
              </el-option>
            </el-select>
          </div>

          <div class="filter-item">
            <label class="filter-label">卖出平台</label>
            <el-select
              v-model="targetPlatform"
              style="width: 150px;"
              @change="handleTargetChange"
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
              <el-option value="steam" label="Steam">
                <div style="display: flex; align-items: center; gap: 8px;">
                  <img :src="steamIcon" style="width: 18px; height: 18px; object-fit: contain;" />
                  <span>Steam</span>
                </div>
              </el-option>
            </el-select>
          </div>
        </div>

        <!-- 第二行 - 筛选参数 -->
        <div style="display: flex; gap: 20px; flex-wrap: wrap; align-items: center; margin-bottom: 16px;">
          <div style="display: flex; align-items: center; gap: 8px;">
            <div style="display: flex; align-items: center; gap: 4px;">
              <label style="font-size: 14px; color: #595959; white-space: nowrap;">价格范围</label>
              <el-tooltip content="筛选指定价格区间的饰品" placement="top">
                <el-icon style="color: #8c8c8c; cursor: help; font-size: 14px;"><QuestionFilled /></el-icon>
              </el-tooltip>
            </div>
            <span style="color: #8c8c8c;">:</span>
            <el-input
              v-model.number="filterParams.min_sell_price"
              type="number"
              placeholder="0"
              style="width: 100px;"
            >
              <template #suffix>¥</template>
            </el-input>
            <span style="color: #8c8c8c;">~</span>
            <el-input
              v-model.number="filterParams.max_sell_price"
              type="number"
              placeholder="10000"
              style="width: 100px;"
            >
              <template #suffix>¥</template>
            </el-input>
          </div>

          <div style="display: flex; align-items: center; gap: 8px;">
            <div style="display: flex; align-items: center; gap: 4px;">
              <label style="font-size: 14px; color: #595959; white-space: nowrap;">平台在售量</label>
              <el-tooltip content="筛选卖出平台当前在售数量大于指定值的饰品" placement="top">
                <el-icon style="color: #8c8c8c; cursor: help; font-size: 14px;"><QuestionFilled /></el-icon>
              </el-tooltip>
            </div>
            <span style="color: #8c8c8c;">:</span>
            <el-input
              v-model.number="filterParams.min_sell_num"
              type="number"
              placeholder="0"
              style="width: 100px;"
            >
              <template #suffix>件</template>
            </el-input>
          </div>

          <div style="display: flex; align-items: center; gap: 8px;">
            <div style="display: flex; align-items: center; gap: 4px;">
              <label style="font-size: 14px; color: #595959; white-space: nowrap;">最小价格差</label>
              <el-tooltip content="筛选买卖平台价格差大于指定值的饰品" placement="top">
                <el-icon style="color: #8c8c8c; cursor: help; font-size: 14px;"><QuestionFilled /></el-icon>
              </el-tooltip>
            </div>
            <span style="color: #8c8c8c;">:</span>
            <el-input
              v-model.number="filterParams.min_diff"
              type="number"
              placeholder="0"
              style="width: 100px;"
            >
              <template #suffix>¥</template>
            </el-input>
          </div>
        </div>

        <!-- 第三行 - 搜索、类别和排序 -->
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

          <div style="min-width: 120px;">
            <label style="display: block; margin-bottom: 8px; font-size: 14px; color: #595959; font-weight: 500;">类别</label>
            <el-select
              v-model="category"
              @change="handleCategoryChange"
              style="width: 120px;"
              clearable
              placeholder="全部"
            >
              <el-option label="全部" value="" />
              <el-option label="匕首" value="匕首" />
              <el-option label="手套" value="手套" />
              <el-option label="步枪" value="步枪" />
              <el-option label="手枪" value="手枪" />
              <el-option label="微型冲锋枪" value="微型冲锋枪" />
              <el-option label="霰弹枪" value="霰弹枪" />
              <el-option label="机枪" value="机枪" />
              <el-option label="印花" value="印花" />
              <el-option label="涂鸦" value="涂鸦" />
              <el-option label="探员" value="探员" />
              <el-option label="挂件" value="挂件" />
              <el-option label="音乐盒" value="音乐盒" />
              <el-option label="武器箱" value="武器箱" />
              <el-option label="布章" value="布章" />
            </el-select>
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
            </el-select>
          </div>

          <div>
            <button
              class="btn btn-primary"
              :disabled="skinStore.loading"
              @click="refreshData"
              style="height: 40px; padding: 0 24px; font-size: 14px;"
            >
              {{ skinStore.loading ? '刷新中...' : '确定并搜索' }}
            </button>
          </div>

          <div style="margin-left: auto;">
            <button
              class="btn btn-refresh"
              :disabled="skinStore.loading"
              @click="reloadData"
              style="height: 40px; padding: 0 24px; font-size: 14px;"
            >
              <el-icon style="margin-right: 6px;"><Refresh /></el-icon>
              {{ skinStore.loading ? '刷新中...' : '刷新数据' }}
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
          :data="skinStore.skinItems"
          v-loading="skinStore.loading"
          style="width: 100%"
        >
        <el-table-column type="index" label="#" width="60" />
        
        <el-table-column label="饰品名称" min-width="200" class-name="skin-column">
          <template #default="{ row }">
            <div class="skin-info">
              <img :src="row.image_url" @error="handleImageError" alt="饰品" class="skin-img" />
              <div>
                <div class="skin-name clickable" @click="goToDetail(row.market_hash_name)">{{ row.name }}</div>
                <div style="font-size: 12px; color: #8c8c8c;">{{ row.category }}</div>
              </div>
            </div>
          </template>
        </el-table-column>
        
          <el-table-column :label="sourcePriceLabel" width="160" class-name="price-column-left">
          <template #default="{ row }">
            <div style="display: flex; align-items: center; gap: 6px;">
              <img 
                :src="getPlatformIcon(sourcePlatform)" 
                style="width: 28px; height: 28px; object-fit: contain; cursor: pointer;" 
                :alt="sourcePlatform"
                @click="openPlatformLink(row, sourcePlatform)"
              />
              <div style="display: flex; flex-direction: column; gap: 2px;">
                <span style="color: #52c41a; font-weight: 600; font-size: 14px;">¥{{ formatPrice(row.source_price) }}</span>
                <span style="color: #8c8c8c; font-size: 11px;">{{ formatTimeAgo(row.source_update_time) }}</span>
              </div>
            </div>
          </template>
        </el-table-column>
        
          <el-table-column :label="targetPriceLabel" width="160" class-name="price-column-left">
          <template #default="{ row }">
            <div style="display: flex; align-items: center; gap: 6px;">
              <img 
                :src="getPlatformIcon(targetPlatform)" 
                style="width: 28px; height: 28px; object-fit: contain; cursor: pointer;" 
                :alt="targetPlatform"
                @click="openPlatformLink(row, targetPlatform)"
              />
              <div style="display: flex; flex-direction: column; gap: 2px;">
                <span style="color: #1890ff; font-weight: 600; font-size: 14px;">¥{{ formatPrice(row.target_price) }}</span>
                <span style="color: #8c8c8c; font-size: 11px;">{{ formatTimeAgo(row.target_update_time) }}</span>
              </div>
            </div>
          </template>
        </el-table-column>
        
          <el-table-column width="140" align="center">
          <template #header>
            <div style="display: flex; align-items: center; justify-content: center; gap: 4px;">
              <span>{{ turnOverLabel }}</span>
              <el-tooltip :content="turnOverTooltip" placement="top">
                <el-icon style="color: #8c8c8c; cursor: help; font-size: 14px;"><QuestionFilled /></el-icon>
              </el-tooltip>
            </div>
          </template>
          <template #default="{ row }">
              <span style="color: #595959; font-weight: 500; font-size: 14px;">{{ row.turn_over || 0 }}</span>
          </template>
        </el-table-column>
        
          <el-table-column width="140" align="center">
          <template #header>
            <div style="display: flex; align-items: center; justify-content: center; gap: 4px;">
              <span>{{ sellCountLabel }}</span>
              <el-tooltip :content="sellCountTooltip" placement="top">
                <el-icon style="color: #8c8c8c; cursor: help; font-size: 14px;"><QuestionFilled /></el-icon>
              </el-tooltip>
            </div>
          </template>
          <template #default="{ row }">
              <span style="color: #595959; font-weight: 500; font-size: 14px;">{{ row.sell_count || 0 }}</span>
          </template>
        </el-table-column>
        
          <el-table-column label="价格差" width="120" class-name="price-column-left">
          <template #default="{ row }">
              <span style="color: #faad14; font-weight: 600; font-size: 15px;">¥{{ formatPrice(row.price_diff) }}</span>
          </template>
        </el-table-column>
        
          <el-table-column width="120" class-name="price-column-left">
          <template #header>
            <div style="display: flex; align-items: center; gap: 4px;">
              <span>利润率</span>
              <el-tooltip content="等于价格差除以买入价格" placement="top">
                <el-icon style="color: #8c8c8c; cursor: help; font-size: 14px;"><QuestionFilled /></el-icon>
              </el-tooltip>
            </div>
          </template>
          <template #default="{ row }">
              <span :class="getProfitTagClass(row.profit_rate)" style="font-size: 15px !important;">
                {{ formatPercent(row.profit_rate) }}
              </span>
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
                  <div v-if="hasLowerPriceThanSource(row)" class="lower-platform-tip-inline">有更低价格平台</div>
                  <span class="platform-data-link">平台数据</span>
                </div>
              </template>
              <div style="padding: 8px; max-height: 500px; overflow-y: auto;">
                <!-- 卖出平台数据（抬头） -->
                <div v-if="getTargetPlatformData(row)" style="padding: 12px; background: #f5f7fa; border-radius: 8px; margin-bottom: 12px;">
                  <div style="display: flex; align-items: center; gap: 20px;">
                    <div style="display: flex; align-items: center; gap: 10px;">
                      <img 
                        :src="getPlatformIcon(targetPlatform)" 
                        style="width: 20px; height: 20px; object-fit: contain; cursor: pointer;" 
                        :alt="targetPlatform"
                        @click="openPlatformLink(row, targetPlatform)"
                      />
                      <span style="font-weight: 600; font-size: 14px; color: #262626;">{{ getPlatformDisplayName(targetPlatform) }}</span>
                    </div>
                    <div style="display: flex; gap: 40px; flex: 1; justify-content: center;">
                      <span style="font-size: 13px; white-space: nowrap;">在售价(¥): <strong style="color: #52c41a;">{{ formatPrice(getTargetPlatformData(row)?.sellPrice) }}</strong>({{ getTargetPlatformData(row)?.sellCount }})</span>
                      <span style="font-size: 13px; white-space: nowrap;">求购价(¥): <strong style="color: #ff4d4f;">{{ formatPrice(getTargetPlatformData(row)?.biddingPrice) }}</strong>({{ getTargetPlatformData(row)?.biddingCount }})</span>
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
                      <th style="padding: 10px 12px; text-align: left; font-weight: 600; color: #595959;">价格差(¥)</th>
                    </tr>
                  </thead>
                  <tbody>
                    <tr 
                      v-for="platform in getOtherPlatforms(row)" 
                      :key="platform.platformItemId"
                      style="border-bottom: 1px solid #f0f0f0;"
                    >
                      <td style="padding: 10px 12px;">
                        <div style="display: flex; align-items: center; gap: 8px;">
                          <img 
                            :src="getPlatformIconByName(platform.platformName)" 
                            style="width: 18px; height: 18px; object-fit: contain; cursor: pointer;" 
                            @click="openPlatformLinkByData(platform)"
                            :alt="platform.platformName"
                          />
                          <span>{{ platform.platformName }}</span>
                        </div>
                      </td>
                      <td style="padding: 10px 12px; color: #52c41a; font-weight: 500;">
                        {{ formatPrice(platform.sellPrice) }}<span style="color: #8c8c8c;">({{ platform.sellCount }})</span>
                      </td>
                      <td style="padding: 10px 12px; color: #ff4d4f; font-weight: 500;">
                        {{ formatPrice(platform.biddingPrice) }}<span style="color: #8c8c8c;">({{ platform.biddingCount }})</span>
                      </td>
                      <td style="padding: 10px 12px; color: #faad14; font-weight: 500;">{{ formatPrice(platform.price_diff) }}</td>
                    </tr>
                  </tbody>
                </table>
              </div>
            </el-popover>
            <span v-else style="color: #bfbfbf; font-size: 12px;">-</span>
          </template>
        </el-table-column>
        </el-table>
      </div>

      <!-- 分页 -->
      <div class="pagination">
        <div class="pagination-info">共 {{ skinStore.total }} 条数据</div>
        
        <div class="pagination-controls">
          <div class="page-size">
            <span>每页</span>
            <select v-model="skinStore.pagination.page_size" @change="handleSizeChange">
              <option :value="20">20</option>
              <option :value="50">50</option>
              <option :value="100">100</option>
              <option :value="200">200</option>
            </select>
            <span>条</span>
        </div>
        
          <button class="page-btn" :disabled="skinStore.pagination.page_num === 1" @click="handleCurrentChange(skinStore.pagination.page_num - 1)">‹</button>
          <button
              v-for="page in visiblePages"
              :key="page"
            class="page-btn"
            :class="{ active: page === skinStore.pagination.page_num, ellipsis: page === '...' }"
            :disabled="page === '...'"
              @click="page !== '...' && handleCurrentChange(page as number)"
            >
              {{ page }}
          </button>
          <button class="page-btn" :disabled="skinStore.pagination.page_num >= totalPages" @click="handleCurrentChange(skinStore.pagination.page_num + 1)">›</button>
        </div>
        </div>
      </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { Search, QuestionFilled, Refresh, DataAnalysis } from '@element-plus/icons-vue'
import { useSkinStore } from '@/stores/skin'
import { useSettingsStore } from '@/stores/settings'
import { dataApi } from '@/api'
import { formatPrice, formatPercent, debounce } from '@/utils'
import { showMessage } from '@/utils/message'
import type { PlatformData } from '@/types'
import buffIcon from '@/assets/icons/buff.png'
import uuIcon from '@/assets/icons/uu.png'
import c5Icon from '@/assets/icons/c5.png'
import steamIcon from '@/assets/icons/steam.png'

const router = useRouter()
const skinStore = useSkinStore()
const settingsStore = useSettingsStore()

const searchKeyword = ref('')
// 使用 store 中的平台和排序设置（会自动持久化到 localStorage）
const sortOption = computed({
  get: () => skinStore.sortOption,
  set: (val) => { skinStore.sortOption = val }
})
const sourcePlatform = computed({
  get: () => skinStore.sourcePlatform,
  set: (val) => { skinStore.sourcePlatform = val }
})
const targetPlatform = computed({
  get: () => skinStore.targetPlatform,
  set: (val) => { skinStore.targetPlatform = val }
})
const category = computed({
  get: () => skinStore.category,
  set: (val) => { skinStore.category = val }
})

// 筛选参数
const filterParams = reactive({
  min_sell_num: 0,
  min_diff: 0,
  max_sell_price: 5000,
  min_sell_price: 2,
})

const totalPages = computed(() => {
  return Math.ceil(skinStore.total / skinStore.pagination.page_size)
})

const visiblePages = computed(() => {
  const current = skinStore.pagination.page_num
  const total = totalPages.value
  const pages: (number | string)[] = []
  
  if (total <= 7) {
    for (let i = 1; i <= total; i++) {
      pages.push(i)
    }
  } else {
    pages.push(1)
    
    if (current <= 4) {
      for (let i = 2; i <= 5; i++) {
        pages.push(i)
      }
      pages.push('...')
      pages.push(total)
    } else if (current >= total - 3) {
      pages.push('...')
      for (let i = total - 4; i <= total; i++) {
        pages.push(i)
      }
    } else {
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

const handleSearch = debounce(() => {
  skinStore.getSkinItems({
    search: searchKeyword.value,
    page_num: 1
  })
}, 300)

const handleSortChange = () => {
  // 排序配置在 store 中管理，直接刷新数据
  skinStore.getSkinItems({
    search: searchKeyword.value,
    page_num: 1
  })
}

const handleCategoryChange = () => {
  // 类别配置在 store 中管理，直接刷新数据
  skinStore.getSkinItems({
    search: searchKeyword.value,
    page_num: 1
  })
}

const handleSourceChange = () => {
  // 如果买入平台和卖出平台相同，自动把卖出平台换成另一个
  if (sourcePlatform.value === targetPlatform.value) {
    const platforms = ['buff', 'uu', 'c5', 'steam']
    const otherPlatforms = platforms.filter(p => p !== sourcePlatform.value)
    targetPlatform.value = otherPlatforms[0]
  }
  
  // 刷新数据（store 会自动使用保存的平台设置）
  skinStore.getSkinItems({ 
    search: searchKeyword.value,
    page_num: 1
  })
}

const handleTargetChange = () => {
  // 如果卖出平台和买入平台相同，自动把买入平台换成另一个
  if (targetPlatform.value === sourcePlatform.value) {
    const platforms = ['buff', 'uu', 'c5', 'steam']
    const otherPlatforms = platforms.filter(p => p !== targetPlatform.value)
    sourcePlatform.value = otherPlatforms[0]
  }
  
  // 刷新数据（store 会自动使用保存的平台设置）
  skinStore.getSkinItems({
    search: searchKeyword.value,
    page_num: 1
  })
}

const refreshData = async () => {
  // 先保存设置
  try {
    await settingsStore.updateSettings(filterParams)
  } catch (error) {
    console.error('保存设置失败:', error)
  }
  
  // 再刷新数据（store 会自动使用保存的平台设置）
  await skinStore.getSkinItems({})
}

// 仅刷新数据（不保存设置）
const reloadData = async () => {
  await skinStore.getSkinItems({})
}

// 获取筛选描述
const getFilterDescription = () => {
  const platformNames: Record<string, string> = {
    'buff': 'BUFF',
    'uu': '悠悠',
    'c5': 'C5GAME',
    'steam': 'Steam'
  }
  const sourceName = platformNames[sourcePlatform.value] || sourcePlatform.value
  const targetName = platformNames[targetPlatform.value] || targetPlatform.value
  const priceRange = `价格${filterParams.min_sell_price}-${filterParams.max_sell_price}元`
  const sellNum = filterParams.min_sell_num > 0 ? `，在售量大于${filterParams.min_sell_num}件` : ''
  const priceDiff = filterParams.min_diff > 0 ? `，价格差大于${filterParams.min_diff}元` : ''
  const categoryDesc = category.value ? `，类别：${category.value}` : ''
  
  return `从${sourceName}平台买入饰品，到${targetName}平台卖出（${priceRange}${sellNum}${priceDiff}${categoryDesc}）`
}

// 加载用户设置
const loadSettings = async () => {
  await settingsStore.getSettings()
  Object.assign(filterParams, settingsStore.settings)
}

const handleSizeChange = () => {
  skinStore.getSkinItems({ page_num: 1 })
}

const handleCurrentChange = (page: number) => {
  if (page >= 1 && page <= totalPages.value) {
    skinStore.getSkinItems({ page_num: page })
  }
}

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

// 获取平台图标
const getPlatformIcon = (platform: string) => {
  const iconMap: Record<string, string> = {
    'buff': buffIcon,
    'uu': uuIcon,
    'c5': c5Icon,
    'steam': steamIcon
  }
  return iconMap[platform] || ''
}

// 获取平台选择器显示标签
const getPlatformSelectLabel = (platform: string) => {
  const labelMap: Record<string, string> = {
    'buff': 'BUFF',
    'uu': '悠悠',
    'c5': 'C5GAME',
    'steam': 'Steam'
  }
  return labelMap[platform] || platform
}

// 格式化时间为"X小时前"
const formatTimeAgo = (timestamp: number) => {
  if (!timestamp) return ''
  
  const now = Math.floor(Date.now() / 1000)
  const diff = now - timestamp
  
  if (diff < 60) {
    return '刚刚'
  } else if (diff < 3600) {
    const minutes = Math.floor(diff / 60)
    return `${minutes}分钟前`
  } else if (diff < 86400) {
    const hours = Math.floor(diff / 3600)
    return `${hours}小时前`
  } else {
    const days = Math.floor(diff / 86400)
    return `${days}天前`
  }
}

// 打开平台链接
const openPlatformLink = (row: any, platform: string) => {
  const platformData = row.platform_list?.find((p: any) => {
    const platformNameMap: Record<string, string> = {
      'buff': 'BUFF',
      'uu': '悠悠',
      'c5': 'C5GAME',
      'steam': 'Steam'  // 后端返回的是Steam而不是STEAM
    }
    return p.platformName === platformNameMap[platform]
  })
  
  if (platformData?.link) {
    window.open(platformData.link, '_blank')
  }
}

// 通过平台数据直接打开链接
const openPlatformLinkByData = (platform: PlatformData) => {
  if (platform?.link) {
    window.open(platform.link, '_blank')
  }
}

// 获取卖出平台数据
const getTargetPlatformData = (row: any) => {
  const platformNameMap: Record<string, string> = {
    'buff': 'BUFF',
    'uu': '悠悠',
    'c5': 'C5GAME',
    'steam': 'Steam'  // 后端返回的是Steam而不是STEAM
  }
  return row.platform_list?.find((p: any) => p.platformName === platformNameMap[targetPlatform.value])
}

// 获取其他平台数据（排除卖出平台）
const getOtherPlatforms = (row: any) => {
  const platformNameMap: Record<string, string> = {
    'buff': 'BUFF',
    'uu': '悠悠',
    'c5': 'C5GAME',
    'steam': 'Steam'  // 后端返回的是Steam而不是STEAM
  }
  const targetName = platformNameMap[targetPlatform.value]
  return row.platform_list?.filter((p: any) => p.platformName !== targetName) || []
}

// 判断是否有平台价格低于买入平台价格
const hasLowerPriceThanSource = (row: any) => {
  if (!row.platform_list || !row.source_price || row.source_price <= 0) return false
  // 检查是否有其他平台的卖出价低于买入平台的价格
  return row.platform_list.some((p: any) => 
    p.sellPrice > 0 && p.sellPrice < row.source_price
  )
}

// 获取平台显示名称
const getPlatformDisplayName = (platform: string) => {
  const platformNames: Record<string, string> = {
    'buff': 'BUFF',
    'uu': '悠悠',
    'c5': 'C5GAME',
    'steam': 'Steam'
  }
  return platformNames[platform] || platform
}

// 根据平台名称获取图标
const getPlatformIconByName = (platformName: string) => {
  const iconMap: Record<string, string> = {
    'BUFF': buffIcon,
    '悠悠': uuIcon,
    'C5GAME': c5Icon,
    'STEAM': steamIcon,
    'Steam': steamIcon  // 后端返回的是Steam而不是STEAM
  }
  return iconMap[platformName] || ''
}

// 动态列名
const sourcePriceLabel = computed(() => {
  const platformLabels: Record<string, string> = {
    'buff': 'BUFF最低售价(¥)',
    'uu': '悠悠最低售价(¥)',
    'c5': 'C5最低售价(¥)',
    'steam': 'Steam最低售价(¥)'
  }
  return platformLabels[sourcePlatform.value] || '买入最低售价(¥)'
})

const targetPriceLabel = computed(() => {
  const platformLabels: Record<string, string> = {
    'buff': 'BUFF最低售价(¥)',
    'uu': '悠悠最低售价(¥)',
    'c5': 'C5最低售价(¥)',
    'steam': 'Steam最低售价(¥)'
  }
  return platformLabels[targetPlatform.value] || '卖出最低售价(¥)'
})

const turnOverLabel = computed(() => {
  return '成交量'
})

const turnOverTooltip = computed(() => {
  const platformLabels: Record<string, string> = {
    'buff': 'BUFF',
    'uu': '悠悠',
    'c5': 'C5GAME',
    'steam': 'Steam'
  }
  const platformName = platformLabels[targetPlatform.value] || ''
  return `${platformName}平台近12小时成交量`
})

const sellCountLabel = computed(() => {
  return '在售数量'
})

const sellCountTooltip = computed(() => {
  const platformLabels: Record<string, string> = {
    'buff': 'BUFF',
    'uu': '悠悠',
    'c5': 'C5GAME',
    'steam': 'Steam'
  }
  const platformName = platformLabels[targetPlatform.value] || ''
  return `${platformName}平台目前在售数量`
})

const getProfitTagClass = (rate: number) => {
  // 将小数转换为百分比（rate是0-1的小数）
  const percent = rate * 100
  
  if (percent >= 60) return 'tag tag-danger'    // 红色：>= 60%
  if (percent >= 40) return 'tag tag-warning'   // 橙色：40% - 60%
  if (percent >= 20) return 'tag tag-success'   // 绿色：20% - 40%
  return 'tag tag-primary'                       // 蓝色：< 20%
}

onMounted(async () => {
  await loadSettings()
  skinStore.getSkinItems({})
})
</script>

<style scoped>
/* 所有样式在unified.css中 */
.home-page {
  padding: 0;
  width: 1800px;
  max-width: 100%;
  margin: 0 auto;
}

/* 平台数据单元格 */
.platform-data-cell {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
}

.lower-platform-tip-inline {
  background: linear-gradient(135deg, #ff6b6b 0%, #ee5a5a 100%);
  color: #fff;
  text-align: center;
  padding: 3px 8px;
  border-radius: 4px;
  font-size: 11px;
  font-weight: 500;
  white-space: nowrap;
}

.platform-data-link {
  color: #1890ff;
  cursor: pointer;
  font-size: 13px;
}

.platform-data-link:hover {
  color: #40a9ff;
}

/* 饰品名称点击样式 */
.skin-name.clickable {
  cursor: pointer;
  transition: color 0.2s;
}

.skin-name.clickable:hover {
  color: #1890ff;
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
</style>

