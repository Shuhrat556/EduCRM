import { useState } from "react";

import { useCreateSubjectMutation } from "@/features/subjects/hooks/use-subjects";
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

type SubjectCreateDialogProps = {
  open: boolean;
  onOpenChange: (open: boolean) => void;
};

export function SubjectCreateDialog({
  open,
  onOpenChange,
}: SubjectCreateDialogProps) {
  const [formKey, setFormKey] = useState(0);
  const create = useCreateSubjectMutation();

  async function handleSubmit(values: SubjectFormValues) {
    try {
      await create.mutateAsync(subjectFormToPayload(values));
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
      <DialogContent className="max-h-[min(90vh,720px)] overflow-y-auto sm:max-w-lg">
        <DialogHeader>
          <DialogTitle>Add subject</DialogTitle>
        </DialogHeader>
        <SubjectForm
          key={formKey}
          submitLabel="Create subject"
          onSubmit={handleSubmit}
        />
      </DialogContent>
    </Dialog>
  );
}
