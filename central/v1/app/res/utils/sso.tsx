"use client";

import React from "react";

export function SsoLoginButton() {
  const handleSsoLogin = () => {
    const clientId = "fehmi-portal"; // Your registered client identifier
    const redirectUri = window.location.origin + "/auth/callback"; // Current app's callback
    const state = crypto.randomUUID(); // CSRF protection token

    // Store state in sessionStorage to verify during callback
    sessionStorage.setItem("sso_state", state);

    // Build the SSO Authorization URL pointing to your auth server
    const ssoUrl = new URL("https://auth.fehmi.cloud/api/sso");
    ssoUrl.searchParams.set("client_id", clientId);
    ssoUrl.searchParams.set("redirect_uri", redirectUri);
    ssoUrl.searchParams.set("state", state);
    ssoUrl.searchParams.set("response_type", "code");

    // Redirect user to central authentication server
    window.location.href = ssoUrl.toString();
  };

  return (
    <button
      onClick={handleSsoLogin}
      className="px-4 py-2 bg-slate-900 text-white dark:bg-white dark:text-slate-900 font-medium rounded-lg shadow hover:opacity-90 transition"
    >
      Sign in with Fehmi SSO
    </button>
  );
}