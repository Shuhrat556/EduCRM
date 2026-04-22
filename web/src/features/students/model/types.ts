export type StudentStatus = "active" | "inactive" | "graduated" | "suspended";

export interface Student {
  id: string;
  fullName: string;
  phone: string;
  email: string;
  status: StudentStatus;
  groupId: string | null;
  groupName: string | null;
  photoUrl: string | null;
  createdAt: string;
  updatedAt: string;
}

export type StudentsListParams = {
  page: number;
  pageSize: number;
  search: string;
  status: StudentStatus | "all";
  /** When set, only students in this group (MVP: one group per student). */
  groupId?: string;
};

export type StudentsListResponse = {
  items: Student[];
  total: number;
  page: number;
  pageSize: number;
};

export type StudentGroupOption = {
  id: string;
  name: string;
};

export type StudentCreatePayload = {
  fullName: string;
  phone: string;
  email: string;
  status: StudentStatus;
  groupId?: string | null;
  photoUrl?: string | null;
};

export type StudentUpdatePayload = Partial<StudentCreatePayload>;
