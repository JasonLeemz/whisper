<script>
import {ExclamationCircleOutlined} from "@ant-design/icons-vue";
import axios from "axios";
import DrawerRune from "@/components/DrawerRune.vue";
export default {
  components: { DrawerRune,ExclamationCircleOutlined},
  props: {
    queryResult: Object, // 父组件传递的数据类型
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
    showDrawer(platform, version, id) {
      this.drawer.show++
      this.drawer.isLoading = true

      axios.post('/rune/hero/suit', {
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
  }
}
</script>

<template>
  <a-descriptions>
    <a-descriptions-item>{{ queryResult.tips }}</a-descriptions-item>
  </a-descriptions>
  <div class="result-card" v-for="(item,i) in queryResult.data" :key="i">
    <a-space direction="vertical">
      <a-card :hoverable="true" @click="showDrawer(item.platform,item.version,item.id)">
        <a-card-meta :title="item.name">
          <template #avatar>
            <a-avatar :src="item.icon"/>
          </template>
        </a-card-meta>
        <div class="ant-tag-wrap">
          <a-tag v-for="tag in item.tags" :key="tag.id" color="blue">{{ tag }}</a-tag>
        </div>
        <a-tag class="platform-tag" color="warning">
          <template #icon>
            <ExclamationCircleOutlined/>
          </template>
          {{ item.platform === 0 ? '端游' : '手游' }}
        </a-tag>

        <div class="mainText" v-html="item.desc"></div>
      </a-card>
    </a-space>
  </div>

  <DrawerRune :rune-result="drawer"/>
</template>
