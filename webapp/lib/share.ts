function shareGetRoom(addr: string): string {
  const arr = (new URL(addr)).pathname.split('/')
  if (arr.length > 0) {
    const str = arr[arr.length - 1]
    return /\d/.test(str) ? str : ''
  }
  return ''
}

function generateShare(id: string): string {
  return window.location.href + id
}

export {
  shareGetRoom,
  generateShare,
}

