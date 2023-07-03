import { Meta } from "./pool/data"

interface Item {
  name: string
  // TODO: mime
  type: string
  size: number

  files: Meta[]
}

interface Box extends Item {
  uxid: string

  secret?: string
  action: string
  expire: string
  remain: number
}


export type {
  Box,
  Item,
}
