import { test, expect } from '@playwright/test';

import fs from 'fs';
import md5 from 'md5';

const share = 'http://localhost:8080/1111';
//const share = 'https://send.22333.fun/1111';
const file = 'playwright.config.ts';

test('browser to browser', async ({ page, context }) => {
  await page.goto(share);
  expect(await page.title()).toBe('Filegogo');

  await page.setInputFiles('input#upload', file);

  const page2 = await context.newPage();
  await page2.goto(share);

  const [download] = await Promise.all([
    page2.waitForEvent('download'),
    page2.locator('text=getFile').click()
  ]);

  // wait for download to complete
  const path = await download.path();

  const srcMd5 = md5(await fs.promises.readFile(file));
  const dstMd5 = md5(await fs.promises.readFile(path));
  expect(srcMd5).toBe(dstMd5);
});

