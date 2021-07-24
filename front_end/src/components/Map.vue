<template>
  <div ref="MapDivElement" id="Map" class="">
    <!-- <Circle class="Circle" ref="Circle"></Circle> -->
    <div
      id="openmap"
      class="bg-red-500 w-full flex-grow overflow-scroll"
      style=""
    >
      <div class="bg-yellow-500" style="width: 2000px; height: 2000px"></div>
    </div>
  </div>
</template>

<script lang="ts">
import { Options, Vue } from "vue-class-component";
import Circle from "@/components/Circle.vue";
import L from "leaflet";
import "leaflet/dist/leaflet.css";

@Options({
  components: { Circle },
})
export default class Map extends Vue {
  cercleDivElement!: HTMLDivElement;
  mapDivElement!: HTMLDivElement;
  maxCircleHeight!: number;
  maxCircleWidth!: number;
  mounted() {
    this.cercleDivElement = this.$refs.Circle as HTMLDivElement;
    this.mapDivElement = this.$refs.MapDivElement as HTMLDivElement;
    this.maxCircleHeight = this.mapDivElement.offsetHeight * 0.9;
    this.maxCircleWidth = this.mapDivElement.offsetWidth * 0.9;

    // var map = L.map("openmap").setView([51.505, -0.09], 13);
    // L.tileLayer("https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png", {
    //   attribution:
    //     '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors',
    // }).addTo(map);
  }
  get ComponetIsHorizontal(): boolean {
    return this.mapDivElement.offsetWidth > this.mapDivElement.offsetHeight;
  }
  UpdateRadius(newRadius: number) {
    console.log("newRadius in map: ", newRadius);
    this.updateCircleRadius(newRadius);
  }
  updateCircleRadius(radius: number) {
    if (!this.ComponetIsHorizontal) {
      (this.$refs.Circle as Circle).UpdateRadius(
        (radius / 100) * this.maxCircleWidth + "px"
      );
    } else {
      (this.$refs.Circle as Circle).UpdateRadius(
        (radius / 100) * this.maxCircleHeight + "px"
      );
    }
  }
}
</script>

<style scoped>
#Map {
}
/* .AX {
  width: 100%;
  height: 100%;
} */
.inner-content {
  position: absolute;
  top: 0;
  bottom: 0;
  left: 0;
  right: 0;
}
</style>