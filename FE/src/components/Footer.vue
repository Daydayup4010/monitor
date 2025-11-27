<template>
  <footer class="site-footer">
    <div class="footer-content">
      <div class="footer-info">
        <a href="https://beian.miit.gov.cn/" target="_blank" rel="noopener noreferrer" class="icp-link">
          闽ICP备2024049911号-3
        </a>
        <span class="footer-divider">|</span>
        <span class="contact-info">
          📧 联系我们：<a href="mailto:goods.monitor@foxmail.com">goods.monitor@foxmail.com</a>
        </span>
      </div>
      <div class="footer-copyright">
        © {{ currentYear }} CS Goods. All rights reserved.
      </div>
    </div>
  </footer>
</template>

<script setup lang="ts">
const currentYear = new Date().getFullYear()
</script>

<style scoped>
.site-footer {
  background: linear-gradient(135deg, #f8faff, #e6f7ff);
  border-top: 1px solid rgba(24, 144, 255, 0.15);
  padding: 20px 40px;
  text-align: center;
}

.footer-content {
  max-width: 1200px;
  margin: 0 auto;
}

.footer-info {
  display: flex;
  align-items: center;
  justify-content: center;
  flex-wrap: wrap;
  gap: 8px;
  margin-bottom: 8px;
  font-size: 13px;
  color: #666;
}

.icp-link {
  color: #1890ff;
  text-decoration: none;
  transition: color 0.3s;
}

.icp-link:hover {
  color: #40a9ff;
  text-decoration: underline;
}

.footer-divider {
  color: #d9d9d9;
  margin: 0 8px;
}

.contact-info a {
  color: #1890ff;
  text-decoration: none;
}

.contact-info a:hover {
  text-decoration: underline;
}

.footer-copyright {
  font-size: 12px;
  color: #999;
}

@media (max-width: 768px) {
  .site-footer {
    padding: 16px 20px;
  }

  .footer-info {
    flex-direction: column;
    gap: 6px;
  }

  .footer-divider {
    display: none;
  }
}
</style>

