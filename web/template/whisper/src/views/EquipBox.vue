<script setup>

// const condType = [
//   {label: '召唤师峡谷', value: '召唤师峡谷'},
//   {label: '嚎哭深渊', value: '嚎哭深渊'},
//   {label: '斗魂竞技场', value: '斗魂竞技场'},
// ];
</script>
<script>
import axios from 'axios';

export default {
  data() {
    return {
      formData: {
        platform: 0,
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
      result:[],
    }
  },
  watch: {
    'formData.keywords'(n) {
      // 使用 Axios 发起请求获取服务器数据
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
    }
  },
  methods: {
    filter(event) {
    }
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
          <a-checkbox-group v-model:value="formData.keywords">
            <a-descriptions v-for="(cate, index) in equipTypes" :title="cate.cate" :key="index">
              <a-descriptions-item >
                <a-checkbox v-for="(sub_cate, ii) in cate.sub_cate" :key="ii" :value="sub_cate.keywordsStr">
                  {{ sub_cate.name }}
                </a-checkbox>
              </a-descriptions-item>
            </a-descriptions>
          </a-checkbox-group>
        </a-form>

        <div class="equip-card" v-for="(item,i) in result" :key="i">
          <a-space direction="vertical">
            <a-card>
              <a-card-meta :title="item.name" description="-">
                <template #avatar>
                  <a-avatar :src="item.icon"/>
                </template>
              </a-card-meta>
              <div v-html="item.desc"></div>
            </a-card>
          </a-space>
        </div>
        <a-back-top/>
      </a-layout-content>
    </a-layout>
  </a-space>
</template>

<style>
</style>
