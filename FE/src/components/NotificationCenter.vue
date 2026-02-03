<template>
  <div class="notification-center">
    <el-popover
      ref="popoverRef"
      placement="bottom-end"
      :width="360"
      trigger="click"
      :show-arrow="false"
      popper-class="notification-popover"
    >
      <template #reference>
        <div class="notification-trigger" @click="handleClick">
          <el-icon class="notification-icon"><Bell /></el-icon>
          <span v-if="unreadCount > 0" class="notification-badge">
            {{ unreadCount > 99 ? '99+' : unreadCount }}
          </span>
        </div>
      </template>
      
      <div class="notification-panel">
        <div class="notification-header">
          <span class="notification-title">通知中心</span>
          <el-button 
            v-if="unreadCount > 0" 
            link 
            type="primary" 
            @click="markAllRead"
            :loading="markingAllRead"
          >
            全部已读
          </el-button>
        </div>
        
        <div class="notification-list" v-loading="loading">
          <template v-if="notifications.length > 0">
            <div 
              v-for="item in notifications" 
              :key="item.id"
              class="notification-item"
              :class="{ unread: !item.is_read }"
              @click="handleItemClick(item)"
            >
              <div class="notification-dot" v-if="!item.is_read"></div>
              <div class="notification-content">
                <div class="notification-item-title">
                  <span class="notification-tag">【公告】</span>{{ item.title }}
                </div>
                <div class="notification-item-time">{{ formatTime(item.created_at) }}</div>
              </div>
            </div>
          </template>
          <el-empty v-else description="暂无通知" :image-size="80" />
        </div>
        
        <div class="notification-footer" v-if="notifications.length > 0">
          <el-button link type="primary" @click="goToNotificationCenter">查看全部</el-button>
        </div>
      </div>
    </el-popover>
    
    <!-- 全部通知弹窗 -->
    <el-dialog
      v-model="allNotificationsVisible"
      title="全部通知"
      width="600px"
      class="all-notifications-dialog"
      append-to-body
      destroy-on-close
    >
      <div class="all-notifications-list" v-loading="loadingAll">
        <template v-if="allNotifications.length > 0">
          <div 
            v-for="item in allNotifications" 
            :key="item.id"
            class="notification-item"
            :class="{ unread: !item.is_read }"
            @click="handleItemClick(item)"
          >
            <div class="notification-dot" v-if="!item.is_read"></div>
            <div class="notification-content">
              <div class="notification-item-title">{{ item.title }}</div>
              <div class="notification-item-desc">{{ item.content }}</div>
              <div class="notification-item-time">{{ formatTime(item.created_at) }}</div>
            </div>
          </div>
        </template>
        <el-empty v-else description="暂无通知" :image-size="100" />
      </div>
      
      <!-- 分页 -->
      <div class="pagination-wrapper" v-if="allTotal > pageSize">
        <el-pagination
          v-model:current-page="allPageNum"
          :page-size="pageSize"
          :total="allTotal"
          layout="prev, pager, next"
          @current-change="loadAllNotifications"
        />
      </div>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { Bell } from '@element-plus/icons-vue'
import { notificationApi } from '@/api/notification'
import { showMessage } from '@/utils/message'
import dayjs from 'dayjs'
import relativeTime from 'dayjs/plugin/relativeTime'
import 'dayjs/locale/zh-cn'

dayjs.extend(relativeTime)
dayjs.locale('zh-cn')

const router = useRouter()

interface NotificationItem {
  id: string
  title: string
  content: string
  image_url?: string
  is_read: boolean
  created_at: string
}

const popoverRef = ref()
const loading = ref(false)
const loadingAll = ref(false)
const markingAllRead = ref(false)
const unreadCount = ref(0)
const notifications = ref<NotificationItem[]>([])
const allNotifications = ref<NotificationItem[]>([])
const allNotificationsVisible = ref(false)

const pageSize = 10
const allPageNum = ref(1)
const allTotal = ref(0)

// 轮询定时器
let pollTimer: ReturnType<typeof setInterval> | null = null

// 获取未读数量
const fetchUnreadCount = async () => {
  try {
    const res = await notificationApi.getUnreadCount()
    if (res.code === 1) {
      unreadCount.value = res.data?.unread_count || 0
    }
  } catch (error) {
    console.error('获取未读数量失败:', error)
  }
}

// 获取通知列表
const fetchNotifications = async () => {
  loading.value = true
  try {
    const res = await notificationApi.getList({ page_num: 1, page_size: 5 })
    if (res.code === 1) {
      notifications.value = res.data || []
    }
  } catch (error) {
    console.error('获取通知列表失败:', error)
  } finally {
    loading.value = false
  }
}

// 加载全部通知
const loadAllNotifications = async () => {
  loadingAll.value = true
  try {
    const res = await notificationApi.getList({ 
      page_num: allPageNum.value, 
      page_size: pageSize 
    })
    if (res.code === 1) {
      allNotifications.value = res.data || []
      allTotal.value = res.total || 0
    }
  } catch (error) {
    console.error('获取全部通知失败:', error)
  } finally {
    loadingAll.value = false
  }
}

// 点击触发器
const handleClick = () => {
  fetchNotifications()
}

