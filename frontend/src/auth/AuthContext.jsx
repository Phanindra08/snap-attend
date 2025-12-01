import { createContext, useContext, useState, useEffect } from "react";
import { TOKEN_KEY, USER_KEY } from "../config";

const AuthContext = createContext(null);

export function AuthProvider({ children }) {
  const [user, setUserState] = useState(null);

  useEffect(() => {
    const stored = localStorage.getItem(USER_KEY);
    if (stored) {
      try {
        setUserState(JSON.parse(stored));
      } catch {
        setUserState(null);
      }
    }
  }, []);

  const setUser = (userObj) => {
    setUserState(userObj);
    if (userObj) {
      localStorage.setItem(USER_KEY, JSON.stringify(userObj));
    } else {
      localStorage.removeItem(USER_KEY);
    }
  };

  const logout = () => {
    localStorage.removeItem(TOKEN_KEY);
    localStorage.removeItem(USER_KEY);
    setUserState(null);
    window.location.href = "/login";
  };

  return (
    <AuthContext.Provider value={{ user, setUser, logout }}>
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  return useContext(AuthContext);
}
