import {
  createContext,
  type ReactNode,
  useCallback,
  useContext,
  useMemo,
  useState,
} from "react";

const STORAGE_KEY = "educrm_sidebar_collapsed";

type AdminShellContextValue = {
  sidebarCollapsed: boolean;
  setSidebarCollapsed: (value: boolean) => void;
  toggleSidebar: () => void;
};

const AdminShellContext = createContext<AdminShellContextValue | null>(null);

function readInitialCollapsed(): boolean {
  try {
    return localStorage.getItem(STORAGE_KEY) === "1";
  } catch {
    return false;
  }
}

export function AdminShellProvider({ children }: { children: ReactNode }) {
  const [sidebarCollapsed, setSidebarCollapsedState] = useState(
    readInitialCollapsed,
  );

  const setSidebarCollapsed = useCallback((value: boolean) => {
    setSidebarCollapsedState(value);
    try {
      localStorage.setItem(STORAGE_KEY, value ? "1" : "0");
    } catch {
      /* ignore */
    }
  }, []);

  const toggleSidebar = useCallback(() => {
    setSidebarCollapsedState((prev) => {
      const next = !prev;
      try {
        localStorage.setItem(STORAGE_KEY, next ? "1" : "0");
      } catch {
        /* ignore */
      }
      return next;
    });
  }, []);

  const value = useMemo(
    () => ({
      sidebarCollapsed,
      setSidebarCollapsed,
      toggleSidebar,
    }),
    [sidebarCollapsed, setSidebarCollapsed, toggleSidebar],
  );

  return (
    <AdminShellContext.Provider value={value}>
      {children}
    </AdminShellContext.Provider>
  );
}

export function useAdminShell() {
  const ctx = useContext(AdminShellContext);
  if (!ctx) {
    throw new Error("useAdminShell must be used within AdminShellProvider");
  }
  return ctx;
}
