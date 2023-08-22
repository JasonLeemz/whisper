<script setup lang="ts">
import { RouterLink, RouterView } from 'vue-router'

import { h, ref } from 'vue';
import { SearchOutlined, ClusterOutlined, SettingOutlined,ShareAltOutlined } from '@ant-design/icons-vue';
const current = ref(['SearchBox']);
</script>

<script lang="ts">
import axios from 'axios';

export default {
  data() {
    return {
      lol_version: null,
      lolm_version: null,
    }
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
  }
}
</script>

<template>
  <a-row>
    <a-col :span="24">
      <div class="alert-banner">
        Tips: 当前端游版本为 {{ lol_version }}，手游版本为 {{ lolm_version }}
      </div>
      <a-menu v-model:selectedKeys="current" mode="horizontal" >
        <a-menu-item key='SearchBox' >
          <SearchOutlined />
          <span><RouterLink to="/">检索</RouterLink></span>
        </a-menu-item>
        <a-menu-item key='EquipBox' >
          <ShareAltOutlined />
          <span><RouterLink to="/equip">装备</RouterLink></span>
        </a-menu-item>
      </a-menu>
    </a-col>
  </a-row>
  <div class="blank"></div>
<!--  <a-space direction="vertical" :style="{ width: '100%' }" >-->
<!--    <a-layout>-->
<!--      <a-menu v-model:selectedKeys="current" mode="horizontal" :items="items" />-->
<!--      Tips: 当前端游版本为 {{ lol_version }}，手游版本为 {{ lolm_version }}-->
<!--&lt;!&ndash;      <a-layout-header>&ndash;&gt;-->
<!--&lt;!&ndash;          <a-breadcrumb>&ndash;&gt;-->
<!--&lt;!&ndash;            <a-space>&ndash;&gt;-->
<!--&lt;!&ndash;              <RouterLink to="/">检索</RouterLink>&ndash;&gt;-->
<!--&lt;!&ndash;              <RouterLink to="/about">|</RouterLink>&ndash;&gt;-->
<!--&lt;!&ndash;              <RouterLink to="/equip">装备</RouterLink>&ndash;&gt;-->
<!--&lt;!&ndash;            </a-space>&ndash;&gt;-->
<!--&lt;!&ndash;          </a-breadcrumb>&ndash;&gt;-->
<!--&lt;!&ndash;      </a-layout-header>&ndash;&gt;-->
<!--    </a-layout>-->
<!--  </a-space>-->
  <RouterView />
</template>
