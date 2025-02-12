import {createRouter, createWebHashHistory} from 'vue-router'
import {isAuthed, SessionService} from "@/services/sessionService";
import {useProfileStore} from "@/stores/profileStore";
import {useConversationsStore} from "@/stores/conversationsStore";
import ConversationService from "@/services/conversationService";

const router = createRouter({
	history: createWebHashHistory(import.meta.env.BASE_URL),
	routes: [
		{
			path: '/',
			name: "homepage",
			component: () => import('@/views/home-page.vue'),
			beforeEnter: async () => {
				if (!isAuthed()) return {name: "login"}
				let profileStore = useProfileStore()
				if (!profileStore.getProfile) await SessionService.refresh()
			},
			children: [
				{
					path: "/conversations/:convId",
					name: "conversation",
					component: () => import("@/views/chat-page.vue"),
					beforeEnter: async (to) => {
						let conversationStore = useConversationsStore()
						if (!conversationStore.activeConversation) await ConversationService.getConversation(to.params.convId)
					}
 				},
				{
					path: "conversations/:convId/detail",
					name: "conversationInfo",
					component: () => import("@/views/conversation-detail-page.vue"),
				},
			],
		},
		{
			path: '/login',
			name: "login",
			component: () => import("@/views/login-page.vue"),
			beforeEnter: () => {
				if (isAuthed()) {
					return {name: "homepage"}
				}
			},
		},
	]
})

export default router
