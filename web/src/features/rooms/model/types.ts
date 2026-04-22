export interface Room {
  id: string;
  name: string;
  building: string | null;
  capacity: number;
  notes: string | null;
  createdAt: string;
  updatedAt: string;
}

export type RoomsListParams = {
  page: number;
  pageSize: number;
  search: string;
};

export type RoomsListResponse = {
  items: Room[];
  total: number;
  page: number;
  pageSize: number;
};

export type RoomCreatePayload = {
  name: string;
  building?: string | null;
  capacity: number;
  notes?: string | null;
};

export type RoomUpdatePayload = Partial<RoomCreatePayload>;
