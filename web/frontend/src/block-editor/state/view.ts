import { State } from '@/common/rx'

type ViewStateData = {
  editorWidth: string
}

export class ViewState extends State<ViewStateData> {
  constructor(defaultEditorWidth: string) {
    super({
      editorWidth: defaultEditorWidth,
    })
  }

  setEditorWidth(width: string) {
    this.set({
      editorWidth: width,
    })
  }
}
