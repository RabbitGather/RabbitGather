<template>
  <div
    class="EverythingWrapper flex flex-col items-start relative w-full h-full"
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
          MainView
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
        <MainContact
          class="
            MainContact
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
        ></MainContact>
      </div>
      <ControlBox
        class="
          ControlBox
          static
          w-full
          h-11
          left-0
          bottom-0
          shadow-up
          flex-none
          order-1
          self-stretch
          flex-grow-0
        "
      >
      </ControlBox>
    </div>

    <RadarRadiusRuler
      min=1
      max=100
    
      class="
        RadarRadiusRuler
        flex-row
        static
        w-full
        h-14
        left-0
        bottom-0
        overflow-hidden
        flex-none
        order-1
        self-stretch
        flex-grow-0
      "
    >
    </RadarRadiusRuler>
  </div>
</template>

<script lang="ts" >
import { Options, Vue } from "vue-class-component";
import RadarRadiusRuler from "@/components/RadarRadiusRuler.vue";
import PeerChat from "@/components/PeerChat.vue";
import MainContact from "@/components/MainContact.vue";
import StatusBar from "@/components/StatusBar.vue";
import ControlBox from "@/components/ControlBox.vue";
@Options({
  components: {
    RadarRadiusRuler,
    MainContact,
    StatusBar,
    PeerChat,
    ControlBox,
  },
})
export default class Home extends Vue {}
</script>

<style scoped>
</style>