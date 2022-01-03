<template>
  <v-tooltip bottom>
    <template v-slot:activator="{ on, attrs }">
       <span
           v-bind="attrs"
           v-on="on"
           :class="'d-inline-block text-truncate '+textClass"
           :style="'max-width: '+maxWidth"
       >
         {{ text }}
       </span>
    </template>
    <span> {{ text }}</span>
  </v-tooltip>
</template>

<script>
export default {
  name: "TextTruncate",
  props: {
    text: {
      type: String,
      default: '',
    },
    textClass: {
      type: String,
      default: '',
    },
    maxWidth: {
      type: String,
      default: '150px',
    },
  },
}
</script>

<style scoped>

</style>