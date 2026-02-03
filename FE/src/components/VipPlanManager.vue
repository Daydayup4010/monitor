<template>
  <div class="vip-plan-manager">
    <div class="manager-header">
      <h2>VIP定价管理</h2>
      <el-button type="primary" @click="showAddDialog">
        <el-icon><Plus /></el-icon>
        添加套餐
      </el-button>
    </div>

    <el-table :data="plans" v-loading="loading" stripe>
      <el-table-column prop="id" label="ID" width="80" />
      <el-table-column prop="months" label="时长（月）" width="120">
        <template #default="{ row }">
          <span class="months-text">{{ row.months }}个月</span>
        </template>
      </el-table-column>
      <el-table-column prop="price" label="价格" width="120">
        <template #default="{ row }">
          <span class="price-text">¥{{ row.price.toFixed(2) }}</span>
        </template>
      </el-table-column>
      <el-table-column label="月均价" width="120">
        <template #default="{ row }">
          <span class="avg-price">¥{{ (row.price / row.months).toFixed(2) }}/月</span>
        </template>
      </el-table-column>
      <el-table-column prop="enabled" label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="row.enabled ? 'success' : 'info'">
            {{ row.enabled ? '启用' : '禁用' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="created_at" label="创建时间" width="180">
        <template #default="{ row }">
          {{ formatTime(row.created_at) }}
        </template>
      </el-table-column>
      <el-table-column label="操作" width="180" fixed="right">
        <template #default="{ row }">
          <el-button link type="primary" @click="showEditDialog(row)">编辑</el-button>
          <el-button link type="danger" @click="handleDelete(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <!-- 添加/编辑对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="isEdit ? '编辑套餐' : '添加套餐'"
      width="450px"
      :close-on-click-modal="false"
    >
      <el-form :model="form" :rules="rules" ref="formRef" label-width="100px">
        <el-form-item label="时长（月）" prop="months">
          <el-input-number
            v-model="form.months"
            :min="1"
            :max="120"
            controls-position="right"
            style="width: 200px"
          />
        </el-form-item>
        <el-form-item label="价格" prop="price">
          <el-input-number
            v-model="form.price"
            :min="0.01"
            :precision="2"
            :step="10"
            controls-position="right"
            style="width: 200px"
          />
          <span class="price-unit">元</span>
        </el-form-item>
        <el-form-item label="状态" prop="enabled">
          <el-switch
            v-model="form.enabled"
            active-text="启用"
            inactive-text="禁用"
          />
        </el-form-item>
        <el-form-item label="月均价">
          <span class="calculated-price">
            ¥{{ form.months > 0 ? (form.price / form.months).toFixed(2) : '0.00' }}/月
          </span>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitting">
          {{ isEdit ? '保存' : '添加' }}
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, reactive } from 'vue'
import { Plus } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import api from '@/api'
import dayjs from 'dayjs'

interface VipPlan {
  id: number
  months: number
  price: number
  enabled: boolean
  created_at: string
  updated_at: string
}

const loading = ref(false)
const plans = ref<VipPlan[]>([])
const dialogVisible = ref(false)
const isEdit = ref(false)
const submitting = ref(false)
const formRef = ref<FormInstance>()

const form = reactive({
  id: 0,
  months: 1,
  price: 19.9,
  enabled: true
})

const rules: FormRules = {
  months: [
    { required: true, message: '请输入月数', trigger: 'blur' },
    { type: 'number', min: 1, message: '月数必须大于0', trigger: 'blur' }
  ],
  price: [
    { required: true, message: '请输入价格', trigger: 'blur' },
    { type: 'number', min: 0.01, message: '价格必须大于0', trigger: 'blur' }
  ]
}

// 加载套餐列表
const loadPlans = async () => {
  loading.value = true
  try {
    const res = await api.get('/admin/vip-plans')
    // 响应拦截器已经解包，res 直接是后端返回的数据
    if (res.code === 1) {
      plans.value = res.data || []
    }
  } catch (error) {
    console.error('加载套餐失败:', error)
    ElMessage.error('加载套餐失败')
  } finally {
    loading.value = false
  }
}

// 显示添加对话框
const showAddDialog = () => {
  isEdit.value = false
  form.id = 0
  form.months = 1
  form.price = 19.9
  form.enabled = true
  dialogVisible.value = true
}

// 显示编辑对话框
const showEditDialog = (row: VipPlan) => {
  isEdit.value = true
  form.id = row.id
  form.months = row.months
  form.price = row.price
  form.enabled = row.enabled
  dialogVisible.value = true
}

// 提交表单
const handleSubmit = async () => {
  const valid = await formRef.value?.validate()
  if (!valid) return

  submitting.value = true
  try {
    if (isEdit.value) {
      const res = await api.put('/admin/vip-plan', form)
      if (res.code === 1) {
        ElMessage.success('更新成功')
        dialogVisible.value = false
        loadPlans()
      } else {
        ElMessage.error(res.msg || '更新失败')
      }
    } else {
      const res = await api.post('/admin/vip-plan', form)
      if (res.code === 1) {
        ElMessage.success('添加成功')
        dialogVisible.value = false
        loadPlans()
      } else {
        ElMessage.error(res.msg || '添加失败')
      }
    }
  } catch (error) {
    console.error('操作失败:', error)
  } finally {
    submitting.value = false
  }
}

// 删除套餐
const handleDelete = async (row: VipPlan) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除 ${row.months}个月 的套餐吗？`,
      '删除确认',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )

    const res = await api.delete('/admin/vip-plan', { data: { id: row.id } })
    if (res.code === 1) {
      ElMessage.success('删除成功')
      loadPlans()
    } else {
      ElMessage.error(res.msg || '删除失败')
    }
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除失败:', error)
    }
  }
}

// 格式化时间
const formatTime = (time: string) => {
  return dayjs(time).format('YYYY-MM-DD HH:mm:ss')
}

onMounted(() => {
  loadPlans()
})
</script>

<style scoped>
.vip-plan-manager {
  background: white;
  border-radius: 12px;
  padding: 24px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
}

.manager-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.manager-header h2 {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
  color: #262626;
}

.months-text {
  font-weight: 500;
}

.price-text {
  color: #f56c6c;
  font-weight: 600;
}

.avg-price {
  color: #909399;
  font-size: 13px;
}

.price-unit {
  margin-left: 8px;
  color: #909399;
}

.calculated-price {
  color: #f56c6c;
  font-weight: 500;
}
</style>
