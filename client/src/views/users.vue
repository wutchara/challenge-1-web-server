<template>
  <div>
    <div v-show="isLoading">
      <b-spinner variant="primary" type="grow"></b-spinner>
      <b-spinner variant="primary"></b-spinner>
      <b-spinner variant="primary" type="grow"></b-spinner>
    </div>
    <Users :users="users" />
  </div>
</template>

<script>
import Users from "../components/Users.vue";

export default {
  name: "App",
  components: {
    Users,
  },
  data() {
    return {
      isLoading: true,
      users: [],
    };
  },
  created() {
    this.isLoading = true;
    console.log("this.axios", this.axios);
    this.axios({
      method: "get",
      url: "http://127.0.0.1:3000/users",
      responseType: "arraybuffer",
    })
      .then((response) => {
        // console.log('response', response);

        var enc = new TextDecoder("utf-8");
        return JSON.parse(enc.decode(response.data));
      })
      .then((res) => {
        console.log("res", res);
        this.users = (res?.data || []).reverse();
      })
      .catch((err) => {
        console.error(err);
      })
      .finally(() => {
        console.log("......... Load data compleated .........");
        this.isLoading = false;
      });
  },
};
</script>