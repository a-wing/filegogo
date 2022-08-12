const l1Length = 2
const l2Length = 2

function encode(head: ArrayBuffer, body: ArrayBuffer): ArrayBuffer {
  const data = new Uint8Array(l1Length + l2Length + head.byteLength + body.byteLength)
  const meta = new DataView(data.buffer)
  meta.setUint16(0, head.byteLength)
  meta.setUint16(l1Length, body.byteLength)
  data.set(new Uint8Array(head), l1Length + l2Length)
  data.set(new Uint8Array(body), l1Length + l2Length + head.byteLength)
  return data.buffer
}

function decode(data: ArrayBuffer): ArrayBuffer[] {
  if (data.byteLength < l1Length + l2Length) {
    return [new ArrayBuffer(0), new ArrayBuffer(0)]
  }
  const meta = new DataView(data)
  const l1 = meta.getUint16(0)
  const l2 = meta.getUint16(l1Length)

  if (data.byteLength < l1Length + l2Length + l1 + l2) {
    return [new ArrayBuffer(0), new ArrayBuffer(0)]
  }

  return [
    data.slice(l1Length + l2Length, l1Length + l2Length + l1),
    data.slice(l1Length + l2Length + l1, l1Length + l2Length + l1 + l2),
  ]
}

export {
  encode,
  decode,
}
