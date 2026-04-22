import { MoreHorizontal, Plus, Users } from "lucide-react";
import { useState } from "react";
import { Link, useNavigate } from "react-router-dom";

import type { Group, GroupStatus } from "@/features/groups/model/types";
import { formatGroupFee } from "@/features/groups/lib/format-group-fee";
import {
  GROUP_STATUS_OPTIONS,
  groupStatusLabel,
} from "@/features/groups/lib/group-status";
import {
  useDeleteGroupMutation,
  useGroupsListQuery,
} from "@/features/groups/hooks/use-groups";
import { DeleteGroupDialog } from "@/features/groups/ui/delete-group-dialog";
import { GroupCreateDialog } from "@/features/groups/ui/group-create-dialog";
import { GroupEditDialog } from "@/features/groups/ui/group-edit-dialog";
import { DEFAULT_PAGE_SIZE } from "@/shared/constants/pagination";
import { formatIsoDateDisplay } from "@/shared/lib/format/iso-date-display";
import { useDebouncedValue } from "@/shared/lib/hooks/useDebouncedValue";
import { usePortalBase } from "@/shared/lib/hooks/usePortalBase";
import { PageSizeSelect } from "@/shared/ui/data/page-size-select";
import { TablePagination } from "@/shared/ui/data/table-pagination";
import { AppEmpty } from "@/shared/ui/feedback/app-empty";
import { QueryErrorAlert } from "@/shared/ui/feedback/query-error-alert";
import { SearchField } from "@/shared/ui/forms/search-field";
import { PageHeader } from "@/shared/ui/layout/page-header";
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
  s: GroupStatus,
): "default" | "secondary" | "outline" | "destructive" {
  switch (s) {
    case "active":
      return "secondary";
    case "draft":
      return "outline";
    case "paused":
      return "outline";
    case "completed":
      return "default";
    case "archived":
      return "destructive";
    default:
      return "outline";
  }
}

