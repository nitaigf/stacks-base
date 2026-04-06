import ky from 'ky';
import { authStore } from '../stores/auth';
import type { ApiErrorEnvelope } from '../types/auth';

const apiBaseUrl = import.meta.env.VITE_API_BASE_URL ?? 'http://127.0.0.1:8080';

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
        try {
          const payload = (await response.clone().json()) as ApiErrorEnvelope;
          message = payload.error.message;
        } catch {
          message = response.statusText || message;
        }

        throw new Error(message);
      },
    ],
  },
});