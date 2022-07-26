import FileHash from "./hash"
import { DataChunk, Meta, Hash } from "./data"
import { IFile } from "./file/file"

export default class Pool {
  // htmlDOMfile
  //sender: File | null = null
  sender: IFile | null = null

  // https://developer.mozilla.org/en-US/docs/Web/API/WritableStream/getWriter
  //recver: File | FileDigester | WritableStreamDefaultWriter | null = null
  recver: IFile | null = null

  fileHash: FileHash = new FileHash()

  meta: Meta | null = null
  hash: Hash | null = null

  doneCount: number = 0
  nextCount: number = 0

  OnFinish: () => void = () => {}
  OnProgress: (c :number) => void = (_) => {}

  // [safari] max-message-size: 64 * 1024
  // [chrome, firefox] max-message-size: 256 * 1024
  chunkSize: number = 32 * 1024

  currentSize: number = 0
  pendingSize: number = 0

  setSend(file: IFile) {
    this.sender = file
  }

  setRecv(file: IFile) {
    this.recver = file
  }

  async recvMeta(meta: Meta): Promise<void> {
    this.meta = meta
    await this.recver?.setMeta(meta)
  }

  async sendMeta(): Promise<Meta> {
    if (!this.sender) {
      throw "Not found sender file"
    }

    const meta = await this.sender.getMeta()
    this.meta = meta
    return meta
  }

  sendHash(): Hash {
    if (!this.sender) {
      throw "Not found sender file"
    }

    return {
      name: this.meta?.name || "",
      hash: this.fileHash.sum(),
    }
  }

  recvHash(hash: Hash): boolean {
    return hash.hash === this.fileHash.sum()
  }

  async sendData(c: DataChunk): Promise<ArrayBuffer> {
    if (!this.sender) {
      throw "Not found sender file"
    }

    const data = await this.sender.read(c.offset, c.length)

    this.fileHash.onData(c, data)
    this.OnProgress(this.fileHash.offset)
    return data
  }

  async recvData(c: DataChunk, data: ArrayBuffer): Promise<void> {
    if (!this.recver) {
      throw "Not found recver file"
    }

    this.currentSize += c.length
    this.fileHash.onData(c, data)
    this.OnProgress(this.fileHash.offset)

    return this.recver.write(c.offset, c.length, data)
  }

  next(): DataChunk | null {
    if (!this.meta) {
      throw "Not found recver file"
    }

    if (this.currentSize >= this.meta.size) {
      this.OnFinish()
      return null
    }

    if (this.pendingSize >= this.meta.size) {
      return null
    }

    let length = this.chunkSize
    const next = this.currentSize + this.chunkSize
    if (next > this.meta.size) {
      const n = next - this.meta.size
      length = this.chunkSize - n
    }

    const offset = this.pendingSize

    this.pendingSize += this.chunkSize
    return {
      offset: offset,
      length: length,
    }
  }
}
