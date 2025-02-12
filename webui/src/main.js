import {createApp} from 'vue'
import App from './App.vue'
import router from './router'
import axios from './services/axios.js';
import {createPinia} from 'pinia'
import TheNewConversation from "@/components/TheNewConversation.vue";
import TheConversationList from "@/components/TheConversationList.vue";
import TheProfile from "@/components/TheProfile.vue";
import TheNewGroup from "@/components/TheNewGroup.vue";
import './assets/css/main.css'

const pinia = createPinia()
const app = createApp(App)
app.config.globalProperties.$axios = axios;
app.use(router)
app.use(pinia)

app.directive("clickOutside", {
	mounted: (el, binding) => {
		el.clickOutsideEvent = function (event) {
			if (!(el === event.target || el.contains(event.target))) {
				binding.value(event)
			}
		}
		document.addEventListener('click', el.clickOutsideEvent)
	},
	unmounted: (el) => {
		document.removeEventListener('click', el.clickOutsideEvent)
	}
})

app.component("TheNewConversation", TheNewConversation);
app.component("TheConversationList", TheConversationList);
app.component("TheProfile", TheProfile);
app.component("TheNewGroup", TheNewGroup);

app.mount('#app')
