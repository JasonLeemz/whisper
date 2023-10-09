<script>

export default {
  props: {
    versionResult: {
      show: 0,
      isLoading: false,
      title:'',
      introduction:'',
      data: Object,
    },
  },
  data() {
    return {
      sideDrawer: {
        show: false,
      },
    }
  },
  watch: {
    'versionResult.show'() {
      this.sideDrawer.show = true
    },
    'versionResult.isLoading'(isLoading) {
      this.sideDrawer.isLoading = isLoading
    },
  },
  methods: {}
}
</script>

<template>
  <a-drawer
      v-model:open="sideDrawer.show"
      class="custom-class"
      root-class-name="root-class-name"
      :title="versionResult.title"
      placement="right"
  >

    <a-typography-text type="secondary" class="version-subtitle">{{versionResult.introduction}}</a-typography-text>

    <template v-if="sideDrawer.isLoading">
      <a-skeleton active/>
    </template>

    <template v-if="!sideDrawer.isLoading">
      <a-collapse>
        <a-collapse-panel v-for="(row,cate) in versionResult.data" :key="cate"
                          :header="cate"
                          v-show="row.msg === 'success'" >
          <div class="version-panel">
            <div v-for="(item,idx) in row.data.list" :key="idx" class="version-panel-row">
              <img :src="item.head_url == null?item.img_url:item.head_url" :alt="item.title" class="version-avatar">
              <span class="version-title">{{item.title}} <em>{{item.descirbe}}</em> <em>{{item.content}}</em></span>
              <p class="version-attach-content">{{item.attach_content}}</p>
              <div v-if="item.head_url != null && item.img_url !== item.head_url" class="version-skin-wrap">
                <a-image :src="item.img_url" alt="" class="version-skin" />
              </div>
              <div>
                <div v-for="(spell,listidx) in item.list" :key="listidx">
                  <div class="version-spell-list">
                    <img :src="spell.icon" :alt="spell.title" class="version-spell">
                    <div class="version-spell-desc">
                      {{spell.title}}
                      <br>
                      {{spell.content}}
                    </div>
                  </div>
                </div>
              </div>
              <a-divider />
            </div>
          </div>
        </a-collapse-panel>
      </a-collapse>
    </template>
  </a-drawer>
</template>
