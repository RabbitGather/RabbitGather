<template>
  <!-- <div id="nav"> -->
  <router-link to="/chat">Chat </router-link>
  <router-link to="/login">Login </router-link>
  <router-link to="/"> Home</router-link>

  <router-view />
  <!-- </div> -->
</template>

<style>
* {
  margin: 0;
  padding: 0;
}
body {
  margin: 0;
  padding: 0;
  width: 100%;
  position: relative;
  background: black;
}
#app {
  min-height: 100vh;
  width: 100%;
  height: 100vh;
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  background: white;
}
</style>
