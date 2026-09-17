import { SysPayload } from "@/hooks/useSysInfo";
import { CheckCircle2, XCircle, ShieldAlert, ArrowRight } from "lucide-react";

interface PrerequisitesCheckProps {
  sysInfo?: SysPayload | null;
  onGrantAdmin?: () => void;
}

export function PrerequisitesCheck({ sysInfo = null, onGrantAdmin }: PrerequisitesCheckProps) {
  const isAdmin = sysInfo?.isAdmin ?? false;
  const targetArch = sysInfo?.targetArch || "x64";
  const isX64 = targetArch === "x64" || targetArch === "amd64";

  return (
    <div className="rounded-lg border border-stone-800 bg-stone-950/60 p-4 space-y-4 select-none">
      <div className="flex items-center justify-between">
        <h3 className="text-xs font-semibold uppercase tracking-wider text-cyan-400">
          Prerequisites Check
        </h3>
        {!isAdmin && (
          <span className="flex items-center gap-1 text-[10px] font-medium text-amber-400 bg-amber-950/60 px-2 py-0.5 rounded border border-amber-800/50">
            <ShieldAlert className="h-3 w-3" /> Admin Required
          </span>
        )}
      </div>

      <div className="grid grid-cols-2 gap-3 text-xs">
        {/* Architecture Check */}
        <div className="flex items-center gap-2 text-stone-300">
          {isX64 ? (
            <CheckCircle2 className="h-4 w-4 shrink-0 text-emerald-400" />
          ) : (
            <XCircle className="h-4 w-4 shrink-0 text-rose-400" />
          )}
          <span>Architecture: {targetArch}</span>
        </div>

        {/* Administrator Check */}
        <div className="flex items-center gap-2 text-stone-300">
          {isAdmin ? (
            <CheckCircle2 className="h-4 w-4 shrink-0 text-emerald-400" />
          ) : (
            <XCircle className="h-4 w-4 shrink-0 text-rose-400" />
          )}
          <span>Admin Privilege</span>
        </div>

        {/* PowerShell Check */}
        <div className="flex items-center gap-2 text-stone-300">
          <CheckCircle2 className="h-4 w-4 shrink-0 text-emerald-400" />
          <span>PowerShell 5.1+</span>
        </div>

        {/* Port Check */}
        <div className="flex items-center gap-2 text-stone-300">
          <CheckCircle2 className="h-4 w-4 shrink-0 text-emerald-400" />
          <span>Port 443 Outbound</span>
        </div>
      </div>

      {/* Grant Admin Action Banner */}
      {!isAdmin && (
        <div className="mt-3 flex items-center justify-between rounded-md border border-amber-500/30 bg-amber-500/10 p-2.5 text-xs text-amber-200">
          <div className="flex items-center gap-2">
            <ShieldAlert className="h-4 w-4 text-amber-400 shrink-0" />
            <span>Installer requires administrative privileges to proceed.</span>
          </div>
          <button
            onClick={onGrantAdmin}
            type="button"
            className="flex items-center gap-1.5 rounded bg-amber-500 px-3 py-1.5 font-semibold text-stone-950 hover:bg-amber-400 transition-colors cursor-pointer shrink-0"
          >
            Grant Privilege
            <ArrowRight className="h-3.5 w-3.5" />
          </button>
        </div>
      )}
    </div>
  );
}