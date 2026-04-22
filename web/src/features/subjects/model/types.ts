export type SubjectOption = { id: string; name: string };

export interface Subject {
  id: string;
  name: string;
  /** Short label e.g. MATH-101 */
  code: string | null;
  description: string | null;
  createdAt: string;
  updatedAt: string;
}

export type SubjectsListParams = {
  page: number;
  pageSize: number;
  search: string;
};

export type SubjectsListResponse = {
  items: Subject[];
  total: number;
  page: number;
  pageSize: number;
};

export type SubjectCreatePayload = {
  name: string;
  code?: string | null;
  description?: string | null;
};

export type SubjectUpdatePayload = Partial<SubjectCreatePayload>;
