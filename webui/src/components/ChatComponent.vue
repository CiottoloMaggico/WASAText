<script setup>
import {ref, watch, watchEffect, useTemplateRef, nextTick, onMounted, reactive, computed} from "vue"
import {MessageService} from "../services/message"
import MessageComponent from "./MessageComponent.vue";
import {getAuthentication} from "../services/session";
import NewMessageBar from "./NewMessageBar.vue";

const props = defineProps({
	conversation: Object,
})

const messages = ref([])
const messageContainer = useTemplateRef("message-container")
const newMessage = reactive({
	content: null,
	attachment: null,
	repliedMessage: null,
})

onMounted(() => {
	scrollToBottom()
})

watchEffect(async (onCleanup) => {
	if (props.conversation) {
		await getMessages(props.conversation.id)

	}

	onCleanup(() => {
		initializePage()
	})
})

watch([messages, () => newMessage.attachment], () => {
	nextTick(() => {
		scrollToBottom()
	})
})

async function getMessages(conversationId) {
	try {
		const response = await MessageService.getMessages(conversationId)
		messages.value = response.data.content
	} catch (err) {
		console.log(err.toString())
	}
}

async function sendMessage() {
	try {
		const response = await MessageService.sendMessage(props.conversation.id, newMessage)
	} catch (err) {
		console.log(err.toString())
	} finally {
		clearNewMessage()
		await getMessages(props.conversation.id)
	}
}

async function deleteMessage(message) {
	try {
		const response = await MessageService.deleteMessage(props.conversation.id, message.id)
	} catch (err) {
		console.log(err.toString())
	} finally {
		await getMessages(props.conversation.id)
	}
}

function replyTo(message) {
	newMessage.repliedMessage = message
}

function initializePage() {
	messages.value = []
	clearNewMessage()
}

function clearNewMessage() {
	Object.assign(newMessage, {
		content: null,
		attachment: null,
		repliedMessage: null,
	})
}

function scrollToBottom() {
	messageContainer.value.scrollTo({behavior: "instant", top: messageContainer.value.scrollHeight})
}
</script>

<template>
	<div class="chat-component">
		<div class="header">
			<div class="title">
				{{ conversation.name }}
			</div>
			<div class="info-icon-box">
				<img class="info-icon" src="@/assets/images/information.png" width="512" height="512"/>
			</div>
		</div>
		<div class="body" ref="message-container">
			<message-component v-for="message in messages" :key="message.id" :message="message"
							   :is-author="message.author.uuid == getAuthentication()" @reply="replyTo"
							   @delete="deleteMessage"/>
		</div>
		<div class="footer">
			<new-message-bar :new-message="newMessage" @sendMessage="sendMessage"/>
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
	width: 80%;
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
	padding-top: 2rem;
}

.footer {
	border-bottom-right-radius: var(--MAIN-BORDER-RADIUS);
	display: flex;
	flex-flow: column nowrap;
}
</style>
