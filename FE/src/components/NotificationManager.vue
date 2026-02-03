<template>
  <div class="notification-manager">
    <!-- 操作栏 -->
    <div class="action-bar">
      <el-button type="primary" @click="showCreateDialog">
        <el-icon><Plus /></el-icon>
        发布通知
      </el-button>
    </div>
    
    <!-- 通知列表 -->
    <el-table 
      :data="notifications" 
      v-loading="loading"
      stripe
      class="notification-table"
    >
      <el-table-column prop="title" label="标题" min-width="200">
        <template #default="{ row }">
          <div class="notification-title-cell">
            {{ row.title }}
          </div>
        </template>
      </el-table-column>
      
      <el-table-column prop="content" label="内容" min-width="300">
        <template #default="{ row }">
          <div class="notification-content-cell">
            {{ row.content }}
          </div>
        </template>
      </el-table-column>
      
      <el-table-column label="图片" width="80" align="center">
        <template #default="{ row }">
          <el-image 
            v-if="row.image_url"
            :src="row.image_url" 
            fit="cover"
            :preview-src-list="[row.image_url]"
            style="width: 40px; height: 40px; border-radius: 4px;"
          >
            <template #error>
              <div class="image-placeholder">-</div>
            </template>
          </el-image>
          <span v-else class="no-image">-</span>
        </template>
      </el-table-column>
      
      <el-table-column prop="created_at" label="发布时间" width="180">
        <template #default="{ row }">
          {{ formatTime(row.created_at) }}
        </template>
      </el-table-column>
      
      <el-table-column label="操作" width="120" fixed="right">
        <template #default="{ row }">
          <el-button 
            link 
            type="primary" 
            @click="viewDetail(row)"
          >
            查看
          </el-button>
          <el-popconfirm
            title="确定要删除这条通知吗？"
            confirm-button-text="确定"
            cancel-button-text="取消"
            @confirm="handleDelete(row)"
          >
            <template #reference>
              <el-button link type="danger">删除</el-button>
            </template>
          </el-popconfirm>
        </template>
      </el-table-column>
    </el-table>
    
    <!-- 分页 -->
    <div class="pagination-wrapper">
      <el-pagination
        v-model:current-page="pageNum"
        v-model:page-size="pageSize"
        :total="total"
        :page-sizes="[10, 20, 50]"
        layout="total, sizes, prev, pager, next"
        @size-change="loadNotifications"
        @current-change="loadNotifications"
      />
    </div>
    
    <!-- 创建通知弹窗 -->
    <el-dialog
      v-model="createDialogVisible"
      title="发布通知"
      width="600px"
      :close-on-click-modal="false"
    >
      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-width="80px"
      >
        <el-form-item label="标题" prop="title">
          <el-input 
            v-model="form.title" 
            placeholder="请输入通知标题"
            maxlength="100"
            show-word-limit
          />
        </el-form-item>
        
        <el-form-item label="内容" prop="content">
          <el-input
            v-model="form.content"
            type="textarea"
            placeholder="请输入通知内容"
            :rows="6"
            maxlength="2000"
            show-word-limit
          />
        </el-form-item>
        
        <el-form-item label="图片">
          <div class="image-upload-wrapper">
            <el-upload
              class="image-uploader"
              :show-file-list="false"
              :before-upload="beforeImageUpload"
              :http-request="handleImageUpload"
              accept=".jpg,.jpeg,.png,.gif,.webp"
            >
              <div v-if="form.image_url" class="uploaded-image">
                <el-image 
                  :src="form.image_url" 
                  fit="contain"
                  style="max-width: 200px; max-height: 150px;"
                />
                <div class="image-actions">
                  <el-button 
                    type="danger" 
                    size="small" 
                    circle 
                    @click.stop="removeImage"
                  >
                    <el-icon><Delete /></el-icon>
                  </el-button>
                </div>
              </div>
              <div v-else class="upload-placeholder" v-loading="uploading">
                <el-icon class="upload-icon"><Plus /></el-icon>
                <span>点击上传图片</span>
                <span class="upload-tip">支持 jpg、png、gif、webp，最大 5MB</span>
              </div>
            </el-upload>
          </div>
        </el-form-item>
      </el-form>
      
      <template #footer>
        <el-button @click="createDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleCreate" :loading="creating">
          发布
        </el-button>
      </template>
    </el-dialog>
    
    <!-- 查看详情弹窗 -->
    <el-dialog
      v-model="detailDialogVisible"
      :title="currentNotification?.title"
      width="600px"
    >
      <div class="detail-content">
        {{ currentNotification?.content }}
      </div>
      <div class="detail-image" v-if="currentNotification?.image_url">
        <el-image 
          :src="currentNotification.image_url" 
          fit="contain"
          :preview-src-list="[currentNotification.image_url]"
          style="max-width: 100%; max-height: 300px; margin-top: 16px;"
        />
      </div>
      <div class="detail-time">
        发布时间：{{ currentNotification ? formatTime(currentNotification.created_at) : '' }}
      </div>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { Plus, Delete } from '@element-plus/icons-vue'
import { notificationAdminApi } from '@/api/notification'
import type { NotificationItem } from '@/api/notification'
import { showMessage } from '@/utils/message'
import type { FormInstance, FormRules } from 'element-plus'
import dayjs from 'dayjs'

const loading = ref(false)
const creating = ref(false)
const uploading = ref(false)
const notifications = ref<NotificationItem[]>([])
const pageNum = ref(1)
const pageSize = ref(10)
const total = ref(0)

