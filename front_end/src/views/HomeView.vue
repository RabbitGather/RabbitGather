<template>
  <div
    class="HoneView flex flex-col bg-black h-full w-full max-h-full max-w-full"
  >
    <StatusBar
      class="
        StatusBar
        static
        w-full
        h-11
        left-0
        top-0
        flex-none
        order-none
        flex-shrink-0
      "
    ></StatusBar>
    <router-view class="bg-green-300 flex-grow max-w-full max-h-full" />
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