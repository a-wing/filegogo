// Firefox, Safari
import 'web-streams-polyfill/dist/polyfill.min.js'

import streamSaver from 'streamsaver'

import { IFile } from "./file"
import { Meta } from "../data"

class StreamRecvFile implements IFile {
  // https://developer.mozilla.org/en-US/docs/Web/API/WritableStream/getWriter
  //recver: File | FileDigester | WritableStreamDefaultWriter | null = null
  private file: WritableStreamDefaultWriter | null = null
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
    this.file = streamSaver.createWriteStream(meta.name, {
      size: meta.size,
      //mitm: meta.type,
    }).getWriter()
  }
  async read(_: number, __: number): Promise<ArrayBuffer> {
    return new ArrayBuffer(0)
  }

  async write(offset: number, length: number, data: ArrayBuffer): Promise<void> {
    if (this.bytesReceived !== offset) return
    if (length !== data.byteLength) return

    await this.file?.write(data)
    this.bytesReceived += data.byteLength
    if (this.bytesReceived < this.meta.size) return

    this.close()

  }
  private close() {
    this.file?.close()
  }
}

export {
  StreamRecvFile,
}
