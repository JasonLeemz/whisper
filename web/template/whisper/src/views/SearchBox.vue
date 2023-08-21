<script setup>
import { RightOutlined,InfoCircleOutlined } from '@ant-design/icons-vue';

const wayOptions = [
  { label: '按名字', value: 'name' },
  { label: '按介绍', value: 'description' },
];

const mapOptions = [
  { label: '召唤师峡谷', value: '召唤师峡谷' },
  { label: '嚎哭深渊', value: '嚎哭深渊' },
  { label: '斗魂竞技场', value: '斗魂竞技场' },
];
</script>
<script>
import axios from 'axios';

export default {
  data() {
    return {
      formData: {
        key_words: '',
        platform: '0',
        category:'lol_equipment',
        way:['name','description'],
        map:['召唤师峡谷'],
      },
      tips: '',
      lists: []
    }
  },
  methods: {
    search(event) {

      // 使用 Axios 发起请求获取服务器数据
      axios.post('http://127.0.0.1:8123/query', this.formData)
          .then(response => {
            // 将服务器返回的数据更新到组件的 serverData 字段
            if (response.data.data != null) {
              this.tips = response.data.data.tips;
              this.lists = response.data.data.lists;
            }
          })
          .catch(error => {
            console.error('Error fetching server data:', error);
          });

      // // 方法中的 `this` 指向当前活跃的组件实例
      // alert(`Hello ${this.name}!`)
      // // `event` 是 DOM 原生事件
      // if (event) {
      //   alert(event.target.tagName)
      // }
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
            :label-col="{ span: 8 }"
            :wrapper-col="{ span: 16 }"
            autocomplete="on"
            @finish="onFinish"
            @finishFailed="onFinishFailed"
        >
          <div>
            <a-input-search
                v-model:value="formData.key_words"
                placeholder="搜索…"
                enter-button
                allowClear="true"
                required
                @search="search"
            />
          </div>
          <div class="blank"></div>

          <a-radio-group v-model:value="formData.platform" name="radioGroup">
            <a-radio value="0">英雄联盟端游</a-radio>
            <a-radio value="1">英雄联盟手游</a-radio>
          </a-radio-group>
          <div class="blank"></div>
          <a-space direction="vertical">
          <a-typography-text>更多条件 <RightOutlined /> </a-typography-text>
          </a-space>
          <div class="blank"></div>
          <a-radio-group v-model:value="formData.category" name="radioGroup">
            <a-radio value="lol_equipment">装备</a-radio>
            <a-radio value="lol_heroes">英雄</a-radio>
            <a-radio value="lol_rune">符文</a-radio>
            <a-radio value="lol_skill">召唤师技能</a-radio>
          </a-radio-group>
          <div class="blank"></div>
          <a-checkbox-group v-model:value="formData.way" name="way" :options="wayOptions" />
          <div class="blank"></div>

          <a-checkbox-group v-model:value="formData.map" name="map" :options="mapOptions" />
          <div class="blank"></div>
        </a-form>

        <a-descriptions>
          <a-descriptions-item>{{ tips }}</a-descriptions-item>
        </a-descriptions>
        <a-timeline>
          <a-timeline-item v-for="item in lists" :key="item.id">
            <template #dot>
              <InfoCircleOutlined />
            </template>
            <h4 class="timeline-h4" >
              <a-image
                  :width="35"
                  :height="35"
                  v-bind:src="item.iconPath"
              />
              <span v-html="item.name"></span>
            </h4>

            <div>
              <a-tag v-for="tag in item.tags" :key="tag.id">{{ tag }}</a-tag>
            </div>
            <a-divider />

            <div class="mainText" v-html="item.description"></div>
          </a-timeline-item>
        </a-timeline>

        <a-back-top />
      </a-layout-content>
    </a-layout>
  </a-space>
</template>

