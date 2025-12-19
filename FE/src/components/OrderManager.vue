<template>
  <div class="order-manager">
    <div class="card">
      <div class="card-title">📋 订单管理</div>
      
      <!-- 筛选栏 -->
      <div class="filter-bar">
        <el-input
          v-model="keyword"
          placeholder="搜索订单号、用户名或邮箱"
          clearable
          style="width: 280px;"
          @input="handleSearch"
        >
          <template #prefix>
            <el-icon><Search /></el-icon>
          </template>
        </el-input>

        <el-select v-model="statusFilter" placeholder="订单状态" style="width: 140px;" @change="fetchOrders">
          <el-option label="全部状态" :value="-1" />
          <el-option label="待支付" :value="0" />
          <el-option label="已支付" :value="1" />
          <el-option label="已取消" :value="2" />
          <el-option label="已退款" :value="3" />
        </el-select>

        <el-button type="primary" @click="fetchOrders" :loading="loading">
          <el-icon style="margin-right: 4px;"><Refresh /></el-icon>
          刷新
        </el-button>
      </div>

      <!-- 订单表格 -->
      <div class="table-wrapper">
        <el-table :data="orders" v-loading="loading" style="width: 100%">
          <el-table-column prop="out_trade_no" label="订单号" width="200">
            <template #default="{ row }">
              <span class="order-no">{{ row.out_trade_no }}</span>
            </template>
          </el-table-column>
          
          <el-table-column label="用户" width="200">
            <template #default="{ row }">
              <div class="user-info">
                <div class="username">{{ row.username || '-' }}</div>
                <div class="email">{{ row.email || '-' }}</div>
              </div>
            </template>
          </el-table-column>
          
          <el-table-column label="套餐" width="100" align="center">
            <template #default="{ row }">
              <span class="months-tag">{{ row.months }}个月</span>
            </template>
          </el-table-column>
          
          <el-table-column label="金额" width="100" align="right">
            <template #default="{ row }">
              <span class="amount">¥{{ row.amount.toFixed(2) }}</span>
            </template>
          </el-table-column>
          
          <el-table-column label="状态" width="100" align="center">
            <template #default="{ row }">
              <el-tag :type="getStatusType(row.status)" size="small">
                {{ getStatusText(row.status) }}
              </el-tag>
            </template>
          </el-table-column>
          
          <el-table-column label="创建时间" width="170">
            <template #default="{ row }">
              <span class="time">{{ formatTime(row.created_at) }}</span>
            </template>
          </el-table-column>
          
          <el-table-column label="支付时间" width="170">
            <template #default="{ row }">
              <span class="time" v-if="row.pay_time">{{ formatTime(row.pay_time) }}</span>
              <span class="time empty" v-else>-</span>
            </template>
          </el-table-column>
          
          <el-table-column label="YunGouOS订单号" min-width="180">
            <template #default="{ row }">
              <span class="yun-order" v-if="row.yun_order_no">{{ row.yun_order_no }}</span>
              <span class="yun-order empty" v-else>-</span>
            </template>
          </el-table-column>
        </el-table>
      </div>

      <!-- 分页 -->
      <div class="pagination-wrapper">
        <div class="pagination-info">
          共 {{ total }} 条记录
        </div>
        <el-pagination
          v-model:current-page="pageNum"
          v-model:page-size="pageSize"
          :page-sizes="[10, 20, 50, 100]"
          :total="total"
          layout="sizes, prev, pager, next, jumper"
          @size-change="fetchOrders"
          @current-change="fetchOrders"
        />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Search, Refresh } from '@element-plus/icons-vue'
import { orderApi, type PaymentOrderItem } from '@/api'
import { debounce } from '@/utils'

const orders = ref<PaymentOrderItem[]>([])
const loading = ref(false)
const total = ref(0)
const pageNum = ref(1)
const pageSize = ref(20)
const statusFilter = ref(-1)
const keyword = ref('')

// 获取订单列表
const fetchOrders = async () => {
  loading.value = true
  try {
    const res = await orderApi.getAllOrders({
      page_size: pageSize.value,
      page_num: pageNum.value,
      status: statusFilter.value,
      keyword: keyword.value,
    })
    if (res.code === 1) {
      orders.value = res.data || []
      total.value = res.total || 0
    }
  } catch (error) {
    console.error('获取订单列表失败:', error)
  } finally {
    loading.value = false
  }
}

// 搜索防抖
const handleSearch = debounce(() => {
  pageNum.value = 1
  fetchOrders()
}, 300)

// 获取状态文本
const getStatusText = (status: number) => {
  const statusMap: Record<number, string> = {
    0: '待支付',
    1: '已支付',
    2: '已取消',
    3: '已退款',
  }
  return statusMap[status] || '未知'
}

// 获取状态标签类型
const getStatusType = (status: number) => {
  const typeMap: Record<number, string> = {
    0: 'warning',
    1: 'success',
    2: 'info',
    3: 'danger',
  }
  return typeMap[status] || 'info'
}

// 格式化时间
const formatTime = (time: string) => {
  if (!time) return '-'
  const date = new Date(time)
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit',
  })
}

onMounted(() => {
  fetchOrders()
})
</script>

<style scoped>
.order-manager {
  width: 100%;
}

.card {
  background: white;
  border-radius: 12px;
  padding: 24px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
}

.card-title {
  font-size: 18px;
  font-weight: 600;
  color: #262626;
  margin-bottom: 20px;
  padding-bottom: 16px;
  border-bottom: 1px solid #f0f0f0;
}

.filter-bar {
  display: flex;
  gap: 12px;
  margin-bottom: 20px;
  flex-wrap: wrap;
}

.table-wrapper {
  margin-bottom: 20px;
  border-radius: 8px;
  overflow: hidden;
  border: 1px solid #e8e8e8;
}

.order-no {
  font-family: monospace;
  font-size: 13px;
  color: #1890ff;
}

.user-info {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.username {
  font-weight: 500;
  color: #262626;
}

.email {
  font-size: 12px;
  color: #8c8c8c;
}

.months-tag {
  background: #e6f7ff;
  color: #1890ff;
  padding: 2px 8px;
  border-radius: 4px;
  font-size: 13px;
}

.amount {
  font-weight: 600;
  color: #52c41a;
  font-size: 14px;
}

.time {
  font-size: 13px;
  color: #595959;
}

.time.empty,
.yun-order.empty {
  color: #bfbfbf;
}

.yun-order {
  font-family: monospace;
  font-size: 12px;
  color: #8c8c8c;
}

.pagination-wrapper {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.pagination-info {
  font-size: 14px;
  color: #595959;
}

:deep(.el-table th) {
  background: #fafafa !important;
  font-weight: 600;
  color: #595959;
}

:deep(.el-table__row:hover > td) {
  background: #f5f7fa !important;
}
</style>


