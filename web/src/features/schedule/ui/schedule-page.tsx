import {
  Building2,
  CalendarDays,
  ClipboardCheck,
  Clock,
  MoreHorizontal,
  Plus,
  UserRound,
  Users,
} from "lucide-react";
import { Link } from "react-router-dom";

import { isoDateToday } from "@/features/attendance/lib/iso-date";
import { useMemo, useState } from "react";

import type { ScheduleLesson } from "@/features/schedule/model/types";
import type { Weekday } from "@/features/schedule/lib/weekdays";
import { WEEKDAY_ORDER, WEEKDAYS, weekdayLabel } from "@/features/schedule/lib/weekdays";
import { compareTimeStrings } from "@/features/schedule/lib/time-compare";
import {
  useDeleteLessonMutation,
  useScheduleLessonsQuery,
} from "@/features/schedule/hooks/use-schedule";
import { DeleteLessonDialog } from "@/features/schedule/ui/delete-lesson-dialog";
import { LessonCreateDialog } from "@/features/schedule/ui/lesson-create-dialog";
import { LessonEditDialog } from "@/features/schedule/ui/lesson-edit-dialog";
import { AppEmpty } from "@/shared/ui/feedback/app-empty";
import { Button } from "@/shared/ui/components/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/shared/ui/components/card";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/shared/ui/components/dropdown-menu";
import { ScrollArea } from "@/shared/ui/components/scroll-area";
import { Skeleton } from "@/shared/ui/components/skeleton";
import { cn } from "@/shared/lib/cn";
import { usePortalBase } from "@/shared/lib/hooks/usePortalBase";

function lessonsByWeekday(lessons: ScheduleLesson[]): Record<Weekday, ScheduleLesson[]> {
  const init = {} as Record<Weekday, ScheduleLesson[]>;
  for (const d of WEEKDAY_ORDER) init[d] = [];
  for (const l of lessons) {
    init[l.weekday].push(l);
  }
  for (const d of WEEKDAY_ORDER) {
    init[d].sort((a, b) => compareTimeStrings(a.startTime, b.startTime));
  }
  return init;
}

function LessonBlock({
  lesson,
  portalBase,
  onEdit,
  onDelete,
}: {
  lesson: ScheduleLesson;
  portalBase: string;
  onEdit: () => void;
  onDelete: () => void;
}) {
  return (
    <div
      role="button"
      tabIndex={0}
      onClick={onEdit}
      onKeyDown={(e) => {
        if (e.key === "Enter" || e.key === " ") {
          e.preventDefault();
          onEdit();
        }
      }}
      className={cn(
        "group relative cursor-pointer rounded-lg border border-border/80 bg-card p-3 text-left shadow-sm transition-colors",
        "hover:border-primary/40 hover:bg-muted/40 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring",
      )}
    >
      <div className="flex items-start justify-between gap-2">
        <div className="min-w-0 flex-1 space-y-1">
          <p className="flex items-center gap-1.5 font-semibold tabular-nums text-foreground">
            <Clock className="h-3.5 w-3.5 shrink-0 text-muted-foreground" />
            <span>
              {lesson.startTime}–{lesson.endTime}
            </span>
          </p>
          {lesson.title ? (
            <p className="truncate text-sm font-medium leading-tight">
              {lesson.title}
            </p>
          ) : null}
          <p className="flex items-center gap-1.5 truncate text-xs text-muted-foreground">
            <Users className="h-3 w-3 shrink-0" />
            {lesson.group.name}
          </p>
          <p className="flex items-center gap-1.5 truncate text-xs text-muted-foreground">
            <UserRound className="h-3 w-3 shrink-0" />
            {lesson.teacher.fullName}
          </p>
          <p className="flex items-center gap-1.5 truncate text-xs text-muted-foreground">
            <Building2 className="h-3 w-3 shrink-0" />
            {lesson.room.name}
          </p>
        </div>
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <Button
              type="button"
              variant="ghost"
              size="icon"
              className="h-7 w-7 shrink-0 opacity-70 group-hover:opacity-100"
              aria-label="Lesson actions"
              onClick={(e) => e.stopPropagation()}
            >
              <MoreHorizontal className="h-4 w-4" />
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="end" onClick={(e) => e.stopPropagation()}>
            <DropdownMenuItem asChild>
              <Link
                to={`${portalBase}/attendance/lesson/${lesson.id}?date=${isoDateToday()}`}
                className="flex cursor-pointer items-center gap-2"
              >
                <ClipboardCheck className="h-4 w-4" />
                Attendance & grades
              </Link>
            </DropdownMenuItem>
            <DropdownMenuItem
              onSelect={() => {
                onEdit();
              }}
            >
              Edit
            </DropdownMenuItem>
            <DropdownMenuSeparator />
            <DropdownMenuItem
              className="text-destructive focus:text-destructive"
              onSelect={() => onDelete()}
            >
              Delete
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      </div>
    </div>
  );
}

