# Vendor files

## rjsf-shadcn-default.css

Copied from `@rjsf/shadcn/dist/styling/themes/default.css` (v6.6.2) — the
**default** shadcn sub-theme. The 9 other sub-themes (caffeine, claude, …) and
the `FormThemePicker` that switched between them have been removed: the form now
always uses the default sub-theme and follows the global light/dark toggle, with
its colors overridden to match the app (see `appColorOverrideCss` below).

**Why it's vendored:** `@rjsf/shadcn` is built with Tailwind CSS v4, but this
project uses Tailwind v3. The v4 CSS uses `@layer base`/`@layer theme` and other
v4-only directives that Tailwind v3's PostCSS plugin rejects at build time
(`@layer base is used but no matching @tailwind base directive is present`).

The package's `package.json` `exports` field also doesn't expose the CSS files,
so they can't be imported directly via the package name.

**Why it's injected into a Shadow DOM, not the main document:** the sub-theme CSS
defines the shadcn CSS variables (`--background`, `--foreground`, `--primary`,
`--card`, `--border`, …) under `:root` and `.dark`. This app's own `src/index.css`
defines the **same variable names** under `:root`/`.dark` for the main UI. If the
sub-theme stylesheet were injected into the main document, its `:root`/`.dark`
definitions would override the app's and recolor the sidebar, header and panels.

To avoid that, `src/features/workspace/DynamicForm.tsx` renders the RJSF shadcn
form inside a Shadow DOM (`ShadowFormHost`) and injects the default sub-theme's
CSS into that shadow root. The shadow boundary keeps the Tailwind v4 preflight
and variable definitions fully contained to the form area.

**App-color override (`appColorOverrideCss`):** the default sub-theme uses its own
oklch palette, which differs from the app's neutral-gray HSL palette in
`index.css`. To unify the form with the rest of the app, `rjsfShadcnThemes.ts`
also exports `appColorOverrideCss` — a stylesheet that redefines the shadow's
shadcn variables with the app's HSL values (wrapped in `hsl()`, since the v4
utilities consume `var(--background)` directly, unlike the app's v3 utilities
which wrap with `hsl(var(--background))`). `ShadowFormHost` injects this AFTER
the default sub-theme CSS, so it wins at equal `:host` specificity and the form
area's colors match the app exactly.

**Dark mode:** the sub-theme CSS keys its dark variant off a `.dark` class. The
shadow root doesn't inherit `<html>`'s class, so `ShadowFormHost` mirrors the
app's dark state onto the shadow host element. Flipping the app theme flips the
form's dark variant in lockstep — no separate dark toggle for the form.

**How it's imported:** `rjsfShadcnThemes.ts` imports the file with the `?raw`
suffix (so PostCSS skips it) and exposes `getShadcnDefaultSubThemeCss()` (with
`:root`/`.dark` selectors rewritten to `:host`/`:host(.dark)` for shadow DOM) and
`appColorOverrideCss`. The CSS is kept as raw strings and only injected into the
form's shadow root.

**To update:** re-copy from `node_modules/@rjsf/shadcn/dist/styling/themes/default.css`
after upgrading `@rjsf/shadcn`. If the app's palette in `src/index.css` changes,
update `appColorOverrideCss` in `rjsfShadcnThemes.ts` to match.
