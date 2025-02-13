import api from "./axios";
import ConversationService from "@/services/conversationService";
import UserService from "@/services/userService";
import {useProfileStore} from "@/stores/profileStore";
import {useConversationsStore} from "@/stores/conversationsStore";
import router from "@/router";

const API_AUTHENTICATION_KEY = "auth-key"

function setAuthentication(token) {
	if (token) {
		localStorage.setItem(API_AUTHENTICATION_KEY, token)
	} else {
		localStorage.removeItem(API_AUTHENTICATION_KEY)
	}
}

export function getAuthentication() {
	const token = localStorage.getItem(API_AUTHENTICATION_KEY)
	return token
}

export function isAuthed() {
	if (!getAuthentication()) {
		return false
	}
	return true
}


export const SessionService = Object.freeze({
	get profileStore() {
		return useProfileStore()
	},
	get conversationsStore() {
		return useConversationsStore()
	},
	async refresh() {
		await UserService.refresh()
		await ConversationService.refresh()
	},
	async doLogin(username){
		const response = await api.post("/session", {username: username})

		if (response.status !== 200) {
			throw new Error(response.statusText)
		}
		setAuthentication(response.data.uuid)
		this.profileStore.update(response.data)
		await ConversationService.refresh()
		return response.data
	},
	async logout() {
		setAuthentication(null)
		await router.replace({name: "login"})
		this.profileStore.flush()
		this.conversationsStore.flush()
	}
})

export default SessionService
