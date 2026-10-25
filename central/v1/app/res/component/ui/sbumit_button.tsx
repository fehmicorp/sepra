"use client";

import React from "react";

interface SubmitButtonProps extends React.ButtonHTMLAttributes<HTMLButtonElement> {
  loading?: boolean;
  loadingText?: string;
  children: React.ReactNode;
}

export const SubmitButton: React.FC<SubmitButtonProps> = ({
  loading = false,
  loadingText = "Processing...",
  children,
  ...props
}) => {
  return (
    <button
      type="submit"
      disabled={loading || props.disabled}
      className="w-full bg-[var(--sepra-primary,#0284c7)] text-white font-medium text-sm rounded-xl py-3 px-4 shadow-lg transition-all hover:opacity-90 disabled:opacity-50 cursor-pointer"
      {...props}
    >
      {loading ? loadingText : children}
    </button>
  );
};