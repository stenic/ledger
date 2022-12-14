export interface VersionData {
  id: string;
  application: {
    name: string;
  };
  location: {
    name: string;
  };
  environment: {
    name: string;
  };
  version: string;
  timestamp: string;
}
