import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import "./styles/index.css";
import "./styles/font-faces.css";
import App from "./App.tsx";
import { ThemeProvider } from "@/components/theme-provider";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { Toaster } from "react-hot-toast";

const queryClient = new QueryClient();
const loadingSpinner = document.getElementById("spinner-start");

if (loadingSpinner) {
  loadingSpinner.style.display = "none";
}

createRoot(document.getElementById("root")!).render(
  <StrictMode>
    <QueryClientProvider client={queryClient}>
      <ThemeProvider defaultTheme="dark" storageKey="vite-ui-theme">
        <Toaster
          position="bottom-right"
          reverseOrder={false}
          gutter={8}
          containerClassName=""
          containerStyle={{}}
          toastOptions={{
            className: "",
            duration: 10000,
            style: {
              background: "#27272A",
              color: "#E4E4E7",
              borderRadius: "13px",
            },
          }}
        />
        <App />
      </ThemeProvider>
    </QueryClientProvider>
  </StrictMode>
);
