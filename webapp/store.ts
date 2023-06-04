import { atom } from "jotai"
import { Item } from "./lib/manifest"

interface Manifest extends Item {
  uxid: string
}

const ItemsAtom = atom<Array<Manifest>>([])

export {
  ItemsAtom,
}

export type {
  Manifest,
}
