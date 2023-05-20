
interface Meta {
  name: string
  // TODO: mime
  type: string
  size: number
}

export type {
  Meta
}

export default class Archive {
  //private
  files: File[]
  constructor() {
    this.files = Array<File>()
  }

  private isDuplicate(newFile: File): boolean {
    for (const file of this.files) {
      if (
        newFile.name === file.name &&
        newFile.size === file.size &&
        newFile.lastModified === file.lastModified
      ) {
        return true
      }
    }
    return false
  }

  get manifest(): Meta[] {
    return this.files.map(file => ({
      name: file.name,
      size: file.size,
      type: file.type,
    }))
  }

  addFiles(files: File[]) {
    const newFiles = files.filter(
      file => file.size > 0 && !this.isDuplicate(file)
    )

    //const newSize = newFiles.reduce((total, file) => total + file.size, 0);
    this.files = this.files.concat(newFiles);
  }

  remove(file: File) {
    const index = this.files.indexOf(file)
    if (index > -1) {
      this.files.splice(index, 1)
    }
  }
}
