<template>
  <v-app>
    <v-card class="my-2">
      <v-card-title>
        <v-row>
          <v-col cols="12"
                 sm="6"
          >
            <v-text-field
                v-model="userName"
                label="请输入名称"
                outlined
                dense
            >
              <v-icon
                  slot="append"
                  color="blue"
                  @click="getUsers"
              >
                mdi-magnify
              </v-icon>
            </v-text-field>
          </v-col>
          <v-spacer></v-spacer>
        </v-row>
      </v-card-title>
      <v-data-table
          :headers="headers"
          :items="userItems"
          item-key="id"
          outlined
          :footer-props="{
              showFirstLastPage: true,
              firstIcon: 'mdi-arrow-collapse-left',
              lastIcon: 'mdi-arrow-collapse-right',
              prevIcon: 'mdi-minus',
              nextIcon: 'mdi-plus'
            }"
      >
        <template v-slot:item.actions="{ item }">
          <v-tooltip bottom>
            <template v-slot:activator="{ on, attrs }">
              <v-btn
                  icon
                  v-bind="attrs"
                  v-on="on"
                  class="mr-2"
                  color="blue darken-2"
                  @click="userId=item.id;showAssignDialog=true;"
              >
                <v-icon>
                  mdi-account-key-outline
                </v-icon>
              </v-btn>
            </template>
            <span>分配角色</span>
          </v-tooltip>
        </template>
      </v-data-table>
    </v-card>
    <v-dialog v-model="showAssignDialog">
      <v-card class="mx-auto">
        <v-card-title>分配角色至用户</v-card-title>
        <v-card-text>
          <v-row class="text-left">
            <v-col>
              <v-text-field
                  v-model="search"
                  label="输入关键字搜索角色"
                  clearable
                  dense
                  outlined
                  @input="handleSearch"
                  clear-icon="mdi-close-circle-outline"
              ></v-text-field>
              <v-treeview
                  ref="items"
                  v-model="selection"
                  :items="items"
                  selectable
                  return-object
                  dense
                  hoverable
                  :open.sync="open"
                  :search="search"
              ></v-treeview>
            </v-col>
            <v-divider vertical></v-divider>
            <v-col>
              <v-scroll-x-transition
                  group
                  hide-on-leave
              >
                <v-chip
                    v-for="(node, i) in selection"
                    :key="i"
                    color="primary"
                    outlined
                    small
                    close
                    class="ma-1"
                    @click:close="selection.splice(i,1)"
                >
                  <v-icon
                      left
                      small
                  >
                    mdi-account
                  </v-icon>
                  {{ node.name }}
                </v-chip>
              </v-scroll-x-transition>
            </v-col>
          </v-row>
        </v-card-text>
        <v-card-actions>
          <v-btn
              color="success"
              @click="assign"
              class="text-center"
          >
            保存
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
    <v-overlay :value="overlay">
      <v-progress-circular
          indeterminate
          color="primary"
      ></v-progress-circular>
    </v-overlay>
  </v-app>
</template>

<script>
export default {
  name: "UserList",
  data() {
    return {
      headers: [{
        text: '用户名称',
        align: 'start',
        value: 'name',
      },
        {text: '用户邮箱', value: 'email'},
        {text: '用户电话', value: 'mobile'},
        {text: '操作', value: 'actions', sortable: false},
      ],
      userItems: [],
      userName: '',
      userId:0,
      open: [],
      allOpened: false,
      lastOpen: [],
      search: null,
      selection: [],
      items: [{id: 0, name: '全选', children: []}],
      showAssignDialog:false,
    }
  },
  watch: {
    userId: {
      handler() {
        this.search=null;
        this.lastOpen=null;
        this.open=null;
        this.allOpened=false;
        this.lastOpen=false;
        this.getUserRoles();
      }
    }
  },
  created() {
    this.getUsers();
    this.getAllRoles();
  },
  methods: {
    assign(){
      this.overlay += 1;
      let casbinRoles = [];
      this.selection.forEach(function (element) {
        casbinRoles.push(element.casbin_role)
      });
      this.$http('post','/user/role/assign', {
        id: this.userId,
        casbin_roles: casbinRoles,
      })
          .then(response => {
            let resData = response.data;
            if (resData.code === 0) {
              this.$toast("分配成功！", {
                type: 'success',
                timeout: 2000
              });
              this.showAssignDialog=false;
            } else {
              this.$toast('分配角色出错：' + resData.message, {
                type: 'error',
              });
            }
          })
          .catch(error => {
            console.log(error);
            this.$toast('保存角色出错：服务器出错！', {
              type: 'error',
              timeout: 2000,
            });
          })
          .finally(() => this.overlay -= 1);
    },
    getAllRoles() {
      this.overlay += 1;
      this.$http('post','/role', {name: ''})
          .then(response => {
            let resData = response.data;
            if (resData.code === 0) {
              this.items[0].children = resData.data;
            } else {
              this.$toast('获取角色出错：' + resData.message, {
                type: 'error',
              });
            }
          })
          .catch(error => {
            console.log(error);
            this.$toast('获取角色出错：服务器出错！', {
              type: 'error',
              timeout: 2000,
            });
          })
          .finally(() => this.overlay -= 1);
    },
    getUserRoles(){
      if (this.userId === 0) {
        this.selection = [];
        return;
      }
      this.overlay += 1;
      this.$http('post','/user/roles', {user_id: this.userId})
          .then(response => {
            let resData = response.data;
            if (resData.code === 0) {
              this.selection = resData.data;
              console.log(this.selection)
            } else {
              this.$toast('获取用户角色出错：' + resData.message, {
                type: 'error',
              });
            }
          })
          .catch(error => {
            console.log(error);
            this.$toast('获取用户角色出错：服务器出错！', {
              type: 'error',
              timeout: 2000,
            });
          })
          .finally(() => this.overlay -= 1);
    },
    getUsers() {
      this.overlay += 1;
      this.$http('post','user', {name: this.userName})
          .then(response => {
            let resData = response.data;
            if (resData.code === 0) {
              this.userItems = resData.data;
            } else {
              this.$toast('获取用户出错：' + resData.message, {
                type: 'error',
              });
            }
          })
          .catch(error => {
            console.log(error);
            this.$toast('获取用户出错：服务器出错！', {
              type: 'error',
              timeout: 2000,
            });
          })
          .finally(() => this.overlay -= 1);
    },
    handleSearch: function (val) {
      if (val) {
        if (!this.allOpened) {
          this.lastOpen = this.open;
          this.allOpened = true;
          this.$refs.items.updateAll(true);
        }
      } else {
        this.$refs.items.updateAll(false);
        this.allOpened = false;
        this.open = this.lastOpen;
      }
    },
  }
}
</script>

<style scoped>

</style>