import { GetterTree } from "vuex";
import { AuthState } from "./state";
import { RootState } from "@/store";

export enum AuthGetterTypes {
  isAuthenticated = "isAuthenticated",
  UserSetting = "UserSetting",
}

// export type UserSetting = {

// }

export type AuthGetters = {
  [AuthGetterTypes.isAuthenticated](state: AuthState): boolean;
  [AuthGetterTypes.UserSetting](state: AuthState): UserSetting;
};

// the usr setting
export class UserSetting {
  RadiusRange = {
    Max: 1020,
    Min: 1,
  };
}

function CheckIsAuthenticated(state: AuthState) {
  console.log("state.API_ACCESS_TOKEN : ", state.API_ACCESS_TOKEN);
  return !!state.API_ACCESS_TOKEN;
}

export const getters: GetterTree<AuthState, RootState> & AuthGetters = {
  [AuthGetterTypes.isAuthenticated]: CheckIsAuthenticated,
  [AuthGetterTypes.UserSetting](state: AuthState) {
    console.log("state.UserSetting");

    return new UserSetting();
  },
};
