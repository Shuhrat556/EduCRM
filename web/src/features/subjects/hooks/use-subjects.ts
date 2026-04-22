import {
  useMutation,
  useQuery,
  useQueryClient,
} from "@tanstack/react-query";

import { subjectsApi } from "@/features/subjects/api/subjects-api";
import type {
  SubjectCreatePayload,
  SubjectUpdatePayload,
  SubjectsListParams,
} from "@/features/subjects/model/types";
import { queryKeys } from "@/shared/api/query-keys";

export function useSubjectsListQuery(params: SubjectsListParams) {
  return useQuery({
    queryKey: queryKeys.subjects.list({
      page: params.page,
      pageSize: params.pageSize,
      search: params.search,
    }),
    queryFn: () => subjectsApi.list(params),
    placeholderData: (prev) => prev,
  });
}

export function useSubjectQuery(id: string | null) {
  return useQuery({
    queryKey: queryKeys.subjects.detail(id ?? ""),
    queryFn: () => subjectsApi.get(id!),
    enabled: Boolean(id),
  });
}

export function useCreateSubjectMutation() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (payload: SubjectCreatePayload) => subjectsApi.create(payload),
    onSuccess: () => {
      void qc.invalidateQueries({ queryKey: queryKeys.subjects.all() });
    },
  });
}

export function useUpdateSubjectMutation() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: ({
      id,
      payload,
    }: {
      id: string;
      payload: SubjectUpdatePayload;
    }) => subjectsApi.update(id, payload),
    onSuccess: (_, { id }) => {
      void qc.invalidateQueries({ queryKey: queryKeys.subjects.all() });
      void qc.invalidateQueries({ queryKey: queryKeys.subjects.detail(id) });
    },
  });
}

export function useDeleteSubjectMutation() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => subjectsApi.remove(id),
    onSuccess: () => {
      void qc.invalidateQueries({ queryKey: queryKeys.subjects.all() });
    },
  });
}
