interface DataChunk {
  offset: number
  length: number
}

interface Meta {
  // TODO: name
  file: string
  // TODO: mime
  type: string
  size: number
}

interface Hash {
  file: string
  hash: string
}

export type {
  DataChunk,
  Meta,
  Hash,
}
