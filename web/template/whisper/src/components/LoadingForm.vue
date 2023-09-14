<script>
export default {
  props: {
    SkeletonState: {
      show: false,
      isLoading: false,
    },
  },
  data() {
    return {
      active:true,
      size:'small',
      shape:'circle',
      block:true,
    }
  },
}
</script>

<template>
  <div class="skeleton-row">
    <a-space>
      <a-skeleton-button :active="active" :size="size" :shape="shape" :block="block" />
      <a-skeleton-avatar :active="active" :size="size" :shape="shape" :block="block" />
      <a-skeleton-input style="width: 300px" :active="active" :size="size" :block="block" />
    </a-space>
  </div>
  <div class="skeleton-row">
    <a-space>
      <a-skeleton-button :active="active" :size="size" :shape="shape" :block="block" />
      <a-skeleton-avatar :active="active" :size="size" :shape="shape" :block="block" />
      <a-skeleton-input style="width: 200px" :active="active" :size="size" :block="block" />
    </a-space>
  </div>
  <div class="skeleton-row">
    <a-space>
      <a-skeleton-button :active="active" :size="size" :shape="shape" :block="block" />
      <a-skeleton-avatar :active="active" :size="size" :shape="shape" :block="block" />
      <a-skeleton-input style="width: 400px" :active="active" :size="size" :block="block" />
    </a-space>
  </div>
</template>

<style scoped>
.skeleton-row{
  .ant-skeleton{
    margin-top: 10px;
  }
}
</style>