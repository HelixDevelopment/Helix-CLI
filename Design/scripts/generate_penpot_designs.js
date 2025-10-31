// PenPot Design Generator for HelixCode
// This script would generate PenPot-compatible design files

const penpotDesignSystem = {
  metadata: {
    name: 'HelixCode Design System',
    version: '1.0.0',
    description: 'Stunning enterprise design system for HelixCode'
  },
  
  colors: {
    // Deep Ocean Theme
    deepOcean: {
      primary: {
        '50': '#e3f2fd',
        '100': '#bbdefb',
        '200': '#90caf9',
        '300': '#64b5f6',
        '400': '#42a5f5',
        '500': '#2196f3',
        '600': '#1e88e5',
        '700': '#1976d2',
        '800': '#1565c0',
        '900': '#0d47a1'
      },
      neutral: {
        background: '#0f0f23',
        surface: '#1a1a2e',
        surfaceVariant: '#252547',
        onBackground: '#e2e2e9',
        onSurface: '#d0d0d7',
        outline: '#4a4a6a'
      }
    },
    
    // Cosmic Purple Theme
    cosmicPurple: {
      primary: {
        '50': '#f3e5f5',
        '100': '#e1bee7',
        '200': '#ce93d8',
        '300': '#ba68c8',
        '400': '#ab47bc',
        '500': '#9c27b0',
        '600': '#8e24aa',
        '700': '#7b1fa2',
        '800': '#6a1b9a',
        '900': '#4a148c'
      },
      neutral: {
        background: '#1a0f23',
        surface: '#251a2e',
        surfaceVariant: '#302547',
        onBackground: '#e2dde9',
        onSurface: '#d0cbd7',
        outline: '#5a4a7a'
      }
    },
    
    // Emerald Forest Theme
    emeraldForest: {
      primary: {
        '50': '#e8f5e8',
        '100': '#c8e6c9',
        '200': '#a5d6a7',
        '300': '#81c784',
        '400': '#66bb6a',
        '500': '#4caf50',
        '600': '#43a047',
        '700': '#388e3c',
        '800': '#2e7d32',
        '900': '#1b5e20'
      },
      neutral: {
        background: '#f8faf8',
        surface: '#ffffff',
        surfaceVariant: '#f1f5f1',
        onBackground: '#0f2a0f',
        onSurface: '#1e3b1e',
        outline: '#cbd5cb'
      }
    },
    
    // Sunset Orange Theme
    sunsetOrange: {
      primary: {
        '50': '#fff3e0',
        '100': '#ffe0b2',
        '200': '#ffcc80',
        '300': '#ffb74d',
        '400': '#ffa726',
        '500': '#ff9800',
        '600': '#fb8c00',
        '700': '#f57c00',
        '800': '#ef6c00',
        '900': '#e65100'
      },
      neutral: {
        background: '#faf8f0',
        surface: '#fffaf0',
        surfaceVariant: '#f1ede0',
        onBackground: '#2a1f0f',
        onSurface: '#3b2e1e',
        outline: '#cbd5e1'
      }
    }
  },
  
  typography: {
    fontFamilies: {
      inter: {
        name: 'Inter',
        type: 'sans-serif',
        weights: [400, 500, 600, 700]
      },
      poppins: {
        name: 'Poppins',
        type: 'sans-serif',
        weights: [400, 500, 600, 700]
      },
      jetbrains: {
        name: 'JetBrains Mono',
        type: 'monospace',
        weights: [400, 500, 600, 700]
      }
    },
    
    textStyles: {
      displayLarge: {
        fontFamily: 'poppins',
        fontSize: 57,
        lineHeight: 64,
        fontWeight: 400
      },
      displayMedium: {
        fontFamily: 'poppins',
        fontSize: 45,
        lineHeight: 52,
        fontWeight: 400
      },
      headlineLarge: {
        fontFamily: 'inter',
        fontSize: 32,
        lineHeight: 40,
        fontWeight: 600
      },
      titleLarge: {
        fontFamily: 'inter',
        fontSize: 22,
        lineHeight: 28,
        fontWeight: 500
      },
      bodyLarge: {
        fontFamily: 'inter',
        fontSize: 16,
        lineHeight: 24,
        fontWeight: 400
      },
      codeLarge: {
        fontFamily: 'jetbrains',
        fontSize: 16,
        lineHeight: 24,
        fontWeight: 400
      }
    }
  },
  
  spacing: {
    baseUnit: 8,
    scale: {
      '1': 4,
      '2': 8,
      '3': 12,
      '4': 16,
      '5': 24,
      '6': 32,
      '7': 48,
      '8': 64,
      '9': 96,
      '10': 128
    }
  },
  
  borderRadius: {
    none: 0,
    small: 8,
    medium: 12,
    large: 16,
    xlarge: 20,
    full: 9999
  },
  
  shadows: {
    '1': '0 1px 3px 0 rgba(0, 0, 0, 0.1), 0 1px 2px 0 rgba(0, 0, 0, 0.06)',
    '2': '0 4px 6px -1px rgba(0, 0, 0, 0.1), 0 2px 4px -1px rgba(0, 0, 0, 0.06)',
    '3': '0 10px 15px -3px rgba(0, 0, 0, 0.1), 0 4px 6px -2px rgba(0, 0, 0, 0.05)',
    '4': '0 20px 25px -5px rgba(0, 0, 0, 0.1), 0 10px 10px -5px rgba(0, 0, 0, 0.04)'
  }
};

