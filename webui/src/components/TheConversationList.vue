<script setup>
import TheConversation from "@/components/TheConversation.vue";
import {useConversationsStore} from "@/stores/conversationsStore";
import {storeToRefs} from "pinia";

const {conversations, activeConversation} = storeToRefs(useConversationsStore())

function isSelected(conversation) {
	return (activeConversation.value && conversation.id === activeConversation.value.id)
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
		</div>
	</div>
</template>

<style scoped>
.sidebar-body {
	overflow-y: scroll !important;
}

.sidebar-item {
	border-bottom: 1px solid #e4e4e4;
}

.selected {
	border-bottom: none;
	background: var(--SECONDARY-COLOR);
}


</style>
