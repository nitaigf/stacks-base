import { z } from 'zod';
import type { UserRole, UserStatus } from '../types/auth';

const nameSchema = z.string().trim().min(2, 'Informe pelo menos 2 caracteres.');
const emailSchema = z.string().trim().email('Informe um e-mail valido.');
const passwordSchema = z.string().min(8, 'A senha precisa ter pelo menos 8 caracteres.');
const tokenSchema = z.string().trim().min(1, 'Informe o token recebido.');
const userRoleSchema = z.string().refine((value): value is UserRole => value === 'admin' || value === 'member', {
  message: 'Selecione um papel valido.',
});
const userStatusSchema = z.string().refine((value): value is UserStatus => value === 'active' || value === 'inactive', {
  message: 'Selecione um status valido.',
});

export const registerSchema = z.object({
  name: nameSchema,
  email: emailSchema,
  password: passwordSchema,
});

export const loginSchema = z.object({
  email: emailSchema,
  password: passwordSchema,
});

export const forgotPasswordSchema = z.object({
  email: emailSchema,
});

export const resetPasswordSchema = z
  .object({
    token: tokenSchema,
    newPassword: passwordSchema,
    confirmPassword: z.string(),
  })
  .refine((input) => input.newPassword === input.confirmPassword, {
    path: ['confirmPassword'],
    message: 'As senhas devem coincidir.',
  });

export const changePasswordSchema = z
  .object({
    currentPassword: passwordSchema,
    newPassword: passwordSchema,
    confirmPassword: z.string(),
  })
  .refine((input) => input.newPassword === input.confirmPassword, {
    path: ['confirmPassword'],
    message: 'As senhas devem coincidir.',
  });

export const userCreateSchema = z.object({
  name: nameSchema,
  email: emailSchema,
  password: passwordSchema,
  role: userRoleSchema,
  status: userStatusSchema,
});

export const userUpdateSchema = z.object({
  name: nameSchema,
  email: emailSchema,
  role: userRoleSchema,
});

export type RegisterInput = z.infer<typeof registerSchema>;
export type LoginInput = z.infer<typeof loginSchema>;
export type ForgotPasswordInput = z.infer<typeof forgotPasswordSchema>;
export type ResetPasswordInput = z.infer<typeof resetPasswordSchema>;
export type ChangePasswordInput = z.infer<typeof changePasswordSchema>;
export type UserCreateInput = z.infer<typeof userCreateSchema>;
export type UserUpdateInput = z.infer<typeof userUpdateSchema>;
