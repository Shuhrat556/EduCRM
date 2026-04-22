/// <reference types="vite/client" />

interface ImportMetaEnv {
  readonly VITE_API_BASE_URL: string;
  readonly VITE_APP_NAME?: string;
  readonly VITE_ENABLE_QUERY_DEVTOOLS?: string;
  /** Set to "true" to load rich sample dashboard data without a live API. */
  readonly VITE_DASHBOARD_DEMO?: string;
  /** Persisted mock CRUD for students when API is not available. */
  readonly VITE_STUDENTS_DEMO?: string;
  /** Persisted mock CRUD for teachers when API is not available. */
  readonly VITE_TEACHERS_DEMO?: string;
  /** Persisted mock CRUD for subjects when API is not available. */
  readonly VITE_SUBJECTS_DEMO?: string;
  /** Persisted mock CRUD for rooms when API is not available. */
  readonly VITE_ROOMS_DEMO?: string;
  /** Persisted mock CRUD for groups when API is not available. */
  readonly VITE_GROUPS_DEMO?: string;
  /** Persisted mock weekly lessons when API is not available. */
  readonly VITE_SCHEDULE_DEMO?: string;
  /** Persisted lesson attendance & grades when API is not available. */
  readonly VITE_ATTENDANCE_DEMO?: string;
  /** Mock tuition ledger / receivables when API is not available. */
  readonly VITE_PAYMENTS_DEMO?: string;
  /** Sample AI analytics payload for the admin AI Insights page. */
  readonly VITE_AI_ANALYTICS_DEMO?: string;
}

interface ImportMeta {
  readonly env: ImportMetaEnv;
}
