interface IO {
  send(data: any): void
  onmessage: (data: any) => void
}

class Virtual {
  send(_: any): void {}
  onmessage: (data: any) => void = (_: any) => {}
}

function IOVirtual (num: number): IO[] {
  let io: IO[] = []
  for (let i = 0; i < num; i++) {
    io.push(new Virtual())
  }

  for (let i = 0; i < num; i++) {
    io[i].send = (data: any) => {
      for (let j = 0; j < num; j++) {
        if (i !== j) {
          io[j].onmessage(data)
        }
      }
    }
  }

  return io
}

export default IOVirtual
