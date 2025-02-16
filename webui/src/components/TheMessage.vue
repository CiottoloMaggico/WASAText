<script setup>
import {getApiUrl} from "../services/axios";
import {watch, computed, onBeforeMount, reactive, ref, useTemplateRef} from "vue";
import EmojiPicker from 'vue3-emoji-picker'
import 'vue3-emoji-picker/css'
import MessageService from "@/services/messageService";
import ShowCommentsModal from "@/components/ShowCommentsModal.vue";
import ForwardingModal from "@/components/ForwardingModal.vue";

const props = defineProps({
	message: Object,
	isAuthor: Boolean,
	messageContainer: Object
})
const emits = defineEmits(["update", "wantReply"])

const emojiPickerToggle = useTemplateRef("emoji-picker-toggle")
const emojiPickerCoords = ref({
	top: '0',
	left: '0',
})
const showEmojiPicker = ref(false)
const showCommentsModal = ref(false)
const showForwardModal = ref(false)
const comments = ref([])


const repliedMessage = reactive({})

const sendAt = computed(() => {
	return new Date(props.message.sendAt).toLocaleString([], {dateStyle: 'short', timeStyle: 'short'})
})

onBeforeMount(async () => {
	await getComments()
	if (!props.message.repliedMessageId) {
		return
	}
	await getRepliedMessage()
})

watch(showCommentsModal, async (newVal, oldVal) => {
	if (!oldVal && newVal) {
		await getComments()
	}
})


async function deleteMessage() {
	try {
		await MessageService.deleteMessage(props.message)
	} catch (err) {
		console.error(err.toString())
	} finally {
		emits("update")
	}
}

async function getRepliedMessage() {
	try {
		const data = await MessageService.getMessage(props.message.conversationId, props.message.repliedMessageId)
		Object.assign(repliedMessage, data)
	} catch (err) {
		console.error(err.toString())
	}
}

async function getComments() {
	try {
		const data = await MessageService.getComments(props.message)
		comments.value = data
	} catch (e) {
		console.error(e)
	}
}

async function setComment(emoji) {
	try {
		await MessageService.commentMessage(props.message, emoji.i)
		await getComments()
	} catch (e) {
		console.error(e.toString())
	} finally {
		closeEmojiPicker()
	}
}


async function deleteComment() {
	try {
		await MessageService.uncommentMessage(props.message)
		await getComments()
	} catch (e) {
		console.error(e)
	}
}


function closeEmojiPicker() {
	if (showEmojiPicker.value) {
		showEmojiPicker.value = false
	}
}

function emojiPickerPosition() {
	if (!emojiPickerToggle.value) {
		return
	}
	let {height, y} = props.messageContainer.getBoundingClientRect()
	let halfContainerHeight = height / 2
	let emojiPickerParent = emojiPickerToggle.value.getBoundingClientRect()
	let toggleY = emojiPickerParent.y
	let top = emojiPickerParent.top + 24, left = emojiPickerParent.left + 24

	if (toggleY > y + halfContainerHeight) {
		top = emojiPickerParent.top - 320
	}

	if (props.isAuthor) {
		left = emojiPickerParent.left - 280
	}

	emojiPickerCoords.value.top = `${top}px`
	emojiPickerCoords.value.left = `${left}px`
	showEmojiPicker.value = !showEmojiPicker.value
}
</script>

