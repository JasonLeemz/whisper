<script setup>
import {RightOutlined, SearchOutlined} from '@ant-design/icons-vue';

const wayOptions = [
  {label: '按名字', value: 'name'},
  {label: '按介绍', value: 'description'},
];

const mapOptions = [
  {label: '召唤师峡谷', value: '召唤师峡谷'},
  {label: '嚎哭深渊', value: '嚎哭深渊'},
  {label: '斗魂竞技场', value: '斗魂竞技场'},
];
</script>
<script>
import axios from 'axios';
import {message} from 'ant-design-vue';
import ListSkeleton from "@/components/ListSkeleton.vue";
import ListEquip from "@/components/ListEquip.vue";
import ListHeroes from "@/components/ListHeroes.vue";
import ListRune from "@/components/ListRune.vue";
import ListSkill from "@/components/ListSkill.vue";

export default {
  emits: ['loadingEvent'],
  components: {
    ListSkeleton, ListEquip, ListHeroes, ListRune, ListSkill
  },
  data() {
    return {
      formData: {
        key_words: '',
        platform: '0',
        category: 'lol_equipment',
        way: ['name', 'description'],
        map: ['召唤师峡谷'],
        more_cond_show: true,
      },
      SkeletonState:{
        show: false,
        isLoading: false,
      },
      query: {
        tips: '',
        equip: {
          data: {},
        },
        hero: {
          data: {},
        },
        rune: {
          data: {},
        },
        skill: {
          data: {},
        },
      }
    }
  },
  methods: {
    search() {
      this.$emit('loadingEvent', 0)

      if (this.formData.key_words === '') {
        message.error({
          top: `100px`,
          duration: 2,
          maxCount: 3,
          content: '请输入查询内容',
        })
        return
      }

      this.SkeletonState.show = true
      this.$emit('loadingEvent', 30)
      // 使用 Axios 发起请求获取服务器数据
      axios.post('/query', this.formData)
          .then(response => {
            // 将服务器返回的数据更新到组件的 serverData 字段
            if (response.data.data != null) {
              if (this.formData.category === 'lol_equipment') {
                this.query.equip.data = response.data.data
              } else if (this.formData.category === 'lol_heroes') {
                this.query.hero.data = response.data.data
              } else if (this.formData.category === 'lol_rune') {
                this.query.rune.data = response.data.data
              } else if (this.formData.category === 'lol_skill') {
                this.query.skill.data = response.data.data
              }
              this.query.tips = response.data.data.tips
            }
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
  },
  mounted() {
  }
}
</script>


<template>
  <a-space direction="vertical" :style="{ width: '100%' }" class="wrap">
    <a-layout>
      <a-layout-content>
        <a-form
            :model="formData"
            name="search-box"
        >
          <div>
            <a-space-compact block>
              <a-input
                  name="search-input"
                  v-model:value="formData.key_words"
                  placeholder="搜索…"
                  allowClear
                  @pressEnter="search"
              />
              <a-select
                  ref="select"
                  v-model:value="formData.platform"
              >
                <a-select-option value="0">端游</a-select-option>
                <a-select-option value="1">手游</a-select-option>
              </a-select>
              <a-button type="primary" @click="search">
                <SearchOutlined/>
              </a-button>
            </a-space-compact>
          </div>
          <div class="blank"></div>
          <a-radio-group v-model:value="formData.category" name="radioGroup">
            <a-radio value="lol_equipment">装备</a-radio>
            <a-radio value="lol_heroes">英雄</a-radio>
            <a-radio value="lol_rune">符文</a-radio>
            <a-radio value="lol_skill">召唤师技能</a-radio>
          </a-radio-group>
          <div class="blank"></div>
          <div class="blank"></div>
          <a-space @click="formData.more_cond_show = !formData.more_cond_show" direction="vertical">
            <a-typography-text>更多条件
              <RightOutlined :rotate="formData.more_cond_show ? 90 : 0"/>
            </a-typography-text>
          </a-space>
          <Transition>
            <div v-show="formData.more_cond_show">
              <div>
                <a-checkbox-group v-model:value="formData.way" name="way" :options="wayOptions"/>
              </div>
              <div>
                <a-checkbox-group v-model:value="formData.map" name="map" :options="mapOptions"/>
              </div>
            </div>
          </Transition>
        </a-form>

        <ListSkeleton
            v-if="SkeletonState.show"
            :skeleton-state="SkeletonState"/>

        <template v-if="!SkeletonState.show">
          <ListEquip
              v-if="formData.category==='lol_equipment'"
              :query-result="query.equip.data"
          />

          <ListHeroes
              v-if="formData.category==='lol_heroes'"
              :query-result="query.hero.data"
          />

          <ListRune
              v-if="formData.category==='lol_rune'"
              :query-result="query.rune.data"
          />

          <ListSkill
              v-if="formData.category==='lol_skill'"
              :query-result="query.skill.data"
          />
        </template>
        <a-back-top/>
      </a-layout-content>
    </a-layout>
  </a-space>
</template>

<style scoped>
</style>