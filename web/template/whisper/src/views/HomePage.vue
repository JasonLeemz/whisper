<script setup>
</script>
<script>
import axios from 'axios';
// import {FilterOutlined, AppleOutlined, AndroidOutlined} from "@ant-design/icons-vue";
// import LoadingList from "@/components/LoadingList.vue";
import ListVersion from "@/components/ListVersion.vue";

export default {
  emits: ['loadingEvent'],
  components: {
     ListVersion,
  },
  data() {
    return {
      skins: [],
      platform: Number(0),
      detailVersion: "",
      SkeletonState: {
        show: false,
        isLoading: false,
      },
      loadingState: {
        condCheckbox: true,
      },
      query: {
        tips: '',
        data: [],
      }
    }
  },
  methods: {
    heroSkins() {
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
    getVersionList() {
      this.$emit('loadingEvent', 0)

      this.SkeletonState.show = true
      this.$emit('loadingEvent', 30)

      axios.post('/version/list', {
        'platform': this.platform,
      })
          .then(response => {
            // 将服务器返回的数据更新到组件的 serverData 字段
            this.query.data = response.data.data.data;
            this.query.tips = response.data.data.tips;
          })
          .catch(error => {
            console.error('Error fetching server data:', error);
          }).finally(() => {
            this.SkeletonState.show = false
            this.$emit('loadingEvent', 100)
          }
      );
    }
  },
  created() {
    // this.heroSkins()
    this.getVersionList()
  },
  watch: {
    'platform'() {
      this.getVersionList()
    },
  },
  mounted() {

  }
}
</script>


<template>
  <!--  <div class="ant-carousel-wrap">-->
  <!--    <div class="carousel-border-radius">-->
  <!--      <a-carousel  :dot-position="'top'" :dots="false">-->
  <!--        <template v-for="(data,i) in skins" :key="i">-->
  <!--          <div class="carousel-wrap">-->
  <!--            <div class="carousel-mask"></div>-->
  <!--            <img :src="data.mainImg" class="carousel-img" />-->
  <!--            <div class="carousel-content">-->

  <!--              <a-tag color="green" class="carousel-tip">{{i+1}}/{{skins.length}}</a-tag>-->

  <!--              <div class="name-mask">-->
  <!--                <em class="carousel-name">{{data.heroName}}</em>-->
  <!--                <em class="carousel-title">{{data.skinName}}</em>-->
  <!--              </div>-->

  <!--              &lt;!&ndash;          <p class="carousel-desc">{{data.description}}</p>&ndash;&gt;-->
  <!--            </div>-->
  <!--          </div>-->
  <!--        </template>-->
  <!--      </a-carousel>-->
  <!--    </div>-->
  <!--  </div>-->

  <div class="version-platform-switch">
    <a-radio-group v-model:value="platform" button-style="solid" size="small" name="platform">
      <a-radio-button :value="0">端游</a-radio-button>
      <a-radio-button :value="1">手游</a-radio-button>
    </a-radio-group>
  </div>

  <a-space direction="vertical" :style="{ width: '100%' }" class="wrap">
    <a-layout>
      <a-layout-content>
        <template v-if="!SkeletonState.show">
          <ListVersion
              :version-result="query"
          />
        </template>
      </a-layout-content>
    </a-layout>
  </a-space>
</template>

<style scoped>

</style>