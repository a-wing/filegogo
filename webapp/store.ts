import { atom } from "jotai"
import { Box } from "./libfgg"

const ItemsAtom = atom<Array<Box>>([])

interface Config {
  iceServers: RTCIceServer[]
}

const ConfigAtom = atom<Config>

export {
  ItemsAtom,
  ConfigAtom,
}
