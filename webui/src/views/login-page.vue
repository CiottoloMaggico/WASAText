<script setup>
import {ref} from "vue"
import {SessionService} from "../services/session";
import router from "../router";

const username = ref("")

async function login() {
	try {
		const response = await SessionService.doLogin(username.value)
		router.push("/")
	} catch (error) {
		console.error(error.toString())
	}
}

</script>

<template>
	<main class="login-page">
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
	</main>
</template>

<style scoped>
.login-page {
	display: flex;
	flex-flow: column nowrap;
	justify-content: center;
	align-items: center;
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
	padding-left: 5%
}

.login-box-title {
	font-weight: bolder;
	color: #fff;
}

.input-label {
	font-weight: bold;
	color: #fff
}
</style>
