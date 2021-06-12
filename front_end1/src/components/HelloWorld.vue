<template>
  <div class="SendArticleWithPositioning">
    <form @submit.prevent="submitForm">
      <label for="title">Title</label>
      <input id="title" type="text" v-model="titleInput" />
      <label for="content">Content</label>
      <textarea
        id="content"
        type="text"
        v-model="contentInput"
        style="resize: none"
      />
      <button type="submit">Submit</button>
      <p>Result : {{ resMessage }}</p>
    </form>
  </div>
</template>

<script lang="ts">
import { Options, Vue } from "vue-class-component";
type Position = { latitude: number; longitude: number };

@Options({})
export default class SendArticleWithPositioning extends Vue {
  titleInput = "";
  contentInput = "";
  resMessage = "";
  submitForm() {
    console.log("titleInput : ", this.titleInput);
    console.log("contentInput : ", this.contentInput);
    try {
      let currentPosition!: Position;
      currentPosition = this.getPosition();
      console.log("Position : ", currentPosition);
      // console.log("Position-latitude : ", currentPosition.latitude);
      // console.log("Position-longitude : ", currentPosition.longitude);
    } catch (e) {
      console.log("Position-Error : ", e);
    }
  }

  getPosition(): Position {
    if (!navigator.geolocation) {
      throw "Geolocation is not supported by your browser";
    } else {
      let res!: Position;
      // fsdfskl;
      navigator.geolocation.getCurrentPosition(
        (position: GeolocationPosition): void => {
          console.log(
            "position : ",
            position.coords.latitude,
            " , ",
            position.coords.longitude
          );
          res = {
            latitude: position.coords.latitude,
            longitude: position.coords.longitude,
          };
        },
        () => {
          throw "Unable to retrieve your location";
        }
      );
      return res;
    }
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
.SendArticleWithPositioning {
  height: 100vh;
  align-items: center;
  display: flex;
  justify-content: center;
  background: rgb(255, 171, 171);
}
.SendArticleWithPositioning > form * {
  width: 100%;
}
.SendArticleWithPositioning > form {
  width: 500px;
}
</style>
