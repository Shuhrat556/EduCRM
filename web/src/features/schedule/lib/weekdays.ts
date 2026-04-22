export const WEEKDAYS = [
  { value: "monday", label: "Monday" },
  { value: "tuesday", label: "Tuesday" },
  { value: "wednesday", label: "Wednesday" },
  { value: "thursday", label: "Thursday" },
  { value: "friday", label: "Friday" },
  { value: "saturday", label: "Saturday" },
  { value: "sunday", label: "Sunday" },
] as const;

export type Weekday = (typeof WEEKDAYS)[number]["value"];

export const WEEKDAY_ORDER: Weekday[] = WEEKDAYS.map((w) => w.value);

export function weekdayLabel(w: Weekday): string {
  return WEEKDAYS.find((d) => d.value === w)?.label ?? w;
}
