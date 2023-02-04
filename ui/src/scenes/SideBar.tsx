import { Box } from "@mui/material";
import Drawer from "@mui/material/Drawer";
import Toolbar from "@mui/material/Toolbar";
import List from "@mui/material/List";
import Divider from "@mui/material/Divider";
import ListItem from "@mui/material/ListItem";
import ListItemButton from "@mui/material/ListItemButton";
import ListItemIcon from "@mui/material/ListItemIcon";
import { HouseOutlined } from "@mui/icons-material";
import ListItemText from "@mui/material/ListItemText";
import BarChartOutlinedIcon from "@mui/icons-material/BarChartOutlined";
import { Link } from "react-router-dom";
import { useLocation } from "react-router-dom";
import SettingsIcon from "@mui/icons-material/Settings";
import LibraryBooksIcon from "@mui/icons-material/LibraryBooks";
import FormatListBulletedOutlinedIcon from "@mui/icons-material/FormatListBulletedOutlined";
import { useTranslation } from "react-i18next";

const drawerWidth = 240;
const menu = {
  main: [
    {
      title: "menu_dashboard",
      icon: <HouseOutlined />,
      to: "/",
    },
    {
      title: "menu_chart",
      icon: <BarChartOutlinedIcon />,
      to: "/chart",
    },
  ],
  secondary: [
    {
      title: "menu_versions",
      icon: <LibraryBooksIcon />,
      to: "/versions",
    },
    {
      title: "menu_last",
      icon: <LibraryBooksIcon />,
      to: "/last",
    },
    {
      title: "menu_feed",
      icon: <FormatListBulletedOutlinedIcon />,
      to: "/feed",
    },
  ],
  bottom: [
    {
      title: "menu_settings",
      icon: <SettingsIcon />,
      to: "/settings",
    },
  ],
};

const SideBar = () => {
  const location = useLocation();
  const { t } = useTranslation();

  return (
    <Box>
      <Drawer
        variant="permanent"
        sx={{
          width: drawerWidth,
          flexShrink: 0,
          [`& .MuiDrawer-paper`]: {
            width: drawerWidth,
            boxSizing: "border-box",
            backgroundColor: "#1F2A40",
          },
          [`& .MuiListItemIcon-root, & .MuiListItemText-root`]: {
            color: "#fff !important",
          },
          [`& .MuiListItemButton-root:hover .MuiListItemIcon-root, & .MuiListItemButton-root:hover .MuiListItemText-root`]:
            { color: "#868dfb !important" },
          [`& .Mui-selected .MuiListItemIcon-root, & .Mui-selected .MuiListItemText-root`]:
            { color: "#6870fa !important" },
        }}
      >
        <Toolbar />
        <Box sx={{ overflow: "auto" }}>
          <List>
            {menu.main.map((item, index) => (
              <ListItem key={index} disablePadding>
                <ListItemButton
                  component={Link}
                  to={item.to}
                  style={{ textDecoration: "none" }}
                  selected={item.to === location.pathname}
                >
                  <ListItemIcon>{item.icon}</ListItemIcon>
                  <ListItemText primary={t(item.title)} />
                </ListItemButton>
              </ListItem>
            ))}
          </List>

          <List>
            {menu.secondary.map((item, index) => (
              <ListItem key={index} disablePadding>
                <ListItemButton
                  component={Link}
                  to={item.to}
                  style={{ textDecoration: "none" }}
                  selected={item.to === location.pathname}
                >
                  <ListItemIcon>{item.icon}</ListItemIcon>
                  <ListItemText primary={t(item.title)} />
                </ListItemButton>
              </ListItem>
            ))}
          </List>
          <Box style={{ position: "absolute", bottom: "0", width: "100%" }}>
            <Divider />
            <List>
              {menu.bottom.map((item, index) => (
                <ListItem key={index} disablePadding>
                  <ListItemButton
                    component={Link}
                    to={item.to}
                    style={{ textDecoration: "none" }}
                    selected={item.to === location.pathname}
                  >
                    <ListItemIcon>{item.icon}</ListItemIcon>
                    <ListItemText primary={t(item.title)} />
                  </ListItemButton>
                </ListItem>
              ))}
            </List>
          </Box>
        </Box>
      </Drawer>
    </Box>
  );
};

export default SideBar;
