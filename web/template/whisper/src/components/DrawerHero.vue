<script>
import {YoutubeOutlined,CaretRightOutlined} from "@ant-design/icons-vue";
import {ref} from "vue";

export default {
  components: {YoutubeOutlined,CaretRightOutlined},
  props: {
    heroResult: {
      show: 0,
      isLoading: false,
      background: '',
      data: {},
    }, // 父组件传递的数据
  },
  data() {
    return {
      loadingImgs: new Map(),
      sideDrawer: {
        show: false,
        isLoading: false,
        title: '推荐出装',
        activeKey: ref(),
        panel: {
          foldAll: true,
          foldBtn: {
            top: ref(50)
          },
          panelKeys: [],
          panelKeysStrategyList: ["strategy-list"],
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
        },
        htmlContent:"",
      },
    }
  },
  watch: {
    'heroResult.show'() {
      this.sideDrawer.show = true
    },
    'heroResult.isLoading'(isLoading) {
      this.sideDrawer.isLoading = isLoading
    },
    'heroResult.data.result'(data) {
      for (let postypes in data.equips) {
        for (let type in data.equips[postypes]) {
          if (type === 'rune' || type === 'core') {
            this.sideDrawer.panel.panelKeys.push(type + '-' + postypes)
          }
        }
      }
    }
  },
  methods: {
    stopLoading(idx) {
      this.loadingImgs.set("load_" + idx, false);
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
    jumpApp(source,jump_url,video_id){
      let ua = navigator.userAgent.toLowerCase();
      let isAndroid = ua.indexOf('android') > -1 || ua.indexOf('linux') > -1;
      let mobileAgent = ["iphone", "ipod", "ipad", "android", "mobile", "blackberry", "webos", "incognito", "webmate", "bada", "nokia", "lg", "ucweb", "skyfire"];
      let isMobile = false;
      for (let i = 0; i < mobileAgent.length; i++) {
        if (ua.indexOf(mobileAgent[i]) !== -1)
        {
          isMobile = true;
          break;
        }
      }

      if (source === 0){
        let h5Link = jump_url
        let schemeLink = "bilibili://video/"+video_id

        if(isAndroid){
          //android
          this.sideDrawer.htmlContent = "<iframe src="+schemeLink+" style='display:none' target='' ></iframe>"
          setTimeout(function(){
            // window.location = h5Link
            window.open(h5Link, "_blank");
          },600);
        }else if (!isMobile){
          window.open(h5Link, "_blank");
        }else{
          //ios
          window.location = schemeLink;
          setTimeout(function(){
            window.location = h5Link
            // window.open(h5Link, "_blank");
          },25);
        }
      }
    }
  },
  created() {
  },
  computed: {},
  mounted() {
    this.panelSwitcher()
  }
}
</script>
<!--v-bind:style="{'background-image': 'linear-gradient(to bottom, rgba(255,255,255,1),rgba(255,255,255,1), rgba(255,255,255,1),rgba(255,255,255,.5)),url('+heroResult.background+')'}"-->

<template>
  <a-drawer
      v-model:open="sideDrawer.show"
      :title="sideDrawer.title"
      placement="right"
      width="100%"
  >
    <template v-if="sideDrawer.isLoading">
      <a-skeleton active/>
    </template>

    <template v-if="!sideDrawer.isLoading">
      <!-- 展开所有/收起全部 START-->
            <a-button v-if="Object.keys(heroResult.data.result.equips).length !== 0"
                      type="primary" size="small"
                      @click="panelSwitcher"
                      class="foldall-btn">
              {{ sideDrawer.panel.foldAll ? '展开常用' : '收起常用' }}
            </a-button>
      <!-- 展开所有/收起全部 END-->

      <a-empty
          v-if="Object.keys(heroResult.data.result.equips).length === 0"
          description="当前版本该英雄缺乏足够的样本数据"/>

      <template v-for="(equips,pos) in heroResult.data.result.equips" :key="pos">

        <!-- 端游:上单/打野... 手游:主播推荐 START-->
        <!-- 端游-->
        <h4 v-if="heroResult.data.result.platform===0">{{ sideDrawer.mapPos[pos] ? sideDrawer.mapPos[pos] : pos }}</h4>

        <!-- 手游-->
        <template v-if="heroResult.data.result.platform===1">
          <a-popover placement="topLeft">
            <template #content>
              <h4 class="hero-suit popover-title">{{ pos }}</h4>
              <span class="hero-suit popover-content" v-html="heroResult.data.result.ext_info.recommend_reason[pos]"></span>
            </template>
            <h4 class="pos-title">
              {{ pos }}
              <h6>
                ({{ heroResult.data.result.ext_info.author_info[pos].name }})
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
            <sub v-if="heroResult.data.result.platform===0">
              <em>登场率</em>
              <em>胜率</em>
            </sub>
            <!-- 胜率和登场率的标题... END-->

            <div v-for="(equips,equipsidx) in row" :key="equipsidx" class="equip-row">
              <!-- 每一项内容的行:出门装/鞋子... START-->
              <div class="equip-item-wrap" :class="heroResult.data.result.platform===0?'':'equip-item-wrap-lolm'">
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
                    <div class="rune-desc-img">
                      <em class="rune-type" v-if="equipidx===0 && equip.rune_type !== ''">{{ equip.rune_type }}</em>
                      <div :class="rowidx==='rune'?'rune-desc-img-wrap':''" >
                        <img :src="equip.icon" :alt="equip.name" :class="rowidx==='rune'?'equip-icon-rune':'equip-icon'">
                        <em class="rune-name" v-if="rowidx==='rune'">{{ equip.name }}</em>
                      </div>
                    </div>
                  </a-popover>
                  <!-- 具体的每一项 END-->
                </span>
              </div>

              <!-- 出场率/胜率 START-->
              <span v-if="heroResult.data.result.platform===0" class="data-statistics">
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

    <a-divider orientation="left">大神攻略</a-divider>
    <a-empty
        v-if="Object.keys(heroResult.data.feed).length === 0"
        description="大神难产中..."/>
    <template v-if="!sideDrawer.isLoading && Object.keys(heroResult.data.feed).length > 0">
      <!-- 推荐列表 START-->
      <a-collapse  v-model:activeKey="sideDrawer.panel.panelKeysStrategyList">
        <a-collapse-panel key="strategy-list" header="点击链接跳转观看">
          <div v-for="(item,i) in heroResult.data.feed" :key="i" class="strategy-wrap">
            <span class="jump-wrap" @click="jumpApp(item.source,item.jump_url,item.video_id)">
              <div class="strategy-main-img-wrap">
                <img :src="item.main_img" :alt="item.title" class="strategy-main-img" />
                <div class="card-state">
                  <i class="play-times"><YoutubeOutlined /> {{ item.played }}</i>
                  <i class="play-length">{{ item.length }}</i>
                </div>
              </div>
              <p class="strategy-title">{{ item.title }}</p>
              <p class="strategy-subtitle">{{ item.subtitle }}</p>
              <p class="strategy-author">
                {{ item.author }}
                <i class="strategy-public-date">{{ item.public_date }}</i>
              </p>

            </span>
          </div>
        </a-collapse-panel>
      </a-collapse>
      <!-- 推荐列表 END-->
    </template>
    <template>
      <div v-html="this.sideDrawer.htmlContent"></div>
    </template>
  </a-drawer>
</template>

<style scoped>
</style>