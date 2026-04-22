import { zodResolver } from "@hookform/resolvers/zod";
import { Loader2 } from "lucide-react";
import { useEffect } from "react";
import { useForm } from "react-hook-form";

import { useGroupsListQuery } from "@/features/groups/hooks/use-groups";
import {
  lessonFormSchema,
  type LessonFormValues,
} from "@/features/schedule/model/lesson-schema";
import { WEEKDAYS } from "@/features/schedule/lib/weekdays";
import { useRoomsListQuery } from "@/features/rooms/hooks/use-rooms";
import { useTeachersListQuery } from "@/features/teachers/hooks/use-teachers";
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

export type LessonFormProps = {
  defaultValues?: Partial<LessonFormValues>;
  onSubmit: (values: LessonFormValues) => void | Promise<void>;
  submitLabel?: string;
  id?: string;
};

function toFormValues(
  partial?: Partial<LessonFormValues>,
): LessonFormValues {
  return {
    weekday: partial?.weekday ?? "monday",
    startTime: partial?.startTime ?? "09:00",
    endTime: partial?.endTime ?? "10:00",
    roomId: partial?.roomId ?? "",
    teacherId: partial?.teacherId ?? "",
    groupId: partial?.groupId ?? "",
    title: partial?.title ?? "",
  };
}

export function LessonForm({
  defaultValues,
  onSubmit,
  submitLabel = "Save",
  id,
}: LessonFormProps) {
  const form = useForm<LessonFormValues>({
    resolver: zodResolver(lessonFormSchema),
    defaultValues: toFormValues(defaultValues),
  });

  const { data: teachersData, isLoading: teachersLoading } =
    useTeachersListQuery({
      page: 1,
      pageSize: 100,
      search: "",
      status: "all",
    });
  const { data: roomsData, isLoading: roomsLoading } = useRoomsListQuery({
    page: 1,
    pageSize: 100,
    search: "",
  });
  const { data: groupsData, isLoading: groupsLoading } = useGroupsListQuery({
    page: 1,
    pageSize: 100,
    search: "",
    status: "all",
  });

  const teachers = teachersData?.items ?? [];
  const rooms = roomsData?.items ?? [];
  const groups = groupsData?.items ?? [];

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
          name="title"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Title (optional)</FormLabel>
              <FormControl>
                <Input
                  placeholder="e.g. Mathematics — Chapter 4"
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
          name="weekday"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Weekday</FormLabel>
              <Select value={field.value} onValueChange={field.onChange}>
                <FormControl>
                  <SelectTrigger>
                    <SelectValue />
                  </SelectTrigger>
                </FormControl>
                <SelectContent>
                  {WEEKDAYS.map((d) => (
                    <SelectItem key={d.value} value={d.value}>
                      {d.label}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
              <FormMessage />
            </FormItem>
          )}
        />

        <div className="grid grid-cols-2 gap-3">
          <FormField
            control={form.control}
            name="startTime"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Start</FormLabel>
                <FormControl>
                  <Input type="time" step={300} {...field} />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />
          <FormField
            control={form.control}
            name="endTime"
            render={({ field }) => (
              <FormItem>
                <FormLabel>End</FormLabel>
                <FormControl>
                  <Input type="time" step={300} {...field} />
                </FormControl>
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
                value={field.value}
                onValueChange={field.onChange}
                disabled={roomsLoading}
              >
                <FormControl>
                  <SelectTrigger>
                    <SelectValue placeholder="Select room" />
                  </SelectTrigger>
                </FormControl>
                <SelectContent>
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
          name="teacherId"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Teacher</FormLabel>
              <Select
                value={field.value}
                onValueChange={field.onChange}
                disabled={teachersLoading}
              >
                <FormControl>
                  <SelectTrigger>
                    <SelectValue placeholder="Select teacher" />
                  </SelectTrigger>
                </FormControl>
                <SelectContent>
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
          name="groupId"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Group</FormLabel>
              <Select
                value={field.value}
                onValueChange={field.onChange}
                disabled={groupsLoading}
              >
                <FormControl>
                  <SelectTrigger>
                    <SelectValue placeholder="Select group" />
                  </SelectTrigger>
                </FormControl>
                <SelectContent>
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
