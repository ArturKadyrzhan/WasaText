<script setup>
const apiUrl = __API_URL__; // Directly using the constant from Vite config
</script>

<script>
import {getId, getToken} from "../store/auth";

export default {
  data() {
    return {
      photoPath:"",
      socket: null,
      userInfo: {},
      groupInfo: {},
      messages: [],
      newMessage: "",
      isModalVisible: false,
      searchQuery: '',
      userSearchResult: [],
      selectedUsers: [],
      showContextMenu: false,
      contextMenuPosition: { x: 0, y: 0 },
      selectedMessage: null,
      showEmojiPicker: false,
      emojis: ["😀", "😂", "😍", "😎", "😢", "😡", "👍", "👎", "🔥", "❤️", "💯"],
      selectedEmoji: "",
      selectedIndex: null,
      showForwardModal:false,
      selectedForwardedMessage: null,
      showReplyModal: false,
      selectedReplyMessage: null,
      forwardedById:0,
    };
  },
  beforeUnmount() {
    if (this.socket) {
      this.socket.close();
    }
  },
  methods: {
    formatTime(timestamp) {
      const options = {
        day: '2-digit',
        month: 'short',
        year: 'numeric',
        hour: '2-digit',
        minute: '2-digit',
      };
      return new Date(timestamp).toLocaleDateString(undefined, options);
    },
    goBack() {
      this.$router.push("/");
    },
    async getMessages(id) {
      const isGroup = this.groupInfo.ID === id;
      console.log(id)
      console.log(this.groupInfo.ID)
      console.log(isGroup)

      try {
        let response = await this.$axios.post('/messages',
            {
              id: id,
              isGroup: isGroup
            },
            {
              headers: {
                'Authorization': `Bearer ${getToken()}`,
                'Content-Type': 'application/json'
              }
            }
        );
        this.messages = response.data['messages'].map(message => ({
          id: message.message_id,
          message: message.message,
          isPhoto: message.is_photo,
          isSent: message.user_id == getId(),
          senderId:message.user_id,
          forwarded_by_username:message.forwarded_by_username ?? "",
          username: message.username ?? "",
          createdAt: message.createdAt ?? "",
          isReceived:true,
          isRead:message.is_read,
          emoji:message.emoji,
          replyTo:message.replied_message,
          replied: message.replied_message && message.replied_message.message_id != 0,
        }));
        console.log('messages forwarded',this.messages)
      } catch (error) {
        if (error.response) {
          if (error.response.status === 401) {
            console.error('Unauthorized: Token may be invalid or expired.');
            localStorage.removeItem('authToken');
          } else {
            console.error('Error fetching messages:', error.response.data);
          }
        } else if (error.request) {
          console.error('No response received:', error.request);
        } else {
          console.error('Error setting up request:', error.message);
        }
      }
    },
    async sendMessage() {
      if (this.newMessage.trim()) {
        const isGroup = !!this.groupInfo.ID;
        let isReceived = false;
        const repliedMessageId = this.selectedReplyMessage?.id || 0;
        try {
          const response = await this.$axios.post('/message/send', {
            text: this.newMessage,
            toUserId: this.userInfo.ID,
            isGroup:isGroup,
            groupId:this.groupInfo.ID,
            repliedMessageId: repliedMessageId,
          }, {
            headers: {
              'Authorization': `Bearer ${getToken()}`
            },
          });

          if (response.data) {
            isReceived = true
          }
          console.log('Message sent:', response.data);
        } catch (error) {
          console.error('Error sending message:', error);
          this.messages.pop();
        }
        const messageToSend = {
          username:this.authUsername,
          message: this.newMessage,
          isSent: true,
          isGroup: isGroup,
          createdAt: Date(),
          isReceived:isReceived,
        };
        this.messages.push(messageToSend);
        this.newMessage = "";
        this.$nextTick(() => {
          const chatMessages = this.$refs.chatMessages;
          chatMessages.scrollTop = chatMessages.scrollHeight;
        });
        this.showReplyModal = false;
      }
      window.location.reload();
    },
    addUserToGroup(user) {
      if (!this.selectedUsers.some(u => u.ID === user.ID)) {
        this.selectedUsers.push(user);
      }
    },
    removeUserFromGroup(user) {
      this.selectedUsers = this.selectedUsers.filter(u => u.ID !== user.ID);
    },
    addUsers() {
      this.isModalVisible = true;
    },
    closeModal() {
      this.isModalVisible = false;
    },
    closeForwardModal() {
      this.showForwardModal = false;
    },
    async findUsers() {
      if (!this.searchQuery) {
        this.userSearchResult = [];
        return;
      }
      try {
        let response = await this.$axios.get(`/users?search=${this.searchQuery}`, {
        headers: {
          'Authorization': `Bearer ${getToken()}`
        }
      });
      this.userSearchResult = response.data.users;
    } catch (error) {
      console.error("Error fetching users:", error);
      this.userSearchResult = [];
    }
  },
  async addNewUsers() {
    try {
      await this.$axios.post(
          "/group/users",
          {
            groupId: this.groupInfo.ID,
            users: this.selectedUsers
          },
          {
            headers: {
              Authorization: `Bearer ${getToken()}`,
            },
          }
      );
      this.isModalVisible = false;
      alert("User added")
    } catch (error) {
      console.error("Error add users to group:", error);
    }
  },
  triggerFileInput() {
    this.$refs.fileInput.click();
  },
  async uploadAndSendPhoto(event) {
    const file = event.target.files[0];
    const isGroup = !!this.groupInfo.ID;

    if (!file) {
      alert('No file selected!');
      return;
    }
    const formData = new FormData();
    formData.append('file', file);
    formData.append('isGroup', isGroup);
    formData.append('groupId', this.groupInfo.ID);
    formData.append('toUserId', this.userInfo.ID);
    console.log(formData)
    try {
      const response = await this.$axios.post("/message/send-photo", formData, {
        headers: {
           Authorization: `Bearer ${getToken()}`,
          "Content-Type": "multipart/form-data",
        },
      });

      this.photoPath = response.data['photoPath']
    } catch (error) {
      console.error("Error creating group:", error);
    }
    window.location.reload();
  },
  markAsRead() {
    const isGroup = !!this.groupInfo.ID;
    try {
       this.$axios.post(
          "/message/read",
          {
            groupId: this.groupInfo.ID,
            toUserId: this.userInfo.ID,
            isGroup:isGroup,
          },
          {
            headers: {
              Authorization: `Bearer ${getToken()}`,
            },
          }
      );
    } catch (error) {
      console.error("Error creating group:", error);
    }
  },
  openContextMenu(event, message, index) {
    this.selectedIndex = index;
    event.preventDefault();
    this.showContextMenu = true;
    this.contextMenuPosition = { x: event.clientX, y: event.clientY };
    this.selectedMessage = message;

    document.addEventListener("click", this.closeContextMenu);
  },
  closeContextMenu() {
    this.showContextMenu = false;
    document.removeEventListener("click", this.closeContextMenu);
  },
  async deleteMessage(message) {
    try {
      const response = await this.$axios.post(
          "/message",
          {
            id: message.id,
          },
          {
            headers: {
              Authorization: `Bearer ${getToken()}`,
            },
          }
      );
      if(response.data) {
        alert("Message deleted")
      }

    } catch (error) {
      console.error("Error creating group:", error);
    }
    window.location.reload()
  },
  commentMessage(message) {
    this.showEmojiPicker = true; // Show emoji picker
    this.showContextMenu = false; // Hide context menu
  }
  ,
    async uncommentMessage(message) {
      try {
        const response = await this.$axios.post(
            "/message/uncomment",
            {
              id: message.id,
            },
            {
              headers: {
                Authorization: `Bearer ${getToken()}`,
              },
            }
        );
        if(response.data) {
          alert("Reaction deleted")
        } else {
          alert("Reaction couldn't be deleted")
        }

      } catch (error) {
        console.error("Error creating group:", error);
      }
      window.location.reload()
    },
    async forwardMessage() {
      try {
        const response = await this.$axios.post(
            "/message/forward",
            {
              id: this.selectedForwardedMessage.id,
              text: this.selectedForwardedMessage.message,
              isPhoto: this.selectedForwardedMessage.isPhoto,
              users: this.selectedUsers,
              forwardedById:this.selectedForwardedMessage.senderId,
            },
            {
              headers: {
                Authorization: `Bearer ${getToken()}`,
              },
            }
        );
        if(response.data) {
          alert("Message forwarded")
        } else {
          alert("Message couldn't be forwarded")
        }

      } catch (error) {
        console.error("Error creating group:", error);
      }
    },
    forwardMessageModal(message) {
      this.showContextMenu = false;
      this.showForwardModal = true;
      this.selectedForwardedMessage = message;
    },
    replyToMessageModal(message) {
      this.showContextMenu = false;
      this.showReplyModal = true;
      this.selectedReplyMessage = message;
    },
    cancelReply() {
      this.showReplyModal = false;
    },

    selectEmoji(emoji, messageId) {
      console.log("Selected emoji:", emoji, messageId);
      this.selectedEmoji = emoji;
      this.showEmojiPicker = false;

      try {
         this.$axios.post(
            "/message/comment",
            {
              messageId: messageId,
              emoji: emoji,
            },
            {
              headers: {
                Authorization: `Bearer ${getToken()}`,
              },
            }
        );
      } catch (error) {
        console.error("Error creating group:", error);
      }
      window.location.reload()
    },
},
mounted() {
  if (this.$route.query.entity) {
    try {
      const entity = JSON.parse(this.$route.query.entity);
      const isGroup = this.$route.query.isGroup === "true";

      if (isGroup) {
        this.groupInfo = entity;
      } else {
        this.userInfo = entity;
      }

      this.markAsRead();
    } catch (error) {
      console.error('Error parsing conversation data:', error);
    }
  }
},
watch: {
  userInfo(newValue) {
    if (newValue && newValue.ID) {
      this.getMessages(newValue.ID);
    }
  },
  groupInfo(newValue) {
    if (newValue && newValue.ID) {
      this.getMessages(newValue.ID);
    }
  },
},
};
</script>
<template>
  <link href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.0.0-beta3/css/all.min.css" rel="stylesheet">

  <div class="conversation-page">
    <!-- Header -->
    <header class="chat-header">
      <div class="back-button" @click="goBack">
        ←
      </div>
      <div class="chat-title">
        {{ userInfo.Username ?? groupInfo.Name }}
      </div>
      <div v-if="groupInfo && groupInfo.Name">
        <div class="add-user-button" @click="addUsers">
          +
        </div>
      </div>
    </header>

    <!-- Chat messages container -->
    <div class="chat-messages" ref="chatMessages">
      <div
          v-for="(message, index) in messages"
          :key="index"
          :class="['chat-message', message.isSent ? 'sent' : 'received']"
          @contextmenu.prevent="openContextMenu($event, message, index)">

        <!-- shows sender username-->

        <!-- Replied Message Section -->
        <div v-if="message.replied" class="replied-message">
          <div class="replied-content">
            <span class="replied-text">
          {{ message.replyTo.content }}
        </span>
          </div>
        </div>

        <!-- shows sender username-->
