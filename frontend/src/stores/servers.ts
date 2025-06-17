import { defineStore } from 'pinia'
import API from '@/api/server';

export const useServerStore = defineStore('servers', {
  state: () => ({
    servers: [] as Server[],
  }),
  actions: {
    async fetchIndexAsync (workspaceId: number) {
      const resp = await API.index(workspaceId);
      this.servers = resp.data.Data;
    },
    async createAsync (workspaceId: number, server: Server) {
      await API.create(workspaceId, server);
    },
    async updateAsync (workspaceId: number, server: Server) {
      await API.update(workspaceId, server.ID!, server);
    },
    async deleteAsync (workspaceId: number, id: number) {
      await API.delete(workspaceId, id);
    },
  },
})
