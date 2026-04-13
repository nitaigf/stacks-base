import { test, expect } from '@playwright/test';

test.describe('Admin Audit Logs', () => {
  const adminEmail = 'admin@stacks-base.local';
  const adminPassword = 'Admin@123456';

  test.beforeEach(async ({ page }) => {
    await page.goto('/auth/login');
    await page.getByLabel('E-mail').fill(adminEmail);
    await page.getByLabel('Senha').fill(adminPassword);
    await page.getByRole('button', { name: 'Entrar', exact: true }).click();
    await expect(page).toHaveURL('/app');
    
    await page.getByRole('button', { name: 'Area admin' }).click();
    await page.getByRole('link', { name: 'Auditoria' }).click();
    await expect(page).toHaveURL('/admin/audit-logs');
  });

  test('lists audit logs', async ({ page }) => {
    await expect(page.getByText('Monitoramento de auditoria')).toBeVisible();
    await expect(page.locator('table.data-table')).toBeVisible();
  });

  test('filters audit logs', async ({ page }) => {
    await page.getByLabel('Acao').fill('login');
    await page.getByRole('button', { name: 'Aplicar filtros' }).click();
    
    await expect(page.locator('table.data-table tbody tr')).toContainText(['login']);
  });
});
