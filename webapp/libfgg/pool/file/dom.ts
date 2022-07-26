import { IFile } from "./file"
import { Meta } from "../data"

class DomSendFile implements IFile {
  private file: File
  constructor(file: File) {
    this.file = file
  }
  async getMeta(): Promise<Meta> {
    return {
      name: this.file.name,
      type: this.file.type,
      size: this.file.size,
    }
  }
  async setMeta(_: Meta): Promise<void> {}
  async read(offset: number, length: number): Promise<ArrayBuffer> {
    // https://developer.mozilla.org/en-US/docs/Web/API/Blob/slice#parameters
    // slice(start?: number, end?: number, contentType?: string)
    return this.file.slice(offset, offset + length).arrayBuffer()
  }
  async write(_: number, __: number, ___: ArrayBuffer): Promise<void> {}
}

// Inspired https://github.com/RobinLinus/snapdrop/blob/724f0af576852517bbea96a5d41302719ea514ec/client/scripts/network.js#L476-L505
class DomRecvFile implements IFile {
  private buffer: Array<ArrayBuffer> = []
  private bytesReceived: number = 0
  private meta: Meta = {
    name: "",
    type: "application/octet-stream",
    size: 0
  }

  async getMeta(): Promise<Meta> {
    return this.meta
  }
  async setMeta(meta: Meta): Promise<void> {
    this.meta = meta
  }
  async read(_: number, __: number): Promise<ArrayBuffer> {
    return new ArrayBuffer(0)
  }

  async write(offset: number, length: number, data: ArrayBuffer): Promise<void> {
    if (this.bytesReceived !== offset) return
    if (length !== data.byteLength) return

    this.buffer.push(data)
    this.bytesReceived += data.byteLength
    if (this.bytesReceived < this.meta.size) return

    this.close()
  }
  private close() {
    // we are done
    let blob = new Blob(this.buffer, { type: this.meta.type })

    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = this.meta.name
    document.documentElement.appendChild(a)
    a.click()
    document.documentElement.removeChild(a)
  }
}

export {
  DomSendFile,
  DomRecvFile,
}