export function GroupsPage() {
  const base = usePortalBase();
  const navigate = useNavigate();
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(DEFAULT_PAGE_SIZE);
  const [searchInput, setSearchInput] = useState("");
  const debouncedSearch = useDebouncedValue(searchInput, 350);
  const [statusFilter, setStatusFilter] = useState<GroupStatus | "all">("all");

  const listParams = {
    page,
    pageSize,
    search: debouncedSearch,
    status: statusFilter,
  };

  const { data, isLoading, isFetching, isError, error, refetch } =
    useGroupsListQuery(listParams);
  const deleteMutation = useDeleteGroupMutation();

  const [createOpen, setCreateOpen] = useState(false);
  const [editGroup, setEditGroup] = useState<Group | null>(null);
  const [deleteTarget, setDeleteTarget] = useState<Group | null>(null);

  const total = data?.total ?? 0;

  return (
    <div className="space-y-6">
      <PageHeader
        title="Groups"
        description="Cohorts, fees, resources, and assignments. Students belong to one group at a time."
        actions={
          <Button
            type="button"
            className="w-full gap-2 sm:w-auto"
            onClick={() => setCreateOpen(true)}
          >
            <Plus className="h-4 w-4" />
            Create group
          </Button>
        }
      />

      <Card className="border-border/80">
        <CardHeader className="pb-4">
          <CardTitle className="text-base">All groups</CardTitle>
          <CardDescription>
            Search by name. Open a row for roster and schedule preview.
          </CardDescription>
          <div className="flex flex-col gap-3 pt-2 sm:flex-row sm:items-center sm:flex-wrap">
            <SearchField
              value={searchInput}
              onValueChange={(v) => {
                setSearchInput(v);
                setPage(1);
              }}
              aria-label="Search groups"
            />
            <Select
              value={statusFilter}
              onValueChange={(v) => {
                setStatusFilter(v as GroupStatus | "all");
                setPage(1);
              }}
            >
              <SelectTrigger className="w-full sm:w-[180px]">
                <SelectValue placeholder="Status" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">All statuses</SelectItem>
                {GROUP_STATUS_OPTIONS.map((o) => (
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
                  <TableHead>Group</TableHead>
                  <TableHead className="hidden md:table-cell">Teacher</TableHead>
                  <TableHead className="hidden lg:table-cell">Subject</TableHead>
                  <TableHead className="hidden xl:table-cell">Room</TableHead>
                  <TableHead className="text-right tabular-nums">Fee / mo</TableHead>
                  <TableHead>Status</TableHead>
                  <TableHead className="hidden sm:table-cell">Dates</TableHead>
                  <TableHead className="w-[52px]" />
                </TableRow>
              </TableHeader>
              <TableBody>
                {isLoading ? (
                  Array.from({ length: 5 }).map((_, i) => (
                    <TableRow key={i}>
                      <TableCell>
                        <Skeleton className="h-4 w-36" />
                      </TableCell>
                      <TableCell className="hidden md:table-cell">
                        <Skeleton className="h-4 w-28" />
                      </TableCell>
                      <TableCell className="hidden lg:table-cell">
                        <Skeleton className="h-4 w-24" />
                      </TableCell>
                      <TableCell className="hidden xl:table-cell">
                        <Skeleton className="h-4 w-20" />
                      </TableCell>
                      <TableCell className="text-right">
                        <Skeleton className="ml-auto h-4 w-16" />
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
                    <TableCell colSpan={8} className="h-48 border-0 p-0">
                      <AppEmpty
                        icon={Users}
                        title="No groups match"
                        description="Try another filter or create a new group."
                        className="border-0 bg-transparent py-8"
                      />
                    </TableCell>
                  </TableRow>
                ) : (
                  data!.items.map((row) => (
                    <TableRow
                      key={row.id}
                      className="cursor-pointer"
                      onClick={() => navigate(`${base}/groups/${row.id}`)}
                    >
                      <TableCell className="font-medium">{row.name}</TableCell>
                      <TableCell className="hidden text-muted-foreground md:table-cell">
                        {row.teacher?.fullName ?? "—"}
                      </TableCell>
                      <TableCell className="hidden text-muted-foreground lg:table-cell">
                        {row.subject?.name ?? "—"}
                      </TableCell>
                      <TableCell className="hidden text-muted-foreground xl:table-cell">
                        {row.room?.name ?? "—"}
                      </TableCell>
                      <TableCell className="text-right tabular-nums text-muted-foreground">
                        {formatGroupFee(row.monthlyFee)}
                      </TableCell>
                      <TableCell>
                        <Badge
                          variant={statusBadgeVariant(row.status)}
                          className="font-normal capitalize"
                        >
                          {groupStatusLabel(row.status)}
                        </Badge>
                      </TableCell>
                      <TableCell className="hidden text-xs text-muted-foreground sm:table-cell">
                        {formatIsoDateDisplay(row.startDate)}
                        {row.endDate
                          ? ` → ${formatIsoDateDisplay(row.endDate)}`
                          : ""}
                      </TableCell>
                      <TableCell onClick={(e) => e.stopPropagation()}>
                        <DropdownMenu>
                          <DropdownMenuTrigger asChild>
                            <Button
                              type="button"
                              variant="ghost"
                              size="icon"
                              className="h-8 w-8"
                              aria-label={`Actions for ${row.name}`}
                            >
                              <MoreHorizontal className="h-4 w-4" />
                            </Button>
                          </DropdownMenuTrigger>
                          <DropdownMenuContent align="end">
                            <DropdownMenuItem asChild>
                              <Link to={`${base}/groups/${row.id}`}>
                                View details
                              </Link>
                            </DropdownMenuItem>
                            <DropdownMenuItem
                              onSelect={() => setEditGroup(row)}
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

      <GroupCreateDialog open={createOpen} onOpenChange={setCreateOpen} />
      <GroupEditDialog
        group={editGroup}
        open={Boolean(editGroup)}
        onOpenChange={(o) => {
          if (!o) setEditGroup(null);
        }}
      />
      <DeleteGroupDialog
        group={deleteTarget}
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
