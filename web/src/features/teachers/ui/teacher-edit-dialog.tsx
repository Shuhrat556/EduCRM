import { useMemo } from "react";

import {
  useUpdateTeacherMutation,
  useUploadTeacherPhotoMutation,
} from "@/features/teachers/hooks/use-teachers";
import type { TeacherFormValues } from "@/features/teachers/model/teacher-schema";
import type { Teacher } from "@/features/teachers/model/types";
import {
  TeacherForm,
  type TeacherPhotoSubmit,
} from "@/features/teachers/ui/teacher-form";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "@/shared/ui/components/dialog";

type TeacherEditDialogProps = {
  teacher: Teacher | null;
  open: boolean;
  onOpenChange: (open: boolean) => void;
};

export function TeacherEditDialog({
  teacher,
  open,
  onOpenChange,
}: TeacherEditDialogProps) {
  const update = useUpdateTeacherMutation();
  const upload = useUploadTeacherPhotoMutation();

  const formDefaults = useMemo(
    () =>
      teacher
        ? {
            fullName: teacher.fullName,
            phone: teacher.phone,
            email: teacher.email,
            status: teacher.status,
            groupIds: [...teacher.groupIds],
            subjectIds: [...teacher.subjectIds],
          }
        : undefined,
    [teacher],
  );

  async function handleSubmit(
    values: TeacherFormValues,
    photo: TeacherPhotoSubmit,
  ) {
    if (!teacher) return;
    try {
      await update.mutateAsync({
        id: teacher.id,
        payload: {
          fullName: values.fullName,
          phone: values.phone,
          email: values.email,
          status: values.status,
          groupIds: values.groupIds,
          subjectIds: values.subjectIds,
          photoUrl: photo.clearExisting ? null : undefined,
        },
      });
      if (photo.file) {
        await upload.mutateAsync({ id: teacher.id, file: photo.file });
      }
      onOpenChange(false);
    } catch {
      /* keep open */
    }
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-h-[min(90vh,760px)] overflow-y-auto sm:max-w-xl">
        <DialogHeader>
          <DialogTitle>Edit teacher</DialogTitle>
        </DialogHeader>
        {teacher ? (
          <TeacherForm
            formId={teacher.id}
            key={teacher.id}
            defaultValues={formDefaults}
            existingPhotoUrl={teacher.photoUrl}
            submitLabel="Save changes"
            onSubmit={handleSubmit}
          />
        ) : null}
      </DialogContent>
    </Dialog>
  );
}
