<script setup>
import {getApiUrl} from "../services/axios";
import {computed, ref, useTemplateRef} from "vue";
import EmojiPicker from 'vue3-emoji-picker'
import 'vue3-emoji-picker/css'
import MessageService from "@/services/messageService";
import ShowCommentsModal from "@/components/ShowCommentsModal.vue";
import ForwardingModal from "@/components/ForwardingModal.vue";
import {isAuthed} from "@/services/sessionService";

const props = defineProps({
	message: Object,
	isAuthor: Boolean,
	messageContainer: Object
})
const emits = defineEmits(["update", "wantReply"])

const emojiPickerToggle = useTemplateRef("emoji-picker-toggle");
const showEmojiPicker = ref(false)
const showCommentsModal = ref(false)
const showForwardModal = ref(false)

const sendAt = computed(() => {
	return new Date(props.message.sendAt).toLocaleString([], {dateStyle: 'short', timeStyle: 'short'})
})


async function setComment(emoji) {
	try {
		await MessageService.commentMessage(props.message, emoji.i)
	} catch (e) {
		console.error(e.toString())
	} finally {
		closeEmojiPicker()
	}
}

async function deleteMessage() {
	try {
		await MessageService.deleteMessage(props.message)
	} catch (err) {
		console.log(err.toString())
	} finally {
		emits("update")
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
	let toggleY = emojiPickerToggle.value.getBoundingClientRect().y

	let direction = 'high-picker'
	if (toggleY < y + halfContainerHeight) {
		direction = 'low-picker'
	}

	if (props.isAuthor) {
		return direction + ' author-picker'
	}
	return direction
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
		<div class="message-box" :class="{'author-message-box': isAuthor}">
			<div class="message-side" :class="{'author-message-side' : isAuthor}">
				<div class="message">
					<div class="body">
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
				</div>
				<span class="comment-btn" v-click-outside="closeEmojiPicker">
						<img class="svg-icon" ref="emoji-picker-toggle" src="@/assets/images/emoji.svg"
							 alt="add comment to the message"
							 @click="showEmojiPicker = !showEmojiPicker"/>
						<emoji-picker v-show="showEmojiPicker" :class="emojiPickerPosition()" class="emoji-picker"
									  :native="true" @select="setComment"/>
				</span>
			</div>
			<div class="options-box">
				<div class="option" :data-bs-target="`#comments-modal-${message.id}`" data-bs-toggle="modal"
					 @click="showCommentsModal = true">
					<img class="option-icon" src="@/assets/images/emoji.svg" alt=""/>
				</div>
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
	<show-comments-modal :message="message" :show="showCommentsModal" @close="showCommentsModal = false"/>
	<forwarding-modal :message="message" :show="showForwardModal" @close="showForwardModal = false"/>
</template>

<style scoped>
.svg-icon {
	width: 100%;
	height: 100%;
}

.comment-btn {
	position: relative;
	width: 1.5rem;
	height: 1.5rem;
	align-self: center;
	margin-right: 5px;
}

.emoji-picker {
	position: absolute;
	left: 0;
}

.low-picker {
	top: 0;
}

.high-picker {
	top: -320px;
}

.author-picker {
	left: -280px
}

.author-message-box {
	justify-content: end;
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
	scale: 0.8
}

.message-box {
	display: flex;
	flex-flow: row nowrap;
	height: 100%;
	width: 100%;
}

.message {
	display: flex;
	flex-flow: column nowrap;
	padding: 5px;
	max-width: 60%;
}

.message-side {
	display: flex;
	width: 100%;
	height: 100%;
	flex-flow: row nowrap;
	justify-content: start;
}

.author-message-side {
	justify-content: end;
	flex-direction: row-reverse;
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
}

.content {
	margin-bottom: 5px;
}

.message-text {
	font-size: 1.2rem;
}

.message-image-box {
	max-width: 336px;
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

</style>
