import Header from "./components2/header"
import SendArea from "./components2/send-area"
import ArchiveArea from "./components2/archive-area"
import Recver from "./components2/recver"

function App() {
  return (
    <>
      <Header/>
      <main className="container mx-auto rounded-xl shadow-xl p-8 flex">
        { false
          ? <>
              <Recver/>
            </>
          : <>
              <section className="w-full h-full">
                <SendArea/>
              </section>
              <section className="w-full h-full">
                { false
                  ? <div className="m-8" >
                      <h1 className="text-4xl font-bold">Simple, private file sharing</h1>
                      <p>Send lets you share files with end-to-end encryption and a link that automatically expires. So you can keep what you share private and make sure your stuff doesn’t stay online forever.</p>
                    </div>
                  : <ArchiveArea/>
                }
              </section>
            </>
        }
      </main>
    </>
  )
}

export default App
