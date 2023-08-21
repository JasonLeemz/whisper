<script setup>
import { RouterLink, RouterView } from 'vue-router'

import { h, ref } from 'vue';
import { SearchOutlined, ClusterOutlined, SettingOutlined } from '@ant-design/icons-vue';
const current = ref(['SearchBox']);
const items = ref([
  {
    key: 'SearchBox',
    icon: () => h(SearchOutlined),
    label: 'SearchBox',
    title: 'SearchBox',
  },
  {
    key: 'Equipment',
    icon: () => h(ClusterOutlined),
    label: 'Equipment',
    title: 'Equipment',
  },
  {
    key: 'Select Version',
    icon: () => h(SettingOutlined),
    label: 'Select Version',
    title: 'Select Version',
    children: [
      {
        type: 'group',
        label: '当前端游版本为 {{ lol_version }}',
        children: [
          {
            label: '13.14',
            key: 'setting:13.14',
          },
          {
            label: '13.14',
            key: 'setting:13.14',
          },
        ],
      },
      {
        type: 'group',
        label: '当前手游版本为 {{ lolm_version }}',
        children: [
          {
            label: '13.14',
            key: 'setting:13.14',
          },
          {
            label: '13.14',
            key: 'setting:13.14',
          },
        ],
      },
    ],
  }
]);
</script>
<script>
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
    axios.get('http://127.0.0.1:8123/version')
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
      <a-menu v-model:selectedKeys="current" mode="horizontal" :items="items" />
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
