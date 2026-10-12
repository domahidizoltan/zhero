---
paths:
  - "template/*"
---

# Presentation layer Rules

- Use kebab-case format for HTML tag IDs and class names.
- Add template names to `template/template.go` for backend reuse.
- Use FontAwesome icons; try to avoid generating SVG codes for icons.
- Make responsible layouts.
- Use theme classes what supports dark and light themes.
- Keep minimal JavaScript in `.hbs` files; extract JavaScript functions to `template/<domain>/<domain>.js` files.
- Use native and DaisyUI validators for HTML forms.
- Don't create custom HTML components; use native and DaisyUI components.
- Split complex templates into partials.
- Ensure frontend interactions are primarily driven by HTMX where possible.