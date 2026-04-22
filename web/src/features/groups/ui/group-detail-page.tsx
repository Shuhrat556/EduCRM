import {
  ArrowLeft,
  Building2,
  Calendar,
  Clock,
  Pencil,
  Trash2,
  UserRound,
} from "lucide-react";
import { useState } from "react";
import { Link, useNavigate, useParams } from "react-router-dom";

import { formatGroupFee } from "@/features/groups/lib/format-group-fee";
import { groupStatusLabel } from "@/features/groups/lib/group-status";
import {
  useDeleteGroupMutation,
  useGroupQuery,
} from "@/features/groups/hooks/use-groups";
import { DeleteGroupDialog } from "@/features/groups/ui/delete-group-dialog";
import { GroupEditDialog } from "@/features/groups/ui/group-edit-dialog";
import { studentStatusLabel } from "@/features/students/lib/student-status";
import { useStudentsListQuery } from "@/features/students/hooks/use-students";
import type { StudentStatus } from "@/features/students/model/types";
import { Badge } from "@/shared/ui/components/badge";
import { Button } from "@/shared/ui/components/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/shared/ui/components/card";
import { Skeleton } from "@/shared/ui/components/skeleton";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/shared/ui/components/table";
import { formatIsoDateDisplay } from "@/shared/lib/format/iso-date-display";
import { usePortalBase } from "@/shared/lib/hooks/usePortalBase";
import { QueryErrorAlert } from "@/shared/ui/feedback/query-error-alert";

function studentBadgeVariant(
  s: StudentStatus,
): "default" | "secondary" | "outline" | "destructive" {
  switch (s) {
    case "active":
      return "secondary";
    case "inactive":
      return "outline";
    case "graduated":
      return "default";
    case "suspended":
      return "destructive";
    default:
      return "outline";
  }
}

