<script setup>
const apiUrl = __API_URL__; // Directly using the constant from Vite config
</script>


<script>
import {getToken} from "../store/auth";


export default {
  data() {
    return {
      searchQuery: "",  // for search input for finding users.
      userSearchResult: [], // for results
      groupName: "", // Stores the name of the created group
      selectedUsers: [], // Stores users added to the group.
      profilePhoto: null, // Stores the uploaded group profile photo.
      isEditing: false, // Boolean to track if the group is being edited.
      updatedGroupName: "", // Stores the new group name (for editing)
      updatedProfilePhoto: null, // Stores the new profile photo (for editing).
      selectedGroup: null,

    };
  },
  methods: {
    //find users to while adding new users to group
    async findUsers() {
      if (!this.searchQuery) {
        this.userSearchResult = []; // if searchQuery is empty → If yes, clears userSearchResult and exits.
        return;
      }
      try { // we send GET request to /users?search=query → Searches for users.
        let response = await this.$axios.get(`/users?search=${this.searchQuery}`, {
        headers: { //  authentication token
          'Authorization': `Bearer ${getToken()}`
        }
      });
      //If successful, updates userSearchResult.
      this.userSearchResult = response.data.users;
    } catch (error) {
      // If an error occurs, logs the error and empty search results.
      console.error("Error fetching users:", error);
      this.userSearchResult = [];
    }
  },
  // Checks if the user is already in selectedUsers.If not, adds them to the list
  addUserToGroup(user) {
    if (!this.selectedUsers.find((u) => u.ID === user.ID)) {
      //Checks if the user is already in selectedUsers, if not we add him
      this.selectedUsers.push(user);
    }
  },
    // removes-user
  removeUserFromGroup(user) {
    this.selectedUsers = this.selectedUsers.filter((u) => u.ID !== user.ID);
  },
  // Capturing the selected file when the user uploads an image.
  handleProfilePhotoUpload(event) {
    this.profilePhoto = event.target.files[0];
  },

  //create-group.go
  async createGroup() {
    // Validates required fields (
    if (!this.groupName || this.selectedUsers.length === 0 || !this.profilePhoto) {
      alert("Please provide a group name and invite at least one user and upload group photo.");
      return;
    }
    // we create Form data object, fill required fields
    const formData = new FormData(); //
    formData.append("groupName", this.groupName);
    formData.append("selectedUsers", JSON.stringify(this.selectedUsers));

    if (this.profilePhoto) {
      formData.append("profilePhoto", this.profilePhoto);
    }
    try {
      await this.$axios.post("/group", formData, {
        headers: {
          Authorization:` Bearer ${getToken()}`,
          "Content-Type": "multipart/form-data", // WHAT I send
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