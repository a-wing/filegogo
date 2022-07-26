import { IFile } from "./file"
import { Meta } from "../data"

export default class VirtualFile implements IFile {
  name: string
  size: number
  constructor(name: string, size: number) {
    this.name = name
    this.size = size
  }
  async getMeta(): Promise<Meta> {
    return {
      name: this.name,
      size: this.size,
      type: "text/plain",
    }
  }
  async setMeta(_: Meta): Promise<void> {}
  async read(_: number, length: number): Promise<ArrayBuffer> {
    return (new TextEncoder()).encode(makeid(length)).buffer
  }
  async write(_: number, __: number, ___: ArrayBuffer): Promise<void> {}
}

// https://stackoverflow.com/questions/1349404/generate-random-string-characters-in-javascript
function makeid(length: number) {
  let result           = ''
  let characters       = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789'
  let charactersLength = characters.length
  for (let i = 0; i < length; i++) {
    result += characters.charAt(Math.floor(Math.random() *
      charactersLength))
  }
  return result
}
