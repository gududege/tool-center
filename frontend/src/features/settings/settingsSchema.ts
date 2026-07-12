import type { UiSchema } from '@rjsf/utils'

type TranslateFn = (key: string) => string

export function getSettingsSchema(t?: TranslateFn) {
  const _ = (key: string, fallback: string) => t?.(key) ?? fallback
  return {
    type: "object" as const,
    properties: {
      theme: {
        type: "string" as const,
        title: _('settings.fields.theme', 'Theme'),
        enum: ["dark", "light", "system"],
        enumNames: ["Dark", "Light", "System"],
      },
      formTheme: {
        type: "string" as const,
        title: _('settings.fields.form_theme', 'Form Theme'),
        enum: ["default", "caffeine", "claude", "clean-slate", "amethyst-haze", "neo-brutalism", "pastel-dreams", "soft-pop", "twitter", "vercel"],
        enumNames: ["Default", "Caffeine", "Claude", "Clean Slate", "Amethyst Haze", "Neo Brutalism", "Pastel Dreams", "Soft Pop", "Twitter", "Vercel"],
      },
      language: {
        type: "string" as const,
        title: _('settings.fields.language', 'Language'),
        enum: ["en", "zh"],
        enumNames: ["English", "中文"],
      },
      pluginDirectory: {
        type: "string" as const,
        title: _('settings.fields.plugin_directory', 'Plugin Directory'),
      },
      runnerPaths: {
        type: "object" as const,
        title: _('settings.fields.runner_paths', 'Interpreter Paths'),
        properties: {
          cmd:        { type: "string" as const, title: "CMD" },
          bat:        { type: "string" as const, title: "BAT" },
          powershell: { type: "string" as const, title: "PowerShell" },
          sh:         { type: "string" as const, title: "Shell" },
          python:     { type: "string" as const, title: "Python" },
          nodejs:     { type: "string" as const, title: "Node.js" },
        },
      },
      runnerOptions: {
        type: "object" as const,
        title: _('settings.fields.runner_options', 'Interpreter Options'),
        description: _('settings.fields.runner_options_desc', 'Arguments inserted between interpreter and target script'),
        properties: {
          cmd:        { type: "string" as const, title: "CMD Options" },
          bat:        { type: "string" as const, title: "BAT Options" },
          powershell: { type: "string" as const, title: "PowerShell Options" },
          sh:         { type: "string" as const, title: "Shell Options" },
          python:     { type: "string" as const, title: "Python Options" },
          nodejs:     { type: "string" as const, title: "Node.js Options" },
        },
      },
    },
  }
}

export const settingsUiSchema: UiSchema = {
  pluginDirectory: {
    "ui:widget": "folderPicker",
  },
  runnerPaths: {
    cmd:        { "ui:widget": "filePicker" },
    bat:        { "ui:widget": "filePicker" },
    powershell: { "ui:widget": "filePicker" },
    sh:         { "ui:widget": "filePicker" },
    python:     { "ui:widget": "filePicker" },
    nodejs:     { "ui:widget": "filePicker" },
  },
  "ui:submitButtonOptions": {
    norender: true,
  },
  "ui:order": ["theme", "formTheme", "language", "pluginDirectory", "runnerPaths", "runnerOptions", "*"],
}
