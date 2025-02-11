<script setup>
import {ref, toRefs, watch} from "vue"
import MessageService from "@/services/messageService";
import {getApiUrl} from "@/services/axios";
import UserConversationService from "@/services/userConversationService";
import router from "@/router";

const props = defineProps(["message", "show"])
const emits = defineEmits(["close"]);

const conversations = ref([])

const {message, show} = toRefs(props)

watch(show, async (newVal, oldVal) => {
	if (!oldVal && newVal) {
		await getConversations()
	}
})

async function getConversations() {
	try {
		const response = await UserConversationService.getConversations(
			{filter: `id ne ${message.value.conversationId}`}
		)
		conversations.value = response.data.content
	} catch (e) {
		console.error(e)
	}
}

async function forwardMessage(destConversation) {
	try {
		await MessageService.forwardMessage(message.value, destConversation)
		await router.push({name: "conversation", params: {convId: destConversation.id}})
	} catch (e) {
		console.error(e)
	}
	closeModal()
}

function closeModal() {
	emits("close")
}
</script>

<template>
	<div
		class="modal fade"
		:ref="`forward-modal-${message.id}`"
		:id="`forward-modal-${message.id}`"
		tabindex="-1"
		aria-hidden="true"
		data-bs-backdrop="static"
		data-bs-keyboard="false"
	>
		<div class="modal-dialog modal-dialog-centered">
			<div class="modal-content">
				<div class="modal-header">
					<h5 class="modal-title">Forward to</h5>
				</div>
				<div class="modal-body">
					<div v-for="conversation in conversations" :key="conversation.id" class="list-card"
						 data-bs-dismiss="modal" @click="forwardMessage(conversation)">
						<div class="card-content">
							<div class="img-box">
								<img class="img" :src="getApiUrl(conversation.image.fullUrl)"
									 :width="conversation.image.width"
									 :height="conversation.image.height" :alt="`${conversation.name}`"/>
							</div>
							<span class="card-title">
									{{ conversation.name }}
							</span>
						</div>
					</div>
				</div>
				<div class="modal-footer">
					<button type="button" class="btn btn-secondary rounded-pill" data-bs-dismiss="modal"
							@click="closeModal">Close
					</button>
				</div>
			</div>
		</div>
	</div>
</template>

<style scoped>
.modal-body {
	width: 100%;
	height: 400px;
	padding: 0;
	overflow-y: scroll;
	overflow-x: hidden;
}

.list-card {
	display: flex;
	height: 4rem;
	width: 100%;
	flex-flow: row nowrap;
	align-items: center;
	border-bottom: 1px solid #e4e4e4;

	&:hover{
		cursor: pointer;
		background: var(--SECONDARY-COLOR);
		border: none;
	}
}

.card-content {
	padding: .5rem;
	display: flex;
	width: 100%;
	height: 100%;
	flex-flow: row nowrap;
	align-items: center;
	gap: .5rem;
}

.img-box {
	display: flex;
	flex-shrink: 0;
	width: 3rem;
	height: 100%;
	flex-flow: row nowrap;
	align-items: center;
	gap: .5rem;
}

.img {
	width: 100%;
	height: 100%;
	border-radius: 50%;
	object-fit: cover;
}

.card-title {
	font-size: 1.2rem;
}
</style>

