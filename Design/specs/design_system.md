# HelixCode Design System

## Design Philosophy

**Goal**: Create the most stunning, enterprise-quality design system that surpasses all existing AI development tools in user experience, aesthetics, and functionality.

## Core Design Principles

### 1. Uninterrupted Workflows
- Running processes continue seamlessly in background
- Real-time updates without blocking user interaction
- Non-intrusive notifications and status indicators
- Smooth transitions between different modes

### 2. Stunning Theme System
- Multiple breathtaking themes that make choice difficult
- Perfect balance between aesthetics and functionality
- System theme detection with manual override
- Consistent experience across all platforms

### 3. Enterprise Quality
- Professional appearance suitable for corporate environments
- Accessibility compliance (WCAG 2.1 AA)
- Responsive design for all screen sizes
- Intuitive navigation and information hierarchy

## Color System

### Primary Color Palette

#### Deep Ocean (Default Dark)
```css
:root {
  --primary-50: #e3f2fd;
  --primary-100: #bbdefb;
  --primary-200: #90caf9;
  --primary-300: #64b5f6;
  --primary-400: #42a5f5;
  --primary-500: #2196f3; /* Main brand color */
  --primary-600: #1e88e5;
  --primary-700: #1976d2;
  --primary-800: #1565c0;
  --primary-900: #0d47a1;
}
```

#### Cosmic Purple
```css
:root {
  --primary-50: #f3e5f5;
  --primary-100: #e1bee7;
  --primary-200: #ce93d8;
  --primary-300: #ba68c8;
  --primary-400: #ab47bc;
  --primary-500: #9c27b0; /* Main brand color */
  --primary-600: #8e24aa;
  --primary-700: #7b1fa2;
  --primary-800: #6a1b9a;
  --primary-900: #4a148c;
}
```

#### Emerald Forest
```css
:root {
  --primary-50: #e8f5e8;
  --primary-100: #c8e6c9;
  --primary-200: #a5d6a7;
  --primary-300: #81c784;
  --primary-400: #66bb6a;
  --primary-500: #4caf50; /* Main brand color */
  --primary-600: #43a047;
  --primary-700: #388e3c;
  --primary-800: #2e7d32;
  --primary-900: #1b5e20;
}
```

#### Sunset Orange
```css
:root {
  --primary-50: #fff3e0;
  --primary-100: #ffe0b2;
  --primary-200: #ffcc80;
  --primary-300: #ffb74d;
  --primary-400: #ffa726;
  --primary-500: #ff9800; /* Main brand color */
  --primary-600: #fb8c00;
  --primary-700: #f57c00;
  --primary-800: #ef6c00;
  --primary-900: #e65100;
}
```

### Neutral Colors

#### Dark Theme
```css
:root {
  --background: #0f0f23;
  --surface: #1a1a2e;
  --surface-variant: #252547;
  --on-background: #e2e2e9;
  --on-surface: #d0d0d7;
  --on-surface-variant: #b8b8c5;
  --outline: #4a4a6a;
  --outline-variant: #353553;
}
```

#### Light Theme
```css
:root {
  --background: #f8fafc;
  --surface: #ffffff;
  --surface-variant: #f1f5f9;
  --on-background: #0f172a;
  --on-surface: #1e293b;
  --on-surface-variant: #475569;
  --outline: #cbd5e1;
  --outline-variant: #e2e8f0;
}
```

## Typography System

### Font Families
- **Primary**: Inter (Sans-serif)
- **Monospace**: JetBrains Mono
- **Display**: Poppins (for headings)

### Scale
```css
:root {
  --display-large: 57px/64px Poppins;
  --display-medium: 45px/52px Poppins;
  --display-small: 36px/44px Poppins;
  
  --headline-large: 32px/40px Inter;
  --headline-medium: 28px/36px Inter;
  --headline-small: 24px/32px Inter;
  
  --title-large: 22px/28px Inter;
  --title-medium: 16px/24px Inter;
  --title-small: 14px/20px Inter;
  
  --body-large: 16px/24px Inter;
  --body-medium: 14px/20px Inter;
  --body-small: 12px/16px Inter;
  
  --label-large: 14px/20px Inter;
  --label-medium: 12px/16px Inter;
  --label-small: 11px/16px Inter;
  
  --code-large: 16px/24px JetBrains Mono;
  --code-medium: 14px/20px JetBrains Mono;
  --code-small: 12px/16px JetBrains Mono;
}
```

