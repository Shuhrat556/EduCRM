export type PaymentType = "cash" | "card" | "bank_transfer" | "other";

export type BillingStatus =
  | "paid"
  | "partial"
  | "unpaid"
  | "overdue"
  | "waived"
  | "no_fee";

export interface ReceivableRow {
  studentId: string;
  studentName: string;
  groupId: string | null;
  groupName: string | null;
  periodMonth: string;
  /** Expected tuition for the period (from group). */
  expectedAmount: number;
  /** Sum of payments recorded for this student and period. */
  paidAmount: number;
  /** expectedAmount - paidAmount (0 if waived / no_fee). */
  balance: number;
  /** Last calendar day of the billing month (YYYY-MM-DD). */
  dueDate: string;
  status: BillingStatus;
  tuitionExempt: boolean;
}

export interface ReceivablesResponse {
  periodMonth: string;
  rows: ReceivableRow[];
  summary: {
    totalExpected: number;
    totalPaid: number;
    totalOutstanding: number;
    overdueCount: number;
    /** Students with expected tuition &gt; 0 (excludes no_fee only). */
    billedStudentCount: number;
  };
}

export interface PaymentRecord {
  id: string;
  studentId: string;
  studentName: string;
  amount: number;
  type: PaymentType;
  periodMonth: string;
  note: string | null;
  recordedAt: string;
}

export type PaymentCreatePayload = {
  studentId: string;
  amount: number;
  type: PaymentType;
  periodMonth: string;
  note?: string | null;
};

export type ReceivablesListParams = {
  periodMonth: string;
  status: BillingStatus | "all";
  search: string;
};
