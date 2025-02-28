<script setup>
const apiUrl = __API_URL__; // Directly using the constant from Vite config
</script>


<script>
import {getToken} from "../store/auth";


export default {
  data() {
    return {
      searchQuery: "",
      userSearchResult: [],
      groupName: "",
      selectedUsers: [],
      profilePhoto: null,
      isEditing: false,
      updatedGroupName: "",
      updatedProfilePhoto: null,
      selectedGroup: null,

    };
  },
  methods: {
    //find users to while adding new users to group
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
  addUserToGroup(user) {
    if (!this.selectedUsers.find((u) => u.ID === user.ID)) {
      this.selectedUsers.push(user);
    }
  },
  removeUserFromGroup(user) {
    this.selectedUsers = this.selectedUsers.filter((u) => u.ID !== user.ID);
  },
  handleProfilePhotoUpload(event) {
    this.profilePhoto = event.target.files[0];
  },

  //create-group.go
  async createGroup() {
    if (!this.groupName || this.selectedUsers.length === 0 || !this.profilePhoto) {
      alert("Please provide a group name and invite at least one user and upload group photo.");
      return;
    }
    const formData = new FormData();
    formData.append("groupName", this.groupName);
    formData.append("selectedUsers", JSON.stringify(this.selectedUsers));

    if (this.profilePhoto) {
      formData.append("profilePhoto", this.profilePhoto);
    }
    try {
      await this.$axios.post("/group", formData, {
        headers: {
          Authorization:` Bearer ${getToken()}`,
          "Content-Type": "multipart/form-data",
        },
      });
      window.location.reload();
    } catch (error) {
      console.error("Error creating group:", error);
    }
  },
},
};
</script>

<template>
  <div id="app">
    <div class="create-group">
      <h3>Create Group</h3>
      <form @submit.prevent="createGroup">
        <div>
          <label for="profilePhoto">Group Profile Photo:</label>
          <input type="file" id="profilePhoto" @change="handleProfilePhotoUpload" />
        </div>
        <div>
          <label for="groupName">Group Name:</label>
          <input type="text" id="groupName" v-model="groupName" required />
        </div>
        <h3>Invite Users</h3>
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
                <button type="button" @click="addUserToGroup(user)">Add to Group</button> <!-- Changed type to button -->
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
                <button type="button" @click="removeUserFromGroup(user)">Remove</button> <!-- Changed type to button -->
              </div>
            </li>
          </ul>
        </div>

        <button type="submit">Create Group</button>
      </form>
    </div>


    <router-view>
    </router-view>
  </div>
</template>



<style scoped>
.create-group {
  padding: 20px;
  border: 1px solid #ccc;
  margin: 20px;
  border-radius: 8px;
}

.create-group input,
.create-group select,
.create-group button {
  margin: 10px 0;
  padding: 10px;
  width: 100%;
  box-sizing: border-box;
}

.create-group button {
  background-color: #0088cc;
  color: white;
  border: none;
  cursor: pointer;
}

.create-group button:hover {
  background-color: #005f99;
}

.search-results ul {
  list-style-type: none;
  padding: 0;
}

.search-item {
  display: flex;
  align-items: center;
  padding: 10px;
  border-bottom: 1px solid #ddd;
}

.search-item button {
  background-color: #28a745;
  color: white;
  padding: 5px 10px;
  border: none;
  cursor: pointer;
}

.search-item button:hover {
  background-color: #218838;
}

.selected-users ul {
  list-style-type: none;
  padding: 0;
}

.selected-user-item {
  display: flex;
  align-items: center;
  padding: 10px;
  border-bottom: 1px solid #ddd;
}

.selected-user-item button {
  background-color: #dc3545;
  color: white;
  padding: 5px 10px;
  border: none;
  cursor: pointer;
}

.selected-user-item button:hover {
  background-color: #c82333;
}
</style>