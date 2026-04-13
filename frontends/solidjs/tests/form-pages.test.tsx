import { cleanup, fireEvent, render, screen, waitFor } from '@solidjs/testing-library';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';

const authServiceMocks = vi.hoisted(() => ({
  forgotPassword: vi.fn(),
  resetPassword: vi.fn(),
  changePassword: vi.fn(),
  logout: vi.fn(),
}));

const userServiceMocks = vi.hoisted(() => ({
  createUser: vi.fn(),
  getUser: vi.fn(),
  updateUser: vi.fn(),
}));

vi.mock('../src/services/auth', () => ({
  forgotPassword: authServiceMocks.forgotPassword,
  resetPassword: authServiceMocks.resetPassword,
  changePassword: authServiceMocks.changePassword,
  logout: authServiceMocks.logout,
}));

vi.mock('../src/services/users', () => ({
  createUser: userServiceMocks.createUser,
  getUser: userServiceMocks.getUser,
  updateUser: userServiceMocks.updateUser,
}));

import { ChangePasswordPage } from '../src/pages/ChangePasswordPage';
import { ForgotPasswordPage } from '../src/pages/ForgotPasswordPage';
import { ResetPasswordPage } from '../src/pages/ResetPasswordPage';
import { UserEditorPage } from '../src/pages/UserEditorPage';

