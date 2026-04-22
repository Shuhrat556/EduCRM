import type { AiAnalyticsSnapshot } from "@/features/ai-analytics/model/types";

/** Rich sample payload for layout review and `VITE_AI_ANALYTICS_DEMO`. */
export const mockAiAnalyticsSnapshot: AiAnalyticsSnapshot = {
  generatedAt: new Date().toISOString(),
  source: "demo",
  adminKpis: [
    { label: "Students flagged", value: "12", hint: "Vs. last week · model estimate" },
    { label: "Avg. attendance (at-risk)", value: "78%", variant: "warning" },
    { label: "Open AI actions", value: "5", hint: "Awaiting admin review" },
    { label: "Teachers coached", value: "3", hint: "Recommendations sent" },
  ],
  highlightBullets: [
    "Payment delays cluster in Grade 10 cohort — consider a targeted reminder campaign.",
    "Two teachers show repeat low engagement scores in afternoon blocks.",
    "Model confidence is higher for attendance-based signals than for grade inference.",
  ],
  overdueStudents: {
    headline: "Tuition and fee risk — overdue or partial payers the model is watching.",
    items: [
      {
        id: "s1",
        name: "Aisha Karimova",
        summary: "Balance 60 days+; historically pays after second reminder.",
        riskScore: 72,
        tags: ["payments", "at_risk", "follow_up"],
      },
      {
        id: "s2",
        name: "Jonas Meyer",
        summary: "Partial payment last period; attendance dipped in parallel.",
        riskScore: 65,
        tags: ["payments", "attendance", "at_risk"],
      },
      {
        id: "s3",
        name: "Sofia Ivanova",
        summary: "New enroll; first invoice unpaid past due — low data confidence.",
        riskScore: 54,
        tags: ["payments", "engagement"],
      },
    ],
  },
  weakStudents: {
    headline: "Academic / engagement weakness — not yet accounting for IEP context.",
    items: [
      {
        id: "s4",
        name: "Leo Park",
        summary: "Below cohort median on formative scores; attendance stable.",
        riskScore: 58,
        tags: ["grades", "engagement"],
      },
      {
        id: "s5",
        name: "Maya Ndlovu",
        summary: "Spike in absences correlated with missing assignments.",
        riskScore: 61,
        tags: ["attendance", "grades", "at_risk"],
      },
    ],
  },
  teacherRecommendations: {
    headline: "Operational nudges generated from timetables and cohort outcomes.",
    items: [
      {
        teacherId: "t1",
        teacherName: "Elena Vogel",
        recommendation:
          "Try shorter checks-for-understanding in second half of double periods — engagement scores trail morning blocks.",
        priority: "medium",
        tags: ["engagement", "follow_up"],
      },
      {
        teacherId: "t2",
        teacherName: "Marcus Chen",
        recommendation:
          "Coordinate with homeroom on two students with overlapping attendance dips (see alerts).",
        priority: "high",
        tags: ["attendance", "behavior"],
      },
    ],
  },
  alerts: {
    headline: "High-signal warnings — always verify before outreach.",
    items: [
      {
        id: "a1",
        title: "Possible disengagement pattern",
        body: "Four consecutive absences from lab sessions without excuse notes on file.",
        severity: "high",
        studentName: "Maya Ndlovu",
        tags: ["attendance", "behavior"],
      },
      {
        id: "a2",
        title: "Fee escalation risk",
        body: "Second straight month partial pay combined with new negative attendance trend.",
        severity: "critical",
        studentName: "Jonas Meyer",
        tags: ["payments", "at_risk"],
      },
    ],
  },
};

export function emptyAiAnalyticsSnapshot(
  message?: string,
): AiAnalyticsSnapshot {
  const now = new Date().toISOString();
  return {
    generatedAt: now,
    source: "empty",
    sourceDetail: message,
    adminKpis: [],
    highlightBullets: [],
    overdueStudents: { headline: "", items: [] },
    weakStudents: { headline: "", items: [] },
    teacherRecommendations: { headline: "", items: [] },
    alerts: { headline: "", items: [] },
  };
}
