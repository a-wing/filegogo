import copy from 'copy-to-clipboard'
import './Address.css'

function Address(props: {address: string}) {
  return (
    <div className="App-address">
      <p className="App-address-text" >{ props.address }</p>
      <button className="App-address-button" onClick={ () => { copy(props.address) } } >COPY</button>
    </div>
  )
}

export default Address
