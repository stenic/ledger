import { Box } from "@mui/material";
import Header from "../../components/Header";
import { useAuth } from "react-oidc-context";

const Dashboard = () => {
  const auth = useAuth();

  return (
    <Box m="20px">
      <Header title="Dashboard" subtitle="Welcome to your dashboard" />
      dash Hello {auth.user?.profile.preferred_username}{" "}
    </Box>
  );
};

export default Dashboard;
