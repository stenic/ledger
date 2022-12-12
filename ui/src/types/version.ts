export interface VersionData {
  id: string;
  application: {
    name: string;
  };
  environment: {
    name: string;
  };
  version: string;
  timestamp: string;
}
