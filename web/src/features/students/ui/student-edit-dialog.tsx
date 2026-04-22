import { useMemo } from "react";

import {
  useGroupOptionsQuery,
  useUpdateStudentMutation,
} from "@/features/students/hooks/use-students";
import type { StudentFormValues } from "@/features/students/model/student-schema";
import type { Student } from "@/features/students/model/types";
import { StudentForm } from "@/features/students/ui/student-form";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "@/shared/ui/components/dialog";

type StudentEditDialogProps = {
  student: Student | null;
  open: boolean;
  onOpenChange: (open: boolean) => void;
};

export function StudentEditDialog({
  student,
  open,
  onOpenChange,
}: StudentEditDialogProps) {
  const { data: groups = [], isLoading: groupsLoading } = useGroupOptionsQuery();
  const update = useUpdateStudentMutation();

  const formDefaults = useMemo(
    () =>
      student
        ? {
            fullName: student.fullName,
            phone: student.phone,
            email: student.email,
            status: student.status,
            groupId: student.groupId,
            photoUrl: student.photoUrl ?? "",
          }
        : undefined,
    [student],
  );

  async function handleSubmit(values: StudentFormValues) {
    if (!student) return;
    try {
      await update.mutateAsync({
        id: student.id,
        payload: {
          fullName: values.fullName,
          phone: values.phone,
          email: values.email,
          status: values.status,
          groupId: values.groupId,
          photoUrl: values.photoUrl?.trim() || null,
        },
      });
      onOpenChange(false);
    } catch {
      /* keep dialog open */
    }
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-h-[min(90vh,720px)] overflow-y-auto sm:max-w-lg">
        <DialogHeader>
          <DialogTitle>Edit student</DialogTitle>
        </DialogHeader>
        {student ? (
          <StudentForm
            id={student.id}
            key={student.id}
            defaultValues={formDefaults}
            groups={groups}
            groupsLoading={groupsLoading}
            submitLabel="Save changes"
            onSubmit={handleSubmit}
          />
        ) : null}
      </DialogContent>
    </Dialog>
  );
}
