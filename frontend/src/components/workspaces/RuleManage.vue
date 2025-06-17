<template>
  <TopInfo :workspace="workspaceStore.currentWorkspace">
    <el-button type="primary" @click="handleAddService">Add Rule</el-button>
  </TopInfo>

  <div class="section-box-dark mb-3">
    <el-collapse accordion>
      <el-collapse-item name="1">
        <template #title="{ isActive }">
          <div :class="['title-wrapper', { 'is-active': isActive }]">
            Consistency
            <el-icon class="header-icon">
              <info-filled />
            </el-icon>
          </div>
        </template>
        <div>
          Consistent with real life: in line with the process and logic of real
          life, and comply with languages and habits that the users are used to;
        </div>
        <div>
          Consistent within interface: all elements should be consistent, such
          as: design style, icons and texts, position of elements, etc.
        </div>
      </el-collapse-item>
      <el-collapse-item title="Feedback" name="2">
        <div>
          Operation feedback: enable the users to clearly perceive their
          operations by style updates and interactive effects;
        </div>
        <div>
          Visual feedback: reflect current state by updating or rearranging
          elements of the page.
        </div>
      </el-collapse-item>
      <el-collapse-item title="Efficiency" name="3">
        <div>
          Simplify the process: keep operating process simple and intuitive;
        </div>
        <div>
          Definite and clear: enunciate your intentions clearly so that the
          users can quickly understand and make decisions;
        </div>
        <div>
          Easy to identify: the interface should be straightforward, which helps
          the users to identify and frees them from memorizing and recalling.
        </div>
      </el-collapse-item>
      <el-collapse-item title="Controllability" name="4">
        <div>
          Decision making: giving advices about operations is acceptable, but do
          not make decisions for the users;
        </div>
        <div>
          Controlled consequences: users should be granted the freedom to
          operate, including canceling, aborting or terminating current
          operation.
        </div>
      </el-collapse-item>
    </el-collapse>
  </div>

  <div class="section-box">
    <el-tree
      :data="state.services"
      emptyText="There is no service"
      node-key="id"
      :expand-on-click-node="false"
      draggable
      :allow-drag="checkAllowDrag"
      :allow-drop="checkAllowDrop"
    >
      <template #default="{node, data}">
        <div class="flex flex-1 justify-between items-center p-2">
          <span>{{ data.name }}</span>
          <div v-if="node.level === 1">
            <el-button type="primary" @click="handleEditService(data)">Edit</el-button>
            <el-button type="danger" @click="handleDeleteService(data)">Delete</el-button>
            <el-button type="success" @click="handleAddRoute(data)">Add Route</el-button>
          </div>
          <div v-if="node.level === 2">
            <el-button type="primary" @click="handleEditRoute(node, data)">Edit</el-button>
            <el-button type="danger" @click="handleDeleteRoute(node, data)">Delete</el-button>
          </div>
        </div>
      </template>

    </el-tree>
  </div>

  <el-drawer v-model="state.serviceForm.showDrawer" direction="rtl">
    <template #header>
      <h4 v-if="state.serviceForm.action === 'update'">Update Rule #{{ state.serviceForm.data.id }}</h4>
      <h4 v-else>Create Rule</h4>
    </template>

    <template #default>

    </template>

    <template #footer>
      <div style="flex: auto">
        <el-button @click="state.serviceForm.showDrawer = false">Cancel</el-button>
        <el-button v-loading.fullscreen.lock="state.serviceForm.loading" type="primary" @click="handleSubmitService">
          {{ state.serviceForm.action == 'create' ? 'Create' : 'Update'}}
        </el-button>
      </div>
    </template>
  </el-drawer>

  <el-drawer v-model="state.routeForm.showDrawer" direction="rtl">
    <template #header>
      <h4 v-if="state.routeForm.action === 'update'">Update Route #{{ state.routeForm.data.id }}</h4>
      <h4 v-else>Create Route</h4>
    </template>

    <template #default>

    </template>

    <template #footer>
      <div style="flex: auto">
        <el-button @click="state.routeForm.showDrawer = false">Cancel</el-button>
        <el-button v-loading.fullscreen.lock="state.routeForm.loading" type="primary" @click="handleSubmitRoute">
          {{ state.routeForm.action == 'create' ? 'Create' : 'Update'}}
        </el-button>
      </div>
    </template>
  </el-drawer>
</template>

<script setup lang="ts">
import _ from 'lodash';
import { reactive, onMounted, ref } from 'vue';
import { useWorkspaceStore } from '@/stores/workspace';
import TopInfo from './TopInfo.vue';

const workspaceStore = useWorkspaceStore()

const state = reactive({
  serviceForm: {
    showDrawer: false,
    loading: false,
    data: {
    } as any,
    rules: {
    },
    action: '',   // create or update
  },
  routeForm: {
    showDrawer: false,
    loading: false,
    currentService: null as any,
    data: {
    } as any,
    rules: {
    },
    action: '',   // create or update
  },
  services: [
    {
      id: '1',
      name: 'Service 1',
      children: [
        {
          id: '1-1',
          name: 'Route 1-1',
        },
        {
          id: '1-2',
          name: 'Route 1-2',
        },
      ],
    },
    {
      id: '2',
      name: 'Service 2',
      children: [
        {
          id: '2-1',
          name: 'Route 2-1',
        },
        {
          id: '2-2',
          name: 'Route 2-2',
        },
      ],
    }
  ],
})

onMounted(() => {
});

const checkAllowDrag = (node: any) => {
  // only allow drag for route
  return node.level === 2
}
const checkAllowDrop = (draggingNode: any, dropNode: any, type: string) => {
  if (dropNode.level === 0) {
    return false
  }
  if (dropNode.level === 1 && type !== 'inner') {
    // route cannot drop to service
    return false
  }
  if (dropNode.level === 2 && type === 'inner') {
    // route cannot contain children
    return false
  }
  console.log(draggingNode, draggingNode.data, dropNode.level, dropNode.data, type)
  return true
}

const handleAddService = () => {
  state.serviceForm.action = 'create'
  state.serviceForm.showDrawer = true
}
const handleEditService = (service: any) => {
  state.serviceForm.action = 'update'
  state.serviceForm.showDrawer = true
  state.serviceForm.data = service
}
const handleDeleteService = (service: any) => {
}
const handleSubmitService = async () => {
}

const handleAddRoute = (service: any) => {
  state.routeForm.action = 'create'
  state.routeForm.showDrawer = true
  state.routeForm.currentService = service
}
const handleEditRoute = (node: any, route: any) => {
  state.routeForm.action = 'update'
  state.routeForm.showDrawer = true
  state.routeForm.data = route
}
const handleDeleteRoute = (node: any, route: any) => {
}
const handleSubmitRoute = async () => {
  
}
</script>

<style scoped>
  :deep(.el-tree-node__content) {
    height: initial;
  }
</style>