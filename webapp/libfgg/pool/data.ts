interface DataChunk {
  offset: number
  length: number
}

interface Meta {
  file: string
  type: string
  size: number
}

interface Hash {
  file: string
  hash: string
}

interface IFile {
  getMeta(): Promise<Meta>
  setMeta(mate: Meta): Promise<void>
  read(offset: number, length: number): Promise<ArrayBuffer>
  write(offset: number, length: number, data: ArrayBuffer): Promise<void>
}

export type {
  DataChunk,
  Meta,
  Hash,
  IFile,
}
