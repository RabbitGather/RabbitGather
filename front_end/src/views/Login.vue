<template>
  <div class="login">
    <!-- <h1>LOGIN PAGE</h1>
    <form class="login_form" @submit.prevent="login">
      <label>User name</label>
      <input v-model="username" type="text" placeholder="UserName" />
      <label v-show="errorMessagByUserNameSideText">{{
        errorMessagByUserNameSideText
      }}</label>

      <label>Password</label>
      <input v-model="password" type="password" placeholder="Password" />
      <label v-show="errorMessagByPasswordSideText">{{
        errorMessagByPasswordSideText
      }}</label>
      <hr />
      <button type="submit">Login</button>
      <label v-show="resErrorMessage">{{ resErrorMessage }}</label>
      <loading
        v-model:active="isLoading"
        :can-cancel="true"
        :on-cancel="onCancel"
        :is-full-page="fullPage"
      />
    </form> -->
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
// import{store} from "vuex"
import { Provide } from 'vue-property-decorator'

import { State, Getter, Action, Mutation, namespace } from "vuex-class";
import {
  action_names,
  AUTH_NAME_SPACE,
  AuthRequestFunc,
} from "../store/auth/types";
// import Loading  from 'vue-loading-overlay';
// const loginModule = namespace('loginStore')

@Options({
  // props: {
  //   msg: String,
  // },
})
export default class LoginPage extends Vue {
  

    // @Provide() private data: any = {
    //     pageName: 'loginPage'
    // }


  // @namespace(AUTH_NAME_SPACE).Action(action_names.AUTH_REQUEST)
  // authRequest!: AuthRequestFunc;
  // isLoading = false;
  // username = "";
  // password = "";
  // resErrorMessage = "";
  // login() {
  //   let formatCorrect = true;
  //   if (this.username === "") {
  //     this.showErrorMessagByUserNameSide("使用者名稱不得為空白");
  //     formatCorrect = false;
  //   }
  //   if (this.password === "") {
  //     this.showErrorMessagByPasswordSide("密碼不得為空白");
  //     formatCorrect = false;
  //   }
  //   if (formatCorrect) {
  //     this.isLoading = true;
  //     this.authRequest({
  //       username: this.username,
  //       password: this.password,
  //       usersystem: "business_support_system",
  //     })
  //       .then(() => {
  //         this.isLoading = false;
  //         this.$router.push("/");
  //       })
  //       .catch((err) => {
  //         this.isLoading = false;
  //         this.resErrorMessage = err;
  //       });
  //   }
  // }
  // errorMessagByUserNameSideText = "";
  // showErrorMessagByUserNameSide(errorMessage: string) {
  //   this.errorMessagByUserNameSideText = errorMessage;
  // }
  // errorMessagByPasswordSideText = "";
  // showErrorMessagByPasswordSide(errorMessage: string) {
  //   this.errorMessagByPasswordSideText = errorMessage;
  // }
}
</script>
