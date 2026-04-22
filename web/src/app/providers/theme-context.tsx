import {
  createContext,
  type ReactNode,
  useCallback,
  useContext,
  useEffect,
  useMemo,
  useState,
} from "react";

export type Theme = "light" | "dark" | "system";

type ThemeContextValue = {
  theme: Theme;
  resolved: "light" | "dark";
  setTheme: (t: Theme) => void;
  toggle: () => void;
};

const ThemeContext = createContext<ThemeContextValue | null>(null);

const STORAGE_KEY = "educrm_theme";

function readStoredTheme(): Theme {
  const stored = localStorage.getItem(STORAGE_KEY) as Theme | null;
  if (stored === "light" || stored === "dark" || stored === "system")
    return stored;
  return "system";
}

export function ThemeProvider({ children }: { children: ReactNode }) {
  const [theme, setThemeState] = useState<Theme>(() => readStoredTheme());
  const [systemDark, setSystemDark] = useState(
    () => window.matchMedia("(prefers-color-scheme: dark)").matches,
  );

  useEffect(() => {
    const mql = window.matchMedia("(prefers-color-scheme: dark)");
    const onChange = () => setSystemDark(mql.matches);
    mql.addEventListener("change", onChange);
    return () => mql.removeEventListener("change", onChange);
  }, []);

  const resolved = useMemo<"light" | "dark">(() => {
    if (theme === "system") return systemDark ? "dark" : "light";
    return theme;
  }, [theme, systemDark]);

  useEffect(() => {
    document.documentElement.classList.toggle("dark", resolved === "dark");
    localStorage.setItem(STORAGE_KEY, theme);
  }, [theme, resolved]);

  const setTheme = useCallback((t: Theme) => setThemeState(t), []);

  const toggle = useCallback(() => {
    setThemeState((prev) => {
      const r = prev === "system" ? (systemDark ? "dark" : "light") : prev;
      return r === "dark" ? "light" : "dark";
    });
  }, [systemDark]);

  const value = useMemo(
    () => ({ theme, resolved, setTheme, toggle }),
    [theme, resolved, setTheme, toggle],
  );

  return (
    <ThemeContext.Provider value={value}>{children}</ThemeContext.Provider>
  );
}

export function useTheme() {
  const ctx = useContext(ThemeContext);
  if (!ctx) throw new Error("useTheme must be used within ThemeProvider");
  return ctx;
}
