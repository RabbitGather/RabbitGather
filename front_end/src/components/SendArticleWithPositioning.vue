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
import axios from "axios";
import { PositionClass, Position } from "../global/Position";
type Article = { title: string; content: string; position: Position };
@Options({
  // props: {
  //   msg: String,
  // },
})
export default class SendArticleWithPositioning extends Vue {
  titleInput = "";
  contentInput = "";
  resMessage = "";
  currentPosition!: Position;
  async submitForm() {
    this.currentPosition = await new PositionClass().getPosition();
    // console.log("Position-latitude : ", this.currentPosition.latitude);
    // console.log("Position-longitude : ", this.currentPosition.longitude);
    // console.log("titleInput : ", this.titleInput);
    // console.log("contentInput : ", this.contentInput);

    this.sendArticleToServer({
      title: this.titleInput,
      content: this.contentInput,
      position: this.currentPosition,
    });
  }
  async sendArticleToServer(article: Article) {
    console.log("article : ", article);
    try {
      const response = await axios.post("/api/post_article", article);
      console.log(response.data["result"]);
      this.resMessage = response.data["result"];
    } catch (error) {
      console.error(error);
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