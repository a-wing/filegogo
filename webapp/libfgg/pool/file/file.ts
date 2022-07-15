import { Meta } from "../data"

interface IFile {
  getMeta(): Promise<Meta>
  setMeta(mate: Meta): Promise<void>
  read(offset: number, length: number): Promise<ArrayBuffer>
  write(offset: number, length: number, data: ArrayBuffer): Promise<void>
}

export type {
  IFile,
}
