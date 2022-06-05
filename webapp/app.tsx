import { useEffect, useState } from 'react'
import styles from './app.module.scss'
import { use100vh } from 'react-div-100vh'

import { getLogLevel, getRoom, shareGetRoom } from './lib/api'
import log, { LogLevelDesc } from 'loglevel'
import history from 'history/browser'

import Index from './components/index'

function App() {
  const [address, setAddress] = useState<string>()

  useEffect(() => {
    log.setLevel(getLogLevel() as LogLevelDesc)

    if (shareGetRoom(window.location.href)) {
      setAddress(window.location.href)
    } else {
      const init = async function() {
        const room = await getRoom()
        if (room) {
          history.push(room)
          setAddress(document.location.origin + '/' + room)
        }
      }

      init()
    }
  }, [])

  return (
    <div className={ styles.app } style={{ height: use100vh() || '100vh' }}>
      <div className={ styles.card }>
        { address
          ? <Index address={ address }></Index>
          : <div>Not Get Share Link</div>
        }
      </div>
    </div>
  )
}

export default App
