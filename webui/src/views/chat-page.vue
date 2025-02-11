<script setup>
import {nextTick, reactive, ref, useTemplateRef, watch, watchEffect} from "vue"
import {MessageService} from "../services/messageService"
import {getAuthentication} from "../services/sessionService";
import TheMessage from "../components/TheMessage.vue";
import NewMessageBar from "../components/NewMessageBar.vue";
import {useProfileStore} from "@/stores/profileStore";
import {storeToRefs} from "pinia";

const profileStore = useProfileStore()
const {activeConversation} = storeToRefs(profileStore);

const messageContainer = useTemplateRef("message-container")

const newMessageReplyTo = ref(null)
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
        const response = await MessageService.setSeen(activeConversation.value)
		messages.value = response.data.content
	} catch (err) {
		console.log(err.toString())
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
}

.footer {
	border-bottom-right-radius: var(--MAIN-BORDER-RADIUS);
	display: flex;
	flex-flow: column nowrap;
}
</style>
