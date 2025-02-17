<script setup>
import TheConversation from "@/components/TheConversation.vue";
import {useConversationsStore} from "@/stores/conversationsStore";
import {storeToRefs} from "pinia";
import ConversationService from "@/services/conversationService";

const {conversations, activeConversation, hasNext} = storeToRefs(useConversationsStore())

function isSelected(conversation) {
	return (activeConversation.value && conversation.id === activeConversation.value.id)
}

async function loadNextPage() {
	try {
		await ConversationService.getNextPage()
	} catch (e) {
		console.error(e.toString())
	}
}
</script>

<template>
	<div class="sidebar-content">
		<div class="sidebar-header">
			<h3 class="header-title">Conversations</h3>
		</div>
		<div class="sidebar-body">
			<router-link v-for="conversation in conversations" :key="conversation.id"
						 :to="{name: 'conversation', params: {convId: conversation.id}}"
						 :class="{'selected' : isSelected(conversation)}"
						 class="sidebar-item">
				<the-conversation :conversation="conversation"/>
			</router-link>
			<div v-if="hasNext" class="sidebar-item load-more-box">
				<div class="btn btn-outline-primary" @click="loadNextPage">Load more</div>
			</div>
		</div>
	</div>
</template>

<style scoped>
.sidebar-body {
	overflow-y: scroll !important;
}

.sidebar-item {
	border-bottom: 1px solid #e4e4e4;
	width: 100%;
}

.selected {
	border-bottom: none;
	background: var(--SECONDARY-COLOR);
}

.load-more-box {
	padding: 0 .5rem;
}


</style>
