// import LedgerAdmin from "./admin";
import { Routes, Route } from "react-router-dom";
import TopBar from "./scenes/TopBar";
import SideBar from "./scenes/SideBar";
import Actions from "./scenes/Actions";
import Dashboard from "./scenes/dashboard";
import Settings from "./scenes/settings";
import Versions from "./scenes/versions";
import LastVersions from "./scenes/last";
import Versions2 from "./scenes/versions2";
import Profile from "./scenes/profile";
import Chart from "./scenes/chart";
import Login from "./scenes/login";
import Feed from "./scenes/feed";
import Toolbar from "@mui/material/Toolbar";
import { useAuth } from "react-oidc-context";
import { SnackbarProvider } from "notistack";
import Websockets from "./components/Websocket";
import "./components/I18n";
import { useTranslation } from "react-i18next";

const App = () => {
  const auth = useAuth();
  const { t } = useTranslation("aa");

  switch (auth.activeNavigator) {
    case "signinSilent":
      return <div>{t("signin_silent")}</div>;
    case "signoutRedirect":
      return <div>{t("app_signout_redirect")}</div>;
  }

  if (auth.isLoading) {
    return <div>{t("app_loading")}</div>;
  }

  if (auth.error) {
    if (auth.error.message === "Token is not active") {
      window.location.reload();
    }
    return <div>{t("app_error", { error: auth.error.message })}</div>;
  }

  if (auth.isAuthenticated) {
    return (
      <>
        <SnackbarProvider maxSnack={3}>
          <div className="app">
            <SideBar />
            <main className="content">
              <TopBar />
              <Toolbar />
              <Routes>
                <Route path="/" element={<Dashboard />} />
                <Route path="/settings" element={<Settings />} />
                <Route path="/versions2" element={<Versions2 />} />
                <Route path="/versions" element={<Versions />} />
                <Route path="/last" element={<LastVersions />} />
                <Route path="/profile" element={<Profile />} />
                <Route path="/feed" element={<Feed />} />
                <Route path="/chart" element={<Chart />} />
              </Routes>
              <Actions />
            </main>
          </div>
          <Websockets />
        </SnackbarProvider>
      </>
    );
  }

  return <Login />;
};

export default App;
