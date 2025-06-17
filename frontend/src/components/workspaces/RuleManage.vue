<template>
  <div class="py-1 px-3">
    <div>
      <el-button type="primary" size="small" @click="handleAdd">Add Rule</el-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import _ from 'lodash';
import { reactive, onMounted, ref } from 'vue';

const props = defineProps<{
  server: Server
}>()

const state = reactive({
  form: {
    showDrawer: false,
    loading: false,
    data: {
    } as any,
    rules: {
    },
    action: '',   // create or update
  },
})

onMounted(() => {
});


const handleAdd = () => {
  state.form.action = 'create'
  state.form.showDrawer = true
}
const handleEdit = (route: any) => {
  state.form.action = 'update'
  state.form.showDrawer = true
  state.form.data = route
}
const handleDelete = (node: any, route: any) => {
}
const handleSubmit = async () => {
  
}
</script>

<style scoped>
</style>