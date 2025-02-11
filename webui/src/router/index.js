import {createRouter, createWebHistory} from 'vue-router'
import {isAuthed} from "@/services/session";
import {useProfileStore} from "@/stores/profileStore";

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
				},
			]
		},
		{
			path: '/login',
			name: "login",
			component: () => import("@/views/login-page.vue"),
			beforeEnter: (to, from) => {
				if (isAuthed()) {
					return false
				}
			},
		},
	]
})

router.beforeEach(
	async (to, from) => {
		if (!isAuthed() && to.name !== 'login') {
			return { name: "login" }
		}
		let profileStore = useProfileStore()
		if (useProfileStore.getProfile == null) {
			await profileStore.refreshProfile()
		}
	}
)

export default router
