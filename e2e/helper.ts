import http from 'http';

import fs from 'fs';
import md5 from 'md5';

const path = '/api/signal/'

async function getRoom(address: string): Promise<string> {
  return new Promise<string>((resolve) => {
    http.get(address + path, res => {
      let body = '';
      res.on('data', chunk => { body += chunk });
      res.on('end', () => {
        const data = JSON.parse(body);
        resolve(data.room);
      })
    })
  })
}

async function checkSum(path: string): Promise<string> {
  const file = await fs.promises.readFile(path)
  return new Promise<string>((resolve) => {
    resolve(md5(file));
  })
}

function genRandomId(): string {
  return Math.random().toString(36).slice(-6)
}

export {
  getRoom,
  checkSum,
  genRandomId,
}
