<script setup>
import ConversationService from "@/services/conversationService";
import router from "../router";
import TheConversationList from "@/components/TheConversationList.vue";
import TheNewGroup from "@/components/TheNewGroup.vue";
import AddParticipantForm from "@/components/AddParticipantForm.vue";

const emits = defineEmits(["switch"])

async function createChat(recipient) {
	let chat
	try {
		chat = await ConversationService.createChat(recipient)
	} catch (error) {
		console.error(error)
		return error
	}
	await router.push({name: "conversation", params: {convId: chat.id}})
	emits("switch", TheConversationList.__name)
}

</script>

<template>
	<div class="sidebar-content">
		<div class="sidebar-header">
			<h3 class="header-title">Add conversation</h3>
		</div>
		<div class="sidebar-body">
			<div class="sidebar-group" @click="emits('switch', TheNewGroup.__name)">
				<div class="sidebar-item">
					<span class="sidebar-item-title">New group</span>
				</div>
			</div>
			<AddParticipantForm :participants="[]" :single-mode="true" @add-participant="createChat"/>
		</div>
	</div>
</template>

<style scoped>
.sidebar-header, .sidebar-group, .sidebar-item {
	border-bottom: 1px solid #e4e4e4;
}

.sidebar-item {
	padding: 0 1rem;
}

</style>
