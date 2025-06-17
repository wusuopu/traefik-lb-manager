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

interface Service {
  ID?: number
  WorkspaceID?: number
  Name: string
  LBServers?: Array<{
    Url: string
    HostName?: string   // 根据 Url 解析而来
    Port?: string
    PathName?: string
    PreservePath?: boolean
    Weight?: number
  }>
  CreatedAt?: string
  UpdatedAt?: string
}
interface ExternalService {
  HostName: string
  Name?: string
  Stack?: string
}


interface Certificate {
  ID?: number
  WorkspaceID?: number
  Name: string
  Domain: string
  Stauts?: string
  Enable: boolean
  ExpiredAt?: string
  CreatedAt?: string
  UpdatedAt?: string
}
