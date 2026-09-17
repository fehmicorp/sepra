import { useState, useEffect } from "react";

export interface AppConfig {
  AppName: string;
  Description: string;
  Icon: string;
  Tagline: string;
  Title: string;
  Width: number;
  Height: number;
	Min: boolean;
	Max: boolean;
	Quit : boolean;
  Version: string;
  Domain: string;
  installDir: string;
}

const DEFAULT_CONFIG: AppConfig = {
  AppName: "",
  Description: "",
  Icon: "",
  Tagline: "",
  Title: "",
  Width: 820,
  Height: 560,
  Min: false,
  Max: false,
  Quit: true,
  Version: "v1.0.0",
  Domain: "fehmicorp.in",
  installDir: "%ProgramFiles%\\fehmi\\agent",
};

export function useAppConfig(): AppConfig {
  const [config, setConfig] = useState<AppConfig>(DEFAULT_CONFIG);

  useEffect(() => {
    async function fetchConfig() {
      try {
        // Step up out of src/ into wailsjs/
        const { GetConfig } = await import("../..//wailsjs/go/main/App");
        if (typeof GetConfig === "function") {
          const cfg = await GetConfig();
          if (cfg) {
            setConfig((prev) => ({
              ...prev,
              ...cfg,
              Width: cfg.Width ?? prev.Width,
              Height: cfg.Height ?? prev.Height,
            }));
          }
        }
      } catch (err: unknown) {
        console.warn("Wails runtime not detected, using fallback config:", err);
      }
    }

    fetchConfig();
  }, []);

  return config;
}