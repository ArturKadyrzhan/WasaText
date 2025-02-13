
  <script>
import { getToken } from "../store/auth";

export default {
  data() {
    return {
      groupId: null,
      updatedGroupName: "",
      updatedProfilePhoto: null,
    };
  },
  async created() {
    this.groupId = parseInt(this.$route.query.groupId, 10);

    if (!this.groupId) {
      alert("Invalid group ID.");
      this.$router.push("/conversation");
    }
  },
  methods: {
    handleProfilePhotoUpload(event) {
      this.updatedProfilePhoto = event.target.files[0];
    },

    async updateGroupName() {
      if (!this.updatedGroupName) {
        alert("Please enter a new group name.");
        return;
      }

      try {
        await this.$axios.post(
            "/group/name",
            {
              id: this.groupId,
              groupName: this.updatedGroupName,
            },
            {
              headers: { Authorization: `Bearer ${getToken()}` },
            }
        );

        alert("Group name updated successfully!");
      } catch (error) {
        console.error("Error updating group name:", error);
        alert("Failed to update group name.");
      }
    },

    async updateGroupPhoto() {
      if (!this.updatedProfilePhoto) {
        alert("Please select a new group photo.");
        return;
      }

      const formData = new FormData();
      formData.append("id", this.groupId);
      formData.append("groupPhoto", this.updatedProfilePhoto);

      try {
        await this.$axios.post("/group/photo", formData, {
          headers: {
            Authorization: `Bearer ${getToken()}`, 
            "Content-Type": "multipart/form-data",
          },
        });

        alert("Group photo updated successfully!");
      } catch (error) {
        console.error("Error updating group photo:", error);
        alert("Failed to update group photo.");
      }
    }
  }
};
</script>


<template>
  <div class="update-group">
    <h3>Update Group</h3>

    <!-- Update Group Name -->
    <form @submit.prevent="updateGroupName">
      <div>
        <label for="groupName">Group Name:</label>
        <input type="text" id="groupName" v-model="updatedGroupName" required />
      </div>
      <button type="submit">Update Group Name</button>
    </form>

    <!-- Update Group Profile Photo -->
    <form @submit.prevent="updateGroupPhoto">
      <div>
        <label for="profilePhoto">Group Profile Photo:</label>
        <input type="file" id="profilePhoto" @change="handleProfilePhotoUpload" />
      </div>
      <button type="submit">Update Group Photo</button>
    </form>
  </div>
</template>

    <style scoped>
    .update-group {
      padding: 20px;
      border: 1px solid #ccc;
      margin: 20px;
      border-radius: 8px;
    }

    .update-group input,
    .update-group button {
      margin: 10px 0;
      padding: 10px;
      width: 100%;
      box-sizing: border-box;
    }

    .update-group button {
      background-color: #0088cc;
      color: white;
      border: none;
      cursor: pointer;
    }

    .update-group button:hover {
      background-color: #005f99;
    }
    </style>