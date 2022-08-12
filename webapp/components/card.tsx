import styles from './card.module.scss'

function Card(props: {name: string, type: string, size: number}) {
  return (
    <div className={ styles.card } >
      <div className={ styles.item } >
        <p className={ `${styles.tag} ${styles.link}` } >Name</p>
        <p className={ `${styles.tag} ${styles.primary}` }>{ props.name }</p>
      </div>
      <div className={ styles.item } >
        <p className={ `${styles.tag} ${styles.warning}` }>Type</p>
        <p className={ `${styles.tag} ${styles.danger}` }>{ props.type }</p>
      </div>
      <div className={ styles.item } >
        <p className={ `${styles.tag} ${styles.success}` }>Size</p>
        <p className={ `${styles.tag} ${styles.info}` }>{ props.size }</p>
      </div>
    </div>
  )
}

export default Card
