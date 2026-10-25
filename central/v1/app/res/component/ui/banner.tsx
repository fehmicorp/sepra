"use client";

import React from "react";
import Icon from "@res/utils/icon";

interface AuthBannerProps {
  iconName?: string;
  title?: string;
}

export const AuthBanner: React.FC<AuthBannerProps> = ({
  iconName = "",
  title = "🚀 Fehmi Authentication Engine",
}) => {
  return (
    <div className="bg-gradient-to-l bg-[var(--bg-banner)] text-[var(--text-banner)] px-4 py-3 rounded-lg text-sm font-medium flex items-center justify-between border border-[var(--border-card)] shadow-sm">
      <span className="flex items-center gap-2">
        {iconName !== "" && <Icon name={iconName} size={16} />} 
        {title}
      </span>
    </div>
  );
};