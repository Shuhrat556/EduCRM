import { z } from "zod";

function timeToMinutes(t: string): number | null {
  const m = /^(\d{1,2}):(\d{2})$/.exec(t.trim());
  if (!m) return null;
  const h = Number(m[1]);
  const min = Number(m[2]);
  if (h < 0 || h > 23 || min < 0 || min > 59) return null;
  return h * 60 + min;
}

const timeSchema = z
  .string()
  .trim()
  .refine((s) => timeToMinutes(s) !== null, "Use a valid time (HH:mm)");

const weekdaySchema = z.enum([
  "monday",
  "tuesday",
  "wednesday",
  "thursday",
  "friday",
  "saturday",
  "sunday",
]);

export const lessonFormSchema = z
  .object({
    weekday: weekdaySchema,
    startTime: timeSchema,
    endTime: timeSchema,
    roomId: z.string().min(1, "Choose a room"),
    teacherId: z.string().min(1, "Choose a teacher"),
    groupId: z.string().min(1, "Choose a group"),
    title: z
      .string()
      .trim()
      .max(120)
      .optional()
      .transform((s) => (s === "" ? undefined : s)),
  })
  .refine((data) => {
    const a = timeToMinutes(data.startTime);
    const b = timeToMinutes(data.endTime);
    if (a === null || b === null) return true;
    return a < b;
  }, {
    message: "End time must be after start time",
    path: ["endTime"],
  });

export type LessonFormValues = z.infer<typeof lessonFormSchema>;
