import {defineStore} from "pinia"
import {useRoute} from "vue-router";

export const useConversationsStore = defineStore("conversationsStore", {
	state: () => ({
		route: useRoute(),
		pagination: {},
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
		getChatByRecipient: (state) => {
			return (recipient) => {return state.conversations.find((conversation) => conversation.name == recipient.username && conversation.type === 'chat')}
		},
		hasNext : (state) => {
			 return state.pagination.nextPage != null
		}
	},
	actions: {
		flush() {
			this.conversations = []
		},
		update(data) {
			this.pagination = data.page
			this.conversations = data.content
		},
		addPage(data) {
			Object.assign(this.pagination, data.page)
			this.conversations = this.conversations.concat(data.content)
		},
		updateConversation(newConversation) {
			let oldConversation = this.conversations.find((c) => c.id == newConversation.id)
			if (!oldConversation) {
				this.conversations.push(newConversation)
				return
			}
			Object.assign(oldConversation, newConversation)
		},
		removeConversation(conversationToRemove) {
			this.conversations.splice(this.conversations.findIndex((c) => c.id == conversationToRemove.id), 1)
		},
		updateLatestMessage(conversation, newMessage) {
			let oldConversation = this.conversations.find((c) => c.id == conversation.id)
			if (!oldConversation.latestMessage) {
				oldConversation.latestMessage = newMessage
			}
			else if (!newMessage) {
				oldConversation.latestMessage = null
			} else {
				oldConversation.latestMessage = {...newMessage}
			}
		},
	}
})
