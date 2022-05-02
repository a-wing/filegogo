import { PlaywrightTestConfig, devices } from '@playwright/test';

const config: PlaywrightTestConfig = {
  forbidOnly: !!process.env.CI,
  retries: process.env.CI ? 2 : 0,
  webServer: {
    command: '../filegogo server',
    port: 8080,
    timeout: 120 * 1000,
    reuseExistingServer: !process.env.CI,
  },
  workers: 1,
  use: {
    trace: 'on-first-retry',
  },
  projects: [
    {
      name: 'chromium',
      use: { ...devices['Desktop Chrome'] },
    },
    //{
    //  name: 'firefox',
    //  use: { ...devices['Desktop Firefox'] },
    //},
  ],
};
export default config;
