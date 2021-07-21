<template>
  <div class="HoneView flex flex-col bg-black h-full w-full">
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
        self-stretch
        flex-grow-0
      "
    ></StatusBar>
    <router-view class="MainContant bg-green-300 flex-grow w-full" />
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