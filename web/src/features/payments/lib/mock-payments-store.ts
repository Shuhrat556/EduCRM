import { currentPeriodMonth } from "@/features/payments/lib/period";
import type { PaymentCreatePayload } from "@/features/payments/model/types";
import type { PaymentRecord } from "@/features/payments/model/types";
import type { PaymentType } from "@/features/payments/model/types";

const STORAGE_KEY = "educrm_mock_payments_v1";

export type PaymentRow = {
  id: string;
  studentId: string;
  studentName: string;
  amount: number;
  type: PaymentType;
  periodMonth: string;
  note: string | null;
  recordedAt: string;
};

type StoreShape = {
  payments: PaymentRow[];
  tuitionExemptByStudent: Record<string, boolean>;
};

function nowIso() {
  return new Date().toISOString();
}

function seed(): StoreShape {
  const t = nowIso();
  const period = currentPeriodMonth();
  return {
    payments: [
      {
        id: "pay_seed_1",
        studentId: "s_seed_1",
        studentName: "Dilnoza Rahimova",
        amount: 600_000,
        type: "bank_transfer",
        periodMonth: period,
        note: "Partial installment",
        recordedAt: t,
      },
      {
        id: "pay_seed_2",
        studentId: "s_seed_2",
        studentName: "Rustam Toshmatov",
        amount: 1_200_000,
        type: "cash",
        periodMonth: period,
        note: "April tuition",
        recordedAt: t,
      },
    ],
    tuitionExemptByStudent: {},
  };
}

function readStore(): StoreShape {
  try {
    const raw = localStorage.getItem(STORAGE_KEY);
    if (!raw) {
      const initial = seed();
      localStorage.setItem(STORAGE_KEY, JSON.stringify(initial));
      return initial;
    }
    const parsed = JSON.parse(raw) as StoreShape;
    if (!parsed?.payments || !parsed.tuitionExemptByStudent) return seed();
    return parsed;
  } catch {
    return seed();
  }
}

function writeStore(s: StoreShape) {
  localStorage.setItem(STORAGE_KEY, JSON.stringify(s));
}

function delay(ms = 200) {
  return new Promise((r) => setTimeout(r, ms));
}

export const mockPaymentsStore = {
  async getStore(): Promise<StoreShape> {
    await delay(80);
    return readStore();
  },

  async listPayments(): Promise<PaymentRow[]> {
    await delay(120);
    return [...readStore().payments].sort(
      (a, b) => b.recordedAt.localeCompare(a.recordedAt),
    );
  },

  async listByStudent(studentId: string): Promise<PaymentRow[]> {
    await delay(120);
    return readStore()
      .payments.filter((p) => p.studentId === studentId)
      .sort((a, b) => b.recordedAt.localeCompare(a.recordedAt));
  },

  sumForStudentPeriod(studentId: string, periodMonth: string): number {
    return readStore().payments
      .filter((p) => p.studentId === studentId && p.periodMonth === periodMonth)
      .reduce((acc, p) => acc + p.amount, 0);
  },

  getTuitionExemptSync(studentId: string): boolean {
    return Boolean(readStore().tuitionExemptByStudent[studentId]);
  },

  async create(
    payload: PaymentCreatePayload & { studentName: string },
  ): Promise<PaymentRecord> {
    await delay();
    const store = readStore();
    const row: PaymentRow = {
      id: crypto.randomUUID(),
      studentId: payload.studentId,
      studentName: payload.studentName,
      amount: payload.amount,
      type: payload.type,
      periodMonth: payload.periodMonth,
      note: payload.note?.trim() || null,
      recordedAt: nowIso(),
    };
    store.payments.unshift(row);
    writeStore(store);
    return { ...row };
  },

  async setTuitionExempt(studentId: string, exempt: boolean): Promise<void> {
    await delay(150);
    const store = readStore();
    if (exempt) {
      store.tuitionExemptByStudent[studentId] = true;
    } else {
      delete store.tuitionExemptByStudent[studentId];
    }
    writeStore(store);
  },
};
