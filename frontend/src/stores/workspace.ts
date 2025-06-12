import { defineStore } from 'pinia'
import API from '@/api/workspace';

export const useWorkspaceStore = defineStore('workspaces', {
  state: () => ({
    workspaces: [] as Workspace[],
    currentWorkspace: null as Workspace | null,
  }),
  actions: {
    async fetchIndexAsync () {
      const resp = await API.index();
      this.workspaces = resp.data.Data;
    },
    async fetchShowAsync (id: number) {
      const resp = await API.show(id);
      this.currentWorkspace = resp.data.Data;
    },
    async createAsync (workspace: Workspace) {
      const resp = await API.create(workspace);
    },
    async updateAsync (workspace: Workspace) {
      await API.update(workspace.ID!, workspace);
    },
    async deleteAsync (id: number) {
      await API.delete(id);
      this.workspaces = this.workspaces.filter((workspace) => workspace.ID !== id);
    },
    setCurrentWorkspace (workspace: Workspace|null) {
      this.currentWorkspace = workspace;
    },
  },
})