<script setup>
import {RouterLink, RouterView} from 'vue-router'

import {ref} from 'vue';
import {SearchOutlined, ShareAltOutlined} from '@ant-design/icons-vue';

let current = ref(['SearchBox']);

if (window.location.pathname === "/equip") {
  current = ref(['EquipBox']);
} else if (window.location.pathname === "/") {
  current = ref(['SearchBox']);
}

</script>

<script>
import axios from 'axios';

export default {
  data() {
    return {
      lol_version: null,
      lolm_version: null,
      loading_percent: 0,
      isLoading: false,
    }
  },
  watch: {
    'loading_percent'(p) {
      if (p === 0 || p === 100){
        this.isLoading = false
      }else{
        this.isLoading = true
      }
    },
  },
  mounted() {
    // 使用 Axios 发起请求获取服务器数据
    axios.get('/version')
        .then(response => {
          // 将服务器返回的数据更新到组件的 serverData 字段
          this.lol_version = response.data.data.lol_version;
          this.lolm_version = response.data.data.lolm_version;
        })
        .catch(error => {
          console.error('Error fetching server data:', error);
        });
  },
  methods: {
    showProgress(p) {
      this.loading_percent = p
      // console.log(p)
    }
  }
}
</script>

<template>
  <Transition name="fade">
    <a-progress
        v-show="isLoading"
        :percent="loading_percent"
        size="small"
        :showInfo="false"
        stroke-linecap="square"/>
  </Transition>

  <a-row>
    <a-col :span="24">
      <div class="alert-banner">
        Tips: 当前端游版本为 {{ lol_version }}，手游版本为 {{ lolm_version }}
      </div>
      <a-menu v-model:selectedKeys="current" mode="horizontal">
        <a-menu-item key='SearchBox'>
          <SearchOutlined/>
          <span><RouterLink to="/">检索</RouterLink></span>
        </a-menu-item>
        <a-menu-item key='EquipBox'>
          <ShareAltOutlined/>
          <span><RouterLink to="/equip">装备</RouterLink></span>
        </a-menu-item>
      </a-menu>
    </a-col>
  </a-row>
  <div class="blank"></div>
  <RouterView @loading-event="showProgress"/>
</template>

<style scoped>
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.5s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 1;
}
</style>