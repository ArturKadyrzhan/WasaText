 <template>
    <div class="login-container" v-if="!checkLogin">
      <h1 class="h2">Login</h1>

      <ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>

      <form @submit.prevent="login" class="login-form">
        <div class="mb-3">
          <label for="username" class="form-label">Username</label>
          <input
              type="text"
              id="username"
              v-model="username"
              class="form-control"
              placeholder="Enter your username"
              required
          />
        </div>
        <div class="mb-3">
          <label for="password" class="form-label">Password</label>
          <input
              type="password"
              id="password"
              v-model="password"
              class="form-control"
              placeholder="Enter your password"
          />
        </div>

        <div class="d-flex justify-content-between align-items-center">
          <button type="submit" class="btn btn-primary" :disabled="loading">
            {{ loading ? "Logging in..." : "Login" }}
          </button>
        </div>
      </form>
    </div>
    <div v-else>
      <header class="navbar navbar-dark sticky-top bg-dark flex-md-nowrap p-0 shadow">
        <a class="navbar-brand col-md-3 col-lg-2 me-0 px-3 fs-6" href="#/">WasaText</a>
        <button class="navbar-toggler position-absolute d-md-none collapsed" type="button" data-bs-toggle="collapse" data-bs-target="#sidebarMenu" aria-controls="sidebarMenu" aria-expanded="false" aria-label="Toggle navigation">
          <span class="navbar-toggler-icon"></span>
        </button>
      </header>

      <div class="container-fluid">
        <div class="row">
          <nav id="sidebarMenu" class="col-md-3 col-lg-2 d-md-block bg-light sidebar collapse">
            <div class="position-sticky pt-3 sidebar-sticky">
              <div class="profile-picture-upload">
                <h4>Your Profile Image </h4>
                <img v-if="profileImage" :src="`${apiUrl}/${profileImage}`" alt="Profile" class="profile-photo" />
                <p v-else> profile image doesn't exist,upload it...</p>
                <label for="username">Your current username:</label>
                <input id="username" v-model="username" type="text" placeholder=" username" />
                <button @click="updateUsername">Update Username</button>
                <p v-if="message">{{ message }}</p>
                <input type="file" @change="onFileChange" accept="image/*" />
                <button class="upload-button" @click="uploadImage" :disabled="!selectedFile">Upload</button>
                <div v-if="previewUrl">
                  <h4>Preview:</h4>
                  <img :src="previewUrl" alt="Profile Preview" class="profile-preview" />
                </div>
              </div>
              <div class="conversations">
                <!-- Search Bar -->
                <div class="user-finder">
                  <input
                      type="text"
                      v-model="searchQuery"
                      placeholder="Search for a user to start a conversation"
                      class="search-bar"
                      @input="findUser"
                  />
                </div>

                <!-- Search Results -->
                <div v-if="userSearchResult.length" class="search-results">
                  <ul>
                    <li v-for="user in userSearchResult" :key="user.ID" class="search-item">
                      <img :src="`${apiUrl}/${user.ProfilePhotoURL}`" alt="Profile" class="profile-photo" />
                      <div class="user-details">
                        <h5>{{ user.Username }}</h5>
                        <button @click="startConversation(user)">Open Chat</button>
                      </div>
                    </li>
                  </ul>
                </div>
                <p v-else-if="searchQuery">No users found!</p>
                
                <h5>Your Conversations:</h5>

                <!-- User Conversations -->
                <ul v-if="users && users.length">
                  <li v-for="user in users" :key="user.ID" class="conversation-item">
                    <img :src="`${apiUrl}/${user.ProfilePhotoURL}`" alt="Profile" class="profile-photo" />
                    <div class="user-details">
                      <h5>{{ user.Username }}</h5>
                      <button @click="startConversation(user)">Open Chat</button>
                    </div>
                  </li>
                </ul>
                <h5>Your Groups:</h5>

                <!-- Group Conversations -->
                <ul v-if="groups && groups.length">
