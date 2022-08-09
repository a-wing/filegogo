import { test, expect } from '@playwright/test';
import { spawn } from 'child_process';

import { getRoom, checkSum } from './helper';

import { address, file } from './config';

test('cli to browser', async ({ page }) => {
  const share = `${address}/${await getRoom(address)}`;

  const ls = spawn('./filegogo', ['send', '-s', share, file]);

  ls.stdout.on('data', (data) => {
    console.log(`stdout: ${data}`);
  });

  ls.stderr.on('data', (data) => {
    console.log(`stdout: ${data}`);
  });

  ls.on('close', (code) => {
    console.log(`child process exited with code ${code}`);
  });

  await page.goto(share);
  expect(await page.title()).toBe('Filegogo');

  const [download] = await Promise.all([
    page.waitForEvent('download'),
    page.locator('text=Download').click()
  ]);

  const path = await download.path();

  expect(await checkSum(file)).toBe(await checkSum(path));
});
