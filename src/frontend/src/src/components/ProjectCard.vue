<script setup lang="ts">
// import { ref } from 'vue'

const props = defineProps<{ msg: string }>()
// const { msg } = toRefs(props);

// const count = ref(0)
</script>

<template>
  <div class="bg-gray-200 p-2 m-2">
    {{ props.msg }}
  </div>
</template>

<style scoped>

</style>
