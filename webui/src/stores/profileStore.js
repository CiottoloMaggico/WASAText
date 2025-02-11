import {defineStore} from "pinia"
import SessionService from "@/services/sessionService";
import UserService from "@/services/userService";
import UserConversationService from "@/services/userConversationService";
import {useRoute} from "vue-router";

export const useProfileStore = defineStore("profileStore", {
	state: () => ({
		route: useRoute(),
		profile: null,
		conversations: [],
	}),
	getters: {
		activeConversation: (state) => {
			let conversation = state.conversations.find((conversation) => conversation.id == state.route.params.convId)
			if (!state.route.params.convId || !conversation) {
				return null
			}

			return conversation
		},
		getProfile: (state) => {
			return state.profile
		},
	},
	actions: {
		async login(username) {
			try {
				let response = await SessionService.doLogin(username)
				this.profile = response.data

				response = await UserConversationService.setDelivered()
				this.conversations = response.data.content
			} catch (error) {
				console.error(error.toString())
				return error
			}
		},
		async refreshProfile() {
			try {
				let response = await UserService.getProfile()
				this.profile = response.data

				response = await UserConversationService.setDelivered()
				this.conversations = response.data.content
			} catch (error) {
				console.error(error.toString())
				return error
			}
		},
		async changeUsername(newUsername) {
			try {
				const response = await UserService.setMyUsername(newUsername)
				this.profile = response.data

			} catch (error) {
				console.error(error.toString())
				return error
			}
		},
		async changeAvatar(newAvatar) {
			try {
				const response = await UserService.setMyPhoto(newAvatar)
				this.profile = response.data

			} catch (err) {
				console.log(err.toString())
				return error
			}
		},
		async getConversations() {
			try {
				const response = await UserConversationService.getConversations({})
				this.conversations = response.data.content
			} catch (err) {
				console.log(err.toString())
				return error
			}
		},
	}
})
