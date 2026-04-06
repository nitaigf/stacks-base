import { z } from 'zod';

export const registerSchema = z.object({
  name: z.string().min(2, 'Informe pelo menos 2 caracteres.'),
  email: z.string().email('Informe um e-mail valido.'),
  password: z.string().min(8, 'A senha precisa ter pelo menos 8 caracteres.'),
});

export const loginSchema = z.object({
  email: z.string().email('Informe um e-mail valido.'),
  password: z.string().min(8, 'A senha precisa ter pelo menos 8 caracteres.'),
});

export type RegisterInput = z.infer<typeof registerSchema>;
export type LoginInput = z.infer<typeof loginSchema>;