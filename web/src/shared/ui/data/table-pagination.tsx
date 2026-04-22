import { ChevronLeft, ChevronRight } from "lucide-react";

import { Button } from "@/shared/ui/components/button";
import { cn } from "@/shared/lib/cn";

export type TablePaginationProps = {
  page: number;
  pageSize: number;
  total: number;
  onPageChange: (page: number) => void;
  /** Disable controls (e.g. while a mutation runs). */
  disabled?: boolean;
  className?: string;
};

export function TablePagination({
  page,
  pageSize,
  total,
  onPageChange,
  disabled,
  className,
}: TablePaginationProps) {
  if (total <= 0) return null;

  const totalPages = Math.max(1, Math.ceil(total / pageSize));
  const from = (page - 1) * pageSize + 1;
  const to = Math.min(page * pageSize, total);

  return (
    <div
      className={cn(
        "flex flex-col items-center justify-between gap-3 sm:flex-row",
        className,
      )}
    >
      <p className="order-2 text-center text-xs text-muted-foreground sm:order-1 sm:text-left">
        Showing{" "}
        <span className="font-medium text-foreground">
          {from}–{to}
        </span>{" "}
        of <span className="font-medium text-foreground">{total}</span>
      </p>
      <div className="order-1 flex items-center gap-2 sm:order-2">
        <Button
          type="button"
          variant="outline"
          size="sm"
          disabled={disabled || page <= 1}
          onClick={() => onPageChange(Math.max(1, page - 1))}
        >
          <ChevronLeft className="h-4 w-4 sm:mr-1" aria-hidden />
          <span className="hidden sm:inline">Previous</span>
        </Button>
        <span className="min-w-[5.5rem] text-center text-sm tabular-nums text-muted-foreground">
          {page} / {totalPages}
        </span>
        <Button
          type="button"
          variant="outline"
          size="sm"
          disabled={disabled || page >= totalPages}
          onClick={() => onPageChange(Math.min(totalPages, page + 1))}
        >
          <span className="hidden sm:inline">Next</span>
          <ChevronRight className="h-4 w-4 sm:ml-1" aria-hidden />
        </Button>
      </div>
    </div>
  );
}
