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
      <p v-show="myPeerID !== ''">
        {{ myPeerID }}
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
import { PeerEventType, ConnectionEventType } from "../lib/peerjs/lib/enums";
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
  myPeerID = "(Not generated yet ...)";
  chatList = [{ name: "OK", message: "MESSAGE" }];
  connectTarget = "";
  private peer!: Peer;
  private lastPeerId?: string;
  private connectionToRemotePeers?: DataConnection;
  startPeerjs() {
    let vueStore = useStore();
    let apiToken = vueStore.state.AUTH.API_ACCESS_TOKEN as string;
    this.peer = new Peer(undefined, {
      host: "peerjs.localhost",
      port: 443,
      path: "/",
      key: "peerjs",
      secure: true,
      token: apiToken,
      config: RTCconnectionToRemotePeersConfiguration,
      debug: 3,
    });

    this.peer.on(PeerEventType.Open, (peerID: string) => {
      console.log("On open event");
      if (peerID === null) {
        console.log("Error : Received null id from peer open");
        return;
      } else {
        this.lastPeerId = this.myPeerID;
      }

      console.log("NewPeerID : " + peerID);
      this.myPeerID = peerID;
      this.status = "Awaiting connection...";
    });

    this.peer.on(
      PeerEventType.Connection,
      (remoteConnection: DataConnection) => {
        console.log("On connection event");

        // Allow only a single connection
        if (this.connectionToRemotePeers && this.connectionToRemotePeers.open) {
          remoteConnection.on("open", () => {
            remoteConnection.send("Already connected to another client");

            setTimeout(() => {
              remoteConnection.close();
            }, 500);
          });
          return;
        }

        this.connectionToRemotePeers = remoteConnection;
        console.log("Connected to: " + this.connectionToRemotePeers.peer);
        this.status = "Connected to: " + this.connectionToRemotePeers.peer;
        this.readyToReceiveMessage();
      }
    );

    this.peer.on(PeerEventType.Disconnected, () => {
      console.log("On disconnected event");

      this.status =
        "Connection to " +
        this.connectionToRemotePeers!.peer +
        " lost. Please reconnect";
      console.log(this.status);

      if (this.lastPeerId) {
        console.log("Try to reconnect  to: " + this.lastPeerId);
        this.peer.reconnect();
      }
    });
    this.peer.on(PeerEventType.Close, () => {
      console.log("On close event");

      this.connectionToRemotePeers = undefined;
      this.status = "Connection destroyed. Please refresh";
      console.log("Connection destroyed");
    });
    this.peer.on(PeerEventType.Error, (err: Error) => {
      console.log("On error event");
      console.log(err);
      // alert("" + err);
    });
  }
  beforeMount() {}

  SubmitMessage() {
    if (this.connectionToRemotePeers && this.connectionToRemotePeers.open) {
      let inputText = this.currentinput;
      this.currentinput = "";
      this.connectionToRemotePeers.send(inputText);
      console.log("Sent: " + inputText);
      this.chatList.push({ name: "ME", message: inputText });
    } else {
      console.log("Connection is closed");
    }
  }

  DoConnect() {
    if (this.connectionToRemotePeers) {
      this.connectionToRemotePeers.close();
    }
    console.log("connectTarget: " + this.connectTarget);
    this.connectionToRemotePeers = this.peer.connect(this.connectTarget, {
      reliable: true,
    });
    let conn = this.connectionToRemotePeers as DataConnection;

    conn.on(ConnectionEventType.Open, () => {
      this.status = "Connected to: " + conn.peer;
      console.log("Connected to: " + conn.peer);
    });
    conn.on(ConnectionEventType.Data, (data) => {
      this.chatList.push({ name: "Remote", message: data });
    });
    conn.on(ConnectionEventType.Close, () => {
      this.status = "Connection closed";
    });
  }
  private readyToReceiveMessage() {
    this.connectionToRemotePeers!.on(ConnectionEventType.Data, (data) => {
      console.log("Data recieved: ", data);
      // var cueString = '<span class="cueMsg">Cue: </span>';
      this.chatList.push({ name: "REMOTE", message: data });
    });
    this.connectionToRemotePeers!.on(ConnectionEventType.Close, () => {
      this.status = "Connection reset<br>Awaiting connection...";
      this.connectionToRemotePeers = undefined;
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