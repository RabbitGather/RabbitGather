<template>
  <div class="flex flex-col w-screen h-screen justify-center items-center">
    <div
      class="marqee w-full h-7 flex-shrink-0 overflow-hidden relative bg-black"
    >
      <ul class="pl-0 flex absolute">
        <li>jdaoisdjaisodjaspidjsaodpasd</li>
        <li>jdaoisdjaisodjaspidjsaodpasd</li>
        <li>jdaoisdjaisodjaspidjsaodpasd</li>
        <li>jdaoisdjaisodjaspidjsaodpasd</li>
        <li>jdaoisdjaisodjaspidjsaodpasd</li>
        <li>jdaoisdjaisodjaspidjsaodpasd</li>
      </ul>
    </div>
    <div class="bg-red-500 flex-grow relative w-full overflow-hidden">
      <img src="/img/map.jpg" alt="" class="object-cover w-full" />
    </div>
    <!-- :class="innerHeight > innerWidth ? 'vertical' : 'horizontal'" -->
  </div>
</template>

<script lang="ts" >
import { Options, Vue } from "vue-class-component";
import TheRabbit from "@/components/TheRabbit.vue";

@Options({
  components: { TheRabbit },
})
export default class WellcomePage extends Vue {
  get innerHeight(): number {
    return window.innerHeight;
  }
  get innerWidth(): number {
    return window.innerWidth;
  }
}
</script>

<style scoped>
.marqee > ul {
  animation: marqee 15s linear infinite;
}

.marqee > ul > li {
  white-space: nowrap;
  margin-right: 2em;
  font-weight: bold;
  color: white;
}

@keyframes marqee {
  0% {
    left: 100%;
    transform: translateX(0%);
  }

  100% {
    left: 0;
    transform: translateX(-100%);
  }
}
.vertical {
  @apply h-full;
}

.horizontal {
  @apply w-full;
}
#bean_curd_skin_right {
}
#bean_curd_skin_left {
}
.a {
  fill: #fff;
}
.b {
  fill: #f385a3;
}
.c {
  fill: #ff9600;
}
.d {
  fill: #feee0d;
}
.e {
  fill: #0053e1;
}
.f {
  fill: #04b0d2;
}
.g {
  fill: #02bb2d;
}
.h {
  fill: #d67498;
}
</style>