import { test, expect } from '@playwright/test';
import { spawn } from 'child_process';

import fs from 'fs';

import { getRoom, checkSum, genRandomId } from './helper';

import { address, file } from './config';

test('browser to cli', async ({ page }) => {
  const share = `${address}/${await getRoom(address)}`;
  await page.goto(share);
  expect(await page.title()).toBe('Filegogo');

  await page.setInputFiles('input#upload', file);

  const path = '/tmp/filegogo-e2e-tmp-' + genRandomId();

  const ls = spawn('./filegogo', ['recv', '-s', share, path]);

  ls.stdout.on('data', (data) => {
    console.log(`stdout: ${data}`);
  });

  ls.stderr.on('data', (data) => {
    console.log(`stdout: ${data}`);
  });

  await (new Promise<void>((resolve) => {
    ls.on('close', (code) => {
      console.log(`child process exited with code ${code}`);
      resolve()
    });
  }))

  expect(await checkSum(file)).toBe(await checkSum(path));

  await fs.promises.rm(path)
});
