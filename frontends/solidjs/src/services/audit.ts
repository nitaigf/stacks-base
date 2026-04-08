import { api } from './api';
import type { AuditLogFilters, AuditLogsEnvelope } from '../types/auth';

export async function listAuditLogs(filters: AuditLogFilters = {}) {
  const searchParams = new URLSearchParams();
  if (filters.page) {
    searchParams.set('page', String(filters.page));
  }
  if (filters.perPage) {
    searchParams.set('perPage', String(filters.perPage));
  }
  if (filters.query) {
    searchParams.set('query', filters.query);
  }
  if (filters.action) {
    searchParams.set('action', filters.action);
  }
  if (filters.resource) {
    searchParams.set('resource', filters.resource);
  }
  return api.get('api/v1/audit-logs', { searchParams }).json<AuditLogsEnvelope>();
}
