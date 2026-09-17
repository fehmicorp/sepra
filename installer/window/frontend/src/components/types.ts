export type InstallStep = "welcome" | "tos" | "config" | "installing" | "complete";

export interface InstallerState {
  step: InstallStep;
  acceptedTos: boolean;
  installPath: string;
  serverUrl: string;
  apiKey: string;
  installAsService: boolean;
  enableAutostart: boolean;
  progress: number;
  currentAction: string;
  logs: string[];
}