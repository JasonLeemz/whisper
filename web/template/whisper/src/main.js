import { createApp } from 'vue'
import Antd from 'ant-design-vue';
import App from './App.vue'
import router from './router'
import 'ant-design-vue/dist/reset.css';
import './assets/main.scss'

import axios from 'axios';
// 设置 Axios 全局配置
// axios.defaults.baseURL = '';
axios.defaults.baseURL = 'http://192.168.31.30:8123';
// axios.defaults.baseURL = 'http://whisper.ybdx.xyz:66';


const app = createApp(App)

app.use(router).use(Antd)

app.mount('#app')
