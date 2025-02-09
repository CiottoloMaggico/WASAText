import {createRouter, createWebHistory} from 'vue-router'

const router = createRouter({
	history: createWebHistory(import.meta.env.BASE_URL),
	routes: [
		{
			path: '/',
			name: "homepage",
			component: () => import('@/views/home-page.vue'),
			children: [
				{
					path: "/profile",
					name: "profile",
					component: () => import("@/views/profile-page.vue"),
				},
				{
					path: "/conversations/:convId",
					name: "conversation",
					component: () => import("@/views/chat-page.vue"),
 				},
				{
					path: "conversations/:convId",
					name: "conversationInfo",
					component: () => import("@/views/conversation-detail-page.vue"),
				}
			]
		},
		{
			path: '/login',
			name: "loginPage",
			component: () => import("@/views/login-page.vue")
		},
	]
})

export default router
