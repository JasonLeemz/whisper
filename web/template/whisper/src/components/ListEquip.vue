<script>
import {ExclamationCircleOutlined} from "@ant-design/icons-vue";
import axios from "axios";

export default {
  components: {ExclamationCircleOutlined},
  props: {
    queryResult: Object, // 父组件传递的数据类型
  },
  data() {
    return {
      sideDrawer: {
        show: false,
        title: '合成路线',
        data: {
          "current": {},
          "from": [],
          "into": [],
          "gapPriceFrom": 0,
          "suit_heroes": [],
        },
      },
    }
  },
  watch: {},
  methods: {
    showDrawer(platform, version, id) {
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
    },
  }
}
</script>

<template>
  <a-descriptions>
    <a-descriptions-item>{{ queryResult.tips }}</a-descriptions-item>
  </a-descriptions>
  <div class="result-card" v-for="(item,i) in queryResult.list" :key="i">
    <a-space direction="vertical">
      <a-card :hoverable="true" @click="showDrawer(item.platform,item.version,item.id)">
        <a-card-meta :title="item.name">
          <template #avatar>
            <a-avatar :src="item.icon"/>
          </template>
        </a-card-meta>
        <div class="ant-tag-wrap">
          <a-tag v-for="tag in item.tags" :key="tag.id" color="blue" :title="tag">{{ tag }}</a-tag>
        </div>
        <a-tag class="platform-tag" color="warning">
          <template #icon>
            <ExclamationCircleOutlined/>
          </template>
          {{ item.platform === 0 ? '端游' : '手游' }}
        </a-tag>

        <div class="mainText" v-html="item.desc"></div>
      </a-card>
    </a-space>
  </div>

  <a-drawer
      v-model:open="sideDrawer.show"
      class="custom-class"
      root-class-name="root-class-name"
      :title="sideDrawer.title"
      placement="right"
  >
    <div class="equip-roadmap equip-into">
      <template v-for="(equip ,index) in sideDrawer.data['into']" :key="index">
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
                      {{ sideDrawer.data['current'].name }}
                    </span>
            <a-tag>
                    <span class="roadmap-item-price">
                      价格:{{ sideDrawer.data['current'].price }}
                    </span>
            </a-tag>
          </div>
          <span v-html="sideDrawer.data['current'].desc"></span>
        </template>
        <img class="equip-roadmap equip-current" :src="sideDrawer.data['current'].icon" alt="">
      </a-popover>
    </a-divider>
    <div class="equip-roadmap equip-from">
      <template v-for="(equip ,index) in sideDrawer.data['from']" :key="index">
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
      <tr v-for="(equip ,index) in sideDrawer.data['from']" :key="index">
        <td>
          <img :src="equip.icon" alt="">
          <span>{{ equip.name }}</span>
          <span>, 价格: <em>{{ equip.price }}</em> </span>
        </td>
      </tr>
      <tr>
        <td>
          总价: <em>{{ sideDrawer.data['current'].price }}</em> , 合成所需: <em>{{ sideDrawer.data.gapPriceFrom }}</em>
        </td>
      </tr>
    </table>

    <h4 class="equip-suit-hero">适配英雄</h4>
    <template v-for="(hero,i) in sideDrawer.data['suit_heroes']" :key="i">
      <a-tooltip :title="hero.name">
        <img class="equip-hero-avatar"
             :src="hero.icon"
             :alt="hero.name"
        />
      </a-tooltip>
    </template>

  </a-drawer>
</template>
