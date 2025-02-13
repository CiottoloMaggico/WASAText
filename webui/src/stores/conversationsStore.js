import {defineStore} from "pinia"
import {useRoute} from "vue-router";

export const useConversationsStore = defineStore("conversationsStore", {
	state: () => ({
		route: useRoute(),
		conversations: [],
	}),
	getters: {
 		activeConversation: (state) => {
			let conversation = state.conversations.find((conversation) => conversation.id === state.route.params.convId)
			if (!state.route.params.convId || !conversation) {
				return null
			}

			return conversation
		},
		getStoredConversation: (state) => {
			return (conversationId) => {
				return state.conversations.find((conversation) => conversationId === conversation.id)
			}
		},
		getChatByRecipient: (state) => {
			return (recipient) => {return state.conversations.find((conversation) => conversation.name === recipient.username && conversation.type === 'chat')}
		},
	},
	actions: {
		flush() {
			this.conversations = []
		},
		updateConversation(newConversation) {
			let oldConversation = this.conversations.find((c) => c.id === newConversation.id)
			if (!oldConversation) {
				this.conversations.push(newConversation)
				return
			}
			Object.assign(oldConversation, newConversation)
		},
		updateConversations(conversations) {
			this.conversations = conversations
		},
		removeConversation(conversationToRemove) {
			this.conversations.splice(this.conversations.findIndex((c) => c.id === conversationToRemove.id), 1)
		},
		updateLatestMessage(conversation, newMessage) {
			let oldConversation = this.conversations.find((c) => c.id === conversation.id)
			if (!oldConversation.latestMessage) {
				oldConversation.latestMessage = newMessage
				return
			}

			oldConversation.latestMessage = {...newMessage}
		},
	}
})
