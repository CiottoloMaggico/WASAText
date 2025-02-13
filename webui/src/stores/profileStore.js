import {defineStore} from "pinia"

export const useProfileStore = defineStore("profileStore", {
	state: () => ({
		profile: null,
	}),
	getters: {
		getProfile: (state) => {
			return state.profile
		},
	},
	actions: {
		update(newProfile) {
			this.profile = {...newProfile}
		},
		flush() {
			this.profile = null
		},
	}
})
