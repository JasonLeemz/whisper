<script>
import {ExclamationCircleOutlined} from "@ant-design/icons-vue";
import axios from "axios";
import DrawerEquip from "@/components/DrawerEquip.vue";

export default {
  components: {DrawerEquip, ExclamationCircleOutlined},
  emits: ['drawerSearch'],
  props: {
    queryResult: {
      tips: '',
      referer: '',
      data: {},
    }, // 父组件传递的数据类型
    formData: Object,
  },
  data() {
    return {
      equips: [],
      drawer: {
        show: 0,
        isLoading: false,
        data: {},
      },
    }
  },
  watch: {
    queryResult:{
      handler(){
        this.equips = this.queryResult.data
        if (this.formData != null && this.formData.key_words != null) {
          for (let i in this.equips) {
            this.equips[i].name = this.highlight(this.formData.key_words, this.equips[i].name)
            this.equips[i].desc = this.highlight(this.formData.key_words, this.equips[i].desc)
          }
        }

        // keywords map["string,string"]bool
        if (this.formData != null && this.formData.keywords != null){
          for (let i in this.equips) {
            // 遍历 Map 的键值对
            for (let keywords in this.formData.keywords){
              if (this.formData.keywords[keywords]){
                let arr = keywords.split(",");
                for (let k in arr) {
                  this.equips[i].name = this.highlight(arr[k], this.equips[i].name)
                  this.equips[i].desc = this.highlight(arr[k], this.equips[i].desc)
                }
              }
            }
          }
        }
      },
      immediate:true,// 这个属性是重点啦
    },
  },
  methods: {
    showDrawer(map, platform, version, id) {
      this.drawer.show++
      this.drawer.isLoading = true

      axios.post('/equip/roadmap', {
            'platform': platform,
            'version': version,
            'id': id,
            'map': [map],
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
    skipSearch(keywords) {
      if (this.queryResult.referer === 'equip-box') {
        return
      }
      this.drawer.show = 0
      this.$emit('drawerSearch', keywords)
    },
    highlight(keywords, text) {
      let newText = text.replace(new RegExp(`(${keywords})`, 'g'), `<em>$1</em>`);
      return newText
    },
  },
  created() {
  },
  computed: {},
  mounted() {

  }
}
</script>

<template>
  <a-descriptions>
    <a-descriptions-item>{{ queryResult.tips }}</a-descriptions-item>
  </a-descriptions>
  <div class="result-card" v-for="(item,i) in equips" :key="i">
    <a-space direction="vertical">
      <a-card :hoverable="true" @click="showDrawer(item.maps, item.platform,item.version,item.id)">
        <a-card-meta>
          <template #title>
            <span v-html="item.name"/>
          </template>
          <template #avatar>
            <a-avatar :src="item.icon"/>
          </template>
        </a-card-meta>
        <div class="ant-tag-wrap">
          <a-tag v-for="tag in item.tags" :key="tag.id" color="blue" :title="tag">{{ tag }}</a-tag>
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

  <DrawerEquip
      @skip-search="skipSearch"
      :equip-result="drawer"/>
</template>
