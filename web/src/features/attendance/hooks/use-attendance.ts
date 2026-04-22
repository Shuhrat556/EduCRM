import {
  useMutation,
  useQuery,
  useQueryClient,
} from "@tanstack/react-query";

import { attendanceApi } from "@/features/attendance/api/attendance-api";
import type { AttendanceBatchPayload } from "@/features/attendance/model/types";
import { queryKeys } from "@/shared/api/query-keys";

export function useAttendanceSessionQuery(
  lessonId: string | null,
  sessionDate: string | null,
) {
  return useQuery({
    queryKey: queryKeys.attendance.session(lessonId ?? "", sessionDate ?? ""),
    queryFn: () => attendanceApi.getSession(lessonId!, sessionDate!),
    enabled: Boolean(lessonId && sessionDate),
  });
}

export function useSaveAttendanceBatchMutation() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (payload: AttendanceBatchPayload) =>
      attendanceApi.saveBatch(payload),
    onSuccess: (_, payload) => {
      void qc.invalidateQueries({
        queryKey: queryKeys.attendance.session(
          payload.lessonId,
          payload.sessionDate,
        ),
      });
    },
  });
}
