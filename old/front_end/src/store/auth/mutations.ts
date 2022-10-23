import { MutationTree } from 'vuex'
import { AuthState } from './state'

export enum AuthMutationTypes {
  SET_API_ACCESS_TOKEN = 'SET_API_ACCESS_TOKEN'
}

export type AuthMutations = {
  [AuthMutationTypes.SET_API_ACCESS_TOKEN](state: AuthState, token: string): void;
}

export const mutations: MutationTree<AuthState> & AuthMutations = {
  [AuthMutationTypes.SET_API_ACCESS_TOKEN] (state, token) {
    state.API_ACCESS_TOKEN = token
  }
}
