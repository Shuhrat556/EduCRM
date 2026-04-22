import { CalendarDays, ClipboardCheck } from "lucide-react";
import { useMemo, useState } from "react";
import { useNavigate } from "react-router-dom";

import { isoDateToday } from "@/features/attendance/lib/iso-date";
import { compareTimeStrings } from "@/features/schedule/lib/time-compare";
import { weekdayLabel, WEEKDAY_ORDER } from "@/features/schedule/lib/weekdays";
import { useScheduleLessonsQuery } from "@/features/schedule/hooks/use-schedule";
import { Button } from "@/shared/ui/components/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/shared/ui/components/card";
import { Input } from "@/shared/ui/components/input";
import { Label } from "@/shared/ui/components/label";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/shared/ui/components/select";
import { Skeleton } from "@/shared/ui/components/skeleton";
import { usePortalBase } from "@/shared/lib/hooks/usePortalBase";

export function AttendanceHomePage() {
  const base = usePortalBase();
  const navigate = useNavigate();
  const { data: lessons, isLoading } = useScheduleLessonsQuery();
  const [lessonId, setLessonId] = useState("");
  const [sessionDate, setSessionDate] = useState(isoDateToday());

  const sorted = useMemo(() => {
    const list = [...(lessons ?? [])];
    const order = new Map(WEEKDAY_ORDER.map((d, i) => [d, i]));
    list.sort((a, b) => {
      const da = order.get(a.weekday) ?? 99;
      const db = order.get(b.weekday) ?? 99;
      if (da !== db) return da - db;
      return compareTimeStrings(a.startTime, b.startTime);
    });
    return list;
  }, [lessons]);

  function openRoster() {
    if (!lessonId) return;
    navigate(`${base}/attendance/lesson/${lessonId}?date=${sessionDate}`);
  }

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-semibold tracking-tight md:text-3xl">
          Attendance & grades
        </h1>
        <p className="text-sm text-muted-foreground">
          Pick a scheduled lesson and session date to record attendance, notes,
          lesson grades, and optional weekly ratings for the whole group.
        </p>
      </div>

      <Card className="max-w-xl border-border/80">
        <CardHeader>
          <CardTitle className="flex items-center gap-2 text-base">
            <ClipboardCheck className="h-4 w-4" />
            Open lesson roster
          </CardTitle>
          <CardDescription>
            Roster lists every student in the lesson&apos;s group. Save updates
            in one batch for the selected date.
          </CardDescription>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="space-y-2">
            <Label htmlFor="lesson">Lesson slot</Label>
            {isLoading ? (
              <Skeleton className="h-9 w-full" />
            ) : (
              <Select value={lessonId} onValueChange={setLessonId}>
                <SelectTrigger id="lesson">
                  <SelectValue placeholder="Choose a lesson from the schedule" />
                </SelectTrigger>
                <SelectContent>
                  {sorted.map((l) => (
                    <SelectItem key={l.id} value={l.id}>
                      {weekdayLabel(l.weekday)} · {l.startTime}–{l.endTime} ·{" "}
                      {l.group.name}
                      {l.title ? ` — ${l.title}` : ""}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            )}
          </div>
          <div className="space-y-2">
            <Label htmlFor="session-date">Session date</Label>
            <div className="flex items-center gap-2">
              <CalendarDays className="h-4 w-4 shrink-0 text-muted-foreground" />
              <Input
                id="session-date"
                type="date"
                value={sessionDate}
                onChange={(e) => setSessionDate(e.target.value)}
                className="max-w-xs"
              />
            </div>
            <p className="text-xs text-muted-foreground">
              Same lesson can have different attendance per calendar date (e.g.
              each Monday).
            </p>
          </div>
          <Button
            type="button"
            className="w-full sm:w-auto"
            disabled={!lessonId || isLoading}
            onClick={openRoster}
          >
            Open roster
          </Button>
        </CardContent>
      </Card>
    </div>
  );
}
