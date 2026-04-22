import { Loader2 } from "lucide-react";

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

type DeleteStudentDialogProps = {
  student: Student | null;
  open: boolean;
  onOpenChange: (open: boolean) => void;
  onConfirm: () => void | Promise<void>;
  loading?: boolean;
};

export function DeleteStudentDialog({
  student,
  open,
  onOpenChange,
  onConfirm,
  loading,
}: DeleteStudentDialogProps) {
  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent hideClose className="sm:max-w-md">
        <DialogHeader>
          <DialogTitle>Delete student?</DialogTitle>
          <DialogDescription>
            This will permanently remove{" "}
            <span className="font-medium text-foreground">
              {student?.fullName ?? "this student"}
            </span>{" "}
            from the directory. This action cannot be undone.
          </DialogDescription>
        </DialogHeader>
        <DialogFooter className="gap-2 sm:gap-0">
          <Button
            type="button"
            variant="outline"
            onClick={() => onOpenChange(false)}
            disabled={loading}
          >
            Cancel
          </Button>
          <Button
            type="button"
            variant="destructive"
            disabled={loading || !student}
            onClick={() => void onConfirm()}
          >
            {loading ? (
              <>
                <Loader2 className="h-4 w-4 animate-spin" />
                Deleting…
              </>
            ) : (
              "Delete"
            )}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
