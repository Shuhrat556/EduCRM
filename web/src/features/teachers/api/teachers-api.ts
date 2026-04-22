import { apiClient } from "@/shared/api/client";
import { mockTeachersStore } from "@/features/teachers/lib/mock-teachers-store";
import type {
  Teacher,
  TeacherCreatePayload,
  TeacherUpdatePayload,
  TeachersListParams,
  TeachersListResponse,
} from "@/features/teachers/model/types";

function useTeachersDemo() {
  return import.meta.env.VITE_TEACHERS_DEMO === "true";
}

/**
 * REST (adjust to your backend):
 * - GET /teachers?page&pageSize&search&status
 * - GET /teachers/:id
 * - POST /teachers
 * - PUT /teachers/:id
 * - DELETE /teachers/:id
 * - POST /teachers/:id/photo (multipart field `file`)
 */
export const teachersApi = {
  list: async (params: TeachersListParams): Promise<TeachersListResponse> => {
    if (useTeachersDemo()) return mockTeachersStore.list(params);
    const { data } = await apiClient.get<TeachersListResponse>("/teachers", {
      params: {
        page: params.page,
        pageSize: params.pageSize,
        search: params.search || undefined,
        status: params.status === "all" ? undefined : params.status,
      },
    });
    return data;
  },

  get: async (id: string): Promise<Teacher> => {
    if (useTeachersDemo()) {
      const t = await mockTeachersStore.getById(id);
      if (!t) throw new Error("Not found");
      return t;
    }
    const { data } = await apiClient.get<Teacher>(`/teachers/${id}`);
    return data;
  },

  create: async (payload: TeacherCreatePayload): Promise<Teacher> => {
    if (useTeachersDemo()) return mockTeachersStore.create(payload);
    const { data } = await apiClient.post<Teacher>("/teachers", payload);
    return data;
  },

  update: async (
    id: string,
    payload: TeacherUpdatePayload,
  ): Promise<Teacher> => {
    if (useTeachersDemo()) return mockTeachersStore.update(id, payload);
    const { data } = await apiClient.put<Teacher>(`/teachers/${id}`, payload);
    return data;
  },

  remove: async (id: string): Promise<void> => {
    if (useTeachersDemo()) return mockTeachersStore.remove(id);
    await apiClient.delete(`/teachers/${id}`);
  },

  uploadPhoto: async (
    id: string,
    file: File,
  ): Promise<{ photoUrl: string }> => {
    if (useTeachersDemo()) return mockTeachersStore.uploadPhoto(id, file);
    const body = new FormData();
    body.append("file", file);
    const { data } = await apiClient.post<{ photoUrl: string }>(
      `/teachers/${id}/photo`,
      body,
      { headers: { "Content-Type": "multipart/form-data" } },
    );
    return data;
  },
};
