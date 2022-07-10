import SparkMD5 from 'spark-md5'

import { DataChunk } from './data'

export default class FileHash {
  hash: SparkMD5.ArrayBuffer = new SparkMD5.ArrayBuffer()
  offset: number = 0

  onData(c: DataChunk, data: ArrayBuffer) {
    if (this.offset === c.offset) {
      this.hash.append(data)
      this.offset += data.byteLength
    }
  }

  sum(): string {
    return this.hash.end()
  }
}
