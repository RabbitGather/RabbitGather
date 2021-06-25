<template>
  <div class="PerChat">
    <h1>Send Article With Positioning</h1>
    <div class="status_box">
      <p>Connection ID: <br />{{ recvId }}</p>
      <!-- <p>Status : {{ status }}</p> -->
    </div>
    <br />
    <div class="message_box">
      <p>{{ message }}</p>
    </div>
  </div>
</template>

<script lang="ts">
import { Options, Vue } from "vue-class-component";
import Peer, { DataConnection } from "peerjs";

// import axios from "axios";
// import { PositionClass, PositionPoint } from "../global/Positions";
// type Article = { title: string; content: string; position: PositionPoint };
@Options({
  // props: {
  //   msg: String,
  // },
})
export default class PerChat extends Vue {
  private peer = new Peer(undefined, { debug: 2 });
  private lastPeerId = "";
  recvId = "";
  message = "";
  conn!: DataConnection | null;
  // status = "Reddy";
  created() {
    this.peer.on("open", (id) => {
      // Workaround for peer.reconnect deleting previous id
      if (this.peer.id === null) {
        console.log("Error : Received null id from peer open");
        this.peer.id = this.lastPeerId;
      } else {
        this.lastPeerId = this.peer.id;
      }

      console.log("ID: " + this.peer.id);
      this.recvId = this.peer.id;
      this.message = "Awaiting connection...";
    });
    this.peer.on("connection", (c) => {
      // Allow only a single connection
      if (this.conn && this.conn.open) {
        c.on("open", function () {
          c.send("Already connected to another client");
          setTimeout(() => {
            c.close();
          }, 500);
        });
        return;
      }

      this.conn = c;
      console.log("Connected to: " + this.conn.peer);
      this.message = "Connected";
      this.ready();
    });
  }
  ready() {
    if (this.conn === null) {
      console.log("Error : conn is null");
      return;
    }
    this.conn.on("data", (data) => {
      console.log("Data recieved : ", data);
      this.message +=
        (typeof data === typeof "" ? data : "(the data is not string)") + "\n";
    });
    this.conn.on("close", () => {
      this.message = "Connection reset<br>Awaiting connection...";
      // status.innerHTML = "Connection reset<br>Awaiting connection...";
      this.conn = null;
    });
    this.peer.on("disconnected", () => {
      this.message = "Connection lost. Please reconnect";
      console.log("Connection lost. Please reconnect");

      // Workaround for peer.reconnect deleting previous id
      this.peer.id = this.lastPeerId;
      // this.peer._lastServerId = this.lastPeerId;
      this.peer.reconnect();
    });
    this.peer.on("close", () => {
      this.conn = null;
      this.message = "Connection destroyed. Please refresh";
      console.log("Connection destroyed");
    });
    this.peer.on("error", (err) => {
      console.log(err);
      alert("" + err);
    });
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
.PerChat {
  flex-direction: column;

  height: 100vh;
  align-items: center;
  display: flex;
  justify-content: center;
  background: rgb(255, 171, 171);
}
.PerChat * {
  float: left;
  width: 100%;
}
</style>