import {createApp, reactive} from 'vue'
import App from './App.vue'
import router from './router'
import axios from './services/axios.js';
import ErrorMsg from './components/ErrorMsg.vue'
import LoadingSpinner from './components/LoadingSpinner.vue'
import {createPinia} from 'pinia'
import './assets/css/main.css'
import TheNewConversation from "@/components/TheNewConversation.vue";
import TheConversationList from "@/components/TheConversationList.vue";
import TheProfile from "@/components/TheProfile.vue";
import TheNewGroup from "@/components/TheNewGroup.vue";

const pinia = createPinia()
const app = createApp(App)
app.config.globalProperties.$axios = axios;
app.component("ErrorMsg", ErrorMsg);
app.component("LoadingSpinner", LoadingSpinner);
app.component("TheNewConversation", TheNewConversation);
app.component("TheConversationList", TheConversationList);
app.component("TheProfile", TheProfile);
app.component("TheNewGroup", TheNewGroup);
app.use(router)
app.use(pinia)
app.mount('#app')
