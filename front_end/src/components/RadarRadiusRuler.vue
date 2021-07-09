<template>
  <div
    v-dragscroll.x
    id="RadarRadiusRuler"
    class=""
    ref="scrollboxcontainer"
    v-on:dragscrollstart="scrollboxcontainer_dragscrollstart"
    v-on:dragscrollmove="
      scrollboxcontainer_dragscrollmove($event.detail.deltaX)
    "
    v-on:dragscrollend="scrollboxcontainer_dragscrollend"
  >
    <div class="Triangle" ref="triangle"></div>
    <div class="bg-gray-400 h-full w-max flex flex-row" ref="scrollbox">
      <div class="" style="width: 50vw" value="0"></div>
      <div
        class="self-end w-auto flex flex-row items-baseline"
        v-for="(points, index) in TenScalePoints"
        :key="index"
      >
        <div
          v-show="points[0]"
          class="VerticalLine-short"
          v-bind:value="points[0]"
        ></div>
        <div
          v-show="points[1]"
          v-bind:value="points[1]"
          class="VerticalLine-short"
        ></div>
        <div
          v-show="points[2]"
          v-bind:value="points[2]"
          class="VerticalLine-short"
        ></div>
        <div
          v-show="points[3]"
          v-bind:value="points[3]"
          class="VerticalLine-short"
        ></div>
        <div
          v-show="points[4]"
          v-bind:value="points[4]"
          class="VerticalLine-long"
        ></div>
        <div
          v-show="points[5]"
          v-bind:value="points[5]"
          class="VerticalLine-short"
        ></div>
        <div
          v-show="points[6]"
          v-bind:value="points[6]"
          class="VerticalLine-short"
        ></div>
        <div
          v-show="points[7]"
          v-bind:value="points[7]"
          class="VerticalLine-short"
        ></div>
        <div
          v-show="points[8]"
          v-bind:value="points[8]"
          class="VerticalLine-short"
        ></div>
        <div
          v-show="points[9]"
          v-bind:value="points[9]"
          class="VerticalLine-node"
        >
          <div class="Node">
            <p class="PointText absolute bottom-0 text-xl">{{ points[9] }}</p>
          </div>
        </div>
      </div>
      <div class="" style="width: 50vw"></div>
    </div>
  </div>
</template>

<script lang="ts">
import { Options, Vue, prop } from "vue-class-component";

class Props {
  // optional prop
  min = prop<number>({ required: true });
  max = prop<number>({ required: true });
}
let htmlfontsize: number;
@Options({})
export default class RadarRadiusRuler extends Vue.with(Props) {
  TenScalePoints: number[][] = [];
  pointValue = 0;
  beforeMount() {
    // console.log(this.min);
    // console.log(this.max);
    let tempArray: number[] = [];
    for (let i = this.min; i <= this.max; i++) {
      if (tempArray.length == 10) {
        this.TenScalePoints.push([...tempArray]);
        tempArray = [];
      } else {
        tempArray.push(i);
      }
    }
    this.TenScalePoints.push(tempArray);
  }
  mounted() {
    htmlfontsize = parseInt(
      getComputedStyle(document.getElementsByTagName("html")[0], null).fontSize
    );
  }

  scrollboxcontainer_dragscrollstart() {}
  scrollboxcontainer_dragscrollmove(deltaX: number) {}
  remToPx(rem: number): number {
    return rem * htmlfontsize;
  }
  scrollboxcontainer_dragscrollend() {
    let scrollboxcontainer = this.$refs.scrollboxcontainer as HTMLDivElement;
    let triangle = this.$refs.triangle as HTMLDivElement;
    const targitBar = document.elementsFromPoint(
      triangle.offsetLeft + this.remToPx(0.36),
      window.innerHeight - 5
    )[0] as HTMLDivElement;

    const helfWindow = window.innerWidth / 2;
    let gape =
      (targitBar.offsetLeft == 0 ? helfWindow : targitBar.offsetLeft) +
      this.remToPx(0.5);
    let current = scrollboxcontainer.scrollLeft + helfWindow - 3;
    let gapeans = current > gape;
    scrollboxcontainer.scrollLeft +=
      gape - current + (gapeans ? this.remToPx(0.5) : -1 * this.remToPx(0.5));

    let theValue =
      parseInt(targitBar.attributes.getNamedItem("value")!.value) +
      (gapeans ? 1 : 0);
    if (theValue > this.max) {
      theValue = this.max;
    }
    if (this.pointValue != theValue) {
      this.pointValue = theValue;
      console.log("PointValue: " + this.pointValue);
    }
  }
}
</script>

<style scoped>
.VerticalLine-short {
  @apply h-3
          border-black
          w-4;
  border-left-width: 0.2rem;
}
.VerticalLine-long {
  @apply h-5
          border-black
          w-4;
  border-left-width: 0.2rem;
}
.VerticalLine-node {
  @apply h-5
          border-black
          w-4;
  border-left-width: 0.2rem;
}
.Node {
  @apply w-0 h-0 border-b-0 border-t-8 border-r-4 border-l-4 border-transparent self-end relative;
  left: -0.36rem;
  border-top-color: black;
}
.Triangle {
  @apply absolute block left-0 right-0 m-auto w-0 h-0 border-b-0  border-transparent self-end;
  left: -0.36rem;
  border-top-color: blue;
  border-top-width: 12px;
  border-right-width: 8px;
  border-left-width: 8px;
}
.PointText {
  left: -0.75rem;
  bottom: 0.2rem;
}
#RadarRadiusRuler {
}
</style>