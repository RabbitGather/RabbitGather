import { GetterTree } from "vuex";
import { AuthState } from "./state";
import { RootState } from "@/store";

export enum AuthGetterTypes {
  isAuthenticated = "isAuthenticated",
}

export type AuthGetters = {
  [AuthGetterTypes.isAuthenticated](state: AuthState): boolean;
};

export const getters: GetterTree<AuthState, RootState> & AuthGetters = {
  [AuthGetterTypes.isAuthenticated](state: AuthState) {
      console.log("state.API_ACCESS_TOKEN : ", state.API_ACCESS_TOKEN)
    return !!state.API_ACCESS_TOKEN;
  },
};
