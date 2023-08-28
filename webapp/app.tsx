import { useAtom } from "jotai"
import Header from "./components2/header"
import SendArea from "./components2/send-area"
import ArchiveArea from "./components2/archive-area"
import Recver from "./components2/recver"
import { shareGetRoom } from "./lib/share"
import { ItemsAtom } from "./store"

function App() {
  const [files] = useAtom(ItemsAtom)

  return (
    <>
      <Header/>
      <main className="container mx-auto rounded-xl shadow-xl p-8 flex">
        { !!shareGetRoom(window.location.href)
          ? <>
              <Recver/>
            </>
          : <div className="w-full h-full grid grid-cols-1 md:grid-cols-2">
              <section>
                <SendArea/>
              </section>
              <section>
                { !files.length
                  ? <div className="m-8" >
                      <h1 className="text-4xl font-bold">Simple, private file sharing</h1>
                      <p>Filegogo lets you share files with a link that automatically expires. So you can keep what you share private and make sure your stuff doesn’t stay online forever.</p>
                    </div>
                  : <ArchiveArea/>
                }
              </section>
            </div>
        }
      </main>
    </>
  )
}

export default App
