"use client";

import React from "react";

interface InputFieldProps extends React.InputHTMLAttributes<HTMLInputElement> {
  label: string;
  actionLabel?: string;
  onActionClick?: () => void;
}

export const InputField: React.FC<InputFieldProps> = ({
  label,
  actionLabel,
  onActionClick,
  ...props
}) => {
  return (
    <div className="space-y-2">
      <div className="flex justify-between items-center mb-1">
        <label className="block text-xs font-medium text-[var(--text-muted)] tracking-wider">
          {label}
        </label>
        {actionLabel && (
          <button
            type="button"
            onClick={onActionClick}
            className="text-xs text-[var(--sepra-primary,#0284c7)] hover:underline"
          >
            {actionLabel}
          </button>
        )}
      </div>
      <input
        {...props}
        className="w-full bg-[var(--bg-elevated)] border border-[var(--border-color)] text-[var(--text-main)] text-sm rounded-xl px-4 py-3 outline-none focus:border-[var(--sepra-primary)] transition"
      />
    </div>
  );
};