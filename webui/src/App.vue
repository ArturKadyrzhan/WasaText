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
                <h3>Upload Profile Picture</h3>
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
                      <img :src="getProfileImage(user.ProfilePhotoURL)" alt="Profile" class="profile-photo" />
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
                    <img :src="getProfileImage(user.ProfilePhotoURL)" alt="Profile" class="profile-photo" />
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
                    <img :src="getProfileImage(group.GroupPhotoURL)" alt="Group Profile" class="profile-photo" />
                    <div class="user-details">
                      <h5>{{ group.Name }}</h5>
                      <button @click="startConversation(group)">Open Chat</button>
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

 <script setup>
 import { RouterLink, RouterView } from 'vue-router'
 </script>
 <script>
 import {checkLoginStatus, getToken, logIn, logOut} from "./store/auth";
 import {getProfileImage} from "./services/helpers";

 export default {
   data: function() {
     return {
       userSearchResult: [],
       users: [],
       groups:[],
       errormsg: null,
       loading: false,
       username: "",
       password: "",
       checkLogin: false,
       searchQuery: "", // Input for searching users
       selectedFile: null,
       previewUrl: null,
     }
   },
   methods: {
     async findUser() {
       if (!this.searchQuery) {
         this.userSearchResult = [];
         return;
       }

       try {
         let response = await this.$axios.get(`/get-users?search=${this.searchQuery}`, {
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
     async login() {
       this.loading = true;
       this.errormsg = null;
       try {
         let response = await this.$axios.post("/login", {
           username: this.username,
           password: this.password,
         });
         logIn(response.data['token'], response.data['id'], response.data['username'])
         window.location.reload();
       } catch (e) {
         console.log(e,"eror meror")
         this.errormsg = "Login failed. Please check your credentials.";
       }
       this.loading = false;
     },
     async fetchConversations() {
       try {
          console.log('Token:',getToken())
         let response = await this.$axios.post('/get-conversations',null,{
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
     startConversation(userOrGroup) {
       if (!userOrGroup) {
         console.error('Missing user ID. Cannot navigate to conversation.');
         return;
       }
       const isGroup = userOrGroup.Name !== undefined;
       // Getting url manually
       const entity = encodeURIComponent(JSON.stringify(userOrGroup));
       const isGroupParam = isGroup ? 'true' : 'false';

       // changing URL
       window.location.href = `#/conversation?entity=${entity}&isGroup=${isGroupParam}`;
       // Reset page
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
     onFileChange(event) {
       const file = event.target.files[0];
       if (file) {
         this.selectedFile = file;
         this.previewUrl = URL.createObjectURL(file);
       }
     },
     async uploadImage() {
       const formData = new FormData();
       formData.append("profile_picture", this.selectedFile);

       try {
         const response = await this.$axios.post("/upload-profile-picture", formData, {
           headers: {
             "Content-Type": "multipart/form-data",
             Authorization: `Bearer ${getToken()}`,
           },
         });

         console.log("Image uploaded:", response.data);
         alert("Profile picture uploaded successfully!");
         window.location.reload();
       } catch (error) {
         console.error("Error uploading image:", error);
         alert("Failed to upload image.");
       }
     },
     async leaveGroup(group) {
       try {
         await this.$axios.post("/leave-group",{
            group:group
         }, {
           headers: {
             Authorization: `Bearer ${getToken()}`,
           },
         });

         window.location.reload();
       } catch (error) {
         console.error("Error uploading image:", error);
         alert("Failed to upload image.");
       }
     }
   },
   mounted() {
     this.checkLogin = checkLoginStatus()
     if (this.checkLogin){
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
    width: 40px;
    height: 40px;
    border-radius: 50%;
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
