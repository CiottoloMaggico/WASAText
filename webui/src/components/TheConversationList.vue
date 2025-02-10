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
	<div class="sidebar-content">
		<div class="sidebar-header">
			<h3 class="header-title">Conversations</h3>
		</div>
		<div class="sidebar-body">
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
.sidebar-body {
	height: 100%;
	overflow: auto;
	scroll-behavior: smooth;
}


</style>
