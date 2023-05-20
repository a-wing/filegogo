import { Meta } from "../lib/archive"

interface Item {
  name: string
  // TODO: mime
  type: string
  size: number

  files: Meta[]
}

export type {
  Item
}
