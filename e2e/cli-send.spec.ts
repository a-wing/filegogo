import { test, expect } from '@playwright/test';
import { spawn } from 'child_process';

import { getRoom, checkSum } from './helper';

const address = 'http://localhost:8080';
//const address = 'https://send.22333.fun';
const file = 'playwright.config.ts';

test('cli to browser', async ({ page }) => {
  const share = `${address}/${await getRoom('http://localhost:8080')}`;

  const ls = spawn('../filegogo', ['send', '-s', share, file]);

  ls.stdout.on('data', (data) => {
    console.log(`stdout: ${data}`);
  });

  await page.goto(share);
  expect(await page.title()).toBe('Filegogo');

  const [download] = await Promise.all([
    page.waitForEvent('download'),
    page.locator('text=getFile').click()
  ]);

  const path = await download.path();

  expect(await checkSum(file)).toBe(await checkSum(path));
});