<!--        <div class="message-header">-->
<!--          <span class="message-sender">{{ message.username }}</span>-->
<!--        </div>-->
        <div class="message-header">
            <span class="message-sender">
             {{ message.forwarded_by_username ? ` Forwarded from   ${message.forwarded_by_username} ` : message.username }}
            </span>
        </div>

        <!-- Message content -->
        <div class="message-content">
          <div v-if="message.isPhoto">
            <img :src="`${apiUrl}/${message.message}`" alt="Sent Photo" />
          </div>
          <div v-else>
            {{ message.message }}
          </div>
        </div>

        <!-- Time at the bottom left -->
        <div class="message-footer">
          <span class="message-time">{{ formatTime(message.createdAt) }}</span>
          <span>{{message.emoji}}</span>

          <div v-if="message.isSent && message.isRead">
            <i class="fas fa-check-double check-icon"></i>
          </div>
          <div v-else-if="message.isSent">
            <i class="fas fa-check check-icon"></i>
          </div>
        </div>
      </div>
    </div>

    <div
        v-if="showContextMenu"
        class="context-menu"
        :style="{ top: contextMenuPosition.y + 'px', left: contextMenuPosition.x + 'px' }"
    >
      <ul>
        <li @click="deleteMessage(selectedMessage)">Delete</li>
        <li @click="commentMessage(selectedMessage)">Comment</li>
        <li @click="uncommentMessage(selectedMessage)">Uncomment</li>
        <li @click="forwardMessageModal(selectedMessage)">Forward</li>
        <li @click="replyToMessageModal(selectedMessage)">Reply</li>


      </ul>
    </div>

    <div v-if="showEmojiPicker" class="emoji-picker">
      <div v-for="emoji in emojis" :key="emoji" class="emoji" @click="selectEmoji(emoji, selectedMessage.id)">
        {{ emoji }}
      </div>
    </div>

    <!-- Input box -->
    <div class="chat-input">
      <input
          v-model="newMessage"
          type="text"
          placeholder="Type a message..."
          @keyup.enter="sendMessage"
      />
      <button @click="sendMessage">Send</button>
      <input
          type="file"
          ref="fileInput"
          @change="uploadAndSendPhoto"
          accept="image/*"
          style="display: none;"
      />
      <button @click="triggerFileInput">
        <i class="fas fa-paperclip"></i>
      </button>
    </div>

    <div v-if="showReplyModal" class="modal">
      <div class="modal-content">
        <div class="modal-header">Replying to: "{{ selectedReplyMessage.message }}"</div>
        <input v-model="newMessage" placeholder="Type a message..." />
        <div class="modal-buttons">
          <button @click="sendMessage">Send</button>
          <button @click="cancelReply">Cancel</button>
        </div>
      </div>
    </div>

    <div v-if="isModalVisible" class="modal">
      <form @submit.prevent="addNewUsers">

        <div class="modal-content">
          <h3>Find Users</h3>
          <div class="user-finder">
            <input
                type="text"
                v-model="searchQuery"
                placeholder="Search for a user to start a conversation"
                class="search-bar"
                @input="findUsers"
            />
          </div>

          <div v-if="userSearchResult.length" class="search-results">
            <ul>
              <li v-for="user in userSearchResult" :key="user.ID" class="search-item">
                <img :src="`${apiUrl}/${user.ProfilePhotoURL}`" alt="Profile" class="profile-photo" />
                <div class="user-details">
                  <h5>{{ user.Username }}</h5>
                  <button type="button" @click="addUserToGroup(user)">Add to Group</button>
                </div>
              </li>
            </ul>
          </div>
          <div v-if="selectedUsers.length" class="selected-users">
            <h4>Selected Users</h4>
            <ul>
              <li v-for="user in selectedUsers" :key="user.ID" class="selected-user-item">
                <img :src="`${apiUrl}/${user.ProfilePhotoURL}`" alt="Profile" class="profile-photo" />
                <div class="user-details">
                  <h5>{{ user.Username }}</h5>
                  <button type="button" @click="removeUserFromGroup(user)">Remove</button>
                </div>
              </li>
            </ul>
          </div>

          <button @click="closeModal">Close</button>

          <button type="submit">Update Group</button>
        </div>
      </form>
    </div>


    <div v-if="showForwardModal" class="modal">
      <form @submit.prevent="forwardMessage">

        <div class="modal-content">
          <h3>Find Users</h3>
          <div class="user-finder">
            <input
                type="text"
                v-model="searchQuery"
                placeholder="Search for a user to start a conversation"
                class="search-bar"
                @input="findUsers"
            />
          </div>

          <div v-if="userSearchResult.length" class="search-results">
            <ul>
              <li v-for="user in userSearchResult" :key="user.ID" class="search-item">
                <img :src="`${apiUrl}/${user.ProfilePhotoURL}`" alt="Profile" class="profile-photo" />
                <div class="user-details">
                  <h5>{{ user.Username }}</h5>
                  <button type="button" @click="addUserToGroup(user)">Add to Forward Message</button>
                </div>
              </li>
            </ul>
          </div>

          <div v-if="selectedUsers.length" class="selected-users">
            <h4>Selected Users</h4>
            <ul>
              <li v-for="user in selectedUsers" :key="user.ID" class="selected-user-item">
                <img :src="`${apiUrl}/${user.ProfilePhotoURL}`" alt="Profile" class="profile-photo" />
                <div class="user-details">
                  <h5>{{ user.Username }}</h5>
                  <button type="button" @click="removeUserFromGroup(user)">Remove</button>
                </div>
              </li>
            </ul>
          </div>

          <button @click="closeForwardModal">Close</button>

          <button type="submit">Forward Message</button>
        </div>
      </form>
    </div>
  </div>
