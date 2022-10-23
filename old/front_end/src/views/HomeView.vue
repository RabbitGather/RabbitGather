<template>
  <div
    class="
      HoneView
      flex flex-col
      bg-black
      h-screen
      w-screen
      max-h-screen max-w-screen
      min-h-screen min-w-screen
    "
  >
    <StatusBar class="StatusBar w-full h-11 left-0 top-0 flex-none"></StatusBar>
    <router-view class="bg-green-300 flex-grow max-w-full" />
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