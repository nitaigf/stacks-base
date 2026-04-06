import type { User } from '../types/auth';

export type RouteMode = 'login' | 'register';

export type ResolvedRoute =
  | { kind: 'public'; path: '/' }
  | { kind: 'auth'; path: '/auth/login' | '/auth/register'; mode: RouteMode }
  | { kind: 'private'; path: '/app' }
  | { kind: 'admin'; path: '/admin' }
  | { kind: 'error'; path: string; statusCode: 403 | 404 | 500 };

export function resolveRoute(pathname: string, user: User | null): ResolvedRoute {
  switch (pathname) {
    case '/':
      return { kind: 'public', path: '/' };
    case '/auth/login':
      return { kind: 'auth', path: '/auth/login', mode: 'login' };
    case '/auth/register':
      return { kind: 'auth', path: '/auth/register', mode: 'register' };
    case '/app':
      return user ? { kind: 'private', path: '/app' } : { kind: 'auth', path: '/auth/login', mode: 'login' };
    case '/admin':
      if (!user) {
        return { kind: 'auth', path: '/auth/login', mode: 'login' };
      }

      return user.role === 'admin'
        ? { kind: 'admin', path: '/admin' }
        : { kind: 'error', path: '/errors/403', statusCode: 403 };
    case '/errors/403':
      return { kind: 'error', path: pathname, statusCode: 403 };
    case '/errors/500':
      return { kind: 'error', path: pathname, statusCode: 500 };
    default:
      return { kind: 'error', path: '/errors/404', statusCode: 404 };
  }
}

export function navigate(path: string) {
  window.history.pushState({}, '', path);
  window.dispatchEvent(new PopStateEvent('popstate'));
}