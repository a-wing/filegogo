import QRCode from "qrcode"
import { useEffect, useRef } from "react"

function AppQRCode(props: {address: string}) {
  const qrcode = useRef<HTMLCanvasElement>(null)

  useEffect(() => {
    QRCode.toCanvas(qrcode.current, props.address, {
      width: 300
    }, error => {
      if (error) console.error(error)
      console.log('Create QRCode:', props.address)
    })
  }, [props.address])

  return (
    <canvas className='rounded-lg' ref={ qrcode }></canvas>
  )
}

export default AppQRCode
