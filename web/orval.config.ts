import { defineConfig } from 'orval'

export default defineConfig({
  api: {
    input: { target: '../back/api/openapi.yaml' },
    output: {
      mode: 'tags-split',
      target: 'var/api-client/endpoints.ts',
      schemas: 'var/api-client/model',
      client: 'fetch',
      baseUrl: { getBaseUrlFromSpecification: true },
      clean: true,
      override: {
        mutator: { path: 'src/shared/http/fetch-api.ts', name: 'fetchApi' },
      },
    },
  },
})
