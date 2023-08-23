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
        types: [],
      },
      condTypes:[]
    }
  },
  methods: {
    search() {
      // 使用 Axios 发起请求获取服务器数据
      axios.post('/query', this.formData)
          .then(response => {
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
  created() {
  },
  mounted() {
    // 使用 Axios 发起请求获取服务器数据
    axios.get('/equip/types')
        .then(response => {
          // 将服务器返回的数据更新到组件的 serverData 字段
          this.condTypes = response.data.data.types;
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
        <a-descriptions title="筛选条件">
          <a-descriptions-item label="">
            <a-checkbox-group name="equip-type" v-model:value="formData.types" >
              <a-checkbox v-for="ct in condTypes" :value="ct" :key="ct">{{ ct }}</a-checkbox>
            </a-checkbox-group>
          </a-descriptions-item>

        </a-descriptions>

        <a-space direction="vertical" className="equip-card">
          <a-card title="Card" >
            <div>Card contentCard contentCard contentCard contentCard contentCard contentCard content</div>
            <div>Card content</div>
          </a-card>
          <a-card title="Card"  >
            <div>Card content</div>
            <div>Card conteCard contentCard contentCard contentCard contentnt</div>
          </a-card>
          <a-card title="Card"  >
            <div>Card contentCard contentCard contentCard contentCard contentCard contentCard content</div>
            <div>Card content</div>
          </a-card>
          <a-card title="Card"  >
            <div>Card content</div>
            <div>Card conteCard contentCard contentCard contentCard contentnt</div>
          </a-card>
        </a-space>
        <a-back-top/>
      </a-layout-content>
    </a-layout>
  </a-space>
</template>

<style>
</style>
