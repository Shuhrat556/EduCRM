import { Search } from "lucide-react";

import { Input } from "@/shared/ui/components/input";
import { cn } from "@/shared/lib/cn";

export type SearchFieldProps = {
  value: string;
  onValueChange: (value: string) => void;
  placeholder?: string;
  "aria-label"?: string;
  className?: string;
  inputClassName?: string;
  disabled?: boolean;
};

export function SearchField({
  value,
  onValueChange,
  placeholder = "Search…",
  "aria-label": ariaLabel = "Search",
  className,
  inputClassName,
  disabled,
}: SearchFieldProps) {
  return (
    <div className={cn("relative flex-1", className)}>
      <Search
        className="pointer-events-none absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground"
        aria-hidden
      />
      <Input
        className={cn("pl-9", inputClassName)}
        placeholder={placeholder}
        value={value}
        disabled={disabled}
        aria-label={ariaLabel}
        onChange={(e) => onValueChange(e.target.value)}
      />
    </div>
  );
}
