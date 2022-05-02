import { test, expect } from '@playwright/test';

import { getRoom, checkSum } from './helper';

const address = 'http://localhost:8080';
//const address = 'https://send.22333.fun';
const file = 'playwright.config.ts';

test('browser to browser', async ({ page, context }) => {
  const share = `${address}/${await getRoom('http://localhost:8080')}`;
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

  expect(await checkSum(file)).toBe(await checkSum(path));
});
