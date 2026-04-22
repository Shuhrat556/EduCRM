import { MoreHorizontal, Plus, UserRound } from "lucide-react";
import { useState } from "react";

import type { Student, StudentStatus } from "@/features/students/model/types";
import {
  useDeleteStudentMutation,
  useStudentsListQuery,
} from "@/features/students/hooks/use-students";
import {
  STUDENT_STATUS_OPTIONS,
  studentStatusLabel,
} from "@/features/students/lib/student-status";
import { AssignGroupDialog } from "@/features/students/ui/assign-group-dialog";
import { DeleteStudentDialog } from "@/features/students/ui/delete-student-dialog";
import { StudentCreateDialog } from "@/features/students/ui/student-create-dialog";
import { StudentEditDialog } from "@/features/students/ui/student-edit-dialog";
import { StudentProfileSheet } from "@/features/students/ui/student-profile-sheet";
import { DEFAULT_PAGE_SIZE } from "@/shared/constants/pagination";
import { initialsFromName } from "@/shared/lib/format/initials";
import { useDebouncedValue } from "@/shared/lib/hooks/useDebouncedValue";
import { PageSizeSelect } from "@/shared/ui/data/page-size-select";
import { TablePagination } from "@/shared/ui/data/table-pagination";
import { AppEmpty } from "@/shared/ui/feedback/app-empty";
import { QueryErrorAlert } from "@/shared/ui/feedback/query-error-alert";
import { SearchField } from "@/shared/ui/forms/search-field";
import { PageHeader } from "@/shared/ui/layout/page-header";
import { Avatar, AvatarFallback, AvatarImage } from "@/shared/ui/components/avatar";
import { Badge } from "@/shared/ui/components/badge";
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
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/shared/ui/components/select";
import { Skeleton } from "@/shared/ui/components/skeleton";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/shared/ui/components/table";
import { cn } from "@/shared/lib/cn";

