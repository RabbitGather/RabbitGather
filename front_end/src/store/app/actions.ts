import { ActionTree } from "vuex";
import { AppMutations, AppMutationTypes } from "./mutations";
import { AppState } from "./state";
import { RootState, useStore } from "@/store";
import axios from "axios";
import { GenerateActionAugments } from "@/store/util";
import { UserSettings } from "@/store/app/state";
import { config } from "@vue/test-utils";

type ActionAugments = GenerateActionAugments<AppState, AppMutations>;

export enum AppActionTypes {
  StartGetTodos = "START_GET_TODOS",
  GetUserInfo = "GetUserInfo",
}

export type AppActions = {
  [AppActionTypes.StartGetTodos](context: ActionAugments): void;
  [AppActionTypes.GetUserInfo](context: ActionAugments, input: string): void;
};

export const actions: ActionTree<AppState, RootState> & AppActions = {
  async [AppActionTypes.StartGetTodos]({ commit }) {
    try {
      commit(AppMutationTypes.SetLoading, true);
      const { data } = await axios.get(
        "https://jsonplaceholder.typicode.com/todos"
      );
      commit(AppMutationTypes.SuccessGetTodos, data);
      commit(AppMutationTypes.SetLoading, false);
    } catch (err) {
      console.log("Error in AppActionTypes.StartGetTodos: ", err);
    }
  },
  [AppActionTypes.GetUserInfo]: GetUserInfo,
};

async function GetUserInfo(argumnt: ActionAugments): Promise<UserSettings> {
  return {
    basic: { name: "DEBUG_NAME", userid: 1 },
    radaRadius: { MaxRadius: 1000, MinRadius: 5 },
  };
  // let userinfo: UserSettings | undefined = argumnt.state.userInfo;
  // if (userinfo !== undefined) {
  //   return userinfo;
  // } else {
  //   try {
  //     let rootStore = useStore();
  //     userinfo = (
  //       await axios.get<UserSettings | undefined>(
  //         "https://api.meowalien.com/userinfo",
  //         {
  //           headers: {
  //             token: rootStore.state.AUTH.API_ACCESS_TOKEN,
  //           },
  //         }
  //       )
  //     ).data;
  //     if (userinfo === undefined) {
  //       throw "userinfo get from backend is undefined";
  //     }

  //     return userinfo
  //   } catch (e) {
  //     throw e;
  //   }
  // }
}
