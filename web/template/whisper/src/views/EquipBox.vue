<script setup>
import {ExclamationCircleOutlined} from '@ant-design/icons-vue';
</script>
<script>
import axios from 'axios';

export default {
  data() {
    return {
      formData: {
        platform: '0',
        keywords: [],
      },
      equipTypes: [
        {
          cate: '',
          sub_cate: [{
            name: '',
            keywordsStr: '',
            keywordsSlice: [],
          }],
        }
      ],
      platformCN: '端游',
      drawer: true,
      result: [],
      sideDrawer: false,
      sideDrawerData: {
        "current": {},
        "from": [],
        "into": [],
        "gapPriceFrom": 0,
      },
    }
  },
  watch: {
    'formData.keywords'() {
      // 使用 Axios 发起请求获取服务器数据
      if (this.formData.keywords.length === 0) {
        return
      }
      axios.post('/equip/filter', this.formData)
          .then(response => {
            this.result = response.data.data
          })
          .catch(error => {
            console.error('Error fetching server data:', error);
          });
    },
    'formData.platform'(n) {

      this.platformCN = n === '0' ? '端游' : '手游'
      // 使用 Axios 发起请求获取服务器数据
      if (this.formData.keywords.length === 0) {
        return
      }
      axios.post('/equip/filter', this.formData)
          .then(response => {
            this.result = response.data.data
          })
          .catch(error => {
            console.error('Error fetching server data:', error);
          });
    },
  },
  methods: {
    showDrawer(platform, version, id) {
      axios.post('/equip/roadmap', {
        'platform': platform,
        'version': version,
        'id': id,
      }).then(response => {
        this.sideDrawer = true
        this.sideDrawerData = response.data.data
      }).catch(error => {
        console.error('Error fetching server data:', error);
      });
    }
  },
  created() {
  },
  computed: {},
  mounted() {
    // 使用 Axios 发起请求获取服务器数据
    axios.get('/equip/types')
        .then(response => {
          // 将服务器返回的数据更新到组件的 serverData 字段
          this.equipTypes = response.data.data.types;
        })
        .catch(error => {
          console.error('Error fetching server data:', error);
        });
  }
}
</script>
<template>
  <a-space direction="vertical" :style="{ width: '100%' }" class="wrap">
    <a-layout>
      <a-layout-content>
        <a-form :model="formData" name="equip-box">
          <div class="platform-switch">
            <a-radio-group v-model:value="formData.platform" name="radioGroup">
              <a-radio value=0>端游</a-radio>
              <a-radio value=1>手游</a-radio>
            </a-radio-group>
            <a-button type="link" class="cond-drawer" size="small" @click="drawer = !drawer">
              {{ drawer ? '收起条件' : '展开条件' }}
            </a-button>
          </div>

          <a-checkbox-group v-model:value="formData.keywords" v-show="drawer">
            <a-descriptions v-for="(cate, index) in equipTypes" :title="cate.cate" :key="index">
              <a-descriptions-item>
                <a-checkbox v-for="(sub_cate, ii) in cate.sub_cate" :key="ii" :value="sub_cate.keywordsStr">
                  <a-tooltip placement="top" class="equip-box-tooltip">
                    <template #title>
                      <span>{{ sub_cate.keywordsStr }}</span>
                    </template>
                    {{ sub_cate.name }}
                  </a-tooltip>
                </a-checkbox>
              </a-descriptions-item>
            </a-descriptions>
          </a-checkbox-group>
        </a-form>
        <a-drawer
            v-model:open="sideDrawer"
            class="custom-class"
            root-class-name="root-class-name"
            title="合成路线"
            placement="right"
        >
          <div class="equip-roadmap equip-into">
            <template v-for="(equip ,index) in sideDrawerData['into']" :key="index">
              <a-popover placement="bottom" arrow-point-at-center>
                <template #content>
                  <div class="roadmap-item">
                    <span class="roadmap-item-title">
                      {{ equip.name }}
                    </span>
                    <a-tag>
                    <span class="roadmap-item-price">
                      价格:{{ equip.price }}
                    </span>
                    </a-tag>
                  </div>
                  <span v-html="equip.desc"></span>
                </template>
                <img :src="equip.icon">
              </a-popover>
            </template>
          </div>
          <a-divider>
            <a-popover placement="bottom" arrow-point-at-center>
              <template #content>
                <div class="roadmap-item">
                    <span class="roadmap-item-title">
                      {{ sideDrawerData['current'].name }}
                    </span>
                  <a-tag>
                    <span class="roadmap-item-price">
                      价格:{{ sideDrawerData['current'].price }}
                    </span>
                  </a-tag>
                </div>
                <span v-html="sideDrawerData['current'].desc"></span>
              </template>
              <img class="equip-roadmap equip-current" :src="sideDrawerData['current'].icon" alt="">
            </a-popover>

          </a-divider>
          <div class="equip-roadmap equip-from">
            <template v-for="(equip ,index) in sideDrawerData['from']" :key="index">
              <a-popover placement="bottom" arrow-point-at-center>
                <template #content>
                  <div class="roadmap-item">
                    <span class="roadmap-item-title">
                      {{ equip.name }}
                    </span>
                    <a-tag>
                    <span class="roadmap-item-price">
                      价格:{{ equip.price }}
                    </span>
                    </a-tag>
                  </div>
                  <span v-html="equip.desc"></span>
                </template>
                <img :src="equip.icon">
              </a-popover>
            </template>
          </div>

          <table class="roadmap-detail">
            <tr v-for="(equip ,index) in sideDrawerData['from']" :key="index">
              <td>
                <img :src="equip.icon" alt="">
                <span>{{ equip.name }}</span>
                <span>, 价格: <em>{{ equip.price }}</em> </span>
              </td>
            </tr>
            <tr>
              <td>
                总价: <em>{{ sideDrawerData['current'].price }}</em> , 补差价: <em>{{ sideDrawerData.gapPriceFrom }}</em>
              </td>
            </tr>
          </table>
        </a-drawer>

        <div class="equip-card" v-for="(item,i) in result" :key="i">
          <a-space direction="vertical">
            <a-card :hoverable="true" @click="showDrawer(item.platform,item.version,item.id)">
              <a-card-meta :title="item.name" :description="item.plaintext==''?'&nbsp;':item.plaintext">
                <template #avatar>
                  <a-avatar :src="item.icon"/>
                </template>
              </a-card-meta>
              <a-tag color="blue" size="small">价格 {{ item.price }}</a-tag>
              <a-tag class="platform-tag" color="warning">
                <template #icon>
                  <exclamation-circle-outlined/>
                </template>
                {{ platformCN }}
              </a-tag>
              <div class="card-desc" v-html="item.desc"></div>
            </a-card>
          </a-space>
        </div>
        <a-back-top/>
      </a-layout-content>
    </a-layout>
  </a-space>
</template>

<style scoped>
</style>
