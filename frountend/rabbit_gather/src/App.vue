<script setup lang="ts">
import HelloWorld from './components/Home.vue'
import router from '@/router'

</script>

<template>
 <div class="main">
   <router-view/>
 </div>

</template>

<style scoped>
  .main{
    padding: 0;
    background: #747bff;
    display: block;
    width: 100%;
    height: 100%;
  }
</style>
