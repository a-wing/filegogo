import { assert, describe, it } from 'vitest'

import SparkMD5 from 'spark-md5'
import FileHash from "./hash"

describe('file hash test', async () => {
  const data1 = await new Blob(["aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"]).arrayBuffer()
  const data2 = await new Blob(["bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"]).arrayBuffer()
  const data3 = await new Blob(["cccccccccccccccccccccccccccc"]).arrayBuffer()

  const sum = new SparkMD5.ArrayBuffer()
  sum.append(data1)

  const fh = new FileHash()
  fh.onData({
    offset: 0,
    length: data1.byteLength,
  }, data1)


  it('data1 sum', () => {
    assert.equal(fh.sum(), sum.end())
  })

  sum.append(data2)

  fh.onData({
    offset: data1.byteLength,
    length: data2.byteLength,
  }, data2)


  it('data2 sum', () => {
    assert.equal(fh.sum(), sum.end())
  })

  // duplicate data2
  fh.onData({
    offset: data1.byteLength,
    length: data2.byteLength,
  }, data2)

  it('duplicate data2 sum', () => {
    assert.equal(fh.sum(), sum.end())
  })

  // duplicate data1
  fh.onData({
    offset: 0,
    length: data1.byteLength,
  }, data2)

  it('duplicate data1 and data2 sum', () => {
    assert.equal(fh.sum(), sum.end())
  })

  sum.append(data3)

  fh.onData({
    offset: data1.byteLength + data2.byteLength,
    length: data3.byteLength,
  }, data3)

  it('data3 sum', () => {
    assert.equal(fh.sum(), sum.end())
  })
})
