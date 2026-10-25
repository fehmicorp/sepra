"use client";

import React from "react";

export const AuthLayout: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  return (
    <div className="min-h-screen flex flex-col justify-top mt-4 md:mt-10 lg:mt-16 items-center bg-[var(--bg-base,#0f172a)] text-[var(--text-main,#f8fafc)] p-4">
      <div className="w-full max-w-md space-y-6">{children}</div>
    </div>
  );
};