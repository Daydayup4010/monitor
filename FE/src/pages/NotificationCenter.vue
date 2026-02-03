<template>
  <div class="notification-page">
    <!-- 页面标题 -->
    <div class="page-header">
      <h1 class="page-title">消息中心</h1>
    </div>

    <!-- 主体内容 -->
    <div class="notification-wrapper">
      <!-- 顶部操作栏 -->
      <div class="content-header">
        <div class="header-tabs">
          <span 
            class="tab-item" 
            :class="{ active: activeTab === 'all' }"
            @click="activeTab = 'all'"
          >全部</span>
          <span 
            class="tab-item" 
            :class="{ active: activeTab === 'unread' }"
            @click="activeTab = 'unread'"
          >未读 <span class="tab-badge" v-if="unreadCount > 0">{{ unreadCount }}</span></span>
        </div>
        <div class="header-actions">
          <span 
            class="action-btn" 
            :class="{ disabled: unreadCount === 0 }"
            @click="markAllRead"
          >全部已读</span>
        </div>
      </div>

      <!-- 通知列表 -->
      <div class="notifications-list" v-loading="loading">
        <template v-if="filteredNotifications.length > 0">
          <div 
            v-for="item in filteredNotifications" 
            :key="item.id"
            class="list-item"
            :class="{ unread: !item.is_read }"
            @click="viewDetail(item)"
          >
            <div class="item-dot" v-if="!item.is_read"></div>
            <div class="item-body">
              <div class="item-header">
                <span class="item-tag">【公告】</span>
                <span class="item-title">{{ item.title }}</span>
              </div>
              <div class="item-desc">{{ item.content }}</div>
              <div class="item-time" :title="formatFullTime(item.created_at)">
                {{ formatTime(item.created_at) }}
              </div>
            </div>
          </div>
        </template>

        <el-empty v-else description="暂无消息" :image-size="100" />

        <!-- 加载更多 -->
        <div class="load-more" v-if="filteredNotifications.length > 0">
          <span v-if="hasMore" class="load-btn" @click="loadMore">加载更多</span>
          <span v-else class="load-end">到底了……</span>
        </div>
      </div>
    </div>

    <!-- 通知详情弹窗 -->
    <el-dialog
      v-model="detailVisible"
      :title="currentNotification?.title"
      width="560px"
      :close-on-click-modal="false"
    >
      <div class="detail-content">{{ currentNotification?.content }}</div>
      <div class="detail-image" v-if="currentNotification?.image_url">
        <el-image 
          :src="currentNotification.image_url" 
          fit="contain"
          :preview-src-list="[currentNotification.image_url]"
        />
      </div>
      <div class="detail-time">{{ currentNotification ? formatFullTime(currentNotification.created_at) : '' }}</div>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { useRoute } from 'vue-router'
import { notificationApi } from '@/api/notification'
import { showMessage } from '@/utils/message'
import dayjs from 'dayjs'
import relativeTime from 'dayjs/plugin/relativeTime'
import 'dayjs/locale/zh-cn'

dayjs.extend(relativeTime)
dayjs.locale('zh-cn')

const route = useRoute()

interface NotificationItem {
  id: string
  title: string
  content: string
  image_url?: string
  is_read: boolean
  created_at: string
}

const loading = ref(false)
const notifications = ref<NotificationItem[]>([])
const activeTab = ref('all')
const pageNum = ref(1)
const pageSize = 20
const total = ref(0)
const unreadCount = ref(0)
const hasMore = ref(true)

const detailVisible = ref(false)
const currentNotification = ref<NotificationItem | null>(null)

// 根据标签筛选
const filteredNotifications = computed(() => {
  if (activeTab.value === 'unread') {
    return notifications.value.filter(n => !n.is_read)
  }
  return notifications.value
})

// 加载通知列表
const loadNotifications = async (reset = false) => {
  if (reset) {
    pageNum.value = 1
    notifications.value = []
  }
  
  loading.value = true
  try {
    const res = await notificationApi.getList({
      page_num: pageNum.value,
      page_size: pageSize
    })
    if (res.code === 1) {
      if (reset) {
        notifications.value = res.data || []
      } else {
        notifications.value.push(...(res.data || []))
      }
      total.value = res.total || 0
      hasMore.value = notifications.value.length < total.value
    }
  } catch (error) {
    console.error('加载通知失败:', error)
  } finally {
    loading.value = false
  }
}

// 加载更多
const loadMore = () => {
  if (!hasMore.value || loading.value) return
  pageNum.value++
  loadNotifications()
}

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

// 查看详情
const viewDetail = async (item: NotificationItem) => {
  currentNotification.value = item
  detailVisible.value = true
  
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
}

// 全部已读
const markAllRead = async () => {
  if (unreadCount.value === 0) return
  
  try {
    const res = await notificationApi.markAllAsRead()
    if (res.code === 1) {
      notifications.value.forEach(item => item.is_read = true)
      unreadCount.value = 0
      showMessage.success('已全部标记为已读')
    }
  } catch (error) {
    console.error('标记全部已读失败:', error)
  }
}

