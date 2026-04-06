import { createMemo, createSignal, onCleanup, onMount } from 'solid-js';
import { AuthPage } from './pages/AuthPage';
import { AdminPage } from './pages/AdminPage';
import { DashboardPage } from './pages/DashboardPage';
import { ErrorPage } from './pages/ErrorPage';
import { PublicPage } from './pages/PublicPage';
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

        {route().kind === 'public' ? <PublicPage onLogin={goLogin} onRegister={goRegister} /> : null}
        {route().kind === 'auth' ? (
          <AuthPage
            mode={route().mode}
            onModeChange={(mode) => navigate(mode === 'login' ? '/auth/login' : '/auth/register')}
            onAuthenticated={goApp}
            onFatalError={goFatalError}
          />
        ) : null}
        {route().kind === 'private' ? (
          <DashboardPage onLoggedOut={goLogin} onAdmin={goAdmin} onFatalError={goFatalError} />
        ) : null}
        {route().kind === 'admin' ? <AdminPage onBackToApp={goApp} /> : null}
        {route().kind === 'error' ? <ErrorPage statusCode={route().statusCode} onHome={goHome} onLogin={goLogin} /> : null}
      </div>
    </main>
  );
}