<template>
  <div class="home-page">
    <div class="card">
      <div class="card-title">
        <img src="@/assets/icons/data.png" style="height: 20px; width: auto; vertical-align: text-top; margin-right: 8px; object-fit: contain; position: relative; top: 2px;" alt="饰品数据" />
        饰品数据
      </div>

      <!-- 筛选栏 -->
      <div class="filter-bar">
        <div class="filter-item" style="flex: 1;">
          <label class="filter-label">搜索</label>
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
          
        <div class="filter-item">
          <label class="filter-label">类别</label>
            <el-select
              v-model="selectedCategory"
            placeholder="全部"
              clearable
            style="width: 150px;"
            >
              <el-option
                v-for="category in categories"
                :key="category"
                :label="category"
                :value="category"
              />
            </el-select>
          </div>
          
        <div class="filter-item">
          <label class="filter-label">排序</label>
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
          
        <div class="filter-item" style="padding-top: 20px;">
          <button
            class="btn btn-primary"
            :disabled="skinStore.loading"
            @click="refreshData"
            style="height: 42px; padding: 0 24px; font-size: 14px;"
          >
            {{ skinStore.loading ? '刷新中...' : '刷新数据' }}
          </button>
        </div>
    </div>

    <!-- 数据表格 -->
      <div class="table-wrapper">
        <el-table
          :data="skinStore.skinItems"
          v-loading="skinStore.loading"
          style="width: 100%"
        >
        <el-table-column type="index" label="#" width="50" />
        
        <el-table-column label="饰品名称" class-name="skin-column">
          <template #default="{ row }">
            <div class="skin-info">
              <img :src="row.image_url" @error="handleImageError" alt="饰品" class="skin-img" />
              <div>
                <div class="skin-name">{{ row.name }}</div>
                <div style="font-size: 12px; color: #8c8c8c;">{{ row.category }}</div>
              </div>
            </div>
          </template>
        </el-table-column>
        
          <el-table-column label="UU价格" width="120">
          <template #default="{ row }">
              <span style="color: #1890ff; font-weight: 600;">¥{{ formatPrice(row.u_price) }}</span>
          </template>
        </el-table-column>
        
          <el-table-column label="Buff价格" width="120">
          <template #default="{ row }">
              <span style="color: #52c41a; font-weight: 600;">¥{{ formatPrice(row.buff_price) }}</span>
          </template>
        </el-table-column>
        
          <el-table-column label="价格差" width="100">
          <template #default="{ row }">
              <span style="color: #faad14; font-weight: 600;">¥{{ formatPrice(row.price_diff) }}</span>
          </template>
        </el-table-column>
        
          <el-table-column label="利润率" width="100">
          <template #default="{ row }">
              <span :class="getProfitTagClass(row.profit_rate)">
                {{ formatPercent(row.profit_rate) }}
              </span>
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
import { ref, computed, onMounted } from 'vue'
import { useSkinStore } from '@/stores/skin'
import { formatPrice, formatPercent, debounce } from '@/utils'

const skinStore = useSkinStore()

const searchKeyword = ref('')
const selectedCategory = ref('')
const sortOption = ref('default')

const categories = computed(() => {
  const cats = new Set(skinStore.skinItems.map(item => item.category))
  return Array.from(cats).filter(Boolean)
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
  // 搜索逻辑
}, 300)

const handleSortChange = () => {
  const sortMap: Record<string, { field: string; desc: boolean }> = {
    'default': { field: '', desc: false },
    'price_diff_asc': { field: 'price_diff', desc: false },
    'price_diff_desc': { field: 'price_diff', desc: true },
    'profit_rate_asc': { field: 'profit_rate', desc: false },
    'profit_rate_desc': { field: 'profit_rate', desc: true },
  }
  
  const config = sortMap[sortOption.value]
    skinStore.getSkinItems({ 
    sort: config.field,
    desc: config.desc,
      page_num: 1
    })
  }

const refreshData = async () => {
  await skinStore.getSkinItems()
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

const getProfitTagClass = (rate: number) => {
  if (rate > 0.05) return 'tag tag-success'
  if (rate > 0.03) return 'tag tag-warning'
  return 'tag tag-primary'
}

onMounted(() => {
  skinStore.getSkinItems()
})
</script>

<style scoped>
/* 所有样式在unified.css中 */
.home-page {
  padding: 0;
  max-width: 100%;
}
</style>

