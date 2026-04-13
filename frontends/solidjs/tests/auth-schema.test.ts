import { describe, expect, it } from 'vitest';
import {
  changePasswordSchema,
  forgotPasswordSchema,
  loginSchema,
  registerSchema,
  resetPasswordSchema,
  userCreateSchema,
  userUpdateSchema,
} from '../src/schemas/auth';

describe('auth schemas', () => {
  it('accepts a valid register payload', () => {
    const result = registerSchema.safeParse({
      name: 'Nitai',
      email: 'nitai@example.com',
      password: 'password123',
    });

    expect(result.success).toBe(true);
  });

  it('rejects an invalid login payload', () => {
    const result = loginSchema.safeParse({
      email: 'not-an-email',
      password: '123',
    });

    expect(result.success).toBe(false);
  });

  it('trims and accepts a forgot password payload with a valid e-mail', () => {
    const result = forgotPasswordSchema.safeParse({
      email: '  nitai@example.com  ',
    });

    expect(result.success).toBe(true);
    expect(result.success && result.data.email).toBe('nitai@example.com');
  });

  it('rejects reset password when confirmation does not match', () => {
    const result = resetPasswordSchema.safeParse({
      token: 'reset-token',
      newPassword: 'password123',
      confirmPassword: 'password456',
    });

    expect(result.success).toBe(false);
  });

  it('rejects change password when the confirmation is missing', () => {
    const result = changePasswordSchema.safeParse({
      currentPassword: 'password123',
      newPassword: 'password456',
      confirmPassword: '',
    });

    expect(result.success).toBe(false);
  });

  it('accepts a valid user creation payload', () => {
    const result = userCreateSchema.safeParse({
      name: 'Admin User',
      email: 'admin@example.com',
      password: 'password123',
      role: 'admin',
      status: 'active',
    });

    expect(result.success).toBe(true);
  });

  it('rejects a user update payload with an invalid role', () => {
    const result = userUpdateSchema.safeParse({
      name: 'Admin User',
      email: 'admin@example.com',
      role: 'owner',
    });

    expect(result.success).toBe(false);
  });
});
