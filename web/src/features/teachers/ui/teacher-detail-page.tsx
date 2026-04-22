import {
  ArrowLeft,
  Loader2,
  Mail,
  Pencil,
  Phone,
  Trash2,
} from "lucide-react";
import { useEffect, useState } from "react";
import { Link, useNavigate, useParams } from "react-router-dom";

import { useGroupOptionsQuery } from "@/features/students/hooks/use-students";
import {
  useDeleteTeacherMutation,
  useSubjectOptionsQuery,
  useTeacherQuery,
  useUpdateTeacherMutation,
  useUploadTeacherPhotoMutation,
} from "@/features/teachers/hooks/use-teachers";
import { teacherStatusLabel } from "@/features/teachers/lib/teacher-status";
import { DeleteTeacherDialog } from "@/features/teachers/ui/delete-teacher-dialog";
import { TeacherEditDialog } from "@/features/teachers/ui/teacher-edit-dialog";
import { EntityMultiPicker } from "@/shared/ui/registry/entity-multi-picker";
import { ProfilePhotoField } from "@/shared/ui/registry/profile-photo-field";
import { Avatar, AvatarFallback, AvatarImage } from "@/shared/ui/components/avatar";
import { Badge } from "@/shared/ui/components/badge";
import { Button } from "@/shared/ui/components/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/shared/ui/components/card";
import { Skeleton } from "@/shared/ui/components/skeleton";
import { initialsFromName } from "@/shared/lib/format/initials";
import { usePortalBase } from "@/shared/lib/hooks/usePortalBase";
import { QueryErrorAlert } from "@/shared/ui/feedback/query-error-alert";

