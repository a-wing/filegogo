import PizZip from "pizzip"
import { Item, Meta } from "../libfgg/index"

const ArchiveName = "filegogo-archive.zip"
const ArchiveType = "application/zip"

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

  get name(): string {
    return this.files.length > 1 ? ArchiveName : (this.files.length > 0 ? this.files[0].name : "")
  }

  get size(): number {
    return this.files.reduce((total, file) => total + file.size, 0)
  }

  genManifest(): Item {
    let count = this.files.length
    if (count === 0) {
      throw "not file found"
    } else if (count === 1) {
      return {
        name: this.files[0].name,
        type: this.files[0].type,
        size: this.files[0].size,
        files: [],
      }
    } else {
      return {
        name: ArchiveName,
        type: ArchiveType,
        size: this.size,
        files: this.manifest,
      }
    }
  }

  async exportFile(): Promise<File> {
    if (this.files.length === 0) {
      throw "not found file"
    }

    if (this.files.length === 1) {
      return this.files[0]
    }

    const zip = new PizZip()
    await Promise.all(this.files.map(async f => zip.file(f.name, await f.text())))
    return new File([zip.generate({ type: "blob" })], ArchiveName)
  }

  addFiles(files: File[]) {
    const newFiles = files.filter(
      file => file.size > 0 && !this.isDuplicate(file)
    )

    this.files = this.files.concat(newFiles)
  }

  remove(file: File) {
    const index = this.files.indexOf(file)
    if (index > -1) {
      this.files.splice(index, 1)
    }
  }
}
