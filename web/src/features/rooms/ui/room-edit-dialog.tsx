import type { Room } from "@/features/rooms/model/types";
import { useUpdateRoomMutation } from "@/features/rooms/hooks/use-rooms";
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

type RoomEditDialogProps = {
  room: Room | null;
  open: boolean;
  onOpenChange: (open: boolean) => void;
};

export function RoomEditDialog({
  room,
  open,
  onOpenChange,
}: RoomEditDialogProps) {
  const update = useUpdateRoomMutation();

  async function handleSubmit(values: RoomFormValues) {
    if (!room) return;
    try {
      await update.mutateAsync({
        id: room.id,
        payload: roomFormToPayload(values),
      });
      onOpenChange(false);
    } catch {
      /* keep open */
    }
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-h-[min(90vh,720px)] overflow-y-auto sm:max-w-lg">
        <DialogHeader>
          <DialogTitle>Edit room</DialogTitle>
        </DialogHeader>
        {room ? (
          <RoomForm
            id={room.id}
            defaultValues={{
              name: room.name,
              building: room.building ?? "",
              capacity: room.capacity,
              notes: room.notes ?? "",
            }}
            submitLabel="Save changes"
            onSubmit={handleSubmit}
          />
        ) : null}
      </DialogContent>
    </Dialog>
  );
}
