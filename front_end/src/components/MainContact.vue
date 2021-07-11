<template>
  <div
    id="MainContact"
    class="
      bg-indigo-500
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
    {{ radius }}
    <div
      v-bind:style="{ width: getradius, height: getradius }"
      class="bg-yellow-400 rounded-full border-4 border-red-600"
      style="transition: all 0.2s"
    ></div>
  </div>
</template>

<script lang="ts">
import { Options, Vue, prop } from "vue-class-component";
class Props {
  radius = prop<number>({ required: true });
  magnification = prop<number>({ required: true });
}
@Options({})
export default class MainContact extends Vue.with(Props) {
  mounted() {}
  get getradius(): string {
    let r = this.radius;
    return r * this.magnification + "px";
  }
}
</script>

<style scoped>
#MainContact {
}
</style>