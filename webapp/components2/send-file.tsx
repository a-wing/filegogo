import { useState } from "react"
import { useAtom } from "jotai"
import Copy from "copy-to-clipboard"
import Qrcode from "./qr-code"
import { DetailAtom } from "../store"
import { Box } from "../libfgg/index"
import { generateShare } from "../lib/api"

export default (props: { file: Box }) => {
  const file = props.file
  const address = generateShare(file.uxid)
  const [_, setDetail] = useAtom(DetailAtom)
  const [copyState, setCopyState] = useState(false)

  return (
    <>
      <div className="w-full flex flex-row justify-center">
        <Qrcode address={ address }></Qrcode>
      </div>
      <input id="share-url" type="text" readOnly={ true } value={ address } className="block w-full my-4 border-5px rounded-lg leading-loose h-12 px-2 py-1 dark:bg-grey-80" />
      <button className="p-3 w-full block border-1 rounded-md bg-blue-500 text-white font-bold" onClick={ () => Copy(address) && setCopyState(true) }>
        { copyState ? "Copied" : "Copy Link"}
      </button>
      <button className="p-3 w-full rounded-md text-blue-500 font-bold" onClick={ () => setDetail(null) }>Ok</button>
    </>
  )
}
