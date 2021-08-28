import QRCode from 'qrcode'

import styles from './QRCode.module.scss'

import { useRef } from 'react'

function AppQRCode(props: {address: string}) {
  const qrcode = useRef(null)
  QRCode.toCanvas(qrcode.current, props.address, {
    width: 300
  }, error => {
    if (error) console.error(error)
    console.log('Create QRCode:', props.address)
  })

  return (
    <canvas className={ styles.qrcode } ref={ qrcode }></canvas>
  )
}

export default AppQRCode
