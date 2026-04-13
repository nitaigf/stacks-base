import { defineConfig, devices } from '@playwright/test';

const explicitBaseURL = process.env.PLAYWRIGHT_BASE_URL;
const port = Number(process.env.PLAYWRIGHT_PORT ?? 4173);
const baseURL = explicitBaseURL ?? `http://127.0.0.1:${port}`;

export default defineConfig({
  testDir: './tests/e2e',
  timeout: 30_000,
  expect: {
    timeout: 5_000,
  },
  fullyParallel: false,
  reporter: 'list',
  use: {
    baseURL,
    trace: 'on-first-retry',
    screenshot: 'only-on-failure',
    video: 'retain-on-failure',
  },
  ...(explicitBaseURL
    ? {}
    : {
        webServer: {
          command: `npm run build && npm run preview -- --host 127.0.0.1 --port ${port}`,
          url: baseURL,
          reuseExistingServer: true,
          stdout: 'ignore' as const,
          stderr: 'pipe' as const,
        },
      }),
  projects: [
    {
      name: 'chromium',
      use: {
        ...devices['Desktop Chrome'],
        launchOptions: {
          args: ['--disable-web-security'],
        },
      },
    },
  ],
});
