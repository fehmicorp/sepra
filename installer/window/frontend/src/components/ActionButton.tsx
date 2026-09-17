import { WindowMinimise, WindowToggleMaximise, Quit } from "@wailsjs/runtime/runtime";
import { Minus, Square, X } from "lucide-react";

interface ActionButtonsProps {
  allowMinimize?: boolean;
  allowMaximize?: boolean;
  allowQuit?: boolean;
}

export function ActionButtons({
  allowMinimize = true,
  allowMaximize = true,
  allowQuit = true,
}: ActionButtonsProps) {
  return (
    <div className="flex items-center gap-1.5">
      {/* Minimize Button */}
      <button
        disabled={!allowMinimize}
        onClick={() => allowMinimize && WindowMinimise?.()}
        title={allowMinimize ? "Minimize" : "Disabled"}
        className={`group flex h-5 w-5 items-center justify-center rounded-full border-0 outline-none transition-colors ${
          allowMinimize
            ? "bg-slate-700/60 hover:bg-amber-500 "
            : "bg-slate-800/40 opacity-40"
        }`}
      >
        <Minus className="h-3 w-3 text-slate-300 group-hover:text-white" />
      </button>

      {/* Maximize Button */}
      <button
        disabled={!allowMaximize}
        onClick={() => allowMaximize && WindowToggleMaximise?.()}
        title={allowMaximize ? "Maximize" : "Disabled"}
        className={`group flex h-5 w-5 items-center justify-center rounded-full border-0 outline-none transition-colors ${
          allowMaximize
            ? "bg-slate-700/60 hover:bg-emerald-500 "
            : "bg-slate-800/40 opacity-40"
        }`}
      >
        <Square className="h-2.5 w-2.5 text-slate-300 group-hover:text-white" />
      </button>

      {/* Close/Quit Button */}
      <button
        disabled={!allowQuit}
        onClick={() => allowQuit && Quit?.()}
        title={allowQuit ? "Close" : "Disabled"}
        className={`group flex h-5 w-5 items-center justify-center rounded-full border-0 outline-none transition-colors ${
          allowQuit
            ? "bg-slate-700/60 hover:bg-rose-500 "
            : "bg-slate-800/40 opacity-40"
        }`}
      >
        <X className="h-3 w-3 text-slate-300 group-hover:text-white" />
      </button>
    </div>
  );
}