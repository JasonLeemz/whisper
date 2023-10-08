<script>
import axios from "axios";
import DrawerVersion from "@/components/DrawerVersion.vue";

export default {
  components: {DrawerVersion},
  props: {
    versionResult: Object, // 父组件传递的数据类型
  },
  data() {
    return {
      drawer: {
        show: 0,
        isLoading: false,
        data: {},
      },
    }
  },
  watch: {},
  methods: {
    showDrawer(version, platform) {
      this.drawer.show++
      this.drawer.isLoading = true

      axios.post('/version/detail', {
            'platform': platform,
            'version': version,
          }
      ).then(response => {
            this.drawer.data = response.data.data
          }
      ).catch(error => {
            console.error('Error fetching server data:', error);
          }
      ).finally(() => {
            this.drawer.isLoading = false
          }
      );
    },
  },
  mounted() {

  }
}
</script>

<template>
  <p class="result-tips"></p>
  <div class="result-card" v-for="(item,i) in versionResult.data" :key="i">
    <a-space direction="vertical">
      <a-card hoverable @click="showDrawer(item.platform,item.vkey)" style="max-width: 500px">
        <template #cover>
          <img :alt="item.name" :src="item.image"/>
        </template>
        <div class="version-title">
          {{ item.title }}
        </div>
        <a-card-meta :title="item.introduction">
        </a-card-meta>
        <div class="version-public-date" v-html="item.public_date"></div>
      </a-card>
    </a-space>
  </div>

  <DrawerVersion :version-result="drawer"/>
</template>

<style scoped>
.ant-card .ant-card-meta-description {
  padding: 0 !important;
}
</style>