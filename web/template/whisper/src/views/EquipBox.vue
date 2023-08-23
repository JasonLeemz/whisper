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
      condTypes: []
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
            <a-checkbox-group name="equip-type" v-model:value="formData.types">
              <a-checkbox v-for="(ct, index) in condTypes" :value="ct" :key="index">{{ index }}</a-checkbox>
            </a-checkbox-group>
          </a-descriptions-item>
        </a-descriptions>

        <div class="equip-card">
          <a-space direction="vertical">
            <a-card>
              <a-card-meta title="蓝水晶" description="description">
                <template #avatar >
                  <a-avatar src="https://game.gtimg.cn/images/lgamem/act/lrlib/img/EquipIcons/lol_zfdj.png" />
                </template>
              </a-card-meta>
              <div>生命值 +125 最大生命值\n\n该装备为辅助英雄专用，携带时小兵和野怪的击败赏金会减少。队伍中多个同类装备无法同时生效。\n\n每30秒获得1次充能（最多3次）。攻击敌方小兵时，消耗1次充能并处决生命值低于65%的小兵，将给队友提供全额赏金，自己可获得额外65金币，并回复20-80生命值。\n自己不参与小兵赏金分配，可单独获得50%赏金，并会将自己击败小兵的赏金全额提供给最近的队友。\n从野怪获得的赏金减少50%。\n\n任务：此物品获取500金币后转换为山脉壁垒。\n\n</div>
              <div>description</div>
            </a-card>
            <a-card title="蓝水晶">
              <template #avatar>
                <img src="https://game.gtimg.cn/images/lgamem/act/lrlib/img/EquipIcons/lol_zfdj.png" />
              </template>
              <div></div>
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
