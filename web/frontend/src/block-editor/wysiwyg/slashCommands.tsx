import { WidgetsService } from '../widgets'
import { Widget } from '../widgets/core'
import { SlashCommandItem } from './Suggestions'

export class SlashCommandsProvider {
  private _commands: SlashCommandItem[] = []
  private _widgetService: WidgetsService

  get(): SlashCommandItem[] {
    return this._commands
  }

  constructor(widgetService: WidgetsService) {
    this._widgetService = widgetService

    // TODO error handling
    this.load()
  }

  private async load() {
    const widgets = await this._widgetService.getWidgets()
    this._commands = widgets.map(createSlashCommandFromWidget)
  }
}

function createSlashCommandFromWidget(widget: Widget): SlashCommandItem {
  return {
    name: widget.name,
    description: widget.description,
    command: widget.apply.bind(widget),
  }
}
