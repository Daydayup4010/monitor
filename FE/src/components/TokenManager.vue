<template>
  <div class="token-manager">
    <div class="card">
      <div class="card-title">🔑 Token管理</div>
      
      <div class="two-cols">
        <!-- UU Token -->
        <div class="token-section">
          <h4 style="font-size: 16px; color: #262626; margin-bottom: 16px;">UU平台Token</h4>
          
          <div class="form-item">
            <label class="form-label">Authorization</label>
            <el-input
              v-model="uuForm.authorization"
              type="textarea"
              :rows="3"
              placeholder="请输入UU平台的Authorization"
            />
          </div>

          <div class="form-item">
            <label class="form-label">UK</label>
            <el-input
              v-model="uuForm.uk"
              placeholder="请输入UU平台的UK"
            />
          </div>

          <div style="text-align: center; margin-top: 20px;">
            <button class="btn btn-primary" @click="submitUUToken" :disabled="loading">
              {{ loading ? '更新中...' : '更新UU Token' }}
            </button>
          </div>

          <div v-if="tokenStatus.uu !== 'no'" style="margin-top: 16px; padding: 12px; background: #fff7e6; border: 1px solid #ffd591; border-radius: 8px;">
            <p style="font-size: 13px; color: #d48806; margin: 0;">⚠️ Token可能已过期或无效</p>
          </div>
        </div>

        <!-- Buff Token -->
        <div class="token-section">
          <h4 style="font-size: 16px; color: #262626; margin-bottom: 16px;">Buff平台Token</h4>
          
          <div class="form-item">
            <label class="form-label">Session</label>
            <el-input
              v-model="buffForm.session"
              type="textarea"
              :rows="3"
              placeholder="请输入Buff平台的Session"
            />
          </div>

          <div class="form-item">
            <label class="form-label">CSRF Token</label>
            <el-input
              v-model="buffForm.csrf_token"
              placeholder="请输入Buff平台的CSRF Token"
            />
          </div>

          <div style="text-align: center; margin-top: 20px;">
            <button class="btn btn-primary" @click="submitBuffToken" :disabled="loading">
              {{ loading ? '更新中...' : '更新Buff Token' }}
            </button>
          </div>

          <div v-if="tokenStatus.buff !== 'no'" style="margin-top: 16px; padding: 12px; background: #fff7e6; border: 1px solid #ffd591; border-radius: 8px;">
            <p style="font-size: 13px; color: #d48806; margin: 0;">⚠️ Token可能已过期或无效</p>
          </div>
        </div>
      </div>

      <div style="text-align: center; margin-top: 32px;">
        <button class="btn btn-secondary" @click="verifyAllTokens" :disabled="loading">
          验证所有Token
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { platformTokenApi } from '@/api'
import { showMessage } from '@/utils/message'
import type { UUToken, BuffToken, TokenStatus } from '@/types'

const tokenStatus = ref<TokenStatus>({
  uu: 'yes',
  buff: 'yes'
})
const loading = ref(false)

const uuForm = reactive<UUToken>({
  authorization: '',
  uk: '',
})

const buffForm = reactive<BuffToken>({
  session: '',
  csrf_token: '',
})

const submitUUToken = async () => {
  if (!uuForm.authorization || !uuForm.uk) {
    showMessage.warning('请填写完整的UU Token信息')
    return
  }

  loading.value = true
  try {
    const response = await platformTokenApi.updateUUToken(uuForm)
    if (response.code === 1) {
      showMessage.success('UU Token更新成功')
      uuForm.authorization = ''
      uuForm.uk = ''
      await loadTokenStatus()
    }
  } catch (error) {
    console.error('UU token update failed:', error)
  } finally {
    loading.value = false
  }
}

const submitBuffToken = async () => {
  if (!buffForm.session || !buffForm.csrf_token) {
    showMessage.warning('请填写完整的Buff Token信息')
    return
  }

  loading.value = true
  try {
    const response = await platformTokenApi.updateBuffToken(buffForm)
    if (response.code === 1) {
      showMessage.success('Buff Token更新成功')
      buffForm.session = ''
      buffForm.csrf_token = ''
      await loadTokenStatus()
    }
  } catch (error) {
    console.error('Buff token update failed:', error)
  } finally {
    loading.value = false
  }
}

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

onMounted(() => {
  loadTokenStatus()
})
</script>

<style scoped>
/* 所有样式在unified.css中 */
.token-section {
  padding: 20px;
  background: #fafafa;
  border-radius: 8px;
}
</style>



