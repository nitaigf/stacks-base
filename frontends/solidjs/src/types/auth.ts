export type UserRole = 'admin' | 'member';
export type UserStatus = 'active' | 'inactive';

export type User = {
  id: string;
  name: string;
  email: string;
  role: UserRole;
  status: UserStatus;
  deletedAt?: string | null;
  deletedBy?: string | null;
  lastLoginAt?: string | null;
  createdAt: string;
  updatedAt: string;
};

export type AuditLog = {
  id: string;
  actorUserId?: string | null;
  actorName?: string | null;
  actorEmail?: string | null;
  action: string;
  resource: string;
  resourceId?: string | null;
  route: string;
  method: string;
  ipAddress?: string | null;
  userAgent?: string | null;
  metadata: Record<string, unknown>;
  createdAt: string;
};

export type PaginationMeta = {
  page: number;
  perPage: number;
  total: number;
  totalPages: number;
};

export type AuthPayload = {
  accessToken: string;
  user: User;
};

export type AuthEnvelope = {
  data: AuthPayload;
};

export type UserEnvelope = {
  data: User;
};

export type UsersEnvelope = {
  data: User[];
  meta: PaginationMeta;
};

export type AuditLogsEnvelope = {
  data: AuditLog[];
  meta: PaginationMeta;
};

export type MessageEnvelope = {
  data: {
    message: string;
  };
};

export type ApiErrorEnvelope = {
  error: {
    code: string;
    message: string;
    details?: Record<string, string>;
  };
};

export type UserListFilters = {
  page?: number;
  perPage?: number;
  query?: string;
  role?: UserRole | '';
  status?: UserStatus | '';
  includeDeleted?: boolean;
};

export type AuditLogFilters = {
  page?: number;
  perPage?: number;
  query?: string;
  action?: string;
  resource?: string;
};
