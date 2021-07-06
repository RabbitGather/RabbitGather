<template>
  <div class="flex flex-row justify-center">
    <div class="ChatingBox w-3/4 children:bg-gray text-center">
      <h1 class="text-4xl">Real Time Chating Box</h1>
      <div class="border-2 border-black">
        <h2>Connection ID:</h2>
        <p v-show="myPeerID !== ''">
          {{ myPeerID }}
        </p>
      </div>
      <div class="border-black border-2">
        <h2 class="text-2xl">Connect To :</h2>
        <input
          type="text"
          v-model="connectTarget"
          class="border-2 border-black"
        />
        <button @click="DoConnect" class="ml-2 border-2 border-black">
          Connect
        </button>
      </div>

      <div class="border-black border-2">
        <h2>Status:</h2>
        <p>{{ status }}</p>
      </div>
      <div>
        <input
          v-model="currentinput"
          class="border-black border-2"
          type="text"
        />
        <button @click="SubmitMessage" class="ml-2 border-2 border-black">
          Sent
        </button>
      </div>
      <div class="border-2 border-black">
        <h2 class="text-2xl border-b-2 border-black">ChatRoom :</h2>
        <div
          class="chatList text-left"
          v-for="(item, index) in chatList"
          :key="index"
        >
          {{ index }} - {{ item.name }} : {{ item.message }}
        </div>
      </div>
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
import peerjsConfig from "../config/peerjs_config";
// import * as fs from "fs";
/* eslint-disable @typescript-eslint/no-var-requires */
// import fs from "fs";
// let a = file.readFileSync("../config/peerjs.config.json", "utf8");
// console.log(a);
// @Options({
//   components: {
//     // SendArticleWithPositioning,
//     // PerChat,
//   },
// })
export default class RealTimeChatBox extends Vue {
  nickname = "";
  status = "Initializing ...";
  currentinput = "";
  myPeerID = "(Not generated yet ...)";
  chatList = [{ name: "SYSTEM", message: "DEBUG_MESSAGE" }];
  connectTarget = "";
  private peer!: Peer;

  startPeerjs() {
    console.log(peerjsConfig.RunningHost);

    let vueStore = useStore();
    let apiToken = vueStore.state.AUTH.API_ACCESS_TOKEN as string;
    this.peer = new Peer(undefined, {
      host: "peerjs." + peerjsConfig.RunningHost,
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
      }
      console.log("NewPeerID : " + peerID);
      this.myPeerID = peerID;
      this.status = "Awaiting connection...";
    });
    this.peer.on(PeerEventType.Error, (err) => {
      console.log("On Error event");

      console.log("Error : " + err);
      this.status = err;
      // this.status = "Awaiting connection...";
    });
    // Connection from others
    this.peer.on(PeerEventType.Connection, (connection: DataConnection) => {
      console.log("On connection event");
      try {
        this.appendNewConnection(connection);
      } catch (e) {
        // connection.provider
        // console.log(connection.provider)
        this.status = "Fail to Connect to:" + connection.peer + "Error :" + e;
        return;
      }

      this.status = "Connected to: " + connection.peer;
      console.log(this.status);

      connection.on(ConnectionEventType.Data, (data) => {
        console.log("Data recieved: ", data);
        this.chatList.push({ name: connection.peer, message: data });
      });
      connection.on(ConnectionEventType.Close, () => {
        this.status = "Connection reset: " + connection.peer;
        this.removeConnection(connection.peer);
      });
      // }
    });

    this.peer.on(PeerEventType.Disconnected, (id: string) => {
      console.log("On disconnected event");

      this.status = "Disconnected with Server";
      console.log(this.status);
    });

    this.peer.on(PeerEventType.Close, () => {
      console.log("On close event");
      this.allConnections.forEach(
        (connection: DataConnection, connectionId) => {
          connection.close();
        }
      );
      // this.connectionToRemotePeers = undefined;
      this.status = "Connection destroyed. Please refresh";
      console.log("Connection destroyed");
    });

    this.peer.on(PeerEventType.Error, (err: Error) => {
      console.log("On error event");
      console.log(err);
      // alert("" + err);
    });
  }
  private allConnections = new Map();
  appendNewConnection(connection: DataConnection): void {
    this.allConnections.set(connection.peer, connection);
  }
  removeConnection(connectionID: string): void {
    this.allConnections.delete(connectionID);
  }
  beforeMount() {
    this.startPeerjs();
  }

  SubmitMessage() {
    console.log("Enter SubmitMessage");
    console.log("SubmitMessage - allConnections: ", this.allConnections);
    let inputText = this.currentinput;
    this.currentinput = "";
    this.chatList.push({ name: "ME", message: inputText });
    this.allConnections.forEach((connection: DataConnection, connectionID) => {
      if (!(connection && connection.open)) {
        this.removeConnection(connection.peer);
        return;
      }
      connection.send(inputText);
      console.log("Sent To: " + connection.peer + " , Message: " + inputText);
    });
  }

  async DoConnect() {
    let targetID = this.connectTarget;
    this.connectTarget = "";
    console.log("Try to connect to: " + targetID);
    let conn = this.peer.connect(targetID, {
      reliable: true,
    });
    if (conn === undefined) {
      this.status = "Fail to connect to: " + targetID;
      console.log(this.status);
      return;
    }
    let newConnection = conn as DataConnection;
    this.appendNewConnection(newConnection);
    newConnection.on(ConnectionEventType.Open, () => {
      this.status = "Connected to: " + newConnection.peer;
      console.log(this.status);
    });
    newConnection.on(ConnectionEventType.Data, (data: string) => {
      this.chatList.push({ name: newConnection.peer, message: data });
    });
    newConnection.on(ConnectionEventType.Close, () => {
      this.status = "Connection closed: " + newConnection.peer;
    });
  }
}
</script>

<style scope>
.ChatingBox > * {
  float: left;
  width: 100%;
  display: block;
}
/* .chating_box {
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
.status_box {
  border: 1px solid;
  margin-bottom: 5px;
}
.connection_box {
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
} */
</style>