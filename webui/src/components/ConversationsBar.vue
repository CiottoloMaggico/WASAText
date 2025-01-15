<script setup>
import {ref, reactive} from "vue"
import UserConversationService from "../services/userConversation"
import ConversationCard from "./ConversationCard.vue";


const emit = defineEmits(["changeActiveConversation"])
const conversations = ref([])
const activeConversation = reactive({})

await getConversations()

async function getConversations() {
	try {
		const response = await UserConversationService.getConversations()
		conversations.value = response.data.content
	} catch (error) {
		console.error(error.toString())
	}
}

function selectConversation(conversation) {
	Object.assign(activeConversation, conversation)
	emit("changeActiveConversation", activeConversation)
}

</script>

<template>
	<div class="conversation-bar">
		<div class="conversation-header">
			<h3 class="header-title">Conversations</h3>
		</div>
		<div class="conversation-body">
			<conversation-card v-for="conversation in conversations" :key="conversation.id" @click="selectConversation(conversation)"
							   :conversation="conversation" :selected="conversation.id == activeConversation.id"/>
		</div>
		<div class="conversation-footer">
			<div class="add-conversation-btn primary-btn">Create chat</div>
			<div class="add-conversation-btn primary-btn">Create group</div>
		</div>
	</div>
</template>

<style scoped>
.conversation-bar {
	background-color: var(--TERTIARY-COLOR);
	border-radius: var(--MAIN-BORDER-RADIUS) 0 0 var(--MAIN-BORDER-RADIUS);
	display: flex;
	flex-flow: column nowrap;
	height: 100%;
	width: 20%;
	border-right: 1px solid #000
}

.conversation-header {
	align-items: center;
	display: flex;
	height: 4rem;
	padding-left: 10px;
}

.header-title {
	font-weight: bolder;
	margin-bottom: 0;
}

.conversation-body {
	height: 100%;
	overflow: auto;
	scroll-behavior: smooth;
}

.conversation-footer {
	align-items: center;
	display: flex;
	flex-flow: row nowrap;
	height: 4rem;
	justify-content: space-around;
	padding-left: 10px;
	width: 100%;
}

.add-conversation-btn {
	border-radius: var(--MAIN-BORDER-RADIUS);
	padding: 5px;
	text-align: center;
	text-wrap: nowrap;
	width: 125px;
}
</style>
