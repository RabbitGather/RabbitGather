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
    //color: black;
    //background: bisque;
    width: 100%;
    height: 100%;
  }
</style>
