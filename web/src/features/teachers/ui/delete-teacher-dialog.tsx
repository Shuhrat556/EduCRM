import { Loader2 } from "lucide-react";

import type { Teacher } from "@/features/teachers/model/types";
import { Button } from "@/shared/ui/components/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/shared/ui/components/dialog";

type DeleteTeacherDialogProps = {
  teacher: Teacher | null;
  open: boolean;
  onOpenChange: (open: boolean) => void;
  onConfirm: () => void | Promise<void>;
  loading?: boolean;
};

export function DeleteTeacherDialog({
  teacher,
  open,
  onOpenChange,
  onConfirm,
  loading,
}: DeleteTeacherDialogProps) {
  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent hideClose className="sm:max-w-md">
        <DialogHeader>
          <DialogTitle>Delete teacher?</DialogTitle>
          <DialogDescription>
            Permanently remove{" "}
            <span className="font-medium text-foreground">
              {teacher?.fullName ?? "this teacher"}
            </span>
            ? This cannot be undone.
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
            disabled={loading || !teacher}
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
