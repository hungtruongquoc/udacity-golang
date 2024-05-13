/// <reference types="vite/client" />

interface ImportMetaEnv {
    readonly VITE_API_BASE_URL: string;
    // You can define other environment variables here
}

interface ImportMeta {
    readonly env: ImportMetaEnv;
}