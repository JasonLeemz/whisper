<script setup>
import {RouterLink, RouterView} from 'vue-router'
import {LikeOutlined, SearchOutlined, ShareAltOutlined, SoundOutlined} from '@ant-design/icons-vue';
</script>

<script>
import axios from 'axios';
import IconDonate from "@/components/icons/IconDonate.vue";

export default {
  components: {
    IconDonate
  },
  data() {
    return {
      donate: {
        show: false,
        title: '维护说明',
        placement: 'top',
        content: ' 作为一个免费、纯净无广告的开源项目，能得到大家的喜欢真的是意外又惊喜。由于服务器费用昂贵，维护也需要时间和精力。如果经济允许，可以赞助支持一下作者，留言都有认真看，谢谢你们❤️！ ',
      },
      current: ['SearchBox'],
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
      wm_content: {
        content: [],
        font: {
          fontSize: 12,
        },
        rotate: -22,
        gap: [100, 100],
        offset: [],
      },
      loading_percent: 0,
      isLoading: false,
    }
  },
  watch: {
    'loading_percent'(p) {
      this.isLoading = !(p === 0 || p === 100);
    },
    'current'(n, o) {
      let old = o[0]
      if (n[0] === 'Donate') {
        this.donate.show = true
        this.current = [old]
      }
    }
  },
  mounted() {
    if (window.location.pathname === "/equip") {
      this.current = ['EquipBox'];
    } else if (window.location.pathname === "/search") {
      this.current = ['SearchBox'];
    } else if (window.location.pathname === "/") {
      this.current = ['HomePage'];
    }

    // 使用 Axios 发起请求获取服务器数据
    axios.get('/version')
        .then(response => {
          // 将服务器返回的数据更新到组件的 serverData 字段
          this.version_list = response.data.data;
          this.wm_content.content = [
            '端游:' + this.version_list.LOL.version + '(' + this.version_list.LOL.update_time + ')',
            '手游:' + this.version_list.LOLM.version + '(' + this.version_list.LOLM.update_time + ')',
          ]
        })
        .catch(error => {
          console.error('Error fetching server data:', error);
        });
  },
  methods: {
    showProgress(p) {
      this.loading_percent = p
    },
    closeDrawer() {
      this.donate.show = false
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

  <a-watermark v-bind="wm_content" style="position:absolute;width: 100%;height: 100%;top: 0;z-index: -1">
    <div/>
  </a-watermark>

  <a-drawer :title="donate.title" :placement="donate.placement" :open="donate.show" @close="closeDrawer">
    <a-row>
      <a-col :span="6">
        <IconDonate/>
      </a-col>
      <a-col :span="18" class="donate-content">{{ donate.content }}</a-col>
    </a-row>
  </a-drawer>

  <div class="nav-bar">
    <a-menu v-model:selectedKeys="current" mode="horizontal" class="bottom-nav">
      <a-menu-item key='HomePage'>
        <SoundOutlined/>
        <span><RouterLink to="/">资讯</RouterLink></span>
      </a-menu-item>
      <a-menu-item key='SearchBox'>
        <SearchOutlined/>
        <span><RouterLink to="/search">检索</RouterLink></span>
      </a-menu-item>
      <a-menu-item key='EquipBox'>
        <ShareAltOutlined/>
        <span><RouterLink to="/equip">装备</RouterLink></span>
      </a-menu-item>
      <a-menu-item key='Donate'>
        <LikeOutlined/>
        <span>赞赏</span>
      </a-menu-item>
    </a-menu>
  </div>
</template>

<style scoped>
.progress-bar {
  position: fixed;
  width: 100%;
  height: 10px;
  z-index: 99999;
}

.nav-bar {
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

.version-tip {
  display: inline-block;
  width: 40px;
}

.bottom-nav {
  box-shadow: 0 -2px 10px rgba(0, 0, 0, .05);
}

.donate-content {
  padding-left: 10px;
  color: #00aeec;
  font-size: 12px;
  text-align: left;
}
</style>