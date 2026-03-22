import { useFlyioStore } from "@/store/flyio";
import { AlertTriangle } from "lucide-react";

export default function FlyioBanner() {
  const { incidents } = useFlyioStore();

  if (incidents.length === 0) return null;

  const impactOrder: Record<string, number> = {
    critical: 4,
    major: 3,
    minor: 2,
    none: 1,
  };

  const worst = incidents.reduce((prev, curr) =>
    (impactOrder[curr.impact] ?? 0) > (impactOrder[prev.impact] ?? 0) ? curr : prev
  );

  const isSevere = worst.impact === "critical" || worst.impact === "major";

  return (
    <div
      className={`flex items-center gap-2 border-b px-4 py-2 text-sm ${
        isSevere
          ? "border-red-500/30 bg-red-500/10 text-red-400"
          : "border-yellow-500/30 bg-yellow-500/10 text-yellow-400"
      }`}
    >
      <AlertTriangle size={15} className="shrink-0" />
      <span className="font-semibold">Fly.io incident:</span>
      <span>{worst.name}</span>
      <span className="opacity-60">({worst.status})</span>
      {isSevere && (
        <span className="opacity-70">— new deployments may be unavailable</span>
      )}
      <a
        href="https://status.fly.io"
        target="_blank"
        rel="noopener noreferrer"
        className="ml-auto shrink-0 text-xs underline opacity-60 hover:opacity-100"
      >
        status.fly.io ↗
      </a>
    </div>
  );
}
