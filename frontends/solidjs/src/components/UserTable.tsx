import type { User } from '../types/auth';

interface UserTableProps {
  users: User[];
  onStatusChange: (userId: string, status: 'active' | 'blocked') => void;
  onRoleChange: (userId: string, role: 'admin' | 'member') => void;
}

export function UserTable(props: UserTableProps) {
  const handleStatusToggle = (userId: string, currentStatus: 'active' | 'blocked') => {
    const newStatus = currentStatus === 'active' ? 'blocked' : 'active';
    props.onStatusChange(userId, newStatus);
  };

  const handleRoleToggle = (userId: string, currentRole: 'admin' | 'member') => {
    const newRole = currentRole === 'admin' ? 'member' : 'admin';
    props.onRoleChange(userId, newRole);
  };

  return (
    <div class="table-container">
      <table class="data-table">
        <thead>
          <tr>
            <th>Ações</th>
            <th>Nome</th>
            <th>Email</th>
            <th>Papel</th>
            <th>Status</th>
            <th>Criado em</th>
          </tr>
        </thead>
        <tbody>
          {props.users.map((user) => (
            <tr>
              <td>
                <div class="action-buttons">
                  <button
                    class="button button-sm button-secondary"
                    onClick={() => handleStatusToggle(user.id, user.status)}
                    title={user.status === 'active' ? 'Bloquear usuário' : 'Ativar usuário'}
                  >
                    {user.status === 'active' ? '🔒' : '✓'}
                  </button>
                  <button
                    class="button button-sm button-secondary"
                    onClick={() => handleRoleToggle(user.id, user.role)}
                    title={user.role === 'admin' ? 'Rebaixar para membro' : 'Promover para admin'}
                  >
                    {user.role === 'admin' ? '⬇' : '⬆'}
                  </button>
                </div>
              </td>
              <td>{user.name}</td>
              <td>{user.email}</td>
              <td>
                <span class={`badge ${user.role === 'admin' ? 'badge-primary' : 'badge-secondary'}`}>
                  {user.role === 'admin' ? 'Admin' : 'Membro'}
                </span>
              </td>
              <td>
                <span class={`badge ${user.status === 'active' ? 'badge-success' : 'badge-destructive'}`}>
                  {user.status === 'active' ? 'Ativo' : 'Bloqueado'}
                </span>
              </td>
              <td>{new Date(user.createdAt).toLocaleDateString('pt-BR')}</td>
            </tr>
          ))}
        </tbody>
      </table>
      {props.users.length === 0 && (
        <div class="empty-state">
          <p>Nenhum usuário encontrado.</p>
        </div>
      )}
    </div>
  );
}
