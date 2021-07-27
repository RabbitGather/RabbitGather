<template>
  <div class="MapContainer bg-blue-500 flex flex-col">
    <Map
      ref="Map"
      :maxRadius="maxRadius"
      :minRadius="minRadius"
      class="Map bg-purple-400 flex-grow"
    >
    </Map>

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
// import { io, Socket, Manager } from "socket.io-client";

@Options({
  components: { StatusBar, Map, Ruler },
})
export default class MapContainer extends Vue {
  maxRadius = 0;
  minRadius = 0;
  beforeCreate() {
    let conn = new WebSocket("wss://api.meowalien.com/update_listener");

    conn.onclose = (evt: Event) => {
      console.log("onclose evt: ", evt);
    };
    conn.onmessage = (evt: Event) => {
      console.log("onmessage evt: ", evt);
    };

    conn.onerror = () => {};
    conn.onopen = (ev: Event) => {
      console.log("onopen ev: ", ev);
      conn.send("THIS_IS_TEST");
    };

    store
      .dispatch(AllActionTypes.APP.GetUserInfo)
      .then((userinfo: UserSettings) => {
        this.maxRadius = userinfo.radaRadius.MaxRadius;
        this.minRadius = userinfo.radaRadius.MinRadius;
        (this.$refs.Ruler as Ruler).Init(this.maxRadius, this.minRadius);
      });
  }

  CurrentRadius = 0;
  NewRadius(radius: number) {
    this.CurrentRadius = radius;
  }

  // will emit by Ruler
  RadiusUpdate(newRadius: number) {
    this.UpdateMap(newRadius);
  }

  NewRadishGape = 10;
  lastRadius = 0;
  UpdateMap(newRadius: number) {
    // ask radish from server if radius reach the gape
    if (newRadius > this.NewRadishGape) {
      // pull new Radish from server and add in the map componet.
      // let radishes = this.GetNewRadishFromServer(lastRadius, newRadius);
      // (this.$refs.Map as Map).PushNewRadishes(radishes);
    }
    (this.$refs.Map as Map).UpdateRadius(newRadius);
  }
}
</script>

<style scoped>
</style>