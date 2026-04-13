import {
  RouterProvider,
  createRootRoute,
  createRoute,
  createRouter,
  redirect,
  useNavigate,
  useParams,
  useSearch,
} from '@tanstack/solid-router';
import App from './App';
import { PublicPage } from './pages/PublicPage';
import { AuthPage } from './pages/AuthPage';
import { DashboardPage } from './pages/DashboardPage';
import { ErrorPage } from './pages/ErrorPage';
import { ForgotPasswordPage } from './pages/ForgotPasswordPage';
import { ResetPasswordPage } from './pages/ResetPasswordPage';
import { ChangePasswordPage } from './pages/ChangePasswordPage';
import { AdminUsersPage } from './pages/AdminUsersPage';
import { UserEditorPage } from './pages/UserEditorPage';
import { AuditLogsPage } from './pages/AuditLogsPage';
import { PublicLayout } from './layouts/PublicLayout';
import { AuthLayout } from './layouts/AuthLayout';
import { PrivateLayout } from './layouts/PrivateLayout';
import { AdminLayout } from './layouts/AdminLayout';
import { ErrorLayout } from './layouts/ErrorLayout';
import { authStore } from './stores/auth';
import { canAccessAdmin, canAccessPrivate, isAuthenticated } from './utils/access';
import { logout, me } from './services/auth';
import { isApiClientError } from './services/api';

const rootRoute = createRootRoute({
  component: App,
  notFoundComponent: Error404Scene,
});

const indexRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/',
  component: PublicScene,
});

const loginRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/auth/login',
  component: LoginScene,
});

const registerRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/auth/register',
  component: RegisterScene,
});

const forgotPasswordRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/auth/forgot-password',
  component: ForgotPasswordScene,
});

const resetPasswordRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/auth/reset-password',
  component: ResetPasswordScene,
});

const appRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/app',
  beforeLoad: privateGuard,
  component: DashboardScene,
});

const changePasswordRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/app/change-password',
  beforeLoad: privateGuard,
  component: ChangePasswordScene,
});

const adminUsersRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/admin/users',
  beforeLoad: adminGuard,
  component: AdminUsersScene,
});

const adminNewUserRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/admin/users/new',
  beforeLoad: adminGuard,
  component: AdminNewUserScene,
});

const adminEditUserRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/admin/users/$userId',
  beforeLoad: adminGuard,
  component: AdminEditUserScene,
});

const adminAuditLogsRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/admin/audit-logs',
  beforeLoad: adminGuard,
  component: AdminAuditLogsScene,
});

const adminIndexRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/admin',
  beforeLoad: adminGuard,
  loader: () => {
    throw redirect({ to: '/admin/users' });
  },
});

const errorRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/errors/$statusCode',
  component: ErrorScene,
});

const routeTree = rootRoute.addChildren([
  indexRoute,
  loginRoute,
  registerRoute,
  forgotPasswordRoute,
  resetPasswordRoute,
  appRoute,
  changePasswordRoute,
  adminUsersRoute,
  adminNewUserRoute,
  adminEditUserRoute,
  adminAuditLogsRoute,
  adminIndexRoute,
  errorRoute,
]);

export const router = createRouter({
  routeTree,
});

declare module '@tanstack/solid-router' {
  interface Register {
    router: typeof router;
  }
}

export function AppRouter() {
  return <RouterProvider router={router} />;
}

async function privateGuard() {
  if (!canAccessPrivate(authStore.currentUser()) || !authStore.accessToken()) {
    authStore.clearSession();
    throw redirect({ to: '/auth/login' });
  }

  try {
    const user = await authStore.revalidateSession(async () => {
      const response = await me();
      return response.data;
    });

    if (!isAuthenticated(user) || !canAccessPrivate(user)) {
      authStore.clearSession();
      throw redirect({ to: '/auth/login' });
    }

    return user;
  } catch (error) {
    if (isApiClientError(error)) {
      if (error.isUnauthorized) {
        throw redirect({ to: '/auth/login' });
      }
      if (error.isForbidden) {
        throw redirect({ to: '/errors/403' });
      }
      if (error.isServerError) {
        throw redirect({ to: '/errors/500' });
      }
    }

    throw error;
  }
}

async function adminGuard() {
  const user = await privateGuard();
  if (!canAccessAdmin(user)) {
    throw redirect({ to: '/errors/403' });
  }
}

function PublicScene() {
  const navigate = useNavigate();
  return (
    <PublicLayout>
      <PublicPage
        onLogin={() => navigate({ to: '/auth/login' })}
        onRegister={() => navigate({ to: '/auth/register' })}
      />
    </PublicLayout>
  );
}

function LoginScene() {
  return <AuthScene mode="login" />;
}

function RegisterScene() {
  return <AuthScene mode="register" />;
}

function AuthScene(props: { mode: 'login' | 'register' }) {
  const navigate = useNavigate();
  return (
    <AuthLayout onHome={() => navigate({ to: '/' })}>
      <AuthPage
        mode={props.mode}
        onModeChange={(mode) => navigate({ to: mode === 'login' ? '/auth/login' : '/auth/register' })}
        onAuthenticated={() => navigate({ to: '/app' })}
        onForgotPassword={() => navigate({ to: '/auth/forgot-password' })}
        onFatalError={() => navigate({ to: '/errors/500' })}
      />
    </AuthLayout>
  );
}

function ForgotPasswordScene() {
  const navigate = useNavigate();
  return (
    <AuthLayout onHome={() => navigate({ to: '/' })}>
      <ForgotPasswordPage
        onBackToLogin={() => navigate({ to: '/auth/login' })}
        onFatalError={() => navigate({ to: '/errors/500' })}
      />
    </AuthLayout>
  );
}

