import {createRouter, createWebHistory} from 'vue-router'

const router = createRouter({
	history: createWebHistory(import.meta.env.BASE_URL),
	routes: [
		{path: '/', name: "conversationsPage", component: () => import('@/views/homepage.vue') },
		{path: '/login', name: "loginPage", component: () => import("@/views/login-page.vue")},
	]
})

export default router
