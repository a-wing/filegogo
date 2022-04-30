import { test, expect } from '@playwright/test';
import { spawn } from 'child_process';

import fs from 'fs';
import md5 from 'md5';

const share = 'http://localhost:8080/2222';
//const share = 'https://send.22333.fun/2222';
const file = 'playwright.config.ts';

test('cli to browser', async ({ page }) => {
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

  const srcMd5 = md5(await fs.promises.readFile(file));
  const dstMd5 = md5(await fs.promises.readFile(path));
  expect(srcMd5).toBe(dstMd5);
});
