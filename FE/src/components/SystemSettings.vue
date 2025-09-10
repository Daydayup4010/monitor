<template>
  <div class="system-settings">
    <div v-if="settingsStore.loading" class="loading-container">
      <el-loading
        element-loading-text="加载设置中..."
        element-loading-background="transparent"
      />
    </div>
    
    <el-form
      v-else
      ref="formRef"
      :model="form"
      :rules="rules"
      label-width="140px"
      class="settings-form"
      @submit.prevent="submitForm"
    >
      <el-row :gutter="24">
        <el-col :span="12">
          <el-form-item label="最小销售量" prop="min_sell_num">
            <el-input-number
              v-model="form.min_sell_num"
              :min="0"
              :max="10000"
              :step="1"
              controls-position="right"
              class="full-width"
            />
            <div class="form-help">
              设置商品的最小销售量，低于此值的商品将被过滤
            </div>
          </el-form-item>
        </el-col>
        
        <el-col :span="12">
          <el-form-item label="最小价格差异" prop="min_diff">
            <el-input-number
              v-model="form.min_diff"
              :min="0"
              :max="1000"
              :step="0.1"
              :precision="2"
              controls-position="right"
              class="full-width"
            />
            <div class="form-help">
              设置UU价格与Buff价格的最小差异，单位：元
            </div>
          </el-form-item>
        </el-col>
      </el-row>
      
      <el-row :gutter="24">
        <el-col :span="12">
          <el-form-item label="最低销售价格" prop="min_sell_price">
            <el-input-number
              v-model="form.min_sell_price"
              :min="0"
              :max="100000"
              :step="1"
              :precision="2"
              controls-position="right"
              class="full-width"
            />
            <div class="form-help">
              设置商品的最低销售价格，低于此价格的商品将被过滤
            </div>
          </el-form-item>
        </el-col>
        
        <el-col :span="12">
          <el-form-item label="最高销售价格" prop="max_sell_price">
            <el-input-number
              v-model="form.max_sell_price"
              :min="0"
              :max="100000"
              :step="1"
              :precision="2"
              controls-position="right"
              class="full-width"
            />
            <div class="form-help">
              设置商品的最高销售价格，高于此价格的商品将被过滤
            </div>
          </el-form-item>
        </el-col>
      </el-row>
      
      <el-form-item class="form-actions">
        <el-button
          type="primary"
          size="large"
          :loading="settingsStore.loading"
          @click="submitForm"
          class="save-btn"
        >
          <el-icon><Check /></el-icon>
          保存设置
        </el-button>
        <el-button size="large" @click="resetForm">
          <el-icon><RefreshLeft /></el-icon>
          重置
        </el-button>
        <el-button size="large" @click="previewEffect">
          <el-icon><View /></el-icon>
          预览效果
        </el-button>
      </el-form-item>
    </el-form>
    
    <!-- 设置说明 -->
    <el-card class="help-card">
      <template #header>
        <div class="card-header">
          <el-icon size="20"><QuestionFilled /></el-icon>
          <span>设置说明</span>
        </div>
      </template>
      
      <el-row :gutter="24">
        <el-col :span="12">
          <div class="help-section">
            <h4>参数说明：</h4>
            <ul>
              <li><strong>最小销售量：</strong>筛选销售量大于等于此值的商品</li>
              <li><strong>最小价格差异：</strong>筛选UU价格与Buff价格差异大于此值的商品</li>
              <li><strong>价格区间：</strong>筛选在指定价格区间内的商品</li>
            </ul>
          </div>
        </el-col>
        
        <el-col :span="12">
          <div class="help-section">
            <h4>建议配置：</h4>
            <ul>
              <li>新手用户：最小销售量 ≥ 10，价格差异 ≥ 5元</li>
              <li>进阶用户：最小销售量 ≥ 5，价格差异 ≥ 2元</li>
              <li>专业用户：根据市场情况灵活调整参数</li>
            </ul>
          </div>
        </el-col>
      </el-row>
      
      <el-alert
        title="重要提醒"
        type="warning"
        description="修改设置后，新的筛选条件将在下次数据刷新时生效。"
        show-icon
        :closable="false"
        class="alert-section"
      />
    </el-card>
    
    <!-- 当前设置预览 -->
    <el-card v-if="!settingsStore.loading" class="preview-card">
      <template #header>
        <div class="card-header">
          <el-icon size="20"><DataAnalysis /></el-icon>
          <span>当前设置预览</span>
        </div>
      </template>
      
      <div class="preview-content">
        <el-descriptions :column="2" border>
          <el-descriptions-item label="最小销售量">
            <el-tag>{{ form.min_sell_num }} 件</el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="最小价格差异">
            <el-tag type="success">¥{{ formatPrice(form.min_diff) }}</el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="最低销售价格">
            <el-tag type="warning">¥{{ formatPrice(form.min_sell_price) }}</el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="最高销售价格">
            <el-tag type="danger">¥{{ formatPrice(form.max_sell_price) }}</el-tag>
          </el-descriptions-item>
        </el-descriptions>
        
        <div class="filter-description">
          <h4>筛选条件：</h4>
          <p>
            显示销售量 ≥ <strong>{{ form.min_sell_num }}</strong> 件，
            价格差异 ≥ <strong>¥{{ formatPrice(form.min_diff) }}</strong>，
            价格区间在 <strong>¥{{ formatPrice(form.min_sell_price) }}</strong> - 
            <strong>¥{{ formatPrice(form.max_sell_price) }}</strong> 之间的商品
          </p>
        </div>
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted, watch } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useSettingsStore } from '@/stores/settings'
import { formatPrice } from '@/utils'
import type { Settings } from '@/types'

