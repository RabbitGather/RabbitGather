<template>
  <div class="ChatingBox">
    <h1>Real Time Chating Box</h1>
    <h2>Status:</h2>

    <p>{{ status }}</p>
    <div class="connectionBox">
      <h2>Connection ID:</h2>
      <p v-show="connectionID !== ''">
        {{ connectionID }}
      </p>
    </div>
    <input v-model="currentinput" class="messageInput" type="text" />
    <button @click="SubmitMessage">Submit</button>
    <h2>ChatRoom :</h2>
    <div class="chatList" v-for="(item, key, index) in chatList" :key="index">
      {{ index }} - {{ key }} : {{ item }}
    </div>
  </div>
</template>

<script lang="ts" >
import Peer from "peerjs";
import { Options, Vue } from "vue-class-component";
// import SendArticleWithPositioning from "@/components/SendArticleWithPositioning.vue"; // @ is an alias to /src
// import PerChat from "@/components/PeerChat.vue";
@Options({
  components: {
    // SendArticleWithPositioning,
    // PerChat,
  },
})
export default class RealTimeChatBox extends Vue {
  status = "Initializing ...";
  currentinput = "";
  connectionID = "(Not generated yet ...)";
  chatList = { debug: "OK" };
  beforeMount() {
    console.log("BeforeMount");
    // create websocket with server
    // peer = new Peer(null, {
    //   debug: 2,
    // });
  }
  SubmitMessage() {
    let inputText = this.currentinput;
    this.currentinput = "";
    console.log("SubmitMessage : " + inputText);
  }
}
</script>

<style scope>
.chating_box {
  flex-direction: column;
  height: 100vh;
  align-items: center;
  display: flex;
  justify-content: center;
  background: rgb(255, 171, 171);
}
.chating_box > * {
  width: 100%;
}
.connectionBox {
  border: 1px solid;
  margin-bottom: 5px;
}
.messageInput {
  width: 50%;
  margin-right: 3px;
}
.chatList {
  max-width: 750px;
  width: 100%;
  margin-top: 5px;
  border: 1px solid;
  margin: auto;
}
</style>