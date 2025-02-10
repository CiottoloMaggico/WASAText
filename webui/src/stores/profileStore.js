import {defineStore} from "pinia"
import SessionService from "@/services/session";
import UserService from "@/services/userService";
import UserConversationService from "@/services/userConversation";

export const useProfileStore = defineStore("profileStore", {
	state: () => ({
		profile: null,
		conversations: [],
	}),
	getters: {
		getProfile: (state) => {
			return state.profile
		},
	},
	actions: {
		async login(username) {
			try {
				const response = await SessionService.doLogin(username)
				this.profile = response.data
				let err = await this.getConversations()
			} catch (error) {
				console.error(error.toString())
				return error
			}
		},
		async refreshProfile() {
			try {
				const response = await UserService.getProfile()
				this.profile = response.data
				let err = await this.getConversations()
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
				const response = await UserConversationService.getConversations()
				this.conversations = response.data.content
			} catch (err) {
				console.log(err.toString())
				return error
			}
		}
	}
})