<!--                <ul v-if="groups.length">-->
                  <li v-for="group in groups" :key="group.ID" class="conversation-item">
                    <img :src="`${apiUrl}/${group.GroupPhotoURL}`" alt="Group Profile" class="profile-photo" />
                    <div class="user-details">
                      <h5>{{ group.Name }}</h5>
                      <button @click="startConversation(group)">Open Chat</button>
                      <button @click="updateGroup(group.ID)">Update Group</button>
                      <button @click="leaveGroup(group)">Leave Group</button>
                    </div>
                  </li>
                </ul>

                <!-- Fallback message -->
                <p v-if="!users && !groups">No conversations yet!</p>
              </div>
              <h6 class="sidebar-heading d-flex justify-content-between align-items-center px-3 mt-4 mb-1 text-muted text-uppercase">
                <button @click="createGroup" class="nav-link logout-btn">
                  <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#log-out"/></svg>
                  Create Group
                </button>
              </h6>

              <h6 class="sidebar-heading d-flex justify-content-between align-items-center px-3 mt-4 mb-1 text-muted text-uppercase">
                  <button @click="logout" class="nav-link logout-btn">
                    <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#log-out"/></svg>
                    Logout
                  </button>
              </h6>
            </div>
          </nav>

          <main class="col-md-9 ms-sm-auto col-lg-10 px-md-4">
            <RouterView />
          </main>
        </div>
      </div>
    </div>
  </template>


 <script setup> // import dependencies
 import { RouterView } from 'vue-router' // Allows dynamic rendering of components based on the URL.
 const apiUrl = __API_URL__;  //  Using a global variable for the API base URL.
 </script>
 <script>
 import {checkLoginStatus, getToken, logIn, logOut} from "./store/auth";

 export default {
   // Defining Component State
   data: function() {
     return {
       userSearchResult: [],
       users: [], // Stores chat participants.
       groups:[],  // Stores chat participants.
       errormsg: null,
       loading: false,
       username: "", //  Stores login credentials.
       password: "", //  Stores login credentials.
       checkLogin: false,
       searchQuery: "",
       selectedFile: null,  // Handles image uploads.
       previewUrl: null,    // Handles image uploads.
       message: "",
       profileImage: null,
     }
   },
   methods: {
     // get-users.go
     // 1) Searches for users based on searchQuery
     // Calls /users?search=query with an auth token, If search fails, it clears results.
     async findUser() {
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
         this.userSearchResult = response.data['users'];
         console.log(this.userSearchResult, "USER")
       } catch (error) {
         console.error("Error fetching users:", error);
         this.userSearchResult = [];
       }
     },
     // 2)Logging In do-login.go
     // here we send login credentials (username, password) to /session.
     // if login fails, shows an error message.
     async login() {
       this.loading = true;
       this.errormsg = null;
       try {
         let response = await this.$axios.post("/session", {
           username: this.username,
           password: this.password,
         });
         logIn(response.data['token'], response.data['id'], response.data['username'])

       } catch (e) {
         console.log(e,"eror meror")
         this.errormsg = "Login failed. Please check your credentials.";
       }
       this.loading = false;
       window.location.reload();
     },
     //3) get-my-conversations.go, Retrieves your chats with user/groups.
     async fetchConversations() {
       try {
          console.log('Token:',getToken())
         let response = await this.$axios.get('/conversations',{
           headers: {
             'Authorization':`Bearer ${getToken()}`,
           },
         });
         this.users = response.data['result']['users'];
         this.groups = response.data['result']['groups'];
         console.log(response.data)
       } catch (error) {
         console.log(error, "ERROR")
         if (error.response) {
           if (error.response.status === 401) {
             console.error('Unauthorized: Token may be invalid or expired.');
             localStorage.removeItem('authToken');
           } else {
             console.error('Error fetching conversations:', error.response.data);
           }
         } else if (error.request) {
           console.error('No response received:', error.request);
         } else {
           console.error('Error setting up request:', error.message);
         }
       }
     },
     //4) get-user-profile.go, Calls /profile to get username & profile image.
     async fetchProfile() {
       try {
         const response = await this.$axios.get("/profile", {
           headers: { Authorization: `Bearer ${getToken()}` },
         });
         this.username = response.data.profile.Username;
         this.profileImage = response.data.profile.ProfilePhotoURL;
         console.log(response.data)
         console.log(this.username)
         console.log(this.profileImage)
       } catch (error) {
         console.error("Error fetching profile:", error);
         this.message = "Failed to load profile.";
       }
     },
     // Start a chat with someone
     startConversation(userOrGroup) {
       if (!userOrGroup) {
         console.error('Missing user ID. Cannot navigate to conversation.');
         return;
       }
       const isGroup = userOrGroup.Name !== undefined;
       const entity = encodeURIComponent(JSON.stringify(userOrGroup));
       const isGroupParam = isGroup ? 'true' : 'false';
       window.location.href = `#/conversation?entity=${entity}&isGroup=${isGroupParam}`;
       window.location.reload();
     },
     logout() {
       logOut()
       window.location.reload();
     },
     createGroup() {
       this.$router.push({
         name: 'CreateGroup',
       });
     },
     updateGroup(groupId) {
       this.$router.push({
         name: "UpdateGroup",
         query: { groupId }
       });
     },

    // Handles profile picture uploads.Converts file into a preview URL before uploading.
     onFileChange(event) {
       const file = event.target.files[0];
       console.log(file)
       if (file) {
         this.selectedFile = file;
         this.previewUrl = URL.createObjectURL(file);
       }
       else{
         console.log("file doesn't selected")
       }
     },
     // set-my-photo.go
     async uploadImage() {
       const formData = new FormData();
       formData.append("profile_picture", this.selectedFile);
       try {
         const response = await this.$axios.post("/profile/photo", formData, {
           headers: {
             "Content-Type": "multipart/form-data",
             Authorization: `Bearer ${getToken()}`,
           },
         });
         console.log("Image uploaded:", response.data);
         this.previewUrl=response.data.success;
         alert("Profile picture uploaded successfully!");
       } catch (error) {
         console.error("Error uploading image:", error);
         alert("Failed to upload image.");
       }
       window.location.reload();
     },

     //set-my-username.go
     async updateUsername() {
       if (!this.username.trim()) {
         this.message = "Username cannot be empty!";
         return;
       }
       try {
         const response = await this.$axios.post("/profile/username",{
           username:this.username
         },{
           headers: {
             Authorization: `Bearer ${getToken()}`,
           },
       });
         this.username = response.data.username;
         window.location.reload();
         this.message = "Username updated successfully!";
         alert("Username updated successfully!")
       } catch (error) {
         this.message = "This username already exists!!!";
       }
     },
     async leaveGroup(group) {
       try {
         await this.$axios.post("/group/leave",{
            group:group
         }, {
           headers: {
             Authorization: `Bearer ${getToken()}`,
           },
         });
         window.location.href = `#/`;
         window.location.reload();
       } catch (error) {
         console.error("Error uploading image:", error);
         alert("Failed to upload image.");
       }
     }
   },
   mounted() { // Runs when the component is loaded
     this.checkLogin = checkLoginStatus() // Retrieves token from localStorage, If a token exists, it:retrive users,chats
     if (this.checkLogin){
       this.fetchProfile();
       this.fetchConversations();
     }
   }
 }


 </script>

  <style>
  .login-container {
    max-width: 400px;
    margin: 50px auto;
    padding: 20px;
    border: 1px solid #ccc;
    border-radius: 5px;
    background: #fff;
    box-shadow: 0px 4px 6px rgba(0, 0, 0, 0.1);

  }

  .login-form .form-label {
    font-weight: 600;
  }

  .logout-btn {
    background: none;
    border: none;
    color: inherit;
    font: inherit;
    text-align: left;
    display: flex;
    align-items: center;
    font-size: 18px; /* Increase text size */
  }

  .logout-btn svg {
    width: 15px; /* Increase the size of the icon */
    height: 15px; /* Increase the size of the icon */
    margin-right: 10px; /* Increase space between the icon and text */
  }

  .conversations {
    padding: 20px;
  }
  .user-finder {
    margin-bottom: 20px;
    display: flex;
    gap: 10px;
  }
  .search-bar {
    width: 80%
  }
  .search-button {
    display: flex;
    padding: 4px 8px;
    font-size: 0.9rem;
  }
  .search-results {
    margin-bottom: 20px;
  }
  .conversation-item,
  .search-item {
    display: flex;
    align-items: center;
    gap: 10px;
    margin-bottom: 10px;
  }
  .profile-photo {
    width: 80px;
    height: 80px;
    border-radius: 50%;
    align-content:center;
  }
  .user-details {
    flex: 1;
  }
  .profile-picture-upload {
    margin: 25px;
  }
  .upload-button {
    margin-top: 25px;
    margin-bottom: 25px;
  }
  .profile-preview {
    margin-top: 20px;
    max-width: 150px;
    border-radius: 50%;
  }
  </style>
