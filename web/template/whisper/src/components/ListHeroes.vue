<script>
import {ExclamationCircleOutlined} from "@ant-design/icons-vue";
import DrawerHero from "./DrawerHero.vue";
import axios from "axios";

export default {
  components: {ExclamationCircleOutlined, DrawerHero},
  props: {
    queryResult: Object, // 父组件传递的数据类型
    formData: Object,
  },
  data() {
    return {
      heroes:{},
      drawer: {
        show: 0,
        isLoading: false,
        data: {},
      }
    }
  },
  watch: {},
  methods: {
    showDrawer(platform, version, id) {
      if (id === "") {
        return
      }
      this.drawer.show++
      this.drawer.isLoading = true

      axios.post('/hero/suit', {
            'platform': platform,
            'hero_id': id,
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
    highlight(keywords, text) {
      let newText = text.replace(new RegExp(`(${keywords})`, 'g'), `<em>$1</em>`);
      return newText
    },
  },
  mounted() {
    this.heroes = this.queryResult.data
    if (this.formData != null) {
      for (let i in this.heroes) {
        this.heroes[i].name = this.highlight(this.formData.key_words, this.heroes[i].name)
        for (let j in this.heroes[i].spell) {
          this.heroes[i].spell[j].desc = this.highlight(this.formData.key_words, this.heroes[i].spell[j].desc)
        }
      }
    }
  }
}
</script>

<template>
  <a-descriptions>
    <a-descriptions-item>{{ queryResult.tips }}</a-descriptions-item>
  </a-descriptions>
  <div class="result-card" v-for="(item,i) in heroes" :key="i">
    <a-space direction="vertical">
      <a-card :hoverable="true" @click="showDrawer(item.platform,item.version,item.id)">
        <a-card-meta>
          <template #title>
            <span v-html="item.name" />
          </template>
          <template #avatar>
            <a-avatar :src="item.icon" class="hero-icon"/>
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
        <a-divider/>
        <div class="hero-desc mainText">
          <ul>
            <li v-for="(s,index) in item.spell" :key="index">
              <img :src="s.icon" class="spell-icon"/>
              <h6>{{ s.name }}</h6>
              <span>{{ s.sort }}</span>
              <div v-html="s.desc" class="spell-desc"></div>
            </li>
          </ul>
        </div>
      </a-card>
    </a-space>
  </div>

  <DrawerHero :hero-result="drawer"/>
</template>
