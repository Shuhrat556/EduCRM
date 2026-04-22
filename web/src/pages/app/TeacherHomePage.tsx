import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/shared/ui/components/card";

export function TeacherHomePage() {
  return (
    <Card>
      <CardHeader>
        <CardTitle className="text-base">Teacher workspace</CardTitle>
        <CardDescription>
          Optional landing module for teacher-only flows (e.g. class shortcuts).
          The main sidebar already exposes Schedule, Attendance, Groups, and
          Reports.
        </CardDescription>
      </CardHeader>
      <CardContent className="text-sm text-muted-foreground">
        Navigate using the sidebar, or bookmark this path for future teacher
        dashboards.
      </CardContent>
    </Card>
  );
}
