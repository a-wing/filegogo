import { assert, describe, it } from 'vitest'

import { encode, decode } from './protocol'

describe('protocol encode decode test', async () => {
  const sHead = await new Blob(["aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"]).arrayBuffer()
  const sBody = await new Blob(["bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"]).arrayBuffer()

  const data = encode(sHead, sBody)
  const [rHead, rBody] = decode(data)

  it('length should be equal', () => {
    assert.equal(sHead.byteLength, rHead.byteLength)
    assert.equal(sBody.byteLength, rBody.byteLength)
  })

  it('payload should be equal', () => {
    const uint8SHead = new Uint8Array(sHead)
    const uint8RHead = new Uint8Array(rHead)
    for (let i = 0; i < sHead.byteLength; i++) {
      assert.equal(uint8SHead[i], uint8RHead[i])
    }

    const uint8SBody = new Uint8Array(sHead)
    const uint8RBody = new Uint8Array(rHead)
    for (let i = 0; i < sBody.byteLength; i++) {
      assert.equal(uint8SBody[i], uint8RBody[i])
    }
  })
})
