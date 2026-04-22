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
import { Button } from "@/shared/ui/components/button";
import { Shield } from "lucide-react";

const ADMINS = [
  {
    id: "1",
    name: "Jordan Brooks",
    email: "j.brooks@school.org",
    status: "active" as const,
  },
  {
    id: "2",
    name: "Taylor Chen",
    email: "t.chen@school.org",
    status: "invited" as const,
  },
];

export function AdminsPage() {
  const { can } = usePermissions();

  if (!can("super_admin:manage_admins")) {
    return (
      <p className="text-sm text-destructive">
        Only Super Admins can manage administrator accounts.
      </p>
    );
  }

  return (
    <div className="space-y-8">
      <PageHeader
        title="Administrators"
        description="Provision and deactivate institution admins. Hook this table to your directory or admin API."
        actions={
          <Button type="button" size="sm" disabled className="gap-2">
            <Shield className="h-4 w-4" />
            Invite admin
          </Button>
        }
      />

      <Card className="border-border/80">
        <CardHeader>
          <CardTitle className="text-base">Directory</CardTitle>
          <CardDescription>
            Sample rows — replace with{" "}
            <code className="rounded bg-muted px-1 py-0.5 text-xs">
              GET /super-admin/admins
            </code>
            .
          </CardDescription>
        </CardHeader>
        <CardContent className="overflow-x-auto">
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>Name</TableHead>
                <TableHead>Email</TableHead>
                <TableHead className="text-right">Status</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {ADMINS.map((row) => (
                <TableRow key={row.id}>
                  <TableCell className="font-medium">{row.name}</TableCell>
                  <TableCell className="text-muted-foreground">
                    {row.email}
                  </TableCell>
                  <TableCell className="text-right">
                    <Badge
                      variant={
                        row.status === "active" ? "secondary" : "outline"
                      }
                    >
                      {row.status}
                    </Badge>
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </CardContent>
      </Card>
    </div>
  );
}
