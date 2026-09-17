import type { NextConfig } from "next";

const isExport = process.env.OUTPUT === "local" || process.env.OUTPUT === "export";
const outputMode = isExport ? "export" : "standalone";
const rawApp = process.env.APP ? process.env.APP.trim() : "";
const baseHref = rawApp ? `/${rawApp.replace(/^\/+|\/+$/g, '')}` : "";
const nextConfig: NextConfig = {
  output: outputMode as "export" | "standalone",
  basePath: baseHref,
  assetPrefix: baseHref,
  trailingSlash: true,
  images: {
    unoptimized: true,
  },
};

export default nextConfig;