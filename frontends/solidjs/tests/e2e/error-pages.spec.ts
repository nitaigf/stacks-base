import { test, expect } from '@playwright/test';

test.describe('Error Pages', () => {
  test('shows 404 page for unknown routes', async ({ page }) => {
    await page.goto('/unknown-route-123');
    await expect(page.getByText('404')).toBeVisible();
    await expect(page.getByText('Pagina nao encontrada')).toBeVisible();
  });

  test('shows 403 page', async ({ page }) => {
    await page.goto('/errors/403');
    await expect(page.getByText('403')).toBeVisible();
    await expect(page.getByText('Acesso negado')).toBeVisible();
  });

  test('shows 500 page', async ({ page }) => {
    await page.goto('/errors/500');
    await expect(page.getByText('500')).toBeVisible();
    await expect(page.getByText('Erro interno no servidor')).toBeVisible();
  });

  test('error pages link back to home', async ({ page }) => {
    await page.goto('/errors/404');
    await page.getByRole('button', { name: 'Voltar para a home' }).click();
    await expect(page).toHaveURL('/');
  });
});
