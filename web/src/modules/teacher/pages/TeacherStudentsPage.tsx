import { usePermissions } from "@/modules/auth/hooks/usePermissions";
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

const ASSIGNED = [
  { id: "1", name: "Alex Morgan", group: "Grade 10A", status: "active" },
  { id: "2", name: "Jamie Lee", group: "Grade 10A", status: "active" },
  { id: "3", name: "Sam Rivera", group: "Grade 10B", status: "at_risk" },
];

export function TeacherStudentsPage() {
  const { can } = usePermissions();

  return (
    <div className="space-y-8">
      <PageHeader
        title="Assigned students"
        description="Students in your cohorts. Only these learners appear for grading and attendance."
      />

      {!can("teacher:view_assigned_students") ? (
        <p className="text-sm text-destructive">
          Your account cannot view assigned students.
        </p>
      ) : (
        <Card className="border-border/80">
          <CardHeader>
            <CardTitle className="text-base">Roster</CardTitle>
            <CardDescription>
              Sample data — connect{" "}
              <code className="rounded bg-muted px-1 py-0.5 text-xs">
                GET /teachers/me/students
              </code>{" "}
              for production.
            </CardDescription>
          </CardHeader>
          <CardContent className="overflow-x-auto">
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Name</TableHead>
                  <TableHead>Group</TableHead>
                  <TableHead className="text-right">Status</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {ASSIGNED.map((row) => (
                  <TableRow key={row.id}>
                    <TableCell className="font-medium">{row.name}</TableCell>
                    <TableCell className="text-muted-foreground">
                      {row.group}
                    </TableCell>
                    <TableCell className="text-right">
                      <Badge
                        variant={
                          row.status === "at_risk" ? "destructive" : "secondary"
                        }
                      >
                        {row.status.replace("_", " ")}
                      </Badge>
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
