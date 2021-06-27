<template>
  <div class="ChatingBox">
    <h1>Real Time Chating Box</h1>
    <h2>Status:</h2>

    <p>{{ status }}</p>
    <div class="connectionBox">
      <h2>Connection ID:</h2>
      <p v-show="brokeringID !== ''">
        {{ brokeringID }}
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
import Peer, { DataConnection } from "peerjs";
import { Options, Vue } from "vue-class-component";
import RTCconnectionToRemotePeersConfiguration from "@/config/RTCPeerConnectionConfiguration";

@Options({
  components: {
    // SendArticleWithPositioning,
    // PerChat,
  },
})
export default class RealTimeChatBox extends Vue {
  status = "Initializing ...";
  currentinput = "";
  brokeringID = "(Not generated yet ...)";
  chatList = { debug: "OK" };
  private peer!: Peer;
  private lastPeerId?: string;
  private connectionToRemotePeers?: DataConnection;

  beforeMount() {
    console.log("BeforeMount");

    let peer = new Peer(undefined, {
      host: "peerjs.localhost",
      port: 443,
      path: "/",
      key: "peerjs",
      secure: true,
      config: RTCconnectionToRemotePeersConfiguration,
      debug: 3,
    });
    this.peer = peer;
    peer.on("open", (brokeringID: string) => {
      console.log("On open event");
      // Workaround for peer.reconnect deleting previous id
      if (brokeringID === null) {
        console.log("Error : Received null id from peer open");
        return;
        // peer.id = lastPeerId;
      } else {
        this.lastPeerId = brokeringID;
      }

      console.log("brokeringID : " + brokeringID);
      this.brokeringID = brokeringID;
      this.status = "Awaiting connection...";
    });
    peer.on("connection", (remoteConnection: DataConnection) => {
      console.log("On connection event");

      // Allow only a single connection
      if (this.connectionToRemotePeers && this.connectionToRemotePeers.open) {
        remoteConnection.on("open", () => {
          remoteConnection.send(
            "Error : This Peer is already connected to another client"
          );
          setTimeout(function () {
            remoteConnection.close();
          }, 500);
        });
        return;
      }

      this.connectionToRemotePeers = remoteConnection;
      console.log("Connected to: " + this.connectionToRemotePeers.peer);
      this.status = "Connected to: " + this.connectionToRemotePeers.peer;
      // ready();
    });
    peer.on("disconnected", () => {
      console.log("On disconnected event");

      this.status =
        "Connection to " +
        this.connectionToRemotePeers!.peer +
        " lost. Please reconnect";
      console.log(this.status);

      // Workaround for peer.reconnect deleting previous id
      if (this.lastPeerId) {
        console.log("Try to reconnect  to: " + this.lastPeerId);
        peer.id = this.lastPeerId;
        peer.reconnect();
      }
      // peer. = lastPeerId;
    });
    peer.on("close", () => {
      console.log("On close event");

      this.connectionToRemotePeers = undefined;
      this.status = "Connection destroyed. Please refresh";
      console.log("Connection destroyed");
    });
    peer.on("error", (err) => {
      console.log("On error event");

      console.log(err);
      alert("" + err);
    });
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