import {
  AlertTriangle,
  Brain,
  GraduationCap,
  Lightbulb,
  Sparkles,
  UserRound,
} from "lucide-react";
import { Link } from "react-router-dom";

import { useAiAnalyticsSnapshot } from "@/features/ai-analytics/hooks/use-ai-analytics-snapshot";
import type * as Ai from "@/features/ai-analytics/model/types";
import { StatCard } from "@/features/dashboard/ui/stat-card";
import { Badge } from "@/shared/ui/components/badge";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/shared/ui/components/card";
import { ScrollArea } from "@/shared/ui/components/scroll-area";
import { Skeleton } from "@/shared/ui/components/skeleton";
import { cn } from "@/shared/lib/cn";
import { usePortalBase } from "@/shared/lib/hooks/usePortalBase";
import { QueryErrorAlert } from "@/shared/ui/feedback/query-error-alert";
import { PageHeader } from "@/shared/ui/layout/page-header";

const TAG_LABELS: Record<Ai.AiAnalyticsTag, string> = {
  payments: "Payments",
  attendance: "Attendance",
  grades: "Grades",
  behavior: "Behavior",
  engagement: "Engagement",
  at_risk: "At risk",
  follow_up: "Follow up",
};

function severityBadgeVariant(
  s: Ai.AiSeverity,
): "default" | "secondary" | "destructive" | "outline" {
  switch (s) {
    case "critical":
      return "destructive";
    case "high":
      return "destructive";
    case "medium":
      return "default";
    default:
      return "secondary";
  }
}

function severityLabel(s: Ai.AiSeverity) {
  return s.charAt(0).toUpperCase() + s.slice(1);
}

function TagChips({ tags }: { tags: Ai.AiAnalyticsTag[] }) {
  if (!tags.length) return null;
  return (
    <div className="flex flex-wrap gap-1.5 pt-1">
      {tags.map((t) => (
        <Badge key={t} variant="outline" className="font-normal text-[10px]">
          {TAG_LABELS[t] ?? t}
        </Badge>
      ))}
    </div>
  );
}

function SectionSkeleton({ lines = 3 }: { lines?: number }) {
  return (
    <div className="space-y-3">
      {Array.from({ length: lines }).map((_, i) => (
        <div key={i} className="space-y-2 rounded-lg border border-border/60 p-3">
          <Skeleton className="h-4 w-[200px] max-w-[75%]" />
          <Skeleton className="h-3 w-full" />
          <Skeleton className="h-3 w-4/5" />
        </div>
      ))}
    </div>
  );
}

const KPI_ICONS = [AlertTriangle, GraduationCap, Brain, UserRound] as const;

