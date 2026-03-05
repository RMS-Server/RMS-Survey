import { createApp } from 'vue'
import { createPinia } from 'pinia'
import Antd from 'ant-design-vue'
import App from './App.vue'
import router from './router'
import { initGlowEffect } from './utils/glowEffect'
// Import styles in correct order
import './styles/theme.css'
import './styles/glassmorphism.css'
import 'ant-design-vue/dist/reset.css'
import './styles/main.css'
import './styles/responsive.css'

const app = createApp(App)
const pinia = createPinia()

app.use(pinia)
app.use(router)
app.use(Antd)

// Initialize glow effect
initGlowEffect()

app.mount('#app')
