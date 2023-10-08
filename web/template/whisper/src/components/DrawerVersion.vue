<script>

export default {
  props: {
    versionResult: {
      show: 0,
      isLoading: false,
      title:'',
      data: Object,
    },
  },
  data() {
    return {
      sideDrawer: {
        show: false,
        title: 'sss',
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
    'versionResult.title'(text) {
      this.sideDrawer.title = text
    }
  },
  methods: {}
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
      <a-skeleton active/>
    </template>

    <template v-if="!sideDrawer.isLoading">
      <a-collapse>
        <a-collapse-panel v-for="(row,cate) in versionResult.data" :key="cate"
                          :header="cate"
                          v-show="row.msg === 'success'" >
          <div class="version-panel">
            <div v-for="(item,idx) in row.data.list" :key="idx" class="version-panel-row">
              <img :src="item.img_url" :alt="item.title" class="version-avatar">
              <span class="version-title">{{item.title}} | {{item.descirbe}}</span>
              <p class="version-attach-content">{{item.attach_content}}</p>
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
