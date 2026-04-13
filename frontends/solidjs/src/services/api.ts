import ky from 'ky';
import { authStore } from '../stores/auth';
import type { ApiErrorEnvelope } from '../types/auth';

const apiBaseUrl = import.meta.env.VITE_API_BASE_URL ?? 'http://127.0.0.1:8080';

type ApiErrorDetails = Record<string, string> | undefined;

export class ApiClientError extends Error {
  readonly status: number;
  readonly code: string;
  readonly details?: ApiErrorDetails;

  constructor({
    message,
    status,
    code,
    details,
  }: {
    message: string;
    status: number;
    code: string;
    details?: ApiErrorDetails;
  }) {
    super(message);
    this.name = 'ApiClientError';
    this.status = status;
    this.code = code;
    this.details = details;
  }

  get isUnauthorized() {
    return this.status === 401 || this.code === 'unauthorized';
  }

  get isForbidden() {
    return this.status === 403 || this.code === 'forbidden';
  }

  get isServerError() {
    return this.status >= 500 || this.code === 'internal_error';
  }
}

export function isApiClientError(error: unknown): error is ApiClientError {
  return error instanceof ApiClientError;
}

export const api = ky.create({
  prefixUrl: apiBaseUrl,
  credentials: 'include',
  hooks: {
    beforeRequest: [
      (request) => {
        const token = authStore.accessToken();
        if (token) {
          request.headers.set('Authorization', `Bearer ${token}`);
        }
      },
    ],
    afterResponse: [
      async (_request, _options, response) => {
        if (response.ok) {
          return response;
        }

        let message = 'Falha ao processar a requisicao.';
        let code = `http_${response.status}`;
        let details: ApiErrorDetails;
        try {
          const payload = (await response.clone().json()) as ApiErrorEnvelope;
          message = payload.error.message;
          code = payload.error.code;
          details = payload.error.details;
        } catch {
          message = response.statusText || message;
        }

        const error = new ApiClientError({
          message,
          status: response.status,
          code,
          details,
        });

        if (error.isUnauthorized) {
          authStore.clearSession();
        }

        throw error;
      },
    ],
  },
});
