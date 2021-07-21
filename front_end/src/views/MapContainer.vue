<template>
  <div class="MainContainer w-full h-full bg-blue-500 flex flex-col">
    <Map
      class="
        Map
        bg-purple-400
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
      <Radar class="Radar w-20 h-20"></Radar>
    </Map>
    <Rular
      class="
        bg-yellow-400
        RadarRadiusRuler
        flex-row
        static
        w-full
        h-14
        left-0
        bottom-0
        overflow-hidden
        flex-none
        order-1
        self-stretch
        flex-grow-0
      "
    ></Rular>
  </div>
</template>

<script lang="ts" >
import { Options, Vue } from "vue-class-component";
import StatusBar from "@/components/Map.vue";
import Radar from "@/components/Radar.vue";
import RadarRadiusRuler from "@/components/RadarRadiusRuler.vue";

@Options({
  components: { StatusBar, Radar, RadarRadiusRuler },
})
export default class MainContainer extends Vue {}
</script>

<style scoped>
</style>