"use client";

import React from "react";

interface AuthCardProps {
  title: string;
  subtitle: string;
  errorMsg?: string;
  children: React.ReactNode;
}

export const AuthCard: React.FC<AuthCardProps> = ({
  title,
  subtitle,
  errorMsg,
  children,
}) => {
  return (
    <div className="bg-gradient-to-b w-full from-[var(--bg-hero-from)] to-[var(--bg-hero-to)] p-8 rounded-2xl border border-[var(--border-card)] shadow-xl">
      <h1 className="text-xl font-semibold tracking-[.30rem] text-[var(--text-main)] font-mono">
        {title}
      </h1>
      <p className="text-sm text-[var(--text-muted)] mt-1 mb-6">{subtitle}</p>

      {errorMsg && (
        <div className="mb-4 p-3 bg-red-500/10 border border-red-500/30 text-red-400 text-xs rounded-xl">
          {errorMsg}
        </div>
      )}

      {children}
    </div>
  );
};