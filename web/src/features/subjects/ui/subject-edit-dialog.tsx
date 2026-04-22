import type { Subject } from "@/features/subjects/model/types";
import { useUpdateSubjectMutation } from "@/features/subjects/hooks/use-subjects";
import {
  subjectFormToPayload,
  type SubjectFormValues,
} from "@/features/subjects/model/subject-schema";
import { SubjectForm } from "@/features/subjects/ui/subject-form";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "@/shared/ui/components/dialog";

type SubjectEditDialogProps = {
  subject: Subject | null;
  open: boolean;
  onOpenChange: (open: boolean) => void;
};

export function SubjectEditDialog({
  subject,
  open,
  onOpenChange,
}: SubjectEditDialogProps) {
  const update = useUpdateSubjectMutation();

  async function handleSubmit(values: SubjectFormValues) {
    if (!subject) return;
    try {
      await update.mutateAsync({
        id: subject.id,
        payload: subjectFormToPayload(values),
      });
      onOpenChange(false);
    } catch {
      /* keep open */
    }
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-h-[min(90vh,720px)] overflow-y-auto sm:max-w-lg">
        <DialogHeader>
          <DialogTitle>Edit subject</DialogTitle>
        </DialogHeader>
        {subject ? (
          <SubjectForm
            id={subject.id}
            defaultValues={{
              name: subject.name,
              code: subject.code ?? "",
              description: subject.description ?? "",
            }}
            submitLabel="Save changes"
            onSubmit={handleSubmit}
          />
        ) : null}
      </DialogContent>
    </Dialog>
  );
}
