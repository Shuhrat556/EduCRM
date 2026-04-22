import {
  useMutation,
  useQuery,
  useQueryClient,
} from "@tanstack/react-query";

import { groupsApi } from "@/features/groups/api/groups-api";
import type {
  GroupCreatePayload,
  GroupUpdatePayload,
  GroupsListParams,
} from "@/features/groups/model/types";
import { studentsApi } from "@/features/students/api/students-api";
import { queryKeys } from "@/shared/api/query-keys";

function useStudentsDemo() {
  return import.meta.env.VITE_STUDENTS_DEMO === "true";
}

export function useGroupsListQuery(params: GroupsListParams) {
  return useQuery({
    queryKey: queryKeys.groups.list({
      page: params.page,
      pageSize: params.pageSize,
      search: params.search,
      status: params.status,
    }),
    queryFn: () => groupsApi.list(params),
    placeholderData: (prev) => prev,
  });
}

export function useGroupQuery(id: string | null) {
  return useQuery({
    queryKey: queryKeys.groups.detail(id ?? ""),
    queryFn: () => groupsApi.get(id!),
    enabled: Boolean(id),
  });
}

export function useCreateGroupMutation() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (payload: GroupCreatePayload) => groupsApi.create(payload),
    onSuccess: () => {
      void qc.invalidateQueries({ queryKey: queryKeys.groups.all() });
    },
  });
}

export function useUpdateGroupMutation() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: ({
      id,
      payload,
    }: {
      id: string;
      payload: GroupUpdatePayload;
    }) => groupsApi.update(id, payload),
    onSuccess: (_, { id }) => {
      void qc.invalidateQueries({ queryKey: queryKeys.groups.all() });
      void qc.invalidateQueries({ queryKey: queryKeys.groups.detail(id) });
    },
  });
}

export function useDeleteGroupMutation() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: async (id: string) => {
      if (useStudentsDemo()) {
        await studentsApi.unassignAllByGroup(id);
      }
      await groupsApi.remove(id);
    },
    onSuccess: () => {
      void qc.invalidateQueries({ queryKey: queryKeys.groups.all() });
      void qc.invalidateQueries({ queryKey: queryKeys.students.all() });
    },
  });
}