function ResetPasswordScene() {
  const navigate = useNavigate();
  const search = useSearch({ from: '/auth/reset-password' }) as { token?: string };
  return (
    <AuthLayout onHome={() => navigate({ to: '/' })}>
      <ResetPasswordPage
        token={search.token ?? ''}
        onBackToLogin={() => navigate({ to: '/auth/login' })}
        onFatalError={() => navigate({ to: '/errors/500' })}
      />
    </AuthLayout>
  );
}

function DashboardScene() {
  const navigate = useNavigate();
  const currentUser = () => authStore.currentUser();
  return (
    <PrivateLayout
      user={currentUser()}
      section="dashboard"
      onDashboard={() => navigate({ to: '/app' })}
      onLogout={async () => {
        try {
          await logout();
        } finally {
          authStore.clearSession();
          navigate({ to: '/auth/login' });
        }
      }}
      onAdmin={() => navigate({ to: '/admin/users' })}
      onChangePassword={() => navigate({ to: '/app/change-password' })}
    >
      <DashboardPage onFatalError={() => navigate({ to: '/errors/500' })} />
    </PrivateLayout>
  );
}

function ChangePasswordScene() {
  const navigate = useNavigate();
  const currentUser = () => authStore.currentUser();
  return (
    <PrivateLayout
      user={currentUser()}
      section="security"
      onDashboard={() => navigate({ to: '/app' })}
      onLogout={async () => {
        authStore.clearSession();
        navigate({ to: '/auth/login' });
      }}
      onAdmin={() => navigate({ to: '/admin/users' })}
      onChangePassword={() => navigate({ to: '/app/change-password' })}
    >
      <ChangePasswordPage
        onCompleted={() => {
          authStore.clearSession();
          navigate({ to: '/auth/login' });
        }}
        onBack={() => navigate({ to: '/app' })}
        onFatalError={() => navigate({ to: '/errors/500' })}
      />
    </PrivateLayout>
  );
}

function AdminUsersScene() {
  const navigate = useNavigate();
  const currentUser = () => authStore.currentUser();
  return (
    <AdminLayout
      user={currentUser()}
      section="users"
      onBackToApp={() => navigate({ to: '/app' })}
      onUsers={() => navigate({ to: '/admin/users' })}
      onAuditLogs={() => navigate({ to: '/admin/audit-logs' })}
    >
      <AdminUsersPage
        onCreateUser={() => navigate({ to: '/admin/users/new' })}
        onViewUser={(userId: string) => navigate({ to: '/admin/users/$userId', params: { userId } })}
        onFatalError={() => navigate({ to: '/errors/500' })}
      />
    </AdminLayout>
  );
}

function AdminNewUserScene() {
  const navigate = useNavigate();
  const currentUser = () => authStore.currentUser();
  return (
    <AdminLayout
      user={currentUser()}
      section="users"
      onBackToApp={() => navigate({ to: '/app' })}
      onUsers={() => navigate({ to: '/admin/users' })}
      onAuditLogs={() => navigate({ to: '/admin/audit-logs' })}
    >
      <UserEditorPage
        mode="create"
        onBack={() => navigate({ to: '/admin/users' })}
        onSaved={(userId: string) => navigate({ to: '/admin/users/$userId', params: { userId } })}
        onFatalError={() => navigate({ to: '/errors/500' })}
      />
    </AdminLayout>
  );
}

function AdminEditUserScene() {
  const navigate = useNavigate();
  const currentUser = () => authStore.currentUser();
  const params = useParams({ from: '/admin/users/$userId' }) as { userId: string };
  return (
    <AdminLayout
      user={currentUser()}
      section="users"
      onBackToApp={() => navigate({ to: '/app' })}
      onUsers={() => navigate({ to: '/admin/users' })}
      onAuditLogs={() => navigate({ to: '/admin/audit-logs' })}
    >
      <UserEditorPage
        mode="edit"
        userId={params.userId}
        onBack={() => navigate({ to: '/admin/users' })}
        onSaved={() => navigate({ to: '/admin/users/$userId', params: { userId: params.userId } })}
        onFatalError={() => navigate({ to: '/errors/500' })}
      />
    </AdminLayout>
  );
}

function AdminAuditLogsScene() {
  const navigate = useNavigate();
  const currentUser = () => authStore.currentUser();
  return (
    <AdminLayout
      user={currentUser()}
      section="audit"
      onBackToApp={() => navigate({ to: '/app' })}
      onUsers={() => navigate({ to: '/admin/users' })}
      onAuditLogs={() => navigate({ to: '/admin/audit-logs' })}
    >
      <AuditLogsPage onFatalError={() => navigate({ to: '/errors/500' })} />
    </AdminLayout>
  );
}

function ErrorScene() {
  const navigate = useNavigate();
  const params = useParams({ from: '/errors/$statusCode' }) as { statusCode: string };
  const statusCode = Number(params.statusCode) as 403 | 404 | 500;
  return (
    <ErrorLayout>
      <ErrorPage statusCode={statusCode} onHome={() => navigate({ to: '/' })} onLogin={() => navigate({ to: '/auth/login' })} />
    </ErrorLayout>
  );
}

function Error404Scene() {
  const navigate = useNavigate();
  return (
    <ErrorLayout>
      <ErrorPage statusCode={404} onHome={() => navigate({ to: '/' })} onLogin={() => navigate({ to: '/auth/login' })} />
    </ErrorLayout>
  );
}
