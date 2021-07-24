<template>
  <div
    class="
      HoneView
      flex flex-col
      bg-black
      h-full
      w-full
      max-h-full max-w-full
      min-h-full min-w-full
    "
  >
    <StatusBar class="StatusBar w-full h-11 left-0 top-0 flex-none"></StatusBar>
    <router-view class="bg-green-300 flex-grow" />
  </div>
</template>

<script lang="ts" >
import { Options, Vue } from "vue-class-component";
import StatusBar from "@/components/StatusBar.vue";

@Options({
  components: {
    StatusBar,
  },
})
export default class HoneView extends Vue {}
</script>

<style scoped>
.MainContant {
  /* @apply bg-black; */
}
</style>