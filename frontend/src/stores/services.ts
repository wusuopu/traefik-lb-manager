import { defineStore } from 'pinia'
import API from '@/api/service';

export const useServiceStore = defineStore('services', {
  state: () => ({
    services: [] as Service[],
  }),
  actions: {
    async fetchIndexAsync (workspaceId: number) {
      const resp = await API.index(workspaceId);
      this.services = resp.data.Data;
    },
    async createAsync (workspaceId: number, service: Service) {
      const resp = await API.create(workspaceId, service);
    },
    async updateAsync (workspaceId: number, service: Service) {
      const resp = await API.update(workspaceId, service.ID!, service);
    },
    async deleteAsync (workspaceId: number, id: number) {
      await API.delete(workspaceId, id);
    },
  },
})