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
      result:[],
    }
  },
  watch: {
    'formData.keywords'(n) {
      // 使用 Axios 发起请求获取服务器数据
      if (this.formData.keywords.length === 0){
        return
      }
      axios.post('/equip/filter', this.formData)
          .then(response => {
            this.result = response.data.data
            // console.log(response)
            // 将服务器返回的数据更新到组件的 serverData 字段
            // if (response.data.data != null) {
            //   this.tips = response.data.data.tips;
            //   this.lists = response.data.data.lists;
            // }
          })
          .catch(error => {
            console.error('Error fetching server data:', error);
          });
    },
    'formData.platform'(n) {

      this.platformCN = n === '0'?'端游':'手游'
      // 使用 Axios 发起请求获取服务器数据
      if (this.formData.keywords.length === 0){
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
  },
  created() {
  },
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
              {{ drawer? '收起条件' : '展开条件' }}
            </a-button>
          </div>

          <a-checkbox-group v-model:value="formData.keywords" v-show="drawer">
            <a-descriptions v-for="(cate, index) in equipTypes" :title="cate.cate" :key="index">
              <a-descriptions-item >
                <a-checkbox v-for="(sub_cate, ii) in cate.sub_cate" :key="ii" :value="sub_cate.keywordsStr">
                  <a-tooltip placement="top" class="equip-box-tooltip">
                    <template #title>
                      <span>{{sub_cate.keywordsStr}}</span>
                    </template>
                    {{ sub_cate.name }}
                  </a-tooltip>
                </a-checkbox>
              </a-descriptions-item>
            </a-descriptions>
          </a-checkbox-group>
        </a-form>

        <div class="equip-card" v-for="(item,i) in result" :key="i">
          <a-space direction="vertical">
            <a-card :hoverable="true">
              <a-card-meta :title="item.name" :description="item.plaintext==''?'&nbsp;':item.plaintext">
                <template #avatar>
                  <a-avatar :src="item.icon"/>
                </template>
              </a-card-meta>
              <a-tag color="blue" size="small">价格 {{ item.price }}</a-tag>
              <a-tag class="platform-tag" color="warning">
                <template #icon>
                  <exclamation-circle-outlined />
                </template>
                {{platformCN}}
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
