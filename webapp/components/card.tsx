import styles from './card.module.scss'

function Card(props: {name: string, type: string, size: number}) {
  return (
    <div className={ styles.card } >
      <div className={ styles.item } >
        <label className={ `${styles.tag} ${styles.link}` } >Name</label>
        <label className={ `${styles.tag} ${styles.primary}` }>{ props.name }</label>
      </div>
      <div className={ styles.item } >
        <label className={ `${styles.tag} ${styles.warning}` }>Type</label>
        <label className={ `${styles.tag} ${styles.danger}` }>{ props.type }</label>
      </div>
      <div className={ styles.item } >
        <label className={ `${styles.tag} ${styles.success}` }>Size</label>
        <label className={ `${styles.tag} ${styles.info}` }>{ props.size }</label>
      </div>
    </div>
  )
}

export default Card
