# SNE-DRS Command Center: Master Design System

> [!NOTE]
> **From the Desk of the Retiring Lead Designer**
> To whoever inherits this UI: This is not a flashy consumer app. It is an Official Government Dashboard built for disaster response. The aesthetic must remain authoritative, accessible, and highly legible under stressful conditions. I've spent the last few weeks cementing the exact physics, colors, and layout rules to ensure this system remains cohesive. Guard these rules closely.

## 1. Core Aesthetics & Typography

### The Color Palette
We rejected the default dark mode and went with a muted, organic light/dark theme utilizing HSL values. This prevents the UI from looking like a toy and grounds it in reality.
* **Backgrounds & Surfaces:** Uses high-lightness, desaturated greens (`hsl(120, 20%, 97%)`) to reduce eye strain.
* **The "Primary" Split:** 
  * `--primary` (`hsl(132, 43%, 76%)`): A light pastel green. Strictly reserved for buttons and text highlighting.
  * `--primary-dark` (`#2A4724`): A rich forest green. Created because the pastel primary clashed with our grey secondary. Used for interactive states, borders, and tracks to guarantee high contrast.
* **Map Semantic Colors:** Never use raw RGB red or yellow. Use our specific hex codes: `#87231E` (Immediate Risk), `#8F7518` (Moderate), `#354F3E` (Safe), `#156932` (Habitation).

### Typography
* **Font Family:** `Poppins` is mandatory across the entire stack.
* **Weight:** `400` ONLY. We rely on scale for hierarchy, not weight.
* **Color:** All text must be inked in `#0f1510`.
* **Selection:** Text selection (`::selection`) must use the `--primary` background color.

---

## 2. The Primitive Library

### 1. `Button.svelte`
* **Decision:** We use a single, unified button component across the app. 
* **Interaction Physics:** On hover, the background color **does not change**. Instead, the button must physically lift (`translateY(-2px)`) and cast a drop shadow (`0 4px 12px`). This provides tactile feedback without relying on flashy color shifts.
* **Rounding:** `4px` border-radius.

### 2. `Card.svelte` (Data Cards)
* **Decision:** Cards must look like physical, recessed data panels, not floating sheets of paper.
* **Styling:** External drop shadows and colored side-bars are banned. Cards must use an `inset` box shadow, and the background color must be darkened by 4% using `color-mix` to create the recessed depth illusion.

### 3. `Badge.svelte` (Risk Tiers)
* **Decision:** Badges cannot be pill-shaped. They must share the exact `4px` border-radius of the buttons for strict geometrical uniformity.
* **Usage in Cards:** When placed inside a Card, a Badge must be configured as a **Band** (`isBand=true`). It must span 100% of the card's width directly beneath the header, acting as a visual warning stripe.

### 4. `Switch.svelte` & `Checkbox.svelte`
* **Decision:** The pastel `--primary` color was too light for active states and blurred into the grey tracks. You must use `--primary-dark` (`#2A4724`) for the "checked" states.
* **Tick Mark:** The checkmark inside the checkbox must be pure white (`#ffffff`) to pop against the dark green background.

### 5. `Tabs.svelte`
* **Decision:** Unselected tabs must not blend into the background. They use a 60% text-color mix (`color-mix(in srgb, var(--text) 60%, transparent)`) to remain highly visible while deferring to the active tab, which uses a `--primary-dark` underline.

### 6. `Notification.svelte`
* **Decision:** Notifications must behave like critical system alerts, not casual social media toasts.
* **Positioning:** Absolute top-center of the screen.
* **Behavior:** Auto-dismissal is strictly forbidden. The notification must persist on screen until the operator manually dismisses it by clicking the 'X'.
* **Styling:** No thick colored borders. A simple 1px grey border keeps it uniform with the rest of the UI.

### 7. Form Controls (`TextInput`, `Select`, `Slider`, `ProgressBar`)
* **Decision:** Form inputs use a subtle grey border that snaps to `--primary-dark` on focus. The `Slider` thumb must use the dark semantic `--color-habitation` (`#156932`) so it contrasts sharply against the grey track. The `ProgressBar` also uses `--primary-dark` for the fill.

### 8. Extended Data & Form Primitives (`DataTable`, `Textarea`, `SearchInput`, `FileInput`, `RadioGroup`, `ButtonGroup`)
* **Decision:** These extensions strictly follow the core form control styling. Inputs use a subtle grey border snapping to `--primary-dark` on focus. The `FileInput` wraps a hidden native input in the official `Button` to ensure consistent interaction physics.
* **Text Colors:** All placeholder text, secondary icons (like search), and unselected text strictly use a 60% text-color mix (`color-mix(in srgb, var(--text) 60%, transparent)`) rather than `--secondary` to maintain high contrast compliance.

### 9. Structural & Layout Overlays (`Drawer`, `DropdownMenu`, `Popover`)
* **Decision:** Floating elements must not violate the strict "no-float" shadow rules on cards. They utilize a simple 1px `--secondary` border. `Drawer` uses a 40% opaque black backdrop (`rgba(0,0,0,0.4)`) and slides cleanly to hold Data Provenance details.

### 10. Dashboard Feedback & Display (`Callout`, `Metric`, `Tag`, `EmptyState`, `ImageBlock`, `Stepper`)
* **Decision:** Secondary inline alerts (`Callout`) cannot use colored left borders or mix pastels in an illegible way. They use an 85% transparent mix of `--primary` or `--color-moderate` backgrounds. 
* **Empty States & Images:** The `EmptyState` uses a dashed border and `2%` black color-mixed background. `ImageBlock` also uses a `2%` recessed background for citizen-submitted photos.
* **Loading States:** Abrupt spinner pop-ins are discouraged for large blocks. `Skeleton` provides an animated `--secondary` gradient pulse.

### 11. Navigation & Mapping Helpers (`Breadcrumbs`, `Pagination`, `MapLegend`, `Icon`)
* **Decision:** Navigation relies on semantic layout. `Pagination` buttons use standard borders. `Breadcrumbs` use the 60% text mix for separators. `MapLegend` strictly reads from the CSS variables mapped to risk tiers (`--color-redzone`, `--color-moderate`, etc.).

---

## 3. The "Strictly Don't" List

As you take over this system, if you do any of the following, I will come out of retirement just to revert your commit:

1. **DO NOT BOLD TEXT:** Under no circumstances should you use `font-weight: 700` or `<strong>` tags. Hierarchy is established strictly through the mathematical modular scale (font size), not thickness.
2. **DO NOT FLOAT CARDS:** Do not add standard external drop shadows to `Card` components. They are recessed panels.
3. **DO NOT MIX PASTELS:** Do not place `--primary` (light green) directly next to or inside `--secondary` (light grey). They have the same lightness and will turn to visual mud. Use `--primary-dark` for contrast.
4. **DO NOT AUTO-DISMISS ALERTS:** In a disaster response scenario, a user missing a notification because they looked away for 3 seconds is unacceptable. Notifications must be manually dismissed.
5. **DO NOT CHANGE BUTTON COLORS ON HOVER:** Buttons lift. They do not glow or change color.
6. **DO NOT PILL-SHAPE BADGES:** Badges are not pills. They must have a 4px border radius.
7. **DO NOT USE BROWSER DEFAULT SCROLLBARS:** Use the custom floating iOS-style scrollbars defined in `app.css`.
8. **DO NOT USE SECONDARY COLORS FOR TEXT:** All text must derive its ink from the master `--text` variable. If you need muted text (like placeholders or breadcrumb separators), use an opacity `color-mix`, not the raw `--secondary` grey variable.

Good luck, and keep the UI clean.
