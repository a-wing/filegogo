import { atom } from "jotai"
import { Item } from "./lib/manifest"

const ItemsAtom = atom<Array<Item>>([{
    name: "demo",
    type: "application/x-demo",
    size: 123,
    files: [],
  }, {
    name: "demo-zip",
    type: "application/zip",
    size: 1234,
    files: [
      {
        name: "demo-1",
        type: "application/x-demo",
        size: 123,
      }, {
        name: "demo-2",
        type: "application/x-demo",
        size: 1234,
      },
    ],
  }
])

export {
  ItemsAtom,
}
