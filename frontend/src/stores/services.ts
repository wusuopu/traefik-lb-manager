import { defineStore } from 'pinia'
import API from '@/api/service';

export const useServiceStore = defineStore('services', {
  state: () => ({
    services: [] as Service[],
    externals: [] as ExternalService[],
  }),
  actions: {
    async fetchIndexAsync (workspaceId: number) {
      const resp = await API.index(workspaceId);
      this.services = resp.data.Data;
    },
    async createAsync (workspaceId: number, service: Service) {
      await API.create(workspaceId, service);
    },
    async updateAsync (workspaceId: number, service: Service) {
      await API.update(workspaceId, service.ID!, service);
    },
    async deleteAsync (workspaceId: number, id: number) {
      await API.delete(workspaceId, id);
    },
    async fetchExternalIndexAsync (workspaceId: number) {
      try {
        const resp = await API.externalIndex(workspaceId);
        this.externals = resp.data.Data || [];
      } catch (error) {
        this.externals = [];
        throw error
      }
    }
  },
})