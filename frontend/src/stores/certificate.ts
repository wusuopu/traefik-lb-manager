import { defineStore } from 'pinia'
import API from '@/api/certificate';

export const useCertificateStore = defineStore('certificates', {
  state: () => ({
    certificates: [] as Certificate[],
  }),
  actions: {
    async fetchIndexAsync (workspaceId: number) {
      const resp = await API.index(workspaceId);
      this.certificates = resp.data.Data;
    },
    async createAsync (workspaceId: number, certificate: Certificate) {
      const resp = await API.create(workspaceId, certificate);
    },
    async updateAsync (workspaceId: number, certificate: Certificate) {
      const resp = await API.update(workspaceId, certificate.ID!, certificate);
    },
    async deleteAsync (workspaceId: number, id: number) {
      await API.delete(workspaceId, id);
    },
    async renewAsync (workspaceId: number, id: number) {
      await API.renew(workspaceId, id);
    },
  },
})