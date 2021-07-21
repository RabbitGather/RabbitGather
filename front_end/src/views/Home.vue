<template>
  <div
    ref="HomeView"
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
          :radius="CircleRadiusPercentage"
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
    <RadarRadiusRuler
      @point-update="SearchRadiusUpdate"
      :min="MinRadius"
      :max="MaxRadius"
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
import axios from "axios";
import { useStore } from "@/store";
import { UserSetting, AuthGetterTypes } from "@/store/auth/getters";
const store = useStore();

@Options({
  components: {
    RadarRadiusRuler,
    MainContact,
    StatusBar,
    PeerChat,
    ControlBox,
  },
})
/*
view 從後端拉資料，將資料塞到組件內顯示
*/
export default class Home extends Vue {
  // bind on MainContact
  CircleRadiusPercentage = 0;

  // bind on RadarRadiusRuler
  MaxRadius = 0;
  MinRadius = 0;
  UserName = 0;
  // homeView!: HTMLDivElement;
  done = false;
  beforeMount() {
    // console.log(
    //   "AuthGetterTypes.isAuthenticated: ",
    //   useStore().getters[AuthGetterTypes.isAuthenticated]
    // );
    // 取得使用者權限相關設定
    let setting: UserSetting = store.getters[AuthGetterTypes.UserSetting];
    this.MaxRadius = setting.RadiusRange.Max;
    this.MinRadius = setting.RadiusRange.Min;
  }
  mounted() {
    console.log("--- mounted ---");
    // this.homeView = this.$refs.HomeView as HTMLDivElement;
    // console.log("this.HomeView: ", this.homeView);
    this.done = true;
  }

  SearchRadiusUpdate(newRadius: number) {
    if (!this.done) {
      return;
    }
    console.log("SearchRadiusUpdate");
    this.UpdateCircleRadius((newRadius / this.MaxRadius) * 100);
  }

  UpdateCircleRadius(newRadiusPercentage: number) {
    if (!this.done) {
      return;
    }
    console.log("New Radius: ", newRadiusPercentage);
    this.CircleRadiusPercentage = newRadiusPercentage;
  }
}
</script>

<style scoped>
</style>