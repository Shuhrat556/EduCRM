import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/shared/ui/components/card";

export function StudentHomePage() {
  return (
    <Card>
      <CardHeader>
        <CardTitle className="text-base">Student workspace</CardTitle>
        <CardDescription>
          Lightweight area for learner-facing tools. Use the sidebar for
          Schedule, Attendance, Payments, and Settings.
        </CardDescription>
      </CardHeader>
      <CardContent className="text-sm text-muted-foreground">
        Extend this page with grades or assignments when your product scope
        includes a student portal.
      </CardContent>
    </Card>
  );
}
