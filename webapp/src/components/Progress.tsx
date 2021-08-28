
function Progress(props: {percent: number}) {
  return (
    <progress max={ 100 } value={ props.percent } ></progress>
  )
}

export default Progress
