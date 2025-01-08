<template>
  <div class="dashboard_container">
    <h2>Welcome, {{ username }}</h2>
    <button @click="logout">Logout</button>
  </div>
</template>

<script>
export default {
  data() {
    return {
      username: localStorage.getItem("username") || "Guest",
    };
  },
  methods: {
    logout() {
      localStorage.removeItem("token");
      localStorage.removeItem("username");
      this.$router.push("../dashboard");
    },
  },
};
</script>

<style scoped>
.dashboard_container {
  text-align: center;
  margin: 50px auto;
}
</style>

