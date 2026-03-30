import { useFlyioStore } from "@/store/flyio";
import { AlertTriangle } from "lucide-react";

export default function FlyioBanner() {
  const { incidents } = useFlyioStore();
  const active = incidents.filter((i) => !i.resolved);

  if (active.length === 0) return null;

  const impactOrder: Record<string, number> = {
    critical: 4,
    major: 3,
    minor: 2,
    none: 1,
  };

  const worst = active.reduce((prev, curr) =>
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
      <AlertTriangle size={14} className="shrink-0" />
      <span className="font-semibold">Fly.io</span>
      {active.length > 1 && (
        <span
          className={`rounded-full px-1.5 py-0.5 text-xs font-bold leading-none ${
            isSevere
              ? "bg-red-500/20 text-red-400"
              : "bg-yellow-500/20 text-yellow-400"
          }`}
        >
          {active.length}
        </span>
      )}
      <span className="opacity-70">·</span>
      <span className="truncate">{worst.name}</span>
      <span
        className={`shrink-0 rounded px-1.5 py-0.5 text-xs font-medium ${
          isSevere
            ? "bg-red-500/20 text-red-300"
            : "bg-yellow-500/20 text-yellow-300"
        }`}
      >
        {worst.status}
      </span>
      {isSevere && (
        <span className="hidden shrink-0 opacity-60 sm:inline">
          — deployments may be affected
        </span>
      )}
      <a
        href="https://status.fly.io"
        target="_blank"
        rel="noopener noreferrer"
        className="ml-auto shrink-0 text-xs opacity-50 hover:opacity-90"
      >
        status.fly.io ↗
      </a>
    </div>
  );
}
