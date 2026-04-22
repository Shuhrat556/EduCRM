import { Check, ChevronsUpDown } from "lucide-react";
import { useMemo, useState } from "react";

import { Button } from "@/shared/ui/components/button";
import { Checkbox } from "@/shared/ui/components/checkbox";
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/shared/ui/components/popover";
import { ScrollArea } from "@/shared/ui/components/scroll-area";
import { cn } from "@/shared/lib/cn";

export type EntityOption = { id: string; label: string };

type EntityMultiPickerProps = {
  label: string;
  values: string[];
  onChange: (ids: string[]) => void;
  options: EntityOption[];
  placeholder?: string;
  disabled?: boolean;
  emptyHint?: string;
  className?: string;
};

/**
 * Reusable multi-select (e.g. groups, subjects) aligned with admin form patterns.
 */
export function EntityMultiPicker({
  label,
  values,
  onChange,
  options,
  placeholder = "Select…",
  disabled,
  emptyHint = "No options available.",
  className,
}: EntityMultiPickerProps) {
  const [open, setOpen] = useState(false);

  const labelById = useMemo(() => {
    const m = new Map(options.map((o) => [o.id, o.label] as const));
    return m;
  }, [options]);

  const summary =
    values.length === 0
      ? placeholder
      : values.length === 1
        ? labelById.get(values[0]!) ?? "1 selected"
        : `${values.length} selected`;

  return (
    <div className={cn("space-y-2", className)}>
      <span className="text-sm font-medium leading-none">{label}</span>
      <Popover open={open} onOpenChange={setOpen}>
        <PopoverTrigger asChild>
          <Button
            type="button"
            variant="outline"
            role="combobox"
            aria-expanded={open}
            className="h-auto min-h-9 w-full justify-between px-3 py-2 font-normal"
            disabled={disabled || options.length === 0}
          >
            <span
              className={cn(
                "truncate text-left",
                values.length === 0 && "text-muted-foreground",
              )}
            >
              {summary}
            </span>
            <ChevronsUpDown className="h-4 w-4 shrink-0 opacity-50" />
          </Button>
        </PopoverTrigger>
        <PopoverContent className="w-[var(--radix-popover-trigger-width)] p-0" align="start">
          {options.length === 0 ? (
            <p className="p-3 text-sm text-muted-foreground">{emptyHint}</p>
          ) : (
            <ScrollArea className="max-h-60">
              <ul className="p-1">
                {options.map((opt) => {
                  const checked = values.includes(opt.id);
                  return (
                    <li key={opt.id}>
                      <label
                        className="flex cursor-pointer items-center gap-2 rounded-sm px-2 py-2 text-sm hover:bg-accent"
                      >
                        <Checkbox
                          checked={checked}
                          onCheckedChange={(c) => {
                            if (c === true) {
                              if (!values.includes(opt.id)) {
                                onChange([...values, opt.id]);
                              }
                            } else {
                              onChange(values.filter((v) => v !== opt.id));
                            }
                          }}
                        />
                        <span className="flex-1 truncate">{opt.label}</span>
                        {checked ? (
                          <Check className="h-4 w-4 shrink-0 text-primary" />
                        ) : null}
                      </label>
                    </li>
                  );
                })}
              </ul>
            </ScrollArea>
          )}
        </PopoverContent>
      </Popover>
    </div>
  );
}
