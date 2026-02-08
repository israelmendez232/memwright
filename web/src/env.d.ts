/// <reference types="vite/client" />

interface ImportMetaEnv {
  readonly VITE_API_URL: string
  readonly VITE_AUTH_BYPASS?: string
  readonly VITE_REGISTRATION_ENABLED?: string
}

interface ImportMeta {
  readonly env: ImportMetaEnv
}
