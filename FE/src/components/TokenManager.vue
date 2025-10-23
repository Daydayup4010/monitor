<template>
  <div class="token-manager">
    <el-row :gutter="24">
      <!-- UU Token -->
      <el-col :span="12">
        <el-card class="token-card">
          <template #header>
            <div class="card-header">
              <div class="header-info">
                <el-icon size="20" color="#1890ff"><Key /></el-icon>
                <span>UU平台Token</span>
              </div>
              <el-tag :type="isTokenValid('uu') ? 'success' : 'danger'">
                {{ tokenStatus.uu === 'no' ? '有效' : '无效' }}
              </el-tag>
            </div>
          </template>
          
          <el-form
            ref="uuFormRef"
            :model="uuForm"
            :rules="uuRules"
            label-width="120px"
            @submit.prevent="submitUUToken"
          >
            <el-form-item label="Authorization" prop="authorization">
              <el-input
                v-model="uuForm.authorization"
                type="textarea"
                :rows="3"
                placeholder="请输入UU平台的Authorization"
                show-word-limit
                maxlength="500"
              />
            </el-form-item>
            
            <el-form-item label="UK" prop="uk">
              <el-input
                v-model="uuForm.uk"
                placeholder="请输入UU平台的UK"
                show-word-limit
                maxlength="100"
              />
            </el-form-item>
            
            <el-form-item>
              <el-button
                type="primary"
                :loading="loading"
                @click="submitUUToken"
                class="submit-btn"
              >
                更新UU Token
              </el-button>
              <el-button @click="resetUUForm">重置</el-button>
            </el-form-item>
          </el-form>
          
          <el-alert
            v-if="tokenStatus.uu !== 'no'"
            title="Token无效提醒"
            type="warning"
            description="当前UU Token可能已过期或无效，请更新Token以确保数据正常获取。"
            show-icon
            :closable="false"
          />
        </el-card>
      </el-col>
      
      <!-- Buff Token -->
      <el-col :span="12">
        <el-card class="token-card">
          <template #header>
            <div class="card-header">
              <div class="header-info">
                <el-icon size="20" color="#52c41a"><Key /></el-icon>
                <span>Buff平台Token</span>
              </div>
              <el-tag :type="isTokenValid('buff') ? 'success' : 'danger'">
                {{ tokenStatus.buff === 'no' ? '有效' : '无效' }}
              </el-tag>
            </div>
          </template>
          
          <el-form
            ref="buffFormRef"
            :model="buffForm"
            :rules="buffRules"
            label-width="120px"
            @submit.prevent="submitBuffToken"
          >
            <el-form-item label="Session" prop="session">
              <el-input
                v-model="buffForm.session"
                type="textarea"
                :rows="3"
                placeholder="请输入Buff平台的Session"
                show-word-limit
                maxlength="500"
              />
            </el-form-item>
            
            <el-form-item label="CSRF Token" prop="csrf_token">
              <el-input
                v-model="buffForm.csrf_token"
                placeholder="请输入Buff平台的CSRF Token"
                show-word-limit
                maxlength="100"
              />
            </el-form-item>
            
            <el-form-item>
              <el-button
                type="primary"
                :loading="loading"
                @click="submitBuffToken"
                class="submit-btn"
              >
                更新Buff Token
              </el-button>
              <el-button @click="resetBuffForm">重置</el-button>
            </el-form-item>
          </el-form>
          
          <el-alert
            v-if="tokenStatus.buff !== 'no'"
            title="Token无效提醒"
            type="warning"
            description="当前Buff Token可能已过期或无效，请更新Token以确保数据正常获取。"
            show-icon
            :closable="false"
          />
        </el-card>
      </el-col>
    </el-row>
    
    <!-- Token说明 -->
    <el-card class="help-card">
      <template #header>
        <div class="card-header">
          <el-icon size="20"><QuestionFilled /></el-icon>
          <span>Token获取说明</span>
        </div>
      </template>
      
      <el-collapse>
        <el-collapse-item title="如何获取UU平台Token?" name="uu">
          <div class="help-content">
            <ol>
              <li>打开UU平台网站并登录账户</li>
              <li>按F12打开开发者工具</li>
              <li>切换到Network标签页</li>
              <li>刷新页面或进行任何操作</li>
              <li>在请求头中找到Authorization和uk字段</li>
              <li>复制对应的值到上方表单中</li>
            </ol>
            <el-alert
              title="注意事项"
              type="info"
              description="Token具有时效性，一般需要定期更新。建议在Token即将过期前主动更新。"
              show-icon
              :closable="false"
            />
          </div>
        </el-collapse-item>
        
        <el-collapse-item title="如何获取Buff平台Token?" name="buff">
          <div class="help-content">
            <ol>
              <li>打开Buff平台网站并登录账户</li>
              <li>按F12打开开发者工具</li>
              <li>切换到Network标签页</li>
              <li>刷新页面或进行任何操作</li>
              <li>在请求头中找到Cookie中的session字段</li>
              <li>在页面源码或请求中找到csrf_token</li>
              <li>复制对应的值到上方表单中</li>
            </ol>
            <el-alert
              title="安全提醒"
              type="warning"
              description="Token是敏感信息，请妥善保管，不要泄露给他人。定期更新Token可以提高安全性。"
              show-icon
              :closable="false"
            />
          </div>
        </el-collapse-item>
      </el-collapse>
    </el-card>
    
    <!-- 验证按钮 -->
    <div class="verify-section">
      <el-button
        type="primary"
        size="large"
        :loading="loading"
        @click="verifyAllTokens"
        class="verify-btn"
      >
        <el-icon><CircleCheck /></el-icon>
        验证所有Token
      </el-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { platformTokenApi } from '@/api'
