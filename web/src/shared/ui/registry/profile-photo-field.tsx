import { ImagePlus, Trash2 } from "lucide-react";
import { useEffect, useId, useRef, useState } from "react";

import { Avatar, AvatarFallback, AvatarImage } from "@/shared/ui/components/avatar";
import { Button } from "@/shared/ui/components/button";

const MAX_FILE_BYTES = 1_800_000;

type ProfilePhotoFieldProps = {
  /** Used for avatar initials when there is no image. */
  initialsFrom?: string;
  /** Existing remote or data URL (edit mode). */
  existingUrl?: string | null;
  /** Override preview (e.g. object URL from new file). */
  previewUrl?: string | null;
  onFileChange: (file: File | null) => void;
  onReset?: () => void;
  disabled?: boolean;
  helperText?: string;
  error?: string;
};

function initialsFromName(name?: string) {
  if (!name?.trim()) return "?";
  return name
    .split(/\s+/)
    .map((p) => p[0])
    .join("")
    .slice(0, 2)
    .toUpperCase();
}

/**
 * File picker + preview for profile photos. Pair with your API upload mutation.
 */
export function ProfilePhotoField({
  initialsFrom,
  existingUrl,
  previewUrl,
  onFileChange,
  onReset,
  disabled,
  helperText = "JPG, PNG or WebP. Max ~1.8 MB for demo storage.",
  error,
}: ProfilePhotoFieldProps) {
  const inputId = useId();
  const inputRef = useRef<HTMLInputElement>(null);
  const [objectUrl, setObjectUrl] = useState<string | null>(null);

  useEffect(() => {
    return () => {
      if (objectUrl) URL.revokeObjectURL(objectUrl);
    };
  }, [objectUrl]);

  const display =
    previewUrl ??
    objectUrl ??
    existingUrl ??
    null;

  function handleFiles(files: FileList | null) {
    const file = files?.[0];
    if (!file) return;
    if (!file.type.startsWith("image/")) {
      onFileChange(null);
      return;
    }
    if (file.size > MAX_FILE_BYTES) {
      onFileChange(null);
      return;
    }
    if (objectUrl) URL.revokeObjectURL(objectUrl);
    const url = URL.createObjectURL(file);
    setObjectUrl(url);
    onFileChange(file);
  }

  function clear() {
    if (objectUrl) URL.revokeObjectURL(objectUrl);
    setObjectUrl(null);
    onFileChange(null);
    if (inputRef.current) inputRef.current.value = "";
    onReset?.();
  }

  return (
    <div className="space-y-3">
      <div className="flex flex-col items-start gap-4 sm:flex-row sm:items-center">
        <Avatar className="h-24 w-24 border-2 border-dashed border-border">
          {display ? (
            <AvatarImage src={display} alt="" className="object-cover" />
          ) : null}
          <AvatarFallback className="text-lg text-muted-foreground">
            {initialsFromName(initialsFrom)}
          </AvatarFallback>
        </Avatar>
        <div className="flex flex-wrap gap-2">
          <input
            ref={inputRef}
            id={inputId}
            type="file"
            accept="image/jpeg,image/png,image/webp,image/gif"
            className="sr-only"
            disabled={disabled}
            onChange={(e) => {
              handleFiles(e.target.files);
              e.target.value = "";
            }}
          />
          <Button
            type="button"
            variant="outline"
            size="sm"
            className="gap-2"
            disabled={disabled}
            onClick={() => inputRef.current?.click()}
          >
            <ImagePlus className="h-4 w-4" />
            Upload photo
          </Button>
          {(display && (previewUrl || objectUrl || existingUrl)) ? (
            <Button
              type="button"
              variant="ghost"
              size="sm"
              className="gap-2 text-destructive hover:text-destructive"
              disabled={disabled}
              onClick={clear}
            >
              <Trash2 className="h-4 w-4" />
              Remove
            </Button>
          ) : null}
        </div>
      </div>
      {helperText ? (
        <p className="text-xs text-muted-foreground">{helperText}</p>
      ) : null}
      {error ? (
        <p className="text-xs font-medium text-destructive">{error}</p>
      ) : null}
    </div>
  );
}
