<template>
  <div class="MainContainer w-full h-full bg-blue-500 flex flex-col">
    <Map
      ref="Map"
      class="
        Map
        bg-purple-400
        flex flex-col
        items-center
        justify-center
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
    >
    </Map>
    <Ruler
      ref="Ruler"
      @RadiusUpdate="RadiusUpdate"
      class="
        bg-yellow-400
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
    ></Ruler>
  </div>
</template>

<script lang="ts" >
import { Options, Vue } from "vue-class-component";
import StatusBar from "@/components/Map.vue";
import Ruler from "@/components/Ruler.vue";
import Map from "@/components/Map.vue";
import store, { AllActionTypes } from "@/store";
import { UserSettings } from "@/store/app";

@Options({
  components: { StatusBar, Map, Ruler },
})
export default class MainContainer extends Vue {
  beforeCreate() {
    console.log("MainContainer beforeCreate");
    store
      .dispatch(AllActionTypes.APP.GetUserInfo)
      .then((userinfo: UserSettings) => {
        let maxRadius = userinfo.radaRadius.MaxRadius;
        let minRadius = userinfo.radaRadius.MinRadius;
        (this.$refs.Ruler as Ruler).Init(maxRadius, minRadius);
      });
  }

  CurrentRadius = 0;
  NewRadius(radius: number) {
    this.CurrentRadius = radius;
  }
  RadiusUpdate(newRadius: number) {
    this.UpdateMap(newRadius);
  }

  UpdateMap(newRadius: number) {
    // console.log("newRadius: ", newRadius);
    (this.$refs.Map as Map).UpdateRadius(newRadius);
  }
}
</script>

<style scoped>
</style>