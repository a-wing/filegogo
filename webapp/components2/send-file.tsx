import Qrcode from "./qr-code"
import { Manifest } from "../lib/manifest"
import { generateShare } from "../lib/api"

export default (props: { file: Manifest, callback: (_: Manifest | null) => void }) => {
  const file = props.file

  return (
    <>
      <div className="w-full flex flex-row justify-center">
        <Qrcode address={ generateShare(file.uxid) }></Qrcode>
      </div>
      <input id="share-url" type="text" readOnly={ true } value={ generateShare(file.uxid) } className="block w-full my-4 border-5px rounded-lg leading-loose h-12 px-2 py-1 dark:bg-grey-80" />
      <button className="p-3 w-full block border-1 rounded-md bg-blue-500 text-white font-bold" onClick={ () => props.callback(null) }>Ok</button>
    </>
  )
}