export function GroupDetailPage() {
  const base = usePortalBase();
  const { groupId } = useParams<{ groupId: string }>();
  const navigate = useNavigate();
  const id = groupId ?? null;

  const { data: group, isLoading, isError, error, refetch } = useGroupQuery(id);
  const { data: studentsData, isLoading: studentsLoading } =
    useStudentsListQuery(
      {
        page: 1,
        pageSize: 100,
        search: "",
        status: "all",
        groupId: id!,
      },
      { enabled: Boolean(id) },
    );
  const removeGroup = useDeleteGroupMutation();

  const [editOpen, setEditOpen] = useState(false);
  const [deleteOpen, setDeleteOpen] = useState(false);

  if (!id) {
    return (
      <p className="text-sm text-muted-foreground">Invalid group link.</p>
    );
  }

  return (
    <div className="space-y-6">
      <div className="flex flex-wrap items-center gap-3">
        <Button variant="ghost" size="sm" asChild className="gap-2 px-2">
          <Link to={`${base}/groups`}>
            <ArrowLeft className="h-4 w-4" />
            Groups
          </Link>
        </Button>
      </div>

      {isLoading ? (
        <div className="space-y-4">
          <Skeleton className="h-10 w-64" />
          <Skeleton className="h-40 w-full max-w-2xl" />
        </div>
      ) : isError ? (
        <QueryErrorAlert
          error={error}
          title="Could not load group"
          onRetry={() => void refetch()}
        />
      ) : !group ? (
        <p className="text-sm text-muted-foreground">Group not found.</p>
      ) : (
        <>
          <div className="flex flex-col gap-6 lg:flex-row lg:items-start lg:justify-between">
            <div className="space-y-2">
              <h1 className="text-2xl font-semibold tracking-tight md:text-3xl">
                {group.name}
              </h1>
              <div className="flex flex-wrap gap-2">
                <Badge variant="secondary">
                  {groupStatusLabel(group.status)}
                </Badge>
                {group.subject ? (
                  <Badge variant="outline">{group.subject.name}</Badge>
                ) : null}
                {group.room ? (
                  <Badge variant="outline">
                    <Building2 className="mr-1 h-3 w-3" />
                    {group.room.name} · cap {group.room.capacity}
                  </Badge>
                ) : null}
              </div>
              <p className="text-sm text-muted-foreground">
                Monthly fee{" "}
                <span className="font-medium text-foreground">
                  {formatGroupFee(group.monthlyFee)}
                </span>
                · {formatIsoDateDisplay(group.startDate)}
                {group.endDate
                  ? ` → ${formatIsoDateDisplay(group.endDate)}`
                  : " · Open-ended"}
              </p>
              {group.teacher ? (
                <p className="flex items-center gap-2 text-sm text-muted-foreground">
                  <UserRound className="h-4 w-4 shrink-0" />
                  <span>
                    Teacher:{" "}
                    <span className="font-medium text-foreground">
                      {group.teacher.fullName}
                    </span>
                  </span>
                </p>
              ) : (
                <p className="text-sm text-muted-foreground">No teacher assigned.</p>
              )}
            </div>
            <div className="flex flex-wrap gap-2">
              <Button
                type="button"
                variant="outline"
                className="gap-2"
                onClick={() => setEditOpen(true)}
              >
                <Pencil className="h-4 w-4" />
                Edit
              </Button>
              <Button
                type="button"
                variant="destructive"
                className="gap-2"
                onClick={() => setDeleteOpen(true)}
              >
                <Trash2 className="h-4 w-4" />
                Delete
              </Button>
            </div>
          </div>

          <div className="grid gap-6 lg:grid-cols-2">
            <Card className="border-border/80 lg:col-span-2">
              <CardHeader>
                <CardTitle className="flex items-center gap-2 text-base">
                  <Clock className="h-4 w-4" />
                  Schedule preview
                </CardTitle>
                <CardDescription>
                  Sample weekly pattern for this cohort (demo). Replace with
                  live timetable when the schedule module is connected.
                </CardDescription>
              </CardHeader>
              <CardContent>
                <Table>
                  <TableHeader>
                    <TableRow className="hover:bg-transparent">
                      <TableHead>Day</TableHead>
                      <TableHead>Time</TableHead>
                      <TableHead className="hidden sm:table-cell">Subject</TableHead>
                      <TableHead className="hidden sm:table-cell">Room</TableHead>
                    </TableRow>
                  </TableHeader>
                  <TableBody>
                    {group.schedulePreview.length === 0 ? (
                      <TableRow>
                        <TableCell
                          colSpan={4}
                          className="text-sm text-muted-foreground"
                        >
                          No preview slots yet.
                        </TableCell>
                      </TableRow>
                    ) : (
                      group.schedulePreview.map((slot, i) => (
                        <TableRow key={`${slot.weekday}-${i}`}>
                          <TableCell className="font-medium">{slot.weekday}</TableCell>
                          <TableCell className="tabular-nums text-muted-foreground">
                            {slot.startTime}–{slot.endTime}
                          </TableCell>
                          <TableCell className="hidden sm:table-cell">
                            {slot.subjectName}
                          </TableCell>
                          <TableCell className="hidden sm:table-cell">
                            {slot.roomName}
                          </TableCell>
                        </TableRow>
                      ))
                    )}
                  </TableBody>
                </Table>
              </CardContent>
            </Card>

            <Card className="border-border/80 lg:col-span-2">
              <CardHeader>
                <CardTitle className="flex items-center gap-2 text-base">
                  <Calendar className="h-4 w-4" />
                  Students in this group
                </CardTitle>
                <CardDescription>
                  MVP rule: each student appears in at most one group. Assign
                  from the Students module.
                </CardDescription>
              </CardHeader>
              <CardContent>
                <div className="rounded-lg border border-border/80">
                  <Table>
                    <TableHeader>
                      <TableRow className="hover:bg-transparent">
                        <TableHead>Name</TableHead>
                        <TableHead className="hidden md:table-cell">Email</TableHead>
                        <TableHead>Status</TableHead>
                      </TableRow>
                    </TableHeader>
                    <TableBody>
                      {studentsLoading ? (
                        Array.from({ length: 3 }).map((_, i) => (
                          <TableRow key={i}>
                            <TableCell>
                              <Skeleton className="h-4 w-40" />
                            </TableCell>
                            <TableCell className="hidden md:table-cell">
                              <Skeleton className="h-4 w-48" />
                            </TableCell>
                            <TableCell>
                              <Skeleton className="h-5 w-16" />
                            </TableCell>
                          </TableRow>
                        ))
                      ) : studentsData?.items.length === 0 ? (
                        <TableRow>
                          <TableCell
                            colSpan={3}
                            className="py-10 text-center text-sm text-muted-foreground"
                          >
                            No students assigned yet.
                          </TableCell>
                        </TableRow>
                      ) : (
                        studentsData!.items.map((s) => (
                          <TableRow key={s.id}>
                            <TableCell className="font-medium">{s.fullName}</TableCell>
                            <TableCell className="hidden max-w-[220px] truncate text-muted-foreground md:table-cell">
                              {s.email}
                            </TableCell>
                            <TableCell>
                              <Badge
                                variant={studentBadgeVariant(s.status)}
                                className="font-normal capitalize"
                              >
                                {studentStatusLabel(s.status)}
                              </Badge>
                            </TableCell>
                          </TableRow>
                        ))
                      )}
                    </TableBody>
                  </Table>
                </div>
              </CardContent>
            </Card>
          </div>

          <GroupEditDialog
            group={group}
            open={editOpen}
            onOpenChange={setEditOpen}
          />
          <DeleteGroupDialog
            group={group}
            open={deleteOpen}
            onOpenChange={setDeleteOpen}
            loading={removeGroup.isPending}
            onConfirm={async () => {
              await removeGroup.mutateAsync(group.id);
              setDeleteOpen(false);
              navigate(`${base}/groups`, { replace: true });
            }}
          />
        </>
      )}
    </div>
  );
}
