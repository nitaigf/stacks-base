import { createSignal, onMount } from 'solid-js';
import { UserTable } from '../components/UserTable';
import { listUsers, updateUserStatus, updateUserRole } from '../services/users';
import type { User } from '../types/auth';

export function AdminPage() {
  const [users, setUsers] = createSignal<User[]>([]);
  const [loading, setLoading] = createSignal(true);
  const [error, setError] = createSignal<string | null>(null);

  onMount(async () => {
    try {
      setLoading(true);
      // TODO: Implement backend endpoint for listing users
      // const response = await listUsers();
      // setUsers(response.data);
      
      // Mock data for now
      setUsers([
        {
          id: '1',
          name: 'Admin',
          email: 'admin@stacks-base.local',
          role: 'admin',
          status: 'active',
          createdAt: new Date().toISOString(),
          updatedAt: new Date().toISOString(),
        },
        {
          id: '2', 
          name: 'Usuário Teste',
          email: 'test@example.com',
          role: 'member',
          status: 'active',
          createdAt: new Date().toISOString(),
          updatedAt: new Date().toISOString(),
        }
      ]);
    } catch (err) {
      setError('Falha ao carregar usuários');
    } finally {
      setLoading(false);
    }
  });

  const handleStatusChange = async (userId: string, status: 'active' | 'blocked') => {
    try {
      // TODO: Implement backend endpoint
      // await updateUserStatus(userId, status);
      
      setUsers(prev => 
        prev.map(user => 
          user.id === userId ? { ...user, status, updatedAt: new Date().toISOString() } : user
        )
      );
    } catch (err) {
      setError('Falha ao atualizar status do usuário');
    }
  };

  const handleRoleChange = async (userId: string, role: 'admin' | 'member') => {
    try {
      // TODO: Implement backend endpoint
      // await updateUserRole(userId, role);
      
      setUsers(prev => 
        prev.map(user => 
          user.id === userId ? { ...user, role, updatedAt: new Date().toISOString() } : user
        )
      );
    } catch (err) {
      setError('Falha ao atualizar papel do usuário');
    }
  };

  return (
    <section class="surface-card dashboard-card">
      <div class="page-header">
        <span class="metric-pill metric-pill-tag">Pagina administrativa</span>
        <h2 class="form-title">Gerenciamento de Usuários</h2>
        <p class="form-copy">
          Visualize e gerencie todos os usuários do sistema.
        </p>
      </div>

      {error() && (
        <div class="alert alert-error">
          <span>{error()}</span>
        </div>
      )}

      {loading() ? (
        <div class="loading-state">
          <p>Carregando usuários...</p>
        </div>
      ) : (
        <UserTable 
          users={users()} 
          onStatusChange={handleStatusChange}
          onRoleChange={handleRoleChange}
        />
      )}

      <div class="metric-row">
        <span class="metric-pill metric-pill-tag">Total: {users().length} usuários</span>
        <span class="metric-pill metric-pill-tag">Admin only</span>
        <span class="metric-pill metric-pill-tag">403 previsível</span>
      </div>
    </section>
  );
}