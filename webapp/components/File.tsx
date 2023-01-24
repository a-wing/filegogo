import { useRef, useState, ChangeEvent } from 'react'
import styles from './File.module.scss'
import { putBoxFile, getBoxFile, delBoxFile } from '../lib/api'

let tmp: File | undefined

function File(props: {
  recver: boolean,
  isBox: boolean,
  reLoad: () => void,
  percent: number,
  handleFile: (files: FileList | null) => void,
  getFile: () => void }) {

  const hiddenFileInput = useRef<HTMLInputElement>(null)
  const handleClick = () => {
    if (props.recver) {
      if (props.isBox) {
        getBoxFile()
      } else {
        props.getFile()
      }
    } else {
      hiddenFileInput.current?.click?.()
    }

  }

  const [remain, setRemain] = useState<number>(1)
  const [expire, setExpire] = useState<string>('5m')

  const [filename, setFilename] = useState('Select File')

  const handleFile = (files: FileList | null) => {
    props.handleFile(files)
    const file = files?.[0]
    if (file) setFilename(file.name)
  }

  return (
    <>

    { !props.recver
      ? <div style={{
            display: 'flex',
            flexDirection: 'row',
            justifyContent: 'space-around',
        }}>
        <div>
          <label>Remain: </label>
          <select
            value={remain}
            onChange={e => setRemain(Number(e.target.value))}
          >
            <option value="1">1</option>
            <option value="3">3</option>
            <option value="5">5</option>
            <option value="7">7</option>
            <option value="11">11</option>
          </select>
        </div>

        <div>
          <label>Expire: </label>
          <select
            value={expire}
            onChange={e => setExpire(e.target.value)}
          >
            <option value="5m">5m</option>
            <option value="30m">30m</option>
            <option value="1h">1h</option>
            <option value="24h">24h</option>
          </select>
        </div>
        </div>
      : null
    }

      <label className={ styles.button } style={{
        background: 'linear-gradient(to right, #f14668 '+ props.percent +'%, #3ec46d '+ props.percent +'%)'
      }} onClick={ handleClick } >{ props.percent === 0 ? (props.recver ? 'Download' : filename ) : props.percent.toFixed(1) + '%' }</label>
      <input
        className={ styles.input }
        // This id e2e test need
        id="upload"
        type="file"
        ref={ hiddenFileInput }
        onChange={ (ev: ChangeEvent<HTMLInputElement>) => { handleFile(ev.target.files); tmp = ev.target.files?.[0] } }
      />

      { props.recver
        ? <button
            className={ styles.button }
            style={{ backgroundColor: "darksalmon" }}
            onClick={ async () => {
              await delBoxFile()
              props.reLoad()
            }}
          >Clear</button>
        : <button
            className={ styles.button }
            style={{ backgroundColor: "deepskyblue" }}
            onClick={ async () => {
              if (tmp) {
                await putBoxFile(tmp, remain, expire)
                props.reLoad()
              } }}
          >Relay Box</button>
      }
    </>
  )
}

export default File
