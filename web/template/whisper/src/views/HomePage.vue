<script setup>
</script>
<script>
import axios from 'axios';
import { FilterOutlined,AppleOutlined,AndroidOutlined} from "@ant-design/icons-vue";
import LoadingList from "@/components/LoadingList.vue";
import ListVersion from "@/components/ListVersion.vue";

export default {
  emits: ['loadingEvent'],
  components: {
    LoadingList, ListVersion,  FilterOutlined,AppleOutlined,AndroidOutlined,
  },
  data() {
    return {
      skins: [],
      platform: "1",
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
  methods: {},
  created() {
    this.$emit('loadingEvent', 0)

    // 使用 Axios 发起请求获取服务器数据
    axios.post('/hero/skins', {platform: 0, id: ''})
        .then(response => {
          // 将服务器返回的数据更新到组件的 serverData 字段
          this.skins = response.data.data;
        })
        .catch(error => {
          console.error('Error fetching server data:', error);
        });
    this.SkeletonState.show = true
    this.$emit('loadingEvent', 30)
    axios.get('/version/lolm', {})
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
  },
  watch: {},
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

  <a-float-button-group trigger="hover" :style="{ right: '24px' }">
    <template #icon>
      <FilterOutlined />
    </template>
    <a-float-button>
      <template #icon>
        <AppleOutlined />
      </template>
    </a-float-button>
    <a-float-button>
      <template #icon>
        <AndroidOutlined />
      </template>
    </a-float-button>
  </a-float-button-group>

  <a-space direction="vertical" :style="{ width: '100%' }" class="wrap">
    <a-layout>
      <a-layout-content>
        <LoadingList
            v-if="SkeletonState.show"
            :skeleton-state="SkeletonState"/>

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