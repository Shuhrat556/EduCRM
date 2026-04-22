import { apiClient } from "@/shared/api/client";
import { mockStudentsStore } from "@/features/students/lib/mock-students-store";
import type {
  Student,
  StudentCreatePayload,
  StudentUpdatePayload,
  StudentsListParams,
  StudentsListResponse,
} from "@/features/students/model/types";

function useStudentsDemo() {
  return import.meta.env.VITE_STUDENTS_DEMO === "true";
}

/**
 * REST shape (adjustable to your backend):
 * - `GET /students?page=&pageSize=&search=&status=`
 * - `GET /students/:id`
 * - `POST /students`
 * - `PUT /students/:id`
 * - `DELETE /students/:id`
 */
export const studentsApi = {
  list: async (params: StudentsListParams): Promise<StudentsListResponse> => {
    if (useStudentsDemo()) return mockStudentsStore.list(params);
    const { data } = await apiClient.get<StudentsListResponse>("/students", {
      params: {
        page: params.page,
        pageSize: params.pageSize,
        search: params.search || undefined,
        status: params.status === "all" ? undefined : params.status,
        groupId: params.groupId || undefined,
      },
    });
    return data;
  },

  get: async (id: string): Promise<Student> => {
    if (useStudentsDemo()) {
      const s = await mockStudentsStore.getById(id);
      if (!s) throw new Error("Not found");
      return s;
    }
    const { data } = await apiClient.get<Student>(`/students/${id}`);
    return data;
  },

  create: async (payload: StudentCreatePayload): Promise<Student> => {
    if (useStudentsDemo()) return mockStudentsStore.create(payload);
    const { data } = await apiClient.post<Student>("/students", payload);
    return data;
  },

  update: async (
    id: string,
    payload: StudentUpdatePayload,
  ): Promise<Student> => {
    if (useStudentsDemo()) return mockStudentsStore.update(id, payload);
    const { data } = await apiClient.put<Student>(`/students/${id}`, payload);
    return data;
  },

  remove: async (id: string): Promise<void> => {
    if (useStudentsDemo()) return mockStudentsStore.remove(id);
    await apiClient.delete(`/students/${id}`);
  },

  /** Demo/local: clear `groupId` for every student in this group (MVP single-group rule). */
  unassignAllByGroup: async (groupId: string): Promise<void> => {
    if (useStudentsDemo()) return mockStudentsStore.unassignAllByGroup(groupId);
  },
};
