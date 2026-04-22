import { zodResolver } from "@hookform/resolvers/zod";
import { Loader2 } from "lucide-react";
import { useEffect, useState } from "react";
import { useForm } from "react-hook-form";

import { useGroupOptionsQuery } from "@/features/students/hooks/use-students";
import {
  teacherFormSchema,
  type TeacherFormValues,
} from "@/features/teachers/model/teacher-schema";
import { useSubjectOptionsQuery } from "@/features/teachers/hooks/use-teachers";
import { TEACHER_STATUS_OPTIONS } from "@/features/teachers/lib/teacher-status";
import { EntityMultiPicker } from "@/shared/ui/registry/entity-multi-picker";
import { ProfilePhotoField } from "@/shared/ui/registry/profile-photo-field";
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
import { Separator } from "@/shared/ui/components/separator";

export type TeacherPhotoSubmit = {
  file: File | null;
  /** Remove server-side photo on save (no new file). */
  clearExisting: boolean;
};

type TeacherFormProps = {
  defaultValues?: Partial<TeacherFormValues>;
  /** Current profile image (edit mode). */
  existingPhotoUrl?: string | null;
  /** Stable key for edit (e.g. teacher id). */
  formId?: string;
  submitLabel?: string;
  onSubmit: (
    values: TeacherFormValues,
    photo: TeacherPhotoSubmit,
  ) => void | Promise<void>;
};

function toFormValues(partial?: Partial<TeacherFormValues>): TeacherFormValues {
  return {
    fullName: partial?.fullName ?? "",
    phone: partial?.phone ?? "",
    email: partial?.email ?? "",
    status: partial?.status ?? "active",
    groupIds: partial?.groupIds ?? [],
    subjectIds: partial?.subjectIds ?? [],
  };
}

export function TeacherForm({
  defaultValues,
  existingPhotoUrl = null,
  formId,
  submitLabel = "Save",
  onSubmit,
}: TeacherFormProps) {
  const { data: groups = [], isLoading: groupsLoading } = useGroupOptionsQuery();
  const { data: subjects = [], isLoading: subjectsLoading } =
    useSubjectOptionsQuery();

  const [photoFile, setPhotoFile] = useState<File | null>(null);
  const [clearPhoto, setClearPhoto] = useState(false);

  const form = useForm<TeacherFormValues>({
    resolver: zodResolver(teacherFormSchema),
    defaultValues: toFormValues(defaultValues),
  });

  useEffect(() => {
    form.reset(toFormValues(defaultValues));
    setPhotoFile(null);
    setClearPhoto(false);
  }, [defaultValues, form, formId]);

  const groupOptions = groups.map((g) => ({ id: g.id, label: g.name }));
  const subjectOptions = subjects.map((s) => ({ id: s.id, label: s.name }));

  return (
    <Form {...form}>
      <form
        onSubmit={form.handleSubmit(async (values) => {
          await onSubmit(values, {
            file: photoFile,
            clearExisting: clearPhoto,
          });
        })}
        className="space-y-6"
      >
        <ProfilePhotoField
          initialsFrom={form.watch("fullName")}
          existingUrl={clearPhoto ? null : existingPhotoUrl}
          onFileChange={(f) => {
            setPhotoFile(f);
            if (f) setClearPhoto(false);
          }}
          onReset={() => setClearPhoto(true)}
          disabled={form.formState.isSubmitting}
        />

        <Separator />

        <div className="grid gap-4 sm:grid-cols-2">
          <FormField
            control={form.control}
            name="fullName"
            render={({ field }) => (
              <FormItem className="sm:col-span-2">
                <FormLabel>Full name</FormLabel>
                <FormControl>
                  <Input placeholder="Full name" {...field} />
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
                  <Input placeholder="+998 …" {...field} />
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
                  <Input type="email" placeholder="name@school.edu" {...field} />
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
                <Select onValueChange={field.onChange} value={field.value}>
                  <FormControl>
                    <SelectTrigger>
                      <SelectValue placeholder="Status" />
                    </SelectTrigger>
                  </FormControl>
                  <SelectContent>
                    {TEACHER_STATUS_OPTIONS.map((o) => (
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
        </div>

        <div className="grid gap-6 sm:grid-cols-2">
          <FormField
            control={form.control}
            name="groupIds"
            render={({ field }) => (
              <FormItem>
                <FormControl>
                  <EntityMultiPicker
                    label="Groups"
                    values={field.value ?? []}
                    onChange={field.onChange}
                    options={groupOptions}
                    placeholder="Assign groups…"
                    disabled={groupsLoading}
                  />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />
          <FormField
            control={form.control}
            name="subjectIds"
            render={({ field }) => (
              <FormItem>
                <FormControl>
                  <EntityMultiPicker
                    label="Subjects"
                    values={field.value ?? []}
                    onChange={field.onChange}
                    options={subjectOptions}
                    placeholder="Assign subjects…"
                    disabled={subjectsLoading}
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
