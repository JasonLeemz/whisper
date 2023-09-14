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
      block:false,
    }
  },
}
</script>

<template>
  <div class="skeleton" >
    <a-spin />
  </div>
</template>

<style scoped>
.skeleton{
  display: block;
  height: 150px;
  border-radius: 10px;
  background: rgba(0,0,0,.03);
  text-align: center;
  line-height: 150px;
}
</style>