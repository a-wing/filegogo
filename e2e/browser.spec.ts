import { test, expect } from '@playwright/test';

import { checkSum } from './helper';

import { address, file } from './config';

test('browser to browser', async ({ page, context }) => {
  await page.goto(address);
  expect(await page.title()).toBe('Filegogo');

  await page.setInputFiles('input#upload', file);
  await page.getByText("Commit").click();
  const share = await page.locator("#share-url").inputValue();

  const page2 = await context.newPage();
  await page2.goto(share);

  const [download] = await Promise.all([
    page2.waitForEvent('download'),
    page2.getByRole("button").first().click(),
  ]);

  // wait for download to complete
  const path = await download.path();

  expect(await checkSum(file)).toBe(await checkSum(path));
});
