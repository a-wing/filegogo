

// Firefox, Safari
import 'web-streams-polyfill/dist/polyfill.min.js'


import streamSaver from 'streamsaver'

import SparkMD5 from 'spark-md5'

import log from 'loglevel'

interface metaFile {
  file: string
  type: string
  size: number
}

interface metaHash {
  file: string
  hash: string
}

export default class Transfer {
  //metaFile: metaFile
  //metaHash: metaHash
  metaFile: any
  metaHash: any

  file: any
  hash: any

  step: number
  count: number

  onComplete: () => void
  onProgress: (c: number) => void

  complete: boolean
  constructor() {
    // htmlDOMfile
    //this.file = file
    this.file = null
    this.hash = new SparkMD5.ArrayBuffer()

    // safari default
    this.step = 1024 * 64
    // chrome, firefox max-message-size
    // step: 1024 * 256
    this.count = 0

    this.onComplete = () => {}
    this.onProgress = () => {}

    this.metaFile = {}
    this.metaHash = {}

    this.complete = false
  }

  send(file: any) {
    this.file = file
  }

  setMetaFile(meta: metaFile) {
    log.warn(meta)
    let filename = meta.file
    if (meta.file.split("/").length > 0) {
      filename = String(meta.file.split("/").pop())
    }
    this.file = streamSaver.createWriteStream(filename, {
      size: meta.size,
      //mitm: meta.type
    }).getWriter()

    this.metaFile = meta
  }

  //getMetaFile() metaFile {
  getMetaFile() {
    this.metaFile = {
      file: this.file.name,
      type: this.file.type,
      size: this.file.size,
    }
    return this.metaFile
  }


  // => string
  getHash() { return this.hash.end() }

  //getMetaHash() *MetaHash {
  getMetaHash() {
    this.metaHash = {
      file: this.file.name,
      hash: this.getHash()
    }
    return this.metaHash
  }

  // => bool
  verifyHash(meta: metaHash) {
    return meta.hash === this.getHash()
  }

  read(callback: (b: any)=>void, complete: ()=>void) {
    if (this.complete) {
      complete()
    } else {
    const p = this.count
    this.file.slice(p, p + this.step).arrayBuffer()
      .then((buffer: any) => {
        // Md5
        this.hash.append(buffer)
        const c = buffer.byteLength
        this.onProgress(c)
        if (this.count >= this.metaFile.size) {
          this.onComplete()
          this.complete = true
          //this.file.close()
          complete()
        }
        this.count += c

        callback(buffer)
      })
      // @ts-ignore
      .catch(err => console.log(err))
    }
  }

  async write(buffer: ArrayBuffer) {
    if (this.complete) {
      this.onComplete()
      return
    }

    // Md5
    this.hash.append(buffer)
    await this.file.write(new Uint8Array(buffer))
    const c = buffer.byteLength
    this.onProgress(c)
    this.count += c
    if (this.count >= this.metaFile.size) {
      this.onComplete()
      this.complete = true
      this.file.close()
    }
  }
}
