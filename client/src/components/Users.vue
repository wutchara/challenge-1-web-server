<template>
  <div class="p-4">
    <h1>Users List</h1>
    {{ users }}

    <!-- <div class="p-4">
      <b-list-group>
        <b-list-group-item
          button
          class="align-items-start"
          v-for="(user, index) in users"
          :key="index"
        >
          {{ user.name }}
        </b-list-group-item>
      </b-list-group>
    </div> -->

    <div class="p-4">
      <b-card-group deck>
        <div v-for="(user, index) in users" :key="index">
          <b-card :header="user.name">
            <b-form-group
              label="Name"
              class="mb-1"
            >
              <b-form-input
                :value="user.name"
                disabled=true
              ></b-form-input>
            </b-form-group>
            <br />
            <b-form-group
              label="Value"
              class="mb-1"
            >
              <b-form-input
                :value="user.value"
                disabled=true
              ></b-form-input>
            </b-form-group>
          </b-card>
          <br />
        </div>
      </b-card-group>
    </div>
  </div>
</template>

<script>
export default {
  name: "HelloWorld",
  props: {
    users: Array,
  },
};
</script>
