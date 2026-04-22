import type { ReactNode } from "react";

import { cn } from "@/shared/lib/cn";

export type PageHeaderProps = {
  title: string;
  description?: ReactNode;
  /** Primary actions (right side on `sm+`, full width stacked on narrow viewports). */
  actions?: ReactNode;
  /** Extra row under title (badges, context). */
  meta?: ReactNode;
  className?: string;
};

/**
 * Consistent title block for app modules — responsive stacking and spacing.
 */
export function PageHeader({
  title,
  description,
  actions,
  meta,
  className,
}: PageHeaderProps) {
  return (
    <div
      className={cn(
        "flex flex-col gap-4 sm:flex-row sm:items-start sm:justify-between",
        className,
      )}
    >
      <div className="min-w-0 space-y-1">
        <h1 className="text-2xl font-semibold tracking-tight text-foreground md:text-3xl">
          {title}
        </h1>
        {description ? (
          <div className="max-w-2xl text-sm text-muted-foreground md:text-base [&_code]:rounded [&_code]:bg-muted [&_code]:px-1 [&_code]:py-0.5 [&_code]:text-xs">
            {description}
          </div>
        ) : null}
        {meta ? <div className="flex flex-wrap items-center gap-2 pt-1">{meta}</div> : null}
      </div>
      {actions ? (
        <div className="flex w-full shrink-0 flex-col gap-2 sm:w-auto sm:flex-row sm:items-center">
          {actions}
        </div>
      ) : null}
    </div>
  );
}