## Component Specifications

### Buttons

#### Primary Button
- **Height**: 48px
- **Padding**: 24px horizontal
- **Border Radius**: 12px
- **Elevation**: 2px shadow
- **States**: Default, Hover, Focus, Pressed, Disabled

#### Secondary Button
- **Height**: 40px
- **Padding**: 20px horizontal
- **Border**: 1px solid outline color
- **Border Radius**: 8px

### Cards

#### Surface Card
- **Border Radius**: 16px
- **Elevation**: 4px shadow
- **Padding**: 24px
- **Background**: Surface color

#### Elevated Card
- **Border Radius**: 20px
- **Elevation**: 8px shadow
- **Padding**: 32px
- **Background**: Surface variant

### Input Fields

#### Text Field
- **Height**: 56px
- **Border Radius**: 12px
- **Padding**: 16px
- **Label**: Floating label pattern
- **States**: Default, Focus, Error, Disabled

## Layout System

### Grid System
- **Base Unit**: 8px
- **Breakpoints**: Mobile (320px), Tablet (768px), Desktop (1024px), Large (1440px)
- **Columns**: 4 (Mobile), 8 (Tablet), 12 (Desktop)
- **Gutters**: 16px (Mobile), 24px (Tablet/Desktop)

### Spacing Scale
```css
:root {
  --space-1: 4px;
  --space-2: 8px;
  --space-3: 12px;
  --space-4: 16px;
  --space-5: 24px;
  --space-6: 32px;
  --space-7: 48px;
  --space-8: 64px;
  --space-9: 96px;
  --space-10: 128px;
}
```

## Animation System

### Duration
- **Quick**: 150ms
- **Moderate**: 300ms
- **Slow**: 500ms

### Easing Curves
- **Standard**: cubic-bezier(0.4, 0.0, 0.2, 1)
- **Deceleration**: cubic-bezier(0.0, 0.0, 0.2, 1)
- **Acceleration**: cubic-bezier(0.4, 0.0, 1, 1)
- **Sharp**: cubic-bezier(0.4, 0.0, 0.6, 1)

### Transitions
- **Hover**: All 150ms standard
- **Focus**: All 150ms standard
- **Page**: 300ms deceleration
- **Modal**: 300ms deceleration

## Iconography

### Icon Sets
- **Material Symbols** (Primary)
- **Lucide Icons** (Secondary)
- **Custom HelixCode Icons** (Brand-specific)

### Sizes
- **Small**: 16px
- **Medium**: 24px
- **Large**: 32px
- **XLarge**: 48px

## Platform-Specific Adaptations

### Mobile (iOS/Android/Aurora)
- **Navigation**: Bottom navigation bars
- **Gestures**: Swipe to navigate
- **Touch Targets**: Minimum 44px
- **Status Bar**: Platform-specific styling

### Desktop (Windows/macOS/Linux)
- **Navigation**: Sidebar navigation
- **Keyboard**: Extensive shortcut support
- **Window Management**: Multi-window support
- **Menu Bars**: Platform-specific menus

### Web
- **Responsive**: Fluid layouts
- **Progressive Enhancement**: Graceful degradation
- **Accessibility**: Full keyboard navigation
- **Performance**: Optimized loading

## Accessibility Standards

### Color Contrast
- **Normal Text**: 4.5:1 minimum
- **Large Text**: 3:1 minimum
- **UI Components**: 3:1 minimum

### Focus Indicators
- **Visible**: Clear focus rings
- **Consistent**: Same across all components
- **Non-obtrusive**: Doesn't break design

### Screen Reader Support
- **Semantic HTML**: Proper structure
- **ARIA Labels**: Descriptive labels
- **Live Regions**: Dynamic content updates

## Theme Implementation

### Theme Switching
- **System Default**: Auto-detect OS theme
- **Manual Override**: User preference storage
- **Smooth Transition**: Cross-fade animation
- **Persistent**: Remember user choice

### Theme Variants
1. **Deep Ocean** (Default Dark)
2. **Cosmic Purple** (Dark)
3. **Emerald Forest** (Light)
4. **Sunset Orange** (Light)
5. **Midnight Blue** (Dark)
6. **Solar Flare** (Light)

This design system ensures HelixCode provides the most stunning, professional, and user-friendly experience across all platforms while maintaining consistency and accessibility.