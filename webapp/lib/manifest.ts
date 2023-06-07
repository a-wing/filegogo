import { Meta } from "../lib/archive"

interface Item {
  name: string
  // TODO: mime
  type: string
  size: number

  files: Meta[]
}

interface Manifest extends Item {
  uxid: string

  expire: string
  remain: number
}


export type {
  Manifest,
  Item,
}
