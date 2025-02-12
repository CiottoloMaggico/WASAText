<script setup>
import {toRefs, ref, watch} from "vue"
import MessageService from "@/services/messageService";
import {getApiUrl} from "@/services/axios";
import {getAuthentication} from "@/services/sessionService";

const props = defineProps(["message", "show"])
const emits = defineEmits(["close"]);
const comments = ref([])

const {message, show} = toRefs(props)

watch(show, async (newVal, oldVal) => {
	if (!oldVal && newVal) {
		await getComments()
	}
})

async function getComments() {
	try {
		const data = await MessageService.getComments(message.value)
		comments.value = data
	} catch (e) {
		console.error(e)
	}
}

async function deleteComment() {
	try {
		await MessageService.uncommentMessage(message.value)
		await getComments()
	} catch (e) {
		console.error(e)
	}
}

function closeModal() {
	emits("close")
}

</script>

<template>
	<div
		class="modal fade"
		:ref="`comments-modal-${message.id}`"
		:id="`comments-modal-${message.id}`"
		tabindex="-1"
		aria-hidden="true"
		data-bs-backdrop="static"
		data-bs-keyboard="false"
	>
		<div class="modal-dialog modal-dialog-centered">
			<div class="modal-content">
				<div class="modal-header">
					<h5 class="modal-title">Comments</h5>
				</div>
				<div class="modal-body">
					<div v-for="comment in comments" :key="comment.author.uuid" class="comment-card">
						<div class="comment-content">
							<div class="avatar-box">
								<img class="avatar" :src="getApiUrl(comment.author.photo.fullUrl)"
									 :width="comment.author.photo.width"
									 :height="comment.author.photo.height" :alt="`${comment.author.username} comment`"/>
							</div>
							<span class="comment">
									{{ comment.author.username }} : {{ comment.content }}
							</span>
						</div>
						<div class="delete-comment-box" v-if="comment.author.uuid === getAuthentication()">
							<span class="btn btn-outline-danger" @click="deleteComment">Cancella</span>
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
	width: 100% ;
	height: 400px;
	padding: 0;
	overflow-y: scroll;
	overflow-x: hidden;
}

.comment-card {
	display: flex;
	height: 4rem;
	width: 100%;
	flex-flow: row nowrap;
	align-items: center;
	border-bottom: 1px solid #e4e4e4;
}

.comment-content {
	padding: .5rem;
	display: flex;
	width: 100%;
	height: 100%;
	flex-flow: row nowrap;
	align-items: center;
	gap: .5rem;
}

.avatar-box {
	display: flex;
	flex-shrink: 0;
	width: 3rem;
	height: 100%;
	flex-flow: row nowrap;
	align-items: center;
	gap: .5rem;
}

.avatar {
	width: 100%;
	height: 100%;
	border-radius: 50%;
	object-fit: cover;
}

.comment {
	font-size: 1.2rem;
}

.delete-comment-box {
	padding: 0 .5rem;
	height: 4rem;
	display: flex;
	flex-shrink: 0;
	justify-content: center;
	align-items: center;
}
</style>

