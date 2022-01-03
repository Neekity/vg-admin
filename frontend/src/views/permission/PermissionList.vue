<template>
  <v-app>
    <v-card class="my-2">
      <v-card-title>
        <v-row>
          <v-col cols="12"
                 sm="6"
          >
            <v-text-field
                v-model="searchPermissionName"
                label="请输入名称搜索"
                outlined
                dense
            >
              <v-icon
                  slot="append"
                  color="blue"
                  @click="getPermissions"
              >
                mdi-magnify
              </v-icon>
            </v-text-field>
          </v-col>
          <v-spacer></v-spacer>
          <v-btn
              icon
              color="blue darken-2"
              @click="permissionId=0;showEditDialog=true;"
          >
            <v-icon>
              mdi-plus
            </v-icon>
          </v-btn>
        </v-row>
      </v-card-title>
      <v-data-table
          :headers="headers"
          :items="permissionItems"
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
              @click="permissionId=item.id;showEditDialog=true;"
          >
            mdi-pencil
          </v-icon>

        </template>
      </v-data-table>
    </v-card>
    <v-dialog v-model="showEditDialog">
      <v-card class="mx-auto">
        <v-card-title>权限操作</v-card-title>
        <v-card-text>
          <v-row>
            <v-col cols="12" md="4" sm="6">
              <v-text-field
                  v-model="permissionName"
                  dense
                  outlined
                  label="权限名称"
              ></v-text-field>
            </v-col>
            <v-col cols="12" md="4" sm="6">
              <v-text-field
                  v-model="casbinPermission"
                  dense
                  outlined
                  label="casbin权限key"
                  :disabled="!!permissionId"
              ></v-text-field>
            </v-col>
            <v-col cols="12" md="4" sm="6">
              <v-text-field
                  v-model="casbinPermissionType"
                  dense
                  outlined
                  label="casbin权限类型"
                  :disabled="!!permissionId"
              ></v-text-field>
            </v-col>
            <v-col cols="12" md="4" sm="6">
              <v-text-field
                  v-model="permissionRoute"
                  dense
                  outlined
                  label="路由"
              ></v-text-field>
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
  name: "PermissionList",
  data() {
    return {
      headers: [{
        text: '权限名称',
        align: 'start',
        value: 'name',
      },
        {text: 'casbin权限key', value: 'casbin_permission', sortable: false},
        {text: 'casbin权限类型', value: 'casbin_permission_type', sortable: false},
        {text: '路由', value: 'route', sortable: false},
        {text: '操作', value: 'actions', sortable: false},
      ],
      permissionRoute: '',
      searchPermissionName: '',
      permissionId: 0,
      showEditDialog: false,
      overlay: 0,
      permissionName: "",
      casbinPermission: "",
      casbinPermissionType: "",
      permissionItems: [],
    }
  },
  watch: {
    permissionId: {
      handler() {
        this.getPermission();
      }
    }
  },
  created() {
    this.getPermissions()
  },
  methods: {
    getPermission() {
      if (this.permissionId === 0) {
        this.permissionName = '';
        this.casbinPermission = '';
        this.casbinPermissionType = '';
        this.permissionRoute = '';
        return;
      }
      this.overlay += 1;
      this.$http('post', 'permission/detail', {id: this.permissionId})
          .then(response => {
            let resData = response.data;
            if (resData.code === 0) {
              this.casbinPermission = resData.data.casbin_permission;
              this.permissionRoute = resData.data.route;
              this.permissionName = resData.data.name;
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
      this.$http('post', 'permission/store', {
        id: this.permissionId,
        name: this.permissionName,
        casbin_permission: this.casbinPermission,
        route: this.permissionRoute,
        casbin_permission_type: this.casbinPermissionType,
      })
          .then(response => {
            let resData = response.data;
            if (resData.code === 0) {
              this.$toast("保存成功！", {
                type: 'success',
                timeout: 2000
              });
              this.getPermissions();
              this.showEditDialog = false;
            } else {
              this.$toast('保存权限出错：' + resData.message, {
                type: 'error',
              });
            }
          })
          .catch(error => {
            console.log(error);
            this.$toast('保存权限出错：服务器出错！', {
              type: 'error',
              timeout: 2000,
            });
          })
          .finally(() => this.overlay -= 1);
    },

    getPermissions() {
      this.overlay += 1;
      this.$http('post', 'permission', {name: this.searchPermissionName})
          .then(response => {
            let resData = response.data;
            if (resData.code === 0) {
              this.permissionItems = resData.data;
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