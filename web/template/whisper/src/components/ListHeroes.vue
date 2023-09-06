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
        activeKey: ref(),
        panel: {
          foldAll: false,
          foldBtn: {
            top: ref(50)
          },
          panelKeys: [],
          mapEquipType: {
            'out': '出门装',
            'shoe': '鞋子',
            'core': '核心套件',
            'other': '可选装备',
          },
        },
        data: {},
        mapPos: {
          'all': '', // 手游目前没有分路数据
          'top': '上单',
          'mid': '中路',
          'bottom': 'AD',
          'support': '辅助',
          'jungle': '打野',
        }
      },
    }
  },
  watch: {
    'sideDrawer.data'(data) {
      for (let postypes in data.equips) {
        for (let type in data.equips[postypes]) {
          this.sideDrawer.panel.panelKeys.push(type + '-' + postypes)
        }
      }
    }
  },
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
    panelSwitcher() {
      if (this.sideDrawer.panel.foldAll) {
        this.sideDrawer.panel.foldAll = false
        this.sideDrawer.activeKey = ref(this.sideDrawer.panel.panelKeys)
      } else {
        this.sideDrawer.panel.foldAll = true
        this.sideDrawer.activeKey = ref([])
      }
    },
  },
  mounted() {
    this.panelSwitcher()
  },
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

    <a-button v-if="Object.keys(sideDrawer.data.equips).length !== 0"
              type="primary" size="small"
              @click="panelSwitcher"
              class="foldall-btn">
      {{ sideDrawer.panel.foldAll ? '展开所有' : '收起全部' }}
    </a-button>

    <a-empty
        v-if="Object.keys(sideDrawer.data.equips).length === 0"
        description="当前版本该英雄缺乏足够的样本数据"/>

    <template v-for="(equips,pos) in sideDrawer.data.equips" :key="pos">
      <!-- 端游-->
      <h4 v-if="sideDrawer.data.platform===0">{{ sideDrawer.mapPos[pos] ? sideDrawer.mapPos[pos] : pos }}</h4>

      <!-- 手游-->
      <a-popover placement="topLeft">
        <template #content>
          <h4 class="hero-suit popover-title">{{ pos }}</h4>
          <span class="hero-suit popover-content" v-html="sideDrawer.data.ext_info.recommend_reason[pos]"></span>
        </template>
        <h4 v-if="sideDrawer.data.platform===1" class="pos-title">
          {{ pos }}
          <h6>
            ({{ sideDrawer.data.ext_info.author_info[pos].name }})
          </h6>
        </h4>
      </a-popover>

      <a-collapse
          v-model:activeKey="sideDrawer.activeKey"
      >
        <template #expandIcon="{ isActive }">
          <CaretRightOutlined :rotate="isActive ? 90 : 0"/>
        </template>

        <a-collapse-panel v-for="(row,rowidx) in equips" :key="rowidx+'-'+pos"
                          :header="sideDrawer.panel.mapEquipType[rowidx]"
                          v-show="row.length > 0"
                          class="hero-drawer-panel">
          <sub v-if="sideDrawer.data.platform===0">
            <em>登场率</em>
            <em>胜率</em>
          </sub>
          <div v-for="(equips,equipsidx) in row" :key="equipsidx" class="equip-row">
            <div :class="sideDrawer.data.platform===0?'equip-item-wrap':'equip-item-wrap equip-item-wrap-lolm'">
              <span v-for="(equip,equipidx) in equips" :key="equipidx" class="equip-item">
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
            </div>
            <span v-if="sideDrawer.data.platform===0" class="data-statistics">
              <em>{{ equips[0].showrate / 100 }}%</em>
              <em>{{ equips[0].winrate / 100 }}%</em>
            </span>
          </div>
        </a-collapse-panel>
      </a-collapse>
    </template>
  </a-drawer>
</template>