const settingsStore = useSettingsStore()

// 表单引用
const formRef = ref()

// 表单数据
const form = reactive<Settings>({
  min_sell_num: 0,
  min_diff: 0,
  max_sell_price: 10000,
  min_sell_price: 0,
})

// 表单验证规则
const rules = {
  min_sell_num: [
    { required: true, message: '请输入最小销售量', trigger: 'blur' },
    { type: 'number', min: 0, message: '最小销售量不能小于0', trigger: 'blur' }
  ],
  min_diff: [
    { required: true, message: '请输入最小价格差异', trigger: 'blur' },
    { type: 'number', min: 0, message: '最小价格差异不能小于0', trigger: 'blur' }
  ],
  min_sell_price: [
    { required: true, message: '请输入最低销售价格', trigger: 'blur' },
    { type: 'number', min: 0, message: '最低销售价格不能小于0', trigger: 'blur' }
  ],
  max_sell_price: [
    { required: true, message: '请输入最高销售价格', trigger: 'blur' },
    { type: 'number', min: 0, message: '最高销售价格不能小于0', trigger: 'blur' },
    {
      validator: (rule: any, value: number, callback: Function) => {
        if (value <= form.min_sell_price) {
          callback(new Error('最高销售价格必须大于最低销售价格'))
        } else {
          callback()
        }
      },
      trigger: 'blur'
    }
  ]
}

// 加载设置
const loadSettings = async () => {
  await settingsStore.getSettings()
  updateForm()
}

// 更新表单数据
const updateForm = () => {
  Object.assign(form, settingsStore.settings)
}

// 提交表单
const submitForm = async () => {
  if (!formRef.value) return
  
  try {
    await formRef.value.validate()
    await settingsStore.updateSettings(form)
  } catch (error) {
    console.error('Settings validation failed:', error)
  }
}

// 重置表单
const resetForm = () => {
  updateForm()
  formRef.value?.clearValidate()
}

// 预览效果
const previewEffect = () => {
  const { min_sell_num, min_diff, min_sell_price, max_sell_price } = form
  
  ElMessageBox.alert(
    `基于当前设置，系统将筛选出：
    • 销售量 ≥ ${min_sell_num} 件
    • 价格差异 ≥ ¥${formatPrice(min_diff)}
    • 价格区间：¥${formatPrice(min_sell_price)} - ¥${formatPrice(max_sell_price)}
    
    符合以上条件的饰品商品将会显示在数据列表中。`,
    '筛选效果预览',
    {
      confirmButtonText: '了解',
      type: 'info'
    }
  )
}

// 初始化
onMounted(() => {
  loadSettings()
})

// 监听store中的设置变化
watch(
  () => settingsStore.settings,
  () => {
    updateForm()
  },
  { deep: true }
)
</script>

<style scoped>
.system-settings {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.loading-container {
  min-height: 400px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.settings-form {
  padding: 20px 0;
}

.full-width {
  width: 100%;
}

.form-help {
  font-size: 13px;
  color: #8c8c8c;
  margin-top: 8px;
  line-height: 1.5;
  font-style: italic;
}

.form-actions {
  margin-top: 32px;
  text-align: center;
}

.save-btn {
  background: linear-gradient(135deg, #1890ff, #40a9ff);
  border: none;
  border-radius: 8px;
  height: 48px;
  padding: 0 32px;
  font-weight: 600;
  transition: all 0.3s ease;
}

.save-btn:hover {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(24, 144, 255, 0.3);
}

.help-card,
.preview-card {
  margin-top: 24px;
  border-radius: 16px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.help-section {
  padding: 16px;
}

.help-section h4 {
  margin: 0 0 12px 0;
  color: #262626;
}

.help-section ul {
  margin: 0;
  padding-left: 20px;
}

.help-section li {
  margin-bottom: 8px;
  line-height: 1.5;
}

.alert-section {
  margin-top: 16px;
}

.preview-content {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.filter-description {
  padding: 16px;
  background-color: #f8f9fa;
  border-radius: 8px;
}

.filter-description h4 {
  margin: 0 0 8px 0;
  color: #262626;
}

.filter-description p {
  margin: 0;
  line-height: 1.6;
  color: #595959;
}
</style>
