<script setup>
import {InfoCircleOutlined, RightOutlined, SearchOutlined} from '@ant-design/icons-vue';

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
import DrawerEquip from "@/components/DrawerEquip.vue";

export default {
  components: {
    DrawerEquip
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
      query:{
        tips: '',
        lists: [],
      },
      sideDrawer: {
        show: false,
        title: '合成路线',
        data: {
          "current": {},
          "from": [],
          "into": [],
          "gapPriceFrom": 0,
        },
      },
    }
  },
  methods: {
    search() {
      if (this.formData.key_words === '') {
        message.error({
          top: `100px`,
          duration: 2,
          maxCount: 3,
          content: '请输入查询内容',
        })

        return
      }
      // 使用 Axios 发起请求获取服务器数据
      axios.post('/query', this.formData)
          .then(response => {
            // 将服务器返回的数据更新到组件的 serverData 字段
            if (response.data.data != null) {
              this.query.tips = response.data.data.tips;
              this.query.lists = response.data.data.lists;
            }
          })
          .catch(error => {
            console.error('Error fetching server data:', error);
          });
    },
    showDrawer(platform, version, id) {
      if (id === "") {
        return
      }
      if (this.formData.category === "lol_equipment") {
        axios.post('/equip/roadmap', {
          'platform': platform,
          'version': version,
          'id': id,
        }).then(response => {
          this.sideDrawer.show = true
          this.sideDrawer.data = response.data.data
        }).catch(error => {
          console.error('Error fetching server data:', error);
        });
      } else if (this.formData.category === "lol_heroes") {
        axios.post('/hero/suit', {
          'platform': platform,
          'hero_id': id,
        }).then(response => {
          this.sideDrawer.show = true
          console.log(response.data.data)
          // this.sideDrawerData = response.data.data
        }).catch(error => {
          console.error('Error fetching server data:', error);
        });
      } else if (this.formData.category === "lol_rune") {
        // ...
      } else if (this.formData.category === "lol_skill") {
        // ...
      }

    }
  },
  created() {
  },
  mounted() {
  }
}
</script>


<template>
  <DrawerEquip
      :side-drawer="sideDrawer"
  />

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
          <Transition name="fade">
            <div v-show="formData.more_cond_show" >
              <div>
                <a-checkbox-group v-model:value="formData.way" name="way" :options="wayOptions"/>
              </div>
              <div>
                <a-checkbox-group v-model:value="formData.map" name="map" :options="mapOptions"/>
              </div>
            </div>
          </Transition>

        </a-form>

        <a-descriptions>
          <a-descriptions-item>{{ query.tips }}</a-descriptions-item>
        </a-descriptions>
        <a-timeline>
          <a-timeline-item v-for="item in query.lists" :key="item.id" class="ant-card-hoverable">
            <template #dot>
              <InfoCircleOutlined :style="{fontSize: '16px'}"/>
            </template>
            <h4 class="timeline-h4" @click="showDrawer(item.platform,item.version,item.id)">
              <img v-bind:src="item.iconPath"/>
              <span v-html="item.name"></span>
            </h4>

            <div>
              <a-tag v-for="tag in item.tags" :key="tag.id" color="blue">{{ tag }}</a-tag>
            </div>
            <a-divider/>

            <div class="mainText" v-html="item.description"></div>
          </a-timeline-item>
        </a-timeline>

        <a-back-top/>
      </a-layout-content>
    </a-layout>
  </a-space>
</template>

<style scoped>
</style>