"use client";

import React from "react";
import Icon from "@res/utils/icon";

interface SocialProvider {
  id: "google" | "github" | "microsoft";
  label: string;
  iconName: string;
}

const PROVIDERS: SocialProvider[] = [
  { id: "google", label: "Google", iconName: "google_black" },
  { id: "microsoft", label: "Microsoft", iconName: "microsoft_black" },
  { id: "github", label: "GitHub", iconName: "github_black" },
];

interface BaseSocialProps {
  onSelectProvider?: (providerId: string) => void;
  disabled?: boolean;
  page?: string;
  className?: string;
}

// 1. Circle Variant (Compact Icon-Only)
export const SocialButtonsCircle: React.FC<BaseSocialProps> = ({
  onSelectProvider,
  disabled = false,
  page = "Sign in",
  className = "",
}) => {
  return (
    <div className={`flex items-center justify-center gap-4 ${className}`}>
      {PROVIDERS.map((provider) => (
        <button
          key={provider.id}
          type="button"
          disabled={disabled}
          onClick={() => onSelectProvider?.(provider.id)}
          title={`${page} with ${provider.label}`}
          className="w-12 h-12 cursor-pointer rounded-full bg-[var(--bg-elevated,#1e293b)] border border-[var(--border-color,#334155)] text-[var(--text-main,#f8fafc)] flex items-center justify-center transition-all hover:border-[var(--sepra-primary,#0284c7)] hover:bg-[var(--bg-card,#334155)] active:scale-95 disabled:opacity-50"
        >
          {/* Ensure your Icon component supports the iconName provided */}
          <Icon name={provider.iconName} size={20} />
        </button>
      ))}
    </div>
  );
};

// 2. Square Variant (Structured Icon Boxes)
export const SocialButtonsSquare: React.FC<BaseSocialProps> = ({
  onSelectProvider,
  disabled = false,
  page = "Sign in",
  className = "" ,
}) => {
  return (
    <div className={`grid grid-cols-3 gap-3 ${className}`}>
      {PROVIDERS.map((provider) => (
        <button
          key={provider.id}
          type="button"
          disabled={disabled}
          onClick={() => onSelectProvider?.(provider.id)}
          title={`${page} with ${provider.label}`}
          className="h-12 rounded-xl cursor-pointer bg-[var(--bg-elevated,#1e293b)] border border-[var(--border-color,#334155)] text-[var(--text-main,#f8fafc)] flex items-center justify-center gap-2 text-xs font-medium transition-all hover:border-[var(--sepra-primary,#0284c7)] hover:bg-[var(--bg-card,#334155)] active:scale-95 disabled:opacity-50"
        >
          <Icon name={provider.iconName} size={18} />
          <span>{provider.label}</span>
        </button>
      ))}
    </div>
  );
};

// 3. Horizontal Variant (Full Width Stack)
export const SocialButtonsHorizontal: React.FC<BaseSocialProps> = ({
  onSelectProvider,
  disabled = false,
  page = "Sign in",
  className = "",
}) => {
  return (
    <div className={`flex flex-col gap-3 w-full ${className}`}>
      {PROVIDERS.map((provider) => (
        <button
          key={provider.id}
          type="button"
          disabled={disabled}
          onClick={() => onSelectProvider?.(provider.id)}
          title={`${page} with ${provider.label}`}
          className="w-full h-11 px-4 cursor-pointer rounded-xl bg-[var(--bg-elevated,#1e293b)] border border-[var(--border-color,#334155)] text-[var(--text-main,#f8fafc)] flex items-center justify-center gap-3 text-sm font-medium transition-all hover:border-[var(--sepra-primary,#0284c7)] hover:bg-[var(--bg-card,#334155)] active:scale-[0.99] disabled:opacity-50"
        >
          <Icon name={provider.iconName} size={18} />
          <span>Continue with {provider.label}</span>
        </button>
      ))}
    </div>
  );
};

interface SocialDividerProps {
  label?: string;
  className?: string;
}

// Divider Component to separate standard auth from social auth
export const SocialDivider: React.FC<SocialDividerProps> = ({
  label = "OR",
  className = "my-6",
}) => {
  return (
    <div className={`flex flex-row items-center justify-center ${className}`}>
      <div className="w-[40%] border-t border-[var(--border-color,#334155)]" />
      <span className="px-6 text-[10px] font-semibold tracking-widest text-[var(--text-muted,#94a3b8)] uppercase">
        {label}
      </span>
      <div className="w-[40%] border-t border-[var(--border-color,#334155)]" />
    </div>
  );
};