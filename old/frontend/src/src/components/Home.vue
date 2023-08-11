<script setup lang="ts">
import ProjectCard from '@/components/ProjectCard.vue'
// import { ref } from 'vue'

// defineProps<{ msg: string }>()

// const count = ref(0)
</script>

<template>
  <div class="w-full h-full flex flex-col justify-start items-center">
    <div class="flex flex-col gap-2 max-w-20em mt-4">
      <img src="/rabbit.svg" alt="logo" >
      <h1 class="text-3xl">
        RABBIT GATHER
      </h1>
    </div>
    <div class="bg-red-500 w-full flex flex-row">
<ProjectCard msg="54545545"/>
    </div>
  </div>
</template>

<style scoped>

</style>
