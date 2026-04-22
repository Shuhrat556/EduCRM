import { useState } from "react";

import { useCreateRoomMutation } from "@/features/rooms/hooks/use-rooms";
import {
  roomFormToPayload,
  type RoomFormValues,
} from "@/features/rooms/model/room-schema";
import { RoomForm } from "@/features/rooms/ui/room-form";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "@/shared/ui/components/dialog";

type RoomCreateDialogProps = {
  open: boolean;
  onOpenChange: (open: boolean) => void;
};

export function RoomCreateDialog({
  open,
  onOpenChange,
}: RoomCreateDialogProps) {
  const [formKey, setFormKey] = useState(0);
  const create = useCreateRoomMutation();

  async function handleSubmit(values: RoomFormValues) {
    try {
      await create.mutateAsync(roomFormToPayload(values));
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
          <DialogTitle>Add room</DialogTitle>
        </DialogHeader>
        <RoomForm key={formKey} submitLabel="Create room" onSubmit={handleSubmit} />
      </DialogContent>
    </Dialog>
  );
}
