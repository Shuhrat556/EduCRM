import { useState } from "react";

import {
  useCreateTeacherMutation,
  useUploadTeacherPhotoMutation,
} from "@/features/teachers/hooks/use-teachers";
import type { TeacherFormValues } from "@/features/teachers/model/teacher-schema";
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

type TeacherCreateDialogProps = {
  open: boolean;
  onOpenChange: (open: boolean) => void;
};

export function TeacherCreateDialog({
  open,
  onOpenChange,
}: TeacherCreateDialogProps) {
  const [formKey, setFormKey] = useState(0);
  const create = useCreateTeacherMutation();
  const upload = useUploadTeacherPhotoMutation();

  async function handleSubmit(
    values: TeacherFormValues,
    photo: TeacherPhotoSubmit,
  ) {
    try {
      const teacher = await create.mutateAsync({
        fullName: values.fullName,
        phone: values.phone,
        email: values.email,
        status: values.status,
        groupIds: values.groupIds,
        subjectIds: values.subjectIds,
        photoUrl: null,
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
    <Dialog
      open={open}
      onOpenChange={(o) => {
        onOpenChange(o);
        if (o) setFormKey((k) => k + 1);
      }}
    >
      <DialogContent className="max-h-[min(90vh,760px)] overflow-y-auto sm:max-w-xl">
        <DialogHeader>
          <DialogTitle>Add teacher</DialogTitle>
        </DialogHeader>
        <TeacherForm
          key={formKey}
          submitLabel="Create teacher"
          onSubmit={handleSubmit}
        />
      </DialogContent>
    </Dialog>
  );
}
