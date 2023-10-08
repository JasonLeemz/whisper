<script>

export default {
  props: {
    versionResult: {
      show: 0,
      isLoading: false,
      data: {},
    },
  },
  data() {
    return {
      sideDrawer: {
        show: false,
        title: '',
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
      :title="sideDrawer.title"
      placement="right"
  >
    <template v-if="sideDrawer.isLoading">
      <a-skeleton active/>
    </template>

    <template>
      <a-collapse>
        <a-collapse-panel key="1" header="This is panel header 1">
          <p>1</p>
        </a-collapse-panel>
        <a-collapse-panel key="2" header="This is panel header 2">
          <p>2</p>
        </a-collapse-panel>
      </a-collapse>
    </template>
  </a-drawer>
</template>
