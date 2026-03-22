import { Outlet } from "react-router-dom";
import Navbar from "./navbar";
import FlyioBanner from "./flyio-banner";
import { useFlyioStore } from "@/store/flyio";
import { getFlyioStatus } from "@/api/flyio";
import { useEffect } from "react";

export default function Layout() {
  const { setIncidents } = useFlyioStore();

  useEffect(() => {
    const fetchStatus = () => {
      getFlyioStatus()
        .then((data) => setIncidents(data.incidents))
        .catch(() => {});
    };

    fetchStatus();
    const interval = setInterval(fetchStatus, 60_000);
    return () => clearInterval(interval);
  }, [setIncidents]);

  return (
    <div>
      <Navbar />
      <FlyioBanner />
      <Outlet />
    </div>
  );
}
