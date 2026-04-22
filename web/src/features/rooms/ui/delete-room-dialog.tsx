import { Loader2 } from "lucide-react";

import type { Room } from "@/features/rooms/model/types";
import { Button } from "@/shared/ui/components/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/shared/ui/components/dialog";

type DeleteRoomDialogProps = {
  room: Room | null;
  open: boolean;
  onOpenChange: (open: boolean) => void;
  onConfirm: () => void | Promise<void>;
  loading?: boolean;
};

export function DeleteRoomDialog({
  room,
  open,
  onOpenChange,
  onConfirm,
  loading,
}: DeleteRoomDialogProps) {
  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent hideClose className="sm:max-w-md">
        <DialogHeader>
          <DialogTitle>Delete room?</DialogTitle>
          <DialogDescription>
            Permanently remove{" "}
            <span className="font-medium text-foreground">
              {room?.name ?? "this room"}
            </span>
            {room ? (
              <>
                {" "}
                <span className="text-muted-foreground">
                  (capacity {room.capacity})
                </span>
              </>
            ) : null}
            . This cannot be undone.
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
            disabled={loading || !room}
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
