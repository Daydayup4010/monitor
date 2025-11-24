// 系统图标配置 - 集中管理所有图标
// 你可以替换这里的emoji为自定义图标路径

export const ICONS = {
  // 认证页面图标
  login: '🔐',           // 登录页面
  register: '✨',        // 注册页面
  resetPassword: '🔑',   // 重置密码页面
  
  // 功能模块图标
  data: '📊',            // 饰品数据
  settings: '⚙️',        // 设置
  user: '👥',            // 用户管理
  
  // 统计卡片图标
  statUsers: '👥',       // 总用户数
  statVip: '👑',         // VIP用户
  
  // 用户类型标识
  userAdmin: '👨‍💼',     // 管理员
  userVip: '👑',         // VIP会员
  userNormal: '👤',      // 普通用户
  
  // 操作图标
  refresh: '🔄',         // 刷新
  search: '🔍',          // 搜索
  warning: '⚠️',         // 警告
}

// 图标使用说明
export const ICON_USAGE = {
  'ICONS.login': 'Login.vue - 页面标题',
  'ICONS.register': 'Register.vue - 页面标题',
  'ICONS.resetPassword': 'ResetPassword.vue - 页面标题',
  'ICONS.data': 'Home.vue - 卡片标题',
  'ICONS.settings': 'Settings.vue - 卡片标题',
  'ICONS.user': 'UserManager.vue - 卡片标题 & Admin.vue - 左侧菜单',
  'ICONS.statUsers': 'UserManager.vue - 统计卡片',
  'ICONS.statVip': 'UserManager.vue - 统计卡片',
  'ICONS.userAdmin': 'Settings.vue - 管理员徽章',
  'ICONS.userVip': 'Settings.vue - VIP徽章',
  'ICONS.userNormal': 'Settings.vue - 普通用户徽章',
  'ICONS.refresh': 'Home.vue - 刷新按钮',
  'ICONS.search': 'UserManager.vue - 搜索框',
  'ICONS.warning': 'Settings.vue - VIP提示',
}

// 如何替换图标：
// 1. 如果使用emoji：直接修改上面的ICONS对象中的值
// 2. 如果使用图片：将值改为图片路径，如 '/icons/login.svg'
//    然后在组件中使用 <img :src="ICONS.login" /> 而不是直接显示文字
// 3. 如果使用Element Plus图标：导入后使用组件形式

export default ICONS

