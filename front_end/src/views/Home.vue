<template>
  <div class="home">
    <SendArticleWithPositioning />
    <br />
    <PerChat />
  </div>
</template>

<script lang="ts" >
import { Options, Vue } from "vue-class-component";
import SendArticleWithPositioning from "@/components/SendArticleWithPositioning.vue"; // @ is an alias to /src
import PerChat from "@/components/PeerChat.vue";
@Options({
  components: {
    SendArticleWithPositioning,
    PerChat,
  },
})
export default class Home extends Vue {}
</script>
