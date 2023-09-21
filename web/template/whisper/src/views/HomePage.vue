<script setup>
</script>
<script>
import axios from 'axios';

export default {
  emits: [''],
  components: {},
  data() {
    return {
      skins: []
    }
  },
  methods: {},
  created() {
    // 使用 Axios 发起请求获取服务器数据
    axios.post('/hero/skins', {platform: 0, id: ''})
        .then(response => {
          // 将服务器返回的数据更新到组件的 serverData 字段
          this.skins = response.data.data;
        })
        .catch(error => {
          console.error('Error fetching server data:', error);
        });
  },
  watch: {},
  mounted() {

  }
}
</script>


<template>
  <a-carousel  :dot-position="'top'" :dots="false">
    <template v-for="(data,i) in skins" :key="i">
      <div class="carousel-wrap">
        <div class="carousel-mask"></div>
        <img :src="data.mainImg" class="carousel-img">
        <div class="carousel-content">

          <a-tag color="green" class="carousel-tip">{{i+1}}/{{skins.length}}</a-tag>
          <em class="carousel-name">{{data.heroName}}</em>
          <em class="carousel-title">{{data.skinName}}</em>
          <p class="carousel-desc">{{data.description}}</p>
        </div>
      </div>
    </template>
  </a-carousel>
</template>

<style scoped>
:deep(.slick-slide) {
  text-align: center;
  background: url("https://game.gtimg.cn/images/lol/lolstrategy/bg-dhjjc.jpg");
}

</style>