import styles from './Progress.module.scss'

function Progress(props: {percent: number}) {
  return (
    <progress className={ styles.progress } max={ 100 } value={ props.percent } ></progress>
  )
}

export default Progress
