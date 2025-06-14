interface Workspace {
  ID?: number
  Name: string
  Description?: string
  ManagerBaseUrl: string
  Category?: string
  ApiBaseUrl?: string
  ApiKey?: string
  ApiSecret?: string
  CreatedAt?: string
  UpdatedAt?: string
  Entrypoints?: string[]
  TraefikConfig?: string
}
