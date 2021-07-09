<template>
  <div id="StatusBar" class="bg-gray-500"></div>
</template>

<script lang="ts">
import { Options, Vue } from "vue-class-component";

@Options({})
export default class StatusBar extends Vue {}
</script>

<style scoped>
#StatusBar {
}
</style>