</template>












<style scoped>
/* General styles */
.conversation-page {
  display: flex;
  flex-direction: column;
  height: 100vh;
  background-color: #f5f8fb;
}

/* Header styles */
.chat-header {
  display: flex;
  flex-direction: row;
  align-items: center;
  z-index: 1;
  background-color: #000000;
  color: #fff;
  font-size: 25px;
}

.back-button {
  cursor: pointer;
  margin-right: 15px;
  font-size: 25px;
  font-weight: bold;
}

.chat-title {
  flex-grow: 1;
  text-align: center;
  padding: 5px;
}

/* Messages styles */
.chat-messages {
  flex: 1;
  overflow-y: auto;
  padding: 15px;
  display: flex;
  flex-direction: column;
}

/* Input box styles */
.chat-input {
  display: flex;
  padding: 10px;
  border-top: 1px solid #ddd;
  background-color: #fff;
  position: sticky;
  bottom: 0;
}

.chat-input input {
  flex: 1;
  padding: 10px;
  border: 1px solid #ccc;
  border-radius: 20px;
  outline: none;
}

.chat-input button {
  background-color: #0088cc;
  color: white;
  border: none;
  border-radius: 20px;
  padding: 10px 15px;
  margin-left: 10px;
  cursor: pointer;
}

.chat-input button:hover {
  background-color: #005f99;
}

