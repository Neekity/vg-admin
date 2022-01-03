<template>
  <v-container>
    <v-card>
      <v-card-title>
        <v-spacer></v-spacer>
        <v-dialog
            v-model="dialog"
            max-width="800px"
        >
          <template v-slot:activator="{ on, attrs }">
            <v-btn
                color="primary"
                dark
                class="mb-2"
                v-bind="attrs"
                v-on="on"
                icon
            >
              <v-icon>mdi-plus-thick</v-icon>
            </v-btn>
          </template>
          <v-card>
            <v-card-title>
              <span class="headline">菜单</span>
            </v-card-title>
            <v-card-text>
              <v-container>
                <v-row>
                  <v-col
                      cols="12"
                      sm="6"
                      md="3"
                  >
                    <v-text-field
                        type="number"
                        v-model="editedItem.parent_id"
                        label="父级菜单id"
                    ></v-text-field>
                  </v-col>
                  <v-col
                      cols="12"
                      sm="6"
                      md="3"
                  >
                    <v-text-field
                        v-model="editedItem.name"
                        label="名称"
                    ></v-text-field>
                  </v-col>
                  <v-col
                      cols="12"
                      sm="6"
                      md="3"
                  >
                    <v-text-field
                        v-model="editedItem.icon"
                        label="图标"
                    ></v-text-field>
                  </v-col>
                  <v-col
                      cols="12"
                      sm="6"
                      md="3"
                  >
                    <v-text-field
                        v-model="editedItem.path"
                        label="路由"
                    ></v-text-field>
                  </v-col>
                </v-row>
              </v-container>
            </v-card-text>

            <v-card-actions>
              <v-spacer></v-spacer>
              <v-btn
                  color="blue darken-1"
                  text
                  @click="close"
              >
                关闭
              </v-btn>
              <v-btn
                  color="blue darken-1"
                  text
                  @click="save"
              >
                保存
              </v-btn>
            </v-card-actions>
          </v-card>
        </v-dialog>
      </v-card-title>
      <v-card-text>
        <v-treeview
            :open="initiallyOpen"
            activatable
            :items="items"
            item-text="name"
            open-on-click
        >
          <template v-slot:prepend="{ item}">
            <v-icon v-if="!item.icon">
              mdi-folder
            </v-icon>
            <v-icon v-else>
              {{ item.icon }}
            </v-icon>
          </template>
          <template v-slot:append="{ item}">
            <v-btn icon
                   color="blue darken-2"
                   @click="edit(item)">
              <v-icon>mdi-pencil</v-icon>
            </v-btn>
            <v-btn icon
                   color="red darken-2"
                   @click="deleteMenu(item.id)">
              <v-icon>mdi-delete</v-icon>
            </v-btn>
          </template>
        </v-treeview>
      </v-card-text>
    </v-card>
    <v-overlay :value="overlay">
      <v-progress-circular
          indeterminate
          color="primary"
      ></v-progress-circular>
    </v-overlay>
  </v-container>
</template>

<script>
export default {
  data: () => ({
    editedItem: {
      id: 0,
      parent_id: 0,
      name: "",
      icon: "",
      path: "",
    },
    defaultItem: {
      id: 0,
      parent_id: 0,
      name: "",
      icon: "",
      path: ""
    },
    initiallyOpen: [1],
    dialog: false,
    items: [],
    overlay: 0,
  }),
  methods: {
    edit(item) {
      this.dialog = true;
      this.editedItem = item;
    },
    getMenus() {
      this.overlay += 1;
      this.$http('post', '/menu/list').then((response) => {
        let resData = response.data;
        if (resData.code === 0) {
          this.items = resData.data || [];
        } else {
          this.$toast('获取菜单出错：' + resData.message, {
            type: 'error',
          });
        }
      }).catch(error => {
        console.log(error);
        this.$toast('获取菜单出错：服务器出错！', {
          type: 'error',
          timeout: 2000,
        });
      }).finally(() => {
        this.overlay -= 1;
      });
    },
    save() {
      this.overlay += 1;
      console.log(this.editedItem)
      this.editedItem.parent_id = parseInt(this.editedItem.parent_id)
      console.log(this.editedItem)
      this.$http('post', '/menu/store', this.editedItem).then((response) => {
        let resData = response.data;
        if (resData.code === 0) {
          this.$toast('保存菜单成功！', {
            type: 'success',
          });
          this.getMenus();
        } else {
          this.$toast('保存菜单出错：' + resData.message, {
            type: 'error',
          });
        }
      }).catch(error => {
        console.log(error);
        this.$toast('保存菜单出错：服务器出错！', {
          type: 'error',
          timeout: 2000,
        });
      }).finally(() => {
        this.overlay -= 1;
      });
      this.close()
    },
    close() {
      this.dialog = false
    },
    deleteMenu(id) {
      this.overlay += 1;
      this.$http('post', '/menu/delete', {id: id}).then((response) => {
        let resData = response.data;
        if (resData.code === 0) {
          this.$toast('删除菜单成功！', {
            type: 'success',
          });
          this.getMenus();
        } else {
          this.$toast('删除菜单出错：' + resData.message, {
            type: 'error',
          });
        }
      }).catch(error => {
        console.log(error);
        this.$toast('删除菜单出错：服务器出错！', {
          type: 'error',
          timeout: 2000,
        });
      }).finally(() => {
        this.overlay -= 1;
      });
    }
  },
  watch: {
    dialog(val) {
      val || this.close()
    },
  },
  created() {
    this.getMenus();
  },

}
</script>