import { useState, useEffect } from "react";
import { InstallStep } from "../components/types";
import { StartInstallation } from "@wailsjs/go/main/App";
import { EventsOn, EventsOff } from "@wailsjs/runtime/runtime";

export function useInstaller() {
  const [step, setStep] = useState<InstallStep>("welcome");
  const [acceptedTos, setAcceptedTos] = useState(false);
  
  // Configuration State
  const [isPrivate, setIsPrivate] = useState(true);
  const [fqdn, setFqdn] = useState("");
  const [port, setPort] = useState("8443");
  const [apiKey, setApiKey] = useState("");
  const [licenseFile, setLicenseFile] = useState<File | null>(null);
  const [installPath, setInstallPath] = useState("C:\\Program Files\\Sepra\\agent");
  const [installAsService, setInstallAsService] = useState(true);
  const [enableAutostart, setEnableAutostart] = useState(true);

  // Execution & Progress State
  const [progress, setProgress] = useState(0);
  const [currentAction, setCurrentAction] = useState("Initializing installation setup...");
  const [logs, setLogs] = useState<string[]>([]);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (step === "installing") {
      setLogs([]);
      setProgress(0);
      setError(null);

      // Listen for progress updates emitted from Go backend
      EventsOn("install_progress", (data: { percentage: number; message: string; log: string }) => {
        setProgress(data.percentage);
        setCurrentAction(data.message);
        setLogs((prev) => [...prev, data.log]);
      });

      // Listen for completion signal
      EventsOn("install_complete", () => {
        setTimeout(() => setStep("complete"), 800);
      });

      // Listen for errors
      EventsOn("install_error", (errMsg: string) => {
        setError(errMsg);
      });

      // Read license file content if present in Public mode
      const triggerInstallation = async () => {
        let licenseContent = "";
        if (!isPrivate && licenseFile) {
          licenseContent = await licenseFile.text();
        }

        return StartInstallation({
          isPrivate,
          fqdn,
          port,
          apiKey,
          licenseData: licenseContent,
          installPath,
          installAsService,
          enableAutostart,
        });
      };

      triggerInstallation().catch((err) => {
        setError(String(err));
      });

      return () => {
        EventsOff("install_progress");
        EventsOff("install_complete");
        EventsOff("install_error");
      };
    }
  }, [
    step,
    isPrivate,
    fqdn,
    port,
    apiKey,
    licenseFile,
    installPath,
    installAsService,
    enableAutostart,
  ]);

  return {
    step,
    setStep,
    acceptedTos,
    setAcceptedTos,
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
    progress,
    currentAction,
    logs,
    error,
  };
}