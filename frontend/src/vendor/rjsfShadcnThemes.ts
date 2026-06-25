// The RJSF shadcn default sub-theme stylesheet, imported as a raw string.
//
// @rjsf/shadcn ships Tailwind v4 CSS (preflight + @layer theme/base/components/
// utilities + the shadcn CSS variables). Two problems prevent injecting it into
// the main document:
//
//   1. Tailwind v4 directives (@layer base/theme, @property, oklch, color-mix)
//      break this project's Tailwind v3 PostCSS pipeline.
//   2. The sub-theme defines the shadcn variables (--background/--foreground/
//      --primary/…) under :root and .dark — the SAME names this app's index.css
//      defines for the main UI. Injecting globally would override the app's
//      variables and recolor the sidebar/header/panels.
//
// So the default sub-theme CSS is injected inside a Shadow DOM that wraps only
// the RJSF form (see DynamicForm.tsx → ShadowFormHost). The shadow boundary
// contains the Tailwind v4 preflight + variable defs to the form area.
//
// Form/theme unification: the 9 non-default sub-themes and the FormThemePicker
// have been removed. The form now always uses the default sub-theme and follows
// the global light/dark toggle. To make the form's colors match the app exactly
// (the app uses neutral-gray HSL vars in index.css, while the default sub-theme
// uses its own oklch palette), an app-color override stylesheet is injected
// AFTER the default CSS, replacing the shadow's shadcn variables with the app's
// HSL values (wrapped in hsl() since the v4 utilities consume var(--x) directly,
// unlike the app's v3 utilities which wrap with hsl(var(--x))).
//
// Selector rewrite for shadow DOM:
//   - `:root{ --background: … }` (the shadcn color-variable block) → `:host{ … }`,
//     because :root inside a shadow tree does NOT match the shadow host; :host does.
//   - `.dark{ --background: … }` (the dark color-variable block) → `:host(.dark){ … }`,
//     so it applies only when the shadow host carries the .dark class.
//   - `:root,:host{ … }` (Tailwind theme tokens) is left alone — its :host half
//     already works in the shadow.
//
// Dark mode: ShadowFormHost mirrors <html>'s .dark class onto the shadow host
// element, so :host(.dark) flips the form's dark variant in lockstep with the app.

import defaultCss from './rjsf-shadcn-default.css?raw'

// Rewrite :root{} → :host{} and .dark{} → :host(.dark){} for shadow DOM use.
// Only matches the selector immediately before a `{`, so:
//   - `:root,:host{…}` (Tailwind theme tokens) is left alone — its :root is
//     followed by `,` not `{`.
//   - `.dark\:…{` (escaped Tailwind class names) is left alone — the `.dark`
//     there is followed by `\:` not `{`.
function rewriteForShadowDom(css: string): string {
  let out = css.replace(/(^|[^,]):root\{/g, '$1:host{')
  out = out.replace(/(^|[^\\])\.dark\{/g, '$1:host(.dark){')
  return out
}

let rewrittenDefaultCss: string | null = null

// The rewritten default sub-theme CSS (selectors rewritten for shadow DOM).
export function getShadcnDefaultSubThemeCss(): string {
  if (rewrittenDefaultCss) return rewrittenDefaultCss
  rewrittenDefaultCss = rewriteForShadowDom(defaultCss)
  return rewrittenDefaultCss
}

// App color override: replaces the shadow's shadcn color variables with the app's
// HSL values (from src/index.css :root and .dark), wrapped in hsl() because the
// RJSF shadcn Tailwind v4 utilities consume var(--background) directly. Injected
// AFTER the default sub-theme CSS so it wins at equal :host specificity.
// This makes the form area's colors match the rest of the app (neutral gray).
export const appColorOverrideCss = `
:host{
  --background: hsl(0 0% 100%);
  --foreground: hsl(0 0% 3.9%);
  --card: hsl(0 0% 100%);
  --card-foreground: hsl(0 0% 3.9%);
  --popover: hsl(0 0% 100%);
  --popover-foreground: hsl(0 0% 3.9%);
  --primary: hsl(0 0% 9%);
  --primary-foreground: hsl(0 0% 98%);
  --secondary: hsl(0 0% 96.1%);
  --secondary-foreground: hsl(0 0% 9%);
  --muted: hsl(0 0% 96.1%);
  --muted-foreground: hsl(0 0% 45.1%);
  --accent: hsl(0 0% 96.1%);
  --accent-foreground: hsl(0 0% 9%);
  --destructive: hsl(0 84.2% 60.2%);
  --destructive-foreground: hsl(0 0% 98%);
  --border: hsl(0 0% 70%);
  --input: hsl(0 0% 70%);
  --ring: hsl(0 0% 3.9%);
  --radius: 0.5rem;
  /* Tailwind v4 initializes this inside an @supports block that may not apply in
     WebView2 / shadow DOM. Without it, the .border utility sets border-style to
     the initial value (none) and borders disappear. */
  --tw-border-style: solid;
}
:host(.dark){
  --background: hsl(0 0% 3.9%);
  --foreground: hsl(0 0% 98%);
  --card: hsl(0 0% 3.9%);
  --card-foreground: hsl(0 0% 98%);
  --popover: hsl(0 0% 3.9%);
  --popover-foreground: hsl(0 0% 98%);
  --primary: hsl(0 0% 98%);
  --primary-foreground: hsl(0 0% 9%);
  --secondary: hsl(0 0% 14.9%);
  --secondary-foreground: hsl(0 0% 98%);
  --muted: hsl(0 0% 14.9%);
  --muted-foreground: hsl(0 0% 63.9%);
  --accent: hsl(0 0% 14.9%);
  --accent-foreground: hsl(0 0% 98%);
  --destructive: hsl(0 62.8% 30.6%);
  --destructive-foreground: hsl(0 0% 98%);
  --border: hsl(0 0% 30%);
  --input: hsl(0 0% 30%);
  --ring: hsl(0 0% 83.1%);
  --tw-border-style: solid;
}
/* Extra fallback: make sure any element carrying the .border class actually
   renders a solid border, independent of the @layer properties initialization. */
.border {
  border-style: solid;
}
`
