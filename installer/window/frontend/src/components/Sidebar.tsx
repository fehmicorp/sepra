import React, { useState, useEffect } from "react";
import { Monitor, CheckCircle2, ShieldAlert, ShieldCheck } from "lucide-react";
import { InstallStep } from "./types";
import logo from "../assets/images/logo.png";
import { AppConfig } from "@/hooks/useAppConfig";

interface SidebarProps {
  step: InstallStep;
  config: AppConfig;
}

export interface SystemPayload {
  hostname: string;
  targetOs: string;
  targetArch: string;
  targetBuild: string;
  agentVersion: string;
}
export function Sidebar({ step, config }: SidebarProps) {
  const [sysInfo, setSysInfo] = useState<SystemPayload | null>(null);

  useEffect(() => {
    async function fetchSystemInfo() {
      try {
        const { GetSystemInfo } = await import("../../wailsjs/go/main/App");
        if (typeof GetSystemInfo === "function") {
          const info = await GetSystemInfo();
          setSysInfo(info);
        }
      } catch (err) {
        console.warn("Failed to fetch system info, using fallbacks:", err);
      }
    }
    fetchSystemInfo();
  }, []);

  return (
    <div className="flex w-64 flex-col justify-between border-r border-stone-800 bg-stone-950/60 p-6 select-none">
      <div>
        {/* Branding Header */}
        <div className="mb-8 flex items-center gap-3">
          <img src={logo} alt="Logo" className="h-8 w-8 object-contain" />
          <div>
            <h2 className="text-sm font-bold text-white leading-tight">
              {config.AppName}
            </h2>
            <p className="text-[11px] text-slate-400">{config.Tagline}</p>
          </div>
        </div>

        {/* Workflow Steps */}
        <div className="space-y-4">
          <StepItem
            active={step === "welcome"}
            completed={step === "tos" || step === "config" || step === "installing" || step === "complete"}
            title="1. Welcome"
            subtitle="Overview & EULA"
          />
          <StepItem
            active={step === "tos"}
            completed={step === "config" || step === "installing" || step === "complete"}
            title="2. Terms of Service"
            subtitle="End-user Agreement"
          />
          <StepItem
            active={step === "config"}
            completed={step === "installing" || step === "complete"}
            title="3. Configuration"
            subtitle="Server & API Setup"
          />
          <StepItem
            active={step === "installing"}
            completed={step === "complete"}
            title="4. Installation"
            subtitle="Deploying Binaries"
          />
          <StepItem
            active={step === "complete"}
            completed={step === "complete"}
            title="5. Finish"
            subtitle="Complete Setup"
          />
        </div>
      </div>

      {/* Dynamic System Details Card */}
      <div className="rounded-lg border border-stone-800/80 bg-stone-900/50 p-3 text-[11px] text-slate-400 space-y-1.5">
        <div className="flex items-center justify-between border-b border-stone-800/80 pb-1.5 mb-2 font-medium text-slate-300">
          <div className="flex items-center gap-1.5">
            <Monitor className="h-3.5 w-3.5 text-cyan-400" />
            <span>System Details</span>
          </div>
        </div>

        <div className="flex justify-between items-center">
          <span className="text-slate-500">Hostname:</span>
          <span className="font-mono text-slate-200">{sysInfo?.hostname || "Detecting..."}</span>
        </div>

        <div className="flex justify-between items-center">
          <span className="text-slate-500">OS Name:</span>
          <span className="font-mono text-slate-200">
            {sysInfo ? `${sysInfo.targetOs}` : "Detecting..."}
          </span>
        </div>

        <div className="flex justify-between items-center">
          <span className="text-slate-500">OS Build:</span>
          <span className="font-mono text-slate-200">
            {sysInfo ? `${sysInfo.targetBuild}` : "Detecting..."}
          </span>
        </div>

        <div className="flex justify-between items-center">
          <span className="text-slate-500">OS Arch:</span>
          <span className="font-mono text-slate-200">
            {sysInfo ? `${sysInfo.targetArch}` : "Detecting..."}
          </span>
        </div>

        <div className="flex justify-between items-center">
          <span className="text-slate-500">Agent Version:</span>
          <span className="font-mono text-slate-200">{sysInfo?.agentVersion || config.Version}</span>
        </div>
      </div>
    </div>
  );
}

function StepItem({
  title,
  subtitle,
  active,
  completed,
}: {
  title: string;
  subtitle: string;
  active: boolean;
  completed: boolean;
}) {
  return (
    <div className="flex items-start gap-3">
      <div className="mt-0.5">
        {completed && !active ? (
          <CheckCircle2 className="h-4 w-4 text-emerald-400" />
        ) : active ? (
          <div className="h-4 w-4 rounded-full border-2 border-cyan-400 bg-cyan-950 flex items-center justify-center">
            <div className="h-1.5 w-1.5 rounded-full bg-cyan-400" />
          </div>
        ) : (
          <div className="h-4 w-4 rounded-full border border-slate-700 bg-slate-900" />
        )}
      </div>
      <div>
        <p
          className={`text-xs font-semibold leading-none ${
            active ? "text-cyan-400" : completed ? "text-slate-200" : "text-slate-500"
          }`}
        >
          {title}
        </p>
        <p className="text-[10px] text-slate-500 mt-0.5">{subtitle}</p>
      </div>
    </div>
  );
}