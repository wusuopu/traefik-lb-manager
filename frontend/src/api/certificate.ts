import http from './_http';

export default {
  index (workspaceId: number) {
    return http.get(`/workspaces/${workspaceId}/certificates/`);
  },
  create (workspaceId: number, data: any) {
    return http.post(`/workspaces/${workspaceId}/certificates/`, data);
  },
  delete (workspaceId: number, id: number) {
    return http.delete(`/workspaces/${workspaceId}/certificates/${id}`);
  },
}