export function AiAnalyticsPage() {
  const rosterBase = usePortalBase();
  const { data, isPending, isError, error, refetch, isFetching } =
    useAiAnalyticsSnapshot();

  const kpis = data?.adminKpis ?? [];
  const rosterHref =
    rosterBase === "/admin" ? `${rosterBase}/students` : `${rosterBase}/teachers`;

  return (
    <div className="space-y-8">
      <PageHeader
        title="AI Analytics"
        description={
          <>
            Model-generated summaries for admins. Wire{" "}
            <code className="rounded bg-muted px-1 py-0.5 text-xs">
              GET /ai/analytics/snapshot
            </code>{" "}
            to replace placeholders. Always verify before acting on alerts.
          </>
        }
        meta={
          <>
            <Badge variant="secondary" className="gap-1 font-normal">
              <Sparkles className="h-3 w-3" aria-hidden />
              Insights
            </Badge>
            {!isPending && data?.source === "demo" ? (
              <Badge variant="outline" className="font-normal">
                Demo data
              </Badge>
            ) : null}
            {!isPending && data?.source === "empty" ? (
              <Badge variant="outline" className="font-normal">
                Offline preview
              </Badge>
            ) : null}
            {isFetching && !isPending ? (
              <Badge variant="outline" className="font-normal">
                Updating…
              </Badge>
            ) : null}
          </>
        }
      />

      {isError ? (
        <QueryErrorAlert error={error} onRetry={() => void refetch()} />
      ) : null}

      {!isPending && data?.sourceDetail ? (
        <Card className="border-dashed border-primary/30 bg-primary/5">
          <CardContent className="py-3 text-sm text-muted-foreground">
            {data.sourceDetail}
          </CardContent>
        </Card>
      ) : null}

      {/* Admin summary KPIs */}
      <section aria-labelledby="ai-admin-kpis">
        <h2 id="ai-admin-kpis" className="sr-only">
          Admin summary
        </h2>
        <div className="grid gap-4 sm:grid-cols-2 xl:grid-cols-4">
          {isPending
            ? [0, 1, 2, 3].map((i) => (
                <StatCard
                  key={i}
                  title="—"
                  value="—"
                  icon={KPI_ICONS[i % KPI_ICONS.length] ?? Brain}
                  loading
                />
              ))
            : kpis.length === 0
              ? (
                  <div className="col-span-full rounded-lg border border-dashed border-border/80 bg-muted/10 px-4 py-6 text-center text-sm text-muted-foreground">
                    No admin KPIs in the last snapshot. Your API can populate{" "}
                    <code className="rounded bg-muted px-1 py-0.5 text-xs">
                      adminKpis
                    </code>
                    .
                  </div>
                )
              : kpis.map((item, i) => (
                  <StatCard
                    key={`${item.label}-${i}`}
                    title={item.label}
                    value={item.value}
                    subtitle={item.hint}
                    icon={KPI_ICONS[i % KPI_ICONS.length] ?? Brain}
                    variant={item.variant === "warning" ? "warning" : "default"}
                  />
                ))}
        </div>

        <Card className="mt-4 border-border/80">
          <CardHeader className="pb-2">
            <CardTitle className="flex items-center gap-2 text-base font-semibold">
              <Lightbulb className="h-4 w-4 text-primary" aria-hidden />
              Executive highlights
            </CardTitle>
            <CardDescription>
              Bulleted takeaways as returned by the model (empty until the API
              responds).
            </CardDescription>
          </CardHeader>
          <CardContent>
            {isPending ? (
              <div className="space-y-2">
                <Skeleton className="h-4 w-full" />
                <Skeleton className="h-4 w-full" />
                <Skeleton className="h-4 w-4/5" />
              </div>
            ) : (data?.highlightBullets?.length ?? 0) === 0 ? (
              <p className="text-sm text-muted-foreground">
                No highlights yet.
              </p>
            ) : (
              <ul className="list-inside list-disc space-y-2 text-sm text-foreground">
                {data!.highlightBullets.map((b, idx) => (
                  <li key={idx} className="leading-relaxed">
                    {b}
                  </li>
                ))}
              </ul>
            )}
            {!isPending && data?.generatedAt ? (
              <p className="mt-4 text-xs text-muted-foreground">
                Snapshot time:{" "}
                {new Date(data.generatedAt).toLocaleString(undefined, {
                  dateStyle: "medium",
                  timeStyle: "short",
                })}
              </p>
            ) : null}
          </CardContent>
        </Card>
      </section>

      <div className="grid gap-6 xl:grid-cols-2">
        {/* Overdue students */}
        <Card className="min-h-[280px] border-border/80">
          <CardHeader className="pb-2">
            <CardTitle className="flex items-center gap-2 text-base font-semibold">
              <AlertTriangle className="h-4 w-4 text-amber-600 dark:text-amber-400" />
              Overdue students
            </CardTitle>
            <CardDescription>
              {isPending ? (
                <Skeleton className="mt-1 h-3 w-full" />
              ) : (
                data?.overdueStudents.headline || "No summary line from API."
              )}
            </CardDescription>
          </CardHeader>
          <CardContent>
            {isPending ? (
              <SectionSkeleton />
            ) : (data?.overdueStudents.items.length ?? 0) === 0 ? (
              <EmptySection />
            ) : (
              <ScrollArea className="max-h-[340px] pr-3">
                <ul className="space-y-3">
                  {data!.overdueStudents.items.map((s) => (
                    <li
                      key={s.id}
                      className="rounded-lg border border-border/70 bg-muted/20 px-3 py-2.5"
                    >
                      <div className="flex flex-wrap items-baseline justify-between gap-2">
                        <span className="font-medium text-foreground">
                          {s.name}
                        </span>
                        {s.riskScore != null ? (
                          <span className="text-xs tabular-nums text-muted-foreground">
                            Risk {s.riskScore}
                          </span>
                        ) : null}
                      </div>
                      <p className="mt-1 text-sm text-muted-foreground">
                        {s.summary}
                      </p>
                      <TagChips tags={s.tags} />
                    </li>
                  ))}
                </ul>
              </ScrollArea>
            )}
          </CardContent>
        </Card>

        {/* Weak students */}
        <Card className="min-h-[280px] border-border/80">
          <CardHeader className="pb-2">
            <CardTitle className="flex items-center gap-2 text-base font-semibold">
              <GraduationCap className="h-4 w-4 text-primary" />
              Weak students
            </CardTitle>
            <CardDescription>
              {isPending ? (
                <Skeleton className="mt-1 h-3 w-full" />
              ) : (
                data?.weakStudents.headline || "No summary line from API."
              )}
            </CardDescription>
          </CardHeader>
          <CardContent>
            {isPending ? (
              <SectionSkeleton />
            ) : (data?.weakStudents.items.length ?? 0) === 0 ? (
              <EmptySection />
            ) : (
              <ScrollArea className="max-h-[340px] pr-3">
                <ul className="space-y-3">
                  {data!.weakStudents.items.map((s) => (
                    <li
                      key={s.id}
                      className="rounded-lg border border-border/70 bg-muted/20 px-3 py-2.5"
                    >
                      <div className="flex flex-wrap items-baseline justify-between gap-2">
                        <span className="font-medium text-foreground">
                          {s.name}
                        </span>
                        {s.riskScore != null ? (
                          <span className="text-xs tabular-nums text-muted-foreground">
                            Score {s.riskScore}
                          </span>
                        ) : null}
                      </div>
                      <p className="mt-1 text-sm text-muted-foreground">
                        {s.summary}
                      </p>
                      <TagChips tags={s.tags} />
                    </li>
                  ))}
                </ul>
              </ScrollArea>
            )}
          </CardContent>
        </Card>

        {/* Teacher recommendations */}
        <Card className="min-h-[280px] border-border/80">
          <CardHeader className="pb-2">
            <CardTitle className="flex items-center gap-2 text-base font-semibold">
              <UserRound className="h-4 w-4 text-primary" />
              Teacher recommendations
            </CardTitle>
            <CardDescription>
              {isPending ? (
                <Skeleton className="mt-1 h-3 w-full" />
              ) : (
                data?.teacherRecommendations.headline ||
                "No summary line from API."
              )}
            </CardDescription>
          </CardHeader>
          <CardContent>
            {isPending ? (
              <SectionSkeleton lines={2} />
            ) : (data?.teacherRecommendations.items.length ?? 0) === 0 ? (
              <EmptySection />
            ) : (
              <ScrollArea className="max-h-[340px] pr-3">
                <ul className="space-y-3">
                  {data!.teacherRecommendations.items.map((t) => (
                    <li
                      key={t.teacherId}
                      className="rounded-lg border border-border/70 bg-muted/20 px-3 py-2.5"
                    >
                      <div className="flex flex-wrap items-center gap-2">
                        <span className="font-medium text-foreground">
                          {t.teacherName}
                        </span>
                        <Badge
                          variant={severityBadgeVariant(t.priority)}
                          className="text-[10px] font-normal"
                        >
                          {severityLabel(t.priority)}
                        </Badge>
                      </div>
                      <p className="mt-2 text-sm leading-relaxed text-muted-foreground">
                        {t.recommendation}
                      </p>
                      <TagChips tags={t.tags} />
                    </li>
                  ))}
                </ul>
              </ScrollArea>
            )}
          </CardContent>
        </Card>

        {/* Alerts */}
        <Card className="min-h-[280px] border-border/80 border-red-500/15">
          <CardHeader className="pb-2">
            <CardTitle className="flex items-center gap-2 text-base font-semibold">
              <Brain className="h-4 w-4 text-destructive" />
              Warnings &amp; alerts
            </CardTitle>
            <CardDescription>
              {isPending ? (
                <Skeleton className="mt-1 h-3 w-full" />
              ) : (
                data?.alerts.headline || "No summary line from API."
              )}
            </CardDescription>
          </CardHeader>
          <CardContent>
            {isPending ? (
              <SectionSkeleton lines={2} />
            ) : (data?.alerts.items.length ?? 0) === 0 ? (
              <EmptySection />
            ) : (
              <ScrollArea className="max-h-[340px] pr-3">
                <ul className="space-y-3">
                  {data!.alerts.items.map((a) => (
                    <li
                      key={a.id}
                      className={cn(
                        "rounded-lg border px-3 py-2.5",
                        a.severity === "critical" || a.severity === "high"
                          ? "border-destructive/30 bg-destructive/5"
                          : "border-border/70 bg-muted/20",
                      )}
                    >
                      <div className="flex flex-wrap items-start justify-between gap-2">
                        <span className="font-medium text-foreground">
                          {a.title}
                        </span>
                        <Badge
                          variant={severityBadgeVariant(a.severity)}
                          className="shrink-0 text-[10px] font-normal"
                        >
                          {severityLabel(a.severity)}
                        </Badge>
                      </div>
                      {a.studentName ? (
                        <p className="mt-1 text-xs font-medium text-muted-foreground">
                          Student:{" "}
                          <Link
                            to={rosterHref}
                            className="text-primary underline-offset-4 hover:underline"
                          >
                            {a.studentName}
                          </Link>
                        </p>
                      ) : null}
                      <p className="mt-1 text-sm text-muted-foreground">
                        {a.body}
                      </p>
                      <TagChips tags={a.tags} />
                    </li>
                  ))}
                </ul>
              </ScrollArea>
            )}
          </CardContent>
        </Card>
      </div>
    </div>
  );
}

function EmptySection() {
  return (
    <p className="rounded-lg border border-dashed border-border/80 bg-muted/10 py-8 text-center text-sm text-muted-foreground">
      No items in this section yet.
    </p>
  );
}
