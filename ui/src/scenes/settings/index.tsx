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

const LanguageForm = () => {
  const { t } = useTranslation();
  const [language, setLanguage] = useState(i18n.language);

  const handleChange: (event: SelectChangeEvent) => void = (event) => {
    setLanguage(event.target.value);
    i18n.changeLanguage(event.target.value);
  };

  return (
    <FormControl fullWidth>
      <InputLabel id="language-label">{t("settings_language")}</InputLabel>
      <Select
        labelId="language-label"
        value={language}
        label={t("settings_language")}
        onChange={handleChange}
      >
        <MenuItem value={""}>Browser</MenuItem>
        <MenuItem value={"en"}>English</MenuItem>
        <MenuItem value={"nl"}>Nederlands</MenuItem>
      </Select>
    </FormControl>
  );
};

const NotificationForm = () => {
  const { t } = useTranslation();
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
    <FormControlLabel
      control={<Switch checked={notificationEnabled} />}
      label={t("settings_notifications_toggle")}
      onChange={handleNotificationSetting}
    />
  );
};

const Settings = () => {
  const { t } = useTranslation();

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
          <NotificationForm />
        </FormGroup>
        <FormGroup
          sx={{
            gridColumn: "span 6",
          }}
        >
          <Typography variant="h3">{t("settings_language")}</Typography>
          <LanguageForm />
        </FormGroup>
      </Box>
    </Box>
  );
};

export default Settings;
