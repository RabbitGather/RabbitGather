<template>
  <div class="login">
    <h1>LOGIN PAGE</h1>
    <form class="login_form" @submit.prevent="login">
      <label>User name</label>
      <input v-model="username" type="text" placeholder="UserName" />
      <!-- <label v-show="errorMessagByUserNameSideText">{{
        errorMessagByUserNameSideText
      }}</label> -->

      <label>Password</label>
      <input v-model="password" type="password" placeholder="Password" />
      <!-- <label v-show="errorMessagByPasswordSideText">{{
        errorMessagByPasswordSideText
      }}</label> -->
      <hr />
      <button type="submit">Login</button>
      <label v-show="resMessage">{{ resMessage }}</label>
      <!-- /> -->
    </form>
  </div>
</template>
<style scope>
.login_form {
  display: flex;
  flex-direction: column;
  width: 300px;
  padding: 10px;
}
</style>

<script lang="ts">
import { Options, Vue } from "vue-class-component";
import { useStore, AllMutationTypes } from "@/store";
import routes from "@/router";
import axios from "axios";
import { PositionPoint } from "../global/Positions";
import { sleep } from "@/global/Util";
type LoginResponse = {
  ok: boolean;
  err: string;
  token: string;
};

@Options({})
export default class LoginPage extends Vue {
  username = "";
  // errorMessagByUserNameSideText = "";
  password = "";
  // errorMessagByPasswordSideText = "";
  resMessage = "";
  // do login
  async login() {
    let username = this.username;
    let password = this.password;
    let err = await this.checkInputCorrect(username, password);
    if (err !== "") {
      this.resMessage = err;
      return;
    }
    // let position: PositionPoint;
    // // try {
    // //   position = await new PositionClass().getPosition();
    // // } catch (e) {
    // //   this.resMessage =
    // //     e == typeof ""
    // //       ? e
    // //       : "Error : the PositionClass().getPosition() error object should be string";
    // //   return;
    // // }
    let loginResult = await this.sentLoginRequest(username, password);
    if (!loginResult.ok) {
      // err
      this.resMessage = "Fail to login, username or password is wrong";
      return;
    }
    if (loginResult.token === "") {
      //  err
      this.resMessage = "Error, loginResult.token is empty";
      return;
    }
    const store = useStore();
    store.commit(AllMutationTypes.AUTH.SET_API_ACCESS_TOKEN, loginResult.token);
    this.resMessage = "Login success, redirect to home after 3 sec ...";
    // sleep
    await sleep(3000);
    routes.push({ name: "Home" });
    return;
  }
  // sent the login request to server
  private async sentLoginRequest(
    username: any,
    password: any
    // position: PositionPoint
  ): Promise<LoginResponse> {
    return new Promise<LoginResponse>((resolve, reject) => {
      // try {
      axios
        .post<LoginResponse>("/api/login", {
          username: username,
          password: password,
          // position: position,
        })
        .then((resp) => {
          resolve(resp.data);
        })
        .catch((e) => {
          reject("Error when call /api/login : " + e);
        });
    });
  }
  // check if the input format correct
  private async checkInputCorrect(
    username: any,
    password: any
  ): Promise<string> {
    // check
    return "";
  }
}
</script>


<style scoped>
.login {
  flex-direction: column;
  height: 100vh;
  align-items: center;
  display: flex;
  justify-content: center;
  background: rgb(255, 171, 171);
}
.login > form * {
  width: 100%;
}
.login > form {
  width: 500px;
}
.login_form {
  display: block;
  float: left;
}
</style>