export function SchedulePage() {
  const portalBase = usePortalBase();
  const { data, isLoading, isError } = useScheduleLessonsQuery();
  const deleteMut = useDeleteLessonMutation();

  const [createOpen, setCreateOpen] = useState(false);
  const [editLesson, setEditLesson] = useState<ScheduleLesson | null>(null);
  const [deleteTarget, setDeleteTarget] = useState<ScheduleLesson | null>(null);

  const byDay = useMemo(
    () => lessonsByWeekday(data ?? []),
    [data],
  );

  const totalSlots = data?.length ?? 0;

  return (
    <div className="space-y-6">
      <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
        <div>
          <h1 className="text-2xl font-semibold tracking-tight md:text-3xl">
            Schedule
          </h1>
          <p className="text-sm text-muted-foreground">
            Weekly lesson slots: time, room, teacher, and group at a glance.
          </p>
        </div>
        <Button
          type="button"
          className="w-full gap-2 sm:w-auto"
          onClick={() => setCreateOpen(true)}
        >
          <Plus className="h-4 w-4" />
          Add lesson
        </Button>
      </div>

      {isError ? (
        <p className="text-sm text-destructive">
          Could not load schedule. Check your connection or API settings.
        </p>
      ) : null}

      {isLoading ? (
        <div className="grid grid-cols-1 gap-3 sm:grid-cols-2 lg:grid-cols-7 lg:gap-2">
          {WEEKDAY_ORDER.map((d) => (
            <Card key={d} className="overflow-hidden border-border/80">
              <CardHeader className="py-3">
                <Skeleton className="h-4 w-24" />
              </CardHeader>
              <CardContent className="space-y-2">
                <Skeleton className="h-20 w-full" />
                <Skeleton className="h-20 w-full" />
              </CardContent>
            </Card>
          ))}
        </div>
      ) : totalSlots === 0 ? (
        <AppEmpty
          icon={CalendarDays}
          title="No lessons yet"
          description="Add your first slot to build the weekly timetable."
          action={
            <Button type="button" className="gap-2" onClick={() => setCreateOpen(true)}>
              <Plus className="h-4 w-4" />
              Add lesson
            </Button>
          }
        />
      ) : (
        <>
          <Card className="border-border/80">
            <CardHeader className="pb-2">
              <CardTitle className="text-base">Week overview</CardTitle>
              <CardDescription>
                {totalSlots} slot{totalSlots === 1 ? "" : "s"} · click a card
                to edit.
              </CardDescription>
            </CardHeader>
          </Card>

          <div className="grid grid-cols-1 gap-3 sm:grid-cols-2 lg:grid-cols-7 lg:gap-2">
            {WEEKDAYS.map(({ value: day }) => (
              <Card
                key={day}
                className="flex min-h-[220px] flex-col overflow-hidden border-border/80 lg:min-h-[380px]"
              >
                <CardHeader className="border-b border-border/60 bg-muted/30 py-3">
                  <CardTitle className="text-sm font-semibold">
                    {weekdayLabel(day)}
                  </CardTitle>
                  <p className="text-xs text-muted-foreground">
                    {byDay[day].length} slot
                    {byDay[day].length === 1 ? "" : "s"}
                  </p>
                </CardHeader>
                <CardContent className="flex flex-1 flex-col p-0">
                  <ScrollArea className="h-[min(520px,calc(100vh-220px))] lg:h-[min(520px,60vh)]">
                    <div className="space-y-2 p-3">
                      {byDay[day].length === 0 ? (
                        <p className="py-6 text-center text-xs text-muted-foreground">
                          Free day
                        </p>
                      ) : (
                        byDay[day].map((lesson) => (
                          <LessonBlock
                            key={lesson.id}
                            lesson={lesson}
                            portalBase={portalBase}
                            onEdit={() => setEditLesson(lesson)}
                            onDelete={() => setDeleteTarget(lesson)}
                          />
                        ))
                      )}
                    </div>
                  </ScrollArea>
                </CardContent>
              </Card>
            ))}
          </div>
        </>
      )}

      <LessonCreateDialog open={createOpen} onOpenChange={setCreateOpen} />
      <LessonEditDialog
        lesson={editLesson}
        open={Boolean(editLesson)}
        onOpenChange={(o) => {
          if (!o) setEditLesson(null);
        }}
      />
      <DeleteLessonDialog
        lesson={deleteTarget}
        open={Boolean(deleteTarget)}
        onOpenChange={(o) => {
          if (!o) setDeleteTarget(null);
        }}
        loading={deleteMut.isPending}
        onConfirm={async () => {
          if (!deleteTarget) return;
          await deleteMut.mutateAsync(deleteTarget.id);
          setDeleteTarget(null);
        }}
      />
    </div>
  );
}
