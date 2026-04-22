import { Building2, MoreHorizontal, Plus } from "lucide-react";
import { useState } from "react";

import type { Room } from "@/features/rooms/model/types";
import {
  useDeleteRoomMutation,
  useRoomsListQuery,
} from "@/features/rooms/hooks/use-rooms";
import { DeleteRoomDialog } from "@/features/rooms/ui/delete-room-dialog";
import { RoomCreateDialog } from "@/features/rooms/ui/room-create-dialog";
import { RoomEditDialog } from "@/features/rooms/ui/room-edit-dialog";
import { DEFAULT_PAGE_SIZE } from "@/shared/constants/pagination";
import { useDebouncedValue } from "@/shared/lib/hooks/useDebouncedValue";
import { PageSizeSelect } from "@/shared/ui/data/page-size-select";
import { TablePagination } from "@/shared/ui/data/table-pagination";
import { AppEmpty } from "@/shared/ui/feedback/app-empty";
import { QueryErrorAlert } from "@/shared/ui/feedback/query-error-alert";
import { SearchField } from "@/shared/ui/forms/search-field";
import { PageHeader } from "@/shared/ui/layout/page-header";
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

export function RoomsPage() {
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(DEFAULT_PAGE_SIZE);
  const [searchInput, setSearchInput] = useState("");
  const debouncedSearch = useDebouncedValue(searchInput, 350);

  const listParams = { page, pageSize, search: debouncedSearch };
  const { data, isLoading, isFetching, isError, error, refetch } =
    useRoomsListQuery(listParams);
  const deleteMutation = useDeleteRoomMutation();

  const [createOpen, setCreateOpen] = useState(false);
  const [editRoom, setEditRoom] = useState<Room | null>(null);
  const [deleteTarget, setDeleteTarget] = useState<Room | null>(null);

  const total = data?.total ?? 0;

  return (
    <div className="space-y-6">
      <PageHeader
        title="Rooms"
        description="Spaces, capacity, and quick notes for scheduling."
        actions={
          <Button
            type="button"
            className="w-full gap-2 sm:w-auto"
            onClick={() => setCreateOpen(true)}
          >
            <Plus className="h-4 w-4" />
            Add room
          </Button>
        }
      />

      <Card className="border-border/80">
        <CardHeader className="pb-4">
          <CardTitle className="text-base">Inventory</CardTitle>
          <CardDescription>
            Search by room name, building, capacity, or notes.
          </CardDescription>
          <div className="flex flex-col gap-3 pt-2 sm:flex-row sm:items-center sm:flex-wrap">
            <SearchField
              value={searchInput}
              onValueChange={(v) => {
                setSearchInput(v);
                setPage(1);
              }}
              aria-label="Search rooms"
            />
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
                  <TableHead>Room</TableHead>
                  <TableHead className="hidden sm:table-cell">Building</TableHead>
                  <TableHead className="text-right tabular-nums">Capacity</TableHead>
                  <TableHead className="hidden md:table-cell max-w-[240px]">
                    Notes
                  </TableHead>
                  <TableHead className="w-[52px]" />
                </TableRow>
              </TableHeader>
              <TableBody>
                {isLoading ? (
                  Array.from({ length: 5 }).map((_, i) => (
                    <TableRow key={i}>
                      <TableCell>
                        <Skeleton className="h-4 w-32" />
                      </TableCell>
                      <TableCell className="hidden sm:table-cell">
                        <Skeleton className="h-4 w-24" />
                      </TableCell>
                      <TableCell className="text-right">
                        <Skeleton className="ml-auto h-4 w-10" />
                      </TableCell>
                      <TableCell className="hidden md:table-cell">
                        <Skeleton className="h-4 w-40" />
                      </TableCell>
                      <TableCell />
                    </TableRow>
                  ))
                ) : data?.items.length === 0 ? (
                  <TableRow>
                    <TableCell colSpan={5} className="h-48 border-0 p-0">
                      <AppEmpty
                        icon={Building2}
                        title="No rooms found"
                        description="Try another search or add a room."
                        className="border-0 bg-transparent py-8"
                      />
                    </TableCell>
                  </TableRow>
                ) : (
                  data!.items.map((row) => (
                    <TableRow key={row.id}>
                      <TableCell className="font-medium">{row.name}</TableCell>
                      <TableCell className="hidden text-muted-foreground sm:table-cell">
                        {row.building ?? "—"}
                      </TableCell>
                      <TableCell className="text-right tabular-nums font-medium">
                        {row.capacity}
                      </TableCell>
                      <TableCell className="hidden max-w-[240px] truncate text-sm text-muted-foreground md:table-cell">
                        {row.notes ?? "—"}
                      </TableCell>
                      <TableCell>
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
                            <DropdownMenuItem onSelect={() => setEditRoom(row)}>
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

      <RoomCreateDialog open={createOpen} onOpenChange={setCreateOpen} />
      <RoomEditDialog
        room={editRoom}
        open={Boolean(editRoom)}
        onOpenChange={(o) => {
          if (!o) setEditRoom(null);
        }}
      />
      <DeleteRoomDialog
        room={deleteTarget}
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
