<template>
  <div class="bg-red-400" id="Map">
    <Circle class="Circle w-20 h-20"></Circle>
  </div>
</template>

<script lang="ts">
import { Options, Vue } from "vue-class-component";
import Circle from "../components/Circle.vue";
@Options({
  components: { Circle },
})


export default class Map extends Vue {
  UpdateRadius(newRadius: number) {
    console.log("newRadius in map: ", newRadius);
  }
}
</script>

<style scoped>
#Map {
}
</style>