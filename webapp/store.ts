import { atom } from "jotai"
import { Manifest } from "./lib/manifest"

const ItemsAtom = atom<Array<Manifest>>([])

export {
  ItemsAtom,
}
