import { defineStore } from 'pinia'

export const useWorkspaceStore = defineStore('workspaces', {
  state: () => ({
    workspaces: [] as Workspace[],
    currentWorkspace: null as Workspace | null,
  }),
  actions: {
    
  },
})