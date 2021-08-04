import { createRouter, createWebHistory, RouteRecordRaw } from "vue-router";
// import Home from "../views/Home.vue";
// import RealTimeChatBox from "../views/RealTimeChatBox.vue";
import HomeView from "../views/HomeView.vue";
import MapContainer from "../views/MapContainer.vue";
import HelloView from "../views/HelloView.vue";
import RabbitPage from "../views/RabbitPage.vue";
import WellcomePage from "../views/WellcomePage.vue";

const routes: Array<RouteRecordRaw> = [
  {
    path: "/login",
    name: "Login",
    component: () =>
      import(/* webpackChunkName: "about" */ "../views/Login.vue"),
    meta: {
      title: "RabbitGather - Login",
      metaTags: [
        {
          name: "description",
          content: "RabbitGather Login page",
        },
        {
          property: "og:description",
          content: "RabbitGather Login page",
        },
      ],
    },
  },
  {
    path: "/hello",
    component: HelloView,

    children: [
      {
        path: "wellcome",
        component: WellcomePage,
        meta: {
          title: "RabbitGather",
          metaTags: [
            {
              name: "description",
              content: "RabbitGather Main page",
            },
            {
              property: "og:description",
              content: "RabbitGather Main page",
            },
          ],
        },
      },
      {
        path: "rabbit",
        component: RabbitPage,
        meta: {
          title: "RabbitGather",
          metaTags: [
            {
              name: "description",
              content: "RabbitGather Main page",
            },
            {
              property: "og:description",
              content: "RabbitGather Main page",
            },
          ],
        },
      },
    ],
    name: "HelloPage",
    meta: {
      title: "RabbitGather",
      metaTags: [
        {
          name: "description",
          content: "RabbitGather Hello page",
        },
        {
          property: "og:description",
          content: "RabbitGather Hello page",
        },
      ],
    },
  },
  {
    path: "/",
    name: "Home",
    component: HomeView,
    children: [
      {
        path: "",
        component: MapContainer,
        meta: {
          title: "RabbitGather",
          metaTags: [
            {
              name: "description",
              content: "RabbitGather Main page",
            },
            {
              property: "og:description",
              content: "RabbitGather Main page",
            },
          ],
        },
      },
      // {
      //   path: "chat",
      //   component: RealTimeChatBox,
      //   name: "RealTimeChatBox",
      //   meta: {
      //     title: "RabbitGather - RealTimeChatBox",
      //     metaTags: [
      //       {
      //         name: "description",
      //         content: "RabbitGather - RealTimeChatBox",
      //       },
      //       {
      //         property: "og:description",
      //         content: "RabbitGather - RealTimeChatBox",
      //       },
      //     ],
      //   },
      // },
    ],
  },
  // {
  //   path: "/chat",
  //   name: "RealTimeChatBox",
  //   component: RealTimeChatBox,
  //   meta: {
  //     title: "RabbitGather - RealTimeChatBox",
  //     metaTags: [
  //       {
  //         name: "description",
  //         content: "RabbitGather - RealTimeChatBox",
  //       },
  //       {
  //         property: "og:description",
  //         content: "RabbitGather - RealTimeChatBox",
  //       },
  //     ],
  //   },
  // },
];

const router = createRouter({
  history: createWebHistory(process.env.BASE_URL),
  routes,
});

// This callback runs before every route change, including on page load.
router.beforeEach((to, from, next) => {
  // This goes through the matched routes from last to first, finding the closest route with a title.
  // e.g., if we have `/some/deep/nested/route` and `/some`, `/deep`, and `/nested` have titles,
  // `/nested`'s will be chosen.
  const nearestWithTitle = to.matched
    .slice()
    .reverse()
    .find((r) => r.meta && r.meta.title);

  // Find the nearest route element with meta tags.
  const nearestWithMeta = to.matched
    .slice()
    .reverse()
    .find((r) => r.meta && r.meta.metaTags);

  const previousNearestWithMeta = from.matched
    .slice()
    .reverse()
    .find((r) => r.meta && r.meta.metaTags);

  // If a route with a title was found, set the document (page) title to that value.
  if (nearestWithTitle) {
    document.title = nearestWithTitle.meta.title as string;
  } else if (previousNearestWithMeta) {
    document.title = previousNearestWithMeta.meta.title as string;
  }

  // Remove any stale meta tags from the document using the key attribute we set below.
  Array.from(document.querySelectorAll("[data-vue-router-controlled]")).map(
    (el) => el.parentNode!.removeChild(el)
  );

  // Skip rendering meta tags if there are none.
  if (!nearestWithMeta) return next();

  // Turn the meta tag definitions into actual elements in the head.
  (nearestWithMeta.meta as any).metaTags
    .map((tagDef: any) => {
      const tag = document.createElement("meta");

      Object.keys(tagDef).forEach((key) => {
        tag.setAttribute(key, tagDef[key]);
      });

      // We use this to track which meta tags we create so we don't interfere with other ones.
      tag.setAttribute("data-vue-router-controlled", "");

      return tag;
    })
    // Add the meta tags to the document head.
    .forEach((tag: any) => document.head.appendChild(tag));

  next();
});

import { useStore } from "@/store";
import { AuthGetterTypes } from "@/store/auth/getters";

function isAuthenticated(): boolean {
  return useStore().getters[AuthGetterTypes.isAuthenticated];
}

router.beforeEach((to, from, next) => {
  console.log("isAuthenticated() : " + isAuthenticated());
  if (to.name !== "Login" && !isAuthenticated()) {
    next({ name: "Login" });
  } else {
    next();
  }
});

export default router;
