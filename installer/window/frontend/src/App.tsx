import React from "react";
import { TitleBar } from "@/components/TitleBar";
import { useAppConfig } from "@/hooks/useAppConfig";
import { useInstaller } from "@/hooks/useInstaller";
import { useSysInfo } from "@/hooks/useSysInfo";
import { Sidebar } from "@/components/Sidebar";
import { StepWelcome } from "@/components/StepWelcome";
import { StepTos } from "@/components/StepTos";
// import { StepConfig } from "@/components/StepConfig";
// import { StepInstalling } from "@/components/StepInstalling";
// import { StepComplete } from "@/components/StepComplete";
import { Navigation } from "@/components/Navigation";
import { StepConfig } from "./components/StepConfig";

export function App() {
  const appConfig = useAppConfig();
  const installer = useInstaller();
  const sysInfo = useSysInfo();

  const containerStyle: React.CSSProperties = {
    width: appConfig?.Width ? `${appConfig.Width}px` : "900px",
    height: appConfig?.Height ? `${appConfig.Height}px` : "600px",
  };

  // Evaluate prerequisites and current step input requirements
  const isX64 = sysInfo.targetArch === "x64" || sysInfo.targetArch === "amd64";
  // const isWelcomePassed = sysInfo.isAdmin && isX64;

  const isNextEnabled = (() => {
    if (installer.step === "welcome") return true;
    if (installer.step === "tos") return installer.acceptedTos;
    if (installer.step === "config") return Boolean(installer.apiKey.trim());
    return true;
  })();

  return (
    <div className="flex h-screen w-screen items-center justify-center bg-stone-950 font-sans text-slate-100 select-none">
      <div
        style={containerStyle}
        className="relative flex overflow-hidden rounded-xl border border-slate-800 bg-stone-900 shadow-2xl shadow-cyan-950/20"
      >
        {/* Sidebar touches top-to-bottom */}
        <Sidebar step={installer.step} config={appConfig} />

        {/* Main Content Area with TitleBar and Step Views */}
        <div className="flex flex-1 flex-col overflow-hidden">
          <TitleBar />

          <main className="flex flex-1 flex-col justify-between overflow-y-auto bg-stone-900/60 p-8">
            <div className="flex-1">
              {installer.step === "welcome" && <StepWelcome sysInfo={sysInfo} />}

              {installer.step === "tos" && (
                <StepTos
                  accepted={installer.acceptedTos}
                  setAccepted={installer.setAcceptedTos}
                />
              )}

               {installer.step === "config" && (
                <StepConfig
                  isPrivate={installer.isPrivate}
                  setIsPrivate={installer.setIsPrivate}
                  fqdn={installer.fqdn}
                  setFqdn={installer.setFqdn}
                  port={installer.port}
                  setPort={installer.setPort}
                  apiKey={installer.apiKey}
                  setApiKey={installer.setApiKey}
                  licenseFile={installer.licenseFile}
                  setLicenseFile={installer.setLicenseFile}
                  installPath={installer.installPath}
                  setInstallPath={installer.setInstallPath}
                  installAsService={installer.installAsService}
                  setInstallAsService={installer.setInstallAsService}
                  enableAutostart={installer.enableAutostart}
                  setEnableAutostart={installer.setEnableAutostart}
                />
              )}
              {/*
              {installer.step === "installing" && (
                <StepInstalling
                  progress={installer.progress}
                  currentAction={installer.currentAction}
                  logs={installer.logs}
                  error={installer.error}
                />
              )}

              {installer.step === "complete" && <StepComplete />} */}
            </div>

            {/* Centralized Navigation Component */}
            <Navigation
              step={installer.step}
              setStep={installer.setStep}
              next={isNextEnabled}
              onResetTos={() => installer.setAcceptedTos(false)}
              onResetConfig={() => installer.setApiKey("")}
            />
          </main>
        </div>
      </div>
    </div>
  );
}

export default App;