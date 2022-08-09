import { test, expect } from '@playwright/test';

import { getRoom, checkSum } from './helper';

import { address, file } from './config';

test('browser to browser', async ({ page, context }) => {
  const share = `${address}/${await getRoom(address)}`;
  await page.goto(share);
  expect(await page.title()).toBe('Filegogo');

  await page.setInputFiles('input#upload', file);

  const page2 = await context.newPage();
  await page2.goto(share);

  const [download] = await Promise.all([
    page2.waitForEvent('download'),
    page2.locator('text=Download').click()
  ]);

  // wait for download to complete
  const path = await download.path();

  expect(await checkSum(file)).toBe(await checkSum(path));
});
