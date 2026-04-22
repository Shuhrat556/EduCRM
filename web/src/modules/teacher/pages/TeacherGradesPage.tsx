import { usePermissions } from "@/modules/auth/hooks/usePermissions";
import { PageHeader } from "@/shared/ui/layout/page-header";
import { Button } from "@/shared/ui/components/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/shared/ui/components/card";
import { Input } from "@/shared/ui/components/input";
import { Label } from "@/shared/ui/components/label";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/shared/ui/components/table";
import { useState } from "react";

const ROWS = [
  { studentId: "1", name: "Alex Morgan", course: "Mathematics", draft: "" },
  { studentId: "2", name: "Jamie Lee", course: "Mathematics", draft: "" },
];

export function TeacherGradesPage() {
  const { can } = usePermissions();
  const [values, setValues] = useState<Record<string, string>>({});

  if (!can("teacher:assign_grades")) {
    return (
      <p className="text-sm text-destructive">
        You do not have permission to assign grades.
      </p>
    );
  }

  return (
    <div className="space-y-8">
      <PageHeader
        title="Grades"
        description="Enter scores for students assigned to you. Submission calls a batch API when implemented."
      />

      <Card className="border-border/80">
        <CardHeader>
          <CardTitle className="text-base">Assessment draft</CardTitle>
          <CardDescription>
            Values stay local until{" "}
            <code className="rounded bg-muted px-1 py-0.5 text-xs">
              POST /teachers/me/grades
            </code>{" "}
            is wired.
          </CardDescription>
        </CardHeader>
        <CardContent className="space-y-6">
          <div className="overflow-x-auto">
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Student</TableHead>
                  <TableHead>Course</TableHead>
                  <TableHead className="w-[140px]">Score / letter</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {ROWS.map((row) => (
                  <TableRow key={row.studentId}>
                    <TableCell className="font-medium">{row.name}</TableCell>
                    <TableCell className="text-muted-foreground">
                      {row.course}
                    </TableCell>
                    <TableCell>
                      <Label htmlFor={`g-${row.studentId}`} className="sr-only">
                        Grade for {row.name}
                      </Label>
                      <Input
                        id={`g-${row.studentId}`}
                        value={values[row.studentId] ?? ""}
                        onChange={(e) =>
                          setValues((v) => ({
                            ...v,
                            [row.studentId]: e.target.value,
                          }))
                        }
                        placeholder="e.g. 92 or A-"
                        className="h-9"
                      />
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </div>
          <div className="flex flex-wrap gap-2">
            <Button type="button" variant="secondary" disabled>
              Save draft (API pending)
            </Button>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
