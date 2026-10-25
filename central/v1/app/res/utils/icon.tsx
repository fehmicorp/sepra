import React from "react";
import * as Icons from "lucide-react";
import * as LocalIcons from "@res/icons";

export interface IconProps extends React.SVGProps<SVGSVGElement> {
  name: string;
  size?: number | string;
  color?: string;
  className?: string;
}

const createSvgComponent = (svgString: string) => {
  return function CustomSvgIcon({ size = 24, className = "", ...props }: IconProps) {
    // 1. Remove hardcoded fills and width/height attributes
    // 2. Ensure fill="currentColor" is present
    const formattedSvg = svgString
      .replace(/fill="[^"]*"/gi, 'fill="currentColor"')
      .replace(/width="[^"]*"/gi, "")
      .replace(/height="[^"]*"/gi, "")
      .replace(/<svg/i, `<svg width="${size}" height="${size}"`);

    return (
      <span
        style={{ width: size, height: size }}
        className={`inline-flex items-center justify-center shrink-0 text-slate-900 dark:text-white ${className}`}
        dangerouslySetInnerHTML={{ __html: formattedSvg }}
        {...props}
      />
    );
  };
};

// Build icon map combining Lucide icons and custom SVGs
export const iconMap: Record<string, React.ElementType> = Object.entries(Icons).reduce(
  (acc, [key, value]) => {
    if (typeof value === "function" || typeof value === "object") {
      acc[key] = value as React.ElementType;
    }
    return acc;
  },
  {} as Record<string, React.ElementType>
);

// 2. Map ALL local SVG string exports automatically
Object.entries(LocalIcons).forEach(([key, value]) => {
  if (typeof value === "string" && value.includes("<svg")) {
    iconMap[key] = createSvgComponent(value);
  }
});

// Resilient fallback for Lucide version changes (v0.300+ renamed HelpCircle to CircleHelp)
export const FallbackIcon = Icons.CircleHelp || Icons.HelpCircle || (() => null);

export const Icon: React.FC<IconProps> = ({ name, size = 24, ...props }) => {
  const IconComponent = iconMap[name] || FallbackIcon;
  return <IconComponent size={size} {...props} />;
};

export default Icon;