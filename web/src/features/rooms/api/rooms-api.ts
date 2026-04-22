import { apiClient } from "@/shared/api/client";
import { mockRoomsStore } from "@/features/rooms/lib/mock-rooms-store";
import type {
  Room,
  RoomCreatePayload,
  RoomUpdatePayload,
  RoomsListParams,
  RoomsListResponse,
} from "@/features/rooms/model/types";

function useRoomsDemo() {
  return import.meta.env.VITE_ROOMS_DEMO === "true";
}

/**
 * REST (adjust to your backend):
 * - GET /rooms?page&pageSize&search
 * - GET /rooms/:id
 * - POST /rooms
 * - PUT /rooms/:id
 * - DELETE /rooms/:id
 */
export const roomsApi = {
  list: async (params: RoomsListParams): Promise<RoomsListResponse> => {
    if (useRoomsDemo()) return mockRoomsStore.list(params);
    const { data } = await apiClient.get<RoomsListResponse>("/rooms", {
      params: {
        page: params.page,
        pageSize: params.pageSize,
        search: params.search || undefined,
      },
    });
    return data;
  },

  get: async (id: string): Promise<Room> => {
    if (useRoomsDemo()) {
      const r = await mockRoomsStore.getById(id);
      if (!r) throw new Error("Not found");
      return r;
    }
    const { data } = await apiClient.get<Room>(`/rooms/${id}`);
    return data;
  },

  create: async (payload: RoomCreatePayload): Promise<Room> => {
    if (useRoomsDemo()) return mockRoomsStore.create(payload);
    const { data } = await apiClient.post<Room>("/rooms", payload);
    return data;
  },

  update: async (
    id: string,
    payload: RoomUpdatePayload,
  ): Promise<Room> => {
    if (useRoomsDemo()) return mockRoomsStore.update(id, payload);
    const { data } = await apiClient.put<Room>(`/rooms/${id}`, payload);
    return data;
  },

  remove: async (id: string): Promise<void> => {
    if (useRoomsDemo()) return mockRoomsStore.remove(id);
    await apiClient.delete(`/rooms/${id}`);
  },
};
