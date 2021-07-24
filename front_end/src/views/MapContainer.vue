<template>
  <div class="MapContainer bg-blue-500 flex flex-col">
    <Map ref="Map" class="Map bg-purple-400 w-full flex-grow flex flex-col min-h-0"> </Map>

    <Ruler
      ref="Ruler"
      @RadiusUpdate="RadiusUpdate"
      class="h-14 flex-none"
    ></Ruler>
  </div>
  <!-- </div> -->
</template>

<script lang="ts" >
import { Options, Vue } from "vue-class-component";
import StatusBar from "@/components/StatusBar.vue";
import Ruler from "@/components/Ruler.vue";
import Map from "@/components/Map.vue";
import store, { AllActionTypes } from "@/store";
import { UserSettings } from "@/store/app";

@Options({
  components: { StatusBar, Map, Ruler },
})
export default class MapContainer extends Vue {
  beforeCreate() {
    console.log("MapContainer beforeCreate");
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
    // console.log("newRadius: ", newRadius)
    (this.$refs.Map as Map).UpdateRadius(newRadius);
  }
}
</script>

<style scoped>
</style>