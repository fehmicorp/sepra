"use client";

import { useState } from "react";
import { useRouter, useSearchParams } from "next/navigation";
import Link from "next/link";
import { AuthLayout } from "@res/component/ui/layout";
import { AuthBanner } from "@res/component/ui/banner";
import { AuthCard } from "@res/component/ui/card";
import { InputField } from "@res/component/ui/input";
import { SubmitButton } from "@res/component/ui/sbumit_button";

export default function SetPassword() {
  const router = useRouter();
  const searchParams = useSearchParams();

  const [password, setPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");
  const [loading, setLoading] = useState(false);
  const [errorMsg, setErrorMsg] = useState("");

  const userContext = {
    name: searchParams.get("name") || "",
    company: searchParams.get("company") || "fehmi",
    dob: searchParams.get("dob") || "",
    city: searchParams.get("city") || "",
  };

  // Dynamic date sequence generator starting from 2026 moving forward
  const currentYear = new Date().getFullYear();
  const forbiddenYears = Array.from({ length: 15 }, (_, i) => String(currentYear + i));

  // Policy validation checks
  const policyChecks = {
    minLength: password.length >= 8,
    firstCaps: /^[A-Z]/.test(password),
    hasNumber: /\d/.test(password),
    hasSpecial: /[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]/.test(password),
    noForbiddenWords: !/password|passkey/i.test(password),
    noYears: !forbiddenYears.some((yr) => password.includes(yr)),
    noPersonalData: (() => {
      if (!password) return true;
      const lowerPassword = password.toLowerCase();
      const forbiddenTokens = [
        userContext.name,
        userContext.company,
        userContext.dob,
        userContext.city,
      ]
        .filter(Boolean)
        .map((str) => str.toLowerCase());

      return !forbiddenTokens.some((token) => lowerPassword.includes(token));
    })(),
  };

  const handleSetPassword = async (e: React.FormEvent) => {
    e.preventDefault();
    setErrorMsg("");

    if (!password || !confirmPassword) {
      setErrorMsg("Please fill out both fields.");
      return;
    }

    if (password !== confirmPassword) {
      setErrorMsg("Credentials do not match.");
      return;
    }

    if (!policyChecks.minLength) {
      setErrorMsg("Must be at least 8 characters long.");
      return;
    }

    if (!policyChecks.firstCaps) {
      setErrorMsg("Must start with a capital letter.");
      return;
    }

    if (!policyChecks.hasNumber) {
      setErrorMsg("Must contain at least 1 number.");
      return;
    }

    if (!policyChecks.hasSpecial) {
      setErrorMsg("Must contain at least 1 special character.");
      return;
    }

    if (!policyChecks.noForbiddenWords) {
      setErrorMsg("Cannot contain reserved words like 'password' or 'passkey'.");
      return;
    }

    if (!policyChecks.noYears) {
      setErrorMsg("Cannot contain the current year or upcoming years.");
      return;
    }

    if (!policyChecks.noPersonalData) {
      setErrorMsg("Cannot contain personal data (name, company, DOB, or city).");
      return;
    }

    setLoading(true);

    try {
      console.log("Setting credentials successfully");
      // TODO: Submit payload to API
      // await api.setPassword({ password });
      router.push("/auth/signin");
    } catch (err) {
      setErrorMsg("Failed to update credentials. Please try again.");
    } finally {
      setLoading(false);
    }
  };

  return (
    <AuthLayout>
      <AuthBanner title="Fehmi Authentication Engine" iconName="ShieldCheck" />
      <AuthCard
        title="Setup a Credentials"
        subtitle="Create a secure credential for your account"
        errorMsg={errorMsg}
      >
        <form onSubmit={handleSetPassword} className="mt-6 space-y-4">
          <InputField
            label="New Password"
            type="password"
            name="password"
            autoFocus
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            placeholder="••••••••••••"
          />

          <InputField
            label="Confirm Password"
            type="password"
            name="confirmPassword"
            value={confirmPassword}
            onChange={(e) => setConfirmPassword(e.target.value)}
            placeholder="••••••••••••"
          />

          {/* Validation Tracker */}
          <div className="p-3 rounded-xl bg-[var(--bg-elevated,#1e293b)] border border-[var(--border-color,#334155)] space-y-1.5 text-xs">
            <p className="font-semibold text-[var(--text-main)] mb-2">
              Security Requirements:
            </p>
            <ul className="space-y-1">
              <li className={policyChecks.minLength ? "text-emerald-400" : "text-[var(--text-muted)]"}>
                {policyChecks.minLength ? "✓" : "○"} At least 8 characters
              </li>
              <li className={policyChecks.firstCaps ? "text-emerald-400" : "text-[var(--text-muted)]"}>
                {policyChecks.firstCaps ? "✓" : "○"} Starts with a capital letter
              </li>
              <li className={policyChecks.hasNumber ? "text-emerald-400" : "text-[var(--text-muted)]"}>
                {policyChecks.hasNumber ? "✓" : "○"} At least 1 number
              </li>
              <li className={policyChecks.hasSpecial ? "text-emerald-400" : "text-[var(--text-muted)]"}>
                {policyChecks.hasSpecial ? "✓" : "○"} At least 1 special character
              </li>
            </ul>
          </div>

          <div className="pt-2">
            <SubmitButton loading={loading} loadingText="Saving...">
              Save Credentials
            </SubmitButton>
          </div>
        </form>
      </AuthCard>
    </AuthLayout>
  );
}