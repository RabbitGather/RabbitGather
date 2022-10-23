<template>
  <div
    v-bind:style="{ width: circleWidth, height: circleHeight }"
    class="Circle rounded-full border-4 border-red-600"
    style="transition: all 0.2s"
    id="Circle"
  ></div>
</template>


<script lang="ts">
import { Options, Vue, prop } from "vue-class-component";
class Props {
  // optional prop
  // radius = prop<number>({ required: true });
  // max = prop<number>({ required: true });
}

@Options({})
export default class Circle extends Vue.with(Props) {
  circleWidth = "0px";
  circleHeight = "0px";
  UpdateRadius(newRadius: string) {
    console.log("newRadius in Circle:", newRadius);
    this.circleWidth = newRadius;
    this.circleHeight = newRadius;
  }
}
</script>

<style scoped>
#Circle {
}
</style>