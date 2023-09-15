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
      version_list: {
        LOL: {
          version: '',
          update_time: '',
        },
        LOLM: {
          version: '',
          update_time: '',
        },
      },
      loading_percent: 0,
      isLoading: false,
    }
  },
  watch: {
    'loading_percent'(p) {
      if (p === 0 || p === 100) {
        this.isLoading = false
      } else {
        this.isLoading = true
      }
    },
  },
  mounted() {
    // 使用 Axios 发起请求获取服务器数据
    axios.get('/version')
        .then(response => {
          // 将服务器返回的数据更新到组件的 serverData 字段
          this.version_list = response.data.data;
        })
        .catch(error => {
          console.error('Error fetching server data:', error);
        });
  },
  methods: {
    showProgress(p) {
      this.loading_percent = p
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
        stroke-linecap="square"
        class="progress-bar"
    />
  </Transition>
  <RouterView @loading-event="showProgress"/>
  <div class="nav-bar">
    <a-menu v-model:selectedKeys="current" mode="horizontal" class="glass">
      <a-menu-item key='SearchBox'>
        <SearchOutlined/>
        <span><RouterLink to="/">检索</RouterLink></span>
      </a-menu-item>
      <a-menu-item key='EquipBox'>
        <ShareAltOutlined/>
        <span><RouterLink to="/equip">装备</RouterLink></span>
      </a-menu-item>
    </a-menu>

<!--    <p>-->
<!--      Tips: 当前端游版本为-->
<!--      <span :title="version_list.LOL.update_time" class="version-tip">{{ version_list.LOL.version }}</span>-->
<!--      手游版本为-->
<!--      <span :title="version_list.LOLM.update_time" class="version-tip">{{ version_list.LOLM.version }}</span>-->
<!--    </p>-->
  </div>
</template>

<style scoped>
.progress-bar{
  position: fixed;
  width: 100%;
  height: 10px;
  z-index: 99999;
}
.nav-bar{
  position: fixed;
  width: 100%;
  left: 0;
  bottom: 0;
}
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.5s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 1;
}

.version-tip{
  display: inline-block;
  width: 40px;
}

.glass {
  box-shadow: 0 -2px 10px rgba(0,0,0,.05);
}
</style>