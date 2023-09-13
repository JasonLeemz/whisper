<script>
import {CaretRightOutlined} from "@ant-design/icons-vue";
import {ref} from "vue";

export default {
  components: {CaretRightOutlined},
  props: {
    heroResult: {
      show: 0,
      isLoading: false,
      data: {},
    }, // 父组件传递的数据
  },
  data() {
    return {
      sideDrawer: {
        show: false,
        isLoading: false,
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
            'rune': '符文',
            'skill': '召唤师技能',
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
    'heroResult.show'(){
      this.sideDrawer.show = true
    },
    'heroResult.isLoading'(isLoading){
      this.sideDrawer.isLoading = isLoading
    },
    'heroResult.data'(data) {
      for (let postypes in data.equips) {
        for (let type in data.equips[postypes]) {
          this.sideDrawer.panel.panelKeys.push(type + '-' + postypes)
        }
      }
    }
  },
  methods: {
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
  created() {
  },
  computed: {},
  mounted() {
    this.panelSwitcher()
  }
}
</script>

<template>
  <a-drawer
      v-model:open="sideDrawer.show"
      class="custom-class"
      root-class-name="root-class-name"
      :title="sideDrawer.title"
      placement="right"
  >
    <template v-if="sideDrawer.isLoading">
      <a-skeleton active />
    </template>

    <template v-if="!sideDrawer.isLoading">
      <!-- 展开所有/收起全部 START-->
      <a-button v-if="Object.keys(heroResult.data.equips).length !== 0"
                type="primary" size="small"
                @click="panelSwitcher"
                class="foldall-btn">
        {{ sideDrawer.panel.foldAll ? '展开所有' : '收起全部' }}
      </a-button>
      <!-- 展开所有/收起全部 END-->

      <a-empty
          v-if="Object.keys(heroResult.data.equips).length === 0"
          description="当前版本该英雄缺乏足够的样本数据"/>

      <template v-for="(equips,pos) in heroResult.data.equips" :key="pos">

        <!-- 端游:上单/打野... 手游:主播推荐 START-->
        <!-- 端游-->
        <h4 v-if="heroResult.data.platform===0">{{ sideDrawer.mapPos[pos] ? sideDrawer.mapPos[pos] : pos }}</h4>

        <!-- 手游-->
        <template v-if="heroResult.data.platform===1">
          <a-popover placement="topLeft">
            <template #content>
              <h4 class="hero-suit popover-title">{{ pos }}</h4>
              <span class="hero-suit popover-content" v-html="heroResult.data.ext_info.recommend_reason[pos]"></span>
            </template>
            <h4 class="pos-title">
              {{ pos }}
              <h6>
                ({{ heroResult.data.ext_info.author_info[pos].name }})
              </h6>
            </h4>
          </a-popover>
        </template>
        <!-- 端游:上单/打野... 手游:主播推荐 END-->

        <!-- //////////  //////////  //////////  //////////  ////////// -->

        <!-- Panel Wrap START-->
        <a-collapse v-model:activeKey="sideDrawer.activeKey">
          <template #expandIcon="{ isActive }">
            <CaretRightOutlined :rotate="isActive ? 90 : 0"/>
          </template>

          <a-collapse-panel v-for="(row,rowidx) in equips" :key="rowidx+'-'+pos"
                            :header="sideDrawer.panel.mapEquipType[rowidx]?sideDrawer.panel.mapEquipType[rowidx]:rowidx"
                            v-show="row.length > 0"
                            class="hero-drawer-panel">

            <!-- 胜率和登场率的标题... START-->
            <sub v-if="heroResult.data.platform===0">
              <em>登场率</em>
              <em>胜率</em>
            </sub>
            <!-- 胜率和登场率的标题... END-->

            <div v-for="(equips,equipsidx) in row" :key="equipsidx" class="equip-row">
              <!-- 每一项内容的行:出门装/鞋子... START-->
              <div class="equip-item-wrap" :class="heroResult.data.platform===0?'':'equip-item-wrap-lolm'">
                <span v-for="(equip,equipidx) in equips" :key="equipidx" class="equip-item">
                <!-- 具体的每一项 START-->
                  <a-popover placement="bottom" arrow-point-at-center>
                    <template #content>
                      <div class="roadmap-item">
                        <span class="roadmap-item-title">
                          {{ equip.name }}
                        </span>

                        <a-tag v-if="rowidx!=='skill'">
                          <span class="roadmap-item-price" v-if="rowidx!=='rune'">
                            价格:{{ equip.price }}
                          </span>
                          <span class="roadmap-item-rune" v-html="equip.plaintext" v-if="rowidx==='rune'"></span>
                        </a-tag>
                      </div>

                      <!-- 为了将desc内容中的一部分作为subtitle显示，这里做了显示冗余 START-->
                      <span v-html="equip.desc" v-if="rowidx!=='rune'"></span>
                      <!-- 为了将desc内容中的一部分作为subtitle显示，这里做了显示冗余 END-->
                      <span v-html="equip.plaintext" v-if="rowidx==='rune'" class="rune-desc"></span>
                      <span v-html="equip.desc" v-if="rowidx==='rune'" class="rune-desc-long"></span>
                    </template>
                    <img :src="equip.icon" :alt="equip.name" :class="rowidx==='rune'?'equip-icon-rune':'equip-icon'">
                  </a-popover>
                <!-- 具体的每一项 END-->
                </span>
              </div>

              <!-- 出场率/胜率 START-->
              <span v-if="heroResult.data.platform===0" class="data-statistics">
                <em>{{ equips[0].showrate / 100 }}%</em>
                <em>{{ equips[0].winrate / 100 }}%</em>
              </span>
              <!-- 出场率/胜率 END-->

              <!-- 每一项内容的行:出门装/鞋子... END-->
            </div>
          </a-collapse-panel>
        </a-collapse>
        <!-- Panel Wrap END-->

        <!-- //////////  //////////  //////////  //////////  ////////// -->


      </template>

    </template>
  </a-drawer>
</template>