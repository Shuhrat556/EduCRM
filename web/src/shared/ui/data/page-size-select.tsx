import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/shared/ui/components/select";
import {
  DEFAULT_PAGE_SIZE,
  PAGE_SIZE_OPTIONS,
  type PageSizeOption,
} from "@/shared/constants/pagination";
import { cn } from "@/shared/lib/cn";

type PageSizeSelectProps = {
  value: PageSizeOption;
  onChange: (size: PageSizeOption) => void;
  className?: string;
  triggerClassName?: string;
};

export function PageSizeSelect({
  value,
  onChange,
  className,
  triggerClassName,
}: PageSizeSelectProps) {
  return (
    <div className={cn(className)}>
      <Select
        value={String(value)}
        onValueChange={(v) => {
          const n = Number(v);
          const next = PAGE_SIZE_OPTIONS.includes(n as PageSizeOption)
            ? n
            : DEFAULT_PAGE_SIZE;
          onChange(next as PageSizeOption);
        }}
      >
        <SelectTrigger className={cn("w-full sm:w-[130px]", triggerClassName)}>
          <SelectValue />
        </SelectTrigger>
        <SelectContent>
          {PAGE_SIZE_OPTIONS.map((n) => (
            <SelectItem key={n} value={String(n)}>
              {n} / page
            </SelectItem>
          ))}
        </SelectContent>
      </Select>
    </div>
  );
}
