import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/shared/ui/components/card";

export function RoleDemoPage() {
  return (
    <Card>
      <CardHeader>
        <CardTitle className="text-base">Role-protected route</CardTitle>
        <CardDescription>
          You reached this page because your user has{" "}
          <code className="rounded bg-muted px-1 py-0.5 text-xs">admin</code> or{" "}
          <code className="rounded bg-muted px-1 py-0.5 text-xs">
            super_admin
          </code>
          .
        </CardDescription>
      </CardHeader>
      <CardContent className="text-sm text-muted-foreground">
        Replace this placeholder with a real feature entry point when you add
        business screens.
      </CardContent>
    </Card>
  );
}
