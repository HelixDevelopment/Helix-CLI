// Figma Design Generator for HelixCode
// This script would generate Figma-compatible design files

const designSystem = {
  colors: {
    themes: {
      deepOcean: {
        primary: '#2196f3',
        background: '#0f0f23',
        surface: '#1a1a2e',
        onBackground: '#e2e2e9',
        onSurface: '#d0d0d7'
      },
      cosmicPurple: {
        primary: '#9c27b0',
        background: '#1a0f23',
        surface: '#251a2e',
        onBackground: '#e2dde9',
        onSurface: '#d0cbd7'
      },
      emeraldForest: {
        primary: '#4caf50',
        background: '#f8faf8',
        surface: '#ffffff',
        onBackground: '#0f2a0f',
        onSurface: '#1e3b1e'
      },
      sunsetOrange: {
        primary: '#ff9800',
        background: '#faf8f0',
        surface: '#fffaf0',
        onBackground: '#2a1f0f',
        onSurface: '#3b2e1e'
      }
    }
  },
  typography: {
    fontFamilies: {
      inter: 'Inter',
      poppins: 'Poppins',
      jetbrains: 'JetBrains Mono'
    },
    scales: {
      displayLarge: { size: 57, lineHeight: 64 },
      displayMedium: { size: 45, lineHeight: 52 },
      displaySmall: { size: 36, lineHeight: 44 },
      headlineLarge: { size: 32, lineHeight: 40 },
      headlineMedium: { size: 28, lineHeight: 36 },
      headlineSmall: { size: 24, lineHeight: 32 },
      titleLarge: { size: 22, lineHeight: 28 },
      titleMedium: { size: 16, lineHeight: 24 },
      titleSmall: { size: 14, lineHeight: 20 },
      bodyLarge: { size: 16, lineHeight: 24 },
      bodyMedium: { size: 14, lineHeight: 20 },
      bodySmall: { size: 12, lineHeight: 16 }
    }
  },
  spacing: {
    baseUnit: 8,
    scale: [4, 8, 12, 16, 24, 32, 48, 64, 96, 128]
  },
  borderRadius: {
    small: 8,
    medium: 12,
    large: 16,
    xlarge: 20
  }
};

// Application screen definitions
const applicationScreens = {
  terminal: {
    name: 'Terminal UI',
    screens: [
      {
        name: 'Dashboard',
        components: [
          'Header with logo and controls',
          'Navigation tabs',
          'Active sessions list',
          'Worker status panel',
          'Quick actions bar'
        ]
      },
      {
        name: 'Session View',
        components: [
          'Session header with controls',
          'Mode navigation',
          'Project structure tree',
          'AI chat interface',
          'Active tools status'
        ]
      }
    ]
  },
  desktop: {
    name: 'Desktop Application',
    screens: [
      {
        name: 'Main Workspace',
        components: [
          'Menu bar',
          'Project explorer sidebar',
          'Code editor center panel',
          'Worker status sidebar',
          'Tool output panel'
        ]
      },
      {
        name: 'Settings Panel',
        components: [
          'Theme selection grid',
          'Keybindings editor',
          'Worker configuration',
          'Tool settings'
        ]
      }
    ]
  },
  mobile: {
    name: 'Mobile Application',
    screens: [
      {
        name: 'Mobile Dashboard',
        components: [
          'App header with user menu',
          'Quick action buttons',
          'Active sessions cards',
          'Bottom navigation'
        ]
      },
      {
        name: 'Project View',
        components: [
          'File browser',
          'Code viewer',
          'AI chat interface',
          'Build status indicators'
        ]
      }
    ]
  },
  web: {
    name: 'Web Application',
    screens: [
      {
        name: 'Web Workspace',
        components: [
          'Application header',
          'Main navigation',
          'Project sidebar',
          'Editor area',
          'Right panel',
          'Status footer'
        ]
      }
    ]
  }
};

// Generate design specifications
function generateDesignSpecs() {
  const specs = {
    designSystem,
    applications: applicationScreens,
    components: {
      buttons: {
        primary: {
          height: 48,
          padding: [0, 24],
          borderRadius: 12,
          elevation: 2
        },
        secondary: {
          height: 40,
          padding: [0, 20],
          border: 1,
          borderRadius: 8
        }
      },
      cards: {
        surface: {
          borderRadius: 16,
          elevation: 4,
          padding: 24
        },
        elevated: {
          borderRadius: 20,
          elevation: 8,
          padding: 32
        }
      },
      inputs: {
        textField: {
          height: 56,
          borderRadius: 12,
          padding: 16
        }
      }
    }
  };

  return specs;
}

// Export for Figma integration
module.exports = {
  designSystem,
  applicationScreens,
  generateDesignSpecs
};

console.log('HelixCode Design Specifications Generated');
console.log('=========================================');
console.log('Themes:', Object.keys(designSystem.colors.themes));
console.log('Applications:', Object.keys(applicationScreens));

// This would be used with Figma API to create actual design files
// In a real implementation, this would connect to Figma API and create frames, components, etc.