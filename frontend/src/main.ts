import { createApp } from 'vue'
import { createPinia } from 'pinia'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import zhCn from 'element-plus/es/locale/lang/zh-cn'
import * as ElementPlusIconsVue from '@element-plus/icons-vue'

import App from './App.vue'
import router from './router'
import { useSystemStore } from './stores/system'

import './assets/styles/global.scss'

const app = createApp(App)
const pinia = createPinia()

app.use(pinia)

// 注册所有图标
for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
  app.component(key, component)
}

app.use(router)
app.use(ElementPlus, {
  locale: zhCn
})

// 挂载应用（不等待系统配置加载，避免阻塞）
app.mount('#app')

// 异步加载系统配置（不阻塞应用挂载）
const systemStore = useSystemStore()
systemStore.loadConfigs()
