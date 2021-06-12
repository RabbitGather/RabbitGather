<template>
  <div class="home">
    <SendArticleWithPositioning/>
  </div>
</template>

<script lang="ts">
import { Options, Vue } from "vue-class-component";
import SendArticleWithPositioning from "@/components/SendArticleWithPositioning.vue"; // @ is an alias to /src

@Options({
  components: {
    SendArticleWithPositioning,
  },
})
export default class Home extends Vue {}
</script>
