<template>
  <v-app>
    <v-card class="my-2">
      <v-card-title>
        <v-row>
          <v-col cols="12"
                 sm="6"
          >
            <v-text-field
                v-model="searchRoleName"
                label="请输入名称搜索"
                outlined
                dense
            >
              <v-icon
                  slot="append"
                  color="blue"
                  @click="getRoles"
              >
                mdi-magnify
              </v-icon>
            </v-text-field>
          </v-col>
          <v-spacer></v-spacer>
          <v-btn
              icon
              color="blue darken-2"
              @click="roleId=0;showEditDialog=true;"
          >
            <v-icon>
              mdi-plus
            </v-icon>
          </v-btn>
        </v-row>
      </v-card-title>
      <v-data-table
          :headers="headers"
          :items="roleItems"
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

          <v-icon
              small
              class="mr-2"
              color="blue darken-2"
              @click="roleId=item.id;showEditDialog=true;"
          >
            mdi-pencil
          </v-icon>

        </template>
      </v-data-table>
    </v-card>
    <v-dialog v-model="showEditDialog">
      <v-card class="mx-auto">
        <v-card-title>角色操作</v-card-title>
        <v-card-text>
          <v-row>
            <v-col cols="12" md="4" sm="6">
              <v-text-field
                  v-model="roleName"
                  dense
                  outlined
                  label="角色名称"
              ></v-text-field>
            </v-col>
            <v-col cols="12" md="4" sm="6">
              <v-text-field
                  v-model="casbinRole"
                  dense
                  outlined
                  label="casbin角色key"
                  :disabled="!!roleId"
              ></v-text-field>
            </v-col>
          </v-row>
          <v-row class="text-left">
            <v-col md="4" sm="6">
              <v-text-field
                  v-model="search"
                  label="输入关键字搜索权限"
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
            <v-col md="4" sm="6">
              <div class="text-h6 pb-2">分配权限至角色</div>
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
              @click="store"
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
  name: "RoleList",
  data() {
    return {
      headers: [{
        text: '角色名称',
        align: 'start',
        value: 'name',
      },
        {text: 'casbin角色key', value: 'casbin_role', sortable: false},
        {text: '操作', value: 'actions', sortable: false},
      ],
      roleItems: [],
      searchRoleName: '',
      roleId: 0,
      showEditDialog: false,
      overlay: 0,
      open: [],
      allOpened: false,
      lastOpen: [],
      search: null,
      selection: [],
      items: [{id: 0, name: '全选', children: []}],
      roleName: "",
      casbinRole: "",
    }
  },
  watch: {
    roleId: {
      handler() {
        this.search=null;
        this.lastOpen=null;
        this.open=null;
        this.allOpened=false;
        this.lastOpen=false;
        this.getRole();
      }
    }
  },
  created() {
    this.getRoles()
    this.getAllPermissions();
  },
  methods: {
    getAllPermissions() {
      this.overlay += 1;
      this.$http('post','permission', {name: ''})
          .then(response => {
            let resData = response.data;
            if (resData.code === 0) {
              this.items[0].children = resData.data;
            } else {
              this.$toast('获取所有权限出错：' + resData.message, {
                type: 'error',
              });
            }
          })
          .catch(error => {
            console.log(error);
            this.$toast('获取所有权限出错：服务器出错！', {
              type: 'error',
              timeout: 2000,
            });
          })
          .finally(() => this.overlay -= 1);
    },
    getRole() {
      if (this.roleId === 0) {
        this.roleName = '';
        this.casbinRole = '';
        this.selection = [];
        return;
      }
      this.overlay += 1;
      this.$http('post','role/detail', {id: this.roleId})
          .then(response => {
            let resData = response.data;
            if (resData.code === 0) {
              this.casbinRole = resData.data.casbin_role;
              this.roleName = resData.data.name;
              this.selection = resData.data.permissions || [];
            } else {
              this.$toast('获取角色详情出错：' + resData.message, {
                type: 'error',
              });
            }
          })
          .catch(error => {
            console.log(error);
            this.$toast('获取角色详情出错：服务器出错！', {
              type: 'error',
              timeout: 2000,
            });
          })
          .finally(() => this.overlay -= 1);
    },
    store() {
      this.overlay += 1;
      let permissions = [];
      this.selection.forEach(function (element) {
        permissions.push(element.id)
      });
      this.$http('post','role/store', {
        id: this.roleId,
        name: this.roleName,
        casbin_role: this.casbinRole,
        permission_ids: permissions,
      })
          .then(response => {
            let resData = response.data;
            if (resData.code === 0) {
              this.$toast("保存成功！", {
                type: 'success',
                timeout: 2000
              });
              this.getRoles();
              this.showEditDialog=false;
            } else {
              this.$toast('保存角色出错：' + resData.message, {
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
    getRoles() {
      this.overlay += 1;
      this.$http('post','role', {name: this.searchRoleName})
          .then(response => {
            let resData = response.data;
            if (resData.code === 0) {
              this.roleItems = resData.data;
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
  }
}
</script>

<style scoped>

</style>