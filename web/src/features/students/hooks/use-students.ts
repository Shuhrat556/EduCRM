import {
  useMutation,
  useQuery,
  useQueryClient,
} from "@tanstack/react-query";

import { fetchGroupOptions } from "@/features/students/api/groups-api";
import { studentsApi } from "@/features/students/api/students-api";
import type {
  StudentCreatePayload,
  StudentUpdatePayload,
  StudentsListParams,
} from "@/features/students/model/types";
import { queryKeys } from "@/shared/api/query-keys";

export function useStudentsListQuery(
  params: StudentsListParams,
  options?: { enabled?: boolean },
) {
  return useQuery({
    queryKey: queryKeys.students.list({
      page: params.page,
      pageSize: params.pageSize,
      search: params.search,
      status: params.status,
      groupId: params.groupId ?? "",
    }),
    queryFn: () => studentsApi.list(params),
    placeholderData: (prev) => prev,
    enabled: options?.enabled ?? true,
  });
}

export function useStudentQuery(id: string | null) {
  return useQuery({
    queryKey: queryKeys.students.detail(id ?? ""),
    queryFn: () => studentsApi.get(id!),
    enabled: Boolean(id),
  });
}

export function useGroupOptionsQuery() {
  return useQuery({
    queryKey: queryKeys.groups.options(),
    queryFn: fetchGroupOptions,
    staleTime: 300_000,
  });
}

export function useCreateStudentMutation() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (payload: StudentCreatePayload) => studentsApi.create(payload),
    onSuccess: () => {
      void qc.invalidateQueries({ queryKey: queryKeys.students.all() });
    },
  });
}

export function useUpdateStudentMutation() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: ({
      id,
      payload,
    }: {
      id: string;
      payload: StudentUpdatePayload;
    }) => studentsApi.update(id, payload),
    onSuccess: (_, { id }) => {
      void qc.invalidateQueries({ queryKey: queryKeys.students.all() });
      void qc.invalidateQueries({
        queryKey: queryKeys.students.detail(id),
      });
    },
  });
}

export function useDeleteStudentMutation() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => studentsApi.remove(id),
    onSuccess: () => {
      void qc.invalidateQueries({ queryKey: queryKeys.students.all() });
    },
  });
}
