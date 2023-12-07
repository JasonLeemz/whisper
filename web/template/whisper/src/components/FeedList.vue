<script>

import {YoutubeOutlined} from "@ant-design/icons-vue";

export default {
  components: {YoutubeOutlined},
  props: {
    feedList: {
      show: 0,
      isLoading: false,
      data: [],
    },
  },
  data() {
    return {
      strategyList: {
        show: false,
        title: '大神攻略',
      },
      htmlContent:"",
    }
  },
  watch: {
    'feedList.show'() {
      this.strategyList.show = true
    },
    'feedList.isLoading'(isLoading) {
      this.strategyList.isLoading = isLoading
    },
  },
  methods: {
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
          this.htmlContent = "<iframe src="+schemeLink+" style='display:none' target='' ></iframe>"
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
  }
}
</script>

<template>
  <a-divider orientation="left">{{strategyList.title}}</a-divider>
  <a-empty
      v-if="feedList.data.length === 0"
      description="大神难产中..."/>
  <template v-if="feedList.isLoading">
    <a-skeleton active/>
  </template>

  <template v-if="!feedList.isLoading && feedList.data.length !== 0">
    <!-- 推荐列表 START-->
    <a-collapse activeKey="strategy-list">
      <a-collapse-panel key="strategy-list" header="点击链接跳转观看">
        <div v-for="(item,i) in feedList.data" :key="i" class="strategy-wrap">
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
    <div v-html="this.htmlContent"></div>
  </template>
</template>
