<template>
  <div>
    <b-button
    variant="outline-primary"
    href="/users">
    Goto: users
    </b-button>
    <HelloWorld />
  </div>
</template>

<script>
import HelloWorld from '../components/HelloWorld.vue'

export default {
  name: "App",
  components: {
    HelloWorld,
  },
};
</script>