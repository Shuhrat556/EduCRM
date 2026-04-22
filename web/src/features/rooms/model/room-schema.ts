import { z } from "zod";

import type { RoomCreatePayload } from "@/features/rooms/model/types";

export const roomFormSchema = z.object({
  name: z.string().trim().min(1, "Room name is required").max(120),
  building: z.string().trim().max(80),
  capacity: z.coerce
    .number({ invalid_type_error: "Enter capacity" })
    .int("Whole seats only")
    .min(1, "At least 1 seat")
    .max(9999, "Too large"),
  notes: z.string().trim().max(500),
});

export type RoomFormValues = z.infer<typeof roomFormSchema>;

export function roomFormToPayload(values: RoomFormValues): RoomCreatePayload {
  return {
    name: values.name,
    building:
      values.building.trim() === "" ? null : values.building.trim(),
    capacity: values.capacity,
    notes: values.notes.trim() === "" ? null : values.notes.trim(),
  };
}
