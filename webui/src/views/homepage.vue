<script setup>
import ConversationsBar from "../components/ConversationsBar.vue";
import ChatComponent from "../components/ChatComponent.vue";
import {reactive} from "vue";


const currentConversation = reactive({})

function changeConversation(conversation) {
	Object.assign(currentConversation, conversation)
}

</script>

<template>
	<main>
		<conversations-bar @change-active-conversation="changeConversation"/>
		<chat-component :conversation="currentConversation" />
	</main>
</template>

<style>
main {
	display: flex;
	flex-flow: row nowrap;
}

h1 {
	color: black;
}
</style>
