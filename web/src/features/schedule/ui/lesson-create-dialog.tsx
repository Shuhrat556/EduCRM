import { useState } from "react";

import { useCreateLessonMutation } from "@/features/schedule/hooks/use-schedule";
import type { LessonFormValues } from "@/features/schedule/model/lesson-schema";
import { LessonForm } from "@/features/schedule/ui/lesson-form";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "@/shared/ui/components/dialog";

type LessonCreateDialogProps = {
  open: boolean;
  onOpenChange: (open: boolean) => void;
};

export function LessonCreateDialog({
  open,
  onOpenChange,
}: LessonCreateDialogProps) {
  const [formKey, setFormKey] = useState(0);
  const create = useCreateLessonMutation();

  async function handleSubmit(values: LessonFormValues) {
    try {
      await create.mutateAsync({
        weekday: values.weekday,
        startTime: values.startTime,
        endTime: values.endTime,
        roomId: values.roomId,
        teacherId: values.teacherId,
        groupId: values.groupId,
        title: values.title,
      });
      onOpenChange(false);
    } catch {
      /* keep open */
    }
  }

  return (
    <Dialog
      open={open}
      onOpenChange={(o) => {
        onOpenChange(o);
        if (o) setFormKey((k) => k + 1);
      }}
    >
      <DialogContent className="max-h-[min(90vh,760px)] overflow-y-auto sm:max-w-md">
        <DialogHeader>
          <DialogTitle>Add lesson slot</DialogTitle>
        </DialogHeader>
        <LessonForm
          key={formKey}
          submitLabel="Create lesson"
          onSubmit={handleSubmit}
        />
      </DialogContent>
    </Dialog>
  );
}