<template>
	<div class="message-row">
		<div v-if="!isAuthor" class="header">
			<div class="sender-avatar-box">
				<img class="sender-avatar" :src="getApiUrl(message.author.photo.fullUrl)"
					 :width="message.author.photo.width" :height="message.author.photo.height"/>
			</div>
			<div class="sender-name-box">
				<span class="sender-name">{{ message.author.username }}</span>
			</div>
		</div>
		<div class="message-box" :class="{'flex-row-reverse': isAuthor}">
			<div class="message-side" :class="{'align-items-end' : isAuthor}">
				<div class="message" :class="{'flex-row-reverse': isAuthor}">
					<div class="body">
						<div v-if="repliedMessage.id" class="replied-message-box">
							<div class="replied-message-content">
								<div class="replied-message-author">
									{{ repliedMessage.author.username }}
								</div>
								<div class="replied-message-text-box">
									<span v-if="repliedMessage.content"
										  class="replied-message-text">
										{{ repliedMessage.content }}
									</span>
								</div>
							</div>
							<span v-if="repliedMessage.attachment" class="replied-message-image-box">
									<img class="replied-message-image"
										 :src="getApiUrl(repliedMessage.attachment.fullUrl)"/>
								</span>
						</div>
						<div class="content">
							<div v-if="message.attachment !== null" class="message-image-box"
								 :class="{'mb-2' : message.content !== null}">
								<img class="message-image" :src="getApiUrl(message.attachment.fullUrl)"
									 :width="message.attachment.fullUrl.width" :height="message.attachment.fullUrl.height"/>
							</div>
							<span v-if="message.content !== null" class="message-text">
							{{ message.content }}
						</span>
						</div>
						<div class="info-box">
							<span class="send-at">
								{{ sendAt }}
							</span>
							<span class="checkmark-box" v-if="isAuthor">
							<img v-if="message.status === 'delivered'" class="checkmark"
								 src="@/assets/images/Checkmark.png" width="512" height="512"/>
							<img v-else-if="message.status === 'seen'" class="checkmark"
								 src="@/assets/images/seen.png" width="512" height="512"/>
							</span>
						</div>
					</div>
					<span class="comment-btn" ref="emoji-picker-toggle" v-click-outside="closeEmojiPicker">
						<img class="svg-icon" src="@/assets/images/emoji.svg"
							 alt="add comment to the message"
							 @click="emojiPickerPosition()"/>
						<div class="emoji-picker">
							<emoji-picker v-show="showEmojiPicker"
										  :native="true" @select="setComment"/>

						</div>
				</span>
				</div>
				<div v-if="comments.length !== 0" class="reactions-container" :data-bs-target="`#comments-modal-${message.id}`" data-bs-toggle="modal"
					 @click="showCommentsModal = true">
						<span class="reaction" v-for="i in Math.min(3, comments.length)">
							<span>
								{{comments[i-1].content}}
							</span>
						</span>
				</div>
			</div>
			<div class="options-box">
				<div class="option" @click="$emit('wantReply', message)">
					<img class="option-icon" src="@/assets/images/reply.png" alt=""/>
				</div>
				<div class="option" @click="deleteMessage">
					<img class="option-icon" src="@/assets/images/bin.png" alt=""/>
				</div>
				<div class="option" :data-bs-target="`#forward-modal-${message.id}`" data-bs-toggle="modal"
					 @click="showForwardModal = true">
					<img class="option-icon" src="@/assets/images/forward.png" alt=""/>
				</div>
			</div>
		</div>
	</div>
	<show-comments-modal :message="message" :comments="comments" @close="showCommentsModal = false" @delete-comment="deleteComment"/>
	<forwarding-modal :message="message" :show="showForwardModal" @close="showForwardModal = false"/>
</template>

<style scoped>
.svg-icon {
	width: 100%;
	height: 100%;
}

.reactions-container {
	top: -6px;
	position: relative;
	display: flex;
	flex-flow: row nowrap;
	overflow: hidden;
	background-color: var(--PRIMARY-COLOR);
	border-radius: 1rem;
	width: fit-content;
	border: 1px solid #3d6e53;
	padding: 1px;
}

.reaction {
	width: 1.5rem;
	height: 1.5rem;
	display: flex;
	justify-content: center;
	align-items: center;
	flex-shrink: 0;
}

.comment-btn {
	width: 1.5rem;
	height: 1.5rem;
	align-self: center;
	margin-right: 5px;
	flex-shrink: 0;
}

.emoji-picker {
	position: fixed;
	top: v-bind('emojiPickerCoords.top');
	left: v-bind('emojiPickerCoords.left');
	z-index: 1;
}

.author-message {
	flex-direction: row-reverse !important;
}

