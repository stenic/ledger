import {
  Box,
  FormControlLabel,
  FormGroup,
  Switch,
  Typography,
} from "@mui/material";
import Header from "../../components/Header";
import { useEffect, useState } from "react";
import { useSnackbar } from "notistack";

const Settings = () => {
  const { enqueueSnackbar } = useSnackbar();
  const [notificationEnabled, setNotificationEnabled] = useState(
    localStorage.getItem("setting.notification.enabled") === "granted"
  );

  useEffect(() => {
    if (localStorage.getItem("setting.notification.enabled") === null) {
      setNotificationEnabled(Notification.permission === "granted");
    }
  }, []);

  const notificationStorage = async (enabled: boolean) => {
    if (enabled) {
      switch (Notification.permission) {
        case "default":
          await Notification.requestPermission().then((permission) => {
            setNotificationEnabled(permission === "granted");
          });
          return;
        case "denied":
          enqueueSnackbar("Cannot enable notifications!", { variant: "error" });
          setNotificationEnabled(false);
          return;
      }
    }
    localStorage.setItem(
      "setting.notification.enabled",
      enabled ? "granted" : "denield"
    );
  };

  useEffect(() => {
    notificationStorage(notificationEnabled);
  }, [notificationEnabled]);

  const handleNotificationSetting: (
    event: React.SyntheticEvent,
    checked: boolean
  ) => void = (event, checked) => {
    setNotificationEnabled(checked);
  };

  return (
    <Box m="20px">
      <Header title="Settings" subtitle="Application settings" />

      <Box
        sx={{
          display: "grid",
          gap: "10px",
          gridTemplateColumns: "repeat(12, 1fr)",
          "& > div": {
            bgcolor: "background.paper",
            p: 2,
          },
        }}
      >
        <FormGroup
          sx={{
            gridColumn: "span 6",
          }}
        >
          <Typography variant="h3">Notifications</Typography>
          <FormControlLabel
            control={<Switch checked={notificationEnabled} />}
            label="Enable browser notifications"
            onChange={handleNotificationSetting}
          />
        </FormGroup>
      </Box>
    </Box>
  );
};

export default Settings;
