import { createApp } from 'vue'
import Antd from 'ant-design-vue';
import App from './App.vue'
import router from './router'
import 'ant-design-vue/dist/reset.css';
import './assets/main.css'

import axios from 'axios';
// 设置 Axios 全局配置
axios.defaults.baseURL = 'http://127.0.0.1:8123';

const app = createApp(App)

app.use(router).use(Antd)

app.mount('#app')
