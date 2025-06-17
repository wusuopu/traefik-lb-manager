import http from './_http';

export default {
  index (workspaceId: number) {
    return http.get(`/workspaces/${workspaceId}/services/`);
  },
  create (workspaceId: number, data: any) {
    return http.post(`/workspaces/${workspaceId}/services/`, data);
  },
  update (workspaceId: number, id: number, data: any) {
    return http.put(`/workspaces/${workspaceId}/services/${id}`, data);
  },
  delete (workspaceId: number, id: number) {
    return http.delete(`/workspaces/${workspaceId}/services/${id}`);
  },
  externalIndex (workspaceId: number) {
    return http.get(`/workspaces/${workspaceId}/services/external`)
  }
}