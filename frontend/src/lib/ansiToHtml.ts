import anser from 'anser'

export function ansiToHtml(text: string): string {
  return anser.ansiToHtml(anser.escapeForHtml(text), { use_classes: true })
}