import { showMessage } from '@/utils/message'
import type { UUToken, BuffToken, TokenStatus } from '@/types'

// Token状态
const tokenStatus = ref<TokenStatus>({
  uu: 'yes',
  buff: 'yes'
})
const loading = ref(false)

// 表单引用
const uuFormRef = ref()
const buffFormRef = ref()

// UU表单数据
const uuForm = reactive<UUToken>({
  authorization: '',
  uk: '',
})

// Buff表单数据
const buffForm = reactive<BuffToken>({
  session: '',
  csrf_token: '',
})

// UU表单验证规则
const uuRules = {
  authorization: [
    { required: true, message: '请输入Authorization', trigger: 'blur' },
    { min: 10, message: 'Authorization长度不能少于10位', trigger: 'blur' }
  ],
  uk: [
    { required: true, message: '请输入UK', trigger: 'blur' },
    { min: 5, message: 'UK长度不能少于5位', trigger: 'blur' }
  ]
}

// Buff表单验证规则
const buffRules = {
  session: [
    { required: true, message: '请输入Session', trigger: 'blur' },
    { min: 10, message: 'Session长度不能少于10位', trigger: 'blur' }
  ],
  csrf_token: [
    { required: true, message: '请输入CSRF Token', trigger: 'blur' },
    { min: 5, message: 'CSRF Token长度不能少于5位', trigger: 'blur' }
  ]
}

// 提交UU Token
const submitUUToken = async () => {
  if (!uuFormRef.value) return
  
  try {
    await uuFormRef.value.validate()
    loading.value = true
    const response = await platformTokenApi.updateUUToken(uuForm)
    if (response.code === 1) {
      showMessage.success('UU Token更新成功')
      resetUUForm()
      await loadTokenStatus()
    }
  } catch (error) {
    console.error('UU token update failed:', error)
  } finally {
    loading.value = false
  }
}

// 提交Buff Token
const submitBuffToken = async () => {
  if (!buffFormRef.value) return
  
  try {
    await buffFormRef.value.validate()
    loading.value = true
    const response = await platformTokenApi.updateBuffToken(buffForm)
    if (response.code === 1) {
      showMessage.success('Buff Token更新成功')
      resetBuffForm()
      await loadTokenStatus()
    }
  } catch (error) {
    console.error('Buff token update failed:', error)
  } finally {
    loading.value = false
  }
}

// 重置UU表单
const resetUUForm = () => {
  uuForm.authorization = ''
  uuForm.uk = ''
  uuFormRef.value?.clearValidate()
}

// 重置Buff表单
const resetBuffForm = () => {
  buffForm.session = ''
  buffForm.csrf_token = ''
  buffFormRef.value?.clearValidate()
}

// 加载Token状态
const loadTokenStatus = async () => {
  try {
    const response = await platformTokenApi.verifyTokens()
    if (response.code === 1 && response.data) {
      tokenStatus.value = response.data
    }
  } catch (error) {
    console.error('获取Token状态失败:', error)
  }
}

// 验证所有Token
const verifyAllTokens = async () => {
  loading.value = true
  try {
    const response = await platformTokenApi.manualVerifyTokens()
    if (response.code === 1) {
      showMessage.success('Token验证完成')
      if (response.data) {
        tokenStatus.value = response.data
      }
    }
  } catch (error) {
    console.error('验证Token失败:', error)
  } finally {
    loading.value = false
  }
}

// Token是否有效
const isTokenValid = (platform: 'uu' | 'buff'): boolean => {
  return tokenStatus.value[platform] === 'no'
}

// 页面挂载时加载Token状态
onMounted(() => {
  loadTokenStatus()
})
</script>

<style scoped>
.token-manager {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.token-card {
  border-radius: 16px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-info {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: bold;
}

.submit-btn {
  background: linear-gradient(135deg, #1890ff, #40a9ff);
  border: none;
  border-radius: 8px;
  height: 40px;
  padding: 0 24px;
  font-weight: 600;
}

.help-card {
  margin-top: 24px;
  border-radius: 16px;
}

/* 移动端响应式样式 */
@media (max-width: 768px) {
  :deep(.el-row) {
    margin: 0 !important;
  }
  
  :deep(.el-col) {
    width: 100% !important;
    margin-bottom: 16px;
  }
  
  :deep(.el-form-item__label) {
    width: 80px !important;
    font-size: 13px;
  }
  
  .card-header {
    flex-direction: column;
    gap: 8px;
    align-items: flex-start;
  }
  
  .header-info {
    font-size: 14px;
  }
  
  .submit-btn {
    width: 100%;
    margin-bottom: 8px;
  }
  
  :deep(.el-form-item__content) {
    margin-left: 80px !important;
  }
}

@media (max-width: 480px) {
  :deep(.el-form-item__label) {
    width: 70px !important;
    font-size: 12px;
  }
  
  :deep(.el-form-item__content) {
    margin-left: 70px !important;
  }
  
  .header-info {
    font-size: 13px;
  }
  
  :deep(.el-input__wrapper) {
    padding: 8px 12px;
  }
  
  :deep(.el-textarea__inner) {
    padding: 8px 12px;
    font-size: 13px;
  }
}

.help-content {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.help-content ol {
  margin: 0;
  padding-left: 20px;
}

.help-content li {
  margin-bottom: 8px;
  line-height: 1.5;
}

.verify-section {
  display: flex;
  justify-content: center;
  padding: 20px 0;
}

.verify-btn {
  background: linear-gradient(135deg, #52c41a, #73d13d);
  border: none;
  border-radius: 12px;
  height: 48px;
  padding: 0 32px;
  font-size: 16px;
  font-weight: 600;
  transition: all 0.3s ease;
}

.verify-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(82, 196, 26, 0.4);
}
</style>
