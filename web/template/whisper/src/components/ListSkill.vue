<script>
import {ExclamationCircleOutlined} from "@ant-design/icons-vue";
// import axios from "axios";

export default {
  components: {ExclamationCircleOutlined},
  props: {
    queryResult: Object, // 父组件传递的数据类型
  },
  data() {
    return {
      sideDrawer: {
        show: false,
        title: '',
        data: {},
      },
    }
  },
  watch: {},
  methods: {
    showDrawer(platform, version, id) {
      console.log(platform, version, id)
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

  <a-drawer
      v-model:open="sideDrawer.show"
      class="custom-class"
      root-class-name="root-class-name"
      :title="sideDrawer.title"
      placement="right"
  >
    drawer
  </a-drawer>
</template>
