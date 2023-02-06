import AppBar from "@mui/material/AppBar";
import Toolbar from "@mui/material/Toolbar";
import Typography from "@mui/material/Typography";
import * as React from "react";
import Box from "@mui/material/Box";
import MenuItem from "@mui/material/MenuItem";
import Menu from "@mui/material/Menu";
import AccountCircle from "@mui/icons-material/AccountCircle";
import { useAuth } from "react-oidc-context";
import Chip from "@mui/material/Chip";
import { Link } from "react-router-dom";
import Divider from "@mui/material/Divider";
import { useTranslation } from "react-i18next";

const TopBar = () => {
  const auth = useAuth();
  const { t } = useTranslation();

  const [anchorEl, setAnchorEl] = React.useState<HTMLDivElement | undefined>(
    undefined
  );

  const menuId = "primary-search-account-menu";
  const handleProfileMenuOpen = (event: React.MouseEvent<HTMLDivElement>) => {
    setAnchorEl(event.currentTarget);
  };

  const isMenuOpen = Boolean(anchorEl);

  const handleMenuClose = () => {
    setAnchorEl(undefined);
  };

  const renderMenu = (
    <Menu
      anchorEl={anchorEl}
      anchorOrigin={{
        vertical: "top",
        horizontal: "right",
      }}
      id={menuId}
      keepMounted
      transformOrigin={{
        vertical: "top",
        horizontal: "right",
      }}
      open={isMenuOpen}
      onClose={handleMenuClose}
      sx={{ mt: "45px", floodColor: "#1F2A40" }}
    >
      <MenuItem component={Link} to="/profile" onClick={handleMenuClose}>
        {t("app_profile")}
      </MenuItem>
      <MenuItem onClick={handleMenuClose}>{t("app_my_account")}</MenuItem>
      <Divider />
      <MenuItem onClick={() => void auth.removeUser()}>
        {t("app_logout")}
      </MenuItem>
    </Menu>
  );

  return (
    <>
      <AppBar
        position="fixed"
        sx={{ zIndex: (theme) => theme.zIndex.drawer + 1 }}
      >
        <Toolbar>
          <Typography variant="h6" noWrap component="div">
            Ledger
          </Typography>
          <Box sx={{ flexGrow: 1 }} />
          <Box>
            {/* <IconButton
              size="large"
              aria-label="show 4 new mails"
              color="inherit"
            >
              <Badge badgeContent={4} color="error">
                <MailIcon />
              </Badge>
            </IconButton>
            <IconButton
              size="large"
              aria-label="show 17 new notifications"
              color="inherit"
            >
              <Badge badgeContent={17} color="error">
                <NotificationsIcon />
              </Badge>
            </IconButton> */}

            <Chip
              icon={<AccountCircle />}
              label={auth.user?.profile.preferred_username}
              variant="outlined"
              color="secondary"
              onClick={handleProfileMenuOpen}
              sx={{ ml: 2 }}
            />
          </Box>
        </Toolbar>
      </AppBar>
      {renderMenu}
    </>
  );
};

export default TopBar;