// Component definitions for PenPot
const penpotComponents = {
  buttons: {
    primary: {
      name: 'Button / Primary',
      width: 'auto',
      height: 48,
      padding: [0, 24],
      borderRadius: 12,
      background: '{colors.primary.500}',
      textColor: '{colors.neutral.onPrimary}',
      typography: '{typography.titleMedium}',
      shadow: '{shadows.2}'
    },
    secondary: {
      name: 'Button / Secondary',
      width: 'auto',
      height: 40,
      padding: [0, 20],
      borderRadius: 8,
      border: '1px solid {colors.outline}',
      background: 'transparent',
      textColor: '{colors.primary.500}',
      typography: '{typography.titleMedium}'
    }
  },
  
  cards: {
    surface: {
      name: 'Card / Surface',
      width: 320,
      minHeight: 200,
      padding: 24,
      borderRadius: 16,
      background: '{colors.surface}',
      shadow: '{shadows.2}'
    },
    elevated: {
      name: 'Card / Elevated',
      width: 360,
      minHeight: 240,
      padding: 32,
      borderRadius: 20,
      background: '{colors.surfaceVariant}',
      shadow: '{shadows.3}'
    }
  },
  
  inputs: {
    textField: {
      name: 'Input / Text Field',
      width: 280,
      height: 56,
      padding: 16,
      borderRadius: 12,
      background: '{colors.surface}',
      border: '1px solid {colors.outline}',
      typography: '{typography.bodyLarge}'
    }
  },
  
  navigation: {
    sidebar: {
      name: 'Navigation / Sidebar',
      width: 280,
      minHeight: 600,
      background: '{colors.surface}',
      borderRight: '1px solid {colors.outline}'
    },
    bottomBar: {
      name: 'Navigation / Bottom Bar',
      width: '100%',
      height: 80,
      background: '{colors.surface}',
      borderTop: '1px solid {colors.outline}'
    }
  }
};

// Application screen templates
const penpotScreens = {
  terminal: {
    name: 'Terminal UI',
    artboards: [
      {
        name: 'Dashboard',
        width: 800,
        height: 600,
        background: '{colors.background}',
        components: [
          'header',
          'navigation',
          'sessionList',
          'workerStatus',
          'quickActions'
        ]
      },
      {
        name: 'Session View',
        width: 800,
        height: 600,
        background: '{colors.background}',
        components: [
          'sessionHeader',
          'modeNavigation',
          'projectTree',
          'aiChat',
          'toolsStatus'
        ]
      }
    ]
  },
  
  desktop: {
    name: 'Desktop Application',
    artboards: [
      {
        name: 'Main Workspace',
        width: 1440,
        height: 900,
        background: '{colors.background}',
        layout: 'three-column',
        components: [
          'menuBar',
          'projectExplorer',
          'codeEditor',
          'workerPanel',
          'outputPanel'
        ]
      },
      {
        name: 'Settings',
        width: 1200,
        height: 800,
        background: '{colors.background}',
        components: [
          'themeGrid',
          'keybindingsEditor',
          'workerConfig',
          'toolSettings'
        ]
      }
    ]
  },
  
  mobile: {
    name: 'Mobile Application',
    artboards: [
      {
        name: 'Dashboard',
        width: 375,
        height: 812,
        background: '{colors.background}',
        components: [
          'appHeader',
          'quickActions',
          'sessionCards',
          'bottomNav'
        ]
      },
      {
        name: 'Project View',
        width: 375,
        height: 812,
        background: '{colors.background}',
        components: [
          'fileBrowser',
          'codeViewer',
          'aiChat',
          'buildStatus'
        ]
      }
    ]
  },
  
  web: {
    name: 'Web Application',
    artboards: [
      {
        name: 'Web Workspace',
        width: 1440,
        height: 900,
        background: '{colors.background}',
        layout: 'responsive',
        components: [
          'appHeader',
          'mainNav',
          'projectSidebar',
          'editorArea',
          'rightPanel',
          'statusFooter'
        ]
      }
    ]
  }
};

// Generate PenPot design file structure
function generatePenPotDesign() {
  return {
    metadata: penpotDesignSystem.metadata,
    designSystem: {
      colors: penpotDesignSystem.colors,
      typography: penpotDesignSystem.typography,
      spacing: penpotDesignSystem.spacing,
      borderRadius: penpotDesignSystem.borderRadius,
      shadows: penpotDesignSystem.shadows
    },
    components: penpotComponents,
    screens: penpotScreens,
    themes: Object.keys(penpotDesignSystem.colors)
  };
}

// Export for PenPot integration
module.exports = {
  penpotDesignSystem,
  penpotComponents,
  penpotScreens,
  generatePenPotDesign
};

console.log('HelixCode PenPot Design Specifications Generated');
console.log('================================================');
console.log('Design System: Complete');
console.log('Components:', Object.keys(penpotComponents));
console.log('Applications:', Object.keys(penpotScreens));
console.log('Themes:', Object.keys(penpotDesignSystem.colors));

// This design specification can be used to generate actual PenPot files
// through PenPot's API or manual creation in the PenPot application