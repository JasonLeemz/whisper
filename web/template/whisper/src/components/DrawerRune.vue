<script>

export default {
  props: {
    runeResult: {
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
    'runeResult.show'() {
      this.sideDrawer.show = true
    },
    'runeResult.isLoading'(isLoading) {
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
    <!-- 适配英雄 START-->
    <h4 class="equip-suit-hero">适配英雄</h4>
    <template v-for="(hero,i) in runeResult.data['suit_heroes']" :key="i">
      <a-tooltip :title="hero.name">
        <img class="equip-hero-avatar"
             :title="hero.name"
             :src="hero.icon"
             :alt="hero.name"
        />

      </a-tooltip>
    </template>
    <!-- 适配英雄 END-->
  </a-drawer>
</template>
