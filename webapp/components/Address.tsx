import copy from 'copy-to-clipboard'
import styles from './Address.module.scss'

function Address(props: {address: string}) {
  return (
    <div className={ styles.address }>
      <p className={ styles.text } >{ props.address }</p>
      <button className={ `${styles.base} ${styles.button}` } onClick={ () => { copy(props.address) } } >COPY</button>
    </div>
  )
}

export default Address
