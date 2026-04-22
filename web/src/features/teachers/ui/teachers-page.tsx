import { MoreHorizontal, Plus, UserRound } from "lucide-react";
import { useState } from "react";
import { Link, useNavigate } from "react-router-dom";

import type { Teacher, TeacherStatus } from "@/features/teachers/model/types";
import {
  useDeleteTeacherMutation,
  useTeachersListQuery,
} from "@/features/teachers/hooks/use-teachers";
import {
  TEACHER_STATUS_OPTIONS,
  teacherStatusLabel,
} from "@/features/teachers/lib/teacher-status";
import { DeleteTeacherDialog } from "@/features/teachers/ui/delete-teacher-dialog";
import { TeacherCreateDialog } from "@/features/teachers/ui/teacher-create-dialog";
import { TeacherEditDialog } from "@/features/teachers/ui/teacher-edit-dialog";
import { DEFAULT_PAGE_SIZE } from "@/shared/constants/pagination";
import { initialsFromName } from "@/shared/lib/format/initials";
import { useDebouncedValue } from "@/shared/lib/hooks/useDebouncedValue";
import { PageSizeSelect } from "@/shared/ui/data/page-size-select";
import { TablePagination } from "@/shared/ui/data/table-pagination";
import { usePortalBase } from "@/shared/lib/hooks/usePortalBase";
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
  s: TeacherStatus,
): "default" | "secondary" | "outline" | "destructive" {
  switch (s) {
    case "active":
      return "secondary";
    case "inactive":
      return "outline";
    case "on_leave":
      return "destructive";
    default:
      return "outline";
  }
}

export function TeachersPage() {
  const base = usePortalBase();
  const navigate = useNavigate();
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(DEFAULT_PAGE_SIZE);
  const [searchInput, setSearchInput] = useState("");
  const debouncedSearch = useDebouncedValue(searchInput, 350);
  const [statusFilter, setStatusFilter] = useState<TeacherStatus | "all">(
    "all",
  );

  const listParams = {
    page,
    pageSize,
    search: debouncedSearch,
    status: statusFilter,
  };

  const { data, isLoading, isFetching, isError, error, refetch } =
    useTeachersListQuery(listParams);
  const deleteMutation = useDeleteTeacherMutation();

  const [createOpen, setCreateOpen] = useState(false);
  const [editTeacher, setEditTeacher] = useState<Teacher | null>(null);
  const [deleteTarget, setDeleteTarget] = useState<Teacher | null>(null);

  const total = data?.total ?? 0;

  return (
    <div className="space-y-6">
      <PageHeader
        title="Teachers"
        description="Faculty directory, subjects, and group assignments."
        actions={
          <Button
            type="button"
            className="w-full gap-2 sm:w-auto"
            onClick={() => setCreateOpen(true)}
          >
            <Plus className="h-4 w-4" />
            Add teacher
          </Button>
        }
      />

      <Card className="border-border/80">
        <CardHeader className="pb-4">
          <CardTitle className="text-base">Directory</CardTitle>
          <CardDescription>
            Search by name, phone, or email. Open a row for the full profile
            page.
          </CardDescription>
          <div className="flex flex-col gap-3 pt-2 sm:flex-row sm:items-center sm:flex-wrap">
            <SearchField
              value={searchInput}
              onValueChange={(v) => {
                setSearchInput(v);
                setPage(1);
              }}
              aria-label="Search teachers"
            />
            <Select
              value={statusFilter}
              onValueChange={(v) => {
                setStatusFilter(v as TeacherStatus | "all");
                setPage(1);
              }}
            >
              <SelectTrigger className="w-full sm:w-[180px]">
                <SelectValue placeholder="Status" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">All statuses</SelectItem>
                {TEACHER_STATUS_OPTIONS.map((o) => (
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
            <QueryErrorAlert error={error} onRetry={() => void refetch()} />
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
                  <TableHead className="w-[52px]" />
                  <TableHead>Teacher</TableHead>
                  <TableHead className="hidden md:table-cell">Phone</TableHead>
                  <TableHead className="hidden lg:table-cell">Email</TableHead>
                  <TableHead>Status</TableHead>
                  <TableHead className="hidden xl:table-cell">Subjects</TableHead>
                  <TableHead className="hidden xl:table-cell">Groups</TableHead>
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
                      <TableCell className="hidden xl:table-cell">
                        <Skeleton className="h-4 w-32" />
                      </TableCell>
                      <TableCell className="hidden xl:table-cell">
                        <Skeleton className="h-4 w-24" />
                      </TableCell>
                      <TableCell />
                    </TableRow>
                  ))
                ) : data?.items.length === 0 ? (
                  <TableRow>
                    <TableCell colSpan={8} className="h-48 border-0 p-0">
                      <AppEmpty
                        icon={UserRound}
                        title="No teachers match"
                        description="Try another search or add a new teacher."
                        className="border-0 bg-transparent py-8"
                      />
                    </TableCell>
                  </TableRow>
                ) : (
                  data!.items.map((row) => (
                    <TableRow
                      key={row.id}
                      className="cursor-pointer"
                      onClick={() => navigate(`${base}/teachers/${row.id}`)}
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
                          {teacherStatusLabel(row.status)}
                        </Badge>
                      </TableCell>
                      <TableCell className="hidden max-w-[180px] truncate text-xs text-muted-foreground xl:table-cell">
                        {row.subjects.map((s) => s.name).join(", ") || "—"}
                      </TableCell>
                      <TableCell className="hidden max-w-[160px] truncate text-xs text-muted-foreground xl:table-cell">
                        {row.groups.map((s) => s.name).join(", ") || "—"}
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
                            <DropdownMenuItem asChild>
                              <Link to={`${base}/teachers/${row.id}`}>
                                View profile
                              </Link>
                            </DropdownMenuItem>
                            <DropdownMenuItem
                              onSelect={() => setEditTeacher(row)}
                            >
                              Edit
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

      <TeacherCreateDialog open={createOpen} onOpenChange={setCreateOpen} />
      <TeacherEditDialog
        teacher={editTeacher}
        open={Boolean(editTeacher)}
        onOpenChange={(o) => {
          if (!o) setEditTeacher(null);
        }}
      />
      <DeleteTeacherDialog
        teacher={deleteTarget}
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
    </div>
  );
}
