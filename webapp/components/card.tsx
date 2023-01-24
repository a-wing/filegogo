import styles from './card.module.scss'

import moment from 'moment'

function Card(props: {
    name: string,
    type: string,
    size: number,
    remain?: number,
    expire?: string,
  }) {
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
      { props.remain && props.expire
        ? <div style={{
            display: 'flex',
            flexDirection: 'row',
            justifyContent: 'space-around',
          }}>
            <div style={{
              display: 'flex',
              flexDirection: 'row',
            }}>
              <p>Remain: </p>
              <p>{ props.remain }</p>
            </div>
            <div style={{
              display: 'flex',
              flexDirection: 'row',
            }}>
              <p>Expire: </p>
              <p>{ moment.duration(moment(props.expire).valueOf() - moment().valueOf()).humanize() }</p>
            </div>
          </div>
        : null
      }
    </div>
  )
}

export default Card
