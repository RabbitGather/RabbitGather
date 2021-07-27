import { createApp } from "vue";
import App from "./App.vue";
import "./registerServiceWorker";
import router from "./router";
import store from "./store";
import "./index.css";
import { dragscrollNext } from "vue-dragscroll";
// import VueSocketIO from "vue-socket.io";
// import SocketIO from "socket.io-client";

// Vue.use(new VueSocketIO({
//     debug: true,
//     connection: SocketIO('http://metinseylan.com:1992', options), //options object is Optional
//     vuex: {
//       store,
//       actionPrefix: "SOCKET_",
//       mutationPrefix: "SOCKET_"
//     }
//   })
// );

// import PositionFuncs from "./global/Position";

const options = { path: "/rabbit_gather" }; //Options object to pass into SocketIO

const vueapp = createApp(App);
// vueapp.directive("dragscroll", dragscrollNext);
vueapp

  .use(store)
  .use(router)
  // .use(
  //   new VueSocketIO({
  //     debug: true,
  //     connection: SocketIO("https://socket.meowalien.com", options), //options object is Optional
  //     vuex: {
  //       store,
  //       actionPrefix: "SOCKET_",
  //       mutationPrefix: "SOCKET_",
  //     },
  //   })
  // )
  .directive("dragscroll", dragscrollNext)
  .mount("#app");

// App.directive("title", {
//   inserted: function (el, binding) {
//     document.title = el.dataset.title;
//   },
// });
