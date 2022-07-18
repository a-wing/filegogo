import { assert, describe, it } from 'vitest'

import IOVirtual from './virtual'

describe('io virtual test', async () => {
  const data1 = "data1"
  const data2 = "data2"
  const data3 = "data3"

  const [a, b, c] = IOVirtual(3)

  let aD1: any = null
  let aD2: any = null
  a.onmessage = (data: any) => {
    if (aD1 === null) {
      aD1 = data
    } else {
      aD2 = data
    }
  }

  let bD1: any = null
  let bD2: any = null
  b.onmessage = (data: any) => {
    if (bD1 === null) {
      bD1 = data
    } else {
      bD2 = data
    }
  }

  let cD1: any = null
  let cD2: any = null
  c.onmessage = (data: any) => {
    if (cD1 === null) {
      cD1 = data
    } else {
      cD2 = data
    }
  }

  a.send(data1)
  b.send(data2)
  c.send(data3)

  it('a data', () => {
    assert.equal(aD1, data2)
    assert.equal(aD2, data3)
  })

  it('b data', () => {
    assert.equal(bD1, data1)
    assert.equal(bD2, data3)
  })

  it('c data', () => {
    assert.equal(cD1, data1)
    assert.equal(cD2, data2)
  })
})
