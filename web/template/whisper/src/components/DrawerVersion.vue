<script>

export default {
  props: {
    skillResult: {
      show: 0,
      isLoading: false,
      data: {},
    },
  },
  data() {
    return {
      sideDrawer: {
        show: false,
        title: '适配英雄',
      },
    }
  },
  watch: {
    'skillResult.show'() {
      this.sideDrawer.show = true
    },
    'skillResult.isLoading'(isLoading) {
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
      :title="sideDrawer.title"
      placement="right"
  >
    <template v-if="sideDrawer.isLoading">
      <a-skeleton active/>
    </template>

    <a-empty
        v-if="!sideDrawer.isLoading && Object.keys(skillResult.data).length === 0"
        description="当前版本该符文缺乏足够的样本数据"/>

    <!-- 适配英雄 START-->
    <template v-if="!sideDrawer.isLoading && Object.keys(skillResult.data).length > 0">
      <h4 class="equip-suit-hero">适配英雄</h4>
      <template v-for="(hero,i) in skillResult.data" :key="i">
        <a-tooltip :title="hero.name">
          <img class="equip-hero-avatar"
               :title="hero.name"
               :src="hero.icon"
               :alt="hero.name"
          />

        </a-tooltip>
      </template>
    </template>
    <!-- 适配英雄 END-->
  </a-drawer>
</template>
