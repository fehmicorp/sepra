"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";
import { AuthLayout } from "@res/component/ui/layout";
import { AuthBanner } from "@res/component/ui/banner";
import { AuthCard } from "@res/component/ui/card";
import { InputField } from "@res/component/ui/input";
import { SubmitButton } from "@res/component/ui/sbumit_button";
import { SocialButtonsCircle, SocialDivider } from "@res/component/ui/socialbutton";

export default function Signup() {
  const [email, setEmail] = useState("");
  const router = useRouter();
  const [loading, setLoading] = useState(false);
  const [errorMsg, setErrorMsg] = useState("");

  const handleSignup = async (e: React.FormEvent) => {
    e.preventDefault();
    setErrorMsg("");

    const inputEmail = email.trim();

    if (!inputEmail) {
      setErrorMsg("Please enter your email address.");
      return;
    }

    setLoading(true);

    try {
      console.log("Creating account for:", inputEmail);
      // TODO: Call your magic link or registration API endpoint
      // await api.signup({ email: inputEmail });
      router.push("/auth/signin");
    } catch (err) {
      setErrorMsg("Failed to create account. Please try again.");
    } finally {
      setLoading(false);
    }
  };

  return (
    <AuthLayout>
      <AuthBanner title="Fehmi Authentication Engine" iconName="ShieldCheck" />
      <AuthCard
        title="Sign Up"
        subtitle="Create your account to access orchestration console"
        errorMsg={errorMsg}
      >
        <form onSubmit={handleSignup} className="mt-6 space-y-4">
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
            <SubmitButton loading={loading} loadingText="Creating account...">
              Create Account
            </SubmitButton>
          </div>
        </form>

        <SocialDivider className="my-8 mx-6" />
        <SocialButtonsCircle
          className="gap-10"
          page="Sign up"
          onSelectProvider={(id) => console.log("Provider:", id)}
        />

        <div className="mt-6 text-center">
          <span className="text-xs text-[var(--text-muted)]">
            Already have an account?{" "}
          </span>
          <Link
            href="/auth/signin"
            className="text-xs text-[var(--sepra-primary,#0284c7)] hover:underline font-medium"
          >
            Sign In
          </Link>
        </div>
      </AuthCard>
    </AuthLayout>
  );
}