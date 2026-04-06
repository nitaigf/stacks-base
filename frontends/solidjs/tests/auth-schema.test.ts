import { describe, expect, it } from 'vitest';
import { loginSchema, registerSchema } from '../src/schemas/auth';

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
});