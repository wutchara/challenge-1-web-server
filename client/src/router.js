import Vue from 'vue'
import Router from 'vue-router'

import Users from './views/users.vue';
import Home from './views/Home.vue';
import NotFound from './components/404.vue';

Vue.use(Router)

export default new Router({
    mode: 'history',
    base: process.env.BASE_URL,
    routes: [{
        path: '/users',
        name: 'users',
        component: Users
    },
    {
        path: '/',
        name: 'home',
        component: Home
    },
    {
        path: '*',
        component: NotFound
    },
],
});