.message-row {
	flex-shrink: 0;
	display: flex;
	flex-flow: column nowrap;
	width: 100%;
	min-height: 1rem;
	margin-bottom: 10px;
	padding: 0 1rem;
}

.options-box {
	display: flex;
	flex-shrink: 0;
	flex-flow: row nowrap;
	align-items: center;
	height: 100%;
	gap: .7rem;
}

.option {
	width: 1.5rem;
	height: 1.5rem;
	outline: 1px solid black;
	border-radius: 50%;
	cursor: pointer;
}

.option-icon {
	width: 100%;
	height: 100%;
	scale: 0.8;
}

.message-box {
	display: flex;
	flex-flow: row nowrap;
	height: 100%;
	width: 100%;
	justify-content: space-between;
}

.message-side {
	display: flex;
	height: 100%;
	flex-flow: column nowrap;
	justify-content: start;
	overflow: hidden;
	padding: 5px;
	max-width: 60%;
}

.author-message-side {
	justify-content: end;
	flex-direction: row-reverse;
}

.message {
	display: flex;
	flex-flow: row nowrap;
	flex-shrink: 0;
	gap: 1rem;
	width: 100%;
}

.header {
	align-items: center;
	display: flex;
	flex-flow: row nowrap;
	margin-bottom: 0.5rem;
}

.sender-avatar-box {
	flex-shrink: 0;
	height: 3rem;
	width: 3rem;
}

.sender-avatar {
	border-radius: 50%;
	height: 100%;
	width: 100%;
}

.sender-name-box {
	display: flex;
	margin-left: 10px;
	overflow: hidden;
	width: 100%;
}

.sender-name {
	font-size: 1.2rem;
	overflow: hidden;
	text-overflow: ellipsis;
	white-space: nowrap;
	width: 100%;
}

.body {
	background-color: var(--TERTIARY-COLOR);
	border-radius: var(--MAIN-BORDER-RADIUS);
	min-height: 2rem;
	padding: .7rem;
	width: fit-content;
	max-width: 100%;
	overflow: hidden;
}

.content {
	margin-bottom: 5px;
}

.message-text {
	word-break: break-all;
	overflow: hidden;
	display: inline-block;
	font-size: 1.2rem;
	width: 100%;
}

.message-image-box {
	max-width: 100%;
	overflow: hidden;
	background-color: #fff;
	border-radius: 4px;
}

.message-image {
	width: 100%;
}

.info-box {
	display: flex;
	flex-flow: row nowrap;
	align-items: center;
	width: 100%;
	justify-content: space-between;
	gap: 10px;
}

.send-at {
	font-size: .7rem;
}

.checkmark-box {
	width: 1.1rem;
	height: 1.1rem;
}

.checkmark {
	width: 100%;
	height: 100%;
}

.replied-message-box {
	width: 100%;
	height: 4rem;
	background-color: #aaa;
	border-radius: 5px;
	padding: 5px;
	margin-bottom: 5px;
	max-width: 100%;
	display: flex;
	flex-flow: row nowrap;
	align-items: center;
	gap: 10px;
	overflow: hidden;
}

.replied-message-content {
	display: flex;
	flex-flow: column nowrap;
	justify-content: space-between;
	width: 100%;
	overflow: hidden;
	height: 100%;
	flex-grow: 0;
}

.replied-message-author {
	font-size: .9rem;
	flex-shrink: 0;
	overflow: hidden;
	text-overflow: ellipsis;
	white-space: nowrap;
}


.replied-message-text-box {
	width: 100%;
	height: 100%;
	padding: 5px 0;
}

.replied-message-text {
	width: 100%;
	font-size: .8rem;
	display: inline-block;
	overflow: hidden;
	text-overflow: ellipsis;
	white-space: nowrap;
}

.replied-message-image-box {
	width: 3rem;
	height: 3rem;
	flex-shrink: 0;
}

.replied-message-image {
	width: 100%;
	height: 100%;
	border-radius: 5px;
	object-fit: cover;
}

</style>
