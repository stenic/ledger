import React, { useState } from "react";
import { useAuth } from "react-oidc-context";
import { Box, Typography } from "@mui/material";
import Button from "@mui/material/Button";
import TextField from "@mui/material/TextField";
import FormControlLabel from "@mui/material/FormControlLabel";
import Checkbox from "@mui/material/Checkbox";
import Link from "@mui/material/Link";
import Grid from "@mui/material/Grid";
import Card from "@mui/material/Card";
import CardContent from "@mui/material/CardContent";
import Avatar from "@mui/material/Avatar";
import LockOutlinedIcon from "@mui/icons-material/LockOutlined";
import Header from "../../components/Header";
import Divider from "@mui/material/Divider";
import { useTranslation } from "react-i18next";

const LoginForm = ({ defaultOpen = false }: { defaultOpen?: boolean }) => {
  const [formOpen, setFormOpen] = useState(defaultOpen);
  const { t } = useTranslation();

  const openForm = () => {
    setFormOpen(true);
  };

  const handleSubmit = (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    const data = new FormData(event.currentTarget);
    console.log({
      email: data.get("email"),
      password: data.get("password"),
    });
  };

  return (
    <Box component="form" onSubmit={handleSubmit} noValidate sx={{ mt: 1 }}>
      {!formOpen && (
        <Box
          sx={{
            display: "flex",
            justifyContent: "center",
          }}
        >
          <Link href="#" onClick={openForm}>
            {t("login_local")}
          </Link>
        </Box>
      )}

      {formOpen && (
        <>
          <TextField
            margin="dense"
            required
            id="email"
            label={t("login_email")}
            name="email"
            autoComplete="email"
            fullWidth
            autoFocus
          />
          <TextField
            margin="dense"
            required
            fullWidth
            name="password"
            label={t("login_password")}
            type="password"
            id="password"
            autoComplete="current-password"
          />
          <div>
            <FormControlLabel
              control={<Checkbox value="remember" color="primary" />}
              label={t("login_remember")}
            />
          </div>
          <Button
            type="submit"
            fullWidth
            variant="contained"
            sx={{ mt: 3, mb: 2 }}
          >
            {t("login_signin")}
          </Button>
          <Grid container>
            <Grid item xs>
              <Link href="#" variant="body2" color="secondary">
                {t("login_forgot_password")}
              </Link>
            </Grid>
            <Grid item>
              <Link href="#" variant="body2" color="secondary">
                {t("login_signup")}
              </Link>
            </Grid>
          </Grid>
        </>
      )}
    </Box>
  );
};

const LoginOidc = () => {
  const auth = useAuth();
  const { t } = useTranslation();

  return (
    <Button
      type="submit"
      fullWidth
      color="secondary"
      variant="contained"
      sx={{ mt: 3, mb: 0 }}
      onClick={() => void auth.signinRedirect()}
    >
      {t("login_oidc")}
    </Button>
  );
};

const Login = () => {
  const { t } = useTranslation();

  return (
    <Box
      sx={{
        marginTop: 8,
        display: "flex",
        flexDirection: "column",
        alignItems: "center",
      }}
    >
      <Header title="Login" hidden={true} />
      <Card sx={{ minWidth: 275 }}>
        <CardContent>
          <Box
            sx={{
              display: "flex",
              flexDirection: "column",
              alignItems: "center",
            }}
          >
            <Avatar sx={{ m: 1, bgcolor: "secondary.main" }}>
              <LockOutlinedIcon />
            </Avatar>
            <Typography component="h1" variant="h5">
              {t("login_signin")}
            </Typography>
          </Box>

          <LoginOidc />
          <Divider sx={{ mt: 3, mb: 3 }} />
          <LoginForm />
        </CardContent>
      </Card>
    </Box>
  );
};

export default Login;
