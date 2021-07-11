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
          :radius="currentPoint"
          :magnification="6"
          class="MainContact"
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
    <!-- {{ currentPoint }} -->
    <RadarRadiusRuler
      @point-update="RulerPointUpdate"
      :min="1"
      :max="100"
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
export default class Home extends Vue {
  private currentPoint = 0;
  RulerPointUpdate(val: number) {
    // console.log("VAL: " + val);
    this.currentPoint = val;
  }
}
</script>

<style scoped>
</style>