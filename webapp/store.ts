import { atom } from "jotai"
import { Manifest } from "./lib/manifest"

const ItemsAtom = atom<Array<Manifest>>([])


interface Config {
  iceServers: RTCIceServer[]
}

const ConfigAtom = atom<Config>

export {
  ItemsAtom,
  ConfigAtom,
}
