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

const App = () => {
  const auth = useAuth();

  switch (auth.activeNavigator) {
    case "signinSilent":
      return <div>Signing you in...</div>;
    case "signoutRedirect":
      return <div>Signing you out...</div>;
  }

  if (auth.isLoading) {
    return <div>Loading...</div>;
  }

  if (auth.error) {
    if (auth.error.message === "Token is not active") {
      window.location.reload();
    }
    return <div>Oops... {auth.error.message}</div>;
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
