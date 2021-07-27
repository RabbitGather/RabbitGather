<template>
  <div
    v-dragscroll.x
    id="RadarRadiusRuler"
    class="overflow-hidden"
    ref="scrollboxcontainer"
    v-on:dragscrollstart="scrollboxcontainer_dragscrollstart"
    v-on:dragscrollmove="
      scrollboxcontainer_dragscrollmove($event.detail.deltaX)
    "
    v-on:dragscrollend="scrollboxcontainer_dragscrollend"
  >
    <div class="Triangle" ref="triangle"></div>
    <div class="bg-gray-400 h-full w-max flex flex-row" ref="scrollbox">
      <div class="" style="width: 50vw" v-bind:value="min"></div>
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
      <div class="" style="width: 50vw" v-bind:value="max"></div>
    </div>
  </div>
</template>

<script lang="ts">
import { Options, Vue, prop } from "vue-class-component";
import { remToPx } from "@/global/Util";
class Props {
  // optional prop
  // min = prop<number>({ required: true });
  // max = prop<number>({ required: true });
}

let scrollboxcontainer!: HTMLDivElement;
let pointValue = 0;
let triangle!: HTMLDivElement;
let current = 0;
let gapeans = false;
let gape = 0;
const RadiusUpdateEvent = "RadiusUpdate";
function timeout(ms: number) {
  return new Promise((resolve) => setTimeout(resolve, ms));
}
// let updated: boolean;
@Options({})
export default class Ruler extends Vue.with(Props) {
  private TenScalePoints: (number | undefined)[][] = [];

// scrollEvent(){
//   console.lo
// }

  max = 0;
  min = 0;
  async Init(max: number, min: number) {
    console.log("Ruler init Min: ", min);
    console.log("Ruler init Max: ", max);
    this.max = max;
    this.min = min;
    let tempArray: (number | undefined)[] = [];

    // to make sure the long scale is always on multiples of 10
    for (let index = min % 10; index > 1; index--) {
      tempArray.push(undefined);
    }
    for (let i = min; i <= max; i++) {
      tempArray.push(i);
      if (tempArray.length == 10) {
        this.TenScalePoints.push([...tempArray]);
        tempArray = [];
      }
    }
    this.TenScalePoints.push(tempArray);

    this.$nextTick(() => {
      scrollboxcontainer = this.$refs.scrollboxcontainer as HTMLDivElement;
      triangle = this.$refs.triangle as HTMLDivElement;
      let thecenterBar = scrollboxcontainer.querySelector(
        "[value='" + this.max / 2 + "']"
      ) as HTMLDivElement;

      scrollboxcontainer.scrollLeft = 0;
      // scrollboxcontainer.scrollLeft =
      //   thecenterBar.offsetLeft - scrollboxcontainer.offsetWidth / 2 + 3;
      // this.$emit(RadiusUpdateEvent, this.max / 2);
    });
  }

  scrollboxcontainer_dragscrollstart() {}
  scrollboxcontainer_dragscrollmove(deltaX: number) {
    this.updatePoint();
  }

  updatePoint() {
    let thisBar = document.elementsFromPoint(
      triangle.offsetLeft + remToPx(0.36),
      scrollboxcontainer.getBoundingClientRect().bottom - 5
      // window.innerHeight - 5
    )[0] as HTMLDivElement;

    gape =
      (thisBar.offsetLeft == 0
        ? scrollboxcontainer.offsetWidth / 2
        : thisBar.offsetLeft) + remToPx(0.5);

    current =
      scrollboxcontainer.scrollLeft + scrollboxcontainer.offsetWidth / 2 - 3;
    gapeans = current > gape;
    // if (thisBar === targitBar) {
    //   // updated = false;
    //   // return;
    // } else {
    //   targitBar = thisBar;
    // }
    let rawvalue = parseInt(thisBar.attributes.getNamedItem("value")!.value);
    let newValue = rawvalue + (gapeans ? 1 : 0);

    if (newValue <= this.max && pointValue != newValue) {
      pointValue = newValue;
      this.$emit(RadiusUpdateEvent, pointValue);
    }
    // updated = true;
  }
  scrollboxcontainer_dragscrollend() {
    scrollboxcontainer.scrollLeft +=
      gape - current + (gapeans ? remToPx(0.5) : -1 * remToPx(0.5));
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