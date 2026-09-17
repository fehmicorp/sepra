import React from "react";
import { Server, Key, Folder, Globe, FileCheck, Upload, Lock } from "lucide-react";

interface StepConfigProps {
  isPrivate: boolean;
  setIsPrivate: (val: boolean) => void;
  fqdn: string;
  setFqdn: (val: string) => void;
  port: string;
  setPort: (val: string) => void;
  apiKey: string;
  setApiKey: (val: string) => void;
  licenseFile: File | null;
  setLicenseFile: (file: File | null) => void;
  installPath: string;
  setInstallPath: (val: string) => void;
  installAsService: boolean;
  setInstallAsService: (val: boolean) => void;
  enableAutostart: boolean;
  setEnableAutostart: (val: boolean) => void;
}

export function StepConfig({
  isPrivate,
  setIsPrivate,
  fqdn,
  setFqdn,
  port,
  setPort,
  apiKey,
  setApiKey,
  licenseFile,
  setLicenseFile,
  installPath,
  setInstallPath,
  installAsService,
  setInstallAsService,
  enableAutostart,
  setEnableAutostart,
}: StepConfigProps) {
  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files && e.target.files[0]) {
      setLicenseFile(e.target.files[0]);
    }
  };

  return (
    <div className="space-y-4">
      <div>
        <h1 className="text-xl font-bold text-white">Agent Configuration</h1>
        <p className="text-xs text-stone-400">Configure connection mode and target endpoint settings.</p>
      </div>

      {/* Mode Selector Toggle */}
      <div className="grid grid-cols-2 gap-3 rounded-xl border border-stone-800 bg-stone-950 p-1.5">
        <button
          type="button"
          onClick={() => setIsPrivate(true)}
          className={`flex cursor-pointer items-center justify-center gap-2 rounded-lg py-2 text-xs font-medium transition-all ${
            isPrivate
              ? "bg-cyan-600 text-white shadow-md shadow-cyan-950/50"
              : "text-stone-400 hover:text-stone-200"
          }`}
        >
          <Server className="h-4 w-4" />
          Private Server
        </button>
        <button
          type="button"
          onClick={() => setIsPrivate(false)}
          className={`flex cursor-pointer items-center justify-center gap-2 rounded-lg py-2 text-xs font-medium transition-all ${
            !isPrivate
              ? "bg-cyan-600 text-white shadow-md shadow-cyan-950/50"
              : "text-stone-400 hover:text-stone-200"
          }`}
        >
          <Globe className="h-4 w-4" />
          Public Server
        </button>
      </div>

      {/* Dynamic Connection Options */}
      {isPrivate ? (
        <div className="space-y-3">
          <div className="grid grid-cols-3 gap-2">
            <div className="col-span-2">
              <label className="mb-1 block text-xs font-medium text-stone-300">Server FQDN / IP</label>
              <div className="flex items-center rounded-lg border border-stone-700 bg-stone-950 px-3 py-2 focus-within:border-cyan-500">
                <Server className="mr-2 h-4 w-4 text-stone-400" />
                <input
                  type="text"
                  placeholder="agent.sepra.cloud"
                  value={fqdn}
                  onChange={(e) => setFqdn(e.target.value)}
                  className="w-full bg-transparent text-xs text-white outline-none"
                />
              </div>
            </div>
            <div>
              <label className="mb-1 block text-xs font-medium text-stone-300">Port</label>
              <div className="flex items-center rounded-lg border border-stone-700 bg-stone-950 px-3 py-2 focus-within:border-cyan-500">
                <input
                  type="text"
                  placeholder="8443"
                  value={port}
                  onChange={(e) => setPort(e.target.value)}
                  className="w-full bg-transparent text-xs text-white outline-none"
                />
              </div>
            </div>
          </div>

          <div>
            <label className="mb-1 block text-xs font-medium text-stone-300">API Key</label>
            <div className="flex items-center rounded-lg border border-stone-700 bg-stone-950 px-3 py-2 focus-within:border-cyan-500">
              <Key className="mr-2 h-4 w-4 text-stone-400" />
              <input
                type="password"
                placeholder="fhm_live_xxxxxxxxxxxxxxxx"
                value={apiKey}
                onChange={(e) => setApiKey(e.target.value)}
                className="w-full bg-transparent text-xs text-white outline-none"
              />
            </div>
          </div>
        </div>
      ) : (
        <div className="space-y-2">
          <label className="block text-xs font-medium text-stone-300">Upload License Verification File</label>
          <label className="flex cursor-pointer flex-col items-center justify-center rounded-xl border border-dashed border-stone-700 bg-stone-950/60 p-4 transition-all hover:border-cyan-500 hover:bg-stone-900/50">
            {licenseFile ? (
              <div className="flex items-center gap-2 text-xs font-medium text-cyan-400">
                <FileCheck className="h-5 w-5" />
                <span>{licenseFile.name}</span>
              </div>
            ) : (
              <div className="flex flex-col items-center text-center">
                <Upload className="mb-1 h-6 w-6 text-stone-400" />
                <span className="text-xs font-medium text-stone-300">Click to upload license file (.lic, .json)</span>
                <span className="mt-0.5 text-[10px] text-stone-500">Used for verifying gateway authentication token</span>
              </div>
            )}
            <input type="file" accept=".lic,.json,.key" onChange={handleFileChange} className="hidden" />
          </label>
        </div>
      )}

      {/* System Configurations */}
      <div className="space-y-3 pt-1">
        <div>
          <label className="mb-1 block text-xs font-medium text-stone-300">Installation Path</label>
          <div className="flex items-center rounded-lg border border-stone-800 bg-stone-900/50 px-3 py-2 ">
            <Folder className="mr-2 h-4 w-4 text-stone-500" />
            <input
              type="text"
              value={installPath}
              readOnly
              disabled
              className="w-full bg-transparent text-xs text-stone-400 outline-none"
            />
            <Lock className="ml-2 h-3.5 w-3.5 text-stone-500" />
          </div>
        </div>

        <div className="space-y-2 pt-1">
          <label className="flex items-center gap-2 text-xs text-stone-300 ">
            <input
              type="checkbox"
              checked={true}
              className="rounded border-stone-700 bg-stone-800 text-cyan-600 focus:ring-0 opacity-80 "
            />
            <span>
              Install as background Windows Service (<code className="text-cyan-400">wind-sepra-v1</code>)
            </span>
            <span className="text-[10px] text-stone-500 font-mono">(Required)</span>
          </label>

          <label className="flex items-center gap-2 text-xs text-stone-300 ">
            <input
              type="checkbox"
              checked={true}
              className="rounded border-stone-700 bg-stone-800 text-cyan-600 focus:ring-0 opacity-80 "
            />
            <span>Enable automatic startup on Windows boot</span>
            <span className="text-[10px] text-stone-500 font-mono">(Required)</span>
          </label>
        </div>
      </div>
    </div>
  );
}