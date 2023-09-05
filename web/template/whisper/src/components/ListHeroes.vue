<script>
import {CaretRightOutlined, ExclamationCircleOutlined} from "@ant-design/icons-vue";
import axios from "axios";
import {ref} from 'vue';

export default {
  components: {ExclamationCircleOutlined, CaretRightOutlined},
  props: {
    queryResult: Object, // 父组件传递的数据类型
  },
  data() {
    return {
      sideDrawer: {
        show: false,
        title: '推荐出装',
        activeKey: ref(['out', 'shoe', 'core','other']),
        panel:{
          foldAll: false,
          foldBtn :{
            top: ref(50)
          },
        },
        data: {},
        mapPos:{
          'top':'上单',
          'mid':'中路',
          'bottom':'AD',
          'support':'辅助',
          'jungle':'打野',
        }
      },
    }
  },
  watch: {},
  methods: {
    showDrawer(platform, version, id) {
      if (id === "") {
        return
      }

      axios.post('/hero/suit', {
        'platform': platform,
        'hero_id': id,
      }).then(response => {
        this.sideDrawer.show = true
        this.sideDrawer.data = response.data.data
      }).catch(error => {
        console.error('Error fetching server data:', error);
      });
    },
    panelSwitcher(){
      if (this.sideDrawer.panel.foldAll){
        this.sideDrawer.panel.foldAll = false
        this.sideDrawer.activeKey = ref(['out', 'shoe', 'core','other'])
      }else{
        this.sideDrawer.panel.foldAll = true
        this.sideDrawer.activeKey = ref([])
      }
    },
  }
}
</script>

<template>
  <div class="result-card" v-for="(item,i) in queryResult.list" :key="i">
    <a-space direction="vertical">
      <a-card :hoverable="true" @click="showDrawer(item.platform,item.version,item.id)">
        <a-card-meta :title="item.name">
          <template #avatar>
            <a-avatar :src="item.icon" class="hero-icon"/>
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
              <img :src="s.icon" class="spell-icon"/>
              <h6>{{ s.name }}</h6>
              <span>{{ s.sort }}</span>
              <div v-html="s.desc" class="spell-desc"></div>
            </li>
          </ul>
        </div>
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

    <a-button size="small" @click="panelSwitcher" class="foldall-btn">{{sideDrawer.panel.foldAll?'展开所有':'收起全部'}}</a-button>

    <template v-for="(equips,pos) in sideDrawer.data" :key="pos">
      <h4>{{ sideDrawer.mapPos[pos] }}</h4>
      <a-collapse
          v-model:activeKey="sideDrawer.activeKey"
      >
        <template #expandIcon="{ isActive }">
          <CaretRightOutlined :rotate="isActive ? 90 : 0"/>
        </template>
        <a-collapse-panel key="out" header="出门装" class="hero-drawer-panel">
          <sub>
            <em>登场率</em>
            <em>胜率</em>
          </sub>
          <div v-for="(row,rowidx) in equips.out" :key="rowidx" class="equip-row">
            <span v-for="(equip,equipidx) in row" :key="equipidx" class="equip-item">
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
                <img :src="equip.icon" :alt="equip.name" class="equip-icon">
              </a-popover>
            </span>
            <span class="data-statistics">
              <em>{{row[0].showrate/100}}%</em>
              <em>{{row[0].winrate/100}}%</em>
            </span>
          </div>
        </a-collapse-panel>

        <a-collapse-panel key="core" header="核心三件套" class="hero-drawer-panel">
          <sub>
            <em>登场率</em>
            <em>胜率</em>
          </sub>
          <div v-for="(row,rowidx) in equips.core" :key="rowidx" class="equip-row">
            <span v-for="(equip,equipidx) in row" :key="equipidx" class="equip-item">
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
                <img :src="equip.icon" :alt="equip.name" class="equip-icon">
              </a-popover>
            </span>
            <span class="data-statistics">
              <em>{{row[0].showrate/100}}%</em>
              <em>{{row[0].winrate/100}}%</em>
            </span>
          </div>
        </a-collapse-panel>

        <a-collapse-panel key="shoe" header="鞋子" class="hero-drawer-panel">
          <div class="equip-row">
            <span v-for="(equip,rowidx) in equips.shoe" :key="rowidx" class="equip-item">
              <a-popover placement="bottom" arrow-point-at-center>
                <template #content>
                  <div class="roadmap-item">
                          <span class="roadmap-item-title">
                            {{ equip.name }}
                          </span>
                    <a-tag>
                          <span class="roadmap-item-price">
                            胜率:{{equip.winrate/100}}%
                          </span>
                    </a-tag>
                    <a-tag>
                          <span class="roadmap-item-price">
                            价格:{{ equip.price }}
                          </span>
                    </a-tag>
                  </div>
                  <span v-html="equip.desc"></span>
                </template>
                <img :src="equip.icon" :alt="equip.name" class="equip-icon">
              </a-popover>
            </span>
          </div>
        </a-collapse-panel>

        <a-collapse-panel key="other" header="其他装备" class="hero-drawer-panel">
          <div class="equip-row">
            <span v-for="(equip,rowidx) in equips.other" :key="rowidx" class="equip-item">
              <a-popover placement="bottom" arrow-point-at-center>
                <template #content>
                  <div class="roadmap-item">
                          <span class="roadmap-item-title">
                            {{ equip.name }}
                          </span>
                    <a-tag>
                          <span class="roadmap-item-price">
                            胜率:{{equip.winrate/100}}%
                          </span>
                    </a-tag>
                    <a-tag>
                          <span class="roadmap-item-price">
                            价格:{{ equip.price }}
                          </span>
                    </a-tag>
                  </div>
                  <span v-html="equip.desc"></span>
                </template>
                <img :src="equip.icon" :alt="equip.name" class="equip-icon">

              </a-popover>
            </span>
          </div>
        </a-collapse-panel>
      </a-collapse>
    </template>
  </a-drawer>
</template>
