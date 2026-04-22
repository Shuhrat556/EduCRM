import { Loader2 } from "lucide-react";

import type { ScheduleLesson } from "@/features/schedule/model/types";
import { Button } from "@/shared/ui/components/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/shared/ui/components/dialog";

type DeleteLessonDialogProps = {
  lesson: ScheduleLesson | null;
  open: boolean;
  onOpenChange: (open: boolean) => void;
  onConfirm: () => void | Promise<void>;
  loading?: boolean;
};

export function DeleteLessonDialog({
  lesson,
  open,
  onOpenChange,
  onConfirm,
  loading,
}: DeleteLessonDialogProps) {
  const label =
    lesson?.title?.trim() ||
    (lesson
      ? `${lesson.startTime}–${lesson.endTime} · ${lesson.group.name}`
      : "this slot");

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent hideClose className="sm:max-w-md">
        <DialogHeader>
          <DialogTitle>Delete lesson?</DialogTitle>
          <DialogDescription>
            Remove{" "}
            <span className="font-medium text-foreground">{label}</span> from
            the weekly schedule. This cannot be undone.
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
            disabled={loading || !lesson}
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