const createDialogVisible = ref(false)
const detailDialogVisible = ref(false)
const currentNotification = ref<NotificationItem | null>(null)

const formRef = ref<FormInstance>()
const form = reactive({
  title: '',
  content: '',
  image_url: ''
})

const rules: FormRules = {
  title: [
    { required: true, message: '请输入通知标题', trigger: 'blur' },
    { min: 1, max: 100, message: '标题长度不能超过100个字符', trigger: 'blur' }
  ],
  content: [
    { required: true, message: '请输入通知内容', trigger: 'blur' },
    { min: 1, max: 2000, message: '内容长度不能超过2000个字符', trigger: 'blur' }
  ]
}

// 加载通知列表
const loadNotifications = async () => {
  loading.value = true
  try {
    const res = await notificationAdminApi.getAll({
      page_num: pageNum.value,
      page_size: pageSize.value
    })
    if (res.code === 1) {
      notifications.value = res.data || []
      total.value = res.total || 0
    }
  } catch (error) {
    console.error('加载通知列表失败:', error)
  } finally {
    loading.value = false
  }
}

// 显示创建弹窗
const showCreateDialog = () => {
  form.title = ''
  form.content = ''
  form.image_url = ''
  createDialogVisible.value = true
}

// 创建通知
const handleCreate = async () => {
  if (!formRef.value) return
  
  await formRef.value.validate(async (valid) => {
    if (!valid) return
    
    creating.value = true
    try {
      const res = await notificationAdminApi.create({
        title: form.title,
        content: form.content,
        image_url: form.image_url || undefined
      })
      if (res.code === 1) {
        showMessage.success('发布成功')
        createDialogVisible.value = false
        loadNotifications()
      }
    } catch (error) {
      console.error('发布通知失败:', error)
    } finally {
      creating.value = false
    }
  })
}

// 查看详情
const viewDetail = (row: NotificationItem) => {
  currentNotification.value = row
  detailDialogVisible.value = true
}

// 删除通知
const handleDelete = async (row: NotificationItem) => {
  try {
    const res = await notificationAdminApi.delete(row.id)
    if (res.code === 1) {
      showMessage.success('删除成功')
      loadNotifications()
    }
  } catch (error) {
    console.error('删除通知失败:', error)
  }
}

// 格式化时间
const formatTime = (time: string) => {
  return dayjs(time).format('YYYY-MM-DD HH:mm:ss')
}

// 上传前检查
const beforeImageUpload = (file: File) => {
  const isValidType = ['image/jpeg', 'image/png', 'image/gif', 'image/webp'].includes(file.type)
  const isValidSize = file.size / 1024 / 1024 < 5

  if (!isValidType) {
    showMessage.error('只支持 jpg、png、gif、webp 格式的图片')
    return false
  }
  if (!isValidSize) {
    showMessage.error('图片大小不能超过 5MB')
    return false
  }
  return true
}

// 上传图片
const handleImageUpload = async (options: { file: File }) => {
  uploading.value = true
  try {
    const res = await notificationAdminApi.uploadImage(options.file)
    if (res.code === 1 && res.data?.url) {
      form.image_url = res.data.url
      showMessage.success('图片上传成功')
    }
  } catch (error) {
    console.error('上传图片失败:', error)
    showMessage.error('图片上传失败')
  } finally {
    uploading.value = false
  }
}

// 移除图片
const removeImage = () => {
  form.image_url = ''
}

onMounted(() => {
  loadNotifications()
})
</script>

<style scoped>
.notification-manager {
  padding: 20px;
}

.action-bar {
  margin-bottom: 20px;
}

.notification-table {
  border-radius: 8px;
  overflow: hidden;
}

.notification-title-cell {
  font-weight: 500;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.notification-content-cell {
  color: #666;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.pagination-wrapper {
  display: flex;
  justify-content: flex-end;
  margin-top: 20px;
}

.detail-content {
  font-size: 15px;
  line-height: 1.8;
  color: #333;
  white-space: pre-wrap;
  max-height: 400px;
  overflow-y: auto;
}

.detail-time {
  font-size: 13px;
  color: #999;
  margin-top: 20px;
  padding-top: 16px;
  border-top: 1px solid #f0f0f0;
}

.no-image {
  color: #ccc;
}

.image-placeholder {
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #f5f5f5;
  color: #ccc;
  border-radius: 4px;
}

.image-error {
  font-size: 12px;
  color: #f56c6c;
  padding: 10px;
  background: #fef0f0;
  border-radius: 4px;
}

.image-preview {
  margin-top: 8px;
}

.image-upload-wrapper {
  width: 100%;
}

.image-uploader {
  width: 100%;
}

.upload-placeholder {
  width: 200px;
  height: 150px;
  border: 2px dashed #dcdfe6;
  border-radius: 8px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.3s;
  background: #fafafa;
}

.upload-placeholder:hover {
  border-color: #409eff;
  background: #f0f7ff;
}

.upload-icon {
  font-size: 32px;
  color: #8c939d;
  margin-bottom: 8px;
}

.upload-placeholder span {
  font-size: 14px;
  color: #606266;
}

.upload-tip {
  font-size: 12px !important;
  color: #909399 !important;
  margin-top: 4px;
}

.uploaded-image {
  position: relative;
  display: inline-block;
}

.image-actions {
  position: absolute;
  top: 4px;
  right: 4px;
}
</style>
