import { createRouter, createWebHistory } from "vue-router";
import { useUserStore } from "@/stores/user";

const routes = [
  {
    path: "/",
    name: "Home",
    component: () => import("@/views/Home.vue"),
  },
  {
    path: "/login",
    name: "Login",
    component: () => import("@/views/Login.vue"),
  },
  {
    path: "/register",
    name: "Register",
    component: () => import("@/views/Register.vue"),
  },
  {
    path: "/user/:id",
    name: "Profile",
    component: () => import("@/views/Profile.vue"),
  },
  {
    path: "/user/:id/badges",
    name: "UserBadges",
    component: () => import("@/views/UserBadges.vue"),
  },
  {
    path: "/topic/:id",
    name: "TopicDetail",
    component: () => import("@/views/TopicDetail.vue"),
  },
  {
    path: "/new-topic",
    name: "NewTopic",
    component: () => import("@/views/NewTopic.vue"),
    meta: { requiresAuth: true },
  },
  {
    path: "/search",
    name: "Search",
    component: () => import("@/views/Search.vue"),
  },
  {
    path: "/messages",
    name: "Messages",
    component: () => import("@/views/Messages.vue"),
    meta: { requiresAuth: true },
  },
  {
    path: "/notifications",
    name: "Notifications",
    component: () => import("@/views/Notifications.vue"),
    meta: { requiresAuth: true },
  },
  {
    path: "/favorites",
    name: "Favorites",
    component: () => import("@/views/Favorites.vue"),
    meta: { requiresAuth: true },
  },
  {
    path: "/user/:id/follows",
    name: "FollowList",
    component: () => import("@/views/FollowList.vue"),
  },
  {
    path: "/follow-topics",
    name: "FollowTopics",
    component: () => import("@/views/FollowTopics.vue"),
    meta: { requiresAuth: true },
  },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

router.beforeEach((to, from, next) => {
  const userStore = useUserStore();
  if (to.meta.requiresAuth && !userStore.isLoggedIn) {
    next("/login");
  } else {
    next();
  }
});

export default router;
