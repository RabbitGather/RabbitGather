<template>
  <div ref="MapDivElement" id="Map" class="">
    <!-- <Circle
      class="
        absolute
        z-20
        top-1/2
        left-1/2
        transform-gpu
        -translate-y-1/2 -translate-x-1/2
      "
      ref="Circle"
    ></Circle> -->
    <div
      ref="the_map"
      @mousewheel="wheelZoom"
      id="openmap"
      class="bg-red-500 overflow-hidden w-full h-full z-10"
      style="touch-action: none"
    ></div>
  </div>
</template>

<script lang="ts">
import { Options, Vue, prop } from "vue-class-component";
import Circle from "@/components/Circle.vue";
import L, { LeafletEvent, ResizeEvent, ZoomAnimEvent } from "leaflet";
import "leaflet/dist/leaflet.css";
import { GetPosition, PositionPoint } from "@/global/Positions";
import * as t from "@/views/type";
// import PinchZoom from "pinch-zoom-js";
// Long = X, Lat = Y
class Props {
  // optional prop
  maxRadius = prop<number>({ required: true });
  minRadius = prop<number>({ required: true });
}

@Options({
  components: { Circle },
})
export default class Map extends Vue.with(Props) {
  cercleDivElement!: HTMLDivElement;
  mapDivElement!: HTMLDivElement;
  maxCircleHeight!: number;
  maxCircleWidth!: number;
  map!: L.Map;
  circleOnMap!: L.Circle;
  theMapElement!: HTMLDivElement;

  DrawArticleOnMap(article: t.Article) {
    let classes_out = "";
    let classes_main = "";
    let text = "";
    let classes_foot = "";
    let text_foot = "";
    let myIcon = L.divIcon({
      // iconSize: null,
      html: `<div class="${classes_out}">
          <div class="${classes_main}">${text}</div>
          <div class="${classes_foot}">${text_foot}</div>
        </div>`,
      // iconUrl:
      //   "https://raw.githubusercontent.com/pointhi/leaflet-color-markers/master/img/marker-icon-2x-red.png",
      // iconSize: [23, 40],
      // iconAnchor: [22, 94],
      // popupAnchor: [-3, -76],
      // shadowUrl: "my-icon-shadow.png",
      // shadowSize: [68, 95],
      // shadowAnchor: [22, 94],
    });
    L.marker([article.position.y, article.position.x], {
      icon: myIcon,
    }).addTo(this.map);

    this.map.flyTo([article.position.y, article.position.x]);
  }
  wheelZoom(e: WheelEvent) {
    // UP < 0
    // DOWN > 0
    let isUp = e.deltaY < 0;
    let currentZoom = this.map.getZoom();
    if (isUp && currentZoom >= 18) {
      return;
    } else if (!isUp && currentZoom <= 10) {
    }
    // console.log("zoom: ", currentZoom);
    this.map.setZoom(currentZoom + (isUp ? 0.1 : -0.1));
  }

  mounted() {
    this.cercleDivElement = this.$refs.Circle as HTMLDivElement;
    this.mapDivElement = this.$refs.MapDivElement as HTMLDivElement;
    this.maxCircleHeight = this.mapDivElement.offsetHeight * 0.9;
    this.maxCircleWidth = this.mapDivElement.offsetWidth * 0.9;
    this.theMapElement = this.$refs.the_map as HTMLDivElement;
    // let pz = new PinchZoom(this.theMapElement, {
    //   onZoomStart: () => {
    //     console.log("onZoomStart: ");
    //   },
    // });
    // pz.disable();

    // let startPoint = {
    //   latitude: 25.040056717110396,
    //   longitude: 121.51187490970621,
    // };
    // pin to current position
    GetPosition().then((startPoint: PositionPoint) => {
      console.log("currentPosition: ", startPoint);
      this.map = L.map("openmap", {
        // scrollWheelZoom: "center",
        // zoomSnap: 0.1,
        // zoomDelta: 0.1,
        minZoom: 10,
        // boxZoom: false,
        // doubleClickZoom: false,
        // touchZoom: true,
        // zoomControl: true,
        // dragging: true,
        // trackResize: true,
        // scrollWheelZoom: false,
      }).setView([startPoint.latitude, startPoint.longitude], 18);
      L.tileLayer("https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png", {
        attribution:
          '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors',
      }).addTo(this.map);
      this.circleOnMap = L.circle(
        [startPoint.latitude, startPoint.longitude],
        20
      ).addTo(this.map);

      this.map.on("zoomend", (a: LeafletEvent) => {
        console.log("zoomend: ", this.map.getZoom());
      });
      this.MaxRadius = this.NewMaxRadius();

      // this.map.setView(
      //   [currentPosition.latitude, currentPosition.longitude],
      //   18
      // );
      // this.circleOnMap = L.circle(
      //   [currentPosition.latitude, currentPosition.longitude],
      //   1
      // ).addTo(this.map);
    });

    // this.$nextTick(() => {
    //   console.log("maxRadius: ", this.maxRadius);
    //   console.log("minRadius: ", this.minRadius);
    // });

    // 0~18
    // 2~18
    // this.map.on("zoomend", (a: LeafletEvent) => {
    //   console.log(this.map.getZoom());
    // });

    // 初始尺度的最高值
  }
  MaxRadius = 0;

  // Get the max radius in current zoom index
  NewMaxRadius(): number {
    var center = this.map.getCenter();
    let dist = 0;
    if (!this.ComponetIsHorizontal()) {
      var eastBound = this.map.getBounds().getEast();

      var centerEast = L.latLng(center.lat, eastBound);
      dist = center.distanceTo(centerEast);
    } else {
      var northBound = this.map.getBounds().getNorth();
      var centerNorth = L.latLng(northBound, center.lng);
      dist = center.distanceTo(centerNorth);
    }

    return dist * 0.75;
  }

  UpdateCirclePoint(latitude: number, longitude: number) {
    this.circleOnMap.setLatLng([latitude, longitude]);
  }
  UpdateCircleRadius(radius: number) {
    console.log("New Radius: ", radius);
    console.log("MaxRadius: ", this.MaxRadius);
    this.circleOnMap.setRadius(radius);
    if (this.MaxRadius < radius) {
      let currentZoom = this.map.getZoom();
      this.map.setZoom(currentZoom - 0.1);
      // while (!(currentZoom<this.map.getZoom())){}
      this.MaxRadius = this.NewMaxRadius();
    }
    //     else if (this.MinRadius > radius){
    // // 縮小地圖
    //     }
  }

  ComponetIsHorizontal(): boolean {
    return this.mapDivElement.offsetWidth > this.mapDivElement.offsetHeight;
  }

  // will call by parent component
  UpdateRadius(newRadius: number) {
    this.UpdateCircleRadius(newRadius);
  }
  updateCircleRadius(radius: number) {
    if (!this.ComponetIsHorizontal()) {
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