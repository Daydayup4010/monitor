<template>
  <div class="user-manager">
    <!-- 统计卡片 -->
    <div class="stats-grid">
      <el-card class="stat-card">
        <div class="stat-content">
          <div class="stat-icon blue">
            <el-icon><User /></el-icon>
          </div>
          <div class="stat-info">
            <div class="stat-value">{{ total }}</div>
            <div class="stat-label">总用户数</div>
          </div>
        </div>
      </el-card>

      <el-card class="stat-card">
        <div class="stat-content">
          <div class="stat-icon green">
            <el-icon><Medal /></el-icon>
          </div>
          <div class="stat-info">
            <div class="stat-value">{{ vipCount }}</div>
            <div class="stat-label">VIP用户数</div>
          </div>
        </div>
      </el-card>
    </div>

    <!-- 搜索和操作栏 -->
    <div class="action-bar">
      <el-input
        v-model="searchKeyword"
        placeholder="搜索用户名或邮箱"
        style="width: 300px"
        clearable
        @input="handleSearch"
      >
        <template #prefix>
          <el-icon><Search /></el-icon>
        </template>
      </el-input>
    </div>

    <!-- 用户列表表格 -->
    <el-table
      :data="userList"
      v-loading="loading"
      stripe
      style="width: 100%; margin-top: 20px"
      class="user-table"
    >
      <el-table-column prop="id" label="ID" width="280" />
      <el-table-column prop="user_name" label="用户名" width="150" />
      <el-table-column prop="email" label="邮箱" width="200" />
      <el-table-column label="用户类型" width="120">
        <template #default="{ row }">
          <el-tag v-if="row.role === 2" type="danger">管理员</el-tag>
          <el-tag v-else-if="row.role === 1 && isVipValid(row.vip_expiry)" type="success">VIP会员</el-tag>
          <el-tag v-else-if="row.role === 1" type="warning">VIP已过期</el-tag>
          <el-tag v-else type="info">普通用户</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="vip_expiry" label="VIP到期时间" width="180">
        <template #default="{ row }">
          <span v-if="row.vip_expiry && row.role === 1">
            {{ formatDate(row.vip_expiry) }}
          </span>
          <span v-else style="color: #999">-</span>
        </template>
      </el-table-column>
      <el-table-column prop="last_login" label="最后登录" width="180">
        <template #default="{ row }">
          {{ row.last_login ? formatDate(row.last_login) : '-' }}
        </template>
      </el-table-column>
      <el-table-column label="操作" width="200" fixed="right">
        <template #default="{ row }">
          <el-button
            v-if="row.role !== 2"
            type="primary"
            size="small"
            @click="handleRenewVip(row)"
          >
            {{ row.role === 1 && isVipValid(row.vip_expiry) ? '续费VIP' : '开通VIP' }}
          </el-button>
          <el-button
            v-if="row.role !== 2"
            type="danger"
            size="small"
            @click="handleDelete(row)"
          >
            删除
          </el-button>
        </template>
      </el-table-column>
    </el-table>

    <!-- 分页 -->
    <div class="pagination-container">
      <el-pagination
        v-model:current-page="pageNum"
        v-model:page-size="pageSize"
        :page-sizes="[10, 20, 50, 100]"
        :total="total"
        layout="total, sizes, prev, pager, next, jumper"
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
      />
    </div>

    <!-- VIP续费对话框 -->
    <el-dialog
      v-model="renewDialogVisible"
      :title="renewDialogTitle"
      width="500px"
    >
      <el-form :model="renewForm" label-width="100px">
        <el-form-item label="用户名">
          <el-input v-model="currentUser.user_name" disabled />
        </el-form-item>
        <el-form-item label="邮箱">
          <el-input v-model="currentUser.email" disabled />
        </el-form-item>
        <el-form-item label="续费天数">
          <el-input-number
            v-model="renewForm.days"
            :min="1"
            :max="3650"
            controls-position="right"
            style="width: 100%"
          />
          <div style="margin-top: 8px; font-size: 12px; color: #999">
            当前到期时间：{{ currentUser.vip_expiry ? formatDate(currentUser.vip_expiry) : '未开通' }}
          </div>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="renewDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="renewLoading" @click="confirmRenew">
          确认
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessageBox } from 'element-plus'
import { userApi } from '@/api'
import { showMessage } from '@/utils/message'
import type { UserListItem } from '@/types'
import dayjs from 'dayjs'
import { debounce } from '@/utils'

