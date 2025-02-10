<script setup>
import {reactive} from "vue"
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
	<div class="sidebar-content">
		<div class="sidebar-header">
			<h3 class="header-title">Conversations</h3>
		</div>
		<div class="sidebar-body">
			<router-link v-for="conversation in conversations" :key="conversation.id"
						 :to="{name: 'conversation', params: {convId: conversation.id}}"
						 :class="(conversation.id === activeConversation.id) ? 'selected' : ''" class="sidebar-item"
						 @click="selectConversation(conversation)">
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
