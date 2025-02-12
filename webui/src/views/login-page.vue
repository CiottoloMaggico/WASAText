<script setup>
import {ref} from "vue"
import router from "@/router";
import SessionService from "@/services/sessionService";

const username = ref("")

async function login() {
	try {
		await SessionService.doLogin(username.value)
	} catch (error) {
		console.error(error)
		return
	}

	await router.push({name: "homepage"})
}

</script>

<template>
	<div class="login-page">
		<div class="login-box">
			<div class="login-box-content">
				<h1 class="login-box-title mb-4">Login</h1>
				<form @submit.prevent="login">
					<div class="form-group">
						<label for="usernameInput" class="form-label input-label">Username:</label>
						<input id="usernameInput" name="usernameInput" type="text" class="form-control mb-4"
							   placeholder="Username" v-model="username">
						<button type="submit" class="btn primary-btn">Submit</button>
					</div>
				</form>
			</div>
		</div>
	</div>
</template>

<style scoped>
.login-page {
	display: flex;
	flex-flow: column nowrap;
	justify-content: center;
	align-items: center;
	width: 100%;
	height: 100%;
}

.login-box {
	width: 50%;
	height: 500px;
	border-radius: var(--MAIN-BORDER-RADIUS);
	background-color: var(--FOURTH-COLOR);
}

.login-box-content {
	display: flex;
	flex-flow: column nowrap;
	align-items: start;
	justify-content: center;
	height: 100%;
	width: 100%;
	padding-left: 5%;
}

.login-box-title {
	font-weight: bolder;
	color: #fff;
}

.input-label {
	font-weight: bold;
	color: #fff;
}
</style>
