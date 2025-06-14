import http from './_http';

export default {
  index (workspaceId: number) {
    return http.get(`/workspaces/${workspaceId}/certificates/`);
  },
  create (workspaceId: number, data: any) {
    return http.post(`/workspaces/${workspaceId}/certificates/`, data);
  },
  update (workspaceId: number, id: number, data: any) {
    return http.put(`/workspaces/${workspaceId}/certificates/${id}`, data);
  },
  delete (workspaceId: number, id: number) {
    return http.delete(`/workspaces/${workspaceId}/certificates/${id}`);
  },
  renew (workspaceId: number, id: number) {
    return http.put(`/workspaces/${workspaceId}/certificates/${id}/renew`);
  },
}