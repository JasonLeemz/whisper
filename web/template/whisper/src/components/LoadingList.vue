<script>
export default {
  props: {
    SkeletonState: {
      show: false,
      isLoading: false,
    },
  },
  data() {
    return {}
  },
}
</script>

<template>
  <a-skeleton active />
</template>