import { zodResolver } from "@hookform/resolvers/zod";
import { Loader2 } from "lucide-react";
import { useForm } from "react-hook-form";
import { useEffect } from "react";

import {
  studentFormSchema,
  type StudentFormValues,
} from "@/features/students/model/student-schema";
import type { StudentGroupOption } from "@/features/students/model/types";
import { STUDENT_STATUS_OPTIONS } from "@/features/students/lib/student-status";
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
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/shared/ui/components/select";

const NONE_GROUP = "__none__";

export type StudentFormProps = {
  defaultValues?: Partial<StudentFormValues>;
  groups: StudentGroupOption[];
  groupsLoading?: boolean;
  onSubmit: (values: StudentFormValues) => void | Promise<void>;
  submitLabel?: string;
  id?: string;
};

function toFormValues(
  partial?: Partial<StudentFormValues>,
): StudentFormValues {
  return {
    fullName: partial?.fullName ?? "",
    phone: partial?.phone ?? "",
    email: partial?.email ?? "",
    status: partial?.status ?? "active",
    groupId: partial?.groupId ?? null,
    photoUrl: partial?.photoUrl ?? "",
  };
}

export function StudentForm({
  defaultValues,
  groups,
  groupsLoading,
  onSubmit,
  submitLabel = "Save",
  id,
}: StudentFormProps) {
  const form = useForm<StudentFormValues>({
    resolver: zodResolver(studentFormSchema),
    defaultValues: toFormValues(defaultValues),
  });

  useEffect(() => {
    form.reset(toFormValues(defaultValues));
  }, [defaultValues, form, id]);

  return (
    <Form {...form}>
      <form
        onSubmit={form.handleSubmit(async (values) => {
          await onSubmit({
            ...values,
            groupId: values.groupId ?? null,
            photoUrl: values.photoUrl || undefined,
          });
        })}
        className="space-y-4"
      >
        <div className="grid gap-4 sm:grid-cols-2">
          <FormField
            control={form.control}
            name="fullName"
            render={({ field }) => (
              <FormItem className="sm:col-span-2">
                <FormLabel>Full name</FormLabel>
                <FormControl>
                  <Input placeholder="Full legal name" {...field} />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />
          <FormField
            control={form.control}
            name="phone"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Phone</FormLabel>
                <FormControl>
                  <Input placeholder="+998 90 123 45 67" {...field} />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />
          <FormField
            control={form.control}
            name="email"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Email</FormLabel>
                <FormControl>
                  <Input type="email" placeholder="student@school.edu" {...field} />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />
          <FormField
            control={form.control}
            name="status"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Status</FormLabel>
                <Select
                  onValueChange={field.onChange}
                  value={field.value}
                >
                  <FormControl>
                    <SelectTrigger>
                      <SelectValue placeholder="Status" />
                    </SelectTrigger>
                  </FormControl>
                  <SelectContent>
                    {STUDENT_STATUS_OPTIONS.map((o) => (
                      <SelectItem key={o.value} value={o.value}>
                        {o.label}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
                <FormMessage />
              </FormItem>
            )}
          />
          <FormField
            control={form.control}
            name="groupId"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Group</FormLabel>
                <Select
                  disabled={groupsLoading}
                  onValueChange={(v) =>
                    field.onChange(v === NONE_GROUP ? null : v)
                  }
                  value={field.value ?? NONE_GROUP}
                >
                  <FormControl>
                    <SelectTrigger>
                      <SelectValue placeholder="Assign group" />
                    </SelectTrigger>
                  </FormControl>
                  <SelectContent>
                    <SelectItem value={NONE_GROUP}>No group</SelectItem>
                    {groups.map((g) => (
                      <SelectItem key={g.id} value={g.id}>
                        {g.name}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
                <FormMessage />
              </FormItem>
            )}
          />
          <FormField
            control={form.control}
            name="photoUrl"
            render={({ field }) => (
              <FormItem className="sm:col-span-2">
                <FormLabel>Photo URL</FormLabel>
                <FormControl>
                  <Input
                    placeholder="https://…"
                    {...field}
                    value={field.value ?? ""}
                  />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />
        </div>
        <div className="flex justify-end gap-2 pt-2">
          <Button type="submit" disabled={form.formState.isSubmitting}>
            {form.formState.isSubmitting ? (
              <>
                <Loader2 className="h-4 w-4 animate-spin" />
                Saving…
              </>
            ) : (
              submitLabel
            )}
          </Button>
        </div>
      </form>
    </Form>
  );
}
