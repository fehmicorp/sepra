import React from "react";
import { ActionButtons } from "./ActionButton";

export function TitleBar() {
  return (
    <div
      style={{ "--wails-draggable": "drag" } as React.CSSProperties}
      className="flex h-10 w-full items-center justify-between  px-4 select-none"
    >
      <div className="flex items-center gap-2"></div>

      <div
        style={{ "--wails-draggable": "no-drag" } as React.CSSProperties}
        className="flex items-center gap-2"
      >
        <ActionButtons
          allowMinimize={true}
          allowMaximize={false}
          allowQuit={true}
        />
      </div>
    </div>
  );
}