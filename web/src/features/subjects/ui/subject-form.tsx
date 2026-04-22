import { zodResolver } from "@hookform/resolvers/zod";
import { Loader2 } from "lucide-react";
import { useEffect } from "react";
import { useForm } from "react-hook-form";

import {
  subjectFormSchema,
  type SubjectFormValues,
} from "@/features/subjects/model/subject-schema";
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

export type SubjectFormProps = {
  defaultValues?: Partial<SubjectFormValues>;
  onSubmit: (values: SubjectFormValues) => void | Promise<void>;
  submitLabel?: string;
  id?: string;
};

function toFormValues(partial?: Partial<SubjectFormValues>): SubjectFormValues {
  return {
    name: partial?.name ?? "",
    code: partial?.code ?? "",
    description: partial?.description ?? "",
  };
}

export function SubjectForm({
  defaultValues,
  onSubmit,
  submitLabel = "Save",
  id,
}: SubjectFormProps) {
  const form = useForm<SubjectFormValues>({
    resolver: zodResolver(subjectFormSchema),
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
              <FormLabel>Name</FormLabel>
              <FormControl>
                <Input placeholder="e.g. Mathematics" autoComplete="off" {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />
        <FormField
          control={form.control}
          name="code"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Code (optional)</FormLabel>
              <FormControl>
                <Input placeholder="e.g. MATH-101" autoComplete="off" {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />
        <FormField
          control={form.control}
          name="description"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Description (optional)</FormLabel>
              <FormControl>
                <textarea
                  className={cn(
                    "flex min-h-[88px] w-full rounded-md border border-input bg-transparent px-3 py-2 text-base shadow-sm transition-colors placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring disabled:cursor-not-allowed disabled:opacity-50 md:text-sm",
                  )}
                  placeholder="Brief notes for staff"
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
