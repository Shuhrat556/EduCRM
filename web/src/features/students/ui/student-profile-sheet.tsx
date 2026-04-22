import { Mail, Phone, Pencil, UserPlus } from "lucide-react";

import { useStudentQuery } from "@/features/students/hooks/use-students";
import { studentStatusLabel } from "@/features/students/lib/student-status";
import type { Student } from "@/features/students/model/types";
import { Avatar, AvatarFallback, AvatarImage } from "@/shared/ui/components/avatar";
import { Badge } from "@/shared/ui/components/badge";
import { Button } from "@/shared/ui/components/button";
import { ScrollArea } from "@/shared/ui/components/scroll-area";
import { Separator } from "@/shared/ui/components/separator";
import {
  Sheet,
  SheetContent,
  SheetHeader,
  SheetTitle,
} from "@/shared/ui/components/sheet";
import { Skeleton } from "@/shared/ui/components/skeleton";
import {
  Table,
  TableBody,
  TableCell,
  TableRow,
} from "@/shared/ui/components/table";
import { initialsFromName } from "@/shared/lib/format/initials";
import { QueryErrorAlert } from "@/shared/ui/feedback/query-error-alert";

type StudentProfileSheetProps = {
  studentId: string | null;
  open: boolean;
  onOpenChange: (open: boolean) => void;
  onEdit: (student: Student) => void;
  onAssignGroup: (student: Student) => void;
};

export function StudentProfileSheet({
  studentId,
  open,
  onOpenChange,
  onEdit,
  onAssignGroup,
}: StudentProfileSheetProps) {
  const { data: student, isLoading, isError, error, refetch } = useStudentQuery(
    open && studentId ? studentId : null,
  );

  return (
    <Sheet open={open} onOpenChange={onOpenChange}>
      <SheetContent className="flex w-full flex-col sm:max-w-md">
        <SheetHeader className="space-y-1 text-left">
          <SheetTitle>Student profile</SheetTitle>
        </SheetHeader>
        <ScrollArea className="mt-4 flex-1 pr-3">
          {isLoading ? (
            <div className="space-y-4">
              <div className="flex items-center gap-4">
                <Skeleton className="h-16 w-16 rounded-full" />
                <div className="space-y-2">
                  <Skeleton className="h-6 w-40" />
                  <Skeleton className="h-4 w-24" />
                </div>
              </div>
              <Skeleton className="h-32 w-full" />
            </div>
          ) : isError ? (
            <QueryErrorAlert
              error={error}
              onRetry={() => void refetch()}
              className="text-left"
            />
          ) : student ? (
            <div className="space-y-6">
              <div className="flex flex-col items-center gap-3 sm:flex-row sm:items-start">
                <Avatar className="h-20 w-20 border-2 border-border">
                  {student.photoUrl ? (
                    <AvatarImage src={student.photoUrl} alt="" />
                  ) : null}
                  <AvatarFallback className="text-lg">
                    {initialsFromName(student.fullName)}
                  </AvatarFallback>
                </Avatar>
                <div className="min-w-0 flex-1 space-y-1 text-center sm:text-left">
                  <h3 className="text-lg font-semibold leading-tight">
                    {student.fullName}
                  </h3>
                  <div className="flex flex-wrap items-center justify-center gap-2 sm:justify-start">
                    <Badge variant="secondary">
                      {studentStatusLabel(student.status)}
                    </Badge>
                    {student.groupName ? (
                      <Badge variant="outline">{student.groupName}</Badge>
                    ) : (
                      <span className="text-xs text-muted-foreground">
                        No group
                      </span>
                    )}
                  </div>
                </div>
              </div>

              <Separator />

              <Table>
                <TableBody>
                  <TableRow className="border-0 hover:bg-transparent">
                    <TableCell className="w-28 px-0 py-2 text-muted-foreground">
                      <span className="flex items-center gap-1.5 text-xs font-medium uppercase tracking-wide">
                        <Phone className="h-3.5 w-3.5" />
                        Phone
                      </span>
                    </TableCell>
                    <TableCell className="py-2 font-medium">
                      <a
                        href={`tel:${student.phone}`}
                        className="text-primary hover:underline"
                      >
                        {student.phone}
                      </a>
                    </TableCell>
                  </TableRow>
                  <TableRow className="border-0 hover:bg-transparent">
                    <TableCell className="w-28 px-0 py-2 text-muted-foreground">
                      <span className="flex items-center gap-1.5 text-xs font-medium uppercase tracking-wide">
                        <Mail className="h-3.5 w-3.5" />
                        Email
                      </span>
                    </TableCell>
                    <TableCell className="py-2 font-medium break-all">
                      <a
                        href={`mailto:${student.email}`}
                        className="text-primary hover:underline"
                      >
                        {student.email}
                      </a>
                    </TableCell>
                  </TableRow>
                  <TableRow className="border-0 hover:bg-transparent">
                    <TableCell className="px-0 py-2 text-muted-foreground">
                      <span className="text-xs font-medium uppercase tracking-wide">
                        Group
                      </span>
                    </TableCell>
                    <TableCell className="py-2">
                      {student.groupName ?? "—"}
                    </TableCell>
                  </TableRow>
                  <TableRow className="border-0 hover:bg-transparent">
                    <TableCell className="px-0 py-2 text-muted-foreground">
                      <span className="text-xs font-medium uppercase tracking-wide">
                        Photo URL
                      </span>
                    </TableCell>
                    <TableCell className="max-w-[200px] truncate py-2 text-sm text-muted-foreground">
                      {student.photoUrl ?? "—"}
                    </TableCell>
                  </TableRow>
                </TableBody>
              </Table>

              <div className="flex flex-wrap gap-2 pb-6">
                <Button
                  type="button"
                  variant="default"
                  size="sm"
                  className="gap-1.5"
                  onClick={() => onEdit(student)}
                >
                  <Pencil className="h-3.5 w-3.5" />
                  Edit
                </Button>
                <Button
                  type="button"
                  variant="outline"
                  size="sm"
                  className="gap-1.5"
                  onClick={() => onAssignGroup(student)}
                >
                  <UserPlus className="h-3.5 w-3.5" />
                  Assign group
                </Button>
              </div>
            </div>
          ) : (
            <p className="text-sm text-muted-foreground">Student not found.</p>
          )}
        </ScrollArea>
      </SheetContent>
    </Sheet>
  );
}