// 点击通知项 - 跳转到消息中心页面并打开该通知
const handleItemClick = async (item: NotificationItem) => {
  popoverRef.value?.hide()
  
  // 标记为已读
  if (!item.is_read) {
    try {
      await notificationApi.markAsRead({ notification_id: item.id })
      item.is_read = true
      unreadCount.value = Math.max(0, unreadCount.value - 1)
    } catch (error) {
      console.error('标记已读失败:', error)
    }
  }
  
  // 跳转到消息中心页面，传递通知ID
  router.push({ path: '/app/notifications', query: { id: item.id } })
}

// 全部已读
const markAllRead = async () => {
  markingAllRead.value = true
  try {
    const res = await notificationApi.markAllAsRead()
    if (res.code === 1) {
      notifications.value.forEach(item => item.is_read = true)
      allNotifications.value.forEach(item => item.is_read = true)
      unreadCount.value = 0
      showMessage.success('已全部标记为已读')
    }
  } catch (error) {
    console.error('标记全部已读失败:', error)
  } finally {
    markingAllRead.value = false
  }
}

// 查看全部
const viewAll = () => {
  popoverRef.value?.hide()
  allPageNum.value = 1
  allNotificationsVisible.value = true
  loadAllNotifications()
}

// 跳转到消息中心页面
const goToNotificationCenter = () => {
  popoverRef.value?.hide()
  router.push('/app/notifications')
}

// 格式化时间
const formatTime = (time: string) => {
  const date = dayjs(time)
  const now = dayjs()
  
  if (now.diff(date, 'day') < 1) {
    return date.fromNow()
  } else if (now.diff(date, 'day') < 7) {
    return date.format('MM-DD HH:mm')
  } else {
    return date.format('YYYY-MM-DD HH:mm')
  }
}

// 开始轮询
const startPolling = () => {
  // 每60秒轮询一次未读数量
  pollTimer = setInterval(fetchUnreadCount, 60000)
}

// 停止轮询
const stopPolling = () => {
  if (pollTimer) {
    clearInterval(pollTimer)
    pollTimer = null
  }
}

onMounted(() => {
  fetchUnreadCount()
  startPolling()
})

onUnmounted(() => {
  stopPolling()
})
</script>

<style scoped>
.notification-center {
  display: flex;
  align-items: center;
}

.notification-trigger {
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 40px;
  height: 40px;
  border-radius: 50%;
  cursor: pointer;
  transition: all 0.3s;
}

.notification-trigger:hover {
  background: rgba(24, 144, 255, 0.1);
}

.notification-icon {
  font-size: 22px;
  color: #666;
}

.notification-trigger:hover .notification-icon {
  color: #1890ff;
}

.notification-badge {
  position: absolute;
  top: 2px;
  right: 2px;
  min-width: 18px;
  height: 18px;
  padding: 0 5px;
  font-size: 11px;
  font-weight: 600;
  color: #fff;
  background: linear-gradient(135deg, #ff6b6b, #ee5a5a);
  border-radius: 9px;
  display: flex;
  align-items: center;
  justify-content: center;
  line-height: 1;
  box-shadow: 0 2px 6px rgba(255, 107, 107, 0.4);
  animation: badge-pulse 2s infinite;
}

@keyframes badge-pulse {
  0%, 100% {
    transform: scale(1);
  }
  50% {
    transform: scale(1.1);
  }
}

.notification-panel {
  margin: -12px;
}

.notification-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 20px;
  border-bottom: 1px solid #f0f0f0;
}

.notification-title {
  font-size: 16px;
  font-weight: 600;
  color: #333;
}

.notification-list {
  max-height: 360px;
  overflow-y: auto;
}

.notification-item {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  padding: 16px 20px;
  cursor: pointer;
  transition: background 0.2s;
  border-bottom: 1px solid #f5f5f5;
}

.notification-item:last-child {
  border-bottom: none;
}

.notification-item:hover {
  background: #f9fafb;
}

.notification-item.unread {
  background: #f0f7ff;
}

.notification-item.unread:hover {
  background: #e6f4ff;
}

.notification-dot {
  width: 8px;
  height: 8px;
  background: #1890ff;
  border-radius: 50%;
  margin-top: 6px;
  flex-shrink: 0;
}

.notification-content {
  flex: 1;
  min-width: 0;
}

.notification-item-title {
  font-size: 14px;
  font-weight: 500;
  color: #333;
  margin-bottom: 4px;
  line-height: 1.5;
}

.notification-tag {
  color: #667eea;
  font-weight: 600;
}

.notification-item-time {
  font-size: 12px;
  color: #999;
  margin-top: 6px;
}

.notification-footer {
  padding: 12px 20px;
  text-align: center;
  border-top: 1px solid #f0f0f0;
}

/* 全部通知弹窗 */
.all-notifications-list {
  max-height: 500px;
  overflow-y: auto;
}

.pagination-wrapper {
  display: flex;
  justify-content: center;
  margin-top: 16px;
  padding-top: 16px;
  border-top: 1px solid #f0f0f0;
}
</style>

<style>
/* 全局样式覆盖 */
.notification-popover {
  padding: 0 !important;
  border-radius: 12px !important;
  box-shadow: 0 6px 24px rgba(0, 0, 0, 0.12) !important;
}

.all-notifications-dialog .el-dialog__body {
  padding: 0 20px 20px;
}
</style>
