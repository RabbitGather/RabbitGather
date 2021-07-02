<template>
  <div class="ChatingBox">
    <h1>Real Time Chating Box</h1>
    <h2>Connect To :</h2>
    <input type="text" v-model="connectTarget" />
    <button @click="DoConnect">SentMessage</button>
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
    <div class="chatList" v-for="(item, index) in chatList" :key="index">
      {{ index }} - {{ item.name }} : {{ item.message }}
    </div>
  </div>
</template>

<script lang="ts" >
import { Peer } from "../lib/peerjs/lib/peer";
import { DataConnection } from "../lib/peerjs/lib/dataconnection";
import { Options, Vue } from "vue-class-component";
import RTCconnectionToRemotePeersConfiguration from "@/config/RTCPeerConnectionConfiguration";
import { useStore, AllMutationTypes } from "@/store";

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
  chatList = [{ name: "OK", message: "MESSAGE" }];
  private peer!: Peer;
  private lastPeerId?: string;
  private connectionToRemotePeers?: DataConnection;
  connectTarget = "";
  conn?: DataConnection;
  /**
   * Get first "GET style" parameter from href.
   * This enables delivering an initial command upon page load.
   *
   * Would have been easier to use location.hash.
   */
  private getUrlParam(name: string) {
    name = name.replace(/[[]/, "\\[").replace(/[\]]/, "\\]");
    var regexS = "[\\?&]" + name + "=([^&#]*)";
    var regex = new RegExp(regexS);
    var results = regex.exec(window.location.href);
    if (results == null) return null;
    else return results[1];
  }

  SubmitMessage() {
    if (this.conn && this.conn.open) {
      let inputText = this.currentinput;
      this.currentinput = "";
      // console.log("SubmitMessage : " + inputText);
      // var msg = sendMessageBox.value;
      // sendMessageBox.value = "";
      this.conn.send(inputText);
      console.log("Sent: " + inputText);
      this.chatList.push({ name: "ME", message: inputText });
      // addMessage('<span class="selfMsg">Self: </span> ' + msg);
    } else {
      console.log("Connection is closed");
    }
  }

  DoConnect() {
    if (this.conn) {
      this.conn.close();
    }
    console.log("connectTarget: " + this.connectTarget);
    this.conn = this.peer.connect(this.connectTarget, {
      reliable: true,
    });
    let conn = this.conn as DataConnection;

    conn.on("open", () => {
      this.status = "Connected to: " + conn.peer;
      console.log("Connected to: " + conn.peer);

      // Check URL params for comamnds that should be sent immediately
      var command = this.getUrlParam("command");
      if (command) conn.send(command);
    });
    // Handle incoming data (messages only since this is the signal sender)
    conn.on("data", (data) => {
      this.chatList.push({ name: "Remote", message: data });
      // addMessage('<span class="peerMsg">Peer:</span> ' + data);
    });
    conn.on("close", () => {
      this.status = "Connection closed";
    });
  }
  private ready() {
    this.conn!.on("data", (data) => {
      console.log("Data recieved");
      // var cueString = '<span class="cueMsg">Cue: </span>';
      this.chatList.push({ name: "REMOTE", message: data });
    });
    this.conn!.on("close", ()=> {
      this.status = "Connection reset<br>Awaiting connection...";
      this.conn = undefined;
    });
  }
  beforeMount() {
    let vueStore = useStore();
    let apiToken = vueStore.state.AUTH.API_ACCESS_TOKEN as string;
    // routes.
    let apitoken = vueStore.state.AUTH.API_ACCESS_TOKEN;
    this.peer = new Peer(undefined, {
      host: "peerjs.localhost",
      port: 443,
      path: "/",
      key: "peerjs",
      secure: true,
      // token:apiToken,
      config: RTCconnectionToRemotePeersConfiguration,
      debug: 3,
    });
    let peer = this.peer;
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
    peer.on("connection", (c) => {
      // Allow only a single connection
      if (this.conn && this.conn.open) {
        c.on("open", function () {
          c.send("Already connected to another client");
          setTimeout(function () {
            c.close();
          }, 500);
        });
        return;
      }

      this.conn = c;
      console.log("Connected to: " + this.conn!.peer);
      this.status = "Connected";
      this.ready();
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
        // peer.id = this.lastPeerId;
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
    peer.on("error", (err: Error) => {
      console.log("On error event");
      console.log(err);
      alert("" + err);
    });
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