import { test, expect } from '@playwright/test';

test.describe('Admin User Management', () => {
  const adminEmail = 'admin@stacks-base.local';
  const adminPassword = 'Admin@123456';

  test.beforeEach(async ({ page }) => {
    // Login as admin before each test
    await page.goto('/auth/login');
    await page.getByLabel('E-mail').fill(adminEmail);
    await page.getByLabel('Senha').fill(adminPassword);
    await page.getByRole('button', { name: 'Entrar', exact: true }).click();
    await expect(page).toHaveURL('/app');
    
    // Navigate to Admin Users
    await page.getByRole('button', { name: 'Area admin' }).click();
    await expect(page).toHaveURL('/admin/users');
  });

  test('lists users with pagination', async ({ page }) => {
    await expect(page.getByText('Controle operacional de usuarios')).toBeVisible();
    await expect(page.locator('table.data-table')).toBeVisible();
    await expect(page.getByText(/Pagina 1 de/)).toBeVisible();
  });

  test('creates a new user', async ({ page }) => {
    const timestamp = Date.now();
    const newUserEmail = `new-user-${timestamp}@test.local`;

    await page.getByRole('button', { name: 'Novo usuario' }).click();
    await expect(page).toHaveURL('/admin/users/new');

    await page.getByLabel('Nome').fill('Admin Created User');
    await page.getByLabel('E-mail').fill(newUserEmail);
    await page.getByLabel('Papel').selectOption('member');
    await page.getByLabel('Senha inicial').fill('Password123!');
    
    await page.getByRole('button', { name: 'Criar usuario' }).click();
    
    // Should redirect to edit/view page
    await expect(page).toHaveURL(/\/admin\/users\/.+/);
    await expect(page.getByDisplayValue('Admin Created User')).toBeVisible();
  });

  test('filters by search query', async ({ page }) => {
    await page.getByLabel('Busca').fill('Admin');
    await page.getByRole('button', { name: 'Aplicar filtros' }).click();
    
    // Wait for table to update
    await expect(page.locator('table.data-table tbody tr')).toContainText(['Admin']);
  });

  test('triggers exports', async ({ page }) => {
    // We only test that clicking them doesn't crash and maybe show some feedback if implemented
    // Full download testing is a bit more complex in headless CI but we can at least click
    await page.getByRole('button', { name: 'Exportar CSV' }).click();
    await page.getByRole('button', { name: 'Exportar XLSX' }).click();
    await page.getByRole('button', { name: 'Imprimir PDF' }).click();
  });
});
