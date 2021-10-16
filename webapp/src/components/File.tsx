import { useRef, useState, ChangeEvent } from 'react'
import styles from './File.module.scss'

function File(props: {
  recver: boolean,
  percent: number,
  handleFile: (files: FileList | null) => void,
  getFile: () => void }) {

  const hiddenFileInput = useRef<HTMLInputElement>(null)
  const handleClick = () => {
    props.recver
    ? props.getFile()
    : hiddenFileInput.current?.click?.()
  }

  const [filename, setFilename] = useState('Select File')

  const handleFile = (files: FileList | null) => {
    props.handleFile(files)
    const file = files?.[0]
    if (file) setFilename(file.name)
  }

  return (
    <>
      <label className={ styles.button } style={{
        background: 'linear-gradient(to right, #f14668 '+ props.percent +'%, #3ec46d '+ props.percent +'%)'
      }} onClick={ handleClick } >{ props.percent === 0 ? (props.recver ? 'getFile' : filename ) : props.percent.toFixed(1) + '%' }</label>
      <input
        className={ styles.input }
        type="file"
        ref={ hiddenFileInput }
        onChange={ (ev: ChangeEvent<HTMLInputElement>) => { handleFile(ev.target.files) } }
      />
    </>
  )
}

export default File
