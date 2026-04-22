import { useEffect, useState } from "react";
import { Loader2 } from "lucide-react";

import {
  useGroupOptionsQuery,
  useUpdateStudentMutation,
} from "@/features/students/hooks/use-students";
import type { Student } from "@/features/students/model/types";
import { Button } from "@/shared/ui/components/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/shared/ui/components/dialog";
import { Label } from "@/shared/ui/components/label";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/shared/ui/components/select";

const NONE = "__none__";

type AssignGroupDialogProps = {
  student: Student | null;
  open: boolean;
  onOpenChange: (open: boolean) => void;
};

export function AssignGroupDialog({
  student,
  open,
  onOpenChange,
}: AssignGroupDialogProps) {
  const { data: groups = [], isLoading } = useGroupOptionsQuery();
  const update = useUpdateStudentMutation();
  const [value, setValue] = useState<string>(NONE);

  useEffect(() => {
    if (student) setValue(student.groupId ?? NONE);
  }, [student, open]);

  async function save() {
    if (!student) return;
    await update.mutateAsync({
      id: student.id,
      payload: { groupId: value === NONE ? null : value },
    });
    onOpenChange(false);
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-md">
        <DialogHeader>
          <DialogTitle>Assign to group</DialogTitle>
          <DialogDescription>
            Choose a class or cohort for{" "}
            <span className="font-medium text-foreground">
              {student?.fullName ?? "this student"}
            </span>
            .
          </DialogDescription>
        </DialogHeader>
        <div className="space-y-2 py-2">
          <Label htmlFor="assign-group">Group</Label>
          <Select
            disabled={isLoading}
            value={value}
            onValueChange={setValue}
          >
            <SelectTrigger id="assign-group">
              <SelectValue placeholder="Select group" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value={NONE}>No group</SelectItem>
              {groups.map((g) => (
                <SelectItem key={g.id} value={g.id}>
                  {g.name}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        </div>
        <DialogFooter>
          <Button
            type="button"
            variant="outline"
            onClick={() => onOpenChange(false)}
          >
            Cancel
          </Button>
          <Button
            type="button"
            disabled={update.isPending || !student}
            onClick={() => void save()}
          >
            {update.isPending ? (
              <>
                <Loader2 className="h-4 w-4 animate-spin" />
                Saving…
              </>
            ) : (
              "Save"
            )}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
