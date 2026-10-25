"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";
import { AuthLayout } from "@res/component/ui/layout";
import { AuthBanner } from "@res/component/ui/banner";
import { AuthCard } from "@res/component/ui/card";
import { InputField } from "@res/component/ui/input";
import { SubmitButton } from "@res/component/ui/sbumit_button";
import { SocialButtonsCircle, SocialButtonsHorizontal, SocialButtonsSquare, SocialDivider } from "@res/component/ui/socialbutton";
import { SsoLoginButton } from "@res/utils/sso";

export default function SignIn() {
  const [auth, setAuth] = useState({ UserInput: "", Password: "" });
  const router = useRouter();
  const [step, setStep] = useState<1 | 2>(1);
  const [loading, setLoading] = useState(false);
  const [errorMsg, setErrorMsg] = useState("");

  const handleStepOne = async (e: React.FormEvent) => {
    e.preventDefault();
    setErrorMsg("");

    const input = auth.UserInput.trim();
    if (!input) {
      setErrorMsg("Please enter your email, username, or phone number.");
      return;
    }

    setLoading(true);
    try {
      console.log("Checking identifier:", input);
      setStep(2);
    } catch (err) {
      setErrorMsg("User not found. Please check your input.");
    } finally {
      setLoading(false);
    }
  };

  const handleStepTwo = async (e: React.FormEvent) => {
    e.preventDefault();
    setErrorMsg("");

    if (!auth.Password) {
      setErrorMsg("Please enter your password.");
      return;
    }

    setLoading(true);
    try {
      console.log("Authenticating:", { identifier: auth.UserInput });
      router.push("/dashboard");
    } catch (err) {
      setErrorMsg("Invalid password. Please try again.");
    } finally {
      setLoading(false);
    }
  };

  return (
    <AuthLayout>
      <AuthBanner title="Fehmi Authentication Engine" iconName="ShieldCheck" />

      <AuthCard
        title="Sign In"
        subtitle={
          step === 1
            ? "Enter your account identifier to continue"
            : `Welcome back, ${auth.UserInput}`
        }
        errorMsg={errorMsg}
      >
        <form
          onSubmit={step === 1 ? handleStepOne : handleStepTwo}
          className="mt-6 space-y-5"
        >
          {step === 1 ? (
            <InputField
              label="Email & Username"
              type="text"
              autoFocus
              value={auth.UserInput}
              onChange={(e) =>
                setAuth((prev) => ({ ...prev, UserInput: e.target.value }))
              }
              placeholder="admin@youremail.in"
            />
          ) : (
            <InputField
              label="Password"
              type="password"
              autoFocus
              actionLabel="Change user"
              onActionClick={() => {
                setStep(1);
                setErrorMsg("");
              }}
              value={auth.Password}
              onChange={(e) =>
                setAuth((prev) => ({ ...prev, Password: e.target.value }))
              }
              placeholder="•••••••••"
            />
          )}

          <SubmitButton loading={loading}>
            {step === 1 ? "Next" : "Sign In"}
          </SubmitButton>
        </form>

        <SocialDivider className="my-8 mx-6" />
        <SocialButtonsCircle className="gap-10" onSelectProvider={(id) => console.log("Provider:", id)} />

        <div className="mt-6 text-center">
          {step === 2 ? (
            <button
              type="button"
              onClick={() => {
                setErrorMsg("");
                setStep(1);
              }}
              className="text-xs text-[var(--text-muted)] hover:text-[var(--text-main)] transition"
            >
              ← Back to Email
            </button>
          ) : (
            <div className="flex justify-between items-center">
              <div className="text-center">
                <span className="text-xs text-[var(--text-muted)]">
                  Need help?{" "}
                </span>
                <Link
                  href="/auth/reset"
                  className="text-xs text-[var(--sepra-primary,#0284c7)] hover:text-blue-400 hover:underline font-medium"
                >
                  Forgot Password
                </Link>
              </div>
              <Link
                href="/auth/signup"
                className="text-xs text-[var(--text-muted)] hover:text-[var(--text-main)] hover:underline font-medium"
              >
                Create new Account
              </Link>
            </div>
          )}
        </div>
      </AuthCard>
    </AuthLayout>
  );
}