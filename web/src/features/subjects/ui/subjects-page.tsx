import { BookOpen, MoreHorizontal, Plus } from "lucide-react";
import { useState } from "react";

import type { Subject } from "@/features/subjects/model/types";
import {
  useDeleteSubjectMutation,
  useSubjectsListQuery,
} from "@/features/subjects/hooks/use-subjects";
import { DeleteSubjectDialog } from "@/features/subjects/ui/delete-subject-dialog";
import { SubjectCreateDialog } from "@/features/subjects/ui/subject-create-dialog";
import { SubjectEditDialog } from "@/features/subjects/ui/subject-edit-dialog";
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

export function SubjectsPage() {
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(DEFAULT_PAGE_SIZE);
  const [searchInput, setSearchInput] = useState("");
  const debouncedSearch = useDebouncedValue(searchInput, 350);

  const listParams = { page, pageSize, search: debouncedSearch };
  const { data, isLoading, isFetching, isError, error, refetch } =
    useSubjectsListQuery(listParams);
  const deleteMutation = useDeleteSubjectMutation();

  const [createOpen, setCreateOpen] = useState(false);
  const [editSubject, setEditSubject] = useState<Subject | null>(null);
  const [deleteTarget, setDeleteTarget] = useState<Subject | null>(null);

  const total = data?.total ?? 0;

  return (
    <div className="space-y-6">
      <PageHeader
        title="Subjects"
        description="Manage the catalog used across teachers and scheduling."
        actions={
          <Button
            type="button"
            className="w-full gap-2 sm:w-auto"
            onClick={() => setCreateOpen(true)}
          >
            <Plus className="h-4 w-4" />
            Add subject
          </Button>
        }
      />

      <Card className="border-border/80">
        <CardHeader className="pb-4">
          <CardTitle className="text-base">Catalog</CardTitle>
          <CardDescription>
            Search by name, code, or description.
          </CardDescription>
          <div className="flex flex-col gap-3 pt-2 sm:flex-row sm:items-center sm:flex-wrap">
            <SearchField
              value={searchInput}
              onValueChange={(v) => {
                setSearchInput(v);
                setPage(1);
              }}
              aria-label="Search subjects"
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
                  <TableHead>Name</TableHead>
                  <TableHead className="hidden sm:table-cell">Code</TableHead>
                  <TableHead className="hidden md:table-cell max-w-[280px]">
                    Description
                  </TableHead>
                  <TableHead className="w-[52px]" />
                </TableRow>
              </TableHeader>
              <TableBody>
                {isLoading ? (
                  Array.from({ length: 5 }).map((_, i) => (
                    <TableRow key={i}>
                      <TableCell>
                        <Skeleton className="h-4 w-40" />
                      </TableCell>
                      <TableCell className="hidden sm:table-cell">
                        <Skeleton className="h-4 w-16" />
                      </TableCell>
                      <TableCell className="hidden md:table-cell">
                        <Skeleton className="h-4 w-48" />
                      </TableCell>
                      <TableCell />
                    </TableRow>
                  ))
                ) : data?.items.length === 0 ? (
                  <TableRow>
                    <TableCell colSpan={4} className="h-48 border-0 p-0">
                      <AppEmpty
                        icon={BookOpen}
                        title="No subjects found"
                        description="Try a different search or add a subject."
                        className="border-0 bg-transparent py-8"
                      />
                    </TableCell>
                  </TableRow>
                ) : (
                  data!.items.map((row) => (
                    <TableRow key={row.id}>
                      <TableCell className="font-medium">{row.name}</TableCell>
                      <TableCell className="hidden tabular-nums text-muted-foreground sm:table-cell">
                        {row.code ?? "—"}
                      </TableCell>
                      <TableCell className="hidden max-w-[280px] truncate text-sm text-muted-foreground md:table-cell">
                        {row.description ?? "—"}
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
                            <DropdownMenuItem
                              onSelect={() => setEditSubject(row)}
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

      <SubjectCreateDialog open={createOpen} onOpenChange={setCreateOpen} />
      <SubjectEditDialog
        subject={editSubject}
        open={Boolean(editSubject)}
        onOpenChange={(o) => {
          if (!o) setEditSubject(null);
        }}
      />
      <DeleteSubjectDialog
        subject={deleteTarget}
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
