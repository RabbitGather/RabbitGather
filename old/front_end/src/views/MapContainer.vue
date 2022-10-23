<template>
  <div class="MapContainer bg-blue-500 flex flex-col">
    <Map
      ref="Map"
      :maxRadius="maxRadius"
      :minRadius="minRadius"
      class="Map bg-purple-400 flex-grow"
    >
    </Map>

    <Ruler
      ref="Ruler"
      @RadiusUpdate="RadiusUpdate"
      class="h-14 flex-none"
    ></Ruler>
  </div>
  <!-- </div> -->
</template>

<script lang="ts" >
import { Options, Vue } from "vue-class-component";
import StatusBar from "@/components/StatusBar.vue";
import Ruler from "@/components/Ruler.vue";
import Map from "@/components/Map.vue";
import store, { AllActionTypes } from "@/store";
import { UserSettings } from "@/store/app";
// import { io, Socket, Manager } from "socket.io-client";
import { GetPosition, PositionPoint } from "@/global/Positions";
import * as t from "@/views/type";
// import {*} from "@/views/type";

@Options({
  components: { StatusBar, Map, Ruler },
})
export default class MapContainer extends Vue {
  maxRadius = 0;
  minRadius = 0;
  CurrentRadius = 581199;

  beforeCreate() {
    let initListener = async () => {
      let timestamp = Math.floor(Date.now() / 1000);
      let position = await GetPosition();
      let path = `wss://api.meowalien.com/article/listen?timestamp=${timestamp}&radius=${this.CurrentRadius}&position={"x":${position.longitude},"y":${position.latitude}}`;
      console.log(path);
      let conn = new WebSocket(path);
      // Math.floor(Date.now() / 1000)
      // `api.meowalien.com/article/listen?timestamp=1627724266683&radius=3&position={"x":121.3996475828320,"latitude":25.017164133161643}`
      conn.onclose = (evt: Event) => {
        console.log("onclose evt: ", evt);
      };
      conn.onmessage = (evt: MessageEvent) => {
        console.log("onmessage evt: ", evt.data);
        let message = JSON.parse(evt.data) as t.ArticleChangeEvent;
        console.log("message: ", message.event);

        switch (message.event) {
          case "NEW":
            this.DrawPointOnMap({
              timestamp: message.timestamp,
              id: message.id,
              position: {
                x: message.position.x,
                y: message.position.y,
              },
            });
        }
      };

      conn.onerror = (evt: Event) => {
        console.log("onerror evt: ", evt);
      };
      conn.onopen = (ev: Event) => {
        // console.log("onopen ev: ", ev);
        conn.send(JSON.stringify({ action: "MESSAGE", text: "Test message" }));
      };
    };
    initListener();

    // read UserSetting and set the MaxRadius and MinRadius to Ruler components
    // store
    //   .dispatch(AllActionTypes.APP.GetUserInfo)
    //   .then((userinfo: UserSettings) => {
    //     this.maxRadius = userinfo.radaRadius.MaxRadius;
    //     this.minRadius = userinfo.radaRadius.MinRadius;
    //     (this.$refs.Ruler as Ruler).Init(this.maxRadius, this.minRadius);
    //   });
  }

  DrawPointOnMap(article: t.Article) {
    console.log("article: ", article.timestamp);
    console.log("article: ", article.id);
    console.log("article: ", article.position.y);
    console.log("article: ", article.position.x);
    (this.$refs.Map as Map).DrawArticleOnMap(article);
  }
  NewRadius(radius: number) {
    this.CurrentRadius = radius;
  }

  // will emit by Ruler
  RadiusUpdate(newRadius: number) {
    this.UpdateMap(newRadius);
  }

  NewRadishStep = 20;
  NewRadishGape = 10;
  lastRadius = 0;
  UpdateMap(newRadius: number) {
    // ask radish from server if radius reach the gape
    if (newRadius > this.NewRadishGape) {
      this.NewRadishGape += this.NewRadishStep;

      // pull new Radish from server and add in the map componet.
      // let radishes = this.GetNewRadishFromServer(lastRadius, newRadius);
      // (this.$refs.Map as Map).PushNewRadishes(radishes);
    }
    (this.$refs.Map as Map).UpdateRadius(newRadius);
  }
}
</script>

<style scoped>
</style>