const loading = ref(false)
const userList = ref<UserListItem[]>([])
const total = ref(0)
const pageNum = ref(1)
const pageSize = ref(10)
const searchKeyword = ref('')

const renewDialogVisible = ref(false)
const renewLoading = ref(false)
const currentUser = ref<UserListItem>({} as UserListItem)
const renewForm = reactive({
  days: 30,
})

// VIP用户数量
const vipCount = computed(() => {
  return userList.value.filter(user => 
    user.role === 1 && isVipValid(user.vip_expiry)
  ).length
})

// 续费对话框标题
const renewDialogTitle = computed(() => {
  if (currentUser.value.role === 1 && isVipValid(currentUser.value.vip_expiry)) {
    return 'VIP续费'
  }
  return '开通VIP'
})

// 检查VIP是否有效
const isVipValid = (expiryDate?: string): boolean => {
  if (!expiryDate) return false
  return new Date(expiryDate) > new Date()
}

// 格式化日期
const formatDate = (date: string) => {
  return dayjs(date).format('YYYY-MM-DD HH:mm')
}

// 加载用户列表
const loadUserList = async () => {
  loading.value = true
  try {
    const response = await userApi.getUserList({
      page_num: pageNum.value,
      page_size: pageSize.value,
      search: searchKeyword.value,
    })
    if (response.code === 1) {
      userList.value = response.data || []
      total.value = response.total || 0
    }
  } catch (error) {
    console.error('加载用户列表失败:', error)
  } finally {
    loading.value = false
  }
}

// 搜索
const handleSearch = debounce(() => {
  pageNum.value = 1
  loadUserList()
}, 300)

// 页面大小变化
const handleSizeChange = () => {
  pageNum.value = 1
  loadUserList()
}

// 当前页变化
const handleCurrentChange = () => {
  loadUserList()
}

// 续费VIP
const handleRenewVip = (user: UserListItem) => {
  currentUser.value = user
  renewForm.days = 30
  renewDialogVisible.value = true
}

// 确认续费
const confirmRenew = async () => {
  renewLoading.value = true
  try {
    const response = await userApi.renewVip({
      user_id: currentUser.value.id,
      days: renewForm.days,
    })
    if (response.code === 1) {
      showMessage.success('VIP续费成功')
      renewDialogVisible.value = false
      loadUserList()
    }
  } catch (error) {
    console.error('续费失败:', error)
  } finally {
    renewLoading.value = false
  }
}

// 删除用户
const handleDelete = (user: UserListItem) => {
  ElMessageBox.confirm(
    `确定要删除用户 "${user.user_name}" 吗？此操作不可恢复。`,
    '删除用户',
    {
      confirmButtonText: '确定删除',
      cancelButtonText: '取消',
      type: 'warning',
    }
  ).then(async () => {
    try {
      const response = await userApi.deleteUser(user.id)
      if (response.code === 1) {
        showMessage.success('用户删除成功')
        loadUserList()
      }
    } catch (error) {
      console.error('删除失败:', error)
    }
  }).catch(() => {
    // 用户取消
  })
}

onMounted(() => {
  loadUserList()
})
</script>

<style scoped>
.user-manager {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 20px;
  margin-bottom: 10px;
}

.stat-card {
  border-radius: 12px;
  border: none;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
}

.stat-content {
  display: flex;
  align-items: center;
  gap: 16px;
}

.stat-icon {
  width: 50px;
  height: 50px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24px;
  color: white;
}

.stat-icon.blue {
  background: linear-gradient(135deg, #1890ff, #40a9ff);
}

.stat-icon.green {
  background: linear-gradient(135deg, #52c41a, #73d13d);
}

.stat-info {
  flex: 1;
}

.stat-value {
  font-size: 32px;
  font-weight: bold;
  color: #333;
  line-height: 1.2;
}

.stat-label {
  font-size: 14px;
  color: #999;
  margin-top: 4px;
}

.action-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.user-table {
  border-radius: 12px;
  overflow: hidden;
}

.pagination-container {
  display: flex;
  justify-content: center;
  margin-top: 20px;
}

/* 响应式 */
@media (max-width: 768px) {
  .stats-grid {
    grid-template-columns: 1fr;
  }

  .action-bar {
    flex-direction: column;
    gap: 12px;
  }

  .action-bar :deep(.el-input) {
    width: 100% !important;
  }
}
</style>

