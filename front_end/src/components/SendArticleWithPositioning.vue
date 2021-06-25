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
import { PositionClass, PositionPoint } from "../global/Positions";
type Article = { title: string; content: string; position: PositionPoint };
@Options({
  // props: {
  //   msg: String,
  // },
})
export default class SendArticleWithPositioning extends Vue {
  titleInput = "";
  contentInput = "";
  resMessage = "";
  currentPosition!: PositionPoint;
  async submitForm() {
    this.currentPosition = await new PositionClass().getPosition();

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