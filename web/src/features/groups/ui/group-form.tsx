import { zodResolver } from "@hookform/resolvers/zod";
import { Loader2 } from "lucide-react";
import { useEffect } from "react";
import { useForm } from "react-hook-form";

import {
  GROUP_FORM_NONE,
  groupFormSchema,
  type GroupFormValues,
} from "@/features/groups/model/group-schema";
import { GROUP_STATUS_OPTIONS } from "@/features/groups/lib/group-status";
import type { GroupStatus } from "@/features/groups/model/types";
import { useRoomsListQuery } from "@/features/rooms/hooks/use-rooms";
import {
  useSubjectOptionsQuery,
  useTeachersListQuery,
} from "@/features/teachers/hooks/use-teachers";
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

export type GroupFormProps = {
  defaultValues?: Partial<GroupFormValues>;
  onSubmit: (values: GroupFormValues) => void | Promise<void>;
  submitLabel?: string;
  id?: string;
};

function toFormValues(partial?: Partial<GroupFormValues>): GroupFormValues {
  return {
    name: partial?.name ?? "",
    teacherId: partial?.teacherId ?? GROUP_FORM_NONE,
    subjectId: partial?.subjectId ?? GROUP_FORM_NONE,
    roomId: partial?.roomId ?? GROUP_FORM_NONE,
    monthlyFee: partial?.monthlyFee ?? 0,
    startDate: partial?.startDate ?? "",
    endDate: partial?.endDate ?? "",
    status: (partial?.status ?? "draft") as GroupStatus,
  };
}

export function GroupForm({
  defaultValues,
  onSubmit,
  submitLabel = "Save",
  id,
}: GroupFormProps) {
  const form = useForm<GroupFormValues>({
    resolver: zodResolver(groupFormSchema),
    defaultValues: toFormValues(defaultValues),
  });

  const { data: teachersData, isLoading: teachersLoading } =
    useTeachersListQuery({
      page: 1,
      pageSize: 100,
      search: "",
      status: "all",
    });
  const { data: subjects = [], isLoading: subjectsLoading } =
    useSubjectOptionsQuery();
  const { data: roomsData, isLoading: roomsLoading } = useRoomsListQuery({
    page: 1,
    pageSize: 100,
    search: "",
  });

  const teachers = teachersData?.items ?? [];
  const rooms = roomsData?.items ?? [];

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
              <FormLabel>Group name</FormLabel>
              <FormControl>
                <Input placeholder="e.g. Grade 10-A" autoComplete="off" {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />

        <div className="grid gap-4 sm:grid-cols-2">
          <FormField
            control={form.control}
            name="teacherId"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Teacher</FormLabel>
                <Select
                  value={field.value ?? GROUP_FORM_NONE}
                  onValueChange={field.onChange}
                  disabled={teachersLoading}
                >
                  <FormControl>
                    <SelectTrigger>
                      <SelectValue placeholder="Select teacher" />
                    </SelectTrigger>
                  </FormControl>
                  <SelectContent>
                    <SelectItem value={GROUP_FORM_NONE}>None</SelectItem>
                    {teachers.map((t) => (
                      <SelectItem key={t.id} value={t.id}>
                        {t.fullName}
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
            name="subjectId"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Subject</FormLabel>
                <Select
                  value={field.value ?? GROUP_FORM_NONE}
                  onValueChange={field.onChange}
                  disabled={subjectsLoading}
                >
                  <FormControl>
                    <SelectTrigger>
                      <SelectValue placeholder="Select subject" />
                    </SelectTrigger>
                  </FormControl>
                  <SelectContent>
                    <SelectItem value={GROUP_FORM_NONE}>None</SelectItem>
                    {subjects.map((s) => (
                      <SelectItem key={s.id} value={s.id}>
                        {s.name}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
                <FormMessage />
              </FormItem>
            )}
          />
        </div>

        <FormField
          control={form.control}
          name="roomId"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Room</FormLabel>
              <Select
                value={field.value ?? GROUP_FORM_NONE}
                onValueChange={field.onChange}
                disabled={roomsLoading}
              >
                <FormControl>
                  <SelectTrigger>
                    <SelectValue placeholder="Select room" />
                  </SelectTrigger>
                </FormControl>
                <SelectContent>
                  <SelectItem value={GROUP_FORM_NONE}>None</SelectItem>
                  {rooms.map((r) => (
                    <SelectItem key={r.id} value={r.id}>
                      {r.name}
                      {r.building ? ` · ${r.building}` : ""}
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
          name="monthlyFee"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Monthly fee (whole units)</FormLabel>
              <FormControl>
                <Input
                  type="number"
                  min={0}
                  className="tabular-nums"
                  {...field}
                  onChange={(e) => field.onChange(e.target.value)}
                />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />

        <div className="grid gap-4 sm:grid-cols-2">
          <FormField
            control={form.control}
            name="startDate"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Start date</FormLabel>
                <FormControl>
                  <Input type="date" {...field} />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />
          <FormField
            control={form.control}
            name="endDate"
            render={({ field }) => (
              <FormItem>
                <FormLabel>End date (optional)</FormLabel>
                <FormControl>
                  <Input
                    type="date"
                    value={field.value ?? ""}
                    onChange={(e) =>
                      field.onChange(e.target.value || undefined)
                    }
                  />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />
        </div>

        <FormField
          control={form.control}
          name="status"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Status</FormLabel>
              <Select value={field.value} onValueChange={field.onChange}>
                <FormControl>
                  <SelectTrigger>
                    <SelectValue />
                  </SelectTrigger>
                </FormControl>
                <SelectContent>
                  {GROUP_STATUS_OPTIONS.map((o) => (
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

        <p className="text-xs text-muted-foreground">
          One primary teacher per group. The same teacher may lead several
          groups. Students can only be in one group at a time.
        </p>

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