.chat-message {
  padding: 10px;
  margin: 5px 0;
  border-radius: 8px;
  max-width: 60%;
  word-wrap: break-word;
  position: relative;
}

.chat-message.sent {
  background-color: #d1e7dd;
  margin-left: auto;
  text-align: right;
}

.chat-message.received {
  background-color: #f8d7da;
  margin-right: auto;
}

.message-header {
  font-weight: bold;
  margin-bottom: 5px;
}

.message-footer {
  display: flex;
  align-items: center;
  justify-content: space-between; /* Размещаем элементы по краям */
  padding: 5px 10px;
  font-size: 12px;
  color: gray;
}

.check-icon {
  color: gray;          /* Цвет иконки */
  margin-left: 5px;     /* Отступ слева */
  font-size: 14px;      /* Размер иконки */
}

.modal {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
}

.modal-content {
  background: white;
  padding: 20px;
  border-radius: 8px;
  max-width: 500px;
  width: 100%;
}

.search-results, .selected-users {
  margin-top: 20px;
}

.context-menu {
  position: absolute;
  z-index: 1000;
  background: #fff;
  border: 1px solid #ddd;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
  border-radius: 4px;
}

.context-menu ul {
  list-style: none;
  padding: 0;
  margin: 0;
}

.context-menu ul li {
  padding: 8px 12px;
  cursor: pointer;
  transition: background 0.2s;
}

