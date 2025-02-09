<script setup>
import {ref, reactive} from "vue"
import UserConversationService from "../services/userConversation"
import TheConversation from "./TheConversation.vue";
import {storeToRefs} from "pinia";
import {useProfileStore} from "@/stores/profileStore";

const profileStore = useProfileStore()
const {conversations} = storeToRefs(profileStore);
const activeConversation = reactive({})

function selectConversation(conversation) {
	Object.assign(activeConversation, conversation)
}
</script>

<template>
	<div class="conversations-bar">
		<div class="conversations-header">
			<h3 class="header-title">Conversations</h3>
			<div class="actions-box">
				<div class="btn btn-outline-primary">Create chat</div>
				<div class="btn btn-outline-primary">Create group</div>
			</div>
		</div>
		<div class="conversation-body">
			<router-link v-for="conversation in conversations" :key="conversation.id"
						 :to="{name: 'conversation', params: {convId: conversation.id}}"
						 @click="selectConversation(conversation)">
				<the-conversation :conversation="conversation"
								  :selected="conversation.id == activeConversation.id"/>
			</router-link>
		</div>
	</div>
</template>

<style scoped>
.conversations-bar {
	display: flex;
	flex-flow: column nowrap;
	height: 100%;
	width: 100%;
}

.conversations-header {
	display: flex;
	flex-flow: column nowrap;
	padding: 1rem;
	gap: .5rem;
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

.actions-box {
	align-items: center;
	display: flex;
	flex-flow: row nowrap;
	height: 4rem;
	justify-content: space-around;
	width: 100%;
}
</style>
