import { test, expect } from '@playwright/test';

test.describe('Authentication Flows', () => {
  const timestamp = Date.now();
  const testEmail = `user-${timestamp}@test.local`;
  const testPassword = 'Password123!';

  test('register new account', async ({ page }) => {
    await page.goto('/auth/register');
    
    await page.getByLabel('Nome').fill('Register E2E User');
    await page.getByLabel('E-mail').fill(testEmail);
    await page.getByLabel('Senha').fill(testPassword);
    
    await page.getByRole('button', { name: 'Criar conta', exact: true }).click();
    
    await expect(page).toHaveURL('/app');
    await expect(page.getByText('Sessao iniciada com sucesso.')).toBeVisible();
    await expect(page.getByText('Register E2E User')).toBeVisible();
  });

  test('logout from dashboard', async ({ page }) => {
    // Login first
    await page.goto('/auth/login');
    await page.getByLabel('E-mail').fill(testEmail);
    await page.getByLabel('Senha').fill(testPassword);
    await page.getByRole('button', { name: 'Entrar', exact: true }).click();
    
    await expect(page).toHaveURL('/app');
    
    // Click logout
    await page.getByRole('button', { name: 'Sair' }).click();
    
    await expect(page).toHaveURL('/auth/login');
    await expect(page.getByText('Entrar na referencia')).toBeVisible();
  });

  test('login with registered account', async ({ page }) => {
    await page.goto('/auth/login');
    
    await page.getByLabel('E-mail').fill(testEmail);
    await page.getByLabel('Senha').fill(testPassword);
    
    await page.getByRole('button', { name: 'Entrar', exact: true }).click();
    
    await expect(page).toHaveURL('/app');
    await expect(page.getByText('Sessao iniciada com sucesso.')).toBeVisible();
  });

  test('forgot password shows confirmation', async ({ page }) => {
    await page.goto('/auth/login');
    await page.getByRole('button', { name: 'Esqueci minha senha' }).click();
    
    await expect(page).toHaveURL('/auth/forgot-password');
    
    await page.getByLabel('E-mail').fill(testEmail);
    await page.getByRole('button', { name: 'Enviar link' }).click();
    
    await expect(page.getByText(/If the account exists/i)).toBeVisible();
  });

  test('visit /app without auth redirects to login', async ({ page }) => {
    await page.goto('/app');
    await expect(page).toHaveURL('/auth/login');
  });
});
