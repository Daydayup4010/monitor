<template>
  <div class="user-manager">
    <!-- 统计卡片 -->
    <div class="stats">
      <div class="stat-card">
        <div class="stat-icon blue">👥</div>
        <div class="stat-info">
          <div class="stat-value">{{ total }}</div>
          <div class="stat-label">总用户数</div>
        </div>
      </div>

      <div class="stat-card">
        <div class="stat-icon green"><img :src="vipIcon" class="stat-icon-img" /></div>
        <div class="stat-info">
          <div class="stat-value">{{ vipCount }}</div>
          <div class="stat-label">VIP用户</div>
        </div>
      </div>
    </div>

    <!-- 用户管理卡片 -->
    <div class="card">
      <div class="card-title">👥 用户管理</div>

      <div class="search-box">
        <el-input
          v-model="searchKeyword"
          placeholder="🔍 搜索用户名或邮箱..."
          style="width: 320px"
          clearable
          @input="handleSearch"
        />
      </div>

      <div class="table-wrapper">
        <el-table
          :data="userList"
          v-loading="loading"
          style="width: 100%"
        >
          <el-table-column type="index" label="#" width="60" />
          <el-table-column prop="user_name" label="用户名" />
          <el-table-column prop="email" label="邮箱" />
          <el-table-column label="类型" width="120">
            <template #default="{ row }">
              <span v-if="row.role === 2" class="tag tag-danger" style="display: inline-block; white-space: nowrap;">管理员</span>
              <span v-else-if="row.role === 1 && isVipValid(row.vip_expiry)" class="tag tag-success" style="display: inline-block; white-space: nowrap;">VIP会员</span>
              <span v-else-if="row.role === 1" class="tag tag-warning" style="display: inline-block; white-space: nowrap;">VIP已过期</span>
              <span v-else class="tag tag-info" style="display: inline-block; white-space: nowrap;">普通用户</span>
            </template>
          </el-table-column>
          <el-table-column label="VIP到期" width="150">
            <template #default="{ row }">
              <span v-if="row.vip_expiry && isVipValid(row.vip_expiry)">
                {{ formatDate(row.vip_expiry) }}
              </span>
              <span v-else-if="row.vip_expiry && !isVipValid(row.vip_expiry)" style="color: #ff4d4f">
                {{ formatDate(row.vip_expiry) }} (已过期)
              </span>
              <span v-else style="color: #bfbfbf">-</span>
            </template>
          </el-table-column>
          <el-table-column label="最后登录" width="170">
            <template #default="{ row }">
              <span v-if="row.last_login">
                {{ formatDateTime(row.last_login) }}
              </span>
              <span v-else style="color: #bfbfbf">-</span>
            </template>
          </el-table-column>
          <el-table-column label="操作" width="200">
            <template #default="{ row }">
              <div v-if="row.role !== 2" style="display: flex; gap: 8px;">
                <button
                  class="btn btn-primary"
                  style="font-size: 13px; padding: 6px 12px;"
                  @click="handleRenewVip(row)"
                >
                  {{ row.role === 1 && isVipValid(row.vip_expiry) ? '续费VIP' : '开通VIP' }}
                </button>
                <button
                  class="btn btn-secondary"
                  style="font-size: 13px; padding: 6px 12px;"
                  @click="handleDelete(row)"
                >
                  删除
                </button>
              </div>
            </template>
          </el-table-column>
        </el-table>
      </div>

      <!-- 分页 -->
      <div class="pagination">
        <div class="pagination-info">共 {{ total }} 条用户</div>
        
        <div class="pagination-controls">
          <div class="page-size">
            <span>每页</span>
            <select v-model="pageSize" @change="handleSizeChange">
              <option :value="10">10</option>
              <option :value="20">20</option>
              <option :value="50">50</option>
              <option :value="100">100</option>
            </select>
            <span>条</span>
          </div>

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

    <!-- VIP续费对话框 -->
    <el-dialog
      v-model="renewDialogVisible"
      :title="renewDialogTitle"
      width="500px"
    >
      <el-form :model="renewForm" label-width="100px">
        <div class="form-item">
          <label class="form-label">用户名</label>
          <el-input :model-value="currentUser.user_name" disabled />
        </div>
        <div class="form-item">
          <label class="form-label">邮箱</label>
          <el-input :model-value="currentUser.email" disabled />
        </div>
        <div class="form-item">
          <label class="form-label">续费月数</label>
          <el-input-number
            v-model="renewForm.days"
            :min="1"
            :max="120"
            controls-position="right"
            style="width: 100%"
          />
        </div>
      </el-form>
      <template #footer>
        <button class="btn btn-secondary" @click="renewDialogVisible = false">取消</button>
        <button class="btn btn-primary" :disabled="renewLoading" @click="confirmRenew" style="margin-left: 12px;">
          确认
        </button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessageBox } from 'element-plus'
import { userApi } from '@/api'
import { showMessage } from '@/utils/message'
import { debounce } from '@/utils'
import type { UserListItem } from '@/types'
import dayjs from 'dayjs'
import vipIcon from '@/assets/icons/vip.png'

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
  days: 1,  // 默认1个月
})

const vipCount = computed(() => {
  return userList.value.filter(user => 
    user.role === 1 && isVipValid(user.vip_expiry)
  ).length
})

const renewDialogTitle = computed(() => {
  if (currentUser.value.role === 1 && isVipValid(currentUser.value.vip_expiry)) {
    return 'VIP续费'
  }
  return '开通VIP'
})

const totalPages = computed(() => {
  return Math.ceil(total.value / pageSize.value)
})

const visiblePages = computed(() => {
  const current = pageNum.value
  const totalPgs = totalPages.value
  const pages: (number | string)[] = []
  
  if (totalPgs <= 7) {
    for (let i = 1; i <= totalPgs; i++) {
      pages.push(i)
    }
  } else {
    pages.push(1)
    
    if (current <= 4) {
      for (let i = 2; i <= 5; i++) {
        pages.push(i)
      }
      pages.push('...')
      pages.push(totalPgs)
    } else if (current >= totalPgs - 3) {
      pages.push('...')
      for (let i = totalPgs - 4; i <= totalPgs; i++) {
        pages.push(i)
      }
    } else {
      pages.push('...')
      for (let i = current - 1; i <= current + 1; i++) {
        pages.push(i)
      }
      pages.push('...')
      pages.push(totalPgs)
    }
  }
  
  return pages
})

const isVipValid = (expiryDate?: string): boolean => {
  if (!expiryDate) return false
  return new Date(expiryDate) > new Date()
}

const formatDate = (date: string) => {
  return dayjs(date).format('YYYY-MM-DD')
}

const formatDateTime = (date: string) => {
  return dayjs(date).format('YYYY-MM-DD HH:mm')
}

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

const handleSearch = debounce(() => {
  pageNum.value = 1
  loadUserList()
}, 300)

const handleSizeChange = () => {
  pageNum.value = 1
  loadUserList()
}

const handleCurrentChange = (page: number) => {
  if (page >= 1 && page <= totalPages.value) {
    pageNum.value = page
    loadUserList()
  }
}

const handleRenewVip = (user: UserListItem) => {
  currentUser.value = user
  renewForm.days = 1  // 默认1个月
  renewDialogVisible.value = true
}

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
/* 所有样式在unified.css中 */
</style>

