import type { Group } from "@/features/groups/model/types";
import { useUpdateGroupMutation } from "@/features/groups/hooks/use-groups";
import {
  GROUP_FORM_NONE,
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

type GroupEditDialogProps = {
  group: Group | null;
  open: boolean;
  onOpenChange: (open: boolean) => void;
};

export function GroupEditDialog({
  group,
  open,
  onOpenChange,
}: GroupEditDialogProps) {
  const update = useUpdateGroupMutation();

  async function handleSubmit(values: GroupFormValues) {
    if (!group) return;
    try {
      await update.mutateAsync({
        id: group.id,
        payload: groupFormToCreatePayload(values),
      });
      onOpenChange(false);
    } catch {
      /* keep open */
    }
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-h-[min(90vh,820px)] overflow-y-auto sm:max-w-lg">
        <DialogHeader>
          <DialogTitle>Edit group</DialogTitle>
        </DialogHeader>
        {group ? (
          <GroupForm
            id={group.id}
            defaultValues={{
              name: group.name,
              teacherId: group.teacherId ?? GROUP_FORM_NONE,
              subjectId: group.subjectId ?? GROUP_FORM_NONE,
              roomId: group.roomId ?? GROUP_FORM_NONE,
              monthlyFee: group.monthlyFee,
              startDate: group.startDate,
              endDate: group.endDate ?? "",
              status: group.status,
            }}
            submitLabel="Save changes"
            onSubmit={handleSubmit}
          />
        ) : null}
      </DialogContent>
    </Dialog>
  );
}
