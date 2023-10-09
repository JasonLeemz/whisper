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
      loadingImgs: new Map(),
      drawer: {
        show: 0,
        isLoading: false,
        title: '',
        introduction: '',
        data: {},
      },
    }
  },
  watch: {
    versionResult: {
      handler() {
        for (let idx in this.versionResult.data) {
          this.loadingImgs.set("load_" + idx, true);
        }
      },
      immediate: true,// 这个属性是重点啦
    },
  },
  methods: {
    stopLoading(idx) {
      this.loadingImgs.set("load_" + idx, false);
    },
    showDrawer(platform, version, id, desc, introduction) {
      this.drawer.show++
      this.drawer.isLoading = true
      this.drawer.title = desc
      this.drawer.introduction = introduction


      axios.post('/version/detail', {
            'platform': platform,
            'version': version,
            'id': id,
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
  <div class="result-card" v-for="(item,i) in versionResult.data" :key="i">
    <a-space direction="vertical">
      <a-card hoverable
              :loading="loadingImgs.get('load_'+i)"
              @click="showDrawer(item.platform,item.vkey,item.id,item.title,item.introduction)"
              style="max-width: 500px">
        <template #cover>
          <img
              @load="stopLoading(i)"
              :alt="item.name"
              :src="item.image"/>
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

</style>