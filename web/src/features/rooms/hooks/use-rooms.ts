import {
  useMutation,
  useQuery,
  useQueryClient,
} from "@tanstack/react-query";

import { roomsApi } from "@/features/rooms/api/rooms-api";
import type {
  RoomCreatePayload,
  RoomUpdatePayload,
  RoomsListParams,
} from "@/features/rooms/model/types";
import { queryKeys } from "@/shared/api/query-keys";

export function useRoomsListQuery(params: RoomsListParams) {
  return useQuery({
    queryKey: queryKeys.rooms.list({
      page: params.page,
      pageSize: params.pageSize,
      search: params.search,
    }),
    queryFn: () => roomsApi.list(params),
    placeholderData: (prev) => prev,
  });
}

export function useRoomQuery(id: string | null) {
  return useQuery({
    queryKey: queryKeys.rooms.detail(id ?? ""),
    queryFn: () => roomsApi.get(id!),
    enabled: Boolean(id),
  });
}

export function useCreateRoomMutation() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (payload: RoomCreatePayload) => roomsApi.create(payload),
    onSuccess: () => {
      void qc.invalidateQueries({ queryKey: queryKeys.rooms.all() });
    },
  });
}

export function useUpdateRoomMutation() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: ({
      id,
      payload,
    }: {
      id: string;
      payload: RoomUpdatePayload;
    }) => roomsApi.update(id, payload),
    onSuccess: (_, { id }) => {
      void qc.invalidateQueries({ queryKey: queryKeys.rooms.all() });
      void qc.invalidateQueries({ queryKey: queryKeys.rooms.detail(id) });
    },
  });
}

export function useDeleteRoomMutation() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => roomsApi.remove(id),
    onSuccess: () => {
      void qc.invalidateQueries({ queryKey: queryKeys.rooms.all() });
    },
  });
}
