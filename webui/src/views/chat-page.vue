<script setup>
import {nextTick, reactive, ref, useTemplateRef, watchEffect, watch} from "vue"
import {MessageService} from "../services/messageService"
import {getAuthentication} from "../services/sessionService";
import TheMessage from "../components/TheMessage.vue";
import NewMessageBar from "../components/NewMessageBar.vue";
import {storeToRefs} from "pinia";
import {useConversationsStore} from "@/stores/conversationsStore";

const {activeConversation} = storeToRefs(useConversationsStore());

const messageContainer = useTemplateRef("message-container")

const newMessageReplyTo = ref(null)
const pagination = reactive({})
const messages = ref([])

watchEffect(async (onCleanup) => {
    if (activeConversation.value) {
        await getMessages()
	}

	onCleanup(() => {
		initializePage()
	})
})

watch(messages, () => {
	nextTick(() => {
		scrollToBottom()
	})
})

async function getMessages() {
	try {
        const data = await MessageService.setSeen(activeConversation.value)
		Object.assign(pagination, data.page)
		messages.value = data.content
	} catch (err) {
		console.error(err)
	}
}

async function loadNextPage() {
	try {
		const data = await MessageService.getMessages(
			activeConversation.value,
			{
				page: pagination.page+1,
				cursor: pagination.cursor,
			}
		)
		Object.assign(pagination, data.page)
		messages.value = messages.value.concat(data.content)
	} catch (err) {
		console.error(err)
	}
}

function initializePage() {
	messages.value = []
}

function scrollToBottom() {
	messageContainer.value.scrollTo({behavior: "instant", top: messageContainer.value.scrollHeight})
}

function updateReply(message) {
	newMessageReplyTo.value = message
}
</script>

<template>
	<div class="chat-component">
		<div class="header">
			<div class="title">
                {{ activeConversation.name }}
			</div>
            <router-link v-if="activeConversation.type === 'group'" class="info-icon-box"
                         :to="{ name: 'conversationInfo', params: { convId: activeConversation.id } }">
				<img class="info-icon" src="@/assets/images/information.png" width="512" height="512"/>
			</router-link>
		</div>
		<div class="body" ref="message-container">
			<the-message v-for="message in messages" :key="message.id" :message="message"
                         :is-author="message.author.uuid === getAuthentication()"
						 :message-container="messageContainer"
                         @update="getMessages" @want-reply="updateReply"/>
			<div v-if="pagination.nextPage" class="load-more-box">
				<div class="btn btn-outline-primary" @click="loadNextPage">Load more</div>
			</div>
		</div>
		<div class="footer">
            <new-message-bar @update="getMessages" @clear-reply="newMessageReplyTo = null" :conversation="activeConversation" :want-reply="newMessageReplyTo"/>
		</div>
	</div>
</template>

<style scoped>
.chat-component {
	align-items: center;
	display: flex;
	flex-flow: column nowrap;
	height: 100%;
	justify-content: space-between;
	width: 100%;
	flex-grow: 0;
}

.header, .body, .footer {
	width: 100%;
}

.header, .footer {
	flex-shrink: 0;
	background: var(--TERTIARY-COLOR);
	padding: 0 1rem;
}

.header {
	align-items: center;
	border-top-right-radius: var(--MAIN-BORDER-RADIUS);
	display: flex;
	flex-flow: row nowrap;
	height: 4rem;
}

.title {
	font-size: 1.5rem;
	font-weight: bolder;
	width: 100%;
}

.info-icon-box {
	width: 2.5rem;
	height: 2.5rem;
	flex-shrink: 0;
}

.info-icon {
	width: 100%;
	height: 100%;
}

.body {
	background: var(--SECONDARY-COLOR);
	height: 100%;
	overflow-y: scroll;
	overflow-x: hidden;
    display: flex;
    flex-flow: column-reverse nowrap;
	padding-top: 1rem;
}

.load-more-box {
	display: flex;
	width: 100%;
	height: 5rem;
	justify-content: center;
	align-items: center;
	flex-shrink: 0;
}

.footer {
	border-bottom-right-radius: var(--MAIN-BORDER-RADIUS);
	display: flex;
	flex-flow: column nowrap;
}
</style>
