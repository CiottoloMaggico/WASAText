<script setup>
import {computed, reactive, useTemplateRef, watch} from "vue"
import {getApiUrl} from "../services/axios";
import {getAuthentication} from "../services/sessionService";
import MessageService from "@/services/messageService";

const props = defineProps(["conversation", "wantReply"])
const emits = defineEmits(["update", "clearReply"])

const fileUploadElement = useTemplateRef("file-upload")

const newMessage = reactive({
	content: null,
	attachment: null,
	repliedMessage: props.wantReply,
})

const attachmentPreviewUrl = computed(() => {
	if (!newMessage.attachment) {
		return ""
	}
	return URL.createObjectURL(newMessage.attachment)
})
const previewEnabled = computed(() => {
	return newMessage.repliedMessage || newMessage.attachment
})

watch(() => props.wantReply, (newReply) => {
	if (newReply) {
		newMessage.repliedMessage = newReply
	}
})

watch(() => props.conversation.id, () => {
	clearAll()
})

async function sendMessage() {
	if (!newMessage.content && !newMessage.attachment) {
		return
	}

	try {
		await MessageService.sendMessage(props.conversation, newMessage)
	} catch (err) {
		console.log(err.toString())
	} finally {
		emits('update')
		clearAll()
	}
}

function clearAttachment() {
	fileUploadElement.value.value = null
	newMessage.attachment = null
}

function clearReplyTo() {
	newMessage.repliedMessage = null
	emits("clearReply", null)
}

function clearAll() {
	Object.assign(newMessage, {
		content: null,
		attachment: null,
		repliedMessage: null,
	})
	fileUploadElement.value.value = null
	emits("clearReply", null)
}

function fileUploaded() {
	newMessage.attachment = fileUploadElement.value.files.item(0)
}
</script>

<template>
	<div class="new-message-bar" :class="{'preview-enabled' : previewEnabled}">
		<div v-if="previewEnabled" class="message-attachment-row">
			<div v-if="newMessage.attachment" class="attachment-preview-box">
				<span class="cancel-attachment" @click="clearAttachment">
					<img class="cancel-icon" src="@/assets/images/cancel.png" width="30" height="30"/>
				</span>
				<img class="attachment-preview" :src="attachmentPreviewUrl"/>
			</div>
			<div v-if="newMessage.repliedMessage" class="replied-message-box">
				<span class="cancel-attachment" @click="clearReplyTo">
					<img class="cancel-icon" src="@/assets/images/cancel.png" width="30" height="30"/>
				</span>
				<div class="replied-message">
					<div class="replied-message-content">
						<span v-if="newMessage.repliedMessage.author.uuid !== getAuthentication()" class="replied-message-author" >
								{{ newMessage.repliedMessage.author.username }}
						</span>
						<span v-else class="replied-message-author">Me</span>

						<span v-if="newMessage.repliedMessage.content" class="replied-message-text">
								{{ newMessage.repliedMessage.content }}
						</span>
						<span v-else-if="newMessage.repliedMessage.attachment" class="replied-message-text">
								Image
						</span>
					</div>
					<span v-if="newMessage.repliedMessage.attachment" class="replied-message-image-box">
						<img :src="getApiUrl(newMessage.repliedMessage.attachment.fullUrl)"
							 class="replied-message-image"/>
					</span>
				</div>
			</div>
		</div>
		<form @submit.prevent="sendMessage" class="send-message-form">
			<div class="input-file-box">
				<label for="attachment-file">
					<img class="input-file-icon" src="@/assets/images/plus.png" width="512" height="512"/>
				</label>
				<input id="attachment-file" class="d-none" type="file" accept="image/*" @change="fileUploaded"
					   ref="file-upload">
			</div>
			<div class="text-input-box">
				<input id="messagetext" name="messagetext" type="text" class="form-control"
					   placeholder="write a message" v-model="newMessage.content" >
			</div>

			<div class="submit-button-box">
				<label for="submit-button">
					<img class="submit-button-icon" src="@/assets/images/send.png" width="512" height="512"/>
				</label>
				<input id="submit-button" type="submit" class="d-none"/>
			</div>
		</form>
	</div>
</template>

<style scoped>
.preview-enabled {
	height: max-content !important;
}

.new-message-bar {
	display: flex;
	flex-flow: column nowrap;
	height: 4rem;
	justify-content: center;
	min-height: 4rem;
	padding: 10px;
	width: 100%;
}

.message-attachment-row {
	height: 100%;
	width: 100%;
	display: flex;
	flex-flow: row nowrap;
	padding: 10px;
	justify-content: center;
	align-items: center;
	gap: 5rem;
}

.replied-message-box {
	background-color: var(--SECONDARY-COLOR);
	border-radius: var(--MAIN-BORDER-RADIUS);
	flex-shrink: 0;
	width: 20rem;
	height: 5rem;
	max-height: 5rem;
	padding: 10px;
	position: relative;
}


.replied-message {
	width: 100%;
	height: 100%;
	display: flex;
	flex-flow: row nowrap;
	align-items: center;
}


.replied-message-content {
	display: flex;
	flex-flow: column nowrap;
	width: 100%;
	height: 100%;
	justify-content: center;
	align-items: start;
	overflow: hidden;
}

.replied-message-author {
	max-width: 100%;
	overflow: hidden;
	text-overflow: ellipsis;
	white-space: nowrap;
	font-size: 1.2rem;
	margin-bottom: 5px;
	flex-shrink: 0;
}

.replied-message-text {
	font-size: 1.1rem;
	overflow: hidden;
	text-overflow: ellipsis;
	white-space: nowrap;
	width: 100%;
}

.replied-message-image-box {
	width: 3rem;
	height: 3rem;
	flex-shrink: 0;
}

.replied-message-image {
	width: 100%;
	height: 100%;
	object-fit: cover;
	border-radius: 5px;
}

.attachment-preview-box {
	position: relative;
	border-radius: var(--MAIN-BORDER-RADIUS);
	flex-shrink: 0;
	height: 15rem;
	max-height: 15rem;
	max-width: 15rem;
}

.attachment-preview {
	height: 100%;
	width: 100%;
	border-radius: var(--MAIN-BORDER-RADIUS);
	border: 1px solid black;
	object-fit: cover;
}

.cancel-attachment {
	position: absolute;
	width: 2rem;
	height: 2rem;
	top: 0;
	right: 0;
	cursor: pointer;
}

.cancel-icon {
	width: 100%;
	height: 100%;
}

.send-message-form {
	flex-shrink: 0;
	display: flex;
	flex-flow: row nowrap;
	width: 100%;
	align-items: center;
}

.input-file-box {
	flex-shrink: 0;
	width: 2.5rem;
	height: 2.5rem;
	margin-right: 10px;
}

.input-file-icon {
	width: 100%;
	height: 100%;
	cursor: pointer;
}

.text-input-box {
	width: 100%;
	height: 100%;
}

.form-control {
	line-height: 2rem;
}

.submit-button-box {
	flex-shrink: 0;
	width: 2.5rem;
	height: 2.5rem;
	margin-left: 10px;
}

.submit-button-icon {
	width: 100%;
	height: 100%;
	cursor: pointer;
}
</style>

