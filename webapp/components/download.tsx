import styles from './upload.module.scss'

function Download(props: { name: string }) {
  const send = () => {
    window.open(`/raw/${props.name}`)
  }

  return (
    <>
      <button className={ styles.button } onClick={ send } >Download File</button>
    </>
  )
}

export default Download