describe('form pages', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  afterEach(() => {
    cleanup();
  });

  it('prevents forgot password submit when the e-mail is invalid', async () => {
    render(() => <ForgotPasswordPage onBackToLogin={vi.fn()} onFatalError={vi.fn()} />);

    await fireEvent.input(screen.getByLabelText('E-mail'), {
      target: { value: 'invalid-email' },
    });
    await fireEvent.click(screen.getByRole('button', { name: 'Enviar link' }));

    expect(authServiceMocks.forgotPassword).not.toHaveBeenCalled();
    await waitFor(() => {
      expect(screen.getByText('Informe um e-mail valido.')).toBeTruthy();
    });
  });

  it('submits forgot password with a trimmed e-mail', async () => {
    authServiceMocks.forgotPassword.mockResolvedValue({
      data: { message: 'If the account exists, a password recovery message has been sent.' },
    });

    render(() => <ForgotPasswordPage onBackToLogin={vi.fn()} onFatalError={vi.fn()} />);

    await fireEvent.input(screen.getByLabelText('E-mail'), {
      target: { value: '  admin@example.com  ' },
    });
    await fireEvent.click(screen.getByRole('button', { name: 'Enviar link' }));

    await waitFor(() => {
      expect(authServiceMocks.forgotPassword).toHaveBeenCalledWith('admin@example.com');
    });
    expect(screen.getByText('If the account exists, a password recovery message has been sent.')).toBeTruthy();
  });

  it('prevents reset password submit when the confirmation does not match', async () => {
    render(() => <ResetPasswordPage token="reset-token" onBackToLogin={vi.fn()} onFatalError={vi.fn()} />);

    await fireEvent.input(screen.getByLabelText('Nova senha'), {
      target: { value: 'password123' },
    });
    await fireEvent.input(screen.getByLabelText('Confirmar nova senha'), {
      target: { value: 'password456' },
    });
    await fireEvent.click(screen.getByRole('button', { name: 'Redefinir senha' }));

    expect(authServiceMocks.resetPassword).not.toHaveBeenCalled();
    expect(screen.getByText('As senhas devem coincidir.')).toBeTruthy();
  });

  it('logs out after changing the password even when logout cleanup fails', async () => {
    const onCompleted = vi.fn();

    authServiceMocks.changePassword.mockResolvedValue({
      data: { message: 'Password changed successfully. Please sign in again.' },
    });
    authServiceMocks.logout.mockRejectedValue(new Error('session not found'));

    render(() => <ChangePasswordPage onCompleted={onCompleted} onBack={vi.fn()} onFatalError={vi.fn()} />);

    await fireEvent.input(screen.getByLabelText('Senha atual'), {
      target: { value: 'password123' },
    });
    await fireEvent.input(screen.getByLabelText('Nova senha'), {
      target: { value: 'password456' },
    });
    await fireEvent.input(screen.getByLabelText('Confirmar nova senha'), {
      target: { value: 'password456' },
    });
    await fireEvent.click(screen.getByRole('button', { name: 'Salvar nova senha' }));

    await waitFor(() => {
      expect(authServiceMocks.changePassword).toHaveBeenCalledWith('password123', 'password456');
    });
    await waitFor(() => {
      expect(authServiceMocks.logout).toHaveBeenCalledTimes(1);
    });
    await waitFor(() => {
      expect(onCompleted).toHaveBeenCalledTimes(1);
    });
  });

  it('prevents user creation submit when required fields are invalid', async () => {
    render(() => <UserEditorPage mode="create" onBack={vi.fn()} onSaved={vi.fn()} onFatalError={vi.fn()} />);

    await fireEvent.input(screen.getByLabelText('Nome'), {
      target: { value: 'A' },
    });
    await fireEvent.input(screen.getByLabelText('E-mail'), {
      target: { value: 'invalid-email' },
    });
    await fireEvent.input(screen.getByLabelText('Senha inicial'), {
      target: { value: '123' },
    });
    await fireEvent.click(screen.getByRole('button', { name: 'Criar usuario' }));

    expect(userServiceMocks.createUser).not.toHaveBeenCalled();
    await waitFor(() => {
      expect(screen.getByText('Informe pelo menos 2 caracteres.')).toBeTruthy();
      expect(screen.getByText('Informe um e-mail valido.')).toBeTruthy();
      expect(screen.getByText('A senha precisa ter pelo menos 8 caracteres.')).toBeTruthy();
    });
  });

  it('submits user creation with validated payload data', async () => {
    const onSaved = vi.fn();

    userServiceMocks.createUser.mockResolvedValue({
      data: { id: 'user-2' },
    });

    render(() => <UserEditorPage mode="create" onBack={vi.fn()} onSaved={onSaved} onFatalError={vi.fn()} />);

    await fireEvent.input(screen.getByLabelText('Nome'), {
      target: { value: '  Admin User  ' },
    });
    await fireEvent.input(screen.getByLabelText('E-mail'), {
      target: { value: '  admin@example.com  ' },
    });
    await fireEvent.input(screen.getByLabelText('Papel'), {
      target: { value: 'admin' },
    });
    await fireEvent.input(screen.getByLabelText('Senha inicial'), {
      target: { value: 'password123' },
    });
    await fireEvent.click(screen.getByRole('button', { name: 'Criar usuario' }));

    await waitFor(() => {
      expect(userServiceMocks.createUser).toHaveBeenCalledWith({
        name: 'Admin User',
        email: 'admin@example.com',
        password: 'password123',
        role: 'admin',
        status: 'active',
      });
    });
    await waitFor(() => {
      expect(onSaved).toHaveBeenCalledWith('user-2');
    });
  });

  it('loads the user in edit mode and submits validated updates', async () => {
    const onSaved = vi.fn();

    userServiceMocks.getUser.mockResolvedValue({
      data: {
        id: 'user-1',
        name: 'Existing User',
        email: 'existing@example.com',
        role: 'member',
        status: 'active',
        lastLoginAt: null,
        deletedAt: null,
      },
    });
    userServiceMocks.updateUser.mockResolvedValue({
      data: { id: 'user-1' },
    });

    render(() => <UserEditorPage mode="edit" userId="user-1" onBack={vi.fn()} onSaved={onSaved} onFatalError={vi.fn()} />);

    await waitFor(() => {
      expect(userServiceMocks.getUser).toHaveBeenCalledWith('user-1');
    });
    await waitFor(() => {
      expect(screen.getByDisplayValue('Existing User')).toBeTruthy();
    });

    await fireEvent.input(screen.getByLabelText('Nome'), {
      target: { value: '  Updated User  ' },
    });
    await fireEvent.input(screen.getByLabelText('E-mail'), {
      target: { value: '  updated@example.com  ' },
    });
    await fireEvent.input(screen.getByLabelText('Papel'), {
      target: { value: 'admin' },
    });
    await fireEvent.click(screen.getByRole('button', { name: 'Salvar alteracoes' }));

    await waitFor(() => {
      expect(userServiceMocks.updateUser).toHaveBeenCalledWith('user-1', {
        name: 'Updated User',
        email: 'updated@example.com',
        role: 'admin',
      });
    });
    await waitFor(() => {
      expect(onSaved).toHaveBeenCalledWith('user-1');
    });
  });
});
