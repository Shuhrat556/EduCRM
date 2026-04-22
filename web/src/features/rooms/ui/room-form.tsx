import { zodResolver } from "@hookform/resolvers/zod";
import { Loader2 } from "lucide-react";
import { useEffect } from "react";
import { useForm } from "react-hook-form";

import {
  roomFormSchema,
  type RoomFormValues,
} from "@/features/rooms/model/room-schema";
import { Button } from "@/shared/ui/components/button";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/shared/ui/components/form";
import { Input } from "@/shared/ui/components/input";
import { cn } from "@/shared/lib/cn";

export type RoomFormProps = {
  defaultValues?: Partial<RoomFormValues>;
  onSubmit: (values: RoomFormValues) => void | Promise<void>;
  submitLabel?: string;
  id?: string;
};

function toFormValues(partial?: Partial<RoomFormValues>): RoomFormValues {
  return {
    name: partial?.name ?? "",
    building: partial?.building ?? "",
    capacity: partial?.capacity ?? 20,
    notes: partial?.notes ?? "",
  };
}

export function RoomForm({
  defaultValues,
  onSubmit,
  submitLabel = "Save",
  id,
}: RoomFormProps) {
  const form = useForm<RoomFormValues>({
    resolver: zodResolver(roomFormSchema),
    defaultValues: toFormValues(defaultValues),
  });

  useEffect(() => {
    form.reset(toFormValues(defaultValues));
  }, [defaultValues, form, id]);

  return (
    <Form {...form}>
      <form
        onSubmit={form.handleSubmit(async (values) => {
          await onSubmit(values);
        })}
        className="space-y-4"
      >
        <FormField
          control={form.control}
          name="name"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Room name</FormLabel>
              <FormControl>
                <Input placeholder="e.g. Lab A" autoComplete="off" {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />
        <FormField
          control={form.control}
          name="building"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Building (optional)</FormLabel>
              <FormControl>
                <Input
                  placeholder="e.g. Main hall"
                  autoComplete="off"
                  {...field}
                />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />
        <FormField
          control={form.control}
          name="capacity"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Capacity (seats)</FormLabel>
              <FormControl>
                <Input
                  type="number"
                  min={1}
                  max={9999}
                  className="tabular-nums"
                  {...field}
                  onChange={(e) => field.onChange(e.target.value)}
                />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />
        <FormField
          control={form.control}
          name="notes"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Notes (optional)</FormLabel>
              <FormControl>
                <textarea
                  className={cn(
                    "flex min-h-[72px] w-full rounded-md border border-input bg-transparent px-3 py-2 text-base shadow-sm transition-colors placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring disabled:cursor-not-allowed disabled:opacity-50 md:text-sm",
                  )}
                  placeholder="Equipment, access, booking hints"
                  {...field}
                />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />
        <Button
          type="submit"
          className="w-full gap-2 sm:w-auto"
          disabled={form.formState.isSubmitting}
        >
          {form.formState.isSubmitting ? (
            <>
              <Loader2 className="h-4 w-4 animate-spin" />
              Saving…
            </>
          ) : (
            submitLabel
          )}
        </Button>
      </form>
    </Form>
  );
}
