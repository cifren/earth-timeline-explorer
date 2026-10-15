import js from '@eslint/js'
import boundaries from 'eslint-plugin-boundaries'
import jsxA11y from 'eslint-plugin-jsx-a11y'
import reactHooks from 'eslint-plugin-react-hooks'
import { defineConfig, globalIgnores } from 'eslint/config'
import globals from 'globals'
import tseslint from 'typescript-eslint'

const sameContext = { context: '{{ from.captured.context }}' }

export default defineConfig([
  globalIgnores(['dist', 'var', '.cache']),
  {
    files: ['**/*.{ts,tsx}'],
    extends: [
      js.configs.recommended,
      tseslint.configs.recommendedTypeChecked,
      reactHooks.configs.flat.recommended,
      jsxA11y.flatConfigs.recommended,
    ],
    languageOptions: {
      globals: globals.browser,
      parserOptions: { projectService: true, tsconfigRootDir: import.meta.dirname },
    },
  },
  {
    files: ['src/**/*.{ts,tsx}'],
    plugins: { boundaries },
    settings: {
      'import/resolver': { typescript: { project: './tsconfig.app.json' } },
      'boundaries/elements': [
        { type: 'entry', pattern: 'src/*.{ts,tsx}', partialMatch: false },
        { type: 'app', pattern: 'src/app' },
        { type: 'layer', pattern: 'src/bounded-contexts/*/*', capture: ['context', 'layer'] },
        // Files at the root of a context: its public index.ts.
        { type: 'context', pattern: 'src/bounded-contexts/*', capture: ['context'] },
        { type: 'core', pattern: 'src/core/*' },
        { type: 'shared', pattern: 'src/shared/*' },
        { type: 'generated', pattern: 'var/*' },
      ],
    },
    rules: {
      'boundaries/dependencies': [
        'error',
        {
          default: 'disallow',
          policies: [
            { from: { element: { type: 'entry' } }, allow: { to: { element: { type: 'app' } } } },
            {
              from: { element: { type: 'app' } },
              allow: { to: { element: { types: { anyOf: ['app', 'context', 'core', 'shared'] } } } },
            },
            {
              from: { element: { type: 'context' } },
              allow: { to: { element: { type: 'layer', captured: sameContext } } },
            },
            // Layers point inward: presentation → application → infrastructure → domain.
            {
              from: { element: { type: 'layer', captured: { layer: 'presentation' } } },
              allow: {
                to: {
                  element: { type: 'layer', captured: { ...sameContext, layer: '{presentation,application,domain}' } },
                },
              },
            },
            {
              from: { element: { type: 'layer', captured: { layer: 'application' } } },
              allow: {
                to: {
                  element: {
                    type: 'layer',
                    captured: { ...sameContext, layer: '{application,infrastructure,domain}' },
                  },
                },
              },
            },
            {
              from: { element: { type: 'layer', captured: { layer: 'infrastructure' } } },
              allow: {
                to: {
                  element: { type: 'layer', captured: { ...sameContext, layer: '{infrastructure,domain}' } },
                },
              },
            },
            {
              from: { element: { type: 'layer', captured: { layer: 'domain' } } },
              allow: { to: { element: { type: 'layer', captured: { ...sameContext, layer: 'domain' } } } },
            },
            {
              from: { element: { type: 'layer', captured: { layer: 'domain' } } },
              allow: { to: { element: { type: 'core' } } },
            },
            {
              from: { element: { type: 'layer', captured: { layer: '!domain' } } },
              allow: { to: { element: { types: { anyOf: ['core', 'shared'] } } } },
            },
            {
              from: { element: { type: 'layer', captured: { layer: 'infrastructure' } } },
              allow: { to: { element: { type: 'generated' } } },
            },
            {
              from: { element: { type: 'core' } },
              allow: { to: { element: { types: { anyOf: ['core', 'shared'] } } } },
            },
            { from: { element: { type: 'shared' } }, allow: { to: { element: { type: 'shared' } } } },
          ],
        },
      ],
    },
  },
  {
    files: ['**/*.test.{ts,tsx}', 'src/test-setup.ts', 'src/test/**'],
    rules: { 'boundaries/dependencies': 'off' },
  },
])
