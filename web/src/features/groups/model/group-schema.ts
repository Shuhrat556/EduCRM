import { z } from "zod";

import type {
  GroupCreatePayload,
  GroupStatus,
} from "@/features/groups/model/types";

const statusEnum = z.enum([
  "draft",
  "active",
  "paused",
  "completed",
  "archived",
]) satisfies z.ZodType<GroupStatus>;

/** Sentinel value for optional entity selects in the group form. */
export const GROUP_FORM_NONE = "__none__" as const;

export const groupFormSchema = z.object({
  name: z.string().trim().min(1, "Name is required").max(160),
  teacherId: z.string(),
  subjectId: z.string(),
  roomId: z.string(),
  monthlyFee: z.coerce
    .number({ invalid_type_error: "Enter a fee" })
    .int("Use whole numbers")
    .min(0, "Fee cannot be negative"),
  startDate: z.string().min(1, "Start date is required"),
  endDate: z.string().optional(),
  status: statusEnum,
});

export type GroupFormValues = z.infer<typeof groupFormSchema>;

export function groupFormToCreatePayload(
  values: GroupFormValues,
): GroupCreatePayload {
  return {
    name: values.name,
    teacherId:
      values.teacherId === GROUP_FORM_NONE ? null : values.teacherId,
    subjectId:
      values.subjectId === GROUP_FORM_NONE ? null : values.subjectId,
    roomId: values.roomId === GROUP_FORM_NONE ? null : values.roomId,
    monthlyFee: values.monthlyFee,
    startDate: values.startDate,
    endDate:
      values.endDate && values.endDate.trim() !== ""
        ? values.endDate
        : null,
    status: values.status,
  };
}
