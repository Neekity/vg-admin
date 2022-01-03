<template>
      <v-card elevation="2" class="mx-auto my-12" max-width="374">
        <v-card-text class="text-center">
          <v-btn
              color="blue darken-1"
              text
              @click="oauthLogin"
          >
           使用OAuth登录
          </v-btn>
        </v-card-text>
      </v-card>
</template>
<script>
import {mapActions} from 'vuex';
import {gethttp} from "../../store/module/user";

export default {
  data() {
    return {
      subTitle: "登录",
      overlayShow: true,
      showLoginForm: true
    }
  },
  created() {
    if (this.$store.state.user.token && !gethttp('token').token) {
      this.$router.push({
        path: '/example'
      })
    } else if (gethttp('token').token) {
      this.showLoginForm = false;
      this.handleSubmit();
    }
  },
  methods: {
    ...mapActions([
      'handleLogin',
      'getUserInfo'
    ]),
    handleSubmit() {
      this.handleLogin({}).then(() => {
        if (!this.$store.state.user.token) {
          location.href = this.$loginPage;
        } else {
          this.getUserInfo({}).then(() => {
            this.$router.push({
              path: '/example'
            })
          })
        }
      })
    },
    oauthLogin() {
      location.href = this.$loginPage;
    }
  }
}
</script>