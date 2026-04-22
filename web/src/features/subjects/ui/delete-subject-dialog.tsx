import { Loader2 } from "lucide-react";

import type { Subject } from "@/features/subjects/model/types";
import { Button } from "@/shared/ui/components/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/shared/ui/components/dialog";

type DeleteSubjectDialogProps = {
  subject: Subject | null;
  open: boolean;
  onOpenChange: (open: boolean) => void;
  onConfirm: () => void | Promise<void>;
  loading?: boolean;
};

export function DeleteSubjectDialog({
  subject,
  open,
  onOpenChange,
  onConfirm,
  loading,
}: DeleteSubjectDialogProps) {
  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent hideClose className="sm:max-w-md">
        <DialogHeader>
          <DialogTitle>Delete subject?</DialogTitle>
          <DialogDescription>
            This removes{" "}
            <span className="font-medium text-foreground">
              {subject?.name ?? "this subject"}
            </span>{" "}
            from the catalog. Teachers linked to it may need to be updated.
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
            disabled={loading || !subject}
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
