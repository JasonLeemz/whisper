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
      heroes: {},
      drawer: {
        show: 0,
        background: '',
        isLoading: false,
        data: {},
      }
    }
  },
  watch: {
    queryResult: {
      handler() {
        this.heroes = this.queryResult.data
        if (this.formData != null) {
          for (let i in this.heroes) {
            this.heroes[i].name = this.highlight(this.formData.key_words, this.heroes[i].name)
            for (let j in this.heroes[i].spell) {
              this.heroes[i].spell[j].desc = this.highlight(this.formData.key_words, this.heroes[i].spell[j].desc)
            }
          }
        }
      },
      immediate: true,// 这个属性是重点啦
    },
  },
  methods: {
    showDrawer(platform, version, id, mainImg) {
      if (id === "") {
        return
      }
      this.drawer.show++
      this.drawer.isLoading = true
      this.drawer.background = mainImg

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
      return text.replace(new RegExp(`(${keywords})`, 'g'), `<em>$1</em>`)
    },
  },
  mounted() {
  }
}
</script>

<template>

  <p class="result-tips">{{ queryResult.tips }}</p>
  <div class="result-card" v-for="(item,i) in heroes" :key="i">
    <a-space direction="vertical">
      <a-card :hoverable="true"
              @click="showDrawer(item.platform,item.version,item.id,item.icon)"
      >
        <a-card-meta>
          <template #title>
            <span v-html="item.name"/>
          </template>
          <template #avatar>
            <img :src="item.icon" class="hero-icon" :alt="item.name"/>
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
              <span class="spell-icon">
                <img :src="s.icon" :alt="s.name"/>
              </span>
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
