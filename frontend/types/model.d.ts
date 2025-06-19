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

interface Rule {
  ID?: number
  Options: any
  Enable: boolean
  WorkspaceID?: number
  ServerID?: number
  CreatedAt?: string
  UpdatedAt?: string
}

interface Server {
  ID?: number
  WorkspaceID?: number
  Name: string
  Host: string[]
  Enable: boolean
}

interface Service {
  ID?: number
  WorkspaceID?: number
  Name: string
  LBServers?: Array<{
    url: string
    HostName?: string   // 根据 Url 解析而来
    Port?: string
    PathName?: string
    preservePath?: boolean
    weight?: number
  }>
  CreatedAt?: string
  UpdatedAt?: string
}
interface ExternalService {
  HostName: string
  Name?: string
  Stack?: string
  Label?: string
}

interface Middleware {
  ID?: number
  Name: string
  Category: string
  Options: any
  CreatedAt?: string
  UpdatedAt?: string
}

interface Certificate {
  ID?: number
  WorkspaceID?: number
  Name: string
  Domain: string
  Status?: string
  Enable: boolean
  ExpiredAt?: string
  CreatedAt?: string
  UpdatedAt?: string
}
