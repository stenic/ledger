import {
  Box,
  FormControl,
  FormControlLabel,
  FormGroup,
  InputLabel,
  MenuItem,
  Select,
  SelectChangeEvent,
  Switch,
  Typography,
} from "@mui/material";
import Header from "../../components/Header";
import { useEffect, useState } from "react";
import { useSnackbar } from "notistack";
import { useTranslation } from "react-i18next";
import i18n from "i18next";

const Settings = () => {
  const { t } = useTranslation();
  const { enqueueSnackbar } = useSnackbar();
  const [notificationEnabled, setNotificationEnabled] = useState(
    localStorage.getItem("setting.notification.enabled") === "granted"
  );
  const [language, setLanguage] = useState(i18n.language);

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

  const handleChange: (event: SelectChangeEvent) => void = (event) => {
    setLanguage(event.target.value);
    i18n.changeLanguage(event.target.value);
  };

  return (
    <Box m="20px">
      <Header title={t("settings_title")} subtitle={t("settings_subtitle")} />

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
          <Typography variant="h3">{t("settings_notifications")}</Typography>
          <FormControlLabel
            control={<Switch checked={notificationEnabled} />}
            label={t("settings_notifications_toggle")}
            onChange={handleNotificationSetting}
          />
          <Typography variant="h3">{t("settings_language")}</Typography>
          <FormControl fullWidth>
            <InputLabel id="demo-simple-select-label">Language</InputLabel>
            <Select
              labelId="demo-simple-select-label"
              id="demo-simple-select"
              value={language}
              label="Language"
              onChange={handleChange}
            >
              <MenuItem value={""}>Auto-detect</MenuItem>
              <MenuItem value={"en"}>EN</MenuItem>
              <MenuItem value={"nl"}>NL</MenuItem>
            </Select>
          </FormControl>
        </FormGroup>
      </Box>
    </Box>
  );
};

export default Settings;