function statusBadgeVariant(
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

export function StudentsPage() {
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(DEFAULT_PAGE_SIZE);
  const [searchInput, setSearchInput] = useState("");
  const debouncedSearch = useDebouncedValue(searchInput, 350);
  const [statusFilter, setStatusFilter] = useState<StudentStatus | "all">(
    "all",
  );

  const listParams = {
    page,
    pageSize,
    search: debouncedSearch,
    status: statusFilter,
  };

  const { data, isLoading, isFetching, isError, error, refetch } =
    useStudentsListQuery(listParams);
  const deleteMutation = useDeleteStudentMutation();

  const [createOpen, setCreateOpen] = useState(false);
  const [editStudent, setEditStudent] = useState<Student | null>(null);
  const [deleteTarget, setDeleteTarget] = useState<Student | null>(null);
  const [assignStudent, setAssignStudent] = useState<Student | null>(null);
  const [profileId, setProfileId] = useState<string | null>(null);

  const total = data?.total ?? 0;

  return (
    <div className="space-y-6">
      <PageHeader
        title="Students"
        description="Manage learners, groups, and contact details."
        actions={
          <Button
            type="button"
            className="w-full gap-2 sm:w-auto"
            onClick={() => setCreateOpen(true)}
          >
            <Plus className="h-4 w-4" />
            Add student
          </Button>
        }
      />

      <Card className="border-border/80">
        <CardHeader className="pb-4">
          <CardTitle className="text-base">Directory</CardTitle>
          <CardDescription>
            Search by name or phone, filter by enrollment status.
          </CardDescription>
          <div className="flex flex-col gap-3 pt-2 sm:flex-row sm:items-center sm:flex-wrap">
            <SearchField
              value={searchInput}
              onValueChange={(v) => {
                setSearchInput(v);
                setPage(1);
              }}
              placeholder="Search name or phone…"
              aria-label="Search students"
            />
            <Select
              value={statusFilter}
              onValueChange={(v) => {
                setStatusFilter(v as StudentStatus | "all");
                setPage(1);
              }}
            >
              <SelectTrigger className="w-full sm:w-[180px]">
                <SelectValue placeholder="Status" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">All statuses</SelectItem>
                {STUDENT_STATUS_OPTIONS.map((o) => (
                  <SelectItem key={o.value} value={o.value}>
                    {o.label}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
            <PageSizeSelect
              value={pageSize}
              onChange={(n) => {
                setPageSize(n);
                setPage(1);
              }}
            />
          </div>
        </CardHeader>
        <CardContent className="pt-0">
          {isError ? (
            <QueryErrorAlert
              error={error}
              onRetry={() => void refetch()}
            />
          ) : (
            <>
          <div
            className={cn(
              "rounded-lg border border-border/80",
              isFetching && !isLoading && "opacity-80 transition-opacity",
            )}
          >
            <Table>
              <TableHeader>
                <TableRow className="hover:bg-transparent">
                  <TableHead className="w-[52px]"> </TableHead>
                  <TableHead>Student</TableHead>
                  <TableHead className="hidden md:table-cell">Phone</TableHead>
                  <TableHead className="hidden lg:table-cell">Email</TableHead>
                  <TableHead>Status</TableHead>
                  <TableHead className="hidden sm:table-cell">Group</TableHead>
                  <TableHead className="w-[52px]" />
                </TableRow>
              </TableHeader>
              <TableBody>
                {isLoading ? (
                  Array.from({ length: 5 }).map((_, i) => (
                    <TableRow key={i}>
                      <TableCell>
                        <Skeleton className="h-9 w-9 rounded-full" />
                      </TableCell>
                      <TableCell>
                        <Skeleton className="h-4 w-36" />
                      </TableCell>
                      <TableCell className="hidden md:table-cell">
                        <Skeleton className="h-4 w-28" />
                      </TableCell>
                      <TableCell className="hidden lg:table-cell">
                        <Skeleton className="h-4 w-40" />
                      </TableCell>
                      <TableCell>
                        <Skeleton className="h-5 w-16" />
                      </TableCell>
                      <TableCell className="hidden sm:table-cell">
                        <Skeleton className="h-4 w-24" />
                      </TableCell>
                      <TableCell />
                    </TableRow>
                  ))
                ) : data?.items.length === 0 ? (
                  <TableRow>
                    <TableCell colSpan={7} className="h-48 border-0 p-0">
                      <AppEmpty
                        icon={UserRound}
                        title="No students match"
                        description="Try another search or clear filters to see everyone."
                        className="border-0 bg-transparent py-8"
                      />
                    </TableCell>
                  </TableRow>
                ) : (
                  data!.items.map((row) => (
                    <TableRow
                      key={row.id}
                      className="cursor-pointer"
                      onClick={() => setProfileId(row.id)}
                    >
                      <TableCell onClick={(e) => e.stopPropagation()}>
                        <Avatar className="h-9 w-9 border border-border/60">
                          {row.photoUrl ? (
                            <AvatarImage src={row.photoUrl} alt="" />
                          ) : null}
                          <AvatarFallback className="text-[10px]">
                            {initialsFromName(row.fullName)}
                          </AvatarFallback>
                        </Avatar>
                      </TableCell>
                      <TableCell className="font-medium">{row.fullName}</TableCell>
                      <TableCell className="hidden tabular-nums text-muted-foreground md:table-cell">
                        {row.phone}
                      </TableCell>
                      <TableCell className="hidden max-w-[200px] truncate text-muted-foreground lg:table-cell">
                        {row.email}
                      </TableCell>
                      <TableCell>
                        <Badge
                          variant={statusBadgeVariant(row.status)}
                          className="font-normal capitalize"
                        >
                          {studentStatusLabel(row.status)}
                        </Badge>
                      </TableCell>
                      <TableCell className="hidden text-muted-foreground sm:table-cell">
                        {row.groupName ?? "—"}
                      </TableCell>
                      <TableCell onClick={(e) => e.stopPropagation()}>
                        <DropdownMenu>
                          <DropdownMenuTrigger asChild>
                            <Button
                              type="button"
                              variant="ghost"
                              size="icon"
                              className="h-8 w-8"
                              aria-label={`Actions for ${row.fullName}`}
                            >
                              <MoreHorizontal className="h-4 w-4" />
                            </Button>
                          </DropdownMenuTrigger>
                          <DropdownMenuContent align="end">
                            <DropdownMenuItem
                              onSelect={() => setProfileId(row.id)}
                            >
                              View profile
                            </DropdownMenuItem>
                            <DropdownMenuItem
                              onSelect={() => setEditStudent(row)}
                            >
                              Edit
                            </DropdownMenuItem>
                            <DropdownMenuItem
                              onSelect={() => setAssignStudent(row)}
                            >
                              Assign to group
                            </DropdownMenuItem>
                            <DropdownMenuSeparator />
                            <DropdownMenuItem
                              className="text-destructive focus:text-destructive"
                              onSelect={() => setDeleteTarget(row)}
                            >
                              Delete
                            </DropdownMenuItem>
                          </DropdownMenuContent>
                        </DropdownMenu>
                      </TableCell>
                    </TableRow>
                  ))
                )}
              </TableBody>
            </Table>
          </div>

          {!isLoading ? (
            <TablePagination
              className="mt-4"
              page={page}
              pageSize={pageSize}
              total={total}
              disabled={isFetching}
              onPageChange={setPage}
            />
          ) : null}
            </>
          )}
        </CardContent>
      </Card>

      <StudentCreateDialog open={createOpen} onOpenChange={setCreateOpen} />

      <StudentEditDialog
        student={editStudent}
        open={Boolean(editStudent)}
        onOpenChange={(o) => {
          if (!o) setEditStudent(null);
        }}
      />

      <AssignGroupDialog
        student={assignStudent}
        open={Boolean(assignStudent)}
        onOpenChange={(o) => {
          if (!o) setAssignStudent(null);
        }}
      />

      <DeleteStudentDialog
        student={deleteTarget}
        open={Boolean(deleteTarget)}
        onOpenChange={(o) => {
          if (!o) setDeleteTarget(null);
        }}
        loading={deleteMutation.isPending}
        onConfirm={async () => {
          if (!deleteTarget) return;
          await deleteMutation.mutateAsync(deleteTarget.id);
          setDeleteTarget(null);
        }}
      />

      <StudentProfileSheet
        studentId={profileId}
        open={Boolean(profileId)}
        onOpenChange={(o) => {
          if (!o) setProfileId(null);
        }}
        onEdit={(s) => {
          setProfileId(null);
          setEditStudent(s);
        }}
        onAssignGroup={(s) => {
          setProfileId(null);
          setAssignStudent(s);
        }}
      />
    </div>
  );
}
