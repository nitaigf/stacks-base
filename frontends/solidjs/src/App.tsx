import { createMemo, createSignal, onCleanup, onMount } from 'solid-js';
import { AuthPage } from './pages/AuthPage';
import { AdminPage } from './pages/AdminPage';
import { DashboardPage } from './pages/DashboardPage';
import { ErrorPage } from './pages/ErrorPage';
import { PublicPage } from './pages/PublicPage';
import { PublicLayout } from './layouts/PublicLayout';
import { AuthLayout } from './layouts/AuthLayout';
import { PrivateLayout } from './layouts/PrivateLayout';
import { AdminLayout } from './layouts/AdminLayout';
import { ErrorLayout } from './layouts/ErrorLayout';
import { authStore } from './stores/auth';
import { navigate, resolveRoute } from './utils/router';

export default function App() {
  const currentUser = createMemo(() => authStore.currentUser());
  const [pathname, setPathname] = createSignal(window.location.pathname);

  onMount(() => {
    const syncRoute = () => setPathname(window.location.pathname);
    window.addEventListener('popstate', syncRoute);
    onCleanup(() => window.removeEventListener('popstate', syncRoute));
  });

  const route = createMemo(() => resolveRoute(pathname(), currentUser()));

  const goHome = () => navigate('/');
  const goLogin = () => navigate('/auth/login');
  const goRegister = () => navigate('/auth/register');
  const goApp = () => navigate('/app');
  const goAdmin = () => navigate('/admin');
  const goFatalError = () => navigate('/errors/500');

  const handleLogout = async () => {
    try {
      const { logout } = await import('./services/auth');
      await logout();
      authStore.clearSession();
      goLogin();
    } catch {
      goFatalError();
    }
  };

  return (
    <main class="app-shell">
      <div class="page-grid">
        <section class="hero-panel">
          <span class="eyebrow">Stacks Base</span>
          <h1 class="hero-title">SolidJS e Go viram a primeira referencia executavel.</h1>
          <p class="hero-copy">
            Esta baseline demonstra o contrato unico, o design system compartilhado e o fluxo minimo de autenticacao
            sobre PostgreSQL local na porta 5432.
          </p>
          <ul class="stack-list">
            <li>
              <span>Frontend</span>
              <strong>SolidJS</strong>
            </li>
            <li>
              <span>Backend</span>
              <strong>Go net/http</strong>
            </li>
            <li>
              <span>Banco</span>
              <strong>PostgreSQL local 5432</strong>
            </li>
            <li>
              <span>Email</span>
              <strong>Mailtrap SMTP</strong>
            </li>
          </ul>
        </section>

        {route().kind === 'public' ? (
          <PublicLayout>
            <PublicPage onLogin={goLogin} onRegister={goRegister} />
          </PublicLayout>
        ) : null}
        {route().kind === 'auth' ? (
          <AuthLayout onHome={goHome}>
            <AuthPage
              mode={route().mode}
              onModeChange={(mode) => navigate(mode === 'login' ? '/auth/login' : '/auth/register')}
              onAuthenticated={goApp}
              onFatalError={goFatalError}
            />
          </AuthLayout>
        ) : null}
        {route().kind === 'private' ? (
          <PrivateLayout
            user={currentUser()}
            onLogout={handleLogout}
            onAdmin={goAdmin}
            onFatalError={goFatalError}
          >
            <DashboardPage onFatalError={goFatalError} />
          </PrivateLayout>
        ) : null}
        {route().kind === 'admin' ? (
          <AdminLayout onBackToApp={goApp}>
            <AdminPage />
          </AdminLayout>
        ) : null}
        {route().kind === 'error' ? (
          <ErrorLayout>
            <ErrorPage statusCode={route().statusCode} onHome={goHome} onLogin={goLogin} />
          </ErrorLayout>
        ) : null}
      </div>
    </main>
  );
}