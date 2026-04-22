import { useState } from "react";

import { useCreateGroupMutation } from "@/features/groups/hooks/use-groups";
import {
  groupFormToCreatePayload,
  type GroupFormValues,
} from "@/features/groups/model/group-schema";
import { GroupForm } from "@/features/groups/ui/group-form";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "@/shared/ui/components/dialog";

type GroupCreateDialogProps = {
  open: boolean;
  onOpenChange: (open: boolean) => void;
};

export function GroupCreateDialog({
  open,
  onOpenChange,
}: GroupCreateDialogProps) {
  const [formKey, setFormKey] = useState(0);
  const create = useCreateGroupMutation();

  async function handleSubmit(values: GroupFormValues) {
    try {
      await create.mutateAsync(groupFormToCreatePayload(values));
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
      <DialogContent className="max-h-[min(90vh,820px)] overflow-y-auto sm:max-w-lg">
        <DialogHeader>
          <DialogTitle>Create group</DialogTitle>
        </DialogHeader>
        <GroupForm
          key={formKey}
          submitLabel="Create group"
          onSubmit={handleSubmit}
        />
      </DialogContent>
    </Dialog>
  );
}
