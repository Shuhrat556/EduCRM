import {
  useMutation,
  useQuery,
  useQueryClient,
} from "@tanstack/react-query";

import { fetchSubjectOptions } from "@/features/subjects/api/subjects-api";
import { teachersApi } from "@/features/teachers/api/teachers-api";
import type {
  TeacherCreatePayload,
  TeacherUpdatePayload,
  TeachersListParams,
} from "@/features/teachers/model/types";
import { queryKeys } from "@/shared/api/query-keys";

export function useTeachersListQuery(params: TeachersListParams) {
  return useQuery({
    queryKey: queryKeys.teachers.list({
      page: params.page,
      pageSize: params.pageSize,
      search: params.search,
      status: params.status,
    }),
    queryFn: () => teachersApi.list(params),
    placeholderData: (prev) => prev,
  });
}

export function useTeacherQuery(id: string | null) {
  return useQuery({
    queryKey: queryKeys.teachers.detail(id ?? ""),
    queryFn: () => teachersApi.get(id!),
    enabled: Boolean(id),
  });
}

export function useSubjectOptionsQuery() {
  return useQuery({
    queryKey: queryKeys.subjects.options(),
    queryFn: fetchSubjectOptions,
    staleTime: 300_000,
  });
}

export function useCreateTeacherMutation() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (payload: TeacherCreatePayload) => teachersApi.create(payload),
    onSuccess: () => {
      void qc.invalidateQueries({ queryKey: queryKeys.teachers.all() });
    },
  });
}

export function useUpdateTeacherMutation() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: ({
      id,
      payload,
    }: {
      id: string;
      payload: TeacherUpdatePayload;
    }) => teachersApi.update(id, payload),
    onSuccess: (_, { id }) => {
      void qc.invalidateQueries({ queryKey: queryKeys.teachers.all() });
      void qc.invalidateQueries({ queryKey: queryKeys.teachers.detail(id) });
    },
  });
}

export function useDeleteTeacherMutation() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => teachersApi.remove(id),
    onSuccess: () => {
      void qc.invalidateQueries({ queryKey: queryKeys.teachers.all() });
    },
  });
}

export function useUploadTeacherPhotoMutation() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: ({ id, file }: { id: string; file: File }) =>
      teachersApi.uploadPhoto(id, file),
    onSuccess: (_, { id }) => {
      void qc.invalidateQueries({ queryKey: queryKeys.teachers.all() });
      void qc.invalidateQueries({ queryKey: queryKeys.teachers.detail(id) });
    },
  });
}
