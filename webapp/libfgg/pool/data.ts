interface DataChunk {
  offset: number
  length: number
}

interface Meta {
  name: string
  // TODO: mime
  type: string
  size: number
}

interface Hash {
  name: string
  hash: string
}

export type {
  DataChunk,
  Meta,
  Hash,
}
