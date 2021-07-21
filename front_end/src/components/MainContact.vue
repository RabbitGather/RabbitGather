<template>
  <div
    id="MainContact"
    ref="contact_div"
    class="
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
    {{ RadiusString }}
    <div
      ref="cercle"
      v-bind:style="{ width: getWidth(), height: getHeight() }"
      class="rounded-full border-4 border-red-600"
      style="transition: all 0.2s"
    ></div>
  </div>
</template>

<script lang="ts">
import { Options, Vue, prop } from "vue-class-component";
class Props {
  radius = prop<number>({ required: true });
}
@Options({})
export default class MainContact extends Vue.with(Props) {
  contact!: HTMLDivElement;
  horizontal!: boolean;
  circleWidth = 0;
  reddy = false;
  cercle!: HTMLDivElement;
  maxContantHeight!: number;
  maxContantWidth!: number;
  mounted() {
    this.cercle = this.$refs.cercle as HTMLDivElement;
    this.contact = this.$refs.contact_div as HTMLDivElement;
    this.maxContantHeight = this.contact.offsetHeight * 0.9;
    this.maxContantWidth = this.contact.offsetWidth * 0.9;
    this.horizontal = this.contact.offsetWidth > this.contact.offsetHeight;
    this.reddy = true;
  }

  get RadiusString(): string {
    return Number(this.radius).toFixed(0);
  }

  getWidth(): string {  
    if (!this.reddy) {
      return "0%";
    }
    if (!this.horizontal) {
      // console.log("A: ", (this.radius / 100) * this.contact.offsetWidth + "%");
      return (this.radius / 100) * this.maxContantWidth + "px";
    } else {
      // console.log("this.contact.offsetHeight :", this.contact.offsetHeight);
      return this.getHeight();
    }
  }

  getHeight(): string {
    if (!this.reddy) {
      return "0%";
    }
    if (this.horizontal) {
      // console.log("B: ", (this.radius / 100) * this.contact.offsetHeight);
      // console.log("this.contact.offsetHeight: ", this.contact.offsetHeight);
      return (this.radius / 100) * this.maxContantHeight + "px";
    } else {
      return this.getWidth();
    }
  }

  distance(
    lat1: number,
    lon1: number,
    lat2: number,
    lon2: number,
    unit: string
  ) {
    if (lat1 == lat2 && lon1 == lon2) {
      return 0;
    } else {
      var radlat1 = (Math.PI * lat1) / 180;
      var radlat2 = (Math.PI * lat2) / 180;
      var theta = lon1 - lon2;
      var radtheta = (Math.PI * theta) / 180;
      var dist =
        Math.sin(radlat1) * Math.sin(radlat2) +
        Math.cos(radlat1) * Math.cos(radlat2) * Math.cos(radtheta);
      if (dist > 1) {
        dist = 1;
      }
      dist = Math.acos(dist);
      dist = (dist * 180) / Math.PI;
      dist = dist * 60 * 1.1515;
      if (unit == "K") {
        dist = dist * 1.609344;
      }
      if (unit == "N") {
        dist = dist * 0.8684;
      }
      return dist;
    }
  }
  // get getradius(): string {
  //   // let r = this.radius;
  //   return this.radius
  // }
}
</script>

<style scoped>
#MainContact {
}
</style>