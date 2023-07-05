import { atom } from "jotai"
import { Box } from "./libfgg"

const ItemsAtom = atom<Array<Box>>([])
const DetailAtom = atom<Box | null>(null)

export {
  ItemsAtom,
  DetailAtom,
}
