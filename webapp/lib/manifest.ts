import { Meta } from "../lib/archive"

interface Item {
  name: string
  // TODO: mime
  type: string
  size: number

  files: Meta[]

  expire: string
  remain: number
}

export type {
  Item
}
