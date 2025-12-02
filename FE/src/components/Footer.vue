<template>
  <footer class="site-footer">
    <div class="footer-content">
      <!-- 顶部区域：Logo 和链接 -->
      <div class="footer-top">
        <div class="footer-brand">
          <span class="brand-name">CS Goods</span>
          <span class="brand-desc">CS2 饰品数据平台</span>
        </div>
        
        <div class="footer-links">
          <div class="link-group">
            <h4>快速链接</h4>
            <router-link to="/login">登录</router-link>
            <router-link to="/register">注册</router-link>
          </div>
          <div class="link-group">
            <h4>法律条款</h4>
            <router-link to="/user-agreement">用户协议</router-link>
            <router-link to="/privacy-policy">隐私政策</router-link>
          </div>
          <div class="link-group">
            <h4>联系我们</h4>
            <a href="mailto:goods.monitor@foxmail.com">
              <span class="contact-icon">📧</span>
              goods.monitor@foxmail.com
            </a>
          </div>
        </div>
      </div>
      
      <!-- 分割线 -->
      <div class="footer-divider"></div>
      
      <!-- 底部区域：版权和备案 -->
      <div class="footer-bottom">
        <div class="copyright">
          © {{ currentYear }} CS Goods. All rights reserved.
        </div>
        <div class="icp">
          <a href="https://beian.miit.gov.cn/" target="_blank" rel="noopener noreferrer">
            闽ICP备2024049911号-3
          </a>
        </div>
      </div>
    </div>
  </footer>
</template>

<script setup lang="ts">
const currentYear = new Date().getFullYear()
</script>

<style scoped>
.site-footer {
  background: linear-gradient(180deg, #1a1a2e 0%, #16213e 100%);
  color: #a0a0a0;
  padding: 0;
}

.footer-content {
  max-width: 1200px;
  margin: 0 auto;
  padding: 40px 40px 24px;
}

/* 顶部区域 */
.footer-top {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 60px;
}

.footer-brand {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.brand-name {
  font-size: 24px;
  font-weight: 700;
  background: linear-gradient(135deg, #1890ff, #40a9ff);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.brand-desc {
  font-size: 14px;
  color: #666;
}

/* 链接组 */
.footer-links {
  display: flex;
  gap: 60px;
}

.link-group {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.link-group h4 {
  font-size: 14px;
  font-weight: 600;
  color: #fff;
  margin-bottom: 4px;
}

.link-group a {
  font-size: 13px;
  color: #a0a0a0;
  text-decoration: none;
  transition: color 0.3s;
  display: flex;
  align-items: center;
  gap: 6px;
}

.link-group a:hover {
  color: #1890ff;
}

.contact-icon {
  font-size: 14px;
}

/* 分割线 */
.footer-divider {
  height: 1px;
  background: linear-gradient(90deg, transparent, rgba(255,255,255,0.1), transparent);
  margin: 32px 0 20px;
}

/* 底部区域 */
.footer-bottom {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 12px;
}

.copyright {
  color: #666;
}

.icp a {
  color: #666;
  text-decoration: none;
  transition: color 0.3s;
}

.icp a:hover {
  color: #1890ff;
}

/* 响应式 */
@media (max-width: 768px) {
  .footer-content {
    padding: 32px 24px 20px;
  }
  
  .footer-top {
    flex-direction: column;
    gap: 32px;
  }
  
  .footer-links {
    flex-wrap: wrap;
    gap: 32px;
  }
  
  .link-group {
    min-width: 120px;
  }
  
  .footer-bottom {
    flex-direction: column;
    gap: 8px;
    text-align: center;
  }
}

@media (max-width: 480px) {
  .footer-content {
    padding: 24px 16px 16px;
  }
  
  .footer-links {
    gap: 24px;
  }
  
  .brand-name {
    font-size: 20px;
  }
}
</style>
