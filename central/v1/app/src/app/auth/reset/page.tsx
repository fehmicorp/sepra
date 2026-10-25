"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { AuthLayout } from "@res/component/ui/layout";
import { AuthBanner } from "@res/component/ui/banner";
import { AuthCard } from "@res/component/ui/card";
import { InputField } from "@res/component/ui/input";
import { SubmitButton } from "@res/component/ui/sbumit_button";
import Link from "next/link";

export default function ResetPassword() {
  const [email, setEmail] = useState("");
  const router = useRouter();
  const [loading, setLoading] = useState(false);
  const [errorMsg, setErrorMsg] = useState("");
  const [isSubmitted, setIsSubmitted] = useState(false);

  const handleResetPassword = async (e: React.FormEvent) => {
    e.preventDefault();
    setErrorMsg("");

    const inputEmail = email.trim();

    if (!inputEmail) {
      setErrorMsg("Please enter your email address.");
      return;
    }

    setLoading(true);

    try {
      console.log("Requesting password reset for:", inputEmail);
      // TODO: Call your password reset API endpoint
      // await api.resetPassword({ email: inputEmail });
      setIsSubmitted(true);
    } catch (err) {
      setErrorMsg("Failed to send reset link. Please try again.");
    } finally {
      setLoading(false);
    }
  };

  return (
    <AuthLayout>
      <AuthBanner title="Fehmi Authentication Engine" iconName="ShieldCheck" />
      <AuthCard
        title="Reset Password"
        subtitle={
          isSubmitted
            ? "Check your inbox for instructions"
            : "Enter your email to receive a password reset link"
        }
        errorMsg={errorMsg}
      >
        {isSubmitted ? (
          <div className="mt-6 space-y-5 text-center">
            <div className="p-4 rounded-xl bg-[var(--bg-elevated,#1e293b)] border border-[var(--border-color,#334155)] text-xs text-[var(--text-muted)] leading-relaxed">
              We sent a password reset link to{" "}
              <span className="font-semibold text-[var(--text-main)]">
                {email}
              </span>
              . Please check your inbox and follow the instructions.
            </div>

            <SubmitButton onClick={() => router.push("/auth/signin")}>
              Return to Sign In
            </SubmitButton>
          </div>
        ) : (
          <form onSubmit={handleResetPassword} className="mt-6 space-y-4">
            <InputField
              label="Email Address"
              type="email"
              name="email"
              autoFocus
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              placeholder="admin@youremail.in"
            />
            <div className="pt-2">
              <SubmitButton
                loading={loading}
                loadingText="Sending reset link..."
              >
                Send Reset Link
              </SubmitButton>
            </div>
          </form>
        )}

        <div className="mt-6 flex justify-between items-center text-xs text-[var(--text-muted)]">
          <Link
            href="/auth/signin"
            className="hover:text-[var(--text-main)] transition hover:underline"
          >
            ← Back to Sign In
          </Link>
          <Link
            href="/auth/signup"
            className="text-[var(--sepra-primary,#0284c7)] hover:underline font-medium"
          >
            Create Account
          </Link>
        </div>
      </AuthCard>
    </AuthLayout>
  );
}