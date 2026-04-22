import {
  useMutation,
  useQuery,
  useQueryClient,
} from "@tanstack/react-query";

import { scheduleApi } from "@/features/schedule/api/schedule-api";
import type {
  LessonCreatePayload,
  LessonUpdatePayload,
} from "@/features/schedule/model/types";
import { queryKeys } from "@/shared/api/query-keys";

export function useScheduleLessonsQuery() {
  return useQuery({
    queryKey: queryKeys.schedule.lessons(),
    queryFn: () => scheduleApi.listWeek(),
  });
}

export function useScheduleLessonQuery(id: string | null) {
  return useQuery({
    queryKey: queryKeys.schedule.lesson(id ?? ""),
    queryFn: () => scheduleApi.get(id!),
    enabled: Boolean(id),
  });
}

export function useCreateLessonMutation() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (payload: LessonCreatePayload) => scheduleApi.create(payload),
    onSuccess: () => {
      void qc.invalidateQueries({ queryKey: queryKeys.schedule.all() });
    },
  });
}

export function useUpdateLessonMutation() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: ({
      id,
      payload,
    }: {
      id: string;
      payload: LessonUpdatePayload;
    }) => scheduleApi.update(id, payload),
    onSuccess: (_, { id }) => {
      void qc.invalidateQueries({ queryKey: queryKeys.schedule.all() });
      void qc.invalidateQueries({ queryKey: queryKeys.schedule.lesson(id) });
    },
  });
}

export function useDeleteLessonMutation() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => scheduleApi.remove(id),
    onSuccess: () => {
      void qc.invalidateQueries({ queryKey: queryKeys.schedule.all() });
    },
  });
}
