import { apiClient } from "@/shared/api/client";
import { FALLBACK_SUBJECT_OPTIONS } from "@/features/subjects/lib/fallback-subject-options";
import { mockSubjectsStore } from "@/features/subjects/lib/mock-subjects-store";
import type {
  Subject,
  SubjectCreatePayload,
  SubjectOption,
  SubjectUpdatePayload,
  SubjectsListParams,
  SubjectsListResponse,
} from "@/features/subjects/model/types";

function useSubjectsDemo() {
  return import.meta.env.VITE_SUBJECTS_DEMO === "true";
}

function useTeachersDemoOnly() {
  return (
    import.meta.env.VITE_TEACHERS_DEMO === "true" && !useSubjectsDemo()
  );
}

/**
 * REST (adjust to your backend):
 * - GET /subjects?page&pageSize&search
 * - GET /subjects/:id
 * - POST /subjects
 * - PUT /subjects/:id
 * - DELETE /subjects/:id
 * - GET /subjects/options (id + name rows for pickers)
 */
export const subjectsApi = {
  list: async (
    params: SubjectsListParams,
  ): Promise<SubjectsListResponse> => {
    if (useSubjectsDemo()) return mockSubjectsStore.list(params);
    const { data } = await apiClient.get<SubjectsListResponse>("/subjects", {
      params: {
        page: params.page,
        pageSize: params.pageSize,
        search: params.search || undefined,
      },
    });
    return data;
  },

  get: async (id: string): Promise<Subject> => {
    if (useSubjectsDemo()) {
      const s = await mockSubjectsStore.getById(id);
      if (!s) throw new Error("Not found");
      return s;
    }
    const { data } = await apiClient.get<Subject>(`/subjects/${id}`);
    return data;
  },

  create: async (payload: SubjectCreatePayload): Promise<Subject> => {
    if (useSubjectsDemo()) return mockSubjectsStore.create(payload);
    const { data } = await apiClient.post<Subject>("/subjects", payload);
    return data;
  },

  update: async (
    id: string,
    payload: SubjectUpdatePayload,
  ): Promise<Subject> => {
    if (useSubjectsDemo()) return mockSubjectsStore.update(id, payload);
    const { data } = await apiClient.put<Subject>(`/subjects/${id}`, payload);
    return data;
  },

  remove: async (id: string): Promise<void> => {
    if (useSubjectsDemo()) return mockSubjectsStore.remove(id);
    await apiClient.delete(`/subjects/${id}`);
  },
};

/** Shared by teacher assignment UI (`useSubjectOptionsQuery`). */
export async function fetchSubjectOptions(): Promise<SubjectOption[]> {
  if (useSubjectsDemo()) return mockSubjectsStore.getOptions();
  if (useTeachersDemoOnly()) return FALLBACK_SUBJECT_OPTIONS;
  const { data } = await apiClient.get<SubjectOption[]>("/subjects/options");
  return data;
}
