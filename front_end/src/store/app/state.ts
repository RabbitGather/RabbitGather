import { TodoItem } from "./models";

export type AppState = {
  loading: boolean;
  todos: TodoItem[];
  userInfo: UserSettings | undefined;
};

// All Settings of the app
export type UserSettings= {
  basic: {
    name: string;
    userid:number;
  };
  radaRadius: {
    MaxRadius: number;
    MinRadius: number;
  };
}

export const state: AppState = {
  loading: false,
  todos: [],
  userInfo: undefined,
};
