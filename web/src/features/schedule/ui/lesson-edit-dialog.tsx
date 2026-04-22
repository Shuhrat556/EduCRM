import type { ScheduleLesson } from "@/features/schedule/model/types";
import { useUpdateLessonMutation } from "@/features/schedule/hooks/use-schedule";
import type { LessonFormValues } from "@/features/schedule/model/lesson-schema";
import { LessonForm } from "@/features/schedule/ui/lesson-form";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "@/shared/ui/components/dialog";

type LessonEditDialogProps = {
  lesson: ScheduleLesson | null;
  open: boolean;
  onOpenChange: (open: boolean) => void;
};

export function LessonEditDialog({
  lesson,
  open,
  onOpenChange,
}: LessonEditDialogProps) {
  const update = useUpdateLessonMutation();

  async function handleSubmit(values: LessonFormValues) {
    if (!lesson) return;
    try {
      await update.mutateAsync({
        id: lesson.id,
        payload: {
          weekday: values.weekday,
          startTime: values.startTime,
          endTime: values.endTime,
          roomId: values.roomId,
          teacherId: values.teacherId,
          groupId: values.groupId,
          title: values.title,
        },
      });
      onOpenChange(false);
    } catch {
      /* keep open */
    }
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-h-[min(90vh,760px)] overflow-y-auto sm:max-w-md">
        <DialogHeader>
          <DialogTitle>Edit lesson slot</DialogTitle>
        </DialogHeader>
        {lesson ? (
          <LessonForm
            id={lesson.id}
            defaultValues={{
              weekday: lesson.weekday,
              startTime: lesson.startTime,
              endTime: lesson.endTime,
              roomId: lesson.roomId,
              teacherId: lesson.teacherId,
              groupId: lesson.groupId,
              title: lesson.title ?? "",
            }}
            submitLabel="Save changes"
            onSubmit={handleSubmit}
          />
        ) : null}
      </DialogContent>
    </Dialog>
  );
}
