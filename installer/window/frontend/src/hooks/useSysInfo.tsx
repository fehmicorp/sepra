import { useEffect, useState } from "react";

export interface SysPayload {
  isAdmin: boolean;
  targetArch: string;
}

const DEFAULT_CONFIG: SysPayload = {
  isAdmin: false,
  targetArch: "null"
}

export function useSysInfo(): SysPayload {
  const [sysInfo, setSysInfo] = useState<SysPayload>(DEFAULT_CONFIG);
  useEffect(() => {
    async function fetchSysInfo() {
      try {
        const { GetSysInfo } = await import("../..//wailsjs/go/main/App");
        if (typeof GetSysInfo === "function") {
          const info = await GetSysInfo();
          if (info) {
            setSysInfo(info)
          }
        }
      } catch (err) {
        console.warn("Failed to fetch system info, using fallbacks:", err);
      }
    }
    fetchSysInfo();
  }, []);

  return sysInfo;
}