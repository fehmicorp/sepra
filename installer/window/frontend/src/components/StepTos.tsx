import React, { useState } from "react";
import { FileText, ExternalLink, CheckCircle } from "lucide-react";

interface StepTosProps {
  accepted: boolean;
  setAccepted: (value: boolean) => void;
}

export function StepTos({ accepted, setAccepted }: StepTosProps) {
  const [hasScrolledToBottom, setHasScrolledToBottom] = useState(false);

  const handleScroll = (e: React.UIEvent<HTMLDivElement>) => {
    const { scrollTop, scrollHeight, clientHeight } = e.currentTarget;
    if (scrollHeight - scrollTop <= clientHeight + 20) {
      setHasScrolledToBottom(true);
    }
  };

  return (
    <div className="space-y-4">
      <div>
        <h1 className="text-xl font-bold text-white">Terms of Service & EULA</h1>
        <p className="text-xs text-stone-400">
          Please review and accept the End User License Agreement before proceeding.
        </p>
      </div>

      <div
        onScroll={handleScroll}
        className="h-72 overflow-y-auto rounded-lg border border-stone-800 bg-stone-950 p-4 text-xs text-stone-300 space-y-3 leading-relaxed shadow-inner"
      >
        <div className="flex items-center gap-2 text-cyan-400 font-semibold border-b border-stone-800 pb-2">
          <FileText className="h-4 w-4" />
          <span>FEHMICORP END USER LICENSE AGREEMENT (EULA)</span>
        </div>

        <p className="text-stone-400">
          <strong>Last Updated:</strong> September 2026
        </p>

        <p>
          <strong>1. Grant of License:</strong> FehmiCorp grants you a non-exclusive, non-transferable license to deploy and execute the Fehmi Agent software on nodes and endpoints managed within your infrastructure.
        </p>

        <p>
          <strong>2. Telemetry & Data Collection:</strong> The agent collects operational metrics, system health stats, diagnostic telemetry, and performance parameters required for centralized infrastructure monitoring and automated orchestration via <code className="text-cyan-400">fehmicorp.in</code>.
        </p>

        <p>
          <strong>3. Service Elevation:</strong> Installing this software grants system-level privileges (Windows Service) necessary to execute automated maintenance scripts, security patches, and remote commands issued by your authenticated administrative user accounts.
        </p>

        <p>
          <strong>4. Security & Compliance:</strong> All telemetry and control commands are transmitted using TLS 1.3 encryption. You are responsible for safeguarding enrollment keys and API tokens generated from your organization dashboard.
        </p>

        <p>
          <strong>5. Limitation of Liability:</strong> In no event shall FehmiCorp be liable for indirect, incidental, or consequential damages arising out of system configurations or automated service executions performed by this agent.
        </p>
      </div>

      <div className="space-y-3 pt-1">
        <div className="flex items-center justify-between text-xs text-stone-400 border-b border-stone-800/80 pb-2">
          <span className="flex items-center gap-1.5">
            {hasScrolledToBottom ? (
              <span className="text-emerald-400 flex items-center gap-1">
                <CheckCircle className="h-3.5 w-3.5" /> Reviewed terms
              </span>
            ) : (
              "Scroll down to read complete document"
            )}
          </span>
          <div className="flex gap-3">
            <a
              href="https://fehmicorp.in/terms"
              target="_blank"
              rel="noreferrer"
              className="hover:text-cyan-400 transition-colors flex items-center gap-1"
            >
              Open in Browser <ExternalLink className="h-3 w-3" />
            </a>
          </div>
        </div>

        <label className="flex items-start gap-2.5 text-xs text-stone-200 cursor-pointer select-none">
          <input
            type="checkbox"
            checked={accepted}
            onChange={(e) => setAccepted(e.target.checked)}
            className="mt-0.5 rounded border-stone-700 bg-stone-950 text-cyan-600 focus:ring-0 focus:ring-offset-0"
          />
          <span className="leading-tight">
            I have read, understood, and accept the Terms of Service and End User License Agreement.
          </span>
        </label>
      </div>
    </div>
  );
}