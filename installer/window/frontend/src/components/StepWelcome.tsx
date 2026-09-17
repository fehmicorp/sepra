import React from "react";
import { SysPayload } from "@/hooks/useSysInfo";
import { PrerequisitesCheck } from "./PrerequisitesCheck";
import { RelaunchAsAdmin } from "@wailsjs/go/main/App";

interface StepWelcomeProps {
  sysInfo: SysPayload;
}

export function StepWelcome({ sysInfo }: StepWelcomeProps) {
  const handleGrantAdmin = async () => {
    try {
      await RelaunchAsAdmin();
    } catch (err) {
      console.error("Failed to relaunch as admin:", err);
    }
  };

  return (
    <div className="space-y-5">
      <div>
        <h1 className="text-2xl font-bold text-white">Install Fehmi Agent</h1>
        <p className="mt-1 text-sm text-stone-400">
          Deploy system monitoring, telemetry, and infrastructure automation services on this endpoint.
        </p>
      </div>
      <PrerequisitesCheck sysInfo={sysInfo} onGrantAdmin={handleGrantAdmin} />
    </div>
  );
}