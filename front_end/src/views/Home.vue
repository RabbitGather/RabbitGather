<template>
  <div
    class="
      EverythingWrapper
      flex flex-col
      items-start
      relative
      w-full
      h-full
      bg-white
    "
  >
    <div
      class="
        MainContainer
        RadarRadiusRuler
        flex flex-col
        justify-center
        items-center
        static
        w-full
        h-3/4
        left-0
        top-0
        flex-none
        order-none
        self-stretch
        bg-green-600
        flex-grow
      "
    >
      <div
        class="
          MainView
          bg-gray-500
          flex flex-col
          justify-center
          items-center
          static
          w-full
          h-auto
          left-0
          top-0
          flex-none
          order-none
          self-stretch
          flex-grow
        "
      >
        <div
          class="
            StatusBar
            static
            w-full
            h-11
            left-0
            top-0
            bg-gray-500
            flex-none
            order-none
            self-stretch
            flex-grow-0
          "
          style="box-shadow: 0px 4px 4px rgba(0, 0, 0, 0.25)"
        ></div>
        <div
          class="
            MainContact
            bg-red-500
            flex flex-col
            items-center
            static
            w-full
            h-auto
            left-0
            top-0
            overflow-hidden
            flex-none
            order-1
            self-stretch
            flex-grow
          "
        ></div>
      </div>
      <div
        class="
          ControlBox
          static
          w-full
          h-11
          left-0
          bottom-0
          bg-white
          flex-none
          order-1
          self-stretch
          flex-grow-0
        "
        style="box-shadow: 0px -4px 4px rgba(0, 0, 0, 0.25)"
      >
        11111144444444
      </div>
    </div>

    <div
      v-dragscroll.x
      class="
        bg-blue-500
        flex-row
        RadarRadiusRuler
        static
        w-full
        h-12
        left-0
        bottom-0
        overflow-hidden
        flex-none
        order-1
        self-stretch
        flex-grow-0
      "
    >
      <!-- style="width: 20px; height: 20px" -->

      <div class="absolute block left-0 right-0 bg-purple-600 m-auto w-5 h-5">
        <!-- 1111 -->
      </div>
      <div class="bg-yellow-400 h-full" style="width: 2000px">
        sadsd5a4d5a4sd5sa4d6a5sd
      </div>
    </div>
  </div>
</template>

<script lang="ts" >
import { Options, Vue } from "vue-class-component";
import SendArticleWithPositioning from "@/components/SendArticleWithPositioning.vue"; // @ is an alias to /src
import PerChat from "@/components/PeerChat.vue";
@Options({
  components: {
    // SendArticleWithPositioning,
    // PerChat,
  },
})
export default class Home extends Vue {}
</script>

<style scoped>
</style>