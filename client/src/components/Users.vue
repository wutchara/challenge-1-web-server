<template>
  <div class="p-4">
    <h1>Users List</h1>
    <b-overlay :show="isLoading" rounded="sm">
      <b-button class="m-2" variant="success" @click="addNewUser()">Add new User</b-button >
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
            <b-card
              border-variant="primary"
              header-bg-variant="primary"
              header-text-variant="white"
              :header="user.name"
            >
              <b-form-group label="Name" class="mb-1">
                <b-form-input type="text" :value="user.name" :name="'name-' + user.id"></b-form-input>
              </b-form-group>
              <br />
              <b-form-group label="Value" class="mb-1">
                <b-form-input type="number" :value="user.value" :name="'value-' + user.id"></b-form-input>
              </b-form-group>
              <template #footer>
                <b-button class="m-2" variant="outline-warning" @click="editUser(user.id)"
                  >Edit</b-button
                >
                <b-button class="m-2" variant="outline-danger" @click="deleteUser(user.id)"
                  >Delete</b-button
                >
              </template>
            </b-card>
            <br />
          </div>
        </b-card-group>
      </div>
    </b-overlay>
  </div>
</template>

<script>
export default {
  name: "HelloWorld",
  props: {
    users: Array,
  },
  data() {
    return {
      isLoading: false,
    };
  },
  methods: {
    findUser(id) {
      return this.users.find((u) => u.id === id);
    },
    troggleLoading() {
      this.isLoading = !this.isLoading;
    },
    editUser(id) {
      this.troggleLoading();
      const user = this.findUser(id);
      console.log("EDIT:USER", user);

      if (user) {
        // check user data has changed
        const currentData = {
          id,
          name: document.querySelector(`input[name=name-${id}]`).value,
          value: Number(document.querySelector(`input[name=value-${id}]`).value),
        };

        if (currentData.name !== user.name || currentData.value !== user.value) {
          console.log('.....Try to update user.....');
          console.log('currentData', currentData);

          this.axios({
            method: "post",
            url: "http://127.0.0.1:3000/user/" + currentData.id,
            data: JSON.stringify({
              name: currentData.name,
              value: currentData.value
            }),
            responseType: "arraybuffer",
            headers: { 
              'Content-Type': 'application/json'
            },
          }).then(response => {
            var enc = new TextDecoder("utf-8");
            return JSON.parse(enc.decode(response.data));
          }).then(response => {
            console.log('response', response);
            if (!response.status) {
              throw "Cannot update user";
            }
            return response;
          }).catch(e => {
            console.log('ERROR:Update user', e);
          }).finally(() => {
            console.log("......... Load data compleated .........");
            setTimeout(() => {
              this.troggleLoading();
            }, 1000);
          });
        } else {
          this.troggleLoading();
        }
      } else {
        this.troggleLoading();
      }
    },
    deleteUser(id) {
      this.troggleLoading();
      const user = this.findUser(id);
      console.log("DELETE:USER", user);

      if (user) {
        this.axios({
              method: "delete",
              url: "http://127.0.0.1:3000/user/" + id,
              responseType: "arraybuffer",
              headers: { 
                'Content-Type': 'application/json'
              },
            }).then(response => {
              var enc = new TextDecoder("utf-8");
              return JSON.parse(enc.decode(response.data));
            }).then(response => {
              console.log('response', response);
              if (!response.status) {
                throw "Cannot Delete user";
              }

              // remove from array
              const index = this.users.findIndex(u => u.id === response.data.id);
              console.log('index', index);
              console.log('-----this.users', this.users);
              if (index > -1) {
                this.users.splice(index, 1);
              }
              console.log('this.users-----', this.users);

              return response;
            }).catch(e => {
              console.log('ERROR:Delete user', e);
            }).finally(() => {
              console.log("......... Load data compleated .........");
              setTimeout(() => {
                this.troggleLoading();
              }, 1000);
            });
      } else {
        this.troggleLoading();
      }
    },
    addNewUser() {
      this.troggleLoading();
      this.axios({
            method: "put",
            url: "http://127.0.0.1:3000/user",
            data: JSON.stringify({
              name: '',
              value: 0
            }),
            responseType: "arraybuffer",
            headers: { 
              'Content-Type': 'application/json'
            },
          }).then(response => {
            var enc = new TextDecoder("utf-8");
            return JSON.parse(enc.decode(response.data));
          }).then(response => {
            console.log('response', response);
            if (!response.status) {
              throw "Cannot Add new user";
            }

            this.users.push(response.data);
            // scroll down
            window.scrollTo(0,document.body.scrollHeight);

            return response;
          }).catch(e => {
            console.log('ERROR:Add new user', e);
          }).finally(() => {
            console.log("......... Load data compleated .........");
            setTimeout(() => {
              this.troggleLoading();
            }, 1000);
          });
    },
  },
  beforeCreate() {
  },
};
</script>