.context-menu ul li:hover {
  background: #f5f5f5;
}
/* Modal Background */
.modal {
  position: fixed; /* Fixed position to cover the whole screen */
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(0, 0, 0, 0.5); /* Semi-transparent black background */
  display: flex;
  align-items: center;
  justify-content: center;
}

/* Modal Content */
.modal-content {
  background: #fff; /* White background for the modal */
  padding: 20px;
  border-radius: 8px; /* Rounded corners */
  width: 400px;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1); /* Light shadow */
}

/* Header Text */
.modal-header {
  font-size: 18px;
  font-weight: bold;
  margin-bottom: 15px;
}

/* Input Field */
.modal input {
  width: 100%;
  padding: 10px;
  margin-bottom: 15px;
  border: 1px solid #ccc;
  border-radius: 4px;
}

/* Buttons Container */
.modal-buttons {
  display: flex;
  justify-content: flex-end;
  gap: 10px; /* Space between buttons */
}

/* Buttons Style */
.modal button {
  padding: 8px 16px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}

/* Primary Button */
.modal button:first-child {
  background: #4caf50;
  color: white;
}

/* Cancel Button */
.modal button:last-child {
  background: #f44336;
  color: white;
}
.replied-message {
  background: #f0f0f0;
  border-left: 3px solid #4caf50;
  padding: 5px 10px;
  margin-bottom: 8px;
  border-radius: 5px;
}

.replied-content {
  font-size: 12px;
  color: #555;
}

.replied-text {
  font-style: italic;
  color: #666;
}
</style>


