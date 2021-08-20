
function IsShareInit(addr: string): boolean {
  const arr = (new URL(addr)).pathname.split('/')
  if (arr.length > 0) {
    return /\d/.test(arr[arr.length - 1])
  }
  return false
}

export {
  IsShareInit,
}

