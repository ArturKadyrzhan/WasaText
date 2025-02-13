import {createRouter, createWebHashHistory} from 'vue-router';
import ConversationView from "../views/ConversationView.vue";
import CreateGroupView from "../views/CreateGroupView.vue";
import UpdateGroupView from "../views/UpdateGroupView.vue";

const router = createRouter({
	history: createWebHashHistory(import.meta.env.BASE_URL),
	routes: [
		{
			path: '/conversation',
			component: ConversationView,
			name: 'Conversation',
		},
		{
			path: '/create-group',
			component: CreateGroupView,
			name: 'CreateGroup',
		},

		{
			path: '/update-group',
			component: UpdateGroupView,
			name: 'UpdateGroup',
		}
	]
})

export default router
