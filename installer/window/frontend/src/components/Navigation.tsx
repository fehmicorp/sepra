import React from "react";
import { ChevronRight, ChevronLeft, Check } from "lucide-react";
import { InstallStep } from "./types";
import { Quit } from "@wailsjs/runtime/runtime";

interface NavigationProps {
  step: InstallStep;
  setStep: (step: InstallStep) => void;
  prev?: boolean;
  next?: boolean;
  onFinish?: () => void;
  onResetTos?: () => void;
  onResetConfig?: () => void;
}

// Ordered list of wizard steps
const STEP_ORDER: InstallStep[] = [
  "welcome",
  "tos",
  "config",
  "installing",
  "complete",
];

export function Navigation({
  step,
  setStep,
  prev = true,
  next = true,
  onFinish,
  onResetTos,
  onResetConfig,
}: NavigationProps) {
  const currentIndex = STEP_ORDER.indexOf(step);
  const canGoBack = currentIndex > 0 && step !== "installing" && step !== "complete";
  const canGoNext = currentIndex >= 0 && currentIndex < STEP_ORDER.length - 2; // Stops before 'installing'

  const handleBack = () => {
    if (canGoBack) {
      if (step === "tos" && onResetTos) {
        onResetTos(); // e.g. setAcceptedTos(false)
      } else if (step === "config" && onResetConfig) {
        onResetConfig();
      }
      setStep(STEP_ORDER[currentIndex - 1]);
    }
  };

  const handleNext = () => {
    if (canGoNext) {
      setStep(STEP_ORDER[currentIndex + 1]);
    }
  };

  const handleFinish = () => {
    if (onFinish) {
      onFinish();
    } else if (typeof Quit === "function") {
      Quit();
    }
  };

  // Do not render navigation controls during active installation
  if (step === "installing") {
    return null;
  }

  return (
    <div className="mt-auto flex items-center justify-between border-t border-stone-800/80 px-2 pt-4">
      {/* Left Slot: Back Button or Empty Spacer to Maintain Space-Between Layout */}
      {canGoBack ? (
        <button
          type="button"
          onClick={handleBack}
          disabled={!prev}
          className="flex cursor-pointer items-center gap-1.5 rounded-lg border border-stone-700 px-6 py-2 text-xs font-medium text-slate-200 transition-all hover:bg-stone-700 disabled:cursor-not-allowed disabled:opacity-40"
        >
          Back
        </button>
      ) : (
        <div />
      )}

      {/* Right Slot: Proceed or Complete Action */}
      {canGoNext && (
        <button
          type="button"
          onClick={handleNext}
          disabled={!next}
          className="flex cursor-pointer items-center gap-1.5 rounded-lg bg-cyan-600 px-8 py-2 text-xs font-medium text-white shadow-lg shadow-cyan-950/50 transition-all hover:bg-cyan-500 disabled:cursor-not-allowed disabled:opacity-40"
        >
          {step === "welcome" ? "Next" : step === "config" ? "Install" : "Continue"}
        </button>
      )}

      {step === "complete" && (
        <button
          type="button"
          onClick={handleFinish}
          className="flex cursor-pointer items-center gap-1.5 rounded-lg bg-emerald-600 px-6 py-2 text-xs font-medium text-white shadow-lg shadow-emerald-950/50 transition-all hover:bg-emerald-500"
        >
          Finish
        </button>
      )}
    </div>
  );
}