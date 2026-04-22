/** Matches typical AI / LLM output tags for filtering & future chips. */
export type AiAnalyticsTag =
  | "payments"
  | "attendance"
  | "grades"
  | "behavior"
  | "engagement"
  | "at_risk"
  | "follow_up";

export type AiSeverity = "low" | "medium" | "high" | "critical";

export type AiStudentBrief = {
  id: string;
  name: string;
  /** Short model rationale (1–2 lines). */
  summary: string;
  /** Normalized concern score 0–100 if the model supplies one. */
  riskScore?: number;
  tags: AiAnalyticsTag[];
};

export type AiTeacherRecommendation = {
  teacherId: string;
  teacherName: string;
  recommendation: string;
  priority: AiSeverity;
  tags: AiAnalyticsTag[];
};

export type AiStudentAlert = {
  id: string;
  title: string;
  body: string;
  severity: AiSeverity;
  studentId?: string;
  studentName?: string;
  tags: AiAnalyticsTag[];
};

export type AiAdminKpi = {
  label: string;
  value: string;
  hint?: string;
  variant?: "default" | "warning";
};

/**
 * Expected JSON shape from `GET /ai/analytics/snapshot` (adjust to your backend).
 */
export type AiAnalyticsSnapshot = {
  generatedAt: string;
  summaryMarkdown?: string;
  adminKpis: AiAdminKpi[];
  highlightBullets: string[];
  overdueStudents: {
    headline: string;
    items: AiStudentBrief[];
  };
  weakStudents: {
    headline: string;
    items: AiStudentBrief[];
  };
  teacherRecommendations: {
    headline: string;
    items: AiTeacherRecommendation[];
  };
  alerts: {
    headline: string;
    items: AiStudentAlert[];
  };
  /** Set by the client when the API is unavailable or in demo mode. */
  source?: "api" | "demo" | "empty";
  sourceDetail?: string;
};
