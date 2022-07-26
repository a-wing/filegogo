import { assert, describe, it } from 'vitest'

import VirtualFile from "./file/virtual"

import Pool from "./pool"

describe('file pool test', async () => {
  const nameSend = "send_xxx"
  const size = 1024*1024

  const fileSend = new VirtualFile(nameSend, size)
  const fileRecv = new VirtualFile("recv_xxx", size)

  const sender = new Pool
  sender.setSend(fileSend)

  const meta = await sender.sendMeta()

  it('meta file name', () => {
    assert.equal(nameSend, meta.name)
  })

  it('meta file size', () => {
    assert.equal(size, meta.size)
  })

  const recver = new Pool

  recver.setRecv(fileRecv)
  recver.recvMeta(meta)

  const c = recver.next()

  it('recv next data chunk', () => {
    assert.notEqual(c, null)
  })

  if (c) {
    const data = await sender.sendData(c)
    await recver.recvData(c, data)
  }

  const run = async (): Promise<void> => {
    return new Promise((resolve) => {
      const timer = setInterval(async() => {
        const c = recver.next()
        if (c) {
          const data = await sender.sendData(c)
          await recver.recvData(c, data)
        }
      }, 100)

      recver.OnFinish = () => {
        clearInterval(timer)
        resolve()
      }
    })
  }

  await run()

  const hash = sender.sendHash()
  const ok = recver.recvHash(hash)

  it('check sum', () => {
    assert.equal(ok, true)
  })
})
