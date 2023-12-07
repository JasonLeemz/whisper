<script>

import FeedList from "@/components/FeedList.vue";
import axios from "axios";

export default {
  components: {FeedList},
  props: {
    runeResult: {
      show: 0,
      isLoading: false,
      data: {},
    },
  },
  data() {
    return {
      sideDrawer: {
        show: false,
        title: '适配英雄',
      },
      feed:{
        show: 0,
        isLoading: false,
        data:[],
      },
    }
  },
  watch: {
    'runeResult.show'() {
      this.sideDrawer.show = true
    },
    'runeResult.isLoading'(isLoading) {
      this.sideDrawer.isLoading = isLoading
    },
    'runeResult.data'(data) {
      // 获取推荐视频列表
      axios.post('/strategy/rune', {
        'platform': data.platform,
        'id': data.id,
      }).then(response => {
        // 将服务器返回的数据更新到组件的 serverData 字段
        // console.log(response)
        this.feed.show++
        this.feed.isLoading = true
        this.feed.data = response.data.data
      }).catch(error => {
        console.error('Error fetching server data:', error);
      }).finally(() => {
        this.feed.isLoading = false
      });
    },
  },
  methods: {}
}
</script>

<template>
  <a-drawer
      v-model:open="sideDrawer.show"
      class="custom-class"
      root-class-name="root-class-name"
      :title="sideDrawer.title"
      placement="right"
  >
    <template v-if="sideDrawer.isLoading">
      <a-skeleton active/>
    </template>

    <a-empty
        v-if="!sideDrawer.isLoading && runeResult.data.data == null"
        description="当前版本该符文缺乏足够的样本数据"/>

    <!-- 适配英雄 START-->
    <template v-if="!sideDrawer.isLoading && runeResult.data.data != null">
      <h4 class="equip-suit-hero">适配英雄</h4>
      <template v-for="(hero,i) in runeResult.data.data" :key="i">
        <a-tooltip :title="hero.name">
          <img class="equip-hero-avatar"
               :title="hero.name"
               :src="hero.icon"
               :alt="hero.name"
          />

        </a-tooltip>
      </template>
    </template>
    <!-- 适配英雄 END-->

    <FeedList :feed-list="feed"/>
  </a-drawer>
</template>