export function TeacherDetailPage() {
  const base = usePortalBase();
  const { teacherId } = useParams<{ teacherId: string }>();
  const navigate = useNavigate();
  const id = teacherId ?? null;

  const { data: teacher, isLoading, isError, error, refetch } =
    useTeacherQuery(id);
  const { data: groups = [], isLoading: groupsLoading } = useGroupOptionsQuery();
  const { data: subjects = [], isLoading: subjectsLoading } =
    useSubjectOptionsQuery();
  const update = useUpdateTeacherMutation();
  const upload = useUploadTeacherPhotoMutation();
  const removeTeacher = useDeleteTeacherMutation();

  const [groupIds, setGroupIds] = useState<string[]>([]);
  const [subjectIds, setSubjectIds] = useState<string[]>([]);
  const [editOpen, setEditOpen] = useState(false);
  const [deleteOpen, setDeleteOpen] = useState(false);

  useEffect(() => {
    if (teacher) {
      setGroupIds([...teacher.groupIds]);
      setSubjectIds([...teacher.subjectIds]);
    }
  }, [teacher]);

  const groupOptions = groups.map((g) => ({ id: g.id, label: g.name }));
  const subjectOptions = subjects.map((s) => ({ id: s.id, label: s.name }));

  async function saveAssignments() {
    if (!teacher) return;
    await update.mutateAsync({
      id: teacher.id,
      payload: { groupIds, subjectIds },
    });
  }

  if (!id) {
    return (
      <p className="text-sm text-muted-foreground">Invalid teacher link.</p>
    );
  }

  return (
    <div className="space-y-6">
      <div className="flex flex-wrap items-center gap-3">
        <Button variant="ghost" size="sm" asChild className="gap-2 px-2">
          <Link to={`${base}/teachers`}>
            <ArrowLeft className="h-4 w-4" />
            Teachers
          </Link>
        </Button>
      </div>

      {isLoading ? (
        <div className="space-y-4">
          <Skeleton className="h-10 w-64" />
          <Skeleton className="h-40 w-full max-w-2xl" />
        </div>
      ) : isError ? (
        <QueryErrorAlert
          error={error}
          title="Could not load teacher"
          onRetry={() => void refetch()}
        />
      ) : !teacher ? (
        <p className="text-sm text-muted-foreground">Teacher not found.</p>
      ) : (
        <>
          <div className="flex flex-col gap-6 lg:flex-row lg:items-start lg:justify-between">
            <div className="flex flex-col gap-4 sm:flex-row sm:items-center">
              <Avatar className="h-24 w-24 border-2 border-border">
                {teacher.photoUrl ? (
                  <AvatarImage src={teacher.photoUrl} alt="" />
                ) : null}
                <AvatarFallback className="text-xl">
                  {initialsFromName(teacher.fullName)}
                </AvatarFallback>
              </Avatar>
              <div className="space-y-2">
                <h1 className="text-2xl font-semibold tracking-tight md:text-3xl">
                  {teacher.fullName}
                </h1>
                <div className="flex flex-wrap gap-2">
                  <Badge variant="secondary">
                    {teacherStatusLabel(teacher.status)}
                  </Badge>
                  {teacher.subjects.map((s) => (
                    <Badge key={s.id} variant="outline">
                      {s.name}
                    </Badge>
                  ))}
                </div>
              </div>
            </div>
            <div className="flex flex-wrap gap-2">
              <Button
                type="button"
                variant="outline"
                className="gap-2"
                onClick={() => setEditOpen(true)}
              >
                <Pencil className="h-4 w-4" />
                Edit
              </Button>
              <Button
                type="button"
                variant="destructive"
                className="gap-2"
                onClick={() => setDeleteOpen(true)}
              >
                <Trash2 className="h-4 w-4" />
                Delete
              </Button>
            </div>
          </div>

          <div className="grid gap-6 lg:grid-cols-2">
            <Card className="border-border/80">
              <CardHeader>
                <CardTitle className="text-base">Contact</CardTitle>
                <CardDescription>Reach this teacher directly.</CardDescription>
              </CardHeader>
              <CardContent className="space-y-3 text-sm">
                <div className="flex items-center gap-2">
                  <Phone className="h-4 w-4 text-muted-foreground" />
                  <a
                    className="font-medium text-primary hover:underline"
                    href={`tel:${teacher.phone}`}
                  >
                    {teacher.phone}
                  </a>
                </div>
                <div className="flex items-center gap-2">
                  <Mail className="h-4 w-4 text-muted-foreground" />
                  <a
                    className="break-all font-medium text-primary hover:underline"
                    href={`mailto:${teacher.email}`}
                  >
                    {teacher.email}
                  </a>
                </div>
              </CardContent>
            </Card>

            <Card className="border-border/80">
              <CardHeader>
                <CardTitle className="text-base">Assignments</CardTitle>
                <CardDescription>
                  Groups and subjects for timetabling and payroll.
                </CardDescription>
              </CardHeader>
              <CardContent className="space-y-4">
                <EntityMultiPicker
                  label="Groups"
                  values={groupIds}
                  onChange={setGroupIds}
                  options={groupOptions}
                  disabled={groupsLoading}
                />
                <EntityMultiPicker
                  label="Subjects"
                  values={subjectIds}
                  onChange={setSubjectIds}
                  options={subjectOptions}
                  disabled={subjectsLoading}
                />
                <Button
                  type="button"
                  className="inline-flex gap-2"
                  disabled={
                    update.isPending ||
                    (JSON.stringify(groupIds) ===
                      JSON.stringify(teacher.groupIds) &&
                      JSON.stringify(subjectIds) ===
                        JSON.stringify(teacher.subjectIds))
                  }
                  onClick={() => void saveAssignments()}
                >
                  {update.isPending ? (
                    <>
                      <Loader2 className="h-4 w-4 animate-spin" />
                      Saving…
                    </>
                  ) : (
                    "Save assignments"
                  )}
                </Button>
              </CardContent>
            </Card>
          </div>

          <Card className="border-border/80">
            <CardHeader>
              <CardTitle className="text-base">Profile photo</CardTitle>
              <CardDescription>
                Upload replaces the existing image immediately (demo stores a
                data URL). Remove clears the photo.
              </CardDescription>
            </CardHeader>
            <CardContent>
              <ProfilePhotoField
                initialsFrom={teacher.fullName}
                existingUrl={teacher.photoUrl}
                onFileChange={(f) => {
                  if (f) void upload.mutateAsync({ id: teacher.id, file: f });
                }}
                onReset={() => {
                  void update.mutateAsync({
                    id: teacher.id,
                    payload: { photoUrl: null },
                  });
                }}
                disabled={upload.isPending || update.isPending}
              />
            </CardContent>
          </Card>

          <TeacherEditDialog
            teacher={teacher}
            open={editOpen}
            onOpenChange={setEditOpen}
          />
          <DeleteTeacherDialog
            teacher={teacher}
            open={deleteOpen}
            onOpenChange={setDeleteOpen}
            loading={removeTeacher.isPending}
            onConfirm={async () => {
              await removeTeacher.mutateAsync(teacher.id);
              setDeleteOpen(false);
              navigate(`${base}/teachers`, { replace: true });
            }}
          />
        </>
      )}
    </div>
  );
}
