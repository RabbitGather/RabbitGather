import { createApp } from "vue";
import App from "./App.vue";
import "./registerServiceWorker";
import router from "./router";
import store from "./store";
import "./index.css";
import { dragscrollNext } from "vue-dragscroll";


// import PositionFuncs from "./global/Position";

const vueapp = createApp(App);
vueapp.use(store).use(router).mount("#app");
vueapp.directive('dragscroll', dragscrollNext);

// App.directive("title", {
//   inserted: function (el, binding) {
//     document.title = el.dataset.title;
//   },
// });