// 格式化相对时间
const formatTime = (time: string) => {
  const date = dayjs(time)
  const now = dayjs()
  const diffDays = now.diff(date, 'day')
  
  if (diffDays < 1) {
    return date.fromNow()
  } else if (diffDays < 365) {
    return `${diffDays} 天前`
  } else {
    const years = Math.floor(diffDays / 365)
    return `${years} 年前`
  }
}

// 格式化完整时间
const formatFullTime = (time: string) => {
  return dayjs(time).format('YYYY-MM-DD HH:mm:ss')
}

// 根据 ID 查找并打开通知详情
const openNotificationById = async (id: string) => {
  let item = notifications.value.find(n => n.id === id)
  
  if (item) {
    viewDetail(item)
  } else {
    try {
      const res = await notificationApi.getList({ page_num: 1, page_size: 100 })
      if (res.code === 1 && res.data) {
        item = res.data.find((n: NotificationItem) => n.id === id)
        if (item) {
          viewDetail(item)
        }
      }
    } catch (error) {
      console.error('获取通知详情失败:', error)
    }
  }
}

onMounted(async () => {
  await loadNotifications(true)
  fetchUnreadCount()
  
  const notificationId = route.query.id as string
  if (notificationId) {
    openNotificationById(notificationId)
  }
})

watch(() => route.query.id, (newId) => {
  if (newId) {
    openNotificationById(newId as string)
  }
})
</script>

<style scoped>
.notification-page {
  max-width: 800px;
  margin: 0 auto;
  padding: 24px;
}

.page-header {
  margin-bottom: 20px;
}

.page-title {
  font-size: 22px;
  font-weight: 600;
  color: #1a1a1a;
  margin: 0;
}

.notification-wrapper {
  background: #fff;
  border-radius: 8px;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.08);
}

.content-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 20px;
  border-bottom: 1px solid #f0f0f0;
}

.header-tabs {
  display: flex;
  gap: 20px;
}

.tab-item {
  font-size: 14px;
  color: #666;
  cursor: pointer;
  padding-bottom: 4px;
  transition: all 0.2s;
}

.tab-item:hover {
  color: #1890ff;
}

.tab-item.active {
  color: #1a1a1a;
  font-weight: 500;
}

.tab-badge {
  display: inline-block;
  background: #ff4d4f;
  color: #fff;
  font-size: 12px;
  padding: 0 6px;
  border-radius: 10px;
  margin-left: 4px;
}

.action-btn {
  font-size: 14px;
  color: #1890ff;
  cursor: pointer;
}

.action-btn:hover {
  color: #40a9ff;
}

.action-btn.disabled {
  color: #ccc;
  cursor: not-allowed;
}

.notifications-list {
  min-height: 300px;
}

.list-item {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  padding: 16px 20px;
  cursor: pointer;
  transition: background 0.2s;
  border-bottom: 1px solid #f5f5f5;
}

.list-item:last-child {
  border-bottom: none;
}

.list-item:hover {
  background: #fafafa;
}

.list-item.unread {
  background: #f6ffed;
}

.list-item.unread:hover {
  background: #e6ffe6;
}

.item-dot {
  width: 8px;
  height: 8px;
  background: #52c41a;
  border-radius: 50%;
  margin-top: 6px;
  flex-shrink: 0;
}

.item-body {
  flex: 1;
  min-width: 0;
}

.item-header {
  margin-bottom: 6px;
}

.item-tag {
  color: #1890ff;
  font-weight: 500;
  font-size: 14px;
}

.item-title {
  font-size: 14px;
  color: #1a1a1a;
  font-weight: 500;
}

.item-desc {
  font-size: 13px;
  color: #666;
  line-height: 1.6;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  margin-bottom: 8px;
}

.item-time {
  font-size: 12px;
  color: #999;
}

.load-more {
  text-align: center;
  padding: 16px;
}

.load-btn {
  font-size: 13px;
  color: #1890ff;
  cursor: pointer;
}

.load-btn:hover {
  color: #40a9ff;
}

.load-end {
  font-size: 13px;
  color: #999;
}

/* 详情弹窗 */
.detail-content {
  font-size: 15px;
  line-height: 1.8;
  color: #333;
  white-space: pre-wrap;
}

.detail-image {
  margin-top: 16px;
  text-align: center;
}

.detail-image :deep(.el-image) {
  max-width: 100%;
  max-height: 400px;
  border-radius: 4px;
}

.detail-image :deep(.el-image img) {
  max-width: 100%;
  max-height: 400px;
  object-fit: contain;
}

.detail-time {
  font-size: 13px;
  color: #999;
  margin-top: 16px;
  padding-top: 12px;
  border-top: 1px solid #f0f0f0;
}
</style>
