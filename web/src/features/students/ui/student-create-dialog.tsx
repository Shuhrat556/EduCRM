import { useState } from "react";

import {
  useCreateStudentMutation,
  useGroupOptionsQuery,
} from "@/features/students/hooks/use-students";
import type { StudentFormValues } from "@/features/students/model/student-schema";
import { StudentForm } from "@/features/students/ui/student-form";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "@/shared/ui/components/dialog";

type StudentCreateDialogProps = {
  open: boolean;
  onOpenChange: (open: boolean) => void;
};

export function StudentCreateDialog({
  open,
  onOpenChange,
}: StudentCreateDialogProps) {
  const [formKey, setFormKey] = useState(0);
  const { data: groups = [], isLoading: groupsLoading } = useGroupOptionsQuery();
  const create = useCreateStudentMutation();

  async function handleSubmit(values: StudentFormValues) {
    try {
      await create.mutateAsync({
        fullName: values.fullName,
        phone: values.phone,
        email: values.email,
        status: values.status,
        groupId: values.groupId,
        photoUrl: values.photoUrl?.trim() || null,
      });
      onOpenChange(false);
    } catch {
      /* network / validation — keep dialog open */
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
      <DialogContent className="max-h-[min(90vh,720px)] overflow-y-auto sm:max-w-lg">
        <DialogHeader>
          <DialogTitle>Add student</DialogTitle>
        </DialogHeader>
        <StudentForm
          key={formKey}
          groups={groups}
          groupsLoading={groupsLoading}
          submitLabel="Create student"
          onSubmit={handleSubmit}
        />
      </DialogContent>
    </Dialog>
  );
}
