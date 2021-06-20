import { Commit } from "vuex";
import { createStore } from "vuex";

const moduleB = {
  state: () => ({
    loginData: 'login'
  }),
  mutations: { 
    LOGIN_SET(states: any, params: object) {
      states.loginData = params
  }
  },
  actions: {
    loginAction(context: { commit: Commit }, params: object) {
      context.commit('LOGIN_SET', params)
  }
  }
}


export default createStore({
  state: {
    version: "",
  },
  mutations: {},
  actions: {},
  modules: {
    loginStore:moduleB,
  },
});
