import { usePermissions } from "@/modules/auth/hooks/usePermissions";
import { AppEmpty } from "@/shared/ui/feedback/app-empty";
import { AppError } from "@/shared/ui/feedback/app-error";
import { PageHeader } from "@/shared/ui/layout/page-header";
import { Badge } from "@/shared/ui/components/badge";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/shared/ui/components/card";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/shared/ui/components/table";
import { BookOpen } from "lucide-react";

const DEMO_GRADES = [
  { course: "Mathematics", term: "Spring 2026", grade: "A-", credits: 4 },
  { course: "English", term: "Spring 2026", grade: "B+", credits: 3 },
  { course: "Physics Lab", term: "Spring 2026", grade: "A", credits: 2 },
];

export function StudentGradesPage() {
  const { can } = usePermissions();

  if (!can("student:view_own_grades")) {
    return (
      <AppError
        title="Access restricted"
        message="You do not have permission to view grades."
      />
    );
  }

  return (
    <div className="space-y-8">
      <PageHeader
        title="My grades"
        description="Official grades for your enrolled courses. Data is illustrative until your SIS is connected."
      />

      {DEMO_GRADES.length === 0 ? (
        <AppEmpty
          icon={BookOpen}
          title="No grades yet"
          description="When grading sync is enabled, your courses will appear here."
        />
      ) : (
        <Card className="border-border/80">
          <CardHeader>
            <CardTitle className="text-base">Current term</CardTitle>
            <CardDescription>
              Endpoints such as{" "}
              <code className="rounded bg-muted px-1 py-0.5 text-xs">
                GET /students/me/grades
              </code>{" "}
              can replace this sample table.
            </CardDescription>
          </CardHeader>
          <CardContent className="overflow-x-auto">
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Course</TableHead>
                  <TableHead>Term</TableHead>
                  <TableHead className="text-right">Credits</TableHead>
                  <TableHead className="text-right">Grade</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {DEMO_GRADES.map((row) => (
                  <TableRow key={row.course}>
                    <TableCell className="font-medium">{row.course}</TableCell>
                    <TableCell className="text-muted-foreground">
                      {row.term}
                    </TableCell>
                    <TableCell className="text-right tabular-nums">
                      {row.credits}
                    </TableCell>
                    <TableCell className="text-right">
                      <Badge variant="secondary">{row.grade}</Badge>
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </CardContent>
        </Card>
      )}
    </div>
  );
}
