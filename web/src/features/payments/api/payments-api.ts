import { apiClient } from "@/shared/api/client";
import { computeBillingStatus } from "@/features/payments/lib/compute-billing-status";
import { mockPaymentsStore } from "@/features/payments/lib/mock-payments-store";
import { dueDateForPeriodMonth } from "@/features/payments/lib/period";
import type {
  PaymentCreatePayload,
  PaymentRecord,
  ReceivableRow,
  ReceivablesListParams,
  ReceivablesResponse,
} from "@/features/payments/model/types";
import { mockGroupsStore } from "@/features/groups/lib/mock-groups-store";
import { studentsApi } from "@/features/students/api/students-api";

function usePaymentsDemo() {
  return import.meta.env.VITE_PAYMENTS_DEMO === "true";
}

function summarizeReceivables(rows: ReceivableRow[]): ReceivablesResponse["summary"] {
  let totalExpected = 0;
  let totalPaid = 0;
  let totalOutstanding = 0;
  let overdueCount = 0;
  let billedStudentCount = 0;

  for (const r of rows) {
    if (r.tuitionExempt || r.status === "no_fee") continue;
    billedStudentCount += 1;
    totalExpected += r.expectedAmount;
    totalPaid += r.paidAmount;
    if (r.balance > 0) totalOutstanding += r.balance;
    if (r.status === "overdue") overdueCount += 1;
  }

  return {
    totalExpected,
    totalPaid,
    totalOutstanding,
    overdueCount,
    billedStudentCount,
  };
}

async function buildReceivableRows(
  periodMonth: string,
  search: string,
): Promise<ReceivableRow[]> {
  const { items: students } = await studentsApi.list({
    page: 1,
    pageSize: 500,
    search,
    status: "all",
  });

  const rows: ReceivableRow[] = [];

  for (const s of students) {
    const tuitionExempt = mockPaymentsStore.getTuitionExemptSync(s.id);
    const groupRow = s.groupId
      ? mockGroupsStore.getRowByIdSync(s.groupId)
      : undefined;
    const expectedAmount = groupRow?.monthlyFee ?? 0;
    const paidAmount = mockPaymentsStore.sumForStudentPeriod(s.id, periodMonth);
    const rawBalance = expectedAmount - paidAmount;
    const balance = tuitionExempt ? 0 : Math.max(0, rawBalance);
    const dueDate = dueDateForPeriodMonth(periodMonth);
    const status = computeBillingStatus({
      tuitionExempt,
      expectedAmount,
      paidAmount,
      balance: tuitionExempt ? 0 : rawBalance <= 0 ? 0 : rawBalance,
      dueDateYmd: dueDate,
    });

    rows.push({
      studentId: s.id,
      studentName: s.fullName,
      groupId: s.groupId,
      groupName: s.groupName,
      periodMonth,
      expectedAmount: tuitionExempt ? 0 : expectedAmount,
      paidAmount: tuitionExempt ? 0 : paidAmount,
      balance,
      dueDate,
      status,
      tuitionExempt,
    });
  }

  return rows;
}

/**
 * REST (adjust to your backend):
 * - GET `/payments/receivables?periodMonth=&status=&search=`
 * - GET `/payments` — ledger
 * - GET `/payments?studentId=` — history
 * - POST `/payments`
 * - PATCH `/payments/tuition-exempt` body: `{ studentId, exempt }` (super_admin)
 */
export const paymentsApi = {
  listReceivables: async (
    params: ReceivablesListParams,
  ): Promise<ReceivablesResponse> => {
    if (!usePaymentsDemo()) {
      const { data } = await apiClient.get<ReceivablesResponse>(
        "/payments/receivables",
        {
          params: {
            periodMonth: params.periodMonth,
            status: params.status === "all" ? undefined : params.status,
            search: params.search || undefined,
          },
        },
      );
      return data;
    }

    const allRows = await buildReceivableRows(
      params.periodMonth,
      params.search,
    );
    const summary = summarizeReceivables(allRows);
    const rows =
      params.status === "all"
        ? allRows
        : allRows.filter((r) => r.status === params.status);

    return {
      periodMonth: params.periodMonth,
      rows,
      summary,
    };
  },

  listLedger: async (): Promise<PaymentRecord[]> => {
    if (usePaymentsDemo()) {
      return mockPaymentsStore.listPayments();
    }
    const { data } = await apiClient.get<PaymentRecord[]>("/payments");
    return data;
  },

  listByStudent: async (studentId: string): Promise<PaymentRecord[]> => {
    if (usePaymentsDemo()) {
      return mockPaymentsStore.listByStudent(studentId);
    }
    const { data } = await apiClient.get<PaymentRecord[]>(
      `/payments/students/${studentId}`,
    );
    return data;
  },

  create: async (payload: PaymentCreatePayload): Promise<PaymentRecord> => {
    const student = await studentsApi.get(payload.studentId);
    if (usePaymentsDemo()) {
      return mockPaymentsStore.create({
        ...payload,
        studentName: student.fullName,
      });
    }
    const { data } = await apiClient.post<PaymentRecord>("/payments", payload);
    return data;
  },

  setTuitionExempt: async (
    studentId: string,
    exempt: boolean,
  ): Promise<void> => {
    if (usePaymentsDemo()) {
      await mockPaymentsStore.setTuitionExempt(studentId, exempt);
      return;
    }
    await apiClient.patch("/payments/tuition-exempt", { studentId, exempt });
  },
};
