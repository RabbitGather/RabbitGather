<template>
  <router-view class="flex-grow max-w-full" />
</template>

<script lang="ts" >
import { Options, Vue } from "vue-class-component";
import TheRabbit from "@/components/TheRabbit.vue";

@Options({
  components: { TheRabbit },
})
export default class HelloView extends Vue {}
</script>

<style scoped>
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